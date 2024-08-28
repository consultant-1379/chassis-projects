package nfrequester

import (
	"net/url"
	"testing"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

func TestValidateStringTypeForMadatory(t *testing.T) {
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateStringTypeForMadatory(consts.SearchDataTargetNfType)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataTargetNfType) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&target-nf-type=AMF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateStringTypeForMadatory(consts.SearchDataTargetNfType)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataTargetNfType) {
		t.Fatal("Should return not nil, but did!")
	}
}

func TestValidateStringTypeForOptional(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&external-group-identity=aabbccdd-000-00-aa&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateStringTypeForOptional(consts.SearchDataExterGroupID, true)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataExterGroupID) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&target-nf-type=AMF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	_ = queryform1.validateStringTypeForOptional(consts.SearchDataExterGroupID, true)
	if queryform1.GetExistFlag(consts.SearchDataExterGroupID) {
		t.Fatal("Should return not nil, but did!")
	}
}

func TestValidatePlmnType(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&chf-supported-plmn={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validatePlmnType(consts.SearchDataChfSupportedPlmn)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataChfSupportedPlmn) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&chf-supported-plmn={"mcc":"460", "mnc":"0"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validatePlmnType(consts.SearchDataChfSupportedPlmn)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataChfSupportedPlmn) {
		t.Fatal("Should not return nil, but did!")
	}
}

func TestValidatePlmnListType(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"}&target-plmn-list={"mcc":"460"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validatePlmnListType(consts.SearchDataRequesterPlmnList)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err := queryform0.validatePlmnListType(consts.SearchDataTargetPlmnList)
	if err == nil || queryform0.GetExistFlag(consts.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"0"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validatePlmnListType(consts.SearchDataRequesterPlmnList)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataRequesterPlmnList) {
		t.Fatal("Should not return nil, but did!")
	}

	var queryform2 SearchParameterData
	oriqueryForm2, err4 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list=[{"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"11"}]&target-plmn-list={"mcc":"460"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list=[{"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"22"}]`)
	if err4 != nil {
		t.Fatal("url parse error")
	}
	queryform2.InitMember(oriqueryForm2)
	err5 := queryform2.validatePlmnListType(consts.SearchDataRequesterPlmnList)
	if err5 != nil || !queryform2.GetExistFlag(consts.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err6 := queryform2.validatePlmnListType(consts.SearchDataTargetPlmnList)
	if err6 == nil || queryform2.GetExistFlag(consts.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}

	var queryform3 SearchParameterData
	oriqueryForm3, err7 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"11"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"22"}`)
	if err7 != nil {
		t.Fatal("url parse error")
	}
	queryform3.InitMember(oriqueryForm3)
	err8 := queryform3.validatePlmnListType(consts.SearchDataRequesterPlmnList)
	if err8 == nil || queryform3.GetExistFlag(consts.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err9 := queryform3.validatePlmnListType(consts.SearchDataTargetPlmnList)
	if err9 == nil || queryform3.GetExistFlag(consts.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}
}

func TestGetNRFDiscPlmnValue(t *testing.T) {
	var DiscPara SearchParameterData
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"12\"}")
	DiscPara.SetFlag(consts.SearchDataChfSupportedPlmn, true)
	DiscPara.SetValue(consts.SearchDataChfSupportedPlmn, plmnList)
	targetPlmn, _ := DiscPara.GetNRFDiscPlmnValue(consts.SearchDataChfSupportedPlmn)
	if targetPlmn != "46012" {
		t.Fatal("func GetNRFDiscPlmnValue() targetPlmn should be 46012, but not")
	}
}

func TestGetNRFDiscPlmnListValue(t *testing.T) {
	var DiscPara SearchParameterData
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"12\"}")
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	DiscPara.SetFlag(consts.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(consts.SearchDataTargetPlmnList, plmnList)
	targetPlmnList := DiscPara.GetNRFDiscPlmnListValue(consts.SearchDataTargetPlmnList)
	if len(targetPlmnList) != 2 {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmn length should be 2, but not")
	}
	if targetPlmnList[0] != "46012" && targetPlmnList[1] != "46000" {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmnList should match, but not")
	}

	var plmnList2 []string
	plmnList2 = append(plmnList2, "[{\"mcc\":\"460\", \"mnc\":\"12\"},{\"mcc\":\"222\", \"mnc\":\"22\"}]")
	plmnList2 = append(plmnList2, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	DiscPara.SetFlag(consts.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(consts.SearchDataTargetPlmnList, plmnList2)
	targetPlmnList2 := DiscPara.GetNRFDiscPlmnListValue(consts.SearchDataTargetPlmnList)
	if len(targetPlmnList2) != 3 {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmn length should be 3, but not")
	}
	if targetPlmnList[0] != "46012" && targetPlmnList[1] != "22222" && targetPlmnList[2] != "46000" {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmnList should match, but not")
	}
}

func TestValidateListSnssais(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "222222"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListSnssais(consts.SearchDataSnssais)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataSnssais) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":256, "sd": "222222"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListSnssais(consts.SearchDataSnssais)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataSnssais) {
		t.Fatal("Should not return nil, but did!")
	}

	var queryform2 SearchParameterData
	oriqueryForm2, err4 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":1, "sd": "111111"}&snssais=[{"sst":2, "sd": "222222"},{"sst":3, "sd": "333333"}]`)
	if err4 != nil {
		t.Fatal("url parse error")
	}
	queryform2.InitMember(oriqueryForm2)
	err5 := queryform2.validateListSnssais(consts.SearchDataSnssais)
	if err5 != nil || !queryform2.GetExistFlag(consts.SearchDataSnssais) {
		t.Fatal("Should return nil, but not!")
	}
}

func TestValidateServiceName(t *testing.T) {
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(consts.SearchDataServiceName)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataServiceName) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListStringType(consts.SearchDataServiceName)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataServiceName) {
		t.Fatal("Should not return nil, but did!")
	}

	var queryform2 SearchParameterData
	oriqueryForm2, err3 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err3 != nil {
		t.Fatal("url parse error")
	}
	//queryform2.value = oriqueryForm2
	queryform2.InitMember(oriqueryForm2)
	err4 := queryform2.validateListStringType(consts.SearchDataServiceName)
	if err4 != nil || queryform2.GetExistFlag(consts.SearchDataServiceName) {
		t.Fatal("ServiceNames is option, Should return nil!")
	}
}

func TestValidateNsiList(t *testing.T) {
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}&nsi-list=222222`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(consts.SearchDataNsiList)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataNsiList) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":256, "sd": "2"}&nsi-list=`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListStringType(consts.SearchDataNsiList)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataNsiList) {
		t.Fatal("Should not return nil, but did!")
	}
}

func TestGetNRFDiscStringValue(t *testing.T) {
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateGpsiType(consts.SearchDataGpsi)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataGpsi) {
		t.Fatal("Should return nil, but not!")
	}
	gpsi := queryform0.getNRFDiscStringValue(consts.SearchDataGpsi)
	if gpsi != "msisdn-423456789050000" {
		t.Fatal("gpsi should matched, but failed")
	}
}

func TestGetNRFDiscListString(t *testing.T) {
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	var resultServiceNames []string
	resultServiceNames = append(resultServiceNames, "nudr-uecm")
	resultServiceNames = append(resultServiceNames, "nudr-usdm")
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(consts.SearchDataServiceName)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataServiceName) {
		t.Fatal("Should return nil, but not!")
	}
	serviceNames := queryform0.getNRFDiscListString(consts.SearchDataServiceName)
	if len(serviceNames) != len(resultServiceNames) {
		t.Fatal("serviceNames should matched, but failed")
	}
	for i, v := range serviceNames {
		if v != resultServiceNames[i] {
			t.Fatal("serviceNames should matched, but failed")
		}
	}
}
func TestValidateNRFDiscovery(t *testing.T) {
	var DiscPara SearchParameterData

	value := make(map[string][]string)
	DiscPara.InitMember(value)

	err := DiscPara.ValidateNRFDiscovery()
	if err == nil {
		t.Fatal("func ValidateNRFDiscovery() should validate fail, but success")
	}
	var nfTypeArray []string
	nfTypeArray = append(nfTypeArray, "UDR")
	DiscPara.flag[consts.SearchDataTargetNfType] = true
	DiscPara.value[consts.SearchDataTargetNfType] = nfTypeArray

	err1 := DiscPara.ValidateNRFDiscovery()
	if err1 == nil {
		t.Fatal("func ValidateNRFDiscovery() should validate fail, but success")
	}
	var requesterNfTypeArray []string
	requesterNfTypeArray = append(requesterNfTypeArray, "UDM")
	DiscPara.flag[consts.SearchDataRequesterNfType] = true
	DiscPara.value[consts.SearchDataRequesterNfType] = requesterNfTypeArray

	err2 := DiscPara.ValidateNRFDiscovery()
	if err2 != nil {
		t.Fatal("func ValidateNRFDiscovery() should validate success, but fail")
	}

	var serviceNameList []string
	serviceNameList = append(serviceNameList, "service1")
	DiscPara.value[consts.SearchDataServiceName] = serviceNameList
	DiscPara.flag[consts.SearchDataServiceName] = true

	var pgwIndList []string
	pgwIndList = append(pgwIndList, "true")
	DiscPara.value[consts.SearchDataPGWInd] = pgwIndList
	DiscPara.flag[consts.SearchDataPGWInd] = true

	var groupIdList []string
	groupIdList = append(groupIdList, "123")
	DiscPara.value[consts.SearchDataGroupIDList] = groupIdList
	DiscPara.flag[consts.SearchDataGroupIDList] = true

	var requesterNfInstFQDNList []string
	requesterNfInstFQDNList = append(requesterNfInstFQDNList, "123")
	DiscPara.value[consts.SearchDataRequesterNFInstFQDN] = requesterNfInstFQDNList
	DiscPara.flag[consts.SearchDataRequesterNFInstFQDN] = true

	var targetNfFQDNList []string
	targetNfFQDNList = append(targetNfFQDNList, "123")
	DiscPara.value[consts.SearchDataTargetNFFQDN] = targetNfFQDNList
	DiscPara.flag[consts.SearchDataTargetNFFQDN] = true

	var dnnList []string
	dnnList = append(dnnList, "123")
	DiscPara.value[consts.SearchDataDnn] = dnnList
	DiscPara.flag[consts.SearchDataDnn] = true
	err3 := DiscPara.ValidateNRFDiscovery()
	if err3 != nil {
		t.Fatal("func ValidateNRFDiscovery() should validate success, but fail")
	}
}

func TestGetNRFDiscListSnssais(t *testing.T) {
	var DiscPara SearchParameterData
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var snssaisList []string
	snssaisList = append(snssaisList, "{\"sst\":1,\"sd\":\"1\"}")
	snssaisList = append(snssaisList, "[{\"sst\":2,\"sd\":\"2\"}]")
	DiscPara.SetFlag(consts.SearchDataSnssais, true)
	DiscPara.SetValue(consts.SearchDataSnssais, snssaisList)

	snssais := DiscPara.GetNRFDiscListSnssais(consts.SearchDataSnssais)
	if snssais != "[{\"sst\": 1,\"sd\": \"1\"},{\"sst\": 2,\"sd\": \"2\"}]" {
		t.Fatal("func GetNRFDiscListSnssais() snssais should matched, but fail")
	}
}

func TestValidatAccessType(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateAccessType(consts.SearchDataAccessType)
	if err1 == nil || queryform0.GetExistFlag(consts.SearchDataAccessType) {
		t.Fatal("access-type validate should fail, but success!")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err1 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP_ACCESS`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	err2 := queryform1.validateAccessType(consts.SearchDataAccessType)
	if err2 != nil || !queryform1.GetExistFlag(consts.SearchDataAccessType) {
		t.Fatal("access-type validate should success, but fail!")
	}
}

func TestValidateIsSupportedParam(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	params, err1 := queryform0.validateIsSupportedParam()
	if err1 != nil || params != "" {
		t.Fatal("all parameters should be supported, but not!")
	}
	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP&test=test`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	params1, err3 := queryform1.validateIsSupportedParam()
	if err3 == nil || params1 == "" {
		t.Fatal("parameters has unsupported, should validate fail, but success!")
	}
}

func TestValidateGuamiType(t *testing.T) {
	util.PreComplieRegexp()
	var queryform0 SearchParameterData
	oriqueryForm0, err0 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&guami={"plmnId":{"mcc":"460","mnc":"000"},"amfId":"AAff00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateGuamiType(consts.SearchDataGuami)
	if err1 != nil || !queryform0.GetExistFlag(consts.SearchDataGuami) {
		t.Fatal("guami should validated true, but false")
	}

	var queryform1 SearchParameterData
	oriqueryForm1, err2 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&guami={"plmnId":{"mcc":"460","mnc":"000"},"amfId":"AAff001"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateGuamiType(consts.SearchDataGuami)
	if err3 == nil || queryform1.GetExistFlag(consts.SearchDataGuami) {
		t.Fatal("guami should validated false, but true")
	}
}
