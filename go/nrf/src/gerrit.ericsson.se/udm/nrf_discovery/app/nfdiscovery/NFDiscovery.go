package nfdiscovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/client"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/multisite"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscservice"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

var (
	// masterPool *simpleredis.ConnectionPool
	// slavePool  *simpleredis.ConnectionPool
	randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func handleDiscoveryFailure(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct,
	requestnf string, statusCode int, body string) {
	log.Debugf(constvalue.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	nftype := "unknown"
	if len(requestnf) > 0 {
		nftype = requestnf
	}
	pm.Inc(constvalue.NfDiscoveryFailureTotal, constvalue.NfDiscovery, req.Method, strings.Replace(nftype, "-", "_", -1), fmt.Sprintf("%d", statusCode), "_")
	DiscoveryResponseHander(rw, req, nftype, statusCode, body)
	log.Errorf(constvalue.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

func handleDiscoverySuccess(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct,
	requestnf string, statusCode int, body string) {
	log.Debugf(constvalue.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	nftype := "unknown"
	if len(requestnf) > 0 {
		nftype = requestnf
	}
	pm.Inc(constvalue.NfDiscoverySuccessTotal, constvalue.NfDiscovery, req.Method, strings.Replace(nftype, "-", "_", -1), fmt.Sprintf("%d", statusCode))
	DiscoveryResponseHander(rw, req, nftype, statusCode, body)
	log.Debugf(constvalue.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

// DiscoveryResponseHander handle response for  discovery
func DiscoveryResponseHander(rw http.ResponseWriter, req *http.Request, requestnf string, statuscode int, body string) {
	pm.Inc(constvalue.NfDiscoveryRequestsTotal, constvalue.NfDiscovery, req.Method, strings.Replace(requestnf, "-", "_", -1))
	if statuscode != http.StatusOK {
		if internalconf.HTTPWithXVersion {
			rw.Header().Set("X-Version", cm.ServiceVersion)
		}
		rw.Header().Set("Content-Type", "application/problem+json")
		rw.WriteHeader(statuscode)
		if body != "" {
			_, err := rw.Write([]byte(body))
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	} else {
		if internalconf.HTTPWithXVersion {
			rw.Header().Set("X-Version", cm.ServiceVersion)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(statuscode)
		if body != "" {
			_, err := rw.Write([]byte(body))
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	}
}

//NFDiscService for discovery process discovery request
type NFDiscService struct {
	httpInfo *nfdiscservice.HTTPInfo

	RoamingService      *nfdiscservice.NFDiscRoamingService
	AcrossRegionService *nfdiscservice.NFDiscAcrossRegionService
	LocalSearchService  *nfdiscservice.NFDiscLocalSearchService
	ResponseService     *nfdiscservice.NFDiscResponseService
}

func isNeedRoaming(queryForm nfdiscrequest.DiscGetPara) (bool, bool) {
	roaming := false
	localSearch := false
	targetPlmnList := queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
	var homePlmnList []string
	if len(targetPlmnList) > 0 {
		for _, plmn := range cm.NfProfile.PlmnID {
			if len(plmn.Mnc) == 2 {
				homePlmnList = append(homePlmnList, plmn.Mcc+"0"+plmn.Mnc)
			} else {
				homePlmnList = append(homePlmnList, plmn.Mcc+plmn.Mnc)
			}
		}
		log.Debugf("NRF Home Plmn List: %s, target Plmn List: %s", homePlmnList, targetPlmnList)
		for _, targetPlmn := range targetPlmnList {
			var matched bool
			if len(targetPlmn) == 5 {
				plmnArray := []rune(targetPlmn)
				targetPlmn = string(plmnArray[0:3]) + "0" + string(plmnArray[3:])
			}
			for _, homePlmn := range homePlmnList {
				if homePlmn == targetPlmn {
					matched = true
					break
				}
			}
			if matched {
				localSearch = true
			} else {
				roaming = true
			}
		}
		return roaming, localSearch
	}

	return roaming, localSearch
}

func isForwardToHomeNRF(queryForm nfdiscrequest.DiscGetPara) bool {
	if "" == queryForm.GetNRFDiscHnrfURI() {
		return false
	}

	cm.Mutex.RLock()
	for _, v := range cm.DiscNRFSelfAPIURI {
		if v == queryForm.GetNRFDiscHnrfURI() {
			return false
		}
	}
	cm.Mutex.RUnlock()
	return true
}

func forwardToHomeNRF(rw http.ResponseWriter, req *http.Request, queryForm nfdiscrequest.DiscGetPara) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}

	logcontent := &log.LogStruct{SequenceId: sequenceID}
	homeURI := queryForm.GetNRFDiscHnrfURI()

	var resp *httpclient.HttpRespData
	var err error
	body := bytes.NewBufferString("")
	header := make(httpclient.NHeader)

	for k, v := range req.Header {
		header[k] = v[0]
		log.Debugf("key : %s, value: %s", k, v)
	}

	if strings.HasPrefix(homeURI, "https") {
		resp, err = client.NoRedirect_https.HttpDo("GET", homeURI+"nf-instances?"+req.URL.RawQuery, header, body)
	} else if strings.HasPrefix(homeURI, "http") {
		resp, err = client.NoRedirect_h2c.HttpDo("GET", homeURI+"nf-instances?"+req.URL.RawQuery, header, body)
	}

	requestnftype := queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType)
	if err == nil {
		logcontent.RequestDescription = fmt.Sprintf(`NRF Forward request to Home NRF {"target-nf-type":"%s", "requester-nf-type":"%s" "hnrf-uri" : "%s"}`, queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), requestnftype, homeURI)
		logcontent.ResponseDescription = fmt.Sprintf("NRF Forward request to Home NRF success")
		if resp.Header != nil {
			for k, v := range *resp.Header {
				for _, vv := range v {
					rw.Header().Add(k, vv)
				}
			}
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified {
			handleDiscoverySuccess(rw, req, logcontent, requestnftype, resp.StatusCode, string(resp.Body))
			return
		}
		if resp.StatusCode == http.StatusServiceUnavailable || nfdiscutil.StatusCodeDirectReturn(resp.StatusCode) {
			handleDiscoveryFailure(rw, req, logcontent, requestnftype, resp.StatusCode, string(resp.Body))
			return
		}
		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusNotImplemented {
			handleDiscoveryFailure(rw, req, logcontent, requestnftype, http.StatusBadGateway, string(resp.Body))
			return
		}

	}

	logcontent.RequestDescription = fmt.Sprintf(`NRF Forward request to Home NRF {"target-nf-type":"%s", "requester-nf-type":"%s"}`, queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
	logcontent.ResponseDescription = fmt.Sprintf("NRF Forward request to Home NRF Fail")
	problem := problemdetails.ProblemDetails{
		Title: "NRF Forward request to Home NRF Fail",
	}
	handleDiscoveryFailure(rw, req, logcontent, requestnftype, http.StatusBadGateway, problem.ToString())
	log.Warningf("NRF Discovery Forward request to HOME NRF Fail. Home NRF URI: %s", homeURI)
}

//NrfDiscGetHandler handler function
func NrfDiscGetHandler(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debugf("NF Discovery Request comes")
	var problemDetails *problemdetails.ProblemDetails
	oriqueryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		errorInfo := fmt.Sprintf("Parse URL err: %s", err.Error())
		problemDetails = &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		handleDiscoveryFailure(rw, req, logcontent, "", http.StatusBadRequest, problemDetails.ToString())
		return
	}
	for k, v := range oriqueryForm {
		log.Debugf("oriqueryForm Key: %s", k)
		for i, value := range v {
			log.Debugf("oriqueryForm Value[%d]: %s", i, value)
		}
	}
	var queryForm nfdiscrequest.DiscGetPara
	//queryForm.value = oriqueryForm
	queryForm.InitMember(oriqueryForm)
	problem := queryForm.ValidateNRFDiscovery()
	if problem != nil {
		logcontent.RequestDescription = fmt.Sprintf(`{"service-names":%v, "target-nf-type":"%s", "requester-nf-type":"%s"}`, "[]", "", "")
		logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, problem.ToString())
		handleDiscoveryFailure(rw, req, logcontent, "", http.StatusBadRequest, problem.ToString())
		return
	}
	//If request with hnrf-uri parameter, directly forward it to HOMENRF
	if isForwardToHomeNRF(queryForm) {
		forwardToHomeNRF(rw, req, queryForm)
		return
	}
	nfService := &NFDiscService{
		httpInfo:            &nfdiscservice.HTTPInfo{},
		RoamingService:      &nfdiscservice.NFDiscRoamingService{},
		LocalSearchService:  &nfdiscservice.NFDiscLocalSearchService{},
		AcrossRegionService: &nfdiscservice.NFDiscAcrossRegionService{},
		ResponseService:     &nfdiscservice.NFDiscResponseService{},
	}
	nfService.httpInfo.Init(rw, req, queryForm)
	roaming, localSearch := isNeedRoaming(queryForm)
	if roaming && !localSearch {
		nfService.RoamingService.Execute(nfService.httpInfo)
		nfService.ResponseService.Execute(nfService.httpInfo)
		return
	} else if roaming && localSearch {
		nfService.LocalSearchService.Execute(nfService.httpInfo)
		if nfService.httpInfo.GetStatusCode() == http.StatusNotFound {
			nfService.httpInfo.ResetResponseInfo()
			nfService.RoamingService.Execute(nfService.httpInfo)
		}

		nfService.ResponseService.Execute(nfService.httpInfo)
		return
	}

	nfService.LocalSearchService.Execute(nfService.httpInfo)
	if nfService.httpInfo.NeedAcrossRegionSearch() {
		//nfService.httpInfo.ResetResponseInfo()
		nfService.AcrossRegionService.Execute(nfService.httpInfo)
	}
	nfService.ResponseService.Execute(nfService.httpInfo)
}

//NrfDiscSearchGetHandler is a entrance for all request
func NrfDiscSearchGetHandler(rw http.ResponseWriter, req *http.Request) {
	if configmap.MultisiteEnabled && !multisite.GetMonitor().IsActiveSite() {
		logcontent, problemDetails := multisite.GetIsolateMessage()
		handleDiscoveryFailure(rw, req, logcontent, "", http.StatusServiceUnavailable, problemDetails.ToString())
		return
	}

	startedAt := time.Now()
	defer func() {
		pm.Observe(float64(time.Since(startedAt))/float64(time.Second), constvalue.NfRequestDuration, constvalue.NfDiscovery)
		if internalconf.EnableTimeStatistics {
			entTime := time.Now().UnixNano() / 1000000
			dbmgmt.DBLatency.HandlerChannel <- dbmgmt.Latency{ReqStartTime: startedAt.UnixNano() / 1000000, ReqEndTime: entTime}
		}
	}()
	NrfDiscGetHandler(rw, req)
}

//HealthCheckHandler is a health check handler
func HealthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

//TestAllHandler is a test handler
func TestAllHandler(rw http.ResponseWriter, req *http.Request) {
	value := req.Header.Get("Test-Network-Delay")

	if value != "" {
		if d, err := strconv.Atoi(value); err == nil {
			if d != 0 {
				t := randSeed.Intn(d)
				delay := time.Duration(t) * time.Microsecond
				time.Sleep(delay)
			}
		} else {
			log.Warningf("Please input the incorrect value of Header Test-Network-Delay, %s, %s",
				value, err.Error())
		}
	}

	if req.Method == "POST" || req.Method == "PUT" {
		if body, err := ioutil.ReadAll(req.Body); err == nil {
			log.Debugf("%s", string(body))
			if req.Method == "POST" {
				rw.WriteHeader(http.StatusNoContent)
			} else {
				rw.WriteHeader(http.StatusOK)
				_, err := rw.Write(body)
				if err != nil {
					log.Warnf("%v", err)
				}
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}
	} else if req.Method == "GET" {
		rw.WriteHeader(http.StatusOK)
		if body, err := getHostAllEnvInfof(); err == nil {
			_, err := rw.Write(body)
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	}

}

func getHostAllEnvInfof() ([]byte, error) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	return json.MarshalIndent(environment, "", "  ")
}
