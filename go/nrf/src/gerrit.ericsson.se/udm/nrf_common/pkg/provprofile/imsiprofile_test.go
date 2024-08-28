package provprofile

import (
	"com/dbproxy"
	"com/dbproxy/nfmessage/imsiprefixprofile"
	"net/http"
	"net/http/httptest"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt/mock_dbproxy"
	"github.com/golang/mock/gomock"
)

func TestGetImsiProfile(t *testing.T) {
	dbmgmt.InitDB("5000")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbClient := mock_dbproxy.NewMockNFDataManagementServiceClient(ctrl)
	dbclientPool := make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	for i := 0; i < internalconf.DbproxyConnectionNum; i++ {
		dbclientPool[i] = mockDbClient
	}
	dbmgmt.SetDbclientPool(dbclientPool)

	imsi := "450008912100000"

	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	var nfTypeStr = "UDM" + NfTypeSeparator + "AUSF"
	valueInfoSet := []string{}
	valueInfoSet = append(valueInfoSet, "15"+ValueInfoSeparator+PrefixTypeGroupID+ValueInfoSeparator+"gid1"+ValueInfoSeparator+"UDM")
	valueInfoSet = append(valueInfoSet, "15"+ValueInfoSeparator+PrefixTypeGroupID+ValueInfoSeparator+"gid2"+ValueInfoSeparator+nfTypeStr)
	profileGetResponse := &imsiprefixprofile.ImsiprefixProfileGetResponse{Code: 2000, ValueInfo: valueInfoSet}
	imsiprefixProfileGetResponse := &dbproxy.GetResponse_ImsiprefixProfileGetResponse{
		ImsiprefixProfileGetResponse: profileGetResponse,
	}
	getRespData := &dbproxy.GetResponse{
		Data: imsiprefixProfileGetResponse,
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
	imsiSearchResultList := []ImsiSearchResult{}
	code, err := GetImsiProfile(imsi, &imsiSearchResultList)
	if nil != err && code != dbmgmt.DbGetSuccess {
		t.Errorf("TestGetImsiProfile: response error: %s", err.Error())
	}
	if imsiSearchResultList[0].ValueID != "gid1" || imsiSearchResultList[1].ValueID != "gid2" || imsiSearchResultList[0].ValueType != "gid" || imsiSearchResultList[1].ValueType != "gid" {
		t.Errorf("TestGetImsiProfile: check imsi search result failure. imsiSearchResultList: %v", imsiSearchResultList)
	}
}
