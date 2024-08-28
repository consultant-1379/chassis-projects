package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestIsMatchedNFProfileFQDN(t *testing.T) {
	filter := &NFCommonFilter{}
	nfProfile := []byte(`{
	    "fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",

	}`)
	matched := filter.isMatchedNFProfileFQDN("test", nfProfile)
	if matched {
		t.Fatal("fqdn should not be matched, but matched")
	}
	matched2 := filter.isMatchedNFProfileFQDN("nrf.5gc.mnc000.mcc460.3gppnetwork.org", nfProfile)
	if !matched2 {
		t.Fatal("fqdn should be matched, but not matched")
	}
}

func TestIsMatchedNFProfilesStatus(t *testing.T)  {
	filter := &NFCommonFilter{}
	nfProfile := []byte(`{
	    "nfStatus": "REGISTERED",

	}`)
	nfProfile2 := []byte(`{
	    "nfStatus": "SUSPENDED",

	}`)
	matched := filter.isMatchedNFProfilesStatus(nfProfile2)
	if matched {
		t.Fatal("status should not be matched, but matched")
	}
	matched2 := filter.isMatchedNFProfilesStatus(nfProfile)
	if !matched2 {
		t.Fatal("status should be matched, but not matched")
	}
}

func TestIsMatchedNsiList(t *testing.T)  {
	filter := &NFCommonFilter{}
	nfProfile := []byte(`{
	    "nsiList": [
		"12345",
		"22222"
	]

	}`)
	var nsiList []string
	nsiList = append(nsiList, "11111", "22222")
	var nsiList2 []string
	nsiList2 = append(nsiList2, "11111", "33333")
	matched := filter.isMatchedNsiList(nsiList2, nfProfile)
	if matched {
		t.Fatal("nsiList should not be matched, but matched")
	}
	matched2 := filter.isMatchedNsiList(nsiList, nfProfile)
	if !matched2 {
		t.Fatal("nsiList should be matched, but not matched")
	}
}

func TestCommonFilter(t *testing.T) {
	nfProfile := []byte(`{
	  "nfInstanceId": "12345678-9abc-def0-1000-100000000021",
	  "nfType": "UPF",
	  "nfStatus": "REGISTERED",
	  "ipv4Addresses": [
	    "172.16.208.1"
	  ],
	  "nsiList": ["111111"],
	  "allowedNfTypes": ["UDM"],
	  "allowedPlmns": [{"mcc":"460","mnc":"000"}],
	  "allowedNfDomains": ["www.test.com"],
	  "upfInfo": {
	     "sNssaiUpfInfoList": [
		{
		  "sNssai": {
		     "sst": 0,
		     "sd": "000000"
		  },
		  "dnnUpfInfoList": [
		     {
		      "dnn": "01",
		      "dnaiList": ["11"]
		     },
		     {
		      "dnn": "02"
		     }

		  ]
		},
	       {
		  "sNssai": {
		     "sst": 0,
		     "sd": "111111"
		  },
		  "dnnUpfInfoList": [
		     {
		      "dnn": "03",
		      "dnaiList": ["22"]
		     },
		     {
		      "dnn": "04"
		     }

		  ]
		}
	     ],
	     "smfServingArea": ["smf stest"],
	     "interfaceUpfInfoList": [
		{
		  "interfaceType": "N3",
		  "ipv4EndpointAddresses": ["127.0.0.1"],
		  "ipv6EndpointAddresses": ["fe80::0000"],
		  "endpointFqdn": "upf.seli.gic.ericsson.se",
		  "networkInstance": "networkInstance test"
		}
	     ],
	     "iwkEpsInd": true
	  },
	  "nfServices": [
	    {
	      "serviceInstanceId": "nupf-test-01",
	      "serviceName": "nupf-test",
	      "nfServiceStatus": "REGISTERED",
	      "versions": [{
		"apiVersionInUri":"v1",
		"apiFullVersion": "1.R15.1.1 " ,
		"expiry":"2020-07-06T02:54:32Z"}],
	      "scheme": "https",
	      "priority": 100,
	      "fqdn": "seliius03695.seli.gic.ericsson.se",
	      "allowedNfTypes": [
		"UDM"
	      ]

	    }
	  ]
	}
	`)
	commonfilter := &NFCommonFilter{}
	filterInfo := &FilterInfo{}

	req := &nfdiscrequest.DiscGetPara{}
	var nftype []string
	nftype = append(nftype, "UPF")
	value := make(map[string][]string)
	value[constvalue.SearchDataTargetNfType] = nftype
	req.InitMember(value)
	req.SetFlag(constvalue.SearchDataTargetNfType ,true)

	if !commonfilter.filter(nfProfile, req, filterInfo) {
		t.Fatal("filter nftype=UPF should return true, but return false")
	}

	var nsiList []string
	nsiList = append(nsiList, "222222")
	req.SetValue(constvalue.SearchDataNsiList, nsiList)
	req.SetFlag(constvalue.SearchDataNsiList, true)
	if commonfilter.filter(nfProfile, req, filterInfo) {
		t.Fatal("filter nsiList should return false, but return true")
	}

	var nsiList2 []string
	nsiList2 = append(nsiList2, "111111")
	req.SetValue(constvalue.SearchDataNsiList, nsiList2)
	req.SetFlag(constvalue.SearchDataNsiList, true)
	if !commonfilter.filter(nfProfile, req, filterInfo) {
		t.Fatal("filter nsiList should return true, but return false")
	}

	var requesterNfType []string
	requesterNfType = append(requesterNfType, "UDM")
	req.SetValue(constvalue.SearchDataRequesterNfType, requesterNfType)
	req.SetFlag(constvalue.SearchDataRequesterNfType, true)
	if !commonfilter.filter(nfProfile, req, filterInfo) || filterInfo.nfTypeForbiddenInProfile {
		t.Fatal("requesterNfType should be match, but fail")
	}

	filterInfo2 := &FilterInfo{}
	var requesterNfType2 []string
	requesterNfType2 = append(requesterNfType2, "PCF")
	req.SetValue(constvalue.SearchDataRequesterNfType, requesterNfType2)
	req.SetFlag(constvalue.SearchDataRequesterNfType, true)
	if commonfilter.filter(nfProfile, req, filterInfo2) || !filterInfo2.nfTypeForbiddenInProfile {
		t.Fatal("requesterNfType should not be match, but match")
	}

	req2 := &nfdiscrequest.DiscGetPara{}
	req2.InitMember(value)
	req2.SetValue(constvalue.SearchDataTargetNfType, nftype)
	req2.SetFlag(constvalue.SearchDataTargetNfType ,true)
	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\",\"mnc\":\"000\"}")
	req2.SetValue(constvalue.SearchDataRequesterPlmnList, plmnList)
	req2.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	if !commonfilter.filter(nfProfile, req2, filterInfo) || filterInfo.plmnForbiddenInProfile {
		t.Fatal("plmnList should be match, but fail")
	}

	var plmnList2 []string
	plmnList2 = append(plmnList2, "{\"mcc\":\"460\",\"mnc\":\"00\"}")
	req2.SetValue(constvalue.SearchDataRequesterPlmnList, plmnList2)
	req2.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	if commonfilter.filter(nfProfile, req2, filterInfo2) || !filterInfo2.plmnForbiddenInProfile {
		t.Fatal("plmnList should be match, but fail")
	}
}