package nfdiscservice

import (
	"com/dbproxy/nfmessage/nrfprofile"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"testing"
)

func TestMatchByPatternForSupi(t *testing.T) {
	pattern1 := []byte(`
	  {
	    "pattern": "^nai-smartmeter-.+@company\\.com$"
	  }`)

	if !matchByPatternForSupi("nai-smartmeter-11@company.com", []byte(pattern1)) {
		t.Fatalf("should be matched , but failed")
	}

	if matchByPatternForSupi("nai-smartmeter-11@company.comm", []byte(pattern1)) {
		t.Fatalf("should not matched , but matched")
	}
}

func TestMatchbyStartEndForSupi(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	pattern1 := []byte(`
	{
        	"start": "123456789040000",
        	"end": "123456789049999"
      	}`)

	if !matchbyStartEndForSupi("imsi-123456789040001", []byte(pattern1)) {
		t.Fatalf("should be matched , but failed")
	}

	if matchByPatternForSupi("imsi-123456789050000", []byte(pattern1)) {
		t.Fatalf("should not matched , but matched")
	}
}

func TestSortMapByValue(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "7,2"
	m["b"] = "7,1"
	m["c"] = "7,3"
	m["d"] = "5,1"
	m["e"] = "4,1"
	m["f"] = "3,3"
	p := sortMapByValue(m)
	if p[0].Key != "f" || p[1].Key != "e" || p[2].Key != "d" || p[3].Key != "b" || p[4].Key != "a" || p[5].Key != "c" {
		t.Fatal("func sortMapByValue() map should sort by value ase, but not")
	}
	if p.Len() != 6 {
		t.Fatal("func Len() map len should be 3, but not")
	}
	if !p.Less(0, 1) {
		t.Fatal("func Less() index 0 should smaller than 1, but not")
	}
	p.Swap(0, 1)
	if p[0].Key != "e" && p[1].Key != "f" {
		t.Fatal("func Swap() should success, but not")
	}
}

func TestSetAMFFilterForNRFProfile(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	var guamiList []string
	guamiList = append(guamiList, "{\"plmnId\":{\"mcc\":\"460\",\"mnc\":\"000\"},\"amfId\":\"123456\"}")
	value := make(map[string][]string)
	value[constvalue.SearchDataGuami] = guamiList
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataGuami, true)

	var nrfProfileIndex *nrfprofile.NRFProfileIndex
	nrfProfileIndex = &nrfprofile.NRFProfileIndex{}
	setAMFFilterForNRFProfile(DiscPara, nrfProfileIndex)
	if nrfProfileIndex.GetAmfKey1()[0].SubKey1 != "460000" || nrfProfileIndex.GetAmfKey1()[0].SubKey2 != "123456" {
		t.Fatal("func setAMFFilterForNRFProfile() amf guami plmnid shoud be 460000 and amfId should be 123456, but not")
	}

	/*var DiscPara2 nfdiscrequest.DiscGetPara
	var taiList []string
	taiList = append(taiList, "{\"plmnId\":{\"mcc\":\"460\",\"mnc\":\"000\"},\"tac\":\"123456\"}")
	DiscPara2.value = make(map[string][]string)
	DiscPara2.value[constvalue.SearchDataTai] = taiList
	DiscPara2.flag = make(map[string]bool)
	DiscPara2.flag[constvalue.SearchDataTai] = true

	var nrfProfileIndex2 *nrfprofile.NRFProfileIndex
	nrfProfileIndex2 = &nrfprofile.NRFProfileIndex{}
	setAMFFilterForNRFProfile(DiscPara2, nrfProfileIndex2)
	if nrfProfileIndex2.GetAmfKey2()[0].SubKey1 != "460000" || nrfProfileIndex2.GetAmfKey2()[0].SubKey2 != "123456" {
		t.Fatal("func setAMFFilterForNRFProfile() amf tai plmnid shoud be 460000 and tac should be 123456, but not")
	}*/

	var DiscPara3 nfdiscrequest.DiscGetPara
	var amfRegionId []string
	amfRegionId = append(amfRegionId, "123")
	value3 := make(map[string][]string)
	value3[constvalue.SearchDataAmfRegionID] = amfRegionId
	DiscPara3.InitMember(value3)
	DiscPara3.SetFlag(constvalue.SearchDataAmfRegionID, true)

	var nrfProfileIndex3 *nrfprofile.NRFProfileIndex
	nrfProfileIndex3 = &nrfprofile.NRFProfileIndex{}
	setAMFFilterForNRFProfile(DiscPara3, nrfProfileIndex3)
	if nrfProfileIndex3.GetAmfKey3()[0].SubKey1 != "123" {
		t.Fatal("func setAMFFilterForNRFProfile() amf regionId shoud be 123 and, but not")
	}

	var DiscPara4 nfdiscrequest.DiscGetPara
	var amfSetId []string
	amfSetId = append(amfSetId, "1234")
	value4 := make(map[string][]string)
	value4[constvalue.SearchDataAmfSetID] = amfSetId
	DiscPara4.InitMember(value4)
	DiscPara4.SetFlag(constvalue.SearchDataAmfSetID, true)

	var nrfProfileIndex4 *nrfprofile.NRFProfileIndex
	nrfProfileIndex4 = &nrfprofile.NRFProfileIndex{}
	setAMFFilterForNRFProfile(DiscPara4, nrfProfileIndex4)
	if nrfProfileIndex4.GetAmfKey4()[0].SubKey1 != "1234" {
		t.Fatal("func setAMFFilterForNRFProfile() amf setId shoud be 1234 and, but not")
	}
}

func TestSetSMFFilterForNRFProfile(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	var dnnList []string
	dnnList = append(dnnList, "123")
	value := make(map[string][]string)
	value[constvalue.SearchDataDnn] = dnnList
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataDnn, true)

	var nrfProfileIndex *nrfprofile.NRFProfileIndex
	nrfProfileIndex = &nrfprofile.NRFProfileIndex{}
	setSMFFilterForNRFProfile(DiscPara, nrfProfileIndex)
	if nrfProfileIndex.GetSmfKey1()[0].SubKey1 != "123" {
		t.Fatal("func setSMFFilterForNRFProfile() smf dnn shoud be 123, but not")
	}

	var DiscPara2 nfdiscrequest.DiscGetPara
	var pgwList []string
	pgwList = append(pgwList, "1234")
	value2 := make(map[string][]string)
	value2[constvalue.SearchDataPGW] = pgwList
	DiscPara2.InitMember(value2)
	DiscPara2.SetFlag(constvalue.SearchDataPGW, true)

	var nrfProfileIndex2 *nrfprofile.NRFProfileIndex
	nrfProfileIndex2 = &nrfprofile.NRFProfileIndex{}
	setSMFFilterForNRFProfile(DiscPara2, nrfProfileIndex2)
	if nrfProfileIndex2.GetSmfKey2()[0].SubKey1 != "1234" {
		t.Fatal("func setSMFFilterForNRFProfile() smf pgw shoud be 1234, but not")
	}

	/*var DiscPara3 nfdiscrequest.DiscGetPara
	var taiList []string
	taiList = append(taiList, "{\"plmnId\":{\"mcc\":\"460\",\"mnc\":\"000\"},\"tac\":\"123456\"}")
	DiscPara3.value = make(map[string][]string)
	DiscPara3.value[constvalue.SearchDataTai] = taiList
	DiscPara3.flag = make(map[string]bool)
	DiscPara3.flag[constvalue.SearchDataTai] = true

	var nrfProfileIndex3 *nrfprofile.NRFProfileIndex
	nrfProfileIndex3 = &nrfprofile.NRFProfileIndex{}
	setSMFFilterForNRFProfile(DiscPara3, nrfProfileIndex3)
	if nrfProfileIndex3.GetSmfKey3()[0].SubKey1 != "460000" || nrfProfileIndex3.GetSmfKey3()[0].SubKey2 != "123456" {
		t.Fatal("func setSMFFilterForNRFProfile() smf tai plmnid shoud be 460000 and tac should be 123456, but not")
	}*/
}

func TestSetUDMFilterForNRFProfile(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	var routingIndicator []string
	routingIndicator = append(routingIndicator, "123")
	value := make(map[string][]string)
	value[constvalue.SearchDataRoutingIndic] = routingIndicator
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataRoutingIndic, true)

	var nrfProfileIndex *nrfprofile.NRFProfileIndex
	nrfProfileIndex = &nrfprofile.NRFProfileIndex{}
	setUDMFilterForNRFProfile(DiscPara, nrfProfileIndex)
	if nrfProfileIndex.GetUdmKey2()[0].SubKey1 != "123" {
		t.Fatal("func setUDMFilterForNRFProfile() udm routing-indicator shoud be 123, but not")
	}
}

func TestSetAUSFFilterForNRFProfile(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	var routingIndicator []string
	routingIndicator = append(routingIndicator, "123")
	value := make(map[string][]string)
	value[constvalue.SearchDataRoutingIndic] = routingIndicator
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataRoutingIndic, true)

	var nrfProfileIndex *nrfprofile.NRFProfileIndex
	nrfProfileIndex = &nrfprofile.NRFProfileIndex{}
	setAUSFFilterForNRFProfile(DiscPara, nrfProfileIndex)
	if nrfProfileIndex.GetAusfKey2()[0].SubKey1 != "123" {
		t.Fatal("func setAUSFFilterForNRFProfile() ausf routing-indicator shoud be 123, but not")
	}
}

func TestSetPCFFilterForNRFProfile(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	var dnnList []string
	dnnList = append(dnnList, "123")
	value := make(map[string][]string)
	value[constvalue.SearchDataDnn] = dnnList
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataDnn, true)

	var nrfProfileIndex *nrfprofile.NRFProfileIndex
	nrfProfileIndex = &nrfprofile.NRFProfileIndex{}
	setPCFFilterForNRFProfile(DiscPara, nrfProfileIndex)
	if nrfProfileIndex.GetPcfKey1()[0].SubKey1 != "123" {
		t.Fatal("func setPCFFilterForNRFProfile() pcf dnn shoud be 123, but not")
	}
}

func TestIsMatchedGroupIDForNRFProfile(t *testing.T) {
	nfInfo := []byte(`{ "nrfInfo" : {
	"ausfInfoSum": {
		"groupIdList": ["123", "456", "789"]
		"supiRanges": [
			{
				"start": "123456789040000",
				"end": "123456789049999"
			},
			{
				"pattern": "^suci-22345678904\\d{4}$"
			},
			{
				"pattern": "^nai-smartmeter-.+@company\\.com$"
			}
		]
	}
      	}}`)
	nfInfo2 := []byte(`{ "nrfInfo" : {
	"ausfInfoSum": {
		"groupIdList": []
		"supiRanges": [
			{
				"start": "123456789040000",
				"end": "123456789049999"
			},
			{
				"pattern": "^suci-22345678904\\d{4}$"
			},
			{
				"pattern": "^nai-smartmeter-.+@company\\.com$"
			}
		]
	}
      	}}`)
	nfInfo3 := []byte(`{ "nrfInfo" : {
	"ausfInfoSum": {
		"supiRanges": [
			{
				"start": "123456789040000",
				"end": "123456789049999"
			},
			{
				"pattern": "^suci-22345678904\\d{4}$"
			},
			{
				"pattern": "^nai-smartmeter-.+@company\\.com$"
			}
		]
	}
      	}}`)
	var DiscPara1 nfdiscrequest.DiscGetPara
	value := make(map[string][]string)

	var nfTypeArray1 []string
	nfTypeArray1 = append(nfTypeArray1, "AUSF")
	value["target-nf-type"] = nfTypeArray1
	DiscPara1.InitMember(value)
	DiscPara1.SetFlag("target-nf-type", true)
	var groupIdList []string
	groupIdList = append(groupIdList, "123", "345")
	if nfdiscutil.ResultFoundMatch != isMatchedGroupIDForNRFProfile(&DiscPara1, groupIdList, []byte(nfInfo)) {
		t.Fatalf("func isMatchedGroupIDForNRFProfile() should be matched , but failed")
	}
	var groupIdList2 []string
	groupIdList2 = append(groupIdList2, "1234", "4567")
	if nfdiscutil.ResultFoundNotMatch != isMatchedGroupIDForNRFProfile(&DiscPara1, groupIdList2, []byte(nfInfo)) {
		t.Fatalf("func isMatchedGroupIDForNRFProfile() should not be matched , but matched")
	}
	if nfdiscutil.ResultFoundNotMatch != isMatchedGroupIDForNRFProfile(&DiscPara1, groupIdList2, []byte(nfInfo2)) {
		t.Fatalf("func isMatchedGroupIDForNRFProfile() should not be matched , but matched")
	}
	if nfdiscutil.ResultFoundNotMatch != isMatchedGroupIDForNRFProfile(&DiscPara1, groupIdList2, []byte(nfInfo3)) {
		t.Fatalf("func isMatchedGroupIDForNRFProfile() should not be matched , but matched")
	}

	var DiscPara2 nfdiscrequest.DiscGetPara
	value2 := make(map[string][]string)

	var nfTypeArray2 []string
	nfTypeArray2 = append(nfTypeArray2, "AMF")
	value2["target-nf-type"] = nfTypeArray2
	DiscPara2.InitMember(value2)

	if nfdiscutil.ResultError != isMatchedGroupIDForNRFProfile(&DiscPara2, groupIdList, []byte(nfInfo)) {
		t.Fatalf("func isMatchedGroupIDForNRFProfile() return error , but failed")
	}
}

func TestIsMatchedSupiForNRFProfile(t *testing.T) {
	nrfInfo := []byte(`{ "nrfInfo" : {
	"amfInfoSum": {
		"guamiList":[
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"amfId": ""
			}
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		],
		"amfRegionIdList": [
			"",
			""
		],
		"amfSetIdList": [
			"",
			""
		],
	},
	"smfInfoSum": {
		"dnnList": [
			"",
			""
		],
		"pgwFqdnList": [
			"",
			""
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		]
	},
	"udmInfoSum": {
		"supiRanges": [
			{
				"start": "123000000",
				"end": "123999999",
				"pattern": "^imsi-12345\\d{4}$"
			}
		],
		"externalGroupIdentifiersRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"gpsiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"groupIdList": [
			"",
			""
		],
		"routingIndicatorList": [
			"",
			""
		]
	},
	"ausfInfoSum": {
		"groupIdList":[
			"",
			""
		],
		"routingIndicatorList":[
			"",
			""
		],
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		]
	},
	"pcfInfoSum": {
		"dnnList": [
			"",
			""
		]
	}
	}}`)
	DiscPara := &nfdiscrequest.DiscGetPara{}
	var supi []string
	supi = append(supi, "imsi-123454444")
	value := make(map[string][]string)

	value[constvalue.SearchDataSupi] = supi
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataSupi, true)

	var searchTargetNfType []string
	searchTargetNfType = append(searchTargetNfType, "AMF")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)

	if isMatchedSupiForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedSupiForNRFProfile() should return false, but return true")
	}

	var searchTargetNfType2 []string
	searchTargetNfType2 = append(searchTargetNfType2, "UDM")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType2)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)
	if !isMatchedSupiForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedSupiForNRFProfile() should return true, but return false")
	}

	nfdiscutil.PreComplieRegexp()
	var supi2 []string
	supi2 = append(supi2, "imsi-123454444")
	DiscPara.SetValue(constvalue.SearchDataSupi, supi2)
	if !isMatchedSupiForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedSupiForNRFProfile() should return true, but return false")
	}
	var supi3 []string
	supi3 = append(supi3, "imsi-12345444434")
	DiscPara.SetValue(constvalue.SearchDataSupi, supi3)
	if isMatchedSupiForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedSupiForNRFProfile() should return false, but return true")
	}
}

func TestIsMatchedGpsiForNRFPRofile(t *testing.T) {
	nrfInfo := []byte(`{ "nrfInfo" : {
	"amfInfoSum": {
		"guamiList":[
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"amfId": ""
			}
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		],
		"amfRegionIdList": [
			"",
			""
		],
		"amfSetIdList": [
			"",
			""
		],
	},
	"smfInfoSum": {
		"dnnList": [
			"",
			""
		],
		"pgwFqdnList": [
			"",
			""
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		]
	},
	"udmInfoSum": {
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"externalGroupIdentifiersRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"gpsiRanges": [
			{
				"start": "123000000",
				"end": "123999999",
				"pattern": "^imsi-12345\\d{4}$"
			}
		],
		"groupIdList": [
			"",
			""
		],
		"routingIndicatorList": [
			"",
			""
		]
	},
	"ausfInfoSum": {
		"groupIdList":[
			"",
			""
		],
		"routingIndicatorList":[
			"",
			""
		],
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		]
	},
	"pcfInfoSum": {
		"dnnList": [
			"",
			""
		]
	}
	}}`)
	DiscPara := &nfdiscrequest.DiscGetPara{}
	var supi []string
	supi = append(supi, "imsi-123454444")
	value := make(map[string][]string)
	value[constvalue.SearchDataGpsi] = supi
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataGpsi, true)

	var searchTargetNfType []string
	searchTargetNfType = append(searchTargetNfType, "AMF")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)

	if isMatchedGpsiForNRFPRofile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return false, but return true")
	}

	var searchTargetNfType2 []string
	searchTargetNfType2 = append(searchTargetNfType2, "UDM")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType2)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)
	if !isMatchedGpsiForNRFPRofile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return true, but return false")
	}

	nfdiscutil.PreComplieRegexp()
	var supi2 []string
	supi2 = append(supi2, "imsi-123454444")
	DiscPara.SetValue(constvalue.SearchDataGpsi, supi2)
	if !isMatchedGpsiForNRFPRofile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return true, but return false")
	}
	var supi3 []string
	supi3 = append(supi3, "imsi-12345444434")
	DiscPara.SetValue(constvalue.SearchDataGpsi, supi3)
	if isMatchedGpsiForNRFPRofile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return false, but return true")
	}
}

func TestIsMatchedExternalGroupIDForNRFProfile(t *testing.T) {
	nrfInfo := []byte(`{ "nrfInfo" : {
	"amfInfoSum": {
		"guamiList":[
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"amfId": ""
			}
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		],
		"amfRegionIdList": [
			"",
			""
		],
		"amfSetIdList": [
			"",
			""
		],
	},
	"smfInfoSum": {
		"dnnList": [
			"",
			""
		],
		"pgwFqdnList": [
			"",
			""
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		]
	},
	"udmInfoSum": {
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"externalGroupIdentityfiersRanges": [
			{
				"start": "123000000",
				"end": "123999999",
				"pattern": "^imsi-12345\\d{4}$"
			}
		],
		"gpsiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"groupIdList": [
			"",
			""
		],
		"routingIndicatorList": [
			"",
			""
		]
	},
	"ausfInfoSum": {
		"groupIdList":[
			"",
			""
		],
		"routingIndicatorList":[
			"",
			""
		],
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		]
	},
	"pcfInfoSum": {
		"dnnList": [
			"",
			""
		]
	}
	}}`)
	DiscPara := &nfdiscrequest.DiscGetPara{}
	var supi []string
	supi = append(supi, "imsi-123454444")
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	DiscPara.SetFlag(constvalue.SearchDataExterGroupID, true)
	DiscPara.SetValue(constvalue.SearchDataExterGroupID, supi)

	var searchTargetNfType []string
	searchTargetNfType = append(searchTargetNfType, "AMF")
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	if isMatchedExternalGroupIDForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedExternalGroupIDForNRFProfile() should return false, but return true")
	}

	var searchTargetNfType2 []string
	searchTargetNfType2 = append(searchTargetNfType2, "UDM")
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType2)
	if !isMatchedExternalGroupIDForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedExternalGroupIDForNRFProfile() should return true, but return false")
	}

	nfdiscutil.PreComplieRegexp()
	var supi2 []string
	supi2 = append(supi2, "imsi-123454444")
	DiscPara.SetValue(constvalue.SearchDataExterGroupID, supi2)
	if !isMatchedExternalGroupIDForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedExternalGroupIDForNRFProfile() should return true, but return false")
	}
	var supi3 []string
	supi3 = append(supi3, "imsi-12345444434")
	DiscPara.SetValue(constvalue.SearchDataExterGroupID, supi3)
	if isMatchedExternalGroupIDForNRFProfile(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedExternalGroupIDForNRFProfile() should return false, but return true ")
	}
}

func TestGetRegionNRFAddrFromProfile(t *testing.T) {
	nrfprofile := []byte(`{}`)
	if len(getRegionNRFAddrFromProfile(nrfprofile)) > 0 {
		t.Fatal("func getRegionNRFAddrFromProfile() region nrf addr should be null, but not")
	}
	nrfprofile2 := []byte(`{
	    "nfServices": [
		{
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"scheme": "https",
		"fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
		"serviceInstanceId": "nnrf-nfm-01",
		"supportedFeatures": "1F",
		"ipEndPoints": [
		    	{
			    "transport": "TCP",
			    "ipv4Address": "127.0.0.1",
			    "ipv6Address": "127.0.0.2",
			    "port": 443
				}
			    ],
		"apiPrefix": "",
		"priority": 5,
		"versions": [
				{
				    "apiVersionInUri": "v1",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2020-07-06T02:54:32Z"
				},
				{
				    "apiVersionInUri": "v2",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2022-07-06T02:54:32Z"
				}
			    ],
		"capacity": 100,
		"serviceName": "nnrf-disc",
		"allowedNfDomains": [],
		"allowedPlmns": [
				{
				    "mcc": "460",
				    "mnc": "00"
				}
			    ],
		"allowedNfTypes": [ ],
		"allowedNssais": [
				{
				    "sst": 1,
				    "sd": "0"
				}
			    ]
			},
		{
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"scheme": "https",
		"fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
		"serviceInstanceId": "nnrf-nfm-01",
		"supportedFeatures": "1F",
		"ipEndPoints": [
		    	{
			    "transport": "TCP",
			    "ipv4Address": "127.0.0.1",
			    "ipv6Address": "127.0.0.2",
			    "port": 443
				}
			    ],
		"apiPrefix": "",
		"priority": 5,
		"versions": [
				{
				    "apiVersionInUri": "v1",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2020-07-06T02:54:32Z"
				}
			    ],
		"capacity": 100,
		"serviceName": "nnrf-nfm",
		"allowedNfDomains": [],
		"allowedPlmns": [
				{
				    "mcc": "460",
				    "mnc": "00"
				}
			    ],
		"allowedNfTypes": [ ],
		"allowedNssais": [
				{
				    "sst": 1,
				    "sd": "0"
				}
			    ]
			}
	    ],
	    "nrfInfo":{}
	}`)
	if len(getRegionNRFAddrFromProfile(nrfprofile2)) < 0 {
		t.Fatal("func getRegionNRFAddrFromProfile() region nrf addr should be have value, but not")
	}
	addrPriority := getRegionNRFAddrFromProfile(nrfprofile2)
	value, ok := addrPriority["https://127.0.0.1:443/nnrf-disc/v1"]
	if !ok || value != "5,1" {
		t.Fatal("func getRegionNRFAddrFromProfile() ip4addr should be matched, but fail")
	}
	value2, ok2 := addrPriority["https://[127.0.0.2]:443/nnrf-disc/v1"]
	if !ok2 || value2 != "5,2" {
		t.Fatal("func getRegionNRFAddrFromProfile() ip6addr should be matched, but fail")
	}
	value3, ok3 := addrPriority["https://nrf.5gc.mnc000.mcc460.3gppnetwork.org:443/nnrf-disc/v1"]
	if !ok3 || value3 != "5,3" {
		t.Fatal("func getRegionNRFAddrFromProfile() fqdn should be matched, but fail")
	}
	value4, ok4 := addrPriority["https://127.0.0.1:443/nnrf-disc/v2"]
	if !ok4 || value4 != "5,1" {
		t.Fatal("func getRegionNRFAddrFromProfile() ip4addr should be matched, but fail")
	}
	value5, ok5 := addrPriority["https://[127.0.0.2]:443/nnrf-disc/v2"]
	if !ok5 || value5 != "5,2" {
		t.Fatal("func getRegionNRFAddrFromProfile() ip6addr should be matched, but fail")
	}
	value6, ok6 := addrPriority["https://nrf.5gc.mnc000.mcc460.3gppnetwork.org:443/nnrf-disc/v2"]
	if !ok6 || value6 != "5,3" {
		t.Fatal("func getRegionNRFAddrFromProfile() fqdn should be matched, but fail")
	}

	nrfprofile3 := []byte(`{
	    "nfServices": [
		{
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"scheme": "https",
		"fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
		"serviceInstanceId": "nnrf-nfm-01",
		"supportedFeatures": "1F",
		"apiPrefix": "",
		"priority": 5,
		"versions": [
				{
				    "apiVersionInUri": "v1",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2020-07-06T02:54:32Z"
				}
			    ],
		"capacity": 100,
		"serviceName": "nnrf-disc",
		"allowedNfDomains": [],
		"allowedPlmns": [
				{
				    "mcc": "460",
				    "mnc": "00"
				}
			    ],
		"allowedNfTypes": [ ],
		"allowedNssais": [
				{
				    "sst": 1,
				    "sd": "0"
				}
			    ]
			}
	    ],
	    "nrfInfo":{}
	}`)
	if len(getRegionNRFAddrFromProfile(nrfprofile3)) < 0 {
		t.Fatal("func getRegionNRFAddrFromProfile() region nrf addr should be have value, but not")
	}
	addrPriority2 := getRegionNRFAddrFromProfile(nrfprofile3)
	value7, ok7 := addrPriority2["https://nrf.5gc.mnc000.mcc460.3gppnetwork.org:443/nnrf-disc/v1"]
	if !ok7 || value7 != "5,4" {
		t.Fatal("func getRegionNRFAddrFromProfile() fqdn should be matched, but fail")
	}
}

func TestPlmnDiscNRFProfileFilter(t *testing.T) {
	rawNrfProfile := []byte(`{
	    "nsiList": [
		"12345"
		     ],
	    "nfSetId": "5112345",
	    "nfStatus": "REGISTERED",
	    "interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
	    "nfType": "NRF",
	    "fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
	    "locality": "",
	    "plmn": {
		"mcc": "460",
		"mnc": "000"
		},
	    "priority": 10,
	    "nfInstanceId": "0c765084-9cc5-49c6-9876-ae2f5fa2a63f",
	    "ipv6Addresses": [],
	    "capacity": 100,
	    "ipv4Addresses": [ "127.0.0.1"],
	    "sNssais": [
		     {
				"sst": 1,
				"sd": "0"
			    }
			],
	    "ipv6Prefixes": [],
	    "nfServices": [
		{
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"scheme": "https",
		"fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
		"serviceInstanceId": "nnrf-nfm-01",
		"supportedFeatures": "1F",
		"ipEndPoints": [
		    {
				    "transport": "TCP",
				    "ipv4Address": "127.0.0.1",
				    "ipv6Address": "192.168.0.1",
				    "port": 443
				}
			    ],
		"apiPrefix": "",
		"priority": 5,
		"versions": [
				{
				    "apiVersionInUri": "v1",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2020-07-06T02:54:32Z"
				}
			    ],
		"capacity": 100,
		"serviceName": "nnrf-disc",
		"allowedNfDomains": [],
		"allowedPlmns": [
				{
				    "mcc": "460",
				    "mnc": "00"
				}
			    ],
		"allowedNfTypes": [ ],
		"allowedNssais": [
				{
				    "sst": 1,
				    "sd": "0"
				}
			    ]
			}
	    ],
	    "nrfInfo":{
		    "udmInfoSum": {
			"supiRanges": [
				{
					"start": "",
					"end": "",
					"pattern": ""
				}
			],
			"externalGroupIdentifiersRanges": [
				{
					"start": "123000000",
					"end": "123999999",
					"pattern": "^imsi-12345\\d{4}$"
				}
			],
			"gpsiRanges": [
				{
					"start": "",
					"end": "",
					"pattern": ""
				}
			],
			"groupIdList": [
				"",
				""
			],
			"routingIndicatorList": [
				"",
				""
			]
		},
	    }
	}`)
	var DiscPara nfdiscrequest.DiscGetPara
	var supi []string
	supi = append(supi, "imsi-123454444")
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	DiscPara.SetFlag(constvalue.SearchDataExterGroupID, true)
	DiscPara.SetValue(constvalue.SearchDataExterGroupID, supi)

	var searchTargetNfType []string
	searchTargetNfType = append(searchTargetNfType, "UDM")

	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)

	nrfProfileGetResponse := &nrfprofile.NRFProfileGetResponse{}
	nrfProfileInfo := &nrfprofile.NRFProfileInfo{RawNrfProfile: rawNrfProfile}
	nrfProfileGetResponse.NrfProfile = append(nrfProfileGetResponse.NrfProfile, nrfProfileInfo)
	nrfAddrList,_ := plmnDiscNRFProfileFilter(nrfProfileGetResponse, DiscPara)
	if nrfAddrList[0] != "https://127.0.0.1:443/nnrf-disc/v1" || nrfAddrList[1] != "https://[192.168.0.1]:443/nnrf-disc/v1" || nrfAddrList[2] != "https://nrf.5gc.mnc000.mcc460.3gppnetwork.org:443/nnrf-disc/v1" {
		t.Fatal("func plmnDiscNRFProfileFilter() nrfAddrList should be match, but not")
	}

	var DiscPara2 nfdiscrequest.DiscGetPara
	var supi2 []string
	supi2 = append(supi, "imsi-123454")
	value2 := make(map[string][]string)
	DiscPara2.InitMember(value2)

	DiscPara2.SetFlag(constvalue.SearchDataExterGroupID, true)
	DiscPara2.SetValue(constvalue.SearchDataExterGroupID, supi2)

	DiscPara2.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	DiscPara2.SetFlag(constvalue.SearchDataTargetNfType, true)

	nrfProfileGetResponse2 := &nrfprofile.NRFProfileGetResponse{}
	nrfProfileInfo2 := &nrfprofile.NRFProfileInfo{RawNrfProfile: rawNrfProfile}
	nrfProfileGetResponse.NrfProfile = append(nrfProfileGetResponse.NrfProfile, nrfProfileInfo2)
	nrfAddrList2,_ := plmnDiscNRFProfileFilter(nrfProfileGetResponse2, DiscPara2)
	if len(nrfAddrList2) > 0 {
		t.Fatal("func plmnDiscNRFProfileFilter() nrfAddrList should be null, but has value")
	}
}

func TestIsMatchedTaiList(t *testing.T) {
	nrfInfo := []byte(`{ "nrfInfo" : {
	"amfInfoSum": {
		"guamiList":[
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"amfId": ""
			}
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "00",
					"mcc": "460"
				},
				"tac": "111111"
			}
		],
		"amfRegionIdList": [
			"",
			""
		],
		"amfSetIdList": [
			"",
			""
		],
	},
	"smfInfoSum": {
		"dnnList": [
			"",
			""
		],
		"pgwFqdnList": [
			"",
			""
		],
		"taiList": [
			{
				"plmnId": {
					"mnc": "",
					"mcc": ""
				},
				"tac": ""
			}
		]
	},
	"udmInfoSum": {
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"externalGroupIdentifiersRanges": [
			{
				"start": "123000000",
				"end": "123999999",
				"pattern": "^imsi-12345\\d{4}$"
			}
		],
		"gpsiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		],
		"groupIdList": [
			"",
			""
		],
		"routingIndicatorList": [
			"",
			""
		]
	},
	"ausfInfoSum": {
		"groupIdList":[
			"",
			""
		],
		"routingIndicatorList":[
			"",
			""
		],
		"supiRanges": [
			{
				"start": "",
				"end": "",
				"pattern": ""
			}
		]
	},
	"pcfInfoSum": {
		"dnnList": [
			"",
			""
		]
	}
	}}`)
	if !isMatchedTaiList(nrfInfo, "46000", "111111", "AMF") {
		t.Fatalf("tailist should matched, but fail")
	}
	if isMatchedTaiList(nrfInfo, "460000", "111111", "AMF") {
		t.Fatalf("tailist should not matched, but fail")
	}
	if isMatchedTaiList(nrfInfo, "46000", "111111", "UDM") {
		t.Fatalf("tailist should not matched, but fail")
	}
}
