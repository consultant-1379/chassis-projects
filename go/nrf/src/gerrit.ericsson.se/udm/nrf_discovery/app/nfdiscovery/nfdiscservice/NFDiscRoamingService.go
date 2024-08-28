package nfdiscservice

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"com/dbproxy/nfmessage/nrfaddress"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"

	"time"

	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"github.com/buger/jsonparser"
)

//NFDiscRoamingService to process roaming
type NFDiscRoamingService struct {
}

//Execute to execute roaming
func (r *NFDiscRoamingService) Execute(httpInfo *HTTPInfo) {
	addr := r.getHomeNRFAddr(httpInfo)
	if len(addr) == 0 {
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.problem.Title = "Not found Home NRF address"
		httpInfo.logcontent.RequestDescription = "NRF roaming request to Home NRF"
		httpInfo.logcontent.ResponseDescription = "NRF not found Home NRF address"
		httpInfo.body = httpInfo.problem.ToString()
		return
	}

	homeNRFAddr := ""
	urlMap := make(map[string]bool)
	for _, v := range addr {
		url := v + httpInfo.req.RequestURI
		urlMap[url] = true
		homeNRFAddr = homeNRFAddr + "; " + v
	}
	for retryTime := 1; retryTime <= internalconf.HomeNrfForwardRetryTime; retryTime++ {
		var isSearched bool
		isSearched, urlMap = r.doRoming(httpInfo, urlMap)
		if isSearched {
			fm.ClearNRFConnectionFailureAlarm(homeNRFAddr)
			break
		}
		//don't need to sleep 1 second while the third time
		if retryTime < internalconf.HomeNrfForwardRetryTime {
			time.Sleep(time.Duration(internalconf.HomeNrfForwardRetryWaitTime) * time.Second)
		}
		//if roaming fail 3 times, print the log
		if retryTime == internalconf.HomeNrfForwardRetryTime {
			if httpInfo.statusCode == 0 || httpInfo.statusCode >= http.StatusInternalServerError {
				targetPlmnList := httpInfo.queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
				mcc := targetPlmnList[0][:3]
				mnc := targetPlmnList[0][3:]
				additionalKey := constvalue.HomeNRF
				nrfInfo := fmt.Sprintf(constvalue.HomeNRFInfoFormat, mcc, mnc, addr)
				log.Errorf("Forward discovery request to NRF %s failed", nrfInfo)
				fm.RaiseNRFConnectionFailureAlarm(additionalKey, nrfInfo, homeNRFAddr)

				if httpInfo.statusCode == 0 {
					httpInfo.problem.Title = "Get NFProfile from roaming nrf fail"
					httpInfo.logcontent.RequestDescription = "NRF roaming request to Home NRF fail"
					httpInfo.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, "Get NFProfile from roaming nrf fail")
					httpInfo.statusCode = http.StatusBadGateway
				}
				httpInfo.body = httpInfo.problem.ToString()
			}
		}
	}
}

func (r *NFDiscRoamingService) doDiscRoamingSearch(url string, req *http.Request, queryForm nfdiscrequest.DiscGetPara) (response *httpclient.HttpRespData, errinfo error) {
	body := bytes.NewBufferString("")
	header := make(httpclient.NHeader)

	for k, v := range req.Header {
		header[k] = v[0]
		log.Debugf("key : %s, value: %s", k, v)
	}

	return nfdiscutil.DiscHTTPDo("GET", url, header, body, cm.DiscNRFSelfAPIURI)
}

func (r *NFDiscRoamingService) doRoming(httpInfo *HTTPInfo, url map[string]bool) (bool, map[string]bool) {
	log.Debugf("NRF Roaming")
	isSearched := false

	for key, value := range url {
		log.Debugf("HomeNRF URI: %s, %v", key, value)
		if !value {
			continue
		}
		resp, err := r.doDiscRoamingSearch(key, httpInfo.req, httpInfo.queryForm)
		if err == nil {
			httpInfo.logcontent.RequestDescription = fmt.Sprintf(`NRF Roaming {"target-nf-type":"%s", "requester-nf-type":"%s"}`, httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
			httpInfo.logcontent.ResponseDescription = fmt.Sprintf(`"NRF Roaming StatusCode %d Response: %s"`, resp.StatusCode, string(resp.Body))
			//
			//make sure nrf response code is certain, all always response the min status code2
			if httpInfo.statusCode == 0 || (httpInfo.statusCode > resp.StatusCode || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
				for k, v := range *resp.Header {
					for _, vv := range v {
						httpInfo.rw.Header().Add(k, vv)
					}
				}
				if resp.StatusCode == http.StatusTemporaryRedirect {
					httpInfo.statusCode = http.StatusInternalServerError
				} else {
					httpInfo.statusCode = resp.StatusCode
				}
				httpInfo.body = string(resp.Body)
			}

			if resp.StatusCode != 0 && resp.StatusCode < http.StatusInternalServerError {
				url[key] = false
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified {
				isSearched = true
				return isSearched, url
			}

		}
	}

	return isSearched, url
}

func (r *NFDiscRoamingService) getRemotePLMN(queryForm nfdiscrequest.DiscGetPara) []string {
	targetPlmnList := queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
	var homePlmnList []string
	var remotePlmnList []string
	for _, plmn := range cm.NfProfile.PlmnID {
		if len(plmn.Mnc) == 2 {
			homePlmnList = append(homePlmnList, plmn.Mcc+"0"+plmn.Mnc)
		} else {
			homePlmnList = append(homePlmnList, plmn.Mcc+plmn.Mnc)
		}
	}

	for _, targetPlmn := range targetPlmnList {
		if len(targetPlmn) == 5 {
			plmnArray := []rune(targetPlmn)
			targetPlmn = string(plmnArray[0:3]) + "0" + string(plmnArray[3:])
		}
		var matched bool
		for _, homePlmn := range homePlmnList {
			if homePlmn == targetPlmn {
				matched = true
				break
			}
		}

		if !matched {
			remotePlmnList = append(remotePlmnList, targetPlmn)
		}
	}

	return remotePlmnList
}

func (r *NFDiscRoamingService) getHomeNRFAddr(httpInfo *HTTPInfo) []string {
	url := make([]string, 0)

	targetPlmnList := r.getRemotePLMN(httpInfo.queryForm)

	index := &nrfaddress.NRFAddressGetIndex{
		NrfAddressKey1: targetPlmnList,
	}

	filter := &nrfaddress.NRFAddressFilter{
		Index: index,
	}

	nrfAddressFilter := &nrfaddress.NRFAddressGetRequest_Filter{
		Filter: filter,
	}

	nrfAddressGetRequest := &nrfaddress.NRFAddressGetRequest{
		Data: nrfAddressFilter,
	}
	scheme := ""
	addr := ""
	port := int64(0)
	nrfAddr := make([]dbmgmt.NRFAddress, 0)

	nrfAddressResponse, err := dbmgmt.GetNRFAddress(nrfAddressGetRequest)
	if err == nil && nrfAddressResponse.Code == dbmgmt.DbGetSuccess {
		log.Debugf("NRFAddress Result: %s", nrfAddressResponse.NrfAddressData)
		for _, item := range nrfAddressResponse.NrfAddressData {
			_, parserErr := jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
				var err2 error
				scheme, err2 = jsonparser.GetString(value, "scheme")
				if err2 != nil {
					scheme = ""
					return
				}
				//select a available addr from fqdn, ipv4 and ipv6
				var flag bool
				flag, addr = r.getAddrFromNrfAddressData(value)
				if false == flag {
					scheme = ""
					addr = ""
					return
				}

				port, err2 = jsonparser.GetInt(value, "port")
				if err2 != nil {
					scheme = ""
					addr = ""
					port = 0
					return
				}
				nrfAddr = append(nrfAddr, dbmgmt.NRFAddress{Scheme: scheme, Fqdn: addr, Port: int(port)})
			}, "nrfAddresses")
			if parserErr != nil {
				log.Debugf("nrfAddresses parse error, err=%v", parserErr)
			}
		}
	}
	if len(nrfAddr) == 0 {
		log.Debugf("Get scheme or fqdn or port from NRFAddress fail, will use default value do Roaming")
		plmnList := httpInfo.queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
		for _, plmn := range plmnList {
			//mcc, _ := jsonparser.GetString([]byte(httpInfo.queryForm.GetValue()[constvalue.SearchDataTargetPlmnList][i]), constvalue.SearchDataMcc)
			//mnc, _ := jsonparser.GetString([]byte(httpInfo.queryForm.GetValue()[constvalue.SearchDataTargetPlmnList][i]), constvalue.SearchDataMnc)
			mcc := plmn[0:3]
			mnc := plmn[3:]
			if len(mnc) == 2 {
				mnc = "0" + mnc
			}
			fqdn := fmt.Sprintf(constvalue.NRFFqdnFormat, mnc, mcc)
			scheme = cm.NrfCommon.RemoteDefaultSetting.Scheme
			port = int64(cm.NrfCommon.RemoteDefaultSetting.Port)

			nrfAddr = append(nrfAddr, dbmgmt.NRFAddress{Scheme: scheme, Fqdn: fqdn, Port: int(port)})
		}
	}

	for _, v := range nrfAddr {
		if strings.Contains(v.Scheme, "://") {
			url = append(url, (v.Scheme + v.Fqdn + ":" + strconv.Itoa(v.Port)))
		} else {
			url = append(url, (v.Scheme + "://" + v.Fqdn + ":" + strconv.Itoa(v.Port)))
		}

	}
	log.Debugf("Roaming NRF URL: %s", url)
	return url
}

func (r *NFDiscRoamingService) getAddrFromNrfAddressData(value []byte) (bool, string) {
	//select a available addr from fqdn, ipv4 and ipv6
	fqdn, errFqdn := jsonparser.GetString(value, "address", "fqdn")
	if errFqdn == nil {
		return true, fqdn
	}

	ipv4, errIPV4 := jsonparser.GetString(value, "address", "ipv4Address")
	if errIPV4 == nil {
		return true, ipv4
	}
	ipv6, errIPV6 := jsonparser.GetString(value, "address", "ipv6Address")
	if errIPV6 == nil {
		return true, ipv6
	}
	return false, ""
}
