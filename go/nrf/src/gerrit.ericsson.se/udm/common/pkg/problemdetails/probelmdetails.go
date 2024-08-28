package problemdetails

import (
	"encoding/json"
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/log"
	//"github.com/buger/jsonparser"
)

const (
	abnormalProblemDetails = `{
	  "title": "%s"
    }`
)

// ProblemDetails is used to converted to json in response
type ProblemDetails struct {
	Title         string          `json:"title,omitempty"`
	InvalidParams []*InvalidParam `json:"invalidParams,omitempty"`
        Detail        string          `json:"detail,omitempty"`
        Instance      string          `json:"instance,omitempty"`
        Status        int             `json:"status,omitempty"`
        Type          string          `json:"type,omitempty"`
        Cause         string          `json:"cause,omitempty"`
}

// InvalidParam is used to converted to json in response
type InvalidParam struct {
	Param  string `json:"param"`
	Reason string `json:"reason"`
}

// New return a pointer
func New() *ProblemDetails {
	return &ProblemDetails{}
}

// ToString convert struct to json string
func (p *ProblemDetails) ToString() string {
	//var newData []byte
	data, err := json.Marshal(p)
	if err != nil {
		log.Errorf("ProblemDetails.ToString failed! Errors: %v", err)
		return fmt.Sprintf(abnormalProblemDetails, err.Error())
	}
	return string(data)
/*
	_, valueType, _, err := jsonparser.Get(data, "invalidParams")
	if err != nil {
		log.Errorf("ProblemDetails.ToString failed! Errors: %v", err)
		return fmt.Sprintf(abnormalProblemDetails, err.Error())
	}
	if valueType != jsonparser.Array {
		newData = jsonparser.Delete(data, "invalidParams")
	} else {
		newData = data
	}
	return string(newData[:])
*/
}
