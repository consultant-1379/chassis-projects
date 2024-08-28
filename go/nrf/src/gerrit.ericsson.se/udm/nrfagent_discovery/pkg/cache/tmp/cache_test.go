package cache

import (
	"encoding/json"
	"testing"
)

var contentUdmReg = []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
        "plmnList": [
		   {
            "mcc": "460",
            "mnc": "000"
           },
		   {
            "mcc": "560",
            "mnc": "001"
           }
		],
        "sNssais": [{
                "sst": 2,
                "sd": "2"
            },
            {
                "sst": 4,
                "sd": "4"
            }
        ],
        "nsiList": ["100","101","102"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
        "ipv4Addresses": ["172.16.208.1"],
        "ipv6Addresses": ["FF01::1101"],
        "ipv6Prefixes": ["FF01"],
        "capacity": 100,
        "load": 50,
        "locality": "Shanghai",
        "priority": 1,
        "udrInfo": {
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }]
        },
        "udmInfo": {
            "groupId": "udmtest",
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }],
            "routingIndicators": ["1111"]
        },
        "ausfInfo": {
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }],
            "routingIndicators": ["2222"]
        },
        "amfInfo": {
            "amfSetId": "amfSet-01",
            "amfRegionId": "amfRegion-01",
            "guamiList": [{
                "plmn": {
                    "mcc": "460",
                    "mnc": "000"
                },
                "amfId": "amf-01"
            }]
        },
        "smfInfo": {
            "dnnList": [
                "udm-dnn-011",
                "udm-dnn-012"
            ]
        },
        "upfInfo": {
            "sNssaiUpfInfoList": [
                {
                    "sNssai": {
                        "sst": 3,
                        "sd": "sd3"
                    },
                    "dnnUpfInfoList": [
                     {
                        "dnn": "upf-dnn1"
                     },
                     {
                        "dnn": "upf-dnn2"
                     }
                    ]
                }
            ]
        },
        "pcfInfo": {
            "dnnlist":  ["pcf-dnn1","pcf-dnn2"]
        },
        "bsfInfo": {
            "ipv4AdddressRanges": [{
                "start": "172.16.208.0",
                "end": "172.16.208.255"
            }],
            "ipv6PrefixRanges": [{
                "start": "FF01",
                "end": "FF0F"
            }]
        },
        "nfServices": [{
            "serviceInstanceId": "nudm-auth-01",
            "serviceName": "nudm-auth-01",
            "version": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "schema": "https://",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "FF01::1101",
                "ipv6Prefix": "FF01",
                "transport": "TCP",
                "port": 30080
            }],
            "apiPrefix": "nudm-uecm",
            "defaultNotificationSubscriptions": [{
                "notificationType": "N1_MESSAGES",
                "callbackUri": "https://seliius03695.seli.gic.ericsson.se",
                "n1MessageClass": "5GMM",
                "n2InformationClass": "SM"
            }],
            "capacity": 0,
            "load": 50,
            "priority": 0,
            "supportedFeatures": "A1"
        },
        {
            "serviceInstanceId": "nudm-ausf-01",
            "serviceName": "nudm-ausf-01",
            "version": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "schema": "https://",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "FF01::1101",
                "ipv6Prefix": "FF01",
                "transport": "TCP",
                "port": 30080
            }],
            "apiPrefix": "nudm-uecm",
            "defaultNotificationSubscriptions": [{
                "notificationType": "N1_MESSAGES",
                "callbackUri": "https://seliius03695.seli.gic.ericsson.se",
                "n1MessageClass": "5GMM",
                "n2InformationClass": "SM"
            }],
            "capacity": 0,
            "load": 50,
            "priority": 0,
            "supportedFeatures": "A1"
        }]
    }]
}
`)

var content = []byte(`{
	"nfInstanceId": "5g-ausf-01",
	"nfType": "AUSF",
	"nfStatus": "REGISTERED",
	"heartBeatTimer": 120,
	"plmnList": [{
	  "mcc": "460",
	  "mnc": "00"
	}],
	"sNssais": [
	  {
		"sst": 0,
		"sd": "abAB01"
	  },
	  {
		"sst": 1,
		"sd": "abAB01"
	  }
	],
	"upfInfo": {
        "smfServingArea": [
            "00"
        ],
        "sNssaiUpfInfoList": [
            {
                "sNssai": {
                    "sst": 1,
                    "sd": "AaBb01"
                },
                "dnnUpfInfoList": [
                    {
                        "dnn": "01",
                        "dnaiList": [
                            "111",
                            "222"
                        ]
                    },
                    {
                        "dnn": "02"
                    }
                ]
            },
            {
                "sNssai": {
                    "sst": 2,
                    "sd": "AaBb01"
                },
                "dnnUpfInfoList": [
                    {
                        "dnn": "11",
                        "dnaiList": [
                            "111",
                            "333"
                        ]
                    },
                    {
                        "dnn": "22"
                    }
                ]
            }
        ],
        "iwkEpsInd": false
    }, 
	"bsfInfo": {
		"dnnList": ["dnn1"],
		"ipDomainList": ["ericsson.se","ericsson.com"]
	}
  }`)

func TestSpliteSeachResult(t *testing.T) {
	nfinstanceByte, validityPeriod, _ := SpliteSeachResult(contentUdmReg)
	t.Logf("TestSpliteSeachResult validityPeriod value:%d", validityPeriod)
	if nfinstanceByte == nil {
		t.Error("SpliteSeachResult failed and nfProfiles is nil")
	}
	if validityPeriod != 86400 {
		t.Errorf("SpliteSeachResult, expect validityPeriod is 86400, but %d", validityPeriod)
	}
}

func TestIndex(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfInstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfInstanceBytes == nil {
		t.Errorf("TestIndex: SpliteSeachResult fail")
	}

	for _, instance := range nfInstanceBytes {
		instanceId, ok := mcache.indexed(instance, cacheManager.indexGroup)
		if instanceId != "udm-5g-01" && !ok {
			mcache.deIndex("udm-5g-01")
			t.Errorf("TestIndex: Index fail")
		}
	}

	mcache.flush()
}

func TestIndexNew(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfInstanceID, ok := mcache.indexed(content, cacheManager.indexGroup)
	if !ok {
		t.Fatalf("Create index for nfInstanceID[%s] profile failure\n", nfInstanceID)
	}

	mcache.showIndexContent()
	t.Logf("Create index for nfInstanceID[%s] profile success", nfInstanceID)

	ipDomainList := mcache.cacheIndex["ipDomainList"]
	dnaiList := mcache.cacheIndex["dnaiList"]
	iwkEpsInd := mcache.cacheIndex["iwkEpsInd"]

	t.Logf("Index ipDomainList[%+v]", ipDomainList)
	t.Logf("Index dnaiList[%+v]", dnaiList)
	t.Logf("Index iwkEpsInd[%+v]", iwkEpsInd)

	ids := ipDomainList["ericsson.se"]
	if ids.Contains(nfInstanceID) == false {
		t.Fatalf("expect ipDomainList[ericsson.se] index contains %s, but not", nfInstanceID)

	}

	ids = ipDomainList["ericsson.com"]
	if ids.Contains(nfInstanceID) == false {
		t.Fatalf("expect ipDomainList[ericsson.com] index contains %s, but not", nfInstanceID)
	}

	ids = iwkEpsInd["false"]
	if ids.Contains(nfInstanceID) == false {
		t.Fatalf("expect iwkEpsInd[false] index contains %s, but not", nfInstanceID)
	}

	mcache.flush()
}

func TestDeIndex(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfInstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfInstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	instanceId := "udm-5g-01"
	for _, instance := range nfInstanceBytes {
		instanceId, ok := mcache.indexed(instance, cacheManager.indexGroup)
		if instanceId != "udm-5g-01" && !ok {
			t.Errorf("Create index for nfProfile[%s] fail", instanceId)
		}
	}

	t.Logf("Current cacheIndex : %+v\n", mcache.cacheIndex)

	mcache.deIndex(instanceId)

	for _, categoryValue := range mcache.cacheIndex {
		for _, ids := range categoryValue {
			if ids.Contains(instanceId) {
				t.Errorf("After DeIndex, expect the index does not contain %s, but not", instanceId)
			}
		}
	}

	mcache.flush()
}

func TestCached(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	nfInstanceID := "udm-5g-01"
	nfProfile := mcache.profiles[nfInstanceID]
	if nfProfile != nil {
		t.Error("Before cached the profile, expect fetch result is nil, but not")
	}

	for _, instance := range nfinstanceBytes {
		mcache.cached(instance)
	}

	nfProfile = mcache.profiles[nfInstanceID]
	if nfProfile == nil {
		t.Error("After cached the profile, expect fetch result is not nil, but not")
	}
	t.Logf("nfProfile : %s", string(nfProfile))

	mcache.flush()
}

func TestDeCached(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	nfInstanceID := "udm-5g-01"
	nfProfile := mcache.profiles[nfInstanceID]
	if nfProfile != nil {
		t.Error("Before cached the profile, expect fetch result is nil, but not")
	}

	for _, instance := range nfinstanceBytes {
		mcache.cached(instance)
	}

	nfProfile = mcache.profiles[nfInstanceID]
	if nfProfile == nil {
		t.Error("After cached the profile, expect fetch result is not nil, but not")
	}

	mcache.deCached(nfInstanceID)

	nfProfile = mcache.profiles[nfInstanceID]
	if nfProfile != nil {
		t.Error("After cached the profile, expect fetch result is nil, but not")
	}

	mcache.flush()
}

func TestStatusCheck(t *testing.T) {
	setupEnv()
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	status := mcache.getCacheStatus()
	if status != false {
		t.Fatalf("Before ready, expect cache status is false, but %v", status)
	}
	t.Logf("AUSF cache status : %v", status)

	mcache.setCacheStatus(true)
	t.Log("AUSF cache set false")

	status = mcache.getCacheStatus()
	if status != true {
		t.Fatalf("After ready, expect cache status is true, but %v", status)
	}
	t.Logf("AUSF cache status : %v", status)
}

func TestFetchIDs(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	for _, instance := range nfinstanceBytes {
		mcache.cached(instance)
	}

	nfInstanceID := "udm-5g-01"
	ids := mcache.fetchIDs()
	if len(ids) != 1 {
		t.Errorf("Expect cache ids is only one, but not")
	}

	if ids[0] != nfInstanceID {
		t.Errorf("Expect cache ids is %s, but not", nfInstanceID)
	}
	t.Logf("cache ids : %+v", ids)

	mcache.flush()
}

func TestFlush(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	nfInstanceID := "udm-5g-01"
	for _, instance := range nfinstanceBytes {
		id, ok := mcache.indexed(instance, cacheManager.indexGroup)
		if id != nfInstanceID {
			t.Errorf("Expect Index return nfInstanceID is %s, but is %s", nfInstanceID, id)
		}
		if !ok {
			t.Errorf("Expect Index return true, but false")
		}
		mcache.cached(instance)
	}

	mcache.flush()

	//check cash profile after flush
	ok := mcache.probe(nfInstanceID)
	if ok {
		t.Error("After Flush, expect Probe false but true")
	}

	//check index after flush
	for _, categoryValue := range mcache.cacheIndex {
		for _, ids := range categoryValue {
			if ids.Contains(nfInstanceID) {
				t.Errorf("After Flush, expect the Index is empty, but not")
			}
		}
	}
}

func TestFetchProfileByID(t *testing.T) {
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	nfInstanceID := "udm-5g-01"
	nfProfile := mcache.fetchProfileByID(nfInstanceID)
	if nfProfile != nil {
		t.Error("Before cached, expect FetchProfileByID is nil, but not")
	}

	for _, instance := range nfinstanceBytes {
		mcache.cached(instance)
	}

	nfProfile = mcache.fetchProfileByID(nfInstanceID)
	if nfProfile == nil {
		t.Error("After cached, expect FetchProfileByID is not nil, but not")
	}

	mcache.flush()
}

func TestMeetIndexSearch(t *testing.T) { //should move the other file
	requestNfType := "AUSF"
	targetNfType := "UDM"
	mcache := cacheManager.getCache(requestNfType, targetNfType)

	nfinstanceBytes, _, _ := SpliteSeachResult(contentUdmReg)
	if nfinstanceBytes == nil {
		t.Error("SpliteSeachResult fail")
	}

	for _, instance := range nfinstanceBytes {
		mcache.indexed(instance, cacheManager.indexGroup)
		mcache.cached(instance)
	}

	var SearchParameterData = `{"target-nf-type": "udm","service-names":["nudm-auth-01"],"requester-nf-type":"udm"}`
	searchParameter := SearchParameter{}
	err := json.Unmarshal([]byte(SearchParameterData), &searchParameter)
	if err != nil {
		t.Error("TestmeetIndexSearch: Unmarshal falied")
	}

	indexSearcher := indexSearcher{
		mcache,
		&searchParameter,
		cacheManager.indexMapper,
	}
	ok := indexSearcher.meetIndexSearch()
	if !ok {
		t.Errorf("Expect meetIndexSearch is false, but not")
	}

	mcache.flush()
}
