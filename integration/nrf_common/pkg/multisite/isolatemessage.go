package multisite

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
)

//GetIsolateMessage returns IsolateMessage
func GetIsolateMessage() (*log.LogStruct, *problemdetails.ProblemDetails) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}
	logcontent.RequestDescription = "Data not sync with other site"
	logcontent.ResponseDescription = "Data not sync with other site"
	problemDetails := &problemdetails.ProblemDetails{
		Title:  "Service degraded or unavailable",
		Detail: "Data integrity or availability can not be guaranteed for NRF geographical redundancy function.",
	}
	return logcontent, problemDetails
}
