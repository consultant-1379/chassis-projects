package override

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
	"gerrit.ericsson.se/udm/nrf_common/pkg/profileop"
	"gerrit.ericsson.se/udm/nrf_common/pkg/schema"
	jsonpatch "github.com/evanphx/json-patch"
)

// OverrideUpdatedProfile overrides the the updated nf profile requested by NF Service Consumer
func OverrideUpdatedProfile(oldNfProfile *nrfschema.TNFProfileDB, updatedNfProfile *nrfschema.TNFProfile) ([]nrfschema.OverrideInfo, []byte, bool, *problemdetails.ProblemDetails) {
	oldServiceIDOverrideIndexMapper := make(map[string]int, 0)
	newServiceIDOverrideIndexMapper := make(map[string]int, 0)

	for index, service := range oldNfProfile.Body.NfServices {
		serviceID := service.ServiceInstanceId
		oldServiceIDOverrideIndexMapper[serviceID] = index
	}
	log.Debugf("The origin nfProfile offer services : %v", oldServiceIDOverrideIndexMapper)

	for index, service := range updatedNfProfile.NfServices {
		serviceID := service.ServiceInstanceId
		newServiceIDOverrideIndexMapper[serviceID] = index
	}
	log.Debugf("The update nfProfile offer services : %v", newServiceIDOverrideIndexMapper)

	body, _ := json.Marshal(updatedNfProfile)

	var updateBody []byte
	var updateOverrideInfo []nrfschema.OverrideInfo
	var overridden bool
	var err error
	newOverrideInfo := profileop.RebuildNfServiceOverrideInfo(oldServiceIDOverrideIndexMapper, newServiceIDOverrideIndexMapper, oldNfProfile.OverrideInfo)
	if len(newOverrideInfo) != 0 {
		overridden = true
		updateBody, updateOverrideInfo, err = applyOverrideAttributes(body, newOverrideInfo)
		if err != nil {
			log.Debugf("Apply override attributes for nfInstance %s error, %s", updatedNfProfile.NfInstanceId, err.Error())
			problemDetails := &problemdetails.ProblemDetails{
				Title: "Apply override attributes error",
			}

			return nil, nil, false, problemDetails
		}

		log.Debugf("After apply the new overrideInfo : %v", updateOverrideInfo)

		problemDetails := schema.ValidateNfProfile(string(updateBody[:]))
		if problemDetails != nil {
			log.Warnf("After override, the updated nf profile of nfInstance %s is not valid", updatedNfProfile.NfInstanceId)
			problemDetails.Title = fmt.Sprintf("the updated nf profile is not valid")

			return nil, nil, false, problemDetails
		}
	} else {
		updateBody = body
	}

	return updateOverrideInfo, updateBody, overridden, nil
}

func applyOverrideAttributes(body []byte, overrideInfo []nrfschema.OverrideInfo) ([]byte, []nrfschema.OverrideInfo, error) {
	var overrideInfoNew []nrfschema.OverrideInfo
	for _, overrideItem := range overrideInfo {
		patchStr := getPatchContent(overrideItem)

		log.Debugf("Apply patch : %s", patchStr)

		problemDetails := schema.ValidatePatchDocument(patchStr)
		if problemDetails != nil {
			err := fmt.Errorf("PatchDocument is not valid")
			return nil, nil, err
		}

		p, err := jsonpatch.DecodePatch([]byte(patchStr))
		if err == nil {
			var bodyTemp []byte
			bodyTemp, err = p.Apply(body)
			if err != nil {
				log.Warnf("Apply patch[%s] failure, will remove the override attribute. error:%s", patchStr, err.Error())
				bodyTemp = body
				continue
			} else {
				log.Debugf("Apply patch[%s] success\n", patchStr)
				body = bodyTemp
				overrideInfoNew = append(overrideInfoNew, overrideItem)
			}
		} else {
			log.Warnf("Decode patch[%s] failure, please check the patch message and Nrf-Provision will remove the override attribute, error:%s", patchStr, err.Error())
			continue
		}
	}

	return body, overrideInfoNew, nil
}

func getPatchContent(overrideItem nrfschema.OverrideInfo) string {
	ret1, _ := regexp.MatchString("nfStatus", overrideItem.Path)
	ret2, _ := regexp.MatchString("nfServiceStatus", overrideItem.Path)
	ret3, _ := regexp.MatchString("recoveryTime", overrideItem.Path)
	valueIsSting := ret1 || ret2 || ret3
	var patchContent string
	if overrideItem.Action == "replace" {
		if valueIsSting {
			patchContent = fmt.Sprintf("[{\"op\":\"%s\",\"path\":\"%s\",\"value\":\"%s\"}]", overrideItem.Action, overrideItem.Path, overrideItem.Value)
		} else {
			patchContent = fmt.Sprintf("[{\"op\":\"%s\",\"path\":\"%s\",\"value\":%s}]", overrideItem.Action, overrideItem.Path, overrideItem.Value)
		}
	} else if overrideItem.Action == "remove" {
		patchContent = fmt.Sprintf("[{\"op\":\"%s\",\"path\":\"%s\"}]", overrideItem.Action, overrideItem.Path)
	}
	return patchContent
}

// IsOverrideAttrExist is for check overrideAttrList exist in patch body
func IsOverrideAttrExist(patchData []nrfschema.TPatchItem) bool {
	var overrideExist bool = false
	for _, item := range patchData {
		if strings.Contains(item.Path, constvalue.OverrideAttrPath) {
			overrideExist = true
		}
	}
	return overrideExist
}

// IsOverrideAttribute checks whethe a PATCH path is override attribute
func IsOverrideAttribute(path string) bool {
	return strings.Contains(path, constvalue.OverrideAttrPath)
}

// PatchPathNeedOverride checks whether a PATCH path belongs to override list
func PatchPathNeedOverride(path string, overrideList []nrfschema.OverrideInfo) bool {
	for _, override := range overrideList {
		//path /allowedPlmns/0/mnc also need override if allowedPlmns is overrided
		ret, _ := regexp.MatchString("^"+override.Path, path)
		if ret {
			return true
		}
	}
	return false
}

// ConstructUpdatedOverrideInfo constructs overrideInfo for the nf profile updated by PATCH
func ConstructUpdatedOverrideInfo(oldNfProfileInDB *nrfschema.TNFProfileDB, updatedNfProfile *nrfschema.TNFProfile, patchData []nrfschema.TPatchItem, reqFromProv bool) (string, []string, error) {
	if oldNfProfileInDB.Provisioned == constvalue.Cmode_NFRegistered {
		oldServiceIDOverrideIndexMapper := make(map[string]int, 0)
		newServiceIDOverrideIndexMapper := make(map[string]int, 0)

		for index, service := range oldNfProfileInDB.Body.NfServices {
			serviceID := service.ServiceInstanceId
			oldServiceIDOverrideIndexMapper[serviceID] = index
		}
		log.Debugf("The origin nfProfile offer services : %v", oldServiceIDOverrideIndexMapper)

		for index, service := range updatedNfProfile.NfServices {
			serviceID := service.ServiceInstanceId
			newServiceIDOverrideIndexMapper[serviceID] = index
		}
		log.Debugf("The update nfProfile offer services : %v", newServiceIDOverrideIndexMapper)

		newOverrideInfo := profileop.RebuildNfServiceOverrideInfo(oldServiceIDOverrideIndexMapper, newServiceIDOverrideIndexMapper, oldNfProfileInDB.OverrideInfo)

		return constructOverrideInfo(reqFromProv, patchData, newOverrideInfo)
	}

	return "", make([]string, 0), nil
}

//constructOverrideInfo is for construct overrideInfo input for GRPC
func constructOverrideInfo(reqFromProv bool, patchData []nrfschema.TPatchItem, overrideList []nrfschema.OverrideInfo) (string, []string, error) {
	//update oldOverrideInfo
	var overrideAttrList = make([]string, 0)
	var newOverridePathList []string
	var err error
	if reqFromProv {
		newOverridePathList, err = executeOverridePath(patchData, &overrideList)
		if err != nil {
			return "", overrideAttrList, err
		}

		err = createOverrideAttrList(patchData, &overrideList, newOverridePathList)
		if err != nil {
			return "", overrideAttrList, err
		}
	}

	if len(overrideList) == 0 {
		log.Debugf("constructOverrideInfo overrideList is empty.")
		return "", overrideAttrList, nil
	}
	overrideData, err := json.Marshal(overrideList)
	if err != nil {
		err = fmt.Errorf("OverrideList marshal failure, err:%s", err.Error())
		return "", overrideAttrList, err
	}
	for _, item := range overrideList {
		overrideAttrList = append(overrideAttrList, item.Path)
	}
	log.Debugf("constructOverrideInfo overrideInfo: %s", string(overrideData))
	return string(overrideData), overrideAttrList, nil
}

func executeOverridePath(patchData []nrfschema.TPatchItem, overrideList *[]nrfschema.OverrideInfo) ([]string, error) {
	newOverridePathList := make([]string, 0)
	if overrideList == nil {
		err := fmt.Errorf("overrideList is nil")
		return newOverridePathList, err
	}
	for _, item := range patchData {
		if strings.Contains(item.Path, constvalue.OverrideAttrPath) {
			overridePath, err := changeOverrideList(item, overrideList)
			if err != nil {
				return newOverridePathList, err
			}
			if len(overridePath) != 0 {
				newOverridePathList = append(newOverridePathList, overridePath...)
			}
		}
	}
	return newOverridePathList, nil
}

func createOverrideAttrList(patchData []nrfschema.TPatchItem, overrideList *[]nrfschema.OverrideInfo, newOverridePathList []string) error {
	var err error
	for _, item := range patchData {
		if strings.Contains(item.Path, constvalue.OverrideAttrPath) == false {
			var actionStr, valueStr string
			if item.Op == "add" || item.Op == "replace" {
				actionStr = "replace"
			} else if item.Op == "remove" {
				actionStr = "remove"
			} else {
				continue
			}

			index := getIndexFromOverrideList(item.Path, *overrideList)
			if index >= 0 {
				if (item.Op == "add" || item.Op == "replace") && getValueStr(item.Value, &valueStr) == false {
					err = fmt.Errorf("Apply patch error, invaid value for path: %s", item.Path)
					return err
				}
				(*overrideList)[index].Action = actionStr
				(*overrideList)[index].Value = valueStr
			} else if index == -1 && pathExistInNewOverrideBody(item.Path, newOverridePathList) {
				if (item.Op == "add" || item.Op == "replace") && getValueStr(item.Value, &valueStr) == false {
					err = fmt.Errorf("Apply patch error, invaid value for path: %s", item.Path)
					return err
				}
				var override nrfschema.OverrideInfo
				override.Path = item.Path
				override.Action = actionStr
				override.Value = valueStr
				*overrideList = append(*overrideList, override)
			} else {
				//Do nothing
			}

		}
	}
	return nil
}

func getIndexFromOverrideList(path string, overrideList []nrfschema.OverrideInfo) int {
	for i, override := range overrideList {
		if path == override.Path {
			return i
		}
	}
	return -1
}

//pathExistInOverrideBody is to check if path exist in new patchBody overrideAttrList
func pathExistInNewOverrideBody(path string, newOverridePathList []string) bool {
	for _, newPath := range newOverridePathList {
		if path == newPath {
			return true
		}
	}
	return false
}

//changeOverrideList is for change override list and get override path
func changeOverrideList(item nrfschema.TPatchItem, overrideList *[]nrfschema.OverrideInfo) ([]string, error) {
	var overridePath = make([]string, 0)
	if item.Op == "add" && item.Path == constvalue.OverrideAttrPath+"/-" {
		overridePath = append(overridePath, item.Value.(string))
	} else if item.Op == "add" && item.Path == constvalue.OverrideAttrPath {
		for _, path := range (item.Value).([]interface{}) {
			overridePath = append(overridePath, path.(string))
		}
	} else if item.Op == "remove" && item.Path == constvalue.OverrideAttrPath {
		*overrideList = (*overrideList)[:0]
	} else if item.Op == "remove" {
		indexStr := strings.Trim(item.Path, constvalue.OverrideAttrPath+"/")
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 0 || index >= len(*overrideList) {
			err := fmt.Errorf("Apply patch error unable to access invalid index: %s", indexStr)
			return overridePath, err
		}
		(*overrideList) = append((*overrideList)[:index], (*overrideList)[index+1:]...)
	} else {
		err := fmt.Errorf("Invalid PATCH message, Op: %s, path: %s", item.Op, item.Path)
		return overridePath, err
	}
	return overridePath, nil
}

func getValueStr(value interface{}, outputStr *string) bool {
	switch value.(type) {
	case string:
		*outputStr = value.(string)
	case int:
		*outputStr = strconv.Itoa(value.(int))
	case float32:
		*outputStr = strconv.FormatFloat(value.(float64), 'f', -1, 32)
	case float64:
		*outputStr = strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case []interface{}:
		outputData, err := json.Marshal(value.([]interface{}))
		if err != nil {
			log.Errorf("getValueStr Marshal []interface{} failure, err: %s", err.Error())
			return false
		}
		*outputStr = string(outputData)
	case map[string]interface{}:
		outputData, err := json.Marshal(value.(map[string]interface{}))
		if err != nil {
			log.Errorf("getValueStr Marshal map[string]interface{} failure, err: %s", err.Error())
			return false
		}
		*outputStr = string(outputData)
	default:
		return false
	}
	return true
}

// UpdateProvOverrideInfo is for update TNFProfile.ProvisionInfo
func UpdateProvOverrideInfo(updatedNfProfile *nrfschema.TNFProfile, provisioned int32, overrideAttrList []string) bool {
	if updatedNfProfile == nil {
		log.Errorf("handleProvOverrideInfo ptr is nil.")
		return false
	}
	//update oldOverrideInfo
	provInfo := new(nrfschema.TProvInfo)
	if provisioned == constvalue.Cmode_NFRegistered {
		provInfo.CreateMode = constvalue.CMODE_NF_REGISTERED
		provInfo.OverrideAttrList = overrideAttrList

	} else if provisioned == constvalue.Cmode_Provisioned {
		provInfo.CreateMode = constvalue.CMODE_PROVISIONED
	} else {
		log.Errorf("updateProvOverrideInfo provisioned flag is wrong: %d", provisioned)
		return false
	}
	updatedNfProfile.ProvisionInfo = provInfo
	return true

}
