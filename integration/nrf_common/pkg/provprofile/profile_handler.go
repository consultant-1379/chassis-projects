package provprofile

import (
	"net/http"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"crypto/md5"
	"fmt"
	"gerrit.ericsson.se/udm/common/pkg/utils"
)
//ProfileHandler to handle different request for provision
type ProfileHandler interface {
	//PostHandler to process POST request
	PostHandler()
	//PutHandler to process PUT request
	PutHandler()
        //DeleteHandler to process DELETE request
	DeleteHandler()
        //GetHandler to process GET request
	GetHandler()
}

//ProfileContext to store profile information
type ProfileContext struct {
	rw http.ResponseWriter
	req *http.Request
	problemDetails *problemdetails.ProblemDetails
	logcontent     *log.LogStruct
        statusCode int
	IsRegister bool
	body []byte
	profileID string
}


//Init profilecontext
func (p *ProfileContext) Init(rw http.ResponseWriter, req *http.Request, sequenceID, profileID string){
	p.rw = rw
	p.req = req
	p.problemDetails = &problemdetails.ProblemDetails{}
	p.logcontent = &log.LogStruct{SequenceId:sequenceID}
        p.profileID = profileID

}

//GetStatusCode get response status code
func (p *ProfileContext) GetStatusCode() int {
	return p.statusCode
}

//GetLogContent get logcontent
func (p *ProfileContext) GetLogContent() *log.LogStruct{
	return  p.logcontent
}

//GetProblemDetails get problemdetails
func (p *ProfileContext) GetProblemDetails() *problemdetails.ProblemDetails {
	return p.problemDetails
}

//GetProfile to get profile after process
func (p *ProfileContext) GetProfile() []byte{
	return p.body

}


//GenerateID for genereate the provision ID
func GenerateID(content []byte) string {
	formatContent, err := utils.JsonFormatter(content)
	if err == nil {
		md5value := md5.Sum(formatContent)
		md5str := fmt.Sprintf("%x", md5value)

		return md5str
	}

	return ""
}


