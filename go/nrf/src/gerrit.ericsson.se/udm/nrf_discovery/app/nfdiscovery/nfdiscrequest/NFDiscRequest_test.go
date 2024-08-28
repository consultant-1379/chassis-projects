package nfdiscrequest

import (
	"testing"
	"net/url"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

func TestValidateStringTypeForMan(t *testing.T) {
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateStringTypeForMan(constvalue.SearchDataTargetNfType)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataTargetNfType) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&target-nf-type=AMF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateStringTypeForMan(constvalue.SearchDataTargetNfType)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataTargetNfType) {
		t.Fatal("Should return not nil, but did!")
	}
}

func TestValidateStringTypeForOpt(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&external-group-identity=aabbccdd-000-00-aa&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateStringTypeForOpt(constvalue.SearchDataExterGroupID, true)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataExterGroupID) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&target-nf-type=AMF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn={"mcc":"460", "mnc":"00"}&requester-plmn={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	_ = queryform1.validateStringTypeForOpt(constvalue.SearchDataExterGroupID, true)
	if queryform1.GetExistFlag(constvalue.SearchDataExterGroupID) {
		t.Fatal("Should return not nil, but did!")
	}
}

func TestValidatePlmnType(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&chf-supported-plmn={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validatePlmnType(constvalue.SearchDataChfSupportedPlmn)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataChfSupportedPlmn) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&chf-supported-plmn={"mcc":"460", "mnc":"0"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validatePlmnType(constvalue.SearchDataChfSupportedPlmn)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataChfSupportedPlmn) {
		t.Fatal("Should not return nil, but did!")
	}
}

func TestValidatePlmnListType(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"}&target-plmn-list={"mcc":"460"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validatePlmnListType(constvalue.SearchDataRequesterPlmnList)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err := queryform0.validatePlmnListType(constvalue.SearchDataTargetPlmnList)
	if err == nil || queryform0.GetExistFlag(constvalue.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"0"}&requester-plmn-list={"mcc":"460", "mnc":"00"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validatePlmnListType(constvalue.SearchDataRequesterPlmnList)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataRequesterPlmnList) {
		t.Fatal("Should not return nil, but did!")
	}

	var queryform2 DiscGetPara
	oriqueryForm2, err4 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list=[{"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"11"}]&target-plmn-list={"mcc":"460"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list=[{"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"22"}]`)
	if err4 != nil {
		t.Fatal("url parse error")
	}
	queryform2.InitMember(oriqueryForm2)
	err5 := queryform2.validatePlmnListType(constvalue.SearchDataRequesterPlmnList)
	if err5 != nil || !queryform2.GetExistFlag(constvalue.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err6 := queryform2.validatePlmnListType(constvalue.SearchDataTargetPlmnList)
	if err6 == nil || queryform2.GetExistFlag(constvalue.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}

	var queryform3 DiscGetPara
	oriqueryForm3, err7 := url.ParseQuery(`service-names=nausf-auth&target-nf-type=AUSF&requester-nf-type=AMF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&snssais={"sst":0, "sd": "0"}&target-plmn-list={"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"11"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"},{"mcc":"460", "mnc":"22"}`)
	if err7 != nil {
		t.Fatal("url parse error")
	}
	queryform3.InitMember(oriqueryForm3)
	err8 := queryform3.validatePlmnListType(constvalue.SearchDataRequesterPlmnList)
	if err8 == nil || queryform3.GetExistFlag(constvalue.SearchDataRequesterPlmnList) {
		t.Fatal("Should return nil, but not!")
	}
	err9 := queryform3.validatePlmnListType(constvalue.SearchDataTargetPlmnList)
	if err9 == nil || queryform3.GetExistFlag(constvalue.SearchDataTargetPlmnList) {
		t.Fatal("Should not return nil, but yes!")
	}
}

func TestGetNRFDiscPlmnValue(t *testing.T) {
	var DiscPara DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)


	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"12\"}")
	DiscPara.SetFlag(constvalue.SearchDataChfSupportedPlmn, true)
	DiscPara.SetValue(constvalue.SearchDataChfSupportedPlmn, plmnList)
	targetPlmn, _ := DiscPara.GetNRFDiscPlmnValue(constvalue.SearchDataChfSupportedPlmn)
	if targetPlmn != "46012" {
		t.Fatal("func GetNRFDiscPlmnValue() targetPlmn should be 46012, but not")
	}
}

func TestGetNRFDiscPlmnListValue(t *testing.T) {
	var DiscPara DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)


	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"12\"}")
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	DiscPara.SetFlag(constvalue.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataTargetPlmnList, plmnList)
	targetPlmnList := DiscPara.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
	if len(targetPlmnList) != 2 {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmn length should be 2, but not")
	}
	if targetPlmnList[0] != "46012" && targetPlmnList[1] != "46000" {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmnList should match, but not")
	}

	var plmnList2 []string
	plmnList2 = append(plmnList2, "[{\"mcc\":\"460\", \"mnc\":\"12\"},{\"mcc\":\"222\", \"mnc\":\"22\"}]")
	plmnList2 = append(plmnList2, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	DiscPara.SetFlag(constvalue.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataTargetPlmnList, plmnList2)
	targetPlmnList2 := DiscPara.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
	if len(targetPlmnList2) != 3 {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmn length should be 3, but not")
	}
	if targetPlmnList[0] != "46012" && targetPlmnList[1] != "22222" && targetPlmnList[2] != "46000" {
		t.Fatal("func GetNRFDiscPlmnListValue() targetPlmnList should match, but not")
	}
}

func TestValidateListSnssais(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "222222"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListSnssais(constvalue.SearchDataSnssais)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataSnssais) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":256, "sd": "222222"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListSnssais(constvalue.SearchDataSnssais)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataSnssais) {
		t.Fatal("Should not return nil, but did!")
	}


	var queryform2 DiscGetPara
	oriqueryForm2, err4 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":1, "sd": "111111"}&snssais=[{"sst":2, "sd": "222222"},{"sst":3, "sd": "333333"}]`)
	if err4 != nil {
		t.Fatal("url parse error")
	}
	queryform2.InitMember(oriqueryForm2)
	err5 := queryform2.validateListSnssais(constvalue.SearchDataSnssais)
	if err5 != nil || !queryform2.GetExistFlag(constvalue.SearchDataSnssais) {
		t.Fatal("Should return nil, but not!")
	}
}

func TestValidateServiceName(t *testing.T) {
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(constvalue.SearchDataServiceName)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataServiceName) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListStringType(constvalue.SearchDataServiceName)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataServiceName) {
		t.Fatal("Should not return nil, but did!")
	}

	var queryform2 DiscGetPara
	oriqueryForm2, err3 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err3 != nil {
		t.Fatal("url parse error")
	}
	//queryform2.value = oriqueryForm2
	queryform2.InitMember(oriqueryForm2)
	err4 := queryform2.validateListStringType(constvalue.SearchDataServiceName)
	if err4 != nil || queryform2.GetExistFlag(constvalue.SearchDataServiceName) {
		t.Fatal("ServiceNames is option, Should return nil!")
	}
}

func TestValidateNsiList(t *testing.T) {
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}&nsi-list=222222`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(constvalue.SearchDataNsiList)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataNsiList) {
		t.Fatal("Should return nil, but not!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":256, "sd": "2"}&nsi-list=`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	//queryform1.value = oriqueryForm1
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateListStringType(constvalue.SearchDataNsiList)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataNsiList) {
		t.Fatal("Should not return nil, but did!")
	}
}

func TestGetNRFDiscStringValue(t *testing.T) {
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateGpsiType(constvalue.SearchDataGpsi)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataGpsi) {
		t.Fatal("Should return nil, but not!")
	}
	gpsi := queryform0.getNRFDiscStringValue(constvalue.SearchDataGpsi)
	if gpsi != "msisdn-423456789050000" {
		t.Fatal("gpsi should matched, but failed")
	}
}

func TestGetNRFDiscListString(t *testing.T) {
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=msisdn-523456789041234&routing-indicator=1234&snssais={"sst":2, "sd": "2"}`)
	var resultServiceNames []string
	resultServiceNames = append(resultServiceNames, "nudr-uecm")
	resultServiceNames = append(resultServiceNames, "nudr-usdm")
	if err0 != nil {
		t.Fatal("url parse error")
	}
	//queryform0.value = oriqueryForm0
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateListStringType(constvalue.SearchDataServiceName)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataServiceName) {
		t.Fatal("Should return nil, but not!")
	}
	serviceNames := queryform0.getNRFDiscListString(constvalue.SearchDataServiceName)
	if len(serviceNames) != len(resultServiceNames) {
		t.Fatal("serviceNames should matched, but failed")
	}
	for i, v := range serviceNames {
		if v != resultServiceNames[i] {
			t.Fatal("serviceNames should matched, but failed")
		}
	}
}
func TestValidateNRFDiscovery(t *testing.T)  {
	var DiscPara DiscGetPara

	value := make(map[string][]string)
	DiscPara.InitMember(value)

	err := DiscPara.ValidateNRFDiscovery()
	if err == nil {
		t.Fatal("func ValidateNRFDiscovery() should validate fail, but success")
	}
	var nfTypeArray []string
	nfTypeArray = append(nfTypeArray, "UDR")
	DiscPara.flag[constvalue.SearchDataTargetNfType] = true
	DiscPara.value[constvalue.SearchDataTargetNfType] = nfTypeArray

	err1 := DiscPara.ValidateNRFDiscovery()
	if err1 == nil {
		t.Fatal("func ValidateNRFDiscovery() should validate fail, but success")
	}
	var requesterNfTypeArray []string
	requesterNfTypeArray = append(requesterNfTypeArray, "UDM")
	DiscPara.flag[constvalue.SearchDataRequesterNfType] = true
	DiscPara.value[constvalue.SearchDataRequesterNfType] = requesterNfTypeArray

	err2 := DiscPara.ValidateNRFDiscovery()
	if err2 != nil {
		t.Fatal("func ValidateNRFDiscovery() should validate success, but fail")
	}

	var serviceNameList []string
	serviceNameList = append(serviceNameList, "service1")
	DiscPara.value[constvalue.SearchDataServiceName] = serviceNameList
	DiscPara.flag[constvalue.SearchDataServiceName] = true

	var pgwIndList []string
	pgwIndList = append(pgwIndList, "true")
	DiscPara.value[constvalue.SearchDataPGWInd] = pgwIndList
	DiscPara.flag[constvalue.SearchDataPGWInd] = true

	var groupIdList []string
	groupIdList = append(groupIdList, "123")
	DiscPara.value[constvalue.SearchDataGroupIDList] = groupIdList
	DiscPara.flag[constvalue.SearchDataGroupIDList] = true

	var requesterNfInstFQDNList []string
	requesterNfInstFQDNList = append(requesterNfInstFQDNList, "123")
	DiscPara.value[constvalue.SearchDataRequesterNFInstFQDN] = requesterNfInstFQDNList
	DiscPara.flag[constvalue.SearchDataRequesterNFInstFQDN] = true

	var targetNfFQDNList []string
	targetNfFQDNList = append(targetNfFQDNList, "123")
	DiscPara.value[constvalue.SearchDataTargetNFFQDN] = targetNfFQDNList
	DiscPara.flag[constvalue.SearchDataTargetNFFQDN] = true

	var dnnList []string
	dnnList = append(dnnList, "123")
	DiscPara.value[constvalue.SearchDataDnn] = dnnList
	DiscPara.flag[constvalue.SearchDataDnn] = true
	err3 := DiscPara.ValidateNRFDiscovery()
	if err3 != nil {
		t.Fatal("func ValidateNRFDiscovery() should validate success, but fail")
	}
}

func TestGetNRFDiscListSnssais(t *testing.T) {
	var DiscPara DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var snssaisList []string
	snssaisList = append(snssaisList, "{\"sst\":1,\"sd\":\"1\"}")
	snssaisList = append(snssaisList, "[{\"sst\":2,\"sd\":\"2\"}]")
	DiscPara.SetFlag(constvalue.SearchDataSnssais, true)
	DiscPara.SetValue(constvalue.SearchDataSnssais ,snssaisList)

	snssais := DiscPara.GetNRFDiscListSnssais(constvalue.SearchDataSnssais)
	if snssais != "[{\"sst\": 1,\"sd\": \"1\"},{\"sst\": 2,\"sd\": \"2\"}]" {
		t.Fatal("func GetNRFDiscListSnssais() snssais should matched, but fail")
	}
}

func TestGeneratorCacheKey(t *testing.T)  {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&gpsi=msisdn-423456789050000&external-group-identity=groupid-ABcd1234-000-11-aa&routing-indicator=1234&snssais={"sst":1, "sd": "222222"}&target-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&group-id-list=123&pgw-ind=true`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	queryform0.ValidateNRFDiscovery()
	queryform0.generatorCacheKey()


	var queryform1 DiscGetPara
	oriqueryForm1, err1 := url.ParseQuery(`service-names=nudr-usdm&service-names=nudr-uecm&requester-nf-type=AUSF&target-nf-type=UDM&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&external-group-identity=groupid-ABcd1234-000-11-aa&gpsi=msisdn-423456789050000&routing-indicator=1234&snssais={"sd": "222222", "sst":1}&target-plmn-list={"mcc":"460", "mnc":"00"}&requester-plmn-list={"mcc":"460", "mnc":"00"}&pgw-ind=true&group-id-list=123`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	queryform1.ValidateNRFDiscovery()
	queryform1.generatorCacheKey()

	if queryform0.localCacheKey != queryform1.localCacheKey {
		t.Fatalf("func generatorCacheKey should return same value, but fail")
	}
}

func TestValidatAccessType(t *testing.T)  {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateAccessType(constvalue.SearchDataAccessType)
	if err1 == nil || queryform0.GetExistFlag(constvalue.SearchDataAccessType) {
		t.Fatal("access-type validate should fail, but success!")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err1 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP_ACCESS`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	err2 := queryform1.validateAccessType(constvalue.SearchDataAccessType)
	if err2 != nil || !queryform1.GetExistFlag(constvalue.SearchDataAccessType) {
		t.Fatal("access-type validate should success, but fail!")
	}
}

func TestValidateIsSupportedParam(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&service-names=nudr-usdm&target-nf-type=UDM&requester-nf-type=AUSF&access-type=3GPP`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	params, err1 := queryform0.validateIsSupportedParam()
	if err1 != nil || params != "" {
		t.Fatal("all parameters should be supported, but not!")
	}
	var queryform1 DiscGetPara
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

func TestValidateGuamiType(t *testing.T)  {
	nfdiscutil.PreComplieRegexp()
	var queryform0 DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&guami={"plmnId":{"mcc":"460","mnc":"000"},"amfId":"AAff00"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	err1 := queryform0.validateGuamiType(constvalue.SearchDataGuami)
	if err1 != nil || !queryform0.GetExistFlag(constvalue.SearchDataGuami) {
		t.Fatal("guami should validated true, but false")
	}

	var queryform1 DiscGetPara
	oriqueryForm1, err2 := url.ParseQuery(`target-nf-type=UDM&requester-nf-type=AUSF&guami={"plmnId":{"mcc":"460","mnc":"000"},"amfId":"AAff001"}`)
	if err2 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	err3 := queryform1.validateGuamiType(constvalue.SearchDataGuami)
	if err3 == nil || queryform1.GetExistFlag(constvalue.SearchDataGuami) {
		t.Fatal("guami should validated false, but true")
	}
}
