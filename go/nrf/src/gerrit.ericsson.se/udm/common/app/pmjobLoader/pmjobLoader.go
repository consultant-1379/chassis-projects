package pmjobLoader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/buger/jsonparser"
)

var (
	pmBrSkeletonConf   string
	pmBrServiceConf    string
	pmBrSchemaName     string
	pmBulkReporterName string

	cmmService   string
	cmHTTPClient *httpclient.HttpClient
)

const (
	// Default eric-cm-mediator URI
	defaultCmMediatorURI      = "http://eric-cm-mediator:5003/cm/api/v1.1"
	defaultPmBulkReporterName = "adp-gs-pm-br"
	defaultPmBrSchemaName     = "adp-gs-pm-br"

	pmBrSchemaPollingInterval = 15

	// LoadPm indications pm job action: load pm
	LoadPm = "load"
	// UnloadPm indications pm job action: unload pm
	UnLoadPm = "unload"
)

type PatchData struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func Init() {

	var (
		defaultConnections = 2
		defaultTimeout     = 10 * time.Second
		defaultKeepAlive   = true
	)
	cmHTTPClient = httpclient.InitHttpClient(
		httpclient.Connections(defaultConnections),
		httpclient.Timeout(defaultTimeout),
		httpclient.KeepAlive(defaultKeepAlive),
	)
	if cmHTTPClient == nil {
		panic("failed to initialize HTTP Client for cmproxy")
	}

	cmmService = os.Getenv("CMM_SERVICE")
	pmBrSchemaName = os.Getenv("PM_BR_SCHEMA_NAME")
	pmBulkReporterName = os.Getenv("PM_BR_CONF_NAME")
	pmBrSkeletonConf = os.Getenv("PM_BR_SKELETON_FILE")
	pmBrServiceConf = os.Getenv("PM_BR_SERVICE_FILE")

}

func LoadPmJob() {

	for true {
		if checkPmBrSchemaExists() {
			break
		}
		log.Debugf("Wait for Pm Bulk Reporter Schema ready, sleep %d seconds", pmBrSchemaPollingInterval)
		time.Sleep(pmBrSchemaPollingInterval * time.Second)
	}

	if !checkPmBrConfExists() {
		log.Debug("Pm Bulk Reporter Conf not exsits in CM mediator")
		if !postPmBrSkeletonConf() {
			log.Error("Failed to load pm br sckeleton conf ")
		}
	}
	patchPmBrConf(LoadPm)

}

func UnLoadPmJob() {

	patchPmBrConf(UnLoadPm)

}

func checkPmBrSchemaExists() bool {
	pmBrSchemaUrl := getPmBrSchemaUrl()

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	resp, err := cmHTTPClient.HttpDo("GET", pmBrSchemaUrl, header, strings.NewReader(""))

	if err != nil {
		log.Errorf("Failed to check Pm Bulk Reporter schema from CM, error is  %v", err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		log.Debugf("Pm Bulk Reporter schema exists")
		return true
	}
	log.Debugf("Pm Bulk Reporter schema Not exists")
	return false
}

func checkPmBrConfExists() bool {
	pmBrServiceConfUrl := getPmBrServiceConfUrl()

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	resp, err := cmHTTPClient.HttpDo("GET", pmBrServiceConfUrl, header, strings.NewReader(""))
	if err != nil {
		log.Errorf("Failed to check Pm Bulk Reporter Config from CM, error is  %v", err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		log.Debugf("Pm Bulk Reporter Config exists")
		return true
	}
	log.Debugf("Pm Bulk Reporter Config Not exists")
	return false
}

func postPmBrSkeletonConf() bool {

	pmBrSkeletonConfUrl := getPmBrSkeletonConfUrl()

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	fileInBytes, err := loadConf(pmBrSkeletonConf)
	if err != nil {
		return false
	}

	resp, err := cmHTTPClient.HttpDo("POST", pmBrSkeletonConfUrl, header, bytes.NewReader(fileInBytes))
	if err != nil {
		log.Errorf("Send to CM mediator failed, the error is %s", err.Error())
		return false
	}
	log.Debugf("Post Bulk report skeleton conf status code is %d", resp.StatusCode)
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
		return true
	}
	log.Warnf("Post Bulk report skeleton conf failed, the body in response is %s", string(resp.Body))
	return false
}

func patchPmBrConf(action string) {

	pmBrServiceConfUrl := getPmBrServiceConfUrl()

	//Get current pm BR conf from CM
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	resp, err := cmHTTPClient.HttpDo("GET", pmBrServiceConfUrl, header, strings.NewReader(""))
	if err != nil {
		log.Errorf("Fails to get pm br conf from CM, the error is %s", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Fails to get pm br conf from CM, the response from CM is StatusCode:%d, Body:%s", resp.StatusCode, string(resp.Body))
		return
	}

	//Load NF pm BR config from file
	pmBrConfFromfile, err := loadConf(pmBrServiceConf)
	if err != nil {
		return
	}

	var patchData []byte
	if LoadPm == action {
		patchData = constructLoadData(resp.Body, pmBrConfFromfile)
	} else if UnLoadPm == action {
		patchData = constructUnloadData(resp.Body, pmBrConfFromfile)
	} else {
		log.Errorf("parameter action is incorrect")
		return
	}

	if patchData == nil {
		log.Debug("PM job and group are already in CM mediator, no need to send patch to CM!")
		return
	}

	header["Content-Type"] = "application/json-patch+json"
	log.Debugf("Send Patch body %s to CM", string(patchData))
	resp, err = cmHTTPClient.HttpDo("PATCH", pmBrServiceConfUrl, header, bytes.NewReader(patchData))

	if err != nil {
		log.Errorf("Send Patch to CM mediator failed, the error is %s", err.Error())
		return
	}

	log.Debugf("The response from CM is StatusCode:%d, Body:%s", resp.StatusCode, string(resp.Body))
}

func loadConf(confName string) ([]byte, error) {
	log.Debugf("Start to load config file %s", confName)
	fileInBytes, err := ioutil.ReadFile(confName)
	if err != nil {
		log.Errorf("Failed to read config file:%s. %s", confName, err.Error())
		return nil, err
	}
	return fileInBytes, nil
}

func constructLoadData(pmBrInCM, pmBrInConf []byte) []byte {

	var patchDatas []PatchData

	keys := []string{"job", "group"}
	for _, key := range keys {
		_, _ = jsonparser.ArrayEach(pmBrInConf, func(valueInConf []byte, dataType jsonparser.ValueType, offset int, err error) {
			jobNameInConf, _ := jsonparser.GetString(valueInConf, "name")
			toLoad := true

			_, _ = jsonparser.ArrayEach(pmBrInCM, func(valueInCm []byte, dataType jsonparser.ValueType, offset int, err error) {
				jobNameInCM, _ := jsonparser.GetString(valueInCm, "name")
				if jobNameInConf == jobNameInCM {
					toLoad = false
				}
			}, "data", "ericsson-pm:pm", key)

			if toLoad {
				var v interface{}
				_ = json.Unmarshal(valueInConf, &v)
				patchData := PatchData{
					Op:    "add",
					Path:  fmt.Sprintf("/ericsson-pm:pm/%s/-", key),
					Value: v,
				}
				patchDatas = append(patchDatas, patchData)
			}

		}, "data", "ericsson-pm:pm", key)

	}

	// Construct the patch body
	if len(patchDatas) == 0 {
		return nil
	}

	jsonPatchBody := "["
	for i, item := range patchDatas {
		itemInByte, _ := json.Marshal(item)
		if i == 0 {
			jsonPatchBody = jsonPatchBody + string(itemInByte)
		} else {
			jsonPatchBody = jsonPatchBody + "," + string(itemInByte)
		}
	}
	jsonPatchBody = jsonPatchBody + "]"

	return []byte(jsonPatchBody)
}

func constructUnloadData(pmBrInCM, pmBrInConf []byte) []byte {

	jobData := make(map[int]PatchData)
	groupData := make(map[int]PatchData)
	var jobKeys, groupKeys []int

	keys := []string{"job", "group"}
	for _, key := range keys {
		_, _ = jsonparser.ArrayEach(pmBrInConf, func(valueInConf []byte, dataType jsonparser.ValueType, offset int, err error) {
			jobNameInConf, _ := jsonparser.GetString(valueInConf, "name")
			toUnLoad := false

			index := 0
			foundIndex := -1
			_, _ = jsonparser.ArrayEach(pmBrInCM, func(valueInCm []byte, dataType jsonparser.ValueType, offset int, err error) {
				jobNameInCM, _ := jsonparser.GetString(valueInCm, "name")
				if jobNameInConf == jobNameInCM {
					toUnLoad = true
					foundIndex = index
				}
				index = index + 1
			}, "data", "ericsson-pm:pm", key)

			if toUnLoad {
				patchData := PatchData{
					Op:   "remove",
					Path: fmt.Sprintf("/ericsson-pm:pm/%s/%d", key, foundIndex),
				}
				if "job" == key {
					jobData[foundIndex] = patchData
					jobKeys = append(jobKeys, foundIndex)
				} else if "group" == key {
					groupData[foundIndex] = patchData
					groupKeys = append(groupKeys, foundIndex)
				}
			}

		}, "data", "ericsson-pm:pm", key)

	}

	// Construct the patch body, big index item must be removed first.
	jsonPatchBody := ""
	if len(jobKeys) > 0 {
		sort.Ints(jobKeys)
		for _, k := range jobKeys {
			itemInByte, _ := json.Marshal(jobData[k])
			if "" == jsonPatchBody {
				jsonPatchBody = string(itemInByte)
			} else {
				jsonPatchBody = string(itemInByte) + "," + jsonPatchBody
			}
		}
	}

	if len(groupKeys) > 0 {
		sort.Ints(groupKeys)
		for _, k := range groupKeys {
			itemInByte, _ := json.Marshal(groupData[k])
			if "" == jsonPatchBody {
				jsonPatchBody = string(itemInByte)
			} else {
				jsonPatchBody = string(itemInByte) + "," + jsonPatchBody
			}
		}
	}
	jsonPatchBody = "[" + jsonPatchBody + "]"

	return []byte(jsonPatchBody)
}

func getPmBrSchemaUrl() string {
	pmBrSchemaUrlFormat := cmmService + "/schemas/%s"
	return fmt.Sprintf(pmBrSchemaUrlFormat, pmBrSchemaName)
}

func getPmBrServiceConfUrl() string {
	pmBrServiceConfUrlFormat := cmmService + "/configurations/%s"
	return fmt.Sprintf(pmBrServiceConfUrlFormat, pmBulkReporterName)
}

func getPmBrSkeletonConfUrl() string {
	pmBrSkeletonConfUrlFormat := cmmService + "/configurations"
	return fmt.Sprintf(pmBrSkeletonConfUrlFormat)
}
