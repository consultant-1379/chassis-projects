package nfdiscservice

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"net/http"
)

//HTTPInfo to store different service process result
type HTTPInfo struct {
	rw        http.ResponseWriter
	req       *http.Request
	body      string
	queryForm nfdiscrequest.DiscGetPara
        cacheTimeout bool
	header       http.Header

	problem            *problemdetails.ProblemDetails
	logcontent         *log.LogStruct
	acrossRegionSearch bool
	statusCode         int
}

//Init to initial HTTPInfo
func (s *HTTPInfo) Init(rw http.ResponseWriter, req *http.Request, queryForm nfdiscrequest.DiscGetPara) {
	s.rw = rw
	s.req = req
	s.queryForm = queryForm
	s.header = make(http.Header)
	s.problem = &problemdetails.ProblemDetails{}
	s.logcontent = &log.LogStruct{}
}

//GetStatusCode to get http statuscode
func (s *HTTPInfo) GetStatusCode() int {
	return s.statusCode
}

//NeedAcrossRegionSearch to get flag whether need to do across region search
func (s *HTTPInfo) NeedAcrossRegionSearch() bool {
	return s.acrossRegionSearch
}

//ResetResponseInfo should executed when pre-service finished, but execute other service except response-service
func (s *HTTPInfo) ResetResponseInfo() {
	s.statusCode = 0
	s.body = ""
	s.problem.Title = ""
	s.logcontent.ResponseDescription = ""
}

//ForwardResponse Forward response information
type ForwardResponse struct {
	statusCode int
	header     http.Header
	body       string
}