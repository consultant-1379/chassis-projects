package nfrequester

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
	"github.com/buger/jsonparser"
)

//SearchParameterMap is all parameters disc supported
var SearchParameterMap = map[string]bool{
	consts.SearchDataTargetNfType:        true,
	consts.SearchDataRequesterNfType:     true,
	consts.SearchDataServiceName:         true,
	consts.SearchDataRequesterNFInstFQDN: true,
	consts.SearchDataTargetPlmnList:      true,
	consts.SearchDataRequesterPlmnList:   true,
	consts.SearchDataTargetInstID:        true,
	consts.SearchDataTargetNFFQDN:        true,
	consts.SearchDataHnrfURI:             true,
	consts.SearchDataSnssais:             true,
	consts.SearchDataNsiList:             true,
	consts.SearchDataDnn:                 true,
	consts.SearchDataSmfServingArea:      true,
	consts.SearchDataTai:                 true,
	consts.SearchDataAmfRegionID:         true,
	consts.SearchDataAmfSetID:            true,
	consts.SearchDataGuami:               true,
	consts.SearchDataSupi:                true,
	consts.SearchDataUEIPv4Addr:          true,
	consts.SearchDataIPDoamin:            true,
	consts.SearchDataUEIPv6Prefix:        true,
	consts.SearchDataPGWInd:              true,
	consts.SearchDataPGW:                 true,
	consts.SearchDataGpsi:                true,
	consts.SearchDataExterGroupID:        true,
	consts.SearchDataDataSet:             true,
	consts.SearchDataRoutingIndic:        true,
	consts.SearchDataGroupIDList:         true,
	consts.SearchDataDnaiList:            true,
	consts.SearchDataUpfIwkEpsInd:        true,
	consts.SearchDataChfSupportedPlmn:    true,
	consts.SearchDataPreferredLocality:   true,
	consts.SearchDataAccessType:          true,
	consts.SearchDataSupportedFeatures:   true,
	consts.SearchDatacomplexQuery:        true,
}

//SearchParameter NRF Discovery Request Parameters
type SearchParameterData struct {
	value url.Values
	flag  map[string]bool
}

//GetValue to get all parameters
func (sp *SearchParameterData) GetValue() url.Values {
	return sp.value
}

//SetValue to set parameter value
func (sp *SearchParameterData) SetValue(key string, value []string) {
	sp.value[key] = value
}

//SetFlag to set parameters flag wheter it is exist
func (sp *SearchParameterData) SetFlag(key string, flag bool) {
	sp.flag[key] = flag
}

//GetExistFlag to get parameter flag wheter exist
func (sp *SearchParameterData) GetExistFlag(SearchData string) bool {
	return sp.flag[SearchData]
}

//InitMember to initial DiscGetPara
func (sp *SearchParameterData) InitMember(value url.Values) {
	sp.value = value
	sp.flag = make(map[string]bool)
}

func (sp *SearchParameterData) FetchRequesterNfTypeParameter() string {
	return sp.value[consts.SearchDataRequesterNfType][0]
}

func (sp *SearchParameterData) FetchTargetNfTypeParameter() string {
	return sp.value[consts.SearchDataTargetNfType][0]
}

func (sp *SearchParameterData) FetchServcieNamesParameter() []string {
	return sp.value[consts.SearchDataServiceName]
}

func (sp *SearchParameterData) FetchPlmnsParameter() []structs.PlmnID {
	return sp.fetchTargetPlmnListParameter()
}

func (sp *SearchParameterData) FetchRoamPlmnIDParameter() structs.PlmnID {
	plmns := sp.fetchTargetPlmnListParameter()
	if len(plmns) == 0 {
		return structs.PlmnID{}
	}

	return plmns[0]
}

//ValidateNRFDiscovery to validate request's parameter
func (sp *SearchParameterData) ValidateNRFDiscovery() *problemdetails.ProblemDetails {
	invalidParam, err := sp.validateIsSupportedParam()
	if err != nil {
		return sp.setProblemDetails(invalidParam, err)
	}

	err = sp.validateListStringType(consts.SearchDataServiceName)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataServiceName, err)
	}

	err = sp.validateBoolType(consts.SearchDataPGWInd)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataPGWInd, err)
	}

	err = sp.validateBoolType(consts.SearchDataUpfIwkEpsInd)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataUpfIwkEpsInd, err)
	}

	err = sp.validateListStringType(consts.SearchDataGroupIDList)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataGroupIDList, err)
	}

	err = sp.validateNFType(consts.SearchDataTargetNfType)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataTargetNfType, err)
	}

	err = sp.validateNFType(consts.SearchDataRequesterNfType)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataRequesterNfType, err)
	}

	err = sp.validateFQDNType(consts.SearchDataRequesterNFInstFQDN)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataRequesterNFInstFQDN, err)
	}

	err = sp.validateFQDNType(consts.SearchDataTargetNFFQDN)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataTargetNFFQDN, err)
	}

	err = sp.validateNFInstIDType(consts.SearchDataTargetInstID)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataTargetInstID, err)
	}

	err = sp.validatePlmnListType(consts.SearchDataRequesterPlmnList)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataRequesterPlmnList, err)
	}

	err = sp.validatePlmnListType(consts.SearchDataTargetPlmnList)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataTargetPlmnList, err)
	}

	err = sp.validatePlmnType(consts.SearchDataChfSupportedPlmn)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataChfSupportedPlmn, err)
	}

	err = sp.validateListSnssais(consts.SearchDataSnssais)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataSnssais, err)
	}

	err = sp.validateDnnType(consts.SearchDataDnn)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataDnn, err)
	}

	err = sp.validateSmfServingArea(consts.SearchDataSmfServingArea)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataSmfServingArea, err)
	}

	err = sp.validateSupiType(consts.SearchDataSupi)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataSupi, err)
	}

	err = sp.validateAMFRegionID()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataAmfRegionID, err)
	}

	err = sp.validateAMFSetID()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataAmfSetID, err)
	}

	err = sp.validateGuamiType(consts.SearchDataGuami)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataGuami, err)
	}

	err = sp.validateTaiType(consts.SearchDataTai)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataTai, err)
	}

	err = sp.validateIPV4Addr(consts.SearchDataUEIPv4Addr)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataUEIPv4Addr, err)
	}

	err = sp.validateIPDoamin(consts.SearchDataIPDoamin)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataIPDoamin, err)
	}

	err = sp.validateIPV6Prefix(consts.SearchDataUEIPv6Prefix)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataUEIPv6Prefix, err)
	}

	err = sp.validateFQDNType(consts.SearchDataPGW)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataPGW, err)
	}

	err = sp.validateGpsiType(consts.SearchDataGpsi)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataGpsi, err)
	}

	err = sp.validateGroupID()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataExterGroupID, err)
	}

	err = sp.validateDataSet()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataDataSet, err)
	}

	err = sp.validateListStringType(consts.SearchDataDnaiList)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataDnaiList, err)
	}

	err = sp.validateRoutingIndicator()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataRoutingIndic, err)
	}

	err = sp.validateHnrfURI()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataHnrfURI, err)
	}

	err = sp.validateListStringType(consts.SearchDataNsiList)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataNsiList, err)
	}

	err = sp.validateAccessType(consts.SearchDataAccessType)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataAccessType, err)
	}

	err = sp.validatePreferredLocality(consts.SearchDataPreferredLocality)
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataPreferredLocality, err)
	}

	err = sp.validateSupportedFeature()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataSupportedFeatures, err)
	}

	err = sp.validateComplexQuery()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDatacomplexQuery, err)
	}
	//after validate pgw and pgw-ind, then invoke this func
	err = sp.validateException()
	if err != nil {
		return sp.setProblemDetails(consts.SearchDataPGWInd, err)
	}

	return nil
}

//CacheSearchParameterInjection convert to cache searchParameter
func (sp *SearchParameterData) CacheSearchParameterInjection(cacheSearchParameter *cache.SearchParameter) *problemdetails.ProblemDetails {
	if sp.GetExistFlag(consts.SearchDataServiceName) {
		cacheSearchParameter.SetServiceNames(sp.fetchServiceNamesParameter())
	}
	if sp.GetExistFlag(consts.SearchDataTargetNfType) {
		cacheSearchParameter.SetTargetNfType(sp.fetchTargetNfTypeParameter())
	}
	if sp.GetExistFlag(consts.SearchDataRequesterNfType) {
		cacheSearchParameter.SetRequesterNfType(sp.fetchRequesterNfTypeParameter())
	}
	if sp.GetExistFlag(consts.SearchDataTargetPlmnList) {
		cacheSearchParameter.SetTargetPlmnList(sp.fetchTargetPlmnListParameter())
	}
	if sp.GetExistFlag(consts.SearchDataSnssais) {
		cacheSearchParameter.SetSnssai(sp.fetchSnssaisParameter())
	}
	if sp.GetExistFlag(consts.SearchDataSupportedFeatures) {
		cacheSearchParameter.SetSupportedFeatures(sp.fetchSupportedFeatureParameter())
	}
	if sp.GetExistFlag(consts.SearchDataChfSupportedPlmn) {
		cacheSearchParameter.SetChfSupportedPlmn(sp.fetchChfSupportedPlmnParameter())
	}
	if sp.GetExistFlag(consts.SearchDataDnn) {
		cacheSearchParameter.SetDnn(sp.fetchDnnParameter())
	}
	if sp.GetExistFlag(consts.SearchDataSupi) {
		cacheSearchParameter.SetSupi(sp.fetchSupiParameter())
	}
	if sp.GetExistFlag(consts.SearchDataRoutingIndic) {
		cacheSearchParameter.SetRoutingIndicator(sp.fetchRoutingIndicatorParameter())
	}
	if sp.GetExistFlag(consts.SearchDataSmfServingArea) {
		cacheSearchParameter.SetSmfServingArea(sp.fetchSmfServingAreaParameter())
	}
	if sp.GetExistFlag(consts.SearchDataTargetInstID) {
		cacheSearchParameter.SetTargetNfInstanceID(sp.fetchTargetNfIntanceIDParameter())
	}
	if sp.GetExistFlag(consts.SearchDataGpsi) {
		cacheSearchParameter.SetGpsi(sp.fetchGpsiParameter())
	}
	if sp.GetExistFlag(consts.SearchDataExterGroupID) {
		cacheSearchParameter.SetExternalGroupIdentity(sp.fetchExternalGroupIdentityParameter())
	}
	if sp.GetExistFlag(consts.SearchDataDataSet) {
		cacheSearchParameter.SetDataSet(sp.fetchdDataSetParameter())
	}
	if sp.GetExistFlag(consts.SearchDataAccessType) {
		cacheSearchParameter.SetAccessType(sp.fetchAccessTypeParameter())
	}
	if sp.GetExistFlag(consts.SearchDataPreferredLocality) {
		cacheSearchParameter.SetPreferredLocality(sp.fetchPreferredLocalityParameter())
	}
	if sp.GetExistFlag(consts.SearchDataNsiList) {
		cacheSearchParameter.SetNsiList(sp.fetchNsiListParameter())
	}
	if sp.GetExistFlag(consts.SearchDataGroupIDList) {
		cacheSearchParameter.SetGroupIDList(sp.fetchGroupIDListParameter())
	}
	if sp.GetExistFlag(consts.SearchDataIPDoamin) {
		cacheSearchParameter.SetIPDomain(sp.fetchIPDomainParameter())
	}
	if sp.GetExistFlag(consts.SearchDataDnaiList) {
		cacheSearchParameter.SetDnaiList(sp.fetchDnaiListParameter())
	}
	if sp.GetExistFlag(consts.SearchDataUpfIwkEpsInd) {
		cacheSearchParameter.SetUpfIwkEpsInd(sp.fetchUpfIwkEpsIndParameter())
	}

	return nil
}

///////////////////private////////////////////

func (sp *SearchParameterData) setExistFlag(SearchData string, IsExist bool) {
	sp.flag[SearchData] = IsExist
}

func (sp *SearchParameterData) validateStringTypeForMadatory(SearchData string) error {
	if len(sp.value[SearchData]) <= 0 {
		return fmt.Errorf(consts.MadatoryFieldNotExistFormat, SearchData, "SearchRequest")
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	if sp.value[SearchData][0] == "" {
		return fmt.Errorf(consts.FieldEmptyValue, SearchData)
	}
	sp.setExistFlag(SearchData, true)

	return nil
}

func (sp *SearchParameterData) validateStringTypeForOptional(SearchData string, Parttern bool) error {
	if len(sp.value[SearchData]) <= 0 {
		return nil
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	if sp.value[SearchData][0] == "" {
		return fmt.Errorf(consts.FieldEmptyValue, SearchData)
	}

	if Parttern {
		//matched, _ := regexp.MatchString(Parttern, sp.value[SearchData][0])
		matched := util.Compile[SearchData].MatchString(sp.value[SearchData][0])
		if !matched {
			return fmt.Errorf("Parameter %s value can't match parttern", SearchData)
		}
	}

	sp.setExistFlag(SearchData, true)

	return nil
}

func (sp *SearchParameterData) validateAccessType(SearchData string) error {
	if len(sp.value[SearchData]) <= 0 {
		return nil
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	if sp.value[SearchData][0] == "" {
		return fmt.Errorf(consts.FieldEmptyValue, SearchData)
	}

	if sp.value[SearchData][0] != consts.Access3GPP && sp.value[SearchData][0] != consts.NonAccess3GPP {
		return fmt.Errorf("invalid access-type %s, should be 3GPP_ACCESS or NON_3GPP_ACCESS", sp.value[SearchData][0])
	}
	sp.setExistFlag(SearchData, true)

	return nil
}

func (sp *SearchParameterData) validateSupportedFeature() error {
	return sp.validateStringTypeForOptional(consts.SearchDataSupportedFeatures, true)
}

func (sp *SearchParameterData) validateHnrfURI() error {
	//TODO URI regex match
	return sp.validateStringTypeForOptional(consts.SearchDataHnrfURI, true)
}

func (sp *SearchParameterData) validateGroupID() error {
	return sp.validateStringTypeForOptional(consts.SearchDataExterGroupID, true)
}

func (sp *SearchParameterData) validateRoutingIndicator() error {
	return sp.validateStringTypeForOptional(consts.SearchDataRoutingIndic, true)
}

func (sp *SearchParameterData) validateAMFRegionID() error {
	return sp.validateStringTypeForOptional(consts.SearchDataAmfRegionID, false)
}

func (sp *SearchParameterData) validateAMFSetID() error {
	return sp.validateStringTypeForOptional(consts.SearchDataAmfSetID, false)
}

func (sp *SearchParameterData) validateDataSet() error {
	return sp.validateStringTypeForOptional(consts.SearchDataDataSet, true)
}

func (sp *SearchParameterData) validatePlmnType(SearchData string) error {
	if len(sp.value[SearchData]) <= 0 {
		return nil
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	plmn := sp.value[SearchData][0]

	if plmn == "" {
		return fmt.Errorf(consts.FieldEmptyValue, SearchData)
	}

	var mcc string
	var err error
	if SearchData == consts.SearchDataGuami || SearchData == consts.SearchDataTai {
		mcc, err = jsonparser.GetString([]byte(plmn), "plmnId", consts.SearchDataMcc)
	} else {
		mcc, err = jsonparser.GetString([]byte(plmn), consts.SearchDataMcc)
	}
	if err != nil {
		return fmt.Errorf(consts.MadatoryFieldNotExistFormat, "mcc", SearchData)
	}

	matched := util.Compile[consts.SearchDataMcc].MatchString(mcc)
	if !matched {
		return fmt.Errorf("invalid format for mcc in %s", SearchData)
	}
	var mnc string
	if SearchData == consts.SearchDataGuami || SearchData == consts.SearchDataTai {
		mnc, err = jsonparser.GetString([]byte(plmn), "plmnId", consts.SearchDataMnc)
	} else {
		mnc, err = jsonparser.GetString([]byte(plmn), consts.SearchDataMnc)
	}
	if err != nil {
		return fmt.Errorf(consts.MadatoryFieldNotExistFormat, "mnc", SearchData)
	}

	matched = util.Compile[consts.SearchDataMnc].MatchString(mnc)
	if !matched {
		return fmt.Errorf("invalid format for mnc in %s", SearchData)
	}
	sp.setExistFlag(SearchData, true)

	return nil
}

//isJSON is to check if the string is a valid json
func isJSON(s string) bool {
	var js interface{}

	return json.Unmarshal([]byte(s), &js) == nil
}

func (sp *SearchParameterData) validatePlmnListType(SearchData string) error {
	if _, ok := sp.value[SearchData]; ok {
		if len(sp.value[SearchData]) <= 0 {
			return fmt.Errorf(consts.FieldEmptyValue, SearchData)
		}
		for _, plmn := range sp.value[SearchData] {
			if plmn != "" {
				ok := true
				errorInfo := ""
				if !isJSON(plmn) {
					return fmt.Errorf("parse %v json wrong", SearchData)
				}
				if !strings.Contains(plmn, "[") && !strings.Contains(plmn, "]") {
					plmn = "[" + plmn + "]"
				}
				_, err := jsonparser.ArrayEach([]byte(plmn), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
					if !ok {
						return
					}
					mcc, err := jsonparser.GetString([]byte(value), consts.SearchDataMcc)
					if err != nil {
						ok = false
						errorInfo = fmt.Sprintf(consts.MadatoryFieldNotExistFormat, "mcc", SearchData)
					}
					matched := util.Compile[consts.SearchDataMcc].MatchString(mcc)
					if !matched {
						ok = false
						errorInfo = fmt.Sprintf("invalid format for mcc in %s", SearchData)
					}
					mnc, err2 := jsonparser.GetString([]byte(value), consts.SearchDataMnc)
					if err2 != nil {
						ok = false
						errorInfo = fmt.Sprintf(consts.MadatoryFieldNotExistFormat, "mnc", SearchData)
					}
					matched = util.Compile[consts.SearchDataMnc].MatchString(mnc)
					if !matched {
						ok = false
						errorInfo = fmt.Sprintf("invalid format for mnc in %s", SearchData)
					}
				})
				if err != nil {
					return fmt.Errorf("parse %v array wrong: %v", SearchData, err)
				}

				if !ok {
					return fmt.Errorf("%s", errorInfo)
				}
			} else {
				return fmt.Errorf(consts.FieldEmptyValue, SearchData)
			}
		}
	} else {
		return nil
	}

	sp.setExistFlag(SearchData, true)

	return nil
}

func (sp *SearchParameterData) validateGuamiType(SearchData string) error {
	err := sp.validatePlmnType(SearchData)
	if err == nil && sp.GetExistFlag(SearchData) {
		amfID, err := jsonparser.GetString([]byte(sp.value[SearchData][0]), consts.SearchDataAmfID)
		if err != nil {
			sp.setExistFlag(SearchData, false)
			return fmt.Errorf(consts.MadatoryFieldNotExistFormat, "amfid", "Guami")
		}
		matched := util.Compile[consts.SearchDataAmfID].MatchString(amfID)
		if !matched {
			sp.setExistFlag(SearchData, false)
			return fmt.Errorf("Invalid format for amfid in Guami")
		}

		return nil
	}

	return err
}

func (sp *SearchParameterData) validateTaiType(SearchData string) error {
	err := sp.validatePlmnType(SearchData)
	if err == nil && sp.GetExistFlag(SearchData) {
		tac, err := jsonparser.GetString([]byte(sp.value[SearchData][0]), consts.SearchDataTac)
		if err != nil {
			sp.setExistFlag(SearchData, false)
			return fmt.Errorf(consts.MadatoryFieldNotExistFormat, "Tac", "Tai")
		}

		matched := util.Compile[consts.SearchDataTac].MatchString(tac)
		if !matched {
			sp.setExistFlag(SearchData, false)
			return fmt.Errorf("invalid format for amfid in Gumai")
		}

		return nil
	}

	return err
}

func (sp *SearchParameterData) validateIPV4Addr(SearchData string) error {
	//TODO ipv4 addr
	return sp.validateStringTypeForOptional(SearchData, true)
}

func (sp *SearchParameterData) validateIPV6Prefix(SearchData string) error {
	err := sp.validateStringTypeForOptional(SearchData, false)
	if err != nil || sp.GetExistFlag(SearchData) == false {
		return err
	}
	_, _, err = net.ParseCIDR(sp.value[SearchData][0])
	if err != nil {
		sp.setExistFlag(SearchData, false)
		return err
	}

	return nil
}

func (sp *SearchParameterData) validateNFType(SearchData string) error {
	err := sp.validateStringTypeForMadatory(SearchData)
	if err != nil {
		return err
	}

	/*
		if !checkNfTypes(sp.value[SearchData][0]) {
			f.setExistFlag(SearchData, false)
			return fmt.Errorf("%s %s is invaild.", SearchData, f.value[SearchData][0])
		}
	*/

	sp.setExistFlag(SearchData, true)

	return nil

}

func (sp *SearchParameterData) validateGpsiType(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, true)
}

func (sp *SearchParameterData) validateListSnssais(SearchData string) error {
	if len(sp.value[SearchData]) > 0 {
		for _, item := range sp.value[SearchData] {
			if item != "" {
				ok := true
				errorInfo := ""
				if !isJSON(item) {
					return fmt.Errorf("parse %v json wrong", SearchData)
				}
				if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
					item = "[" + item + "]"
				}
				_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
					if !ok {
						return
					}
					sstID, err := jsonparser.GetInt(value, consts.SearchDataSnssaiSst)
					if err != nil {
						ok = false
						errorInfo = fmt.Sprintf("Get sst from snssais wrong: %v", err)
						return
					}
					if sstID < 0 || sstID > 255 {
						ok = false
						errorInfo = fmt.Sprintf("sst out of range")
						return
					}
					//sd is option, but if exist, need meet the pattern
					sd, err := jsonparser.GetString(value, consts.SearchDataSnssaiSd)
					if err == nil {
						matched := util.Compile[consts.SearchDataSnssaiSd].MatchString(sd)
						if !matched {
							ok = false
							errorInfo = fmt.Sprintf("sd exist, but format not right")
							return
						}
					}
				})
				if err != nil {
					return fmt.Errorf("parse snssais array wrong: %v", err)
				}

				if !ok {
					return fmt.Errorf("%s", errorInfo)
				}
			} else {
				return fmt.Errorf(consts.FieldEmptyValue, SearchData)
			}
		}
		sp.setExistFlag(SearchData, true)
	}

	return nil
}

func (sp *SearchParameterData) validateDnnType(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, false)
}

func (sp *SearchParameterData) validateNFInstIDType(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, false)
}

func (sp *SearchParameterData) validateSmfServingArea(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, false)
}

func (sp *SearchParameterData) validateIPDoamin(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, false)
}

func (sp *SearchParameterData) validatePreferredLocality(SearchData string) error {
	return sp.validateStringTypeForOptional(SearchData, false)
}

func (sp *SearchParameterData) validateListString(SearchData string) error {
	if _, exist := sp.value[SearchData]; exist {
		if len(sp.value[SearchData]) <= 0 {
			return fmt.Errorf(consts.ArrayFileldExistEmptyValue, SearchData)
		}

		for _, item := range sp.value[SearchData] {
			if item == "" {
				return fmt.Errorf(consts.FieldEmptyValue, SearchData)
			}
		}

		sp.setExistFlag(SearchData, true)
	}

	return nil
}

func (sp *SearchParameterData) validateListStringType(SearchData string) error {
	switch SearchData {
	case consts.SearchDataServiceName:
		return sp.validateListString(consts.SearchDataServiceName)
	case consts.SearchDataNsiList:
		return sp.validateListString(consts.SearchDataNsiList)
	case consts.SearchDataGroupIDList:
		return sp.validateListString(consts.SearchDataGroupIDList)
	case consts.SearchDataDnaiList:
		return sp.validateListString(consts.SearchDataDnaiList)
	default:
		return fmt.Errorf("Invalid ListString Parameter %s", SearchData)
	}
}

func (sp *SearchParameterData) validateSupiType(SearchData string) error {
	if len(sp.value[SearchData]) <= 0 {
		return nil
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	//supiRegex := "^imsi-[0-9]{5,15}$|nai-.+$"
	supi := sp.value[SearchData][0]

	matched := util.Compile[consts.SearchDataSupi].MatchString(supi)
	if !matched {
		return fmt.Errorf("The %s is %s, doesn't match the regex", SearchData, supi)
	}
	sp.setExistFlag(SearchData, true)

	return nil
}

func (sp *SearchParameterData) validateFQDNType(SearchData string) error {
	switch SearchData {
	case consts.SearchDataRequesterNFInstFQDN:
		return sp.validateStringTypeForOptional(SearchData, false)
	case consts.SearchDataPGW:
		return sp.validateStringTypeForOptional(SearchData, false)
	case consts.SearchDataTargetNFFQDN:
		return sp.validateStringTypeForOptional(SearchData, false)
	default:
		return fmt.Errorf("Invalid FQDN type Parameter %s ", SearchData)
	}
}

func (sp *SearchParameterData) validateIfNoneMatch() error {
	return sp.validateStringTypeForOptional(consts.SearchDataIfNoneMatch, false)
}

func (sp *SearchParameterData) validateBoolType(SearchData string) error {
	if len(sp.value[SearchData]) <= 0 {
		return nil
	}

	if len(sp.value[SearchData]) > 1 {
		return fmt.Errorf(consts.FieldMultipleValue, SearchData)
	}

	if sp.value[SearchData][0] != consts.BoolTrueString && sp.value[SearchData][0] != consts.BoolFalseString {
		return fmt.Errorf("%s value only true or false", SearchData)
	}

	sp.setExistFlag(SearchData, true)

	return nil
}

//this func should be invoked at last on func ValidateNRFDiscovery
func (sp *SearchParameterData) validateException() error {
	//if pgw-ind value is false and pgw exist, should return 400 Bad Request
	if sp.GetExistFlag(consts.SearchDataPGW) {
		value, err := sp.GetNRFDiscBoolValue(consts.SearchDataPGWInd)
		if err == nil && false == value {
			return fmt.Errorf("pgw-ind value is false and pgw exist, should return 400 Bad Request")
		}
	}

	return nil
}

//GetNRFDiscBoolValue to get bool type parameter value
func (sp *SearchParameterData) GetNRFDiscBoolValue(SearchData string) (bool, error) {
	if !sp.GetExistFlag(SearchData) {
		return false, fmt.Errorf("Field not exist")
	}

	if consts.BoolTrueString == sp.value[SearchData][0] {
		return true, nil
	}

	return false, nil
}

func (sp *SearchParameterData) getNRFDiscStringValue(SearchData string) string {
	if !sp.GetExistFlag(SearchData) {
		return ""
	}
	return sp.value[SearchData][0]
}

func (sp *SearchParameterData) getNRFDiscListString(SearchData string) []string {
	if !sp.GetExistFlag(SearchData) {
		return nil
	}

	var retlist []string
	for _, item := range sp.value[SearchData] {
		str := strings.Replace(item, " ", "", -1)
		array := strings.Split(str, ",")
		for _, item2 := range array {
			retlist = append(retlist, item2)
		}
	}

	return retlist
}

//GetNRFDiscPlmnValue to get plmn value
func (sp *SearchParameterData) GetNRFDiscPlmnValue(SearchData string) (string, string) {
	if !sp.GetExistFlag(SearchData) {
		return "", ""
	}
	plmn := sp.value[SearchData][0]
	var mcc string
	var mnc string

	if SearchData == consts.SearchDataGuami || SearchData == consts.SearchDataTai {
		mcc, _ = jsonparser.GetString([]byte(plmn), "plmnId", consts.SearchDataMcc)
		mnc, _ = jsonparser.GetString([]byte(plmn), "plmnId", consts.SearchDataMnc)
		log.Debugf("PLMN: %s %s %s %s %s", plmn, SearchData, consts.SearchDataTai, mcc, mnc)
	} else {
		mcc, _ = jsonparser.GetString([]byte(plmn), consts.SearchDataMcc)
		mnc, _ = jsonparser.GetString([]byte(plmn), consts.SearchDataMnc)
	}
	//if len(mnc) == 2 {
	//	mnc = "0" + mnc
	//}
	switch SearchData {
	//case consts.SearchDataRequesterPlmn:
	//	return mcc, mnc
	//case consts.SearchDataTargetPlmn:
	//case consts.SearchDataTai:
	case consts.SearchDataTai, consts.SearchDataGuami, consts.SearchDataChfSupportedPlmn:
		return mcc + mnc, ""
	default:
		return "", ""
	}
}

//GetNRFDiscPlmnListValue to get plmnList value
func (sp *SearchParameterData) GetNRFDiscPlmnListValue(SearchData string) []string {
	var plmnList []string
	if !sp.GetExistFlag(SearchData) {
		return plmnList
	}
	for _, plmn := range sp.value[SearchData] {
		plmnIds, valueType, _, err := jsonparser.Get([]byte(plmn))
		if err == nil {
			if valueType == jsonparser.Array {
				_, _ = jsonparser.ArrayEach(plmnIds, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					mcc, _ := jsonparser.GetString([]byte(value), consts.SearchDataMcc)
					mnc, _ := jsonparser.GetString([]byte(value), consts.SearchDataMnc)
					plmnList = append(plmnList, mcc+mnc)
				})
			} else if valueType == jsonparser.Object {
				mcc, _ := jsonparser.GetString([]byte(plmnIds), consts.SearchDataMcc)
				mnc, _ := jsonparser.GetString([]byte(plmnIds), consts.SearchDataMnc)
				plmnList = append(plmnList, mcc+mnc)
			}
		}
	}
	return plmnList
}

//GetNRFDiscListSnssais is to get snssias from searchdata
func (sp *SearchParameterData) GetNRFDiscListSnssais(SearchData string) string {
	var snssais string
	if sp.GetExistFlag(SearchData) && len(sp.value[SearchData]) > 0 {
		for _, item := range sp.value[SearchData] {
			if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
				item = "[" + item + "]"
			}
			_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
				sstID, _ := jsonparser.GetInt(value, consts.SearchDataSnssaiSst)
				sdID, _ := jsonparser.GetString(value, consts.SearchDataSnssaiSd)
				if snssais != "" {
					snssais = fmt.Sprintf("%s,{\"sst\": %d,\"sd\": \"%s\"}", snssais, sstID, sdID)
				} else {
					snssais = fmt.Sprintf("{\"sst\": %d,\"sd\": \"%s\"}", sstID, sdID)
				}

			})

			if err != nil {
				return ""
			}
		}
		snssais = "[" + snssais + "]"
	}

	return snssais
}

func (sp *SearchParameterData) validateComplexQuery() error {
	if _, exist := sp.value[consts.SearchDatacomplexQuery]; exist {
		return fmt.Errorf("Not Support Parameter complexQuery")
	}

	return nil
}

//validateIsSupportedParam is to refuse the abnormal params
func (sp *SearchParameterData) validateIsSupportedParam() (string, error) {
	isSupported := true
	var invalidParams string
	for key := range sp.value {
		if SearchParameterMap[key] != true {
			if invalidParams == "" {
				invalidParams = key
			} else {
				invalidParams += "," + key
			}
			isSupported = false
		}
	}
	if !isSupported {
		return invalidParams, fmt.Errorf("Not Support Parameter")
	}

	return invalidParams, nil
}

func (sp *SearchParameterData) fetchServiceNamesParameter() []string {
	serviceNames := make([]string, 0)
	serviceNamesOrigValue := sp.value[consts.SearchDataServiceName]
	if len(serviceNamesOrigValue) > 0 {
		for _, v := range serviceNamesOrigValue {
			if strings.Contains(v, ",") {
				serviceNamesArray := strings.Split(v, ",")
				serviceNames = append(serviceNames, serviceNamesArray...)
			} else {
				serviceNames = append(serviceNames, v)
			}
		}
	}

	return serviceNames
}

func (sp *SearchParameterData) fetchTargetNfTypeParameter() string {
	return sp.value[consts.SearchDataTargetNfType][0]
}

func (sp *SearchParameterData) fetchRequesterNfTypeParameter() string {
	return sp.value[consts.SearchDataRequesterNfType][0]
}

func (sp *SearchParameterData) fetchTargetPlmnListParameter() []structs.PlmnID {
	plmnList := make([]structs.PlmnID, 0)
	targetPlmnList := sp.value[consts.SearchDataTargetPlmnList]
	if len(targetPlmnList) > 0 {
		for _, targetPlmn := range targetPlmnList {
			if targetPlmn != "" {
				if !strings.Contains(targetPlmn, "[") && !strings.Contains(targetPlmn, "]") {
					targetPlmn = "[" + targetPlmn + "]"
				}

				var plmnIDs []structs.PlmnID
				err := json.Unmarshal([]byte(targetPlmn), &plmnIDs)
				if err != nil {
					log.Errorf("Unmarshal %s in %s failed", targetPlmn, consts.SearchDataTargetPlmnList)
					continue
				}

				plmnList = append(plmnList, plmnIDs...)
			} else {
				log.Error("One targetPlmn in targetPlmnList is empty")
				continue
			}
		}
	}

	return plmnList
}

func (sp *SearchParameterData) fetchSnssaisParameter() []cache.SNssai {
	snssais := make([]cache.SNssai, 0)
	snssaiArray := sp.value[consts.SearchDataSnssais]
	if len(snssaiArray) > 0 {
		for _, snssaiItem := range snssaiArray {
			if snssaiItem != "" {
				if !strings.Contains(snssaiItem, "[") && !strings.Contains(snssaiItem, "]") {
					snssaiItem = "[" + snssaiItem + "]"
				}

				var snssai []cache.SNssai
				err := json.Unmarshal([]byte(snssaiItem), &snssai)
				if err != nil {
					log.Errorf("Unmarshal %s in %s failed", snssaiItem, consts.SearchDataSnssais)
					continue
				}

				snssais = append(snssais, snssai...)
			} else {
				log.Error("One snssais in snssaises is empty")
				continue
			}
		}
	}

	return snssais
}

func (sp *SearchParameterData) fetchSupportedFeatureParameter() string {
	return sp.value[consts.SearchDataSupportedFeatures][0]
}

func (sp *SearchParameterData) fetchChfSupportedPlmnParameter() structs.PlmnID {
	chfSupportedPlmn := sp.value[consts.SearchDataChfSupportedPlmn][0]
	if chfSupportedPlmn == "" {
		return structs.PlmnID{}
	}

	var plmn structs.PlmnID
	err := json.Unmarshal([]byte(chfSupportedPlmn), &plmn)
	if err != nil {
		log.Errorf("Unmarshal %s in %s failed", chfSupportedPlmn, consts.SearchDataChfSupportedPlmn)
		return structs.PlmnID{}
	}

	return plmn
}

func (sp *SearchParameterData) fetchDnnParameter() string {
	return sp.value[consts.SearchDataDnn][0]
}

func (sp *SearchParameterData) fetchSupiParameter() string {
	return sp.value[consts.SearchDataSupi][0]
}

func (sp *SearchParameterData) fetchRoutingIndicatorParameter() string {
	return sp.value[consts.SearchDataRoutingIndic][0]
}

func (sp *SearchParameterData) fetchSmfServingAreaParameter() string {
	return sp.value[consts.SearchDataSmfServingArea][0]
}

func (sp *SearchParameterData) fetchTargetNfIntanceIDParameter() string {
	return sp.value[consts.SearchDataTargetInstID][0]
}

func (sp *SearchParameterData) fetchGpsiParameter() string {
	return sp.value[consts.SearchDataGpsi][0]
}

func (sp *SearchParameterData) fetchExternalGroupIdentityParameter() string {
	return sp.value[consts.SearchDataExterGroupID][0]
}

func (sp *SearchParameterData) fetchdDataSetParameter() string {
	return sp.value[consts.SearchDataDataSet][0]
}

func (sp *SearchParameterData) fetchAccessTypeParameter() string {
	return sp.value[consts.SearchDataAccessType][0]
}

func (sp *SearchParameterData) fetchPreferredLocalityParameter() string {
	return sp.value[consts.SearchDataPreferredLocality][0]
}

func (sp *SearchParameterData) fetchNsiListParameter() []string {
	nsiList := make([]string, 0)
	nsiListOriginValue := sp.value[consts.SearchDataNsiList]
	if len(nsiListOriginValue) > 0 {
		for _, v := range nsiListOriginValue {
			if strings.Contains(v, ",") {
				nsiListData := strings.Split(v, ",")
				nsiList = append(nsiList, nsiListData...)
			} else {
				nsiList = append(nsiList, v)
			}
		}

	}

	return nsiList
}

func (sp *SearchParameterData) fetchGroupIDListParameter() []string {
	groupIDList := make([]string, 0)
	groupIDListOriginValue := sp.value[consts.SearchDataGroupIDList]
	if len(groupIDListOriginValue) > 0 {
		for _, v := range groupIDListOriginValue {
			if strings.Contains(v, ",") {
				groupIData := strings.Split(v, ",")
				groupIDList = append(groupIDList, groupIData...)
			} else {
				groupIDList = append(groupIDList, v)
			}
		}
	}

	return groupIDList
}

func (sp *SearchParameterData) fetchIPDomainParameter() string {
	return sp.value[consts.SearchDataIPDoamin][0]
}

func (sp *SearchParameterData) fetchDnaiListParameter() []string {
	dnaiList := make([]string, 0)
	dnaiListOriginValue := sp.value[consts.SearchDataDnaiList]
	if len(dnaiListOriginValue) > 0 {
		for _, v := range dnaiListOriginValue {
			if strings.Contains(v, ",") {
				dnailItems := strings.Split(v, ",")
				dnaiList = append(dnaiList, dnailItems...)
			} else {
				dnaiList = append(dnaiList, v)
			}
		}
	}

	return dnaiList
}

func (sp *SearchParameterData) fetchUpfIwkEpsIndParameter() string {
	return sp.value[consts.SearchDataUpfIwkEpsInd][0]
}

////////////problemDetails////////////

func (sp *SearchParameterData) setProblemDetails(SearchData string, err error) *problemdetails.ProblemDetails {
	var problemDetail *problemdetails.ProblemDetails
	invalidpara := &problemdetails.InvalidParam{}
	problemDetail = &problemdetails.ProblemDetails{}
	problemDetail.InvalidParams = append(problemDetail.InvalidParams, invalidpara)

	problemDetail.Title = err.Error()
	problemDetail.InvalidParams[0].Param = SearchData
	problemDetail.InvalidParams[0].Reason = err.Error()

	if SearchData == consts.SearchDatacomplexQuery {
		problemDetail.Cause = consts.UnSupportedQueryParameter
	}

	return problemDetail
}
