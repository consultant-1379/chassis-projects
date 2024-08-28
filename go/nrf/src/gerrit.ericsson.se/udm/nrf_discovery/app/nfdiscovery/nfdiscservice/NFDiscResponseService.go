package nfdiscservice

import (
	"fmt"
	"net/http"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

//NFDiscResponseService to send response to NF consumer
type NFDiscResponseService struct {
}

//Execute to send response to NF consumer
func (r *NFDiscResponseService) Execute(httpInfo *HTTPInfo) {
	if httpInfo.statusCode == http.StatusOK || httpInfo.statusCode == http.StatusNotModified {
		if httpInfo.cacheTimeout {
			for k, v := range httpInfo.header {
				for _, vv := range v {
					httpInfo.rw.Header().Add(k,vv)
				}
			}
		}
		log.Debugf("http response body: %v", httpInfo.body)
		r.handleDiscoverySuccess(httpInfo.rw, httpInfo.req, httpInfo.logcontent,
			httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType), httpInfo.statusCode, httpInfo.body)
	} else {
		r.handleDiscoveryFailure(httpInfo.rw, httpInfo.req, httpInfo.logcontent,
			httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType), httpInfo.statusCode, httpInfo.problem.ToString())
	}

}

func (r *NFDiscResponseService) handleDiscoveryFailure(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct,
	requestnf string, statusCode int, body string) {
	log.Debugf(constvalue.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	nftype := "unknown"
	if len(requestnf) > 0 {
		nftype = requestnf
	}
	pm.Inc(constvalue.NfDiscoveryFailureTotal, constvalue.NfDiscovery, req.Method, strings.Replace(nftype, "-", "_", -1), fmt.Sprintf("%d", statusCode), "_")
	r.discoveryResponseHander(rw, req, nftype, statusCode, body)
	log.Errorf(constvalue.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

func (r *NFDiscResponseService) handleDiscoverySuccess(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct,
	requestnf string, statusCode int, body string) {
	log.Debugf(constvalue.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	nftype := "unknown"
	if len(requestnf) > 0 {
		nftype = requestnf
	}
	pm.Inc(constvalue.NfDiscoverySuccessTotal, constvalue.NfDiscovery, req.Method, strings.Replace(nftype, "-", "_", -1), fmt.Sprintf("%d", statusCode))
	r.discoveryResponseHander(rw, req, nftype, statusCode, body)
	log.Debugf(constvalue.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

// discoveryResponseHander handle response for  discovery
func (r *NFDiscResponseService) discoveryResponseHander(rw http.ResponseWriter, req *http.Request, requestnf string, statuscode int, body string) {
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
