package nfdiscservice

import (
	"bytes"
	"com/dbproxy/nfmessage/nrfprofile"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/client"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFDiscAcrossRegionService to process across region
type NFDiscAcrossRegionService struct {
}

//Pair A data structure to hold a key/value pair.
type Pair struct {
	Key    string
	Value  int
	Value2 int
}

//PairList A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		if p[i].Value2 == p[j].Value2 {
			return strings.Compare(p[i].Key, p[j].Key) < 0
		}
		return p[i].Value2 < p[j].Value2
	}
	return p[i].Value < p[j].Value
}

// A function to turn a map into a PairList, then sort with value asc and return it.
func sortMapByValue(m map[string]string) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		num := strings.Split(v, ",")
		value, err1 := strconv.Atoi(num[0])
		value2, err2 := strconv.Atoi(num[1])
		if err1 != nil || err2 != nil {
			log.Debugf("string to int error, err1=%v, err2=%v", err1, err2)
		}
		p[i] = Pair{k, value, value2}
		i++
	}
	sort.Sort(p)
	return p
}

//Execute to process across region
func (a *NFDiscAcrossRegionService) Execute(httpInfo *HTTPInfo) {
	log.Debugf("Enter into across region service")
	if cm.NrfCommon.Role == "region-nrf" {
		if nfdiscrequest.GetNRFDiscForward(httpInfo.req) != constvalue.PLMNDiscForwardValue && nfSupportAcrossRegion(httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)) {
			log.Debugf("RegionNRF send request to PLMNNRF")
			a.sendDiscRequestToPlmnNRF(httpInfo)
		} else {
			httpInfo.logcontent.ResponseDescription = "NRF Not Found matched nfprofile in the local and no configure uplevel NRF"
			httpInfo.problem.Title = "NRF Not Found matched nfprofile in the local and no configure uplevel NRF"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusNotFound
			httpInfo.cacheTimeout = false
		}
	} else {
		nrfAddr, nrfAddrInstanceID := a.getRegionNRFAddr(httpInfo)
		nrfAddr = nfdiscutil.FilterAddrWithVersion(nrfAddr, httpInfo.req.RequestURI)
		if len(nrfAddr) != 0 {
			log.Debugf("Plmn nrf hierarchyDiscoveryMode = %s", cm.DiscoveryService.HierarchyMode)
			if cm.DiscoveryService.HierarchyMode == "redirect" {
				a.returnRedirectInfoToRegionNRF(httpInfo, nrfAddr)
			} else {
				a.forwardDiscRequestToRegionNRF(httpInfo, nrfAddr, nrfAddrInstanceID)
			}
		} else {
			if httpInfo.statusCode == 0 {
				httpInfo.logcontent.ResponseDescription = "NRF Not Found matched nfprofile in the local and not found target region NRF"
				httpInfo.problem.Title = "NRF Not Found matched nfprofile in the local and not found target region NRF"
				httpInfo.body = httpInfo.problem.ToString()
				httpInfo.statusCode = http.StatusNotFound
			}
			httpInfo.cacheTimeout = false
		}
	}
}

func nfSupportAcrossRegion(nfType string) bool {
	nfSupportRegion := map[string]bool{
		"UDM":  true,
		"AMF":  true,
		"SMF":  true,
		"AUSF": true,
		"PCF":  true,
	}

	return nfSupportRegion[nfType]
}

func (a *NFDiscAcrossRegionService) sendDiscRequestToPlmnNRF(httpInfo *HTTPInfo) {
	cm.Mutex.RLock()
	if len(cm.PlmnNrfPriority) == 0 {
		cm.Mutex.RUnlock()

		log.Errorf("Region NRF get PLMN NRF address fail")
		httpInfo.logcontent.RequestDescription = "NRF Discovery not find nfprofile and  Region NRF get PLMN NRF address fail"
		httpInfo.logcontent.ResponseDescription = "Region NRF get PLMN NRF address fail"
		httpInfo.problem.Title = "Region NRF get PLMN NRF address fail"
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.cacheTimeout = false
		httpInfo.body = httpInfo.problem.ToString()
		//handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		return
	}
	body := bytes.NewBufferString("")

	header := make(httpclient.NHeader)
	for k, v := range httpInfo.req.Header {
		header[k] = v[0]
		log.Debugf("key : %s, value: %s", k, v)
	}

	forwardResp := make(map[string]*ForwardResponse)
	forwardSuccess := false
	plmnNRFAddr := ""
	for _, priority := range cm.PlmnNrfPriority {
		for _, addr := range cm.PlmnNrfAPIRootMap[priority] {
			plmnNRFAddr = plmnNRFAddr + "; " + addr
			forwardResp[addr] = &ForwardResponse{statusCode: 0, header: make(http.Header), body: ""}
		}
	}

	for retry := 1; retry <= internalconf.RegionNrfForwardRetryTime; retry++ {
		for _, priority := range cm.PlmnNrfPriority {
			for _, apiRoot := range cm.PlmnNrfAPIRootMap[priority] {
				url := fmt.Sprintf("%s%s", apiRoot, httpInfo.req.RequestURI)
				log.Debugf("Region NRF Discovery Send Request to PLMN NRF Discovery: %s", url)
				httpInfo.logcontent.RequestDescription = fmt.Sprintf("Region NRF Discovery send request to PLMN NRF Discovery : %s", url)
				var resp *httpclient.HttpRespData
				var err error
				if strings.HasPrefix(url, "https") {
					resp, err = client.NoRedirect_https.HttpDo("GET", url, header, body)
				} else {
					resp, err = client.NoRedirect_h2c.HttpDo("GET", url, header, body)
				}
				if err == nil {
					fm.ClearNRFConnectionFailureAlarm(plmnNRFAddr)
					forwardSuccess = true
					if resp.StatusCode == http.StatusTemporaryRedirect {
						log.Debugf("receive redirect code 307, Location = %v", resp.Header.Get("Location"))
						location := resp.Header.Get("Location")
						if location == "" {
							log.Errorf("Plmn redirect location is null")
							httpInfo.logcontent.RequestDescription = fmt.Sprintf("PlmnNRF get NRFProfile URI: %s", httpInfo.req.RequestURI)

							httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Not found target region nrf")
							httpInfo.problem.Title = "Not found target region nrf"
							httpInfo.body = httpInfo.problem.ToString()
							httpInfo.statusCode = http.StatusBadGateway
							httpInfo.cacheTimeout = false
						} else {
							a.redirectDiscRequestToRegionNRF(httpInfo, location)
						}
						cm.Mutex.RUnlock()
						return
					}

					if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified {
						for k, v := range *resp.Header {
							for _, vv := range v {
								httpInfo.rw.Header().Add(k, vv)
							}
						}
						httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF receive discovery response code: %d, body: %s", resp.StatusCode, string(resp.Body))
						httpInfo.body = string(resp.Body)
						httpInfo.statusCode = resp.StatusCode
						httpInfo.cacheTimeout = false
						putCacheProfile(httpInfo.req, resp, cm.PlmnNrfAPIRootInstanceIDMap[apiRoot])
						cm.Mutex.RUnlock()
						return
					}
					for k, v := range *resp.Header {
						for _, vv := range v {
							forwardResp[apiRoot].header.Add(k, vv)
						}
					}
					forwardResp[apiRoot].body = string(resp.Body)
					forwardResp[apiRoot].statusCode = resp.StatusCode

				} else {
					forwardResp[apiRoot].statusCode = http.StatusBadGateway
				}
			}
		}
		time.Sleep(time.Duration(internalconf.RegionNrfForwardRetryWaitTime) * time.Second)
	}
	if !forwardSuccess {
		log.Warnf("Sending message to PLMN NRF %s failed, raise alarm", plmnNRFAddr)
		additionalKey := constvalue.PlmnNRF
		alarmInfo := fmt.Sprintf(constvalue.PlmnNRFInfoFormat, plmnNRFAddr)
		fm.RaiseNRFConnectionFailureAlarm(additionalKey, alarmInfo, plmnNRFAddr)
		log.Errorf("RegionNRF send discovery request to PLMN NRF fail")

		if !httpInfo.cacheTimeout {
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF Send discovery request to PLMN NRF Discovery fail")
			httpInfo.problem.Title = "Region NRF Send discovery request to PLMN NRF Discovery fail"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusBadGateway
		}
	} else {
		var key string
		statusCode := 0
		count := 0
		for k, v := range forwardResp {
			if (v.statusCode != 0 && statusCode == 0) || statusCode > v.statusCode {
				statusCode = v.statusCode
				key = k
			}

			if v.statusCode == http.StatusTooManyRequests || (v.statusCode >= http.StatusInternalServerError && v.statusCode <= 599) {
				count = count + 1
			}
		}

		if key == "" || forwardResp[key].statusCode == 0 {
			httpInfo.problem.Title = "PLMN NRF no response or response status code is 0"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusBadGateway
			httpInfo.cacheTimeout = false
		} else {

			if count != len(forwardResp) || !httpInfo.cacheTimeout {
				for k, v := range forwardResp[key].header {
					for _, vv := range v {
						httpInfo.rw.Header().Add(k, vv)
					}
				}
				httpInfo.body = forwardResp[key].body
				httpInfo.statusCode = forwardResp[key].statusCode
				httpInfo.cacheTimeout = false
			}
		}
	}
	cm.Mutex.RUnlock()
}

func (a *NFDiscAcrossRegionService) getRegionNRFAddr(httpInfo *HTTPInfo) (addr []string, addrInstanceID map[string]string) {
	log.Debugf("PLMN NRF get RegionNRF address from nrfprofile")
	nrfProfileGetRequest := &nrfprofile.NRFProfileGetRequest{}
	plmnDiscGRPCGetRequestFilter(httpInfo.queryForm, nrfProfileGetRequest)
	var nrfAddr []string
	nrfProfileResp, err := dbmgmt.GetNRFProfile(nrfProfileGetRequest)
	if err != nil {
		errorInfo := fmt.Sprintf("PlmnNRF get NRFProfile fail. DB error: %v", err)
		httpInfo.logcontent.RequestDescription = fmt.Sprintf("PlmnNRF get NRFProfile URI: %s", httpInfo.req.RequestURI)
		httpInfo.logcontent.ResponseDescription = fmt.Sprintf("%s", errorInfo)
		httpInfo.problem.Title = errorInfo
		httpInfo.body = httpInfo.problem.ToString()
		httpInfo.statusCode = http.StatusInternalServerError

		return nrfAddr, nil
	}

	if nrfProfileResp.Code != dbmgmt.DbGetSuccess {
		errorInfo := "Not found target region nrf"
		httpInfo.logcontent.RequestDescription = fmt.Sprintf("PlmnNRF get NRFProfile URI: %s", httpInfo.req.RequestURI)
		httpInfo.logcontent.ResponseDescription = fmt.Sprintf("%s", errorInfo)
		httpInfo.problem.Title = errorInfo
		httpInfo.body = httpInfo.problem.ToString()
		httpInfo.statusCode = http.StatusNotFound

		return nrfAddr, nil
	}

	nrfAddr, nrfAddrInstanceID := plmnDiscNRFProfileFilter(nrfProfileResp, httpInfo.queryForm)
	if len(nrfAddr) == 0 {
		errorInfo := "Not found target region nrf"
		httpInfo.logcontent.RequestDescription = fmt.Sprintf("PlmnNRF get NRFProfile URI: %s", httpInfo.req.RequestURI)
		httpInfo.logcontent.ResponseDescription = fmt.Sprintf("%s", errorInfo)
		httpInfo.problem.Title = errorInfo
		httpInfo.body = httpInfo.problem.ToString()
		httpInfo.statusCode = http.StatusNotFound
		//handleDiscoveryFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())
		return nrfAddr, nil
	}
	return nrfAddr, nrfAddrInstanceID
}


func (a *NFDiscAcrossRegionService) forwardDiscRequestToRegionNRF(httpInfo *HTTPInfo, addrList []string, addrInstanceID map[string]string) {
	log.Debugf("PLMN NRF send request to RegionNRF")

	header := make(httpclient.NHeader)
	for k, v := range httpInfo.req.Header {
		header[k] = v[0]
	}
	forwardSuccess := false
	forwardResp := make(map[string]*ForwardResponse)
	regionNRFAddr := ""
	addrMap := make(map[string]bool)
	for _, value := range addrList {
		addrMap[value] = true
		regionNRFAddr = regionNRFAddr + "; " + value
		forwardResp[value] = &ForwardResponse{statusCode: 0, header: make(http.Header), body: ""}
	}
	header[constvalue.SearchDataForward] = constvalue.PLMNDiscForwardValue
	requestParam := nfdiscutil.GetRequestParam(httpInfo.req.RequestURI)
	for retry := 1; retry <= internalconf.PlmnNrfForwardRetryTime; retry++ {
		for _, addr := range addrList {
			if !addrMap[addr] {
				continue
			}
			url := fmt.Sprintf("%s%s", addr, requestParam)
			httpInfo.logcontent.RequestDescription = fmt.Sprintf("PLMN NRF Discovery forward discovery request to Region NRF: %s", url)
			resp, err := nfdiscutil.DiscHTTPDo("GET", url, header, bytes.NewBufferString(""), cm.DiscNRFSelfAPIURI)
			if err == nil {
				fm.ClearNRFConnectionFailureAlarm(regionNRFAddr)
				forwardSuccess = true

				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified {
					httpInfo.logcontent.ResponseDescription = fmt.Sprintf("PLMN NRF Discovery forward discovery request to Region NRF, Resp Status Code: %d Body: %s", resp.StatusCode, string(resp.Body))
					for k, v := range *resp.Header {
						for _, vv := range v {
							httpInfo.rw.Header().Add(k, vv)
						}
					}
					httpInfo.body = string(resp.Body)
					httpInfo.statusCode = resp.StatusCode
					httpInfo.cacheTimeout = false
					putCacheProfile(httpInfo.req, resp, addrInstanceID[addr])
					return
				}
				for k, v := range *resp.Header {
					for _, vv := range v {
						forwardResp[addr].header.Add(k, vv)
					}
				}
				forwardResp[addr].body = string(resp.Body)
				forwardResp[addr].statusCode = resp.StatusCode

				if resp.StatusCode != 0 && resp.StatusCode < http.StatusInternalServerError {
					addrMap[addr] = false
				}

			} else {
				forwardResp[addr].statusCode = http.StatusBadGateway
			}
		}

		time.Sleep(time.Duration(internalconf.PlmnNrfForwardRetryWaitTime) * time.Second)

	}

	if !forwardSuccess {
		additionKey := constvalue.RegionNRF
		alarmInfo := fmt.Sprintf(constvalue.RegionNRFInfoFormat, regionNRFAddr)
		fm.RaiseNRFConnectionFailureAlarm(additionKey, alarmInfo, regionNRFAddr)
		log.Errorf("RegionNRF Address: %s", regionNRFAddr)

		if !httpInfo.cacheTimeout {
			httpInfo.problem.Title = "PLMN NRF forward discovery request to Region NRF fail"
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("PLMN NRF forward discovery request to Region NRF fail: %s RequestURI: %s", regionNRFAddr, requestParam)
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusBadGateway
		}
	} else {
		var key string
		statusCode := 0
		count := 0

		for k, v := range forwardResp {
			log.Debugf("All Forward Response: %v, %v", k, v)
			if (v.statusCode != 0 && statusCode == 0) || statusCode > v.statusCode {
				statusCode = v.statusCode
				key = k
			}

			if v.statusCode == http.StatusTooManyRequests || (v.statusCode >= http.StatusInternalServerError && v.statusCode <= 599) {
				count = count + 1
			}

		}

		if key == "" || forwardResp[key].statusCode == 0 {
			httpInfo.statusCode = http.StatusBadGateway
			httpInfo.problem.Title = "PLMN NRF forward discovery request to Region NRF fail"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.cacheTimeout = false
		} else {
			if count != len(forwardResp) || !httpInfo.cacheTimeout {
				for k, v := range forwardResp[key].header {
					for _, vv := range v {
						httpInfo.rw.Header().Add(k, vv)
					}
				}
				httpInfo.body = forwardResp[key].body
				httpInfo.statusCode = forwardResp[key].statusCode
				httpInfo.cacheTimeout = false
			}

		}

	}

}

func (a *NFDiscAcrossRegionService) returnRedirectInfoToRegionNRF(httpInfo *HTTPInfo, addrList []string) {
	log.Debugf("Plmn NRF return 307 Temporary Redirect to region NRF")
	httpInfo.statusCode = http.StatusTemporaryRedirect
	//select the first addr from sorted list by priority asc(lower values indicate a higher priority)
	requestParam := nfdiscutil.GetRequestParam(httpInfo.req.RequestURI)
	httpInfo.rw.Header().Set("Location", fmt.Sprintf("%s%s", addrList[0], requestParam))
	httpInfo.logcontent.ResponseDescription = fmt.Sprintf("PLMN NRF Discovery return region nrf address to Region NRF, Resp Status Code: %d Temporary Redirect Location: %s", httpInfo.statusCode, addrList[0])

}

func (a *NFDiscAcrossRegionService) redirectDiscRequestToRegionNRF(httpInfo *HTTPInfo, location string) {
	log.Debugf("Region NRF send redirect request to RegionNRF, redirect url: %s", location)
	rootURL := nfdiscutil.GetRequestURIRoot(location)
	cm.Mutex.RLock()
	for _, v := range cm.DiscNRFSelfAPIURI {
		if v == rootURL {
			httpInfo.logcontent.ResponseDescription = "redirect region nrf address is the same as self address"
			httpInfo.statusCode = http.StatusBadGateway
			httpInfo.problem.Title = "redirect region nrf address is the same as self address"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.cacheTimeout = false
			cm.Mutex.RUnlock()
			return
		}
	}
	cm.Mutex.RUnlock()
	forwardSuccess := false
	forwardResp := make(map[string]*ForwardResponse)
	forwardResp[location] = &ForwardResponse{statusCode: 0, header: make(http.Header), body: ""}
	header := make(httpclient.NHeader)
	for k, v := range httpInfo.req.Header {
		header[k] = v[0]
	}
	for retry := 1; retry <= internalconf.RegionNrfRedirectRetryTime; retry++ {
		httpInfo.logcontent.RequestDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF: %s", location)
		resp, err := nfdiscutil.DiscHTTPDo("GET", location, header, bytes.NewBufferString(""), cm.DiscNRFSelfAPIURI)
		/*if err == nil && resp.StatusCode == http.StatusServiceUnavailable && retry == internalconf.RegionNrfRedirectRetryTime {
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF, Resp Status Code: %d Body: %s", resp.StatusCode, string(resp.Body))
			httpInfo.statusCode = resp.StatusCode
			httpInfo.body = string(resp.Body)
			httpInfo.statusCode = false
			return
		}*/
		if err == nil && resp.StatusCode == http.StatusTemporaryRedirect {
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF redirect discovery request to Region NRF fail: %s RequestURI: %s, redirect more than 3 times", location, httpInfo.req.RequestURI)
			httpInfo.problem.Title = "Region NRF redirect discovery request to Region NRF fail, redirect more than 3 times"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusInternalServerError
			httpInfo.cacheTimeout = false
			return
		}

		if err == nil { //&& nfdiscutil.StatusCodeDirectReturn(resp.StatusCode) {
			forwardSuccess = true
			fm.ClearNRFConnectionFailureAlarm(rootURL)
			if resp.StatusCode == http.StatusNotModified || resp.StatusCode == http.StatusOK {
				for k, v := range *resp.Header {
					for _, vv := range v {
						httpInfo.rw.Header().Add(k, vv)
					}
				}
				httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF, Resp Status Code: %d Body: %s", resp.StatusCode, string(resp.Body))
				httpInfo.statusCode = resp.StatusCode
				httpInfo.body = string(resp.Body)
				httpInfo.cacheTimeout = false
				putCacheProfile(httpInfo.req, resp, RemoteCacheFromRedirect)
				return
			}
			for k, v := range *resp.Header {
				for _, vv := range v {
					forwardResp[location].header.Add(k, vv)
				}
			}
			forwardResp[location].statusCode = resp.StatusCode
			forwardResp[location].body = string(resp.Body)

		} else {
			forwardResp[location].statusCode = http.StatusBadGateway
		}
		time.Sleep(time.Duration(internalconf.RegionNrfRedirectRetryWaitTime) * time.Second)
	}

	if !forwardSuccess {
		additionKey := constvalue.RegionNRF
		alarmInfo := fmt.Sprintf(constvalue.RegionNRFInfoFormat, location)
		fm.RaiseNRFConnectionFailureAlarm(additionKey, alarmInfo, rootURL)
		log.Errorf("RegionNRF Address: %s", location)
		httpInfo.logcontent.RequestDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF: %s RequestURI: %s", location, httpInfo.req.RequestURI)
		httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF redirect discovery request to Region NRF fail: %s RequestURI: %s", location, httpInfo.req.RequestURI)
		log.Errorf("RegionNRF Address: %s", location)
		if !httpInfo.cacheTimeout {
			httpInfo.logcontent.RequestDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF: %s RequestURI: %s", location, httpInfo.req.RequestURI)
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF redirect discovery request to Region NRF fail: %s RequestURI: %s", location, httpInfo.req.RequestURI)
			httpInfo.problem.Title = "Region NRF redirect discovery request to Region NRF fail"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusBadGateway
		}
	} else {
		if forwardResp[location].statusCode == 0 {
			httpInfo.logcontent.RequestDescription = fmt.Sprintf("Region NRF Discovery redirect discovery request to Region NRF: %s RequestURI: %s", location, httpInfo.req.RequestURI)
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf("Region NRF redirect discovery request to Region NRF fail: %s RequestURI: %s", location, httpInfo.req.RequestURI)
			httpInfo.problem.Title = "Region NRF redirect discovery request to Region NRF fail"
			httpInfo.body = httpInfo.problem.ToString()
			httpInfo.statusCode = http.StatusBadGateway
			httpInfo.cacheTimeout = false
		} else {
			if !((forwardResp[location].statusCode == http.StatusTooManyRequests || (forwardResp[location].statusCode >= http.StatusInternalServerError && forwardResp[location].statusCode <= 599)) && httpInfo.cacheTimeout) {
				for k, v := range forwardResp[location].header {
					for _, vv := range v {
						httpInfo.rw.Header().Add(k, vv)
					}
				}

				httpInfo.body = forwardResp[location].body
				httpInfo.statusCode = forwardResp[location].statusCode
				httpInfo.cacheTimeout = false
			}
		}
	}
}
