package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"net/url"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
)
func TestIsMatchedSupi(t *testing.T) {
       nfdiscutil.PreComplieRegexp()
	udrInfo := []byte(`{
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
      }	`)
	var DiscPara1 nfdiscrequest.DiscGetPara
	var supiArray1 []string
	supiArray1 = append(supiArray1, "imsi-123456789041234")
	value := make(map[string][]string)
	DiscPara1.InitMember(value)
	DiscPara1.SetValue("supi", supiArray1)
	DiscPara1.SetFlag("supi", true)
	var nfTypeArray1 []string
	nfTypeArray1 = append(nfTypeArray1, "UDR")
	DiscPara1.SetFlag("target-nf-type", true)
	DiscPara1.SetValue("target-nf-type", nfTypeArray1)
	if !isMatchedSupi(&DiscPara1, []byte(udrInfo)) {
		t.Fatalf("should be matched , but failed")
	}

	var DiscPara2 nfdiscrequest.DiscGetPara
	var supiArray2 []string
	supiArray2 = append(supiArray2, "imsi-123456789041234")
	value2 := make(map[string][]string)
	DiscPara2.InitMember(value2)
	DiscPara2.SetValue("supi", supiArray2)
	DiscPara2.SetFlag("supi", true)
	var nfTypeArray2 []string
	nfTypeArray2 = append(nfTypeArray2, "UDR")
	DiscPara2.SetFlag("target-nf-type", true)
	DiscPara2.SetValue("target-nf-type", nfTypeArray2)

	if !isMatchedSupi(&DiscPara2, []byte(udrInfo)) {
		t.Fatalf("should be matched , but failed")
	}

	var DiscPara3 nfdiscrequest.DiscGetPara
	var supiArray3 []string
	supiArray3 = append(supiArray3, "imsi-123456789051234")
	value3 := make(map[string][]string)
	DiscPara3.InitMember(value3)
	DiscPara3.SetValue("supi", supiArray3)
	DiscPara3.SetFlag("supi", true)
	var nfTypeArray3 []string
	nfTypeArray3 = append(nfTypeArray3, "UDR")
	DiscPara3.SetFlag("target-nf-type", true)
	DiscPara3.SetValue("target-nf-type", nfTypeArray3)

	if isMatchedSupi(&DiscPara3, []byte(udrInfo)) {
		t.Fatalf("should NOT be matched , but failed")
	}

	var DiscPara4 nfdiscrequest.DiscGetPara
	var supiArray4 []string
	supiArray4 = append(supiArray4, "suci-223456789041234")
	value4 := make(map[string][]string)
	DiscPara4.InitMember(value4)
	DiscPara4.SetValue("supi", supiArray4)
	DiscPara4.SetFlag("supi", true)
	var nfTypeArray4 []string
	nfTypeArray4 = append(nfTypeArray4, "UDR")
	DiscPara4.SetFlag("target-nf-type", true)
	DiscPara4.SetValue("target-nf-type", nfTypeArray4)
	if !isMatchedSupi(&DiscPara4, []byte(udrInfo)) {
		t.Fatalf("should  be matched , but failed")
	}

	var DiscPara5 nfdiscrequest.DiscGetPara
	var supiArray5 []string
	supiArray5 = append(supiArray5, "suci-323456789041234")
	value5 := make(map[string][]string)
	DiscPara5.InitMember(value5)
	DiscPara5.SetValue("supi", supiArray5)
	DiscPara5.SetFlag("supi", true)
	var nfTypeArray5 []string
	nfTypeArray5 = append(nfTypeArray5, "UDR")
	DiscPara5.SetFlag("target-nf-type", true)
	DiscPara5.SetValue("target-nf-type", nfTypeArray5)

	if isMatchedSupi(&DiscPara5, []byte(udrInfo)) {
		t.Fatalf("should  NOT be matched , but failed")
	}

	var DiscPara6 nfdiscrequest.DiscGetPara
	var supiArray6 []string
	supiArray6 = append(supiArray6, "nai-smartmeter-ericsson@company.com")
	value6 := make(map[string][]string)
	DiscPara6.InitMember(value6)
	DiscPara6.SetValue("supi", supiArray6)
	DiscPara6.SetFlag("supi", true)
	var nfTypeArray6 []string
	nfTypeArray6 = append(nfTypeArray6, "UDR")
	DiscPara6.SetFlag("target-nf-type", true)
	DiscPara6.SetValue("target-nf-type", nfTypeArray6)

	if !isMatchedSupi(&DiscPara6, []byte(udrInfo)) {
		t.Fatalf("should  be matched , but failed")
	}

	var DiscPara7 nfdiscrequest.DiscGetPara
	var supiArray7 []string
	supiArray7 = append(supiArray7, "nai-phone-ericsson@company.com")
	value7 := make(map[string][]string)
	DiscPara7.InitMember(value7)
	DiscPara7.SetValue("supi", supiArray7)
	DiscPara7.SetFlag("supi", true)
	var nfTypeArray7 []string
	nfTypeArray7 = append(nfTypeArray7, "UDR")
	DiscPara7.SetFlag("target-nf-type", true)
	DiscPara7.SetValue("target-nf-type", nfTypeArray7)
	if isMatchedSupi(&DiscPara7, []byte(udrInfo)) {
		t.Fatalf("should  NOT be matched , but failed")
	}

}


func TestIsMatchedGroupID(t *testing.T)  {
	nfInfo := []byte(`{
	"groupId": "123"
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
      }	`)
	nfInfo2 := []byte(`{
	"groupId": "123",
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
      }	`)
	var DiscPara1 nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara1.InitMember(value)
	var nfTypeArray1 []string
	nfTypeArray1 = append(nfTypeArray1, "UDR")
	DiscPara1.SetFlag("target-nf-type", true)
	DiscPara1.SetValue("target-nf-type", nfTypeArray1)
	var groupIdList []string
	groupIdList = append(groupIdList, "123", "456")
	if nfdiscutil.ResultFoundMatch != isMatchedGroupID(&DiscPara1, groupIdList, []byte(nfInfo)) {
		t.Fatalf("func isMatchedGroupID() should be matched , but failed")
	}
	var groupIdList2 []string
	groupIdList2 = append(groupIdList2, "1234", "456")
	if nfdiscutil.ResultFoundNotMatch != isMatchedGroupID(&DiscPara1, groupIdList2, []byte(nfInfo)) {
		t.Fatalf("func isMatchedGroupID() should not be matched , but matched")
	}
	if nfdiscutil.ResultFoundNotMatch != isMatchedGroupID(&DiscPara1, groupIdList2, []byte(nfInfo2)) {
		t.Fatalf("func isMatchedGroupID() should not be matched , but matched")
	}
}

func TestIsMatchedGpsi(t *testing.T) {
	nrfInfo := []byte(`{
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
	}`)
	DiscPara := &nfdiscrequest.DiscGetPara{}
	var supi []string
	supi = append(supi, "imsi-123454444")
	value := make(map[string][]string)
	value[constvalue.SearchDataGpsi] = supi
	DiscPara.InitMember(value)
	DiscPara.SetFlag(constvalue.SearchDataGpsi ,true)

	var searchTargetNfType []string
	searchTargetNfType = append(searchTargetNfType, "AMF")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)

	//if isMatchedGpsi(DiscPara, nrfInfo) {
	//	t.Fatal("func isMatchedGpsiForNRFPRofile() should return false, but return true")
	//}

	var searchTargetNfType2 []string
	searchTargetNfType2 = append(searchTargetNfType2, "UDM")
	DiscPara.SetValue(constvalue.SearchDataTargetNfType, searchTargetNfType2)
	DiscPara.SetFlag(constvalue.SearchDataTargetNfType, true)
	if !isMatchedGpsi(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return true, but return false")
	}

	nfdiscutil.PreComplieRegexp()
	var supi2 []string
	supi2 = append(supi2, "imsi-123454444")
	DiscPara.SetValue(constvalue.SearchDataGpsi, supi2)
	if !isMatchedGpsi(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return true, but return false")
	}
	var supi3 []string
	supi3 = append(supi3, "imsi-12345444434")
	DiscPara.SetValue(constvalue.SearchDataGpsi, supi3)
	if isMatchedGpsi(DiscPara, nrfInfo) {
		t.Fatal("func isMatchedGpsiForNRFPRofile() should return false, but return true")
	}
}


func TestIsSnssaisParaOnly(t *testing.T)  {
	nfdiscutil.PreComplieRegexp()
	var queryform0 nfdiscrequest.DiscGetPara
	oriqueryForm0, err0 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&snssais={"sst":2, "sd": "222222"}`)
	if err0 != nil {
		t.Fatal("url parse error")
	}
	queryform0.InitMember(oriqueryForm0)
	queryform0.ValidateNRFDiscovery()
	if !isSnssaisParaOnly(&queryform0) {
		t.Fatal("parameter just has snssais, should return true, but return false")
	}

	var queryform1 nfdiscrequest.DiscGetPara
	oriqueryForm1, err1 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&snssais={"sst":2, "sd": "222222"}&dnn=222222`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform1.InitMember(oriqueryForm1)
	queryform1.ValidateNRFDiscovery()
	if isSnssaisParaOnly(&queryform1) {
		t.Fatal("parameter have snssais and dnn, should return false, but return true")
	}

	var queryform2 nfdiscrequest.DiscGetPara
	oriqueryForm2, err1 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&snssais={"sst":2, "sd": "222222"}&dnai-list=222222`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform2.InitMember(oriqueryForm2)
	queryform2.ValidateNRFDiscovery()
	if isSnssaisParaOnly(&queryform2) {
		t.Fatal("parameter have snssais and dnai, should return false, but return true")
	}

	var queryform3 nfdiscrequest.DiscGetPara
	oriqueryForm3, err1 := url.ParseQuery(`service-names=nudr-uecm&target-nf-type=UDM&requester-nf-type=AUSF&snssais={"sst":2, "sd": "222222"}&dnn=222222&dnai-list=222222`)
	if err1 != nil {
		t.Fatal("url parse error")
	}
	queryform3.InitMember(oriqueryForm3)
	queryform3.ValidateNRFDiscovery()
	if isSnssaisParaOnly(&queryform3) {
		t.Fatal("parameter have snssais dnn and dnai, should return false, but return true")
	}
}

func TestIsAllowedRequesterPlmn(t *testing.T) {

	var requesterPlmnList []string
	requesterPlmnList = append(requesterPlmnList, "460000")
	requesterPlmnList = append(requesterPlmnList, "460003")

	profile := []byte(`{}`)
	if !isAllowedRequesterPlmn(profile, requesterPlmnList, "allowedPlmns", profile) {
		t.Fatal("allowed not exist, should match all, but not")
	}

	profile1 := []byte(`{"allowedPlmns":[{"mcc":"460","mnc":"000"}, {"mcc":"460","mnc":"001"}]}`)
	if !isAllowedRequesterPlmn(profile1, requesterPlmnList, "allowedPlmns", profile1) {
		t.Fatal("allowed should matched, but not ")
	}

	profile2 := []byte(`{"allowedPlmns":[{"mcc":"460","mnc":"001"}], "plmnList" :[{"mcc":"460","mnc":"000"}]}`)
	if !isAllowedRequesterPlmn(profile2, requesterPlmnList, "allowedPlmns", profile2) {
		t.Fatal("PlmnList should matched, but not")
	}

	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "000", Mcc: "460"})
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "111", Mcc: "460"})
	profile3 := []byte(`{"allowedPlmns":[{"mcc":"460","mnc":"001"}]}`)
	if !isAllowedRequesterPlmn(profile3, requesterPlmnList, "allowedPlmns", profile3) {
		t.Fatal("cm plmnid should matched, but not")
	}

}

