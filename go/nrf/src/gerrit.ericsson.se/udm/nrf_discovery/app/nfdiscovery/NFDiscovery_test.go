package nfdiscovery

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

func TestIsNeedRoaming(t *testing.T) {
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "000", Mcc: "460"})
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "111", Mcc: "460"})

	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"111\"}")
	DiscPara.SetFlag(constvalue.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataTargetPlmnList, plmnList)
	roaming, _ := isNeedRoaming(DiscPara)
	if roaming {
		t.Fatal("func isNeedRoaming() should return false, but not")
	}
	var DiscPara2 nfdiscrequest.DiscGetPara
	value2 := make(map[string][]string)
	DiscPara2.InitMember(value2)

	var plmnList2 []string
	plmnList2 = append(plmnList2, "{\"mcc\":\"460\", \"mnc\":\"11\"}")
	plmnList2 = append(plmnList2, "{\"mcc\":\"460\", \"mnc\":\"22\"}")
	DiscPara2.SetFlag(constvalue.SearchDataTargetPlmnList, true)
	DiscPara2.SetValue(constvalue.SearchDataTargetPlmnList, plmnList2)
	roaming, _ = isNeedRoaming(DiscPara2)
	if !roaming {
		t.Fatal("func isNeedRoaming() should return true, but not")
	}
}

func TestPreComplieRegexp(t *testing.T) {
	nfdiscutil.PreComplieRegexp()
	matched := nfdiscutil.Compile[constvalue.SearchDataSupportedFeatures].MatchString("AAAAA")
	if !matched {
		t.Fatalf("Regexp not matched")
	}

	matched = nfdiscutil.Compile[constvalue.SearchDataHnrfURI].MatchString("htttp")
	if matched {
		t.Fatalf("Regexp should not match, but match")
	}

	matched = nfdiscutil.Compile[constvalue.SearchDataExterGroupID].MatchString("aaaaaa")
	if matched {
		t.Fatalf("Regexp not matched")
	}

	matched = nfdiscutil.Compile[constvalue.SearchDataUEIPv4Addr].MatchString("1111.1.1.1")
	if matched {
		t.Fatalf("Should not match, but match")
	}

}

func TestDiscoveryResponseHander(t *testing.T) {
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError
	url := `/nnrf-disc/v1//nf-instances?&service-names=namf-comm&target-nf-type=AMF&requester-nf-type=UDM`
	req := httptest.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	DiscoveryResponseHander(resp, req, "", http.StatusOK, "success")
	if resp.Code != http.StatusOK {
		t.Fatal("response statuscode should be 200, but not")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && string(body) != "success" {
		t.Fatal("response body should be success, but not")
	}
	if resp.Header().Get("Content-Type") != "application/json" {
		t.Fatal("response header should be matched, but not")
	}

	internalconf.HTTPWithXVersion = true
	cm.ServiceVersion = "x11"
	resp2 := httptest.NewRecorder()
	resp2.Code = http.StatusInternalServerError
	req2 := httptest.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req2.Header.Set("Content-Type", "application/json")
	DiscoveryResponseHander(resp2, req2, "", http.StatusBadGateway, "fail")
	if resp2.Code != http.StatusBadGateway {
		t.Fatal("response statuscode should be 502, but not")
	}
	body2, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil && string(body2) != "fail" {
		t.Fatal("response body should be fail, but not")
	}
	if resp2.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatal("response header should be matched, but not")
	}
	if resp2.Header().Get("X-Version") != "x11" {
		t.Fatal("X-Version should be matched, but not")
	}
}
