package dbmgmt

import (
	"com/dbproxy"
	"com/dbproxy/nfmessage/nfprofile"
	"fmt"
	"strings"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt/mock_dbproxy"
	"github.com/golang/mock/gomock"
)

func init() {
	log.SetLevel(log.FatalLevel)
}

//var (
//dbclient dbproxy.NFDataManagementServiceClient
//)

/*func TestPut(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool = make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}

	//create request nfMsg
	profileIndex := &nfprofile.NFProfileIndex{
		NfType:   "NfType",
		NfStatus: "expired",
		Plmn:     "plmn",
	}

	profilePutRequest := &nfprofile.NFProfilePutRequest{NfInstanceId: "123456", Index: profileIndex, RawNfProfile: []byte("body")}
	nfProfilePutRequest := &dbproxy.PutRequest_NfProfilePutRequest{
		NfProfilePutRequest: profilePutRequest,
	}
	putReqData := &dbproxy.PutRequest{
		Data: nfProfilePutRequest,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("grpc error"))
	_, err := Put(putReqData)
	if err == nil {
		t.Fatalf("should get grpc error!")
	}

	//create protocol error nfMsg
	nfProtocolError := &dbproxy.NFProtocolError{Code: 4000}
	protocolError := &dbproxy.NFMessage_ProtocolError{
		ProtocolError: nfProtocolError,
	}

	nfMsgErr := &dbproxy.NFMessage{
		Data: protocolError,
	}
	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgErr, nil)
	_, err1 := Put(putReqData)

	isContain := strings.Contains(err1.Error(), "4000")
	if false == isContain {
		t.Fatalf("should reurnt protocol error code 4000!")
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err2 := Put(putReqData)
	if err2 == nil {
		t.Fatalf("should return err info: db-proxy return invalid messages type")
	}

	profilePutResponse := &nfprofile.NFProfilePutResponse{Code: 2000}
	nfProfilePutResponse := &dbproxy.PutResponse_NfProfilePutResponse{
		NfProfilePutResponse: profilePutResponse,
	}
	putRespData := &dbproxy.PutResponse{
		Data: nfProfilePutResponse,
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
	ret, err3 := Put(putReqData)
	if err3 == nil {
		if (ret.GetNfProfilePutResponse().GetCode()) != 2000 {
			t.Fatalf("NfProfilePutResponse should return code = 2000")
		}
	} else {
		t.Fatalf(" should return NfProfilePutResponse correctly")
	}

	getResponse := &dbproxy.NFResponse_GetResponse{}
	respData = &dbproxy.NFResponse{
		Data: getResponse,
	}
	response = &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp = &dbproxy.NFMessage{
		Data: response,
	}
	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)
	ret, err4 := Put(putReqData)
	if err4 == nil {
		t.Fatalf(" should return err info: db-proxy return invalid respone type")
	}
}

func TestGet(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool = make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}

	//create request nfMsg
	key := &nfprofile.NFProfileKey{
		NfInstanceId: "123456",
	}
	nfProfileGetRequestKey := &nfprofile.NFProfileGetRequest_Key{Key: key}
	profileGetRequest := &nfprofile.NFProfileGetRequest{Data: nfProfileGetRequestKey}
	nfProfileGetRequest := &dbproxy.GetRequest_NfProfileGetRequest{
		NfProfileGetRequest: profileGetRequest,
	}
	getReqData := &dbproxy.GetRequest{
		Data: nfProfileGetRequest,
	}

	//create response nfMsg
	nfProtocolError := &dbproxy.NFProtocolError{Code: 4000}
	protocolError := &dbproxy.NFMessage_ProtocolError{
		ProtocolError: nfProtocolError,
	}

	nfMsgErr := &dbproxy.NFMessage{
		Data: protocolError,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("grpc error"))
	_, err := Get(getReqData)
	if err == nil {
		t.Fatalf("should get grpc error!")
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgErr, nil)
	_, err1 := Get(getReqData)

	isContain := strings.Contains(err1.Error(), "4000")
	if false == isContain {
		t.Fatalf("should reurnt protocol error code 4000!")
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err2 := Get(getReqData)
	if err2 == nil {
		t.Fatalf("should return err info: db-proxy return invalid messages type")
	}

	profileGetResponse := &nfprofile.NFProfileGetResponse{Code: 2000}
	nfProfileGetResponse := &dbproxy.GetResponse_NfProfileGetResponse{
		NfProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: nfProfileGetResponse,
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
	ret, err3 := Get(getReqData)
	if err3 == nil {
		if (ret.GetNfProfileGetResponse().GetCode()) != 2000 {
			t.Fatalf("NfProfileGetResponse should return code = 2000")
		}
	} else {
		t.Fatalf(" should return NfProfileGetResponse correctly")
	}

	putResponse := &dbproxy.NFResponse_PutResponse{}
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
	ret, err4 := Get(getReqData)
	if err4 == nil {
		t.Fatalf(" should return err info: db-proxy return invalid respone type")
	}
}*/

func init() {
	poolLength = 10
}

func TestDelete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool = make([]dbproxy.NFDataManagementServiceClient, poolLength)
	for i := 0; i < poolLength; i++ {
		dbclientPool[i] = mockDbClient
	}

	//create request nfMsg
	profileDelRequest := &nfprofile.NFProfileDelRequest{NfInstanceId: "123456"}
	nfProfileDelRequest := &dbproxy.DelRequest_NfProfileDelRequest{
		NfProfileDelRequest: profileDelRequest,
	}
	delReqData := &dbproxy.DelRequest{
		Data: nfProfileDelRequest,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("grpc error"))
	_, err := Delete(delReqData)
	if err == nil {
		t.Fatalf("should get grpc error!")
	}

	//create protocol error nfMsg
	nfProtocolError := &dbproxy.NFProtocolError{Code: 4000}
	protocolError := &dbproxy.NFMessage_ProtocolError{
		ProtocolError: nfProtocolError,
	}
	nfMsgErr := &dbproxy.NFMessage{
		Data: protocolError,
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgErr, nil)
	_, err1 := Delete(delReqData)

	isContain := strings.Contains(err1.Error(), "4000")
	if !isContain {
		t.Fatalf("should reurnt protocol error code 4000!")
	}

	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err2 := Delete(delReqData)
	if err2 == nil {
		t.Fatalf("should return err info: db-proxy return invalid messages type")
	}

	profileDetResponse := &nfprofile.NFProfileDelResponse{Code: 2000}
	nfProfileDetResponse := &dbproxy.DelResponse_NfProfileDelResponse{
		NfProfileDelResponse: profileDetResponse,
	}
	delRespData := &dbproxy.DelResponse{
		Data: nfProfileDetResponse,
	}
	delResponse := &dbproxy.NFResponse_DelResponse{
		DelResponse: delRespData,
	}
	respData := &dbproxy.NFResponse{
		Data: delResponse,
	}
	response := &dbproxy.NFMessage_Response{
		Response: respData,
	}
	nfMsgResp := &dbproxy.NFMessage{
		Data: response,
	}
	mockDbClient.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nfMsgResp, nil)
	ret, err3 := Delete(delReqData)
	if err3 == nil {
		if (ret.GetNfProfileDelResponse().GetCode()) != 2000 {
			t.Fatalf("NfProfileDelResponse should return code = 2000")
		}
	} else {
		t.Fatalf(" should return NfProfileDelResponse correctly")
	}

	putResponse := &dbproxy.NFResponse_PutResponse{}
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
	_, err4 := Delete(delReqData)
	if err4 == nil {
		t.Fatalf(" should return err info: db-proxy return invalid respone type")
	}
}
