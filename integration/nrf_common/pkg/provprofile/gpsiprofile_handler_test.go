package provprofile

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"github.com/golang/mock/gomock"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt/mock_dbproxy"
	"com/dbproxy"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"net/http/httptest"
	"net/http"
	"bytes"
	"testing"
	"com/dbproxy/nfmessage/gpsiprofile"
)


var gpsiProfile = []byte(`{
"nfType": ["UDM","UDR"],
"gpsiRanges": [
{
    "pattern":"^msisdn-12345\\d{4}$"
}],
"groupId": "shanghai"
}
`)
func TestProvGpsiProfilesPostHandler(t *testing.T) {
	PreComplieRegexp()
	dbmgmt.InitDB("5000")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool := make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}
	dbmgmt.SetDbclientPool(dbclientPool)

	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("POST", "/nnrf-prov/v1/group-profiles/gpsi-groups", bytes.NewBuffer(gpsiProfile))
	req.Header.Set("Content-Type", "application/json")

	profilePutResponse := &gpsiprofile.GpsiProfilePutResponse{Code: 2001}
	gpsiProfilePutResponse := &dbproxy.PutResponse_GpsiProfilePutResponse{
		GpsiProfilePutResponse: profilePutResponse,
	}
	putRespData := &dbproxy.PutResponse{
		Data: gpsiProfilePutResponse,
	}
	putResponse := &dbproxy.NFResponse_PutResponse{
		PutResponse: putRespData,
	}
	respData := &dbproxy.NFResponse{
		Data: putResponse,
	}
	response := &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp := &dbproxy.NFMessage{
		Data: response,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)
	profileHandler := &GpsiProfileHandler{}
	profileHandler.Init(resp, req,  "", "", gpsiProfile, 0)
	profileHandler.PostHandler()
	if profileHandler.context.statusCode != http.StatusCreated {
		t.Errorf("TestNrfProvGpsiProfilesPostHandler: NrfProvGpsiProfilesPostHandler response code %d check fail", resp.Code)
	}
}

func TestProvGpsiProfilePutHandler(t *testing.T) {
	PreComplieRegexp()
	dbmgmt.InitDB("5000")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool := make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}
	dbmgmt.SetDbclientPool(dbclientPool)

	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("PUT", "/nnrf-prov/v1/group-profiles/gpsi-groups?groupProfileID=12345678-provision-Gpsi-0002", bytes.NewBuffer(gpsiProfile))
	req.Header.Set("Content-Type", "application/json")

	gpsiProfileInfo := &gpsiprofile.GpsiProfileInfo{GpsiProfileId: "profileId", GpsiVersion: 0, GpsiProfileData: gpsiProfile}
	gpsiProfileSet := []*gpsiprofile.GpsiProfileInfo{gpsiProfileInfo}
	profileGetResponse := &gpsiprofile.GpsiProfileGetResponse{Code: 2000, GpsiProfileInfo: gpsiProfileSet}
	gpsiProfileGetResponse := &dbproxy.GetResponse_GpsiProfileGetResponse{
		GpsiProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: gpsiProfileGetResponse,
	}
	getResponse := &dbproxy.NFResponse_GetResponse{
		GetResponse: getRespData,
	}
	respData := &dbproxy.NFResponse{
		Data: getResponse,
	}
	response := &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp := &dbproxy.NFMessage{
		Data: response,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)

	profilePutResponse := &gpsiprofile.GpsiProfilePutResponse{Code: 2001}
	gpsiProfilePutResponse := &dbproxy.PutResponse_GpsiProfilePutResponse{
		GpsiProfilePutResponse: profilePutResponse,
	}
	putRespData := &dbproxy.PutResponse{
		Data: gpsiProfilePutResponse,
	}
	putResponse := &dbproxy.NFResponse_PutResponse{
		PutResponse: putRespData,
	}
	respData = &dbproxy.NFResponse{
		Data: putResponse,
	}
	response = &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp = &dbproxy.NFMessage{
		Data: response,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)
	profileHandler := &GpsiProfileHandler{}
	profileHandler.Init(resp, req, "", "", gpsiProfile, 0)
	profileHandler.PutHandler()
	if profileHandler.context.statusCode != http.StatusOK {
		t.Errorf("TestNrfProvGpsiProfilePutHandler: NrfProvGpsiProfilePutHandler response code %d check fail", resp.Code)
	}
}

func TestProvGpsiProfileDeleteHandler(t *testing.T) {
	PreComplieRegexp()
	dbmgmt.InitDB("5000")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool := make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}
	dbmgmt.SetDbclientPool(dbclientPool)

	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("DELETE", "/nnrf-prov/v1/group-profiles/gpsi-groups?groupProfileID=12345678-provision-Gpsi-0002", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	gpsiProfileInfo := &gpsiprofile.GpsiProfileInfo{GpsiProfileId: "profileId", GpsiVersion: 0, GpsiProfileData: gpsiProfile}
	gpsiProfileSet := []*gpsiprofile.GpsiProfileInfo{gpsiProfileInfo}
	profileGetResponse := &gpsiprofile.GpsiProfileGetResponse{Code: 2000, GpsiProfileInfo: gpsiProfileSet}
	gpsiProfileGetResponse := &dbproxy.GetResponse_GpsiProfileGetResponse{
		GpsiProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: gpsiProfileGetResponse,
	}
	getResponse := &dbproxy.NFResponse_GetResponse{
		GetResponse: getRespData,
	}
	respData := &dbproxy.NFResponse{
		Data: getResponse,
	}
	response := &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp := &dbproxy.NFMessage{
		Data: response,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)
	profileDelResponse := &gpsiprofile.GpsiProfileDelResponse{Code: 2000}
	gpsiProfileDelResponse := &dbproxy.DelResponse_GpsiProfileDelResponse{
		GpsiProfileDelResponse: profileDelResponse,
	}
	delRespData := &dbproxy.DelResponse{
		Data: gpsiProfileDelResponse,
	}
	delResponse := &dbproxy.NFResponse_DelResponse{
		DelResponse: delRespData,
	}
	respDelData := &dbproxy.NFResponse{
		Data: delResponse,
	}
	responseDel := &dbproxy.NFMessage_Response{
		Response: respDelData,
	}
	nfMsgDelResp := &dbproxy.NFMessage{
		Data: responseDel,
	}
	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgDelResp, nil)

	profileHandler := &GpsiProfileHandler{}
	profileHandler.Init(resp, req,  "", "", []byte(""), 0)
	profileHandler.DeleteHandler()
	if profileHandler.context.statusCode != http.StatusNoContent {
		t.Errorf("TestNrfProvGpsiProfileDeleteHandler: NrfProvGpsiProfileDeleteHandler response code %d check fail", resp.Code)
	}
}
