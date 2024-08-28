package nfdiscutil

import (
	"net/http"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestPreComplieRegexp(t *testing.T) {
	PreComplieRegexp()
}

func TestGetNFProfileMD5Sum(t *testing.T) {
	customNFProfile := []byte(`{
        "nfProfile": "25d7476ff58fec0e0975a5d0edb9f0ab",
        "nudm-uecm-01": "740bca9170d7756e9c9d0ff96b9041fe",
        "nudm-uecm-02": "5f97289663520ec178dbb6247fb126ce",
        "nudm-sdm-01": "c91f05388da3c67c3d1f94559be06cdf"
    }`)
	originalNFProfile := []byte(`[
            {
                "capacity": 100,
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.1",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-01",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            },
            {
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.2",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-02",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            }
        ]`)
	if "25d7476ff58fec0e0975a5d0edb9f0ab740bca9170d7756e9c9d0ff96b9041fe5f97289663520ec178dbb6247fb126ce" != GetNFProfileMD5Sum(customNFProfile, originalNFProfile) {
		t.Fatalf("func GetNFProfileMD5Sum should get the md5sum matched, but failed")
	}
}

func TestIsAllowedNfType(t *testing.T) {
	nfService := []byte(`{
						    "serviceInstanceId": "srv1",
                             "serviceName": "srv1",
                             "version": [],
                             "schema": "schema1",
							"allowedPlmns": [
                                 {
                                     "mcc": "000",
                                     "mnc": "00"
                                 },
								{
                                     "mcc": "001",
                                     "mnc": "01"
                                 }
                             ],
                             "allowedNfTypes": [
                                 "NRF",
                                 "UDM"
                             ],

					}`)
	if !IsAllowedNfType(nfService, "NRF", constvalue.NFServiceAllowedNFTypes) {
		t.Fatalf("NRF is allowed nf type, but return false !")
	}

	if !IsAllowedNfType(nfService, "UDM", constvalue.NFServiceAllowedNFTypes) {
		t.Fatalf("UDM is allowed nf type, but return false !")
	}

	if IsAllowedNfType(nfService, "AUSF", constvalue.NFServiceAllowedNFTypes) {
		t.Fatalf("AUSF is not allowed nf type, but return true !")
	}
}

func TestIsAllowedNfFQDN(t *testing.T) {
	nfService := []byte(`{
						    "serviceInstanceId": "srv1",
                             "serviceName": "srv1",
                             "version": [],
                             "schema": "schema1",
							"allowedPlmns": [
                                 {
                                     "mcc": "000",
                                     "mnc": "00"
                                 },
								{
                                     "mcc": "001",
                                     "mnc": "01"
                                 }
                             ],
                             "allowedNfTypes": [
                                 "NRF",
                                 "UDM"
                             ],
                             "allowedNfDomains":[
                                 "^seliius\\d{5}.seli.gic.ericsson.se$",
                                 "^seliius\\d{4}.seli.gic.ericsson.se$"
                             ]

					}`)
	if !IsAllowedNfFQDN(nfService, "seliius12121.seli.gic.ericsson.se", constvalue.AllowedNfDomains) {
		t.Fatalf("NF is allowed nf domains, but return false !")
	}

	if !IsAllowedNfFQDN(nfService, "seliius1211.seli.gic.ericsson.se", constvalue.AllowedNfDomains) {
		t.Fatalf("NF is allowed nf domains, but return false !")
	}

	if IsAllowedNfFQDN(nfService, "seliius121271.seli.gic.ericsson.se", constvalue.AllowedNfDomains) {
		t.Fatalf("NF is not allowed nf domains, but return true !")
	}
}

func TestIsAllowedPLMN(t *testing.T) {
	nfService := []byte(`{
						    "serviceInstanceId": "srv1",
                             "serviceName": "srv1",
                             "version": [],
                             "schema": "schema1",
							"allowedPlmns": [
                                 {
                                     "mcc": "000",
                                     "mnc": "00"
                                 },
								{
                                     "mcc": "001",
                                     "mnc": "01"
                                 }
                             ],
                             "allowedNfTypes": [
                                 "NRF",
                                 "UDM"
                             ],

					}`)
	var plmnList []string
	plmnList = append(plmnList, "00000")
	plmnList = append(plmnList, "00020")
	if !IsAllowedPLMN(nfService, plmnList, constvalue.AllowedPlmns) {
		t.Fatalf("mcc=000, mnc=00 is allowed plmn, but return false !")
	}
	var plmnList2 []string
	plmnList2 = append(plmnList2, "00101")
	plmnList2 = append(plmnList2, "00201")
	if !IsAllowedPLMN(nfService, plmnList2, constvalue.AllowedPlmns) {
		t.Fatalf("mcc=001, mnc=01 is allowed plmn, but return false !")
	}
	var plmnList3 []string
	plmnList3 = append(plmnList3, "00001")
	plmnList3 = append(plmnList3, "00002")
	if IsAllowedPLMN(nfService, plmnList3, constvalue.AllowedPlmns) {
		t.Fatalf("mcc=000, mnc=01 is not allowed plmn, but return true !")
	}
}
func TestIsPlmnMatchHomePlmn(t *testing.T) {
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "000", Mcc: "460"})
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "111", Mcc: "460"})

	var plmnList []string
	plmnList = append(plmnList, "46000")
	plmnList = append(plmnList, "46011")
	if !IsPlmnMatchHomePlmn(plmnList) {
		t.Fatal("func IsPlmnMatchHomePlmn() should return true, but not")
	}

	var plmnList2 []string
	plmnList2 = append(plmnList2, "46011")
	plmnList2 = append(plmnList2, "46022")
	if IsPlmnMatchHomePlmn(plmnList2) {
		t.Fatal("func IsPlmnMatchHomePlmn() should return false, but return true")
	}
}

func TestGetRequestParam(t *testing.T) {
	request := "/nnrf-disc/v1/nf-instances?service-names=namf-auth&target-nf-type=AMF&requester-nf-type=UDM&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&tai={\"plmnId\":{\"mcc\":\"310\",\"mnc\":\"0109\"},\"tac\":\"Bc11\"}"
	if GetRequestParam(request) != "/nf-instances?service-names=namf-auth&target-nf-type=AMF&requester-nf-type=UDM&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&tai={\"plmnId\":{\"mcc\":\"310\",\"mnc\":\"0109\"},\"tac\":\"Bc11\"}" {
		t.Fatal("request should substring success, but fail")
	}
}

func TestGetRequestUriRoot(t *testing.T) {
	request := "http://10.111.137.76:3000/nnrf-disc/v1/nf-instances?service-names=namf-auth"
	if GetRequestURIRoot(request) != "http://10.111.137.76:3000/nnrf-disc/v1/" {
		t.Fatal("request should substring success, but fail")
	}
}

func TestGetRequestUriVersion(t *testing.T) {
	request := "http://10.111.137.76:3000/nnrf-disc/v1/nf-instances?service-names=namf-auth"
	if GetRequestURIVersion(request) != "v1" {
		t.Fatal("request version should be v1, but not")
	}
}

func TestFilterAddrWithVersion(t *testing.T) {
	var addrs []string
	addrs = append(addrs, "http://10.111.137.76:3000/nnrf-disc/v1")
	addrs = append(addrs, "http://10.111.137.76:3000/nnrf-disc/v3")
	addrs = append(addrs, "http://10.111.137.76:3000/nnrf-disc/v2")
	addrs = append(addrs, "http://190.168.0.1:3000/nnrf-disc/v1")
	addrs = FilterAddrWithVersion(addrs, "/nnrf-disc/v1/nf-instances?service-names=namf-auth")
	if len(addrs) != 2 {
		t.Fatal("filter addrs length should be 2, but not")
	}
	if addrs[0] != "http://10.111.137.76:3000/nnrf-disc/v1" {
		t.Fatal("addrs[0] should match, but not")
	}
	if addrs[1] != "http://190.168.0.1:3000/nnrf-disc/v1" {
		t.Fatal("addrs[1] should be match, but not")
	}
}

func TestStatusCodeDirectReturn(t *testing.T) {
	statusCode := http.StatusOK
	if !StatusCodeDirectReturn(statusCode) {
		t.Fatal("func statusCodeDirectReturn() should return true, but return false")
	}
	statusCode2 := http.StatusNotFound
	if !StatusCodeDirectReturn(statusCode2) {
		t.Fatal("func statusCodeDirectReturn() should return true, but return false")
	}
	statusCode3 := http.StatusBadRequest
	if !StatusCodeDirectReturn(statusCode3) {
		t.Fatal("func statusCodeDirectReturn() should return true, but return false")
	}
	statusCode4 := http.StatusForbidden
	if !StatusCodeDirectReturn(statusCode4) {
		t.Fatal("func statusCodeDirectReturn() should return true, but return false")
	}
	statusCode5 := http.StatusLengthRequired
	if !StatusCodeDirectReturn(statusCode5) {
		t.Fatal("func statusCodeDirectReturn() should return true, but return false")
	}
	statusCode6 := http.StatusInternalServerError
	if StatusCodeDirectReturn(statusCode6) {
		t.Fatal("func statusCodeDirectReturn() should return false, but return true")
	}
}
