package provprofile

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt/mock_dbproxy"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"net/http/httptest"
	"com/dbproxy/nfmessage/groupprofile"
	"com/dbproxy"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
	"bytes"
)

var groupProfile = []byte(`{
"nfType": ["UDM","AUSF"],
"supiRanges": [
{
    "pattern":"^imsi-12345\\d{4}$"
}],
"groupId": "shanghai"
}
`)
func TestNrfProvGroupProfilesPostHandler(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/nnrf-prov/v1/group-profile", bytes.NewBuffer(groupProfile))
	req.Header.Set("Content-Type", "application/json")

	profilePutResponse := &groupprofile.GroupProfilePutResponse{Code: 2001}
	groupProfilePutResponse := &dbproxy.PutResponse_GroupProfilePutResponse{
		GroupProfilePutResponse: profilePutResponse,
	}
	putRespData := &dbproxy.PutResponse{
		Data: groupProfilePutResponse,
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
	profileHandler := &ImsiProfilHandler{}
	profileHandler.Init(resp, req, "", "", groupProfile, 0)
	profileHandler.PostHandler()
	if profileHandler.context.statusCode != http.StatusCreated {
		t.Errorf("TestNrfProvGroupProfilesPostHandler: NrfProvGroupProfilesPostHandler response code %d check fail", resp.Code)
	}
}

func TestNrfProvGroupProfilePutHandler(t *testing.T) {
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

	req := httptest.NewRequest("PUT", "/nnrf-prov/v1/group-profile?groupProfileID=12345678-provision-GroupSupi-0002", bytes.NewBuffer(groupProfile))
	req.Header.Set("Content-Type", "application/json")

	groupProfileInfo := &groupprofile.GroupProfileInfo{GroupProfileId: "profileId", SupiVersion: 0, GroupProfileData: groupProfile}
	groupProfileSet := []*groupprofile.GroupProfileInfo{groupProfileInfo}
	profileGetResponse := &groupprofile.GroupProfileGetResponse{Code: 2000, GroupProfileInfo: groupProfileSet}
	groupProfileGetResponse := &dbproxy.GetResponse_GroupProfileGetResponse{
		GroupProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: groupProfileGetResponse,
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

	profilePutResponse := &groupprofile.GroupProfilePutResponse{Code: 2001}
	groupProfilePutResponse := &dbproxy.PutResponse_GroupProfilePutResponse{
		GroupProfilePutResponse: profilePutResponse,
	}
	putRespData := &dbproxy.PutResponse{
		Data: groupProfilePutResponse,
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

	profileHandler := &ImsiProfilHandler{}
	profileHandler.Init(resp, req, "", "", groupProfile, 0)
	profileHandler.PutHandler()
	if profileHandler.context.statusCode != http.StatusOK {
		t.Errorf("TestNrfProvGroupProfilePutHandler: NrfProvGroupProfilePutHandler response code %d check fail", resp.Code)
	}
}

func TestNrfProvGroupProfileDeleteHandler(t *testing.T) {
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

	req := httptest.NewRequest("DELETE", "/nnrf-prov/v1/group-profile?groupProfileID=12345678-provision-GroupSupi-0002", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	groupProfileInfo := &groupprofile.GroupProfileInfo{GroupProfileId: "profileId", SupiVersion: 0, GroupProfileData: groupProfile}
	groupProfileSet := []*groupprofile.GroupProfileInfo{groupProfileInfo}
	profileGetResponse := &groupprofile.GroupProfileGetResponse{Code: 2000, GroupProfileInfo: groupProfileSet}
	groupProfileGetResponse := &dbproxy.GetResponse_GroupProfileGetResponse{
		GroupProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: groupProfileGetResponse,
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
	profileDelResponse := &groupprofile.GroupProfileDelResponse{Code: 2000}
	groupProfileDelResponse := &dbproxy.DelResponse_GroupProfileDelResponse{
		GroupProfileDelResponse: profileDelResponse,
	}
	delRespData := &dbproxy.DelResponse{
		Data: groupProfileDelResponse,
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

	profileHandler := &ImsiProfilHandler{}
	profileHandler.Init(resp, req, "", "", []byte(""), 0)
	profileHandler.DeleteHandler()
	if profileHandler.context.statusCode != http.StatusNoContent {
		t.Errorf("TestNrfProvGroupProfileDeleteHandler: NrfProvGroupProfileDeleteHandler response code %d check fail", resp.Code)
	}
}
