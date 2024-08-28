package disc

import (
	"bytes"
	"encoding/json"
	"errors"

	//"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

var cmDataNrfService = []byte(`
{
       "notification-address": {
          "ipv4-address": "127.0.0.1",
          "port": 85,
          "transport": "tcp"
        },
        "nrf": {
          "mode": "active-standby",
          "profile": [
            {
              "id": "nrf-server-0",
              "ipv4-address": ["127.0.0.1"],
              "priority":0,
              "service": [
                {
                  "id": 0,
                  "name": "nnrf-nfm",
                  "scheme": "http",
                  "priority":0,
                  "version": [
                    {
                      "api-version-in-uri": "v1"
                    }
                  ],
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "port": 80,
                      "transport": "tcp"
                    }
                  ],
                  "api-prefix": "nnrf-nfm"
                },
                {
                  "id": 1,
                  "name": "nnrf-disc",
                  "scheme": "http",
                  "priority":0,
                  "version": [
                    {
                      "api-version-in-uri": "v1"
                    }
                  ],
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "port": 80,
                      "transport": "tcp"
                    }
                  ],
                  "api-prefix": "nnrf-disc"
                }
              ]
            }
          ]
        }
    }
`)

var cmDataNrfServicePort3212 = []byte(`
    {
        "notification-address": {
          "ipv4-address": "127.0.0.1",
          "ipv6-address": "",
          "port": "3212",
          "transport": "TCP"
        },
        "nrf": {
          "mode": "active-standby",
          "profile": [
            {
			 "id": "nrf-server-0",
              "ipv4-address": ["127.0.0.1"],
			 "priority": 0,
              "service": [
                {
                  "id": 0,
                  "scheme": "http",
                  "version: [
                    {
                      "api-version-in-uri": "v1"
                    }
                  ],
                  "fqdn": "",
				"priority": 0,
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "ipv6-address": "",
                      "port": "3212",
                      "transport": "TCP"
                    }
                  ],
                  "api-prefix": "nnrf-disc",
                  "name": "nnrf-disc"
                },
                {
                  "id": 1,
                  "scheme": "http",
				"priority": 0,
                  "version": [
                    {
                       "api-version-in-uri": "v1"
                    }
                  ],
                  "fqdn": "",
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "ipv6-address": "",
                      "port": "3212",
                      "transport": "TCP"
                    }
                  ],
                 "api-prefix": "nnrf-nfm",
                  "name": "nnrf-nfm"
                }
              ]
            }
          ]
        }
    }
`)

var cmTargetNfProfile = []byte(`{
      "targetNfProfiles": [
            {
		  "requesterNfType": "AUSF",
		  "requesterNfFqdn": "",
          "targetNfType": "UDM",
          "targetServiceNames": [
                    "nudm-auth-01"
                  ],
		  "notifCondition" : {
				"monitoredAttributes": [
					"/capacity",
					"/priority"
				]
		  }
                },
				{
		  "requesterNfType": "AUSF",
		  "requesterNfFqdn": "",
          "targetNfType": "UDR",
          "targetServiceNames": [
                    "nudr-01"
                  ],
		 "notifCondition" : {
			"unmonitoredAttributes": [
					"/load",
					"/priority"
				]
		  }		 
                }
      ]
    }`)
var cmTargetNfProfileWrong = []byte(`{
      "targetNfProfilesWrong": [
            {
          "targetNfType": "UDM",
          "serviceNames": [
                    "nudm-auth-01"
                  ]
                }
      ]
    }`)

var searchResultUDM2 = []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-02",
        "nfType": "UDM",
		"nfStatus": "REGISTERED",
        "plmnList": [{
            "mcc": "466",
            "mnc": "001"
          },
		 {
		    "mcc": "460",
            "mnc": "001"
		 }
		],
        "sNssais": [{
                "sst": 2,
                "sd": "A00000"
            },
            {
                "sst": 4,
                "sd": "A00000"
            }
        ],
        "nsiList": ["069","002","102"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
        "ipv4Addresses": ["172.16.208.1"],
        "ipv6Addresses": ["1001:da8::36"],
        "capacity": 100,
        "load": 50,
        "locality": "Shanghai",
        "priority": 1,
        "udrInfo": {
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }],
			  "groupId": "udr-02"
        },
		"amfInfo": {
            "amfSetId": "amfSet-01",
            "amfRegionId": "amfRegion-01",
            "guamiList": [{
                "plmnId": {
                    "mcc": "460",
                    "mnc": "000"
                },
               "amfId": "a00001"
              }
			]
        },
		"upfInfo": {
            "sNssaiUpfInfoList": [
                {
                    "sNssai": {
                        "sst": 3,
                        "sd": "A00001"
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
        "udmInfo": {
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }],
            "groupId": "udm-02"
        },
        "ausfInfo": {
            "supiRanges": [{
                "start": "000001",
                "end": "000010"
            }],
            "groupId": "ausf-01"
        },
        "smfInfo": {
			"sNssaiSmfInfoList": [
			 {
			  "sNssai": {
				 "sst": 2,
                  "sd": "A00000"
			  },
			  "dnnSmfInfoList":[
			    {
				   "dnn": "smf-dnn1"
				}
			  ]
			 }
			]
        },
        "pcfInfo": {
            "dnnList": ["pcf-dnn1","pcf-dnn2"]
        },
        "bsfInfo": {
            "ipv4AddressRanges": [{
                "start": "172.16.208.0",
                "end": "172.16.208.255"
            }],
            "ipv6PrefixRanges": [{
                "start": "2001:db8:abcd:12::0/64",
                "end": "2001:db8:abcd:12::0/64"
            }]
        },
        "nfServices": [{
            "serviceInstanceId": "nudm-auth-02",
            "serviceName": "nudm-auth-02",
            "versions": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "scheme": "https://",
			"nfServiceStatus": "REGISTED",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "1001:da8::36",
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
            "serviceInstanceId": "nudm-ausf-02",
            "serviceName": "nudm-ausf-02",
			"nfServiceStatus": "REGISTED",
            "versions": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "scheme": "https://",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "1001:da8::36",
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
            "supportedFeatures": "A2"
        }]
    }]
}
`)

var contentNoValidityPeriod = []byte(`{
		"nfInstances": [
			   {
			     "nfInstanceId": "udm-5g-01",
			     "nfType": "udm",
				 "nfStatus": "REGISTERED",
			     "plmn": {
			       "mcc": "001",
			       "mnc": "001"
			     },
			     "sNssais": [
			       {
			         "sst": 1,
			         "sd": "A00001"
			       },
					{
					 "sst": 11,
			         "sd": "A00011"
					}
			     ],
			     "fqdn": "seliius03696.seli.gic.ericsson.se",
			     "ipv4Addresses": [
			       "127.0.0.1"
			     ],
			     "ipv6Addresses": [
			       "FF01::1101"
			     ],
			     "ipv6Prefixes": [
			       "2001:db8:abcd:12::0/64"
			     ],
			     "capacity": 1,
				 "priority":1,
			     "udrInfo": {
			       "supiRanges": [
			         {
			           "start": "000001",
			           "end": "000010",
			           "pattern": "(^imsi-[0-9]{5,15}$)|(^nai-.+$)"
			         }
			       ]
			     },
			     "amfInfo": {
			       "amfSetId": "amf-01"
			     },
			     "smfInfo": {
			       "dnnList": [
			         "udm-dnn-01",
					 "udm-dnn-011"
			       ],
			       "servingArea": [
			         "udm-servingArea-01"
			       ]
			     },
			     "upfInfo": {
			       "sNssaiUpfInfoList": [
			         {
			           "sNssai": {
			             "sst": 0,
			             "sd": "A00001"
			           },
			           "dnnUpfInfoList": [
			             {
			               "dnn": "udm-dnn"
			             }
			           ]
			         }
			       ]
			     }
			   }
			]
	}`)

var contentWithoutInstance = []byte(`{
		"validityPeriod": 200,
		"nfInstances":
	}`)

var testEvent = "testEvent00"
var testCfgName = "cfgName00"
var testFormat = cmproxy.NtfFormatFull

//var randnum = 1

func StubHTTPDo(resp *httpclient.HttpRespData, err error) {
	client.HTTPDo = func(httpv, method, url string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
		return resp, err
	}
}

func StubGetLeader(ID string) {
	election.GetLeader = func() string {
		return ID
	}
}

func TestCmTargetNfProfilesHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	event := "testEvent"
	cfgName := "cfgName"

	//test format value
	format := cmproxy.NtfFormatPatch
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	_, ok := common.CmGetTargetNfProfile()
	if ok {
		t.Errorf("TestCmTargetNfProfilesHandler: CmTargetNfProfilesHandler format check failure.")
	}

	//test wrong data input
	format = cmproxy.NtfFormatFull
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfileWrong)
	_, ok = common.CmGetTargetNfProfile()
	if ok {
		t.Errorf("TestCmTargetNfProfilesHandler: CmTargetNfProfilesHandler wrong data check failure.")
	}

	//normal process
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	targetNfs, ok := common.CmGetTargetNfProfile()
	if targetNfs[0].TargetNfType != "UDM" || targetNfs[0].TargetServiceNames[0] != "nudm-auth-01" || !ok {
		t.Errorf("TestCmNrfServiceProfilesHandler: CmNrfServiceProfilesHandler UpdateTargetNfProfiles format check failure.")
	}
	if targetNfs[1].TargetNfType != "UDR" || targetNfs[1].TargetServiceNames[0] != "nudr-01" || !ok {
		t.Errorf("TestCmNrfServiceProfilesHandler: CmNrfServiceProfilesHandler UpdateTargetNfProfiles format check failure.")
	}
}

func TestWaitReady(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	os.Setenv("POD_IP", "127.0.0.1")
	StubGetLeader("127.0.0.1")

	event := "testEvent"
	cfgName := "cfgName"

	format := cmproxy.NtfFormatFull
	cmNrfAgentConfHandler(event, cfgName, format, cmDataNrfService)
	waitCMReady()
	//	if !ok {
	//		t.Errorf("TestWaitReady: WaitReady check failure.")
	//	}
}

func TestFetchProfilesByInstanceId(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestCached: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestRequesterNfMessageBusHandler: Cached fail")
		}
	}

	event := "testEvent"
	cfgName := "cfgName"

	client.InitHttpClient()

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	//test format value
	format := cmproxy.NtfFormatFull
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	targetNfSet, ok := common.CmGetTargetNfProfile()
	if !ok {
		t.Errorf("TestFetchProfilesByInstanceId: CmTargetNfProfilesHandler format check failure.")
	}
	cacheManager.SetTargetNf("AUSF", targetNfSet[0])
	cacheManager.SetRequesterFqdn("AUSF", "request_fqdn")
	cmNrfAgentConfHandler(event, cfgName, format, cmDataNrfService)

	handleDiscoveryRequest(&targetNfSet[0], "udm-5g-01")
}

//func TestHandleNewProfile(t *testing.T) {
//	log.SetLevel(log.ErrorLevel)
//	var cm *cache.CacheManager
//	cm.SetCasheConfigUT("../../build/config/cache-index.json")
//	cm = cache.Instance()

//	nfinstanceByte, validityPeriod, ok := cache.SpliteSeachResult(searchResultUDM)
//	if nfinstanceByte == nil || !ok {
//		t.Errorf("TestHandleNewProfile: SpliteSeachResult fail")
//	}
//	for _, instance := range nfinstanceByte {
//		handleNewProfile("AUSF", instance, validityPeriod)
//	}
//	ok = cm.Probe("AUSF", "udm-5g-01")
//	if !ok {
//		t.Errorf("TestHandleNewProfile: HandleNewProfile check fail")
//	}
//	cm.Flush("AUSF")
//}

func TestGetSearchInCache(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, ok := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil || !ok {
		t.Errorf("TestGetSearchInCache: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}
	searchConditions := cache.SearchParameter{}
	searchConditions.SetServiceNames([]string{"nudm-auth-01"})
	searchConditions.SetTargetNfType("UDM")
	searchConditions.SetRequesterNfType("UDM")
	content, ok := cache.Instance().Search("AUSF", "UDM", &searchConditions, false)
	if !ok || len(content) == 0 {
		t.Errorf("TestGetSearchInCache: Search failed")
	}

	cacheManager.Flush("AUSF")
}

//func TestCacheProvisioningHandler(t *testing.T) {
//	log.SetLevel(log.ErrorLevel)

//	resp := httptest.NewRecorder()

//	req := httptest.NewRequest("PUT", "/nfName/nrf-discovery-agent/v1/memcache", bytes.NewBuffer(contentNoValidityPeriod))
//	req.Header.Set("Content-Type", "application/json")

//	cacheProvisioningHandler(resp, req)
//	ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
//	if ok {
//		t.Errorf("TestCacheProvisioningHandler: cacheProvisioningHandler no ValidityPeriod check fail")
//	}

//	req = httptest.NewRequest("PUT", "/nfName/nrf-discovery-agent/v1/memcache", bytes.NewBuffer(contentWithoutInstance))
//	req.Header.Set("Content-Type", "application/json")
//	cacheProvisioningHandler(resp, req)
//	ok = cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
//	if ok {
//		t.Errorf("TestCacheProvisioningHandler: cacheProvisioningHandler no Instance check fail")
//	}

//	req = httptest.NewRequest("PUT", "/nfName/nrf-discovery-agent/v1/memcache", bytes.NewBuffer(searchResultUDM))
//	req.Header.Set("Content-Type", "application/json")
//	resp2 := httptest.NewRecorder()
//	resp2.Code = http.StatusInternalServerError
//	cacheProvisioningHandler(resp2, req)
//	_ = cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
//	if resp2.Code != http.StatusCreated {
//		t.Errorf("TestCacheProvisioningHandler: cacheProvisioningHandler response code %d check fail", resp.Code)
//	}

//	cacheManager.Flush("AUSF")
//}

func TestDiscoverRequestHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	cacheManager.Flush("AUSF")
	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached fail")
		}
	}

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	})
	cacheManager.SetRequesterFqdn("AUSF", "AUSF.ericsson.com")

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	nfDiscoveryRequestHandler(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestDiscoverRequestHandler: nfDiscoveryRequestHandler response code %d check fail", resp.Code)
	}

	cacheManager.Flush("AUSF")
	cacheManager.DeleteTargetNf("AUSF")
}

func TestDiscoverRequestHandlerForTargetPlmnList(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	util.PreComplieRegexp()
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	cacheManager.Flush("AUSF")
	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached fail")
		}
	}

	nfinstanceByte, _, _ = cache.SpliteSeachResult(searchResultUDM2)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult udm2 fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached umd2 fail")
		}
	}

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	})
	cacheManager.SetRequesterFqdn("AUSF", "AUSF.ericsson.com")

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com&target-plmn-list=[{\"mcc\":\"460\",\"mnc\":\"001\"},{\"mcc\":\"460\",\"mnc\":\"002\"}]&target-plmn-list={\"mcc\":\"460\",\"mnc\":\"003\"}", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	nfDiscoveryRequestHandler(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestDiscoverRequestHandler: nfDiscoveryRequestHandler response code %d check fail", resp.Code)
	}

	cacheManager.Flush("AUSF")
	cacheManager.DeleteTargetNf("AUSF")
}

func TestDiscoverRequestHandlerForNsiList(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	cacheManager.Flush("AUSF")
	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached fail")
		}
	}

	nfinstanceByte, _, _ = cache.SpliteSeachResult(searchResultUDM2)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult udm2 fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached umd2 fail")
		}
	}

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01", "nudm-auth-02"},
	})
	cacheManager.SetRequesterFqdn("AUSF", "AUSF.ericsson.com")

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	//req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01,nudm-auth-02&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com&nsi-list=069,102&nsi-list=070,103&nsi-list=071", bytes.NewBuffer(nobody))
	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com&nsi-list=069,102&nsi-list=070,103&nsi-list=071", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	nfDiscoveryRequestHandler(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestDiscoverRequestHandler: nfDiscoveryRequestHandler response code %d check fail", resp.Code)
	}

	cacheManager.Flush("AUSF")
	cacheManager.DeleteTargetNf("AUSF")
}

func TestDiscoverRequestHandlerForGroupID(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	cacheManager.Flush("AUSF")
	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached fail")
		}
	}

	nfinstanceByte, _, _ = cache.SpliteSeachResult(searchResultUDM2)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult udm2 fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached umd2 fail")
		}
	}

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	})
	cacheManager.SetRequesterFqdn("AUSF", "AUSF.ericsson.com")

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com&group-id-list=udm-01,udm-02", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	nfDiscoveryRequestHandler(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestDiscoverRequestHandler: nfDiscoveryRequestHandler response code %d check fail", resp.Code)
	}

	cacheManager.Flush("AUSF")
	cacheManager.DeleteTargetNf("AUSF")
}

func TestDiscoverRequestHandlerForTargetSnssais(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	cacheManager.Flush("AUSF")
	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached fail")
		}
	}

	nfinstanceByte, _, _ = cache.SpliteSeachResult(searchResultUDM2)
	if nfinstanceByte == nil {
		t.Errorf("TestDiscoverRequestHandler: SpliteSeachResult udm2 fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestDiscoverRequestHandler: Cached umd2 fail")
		}
	}

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	})
	cacheManager.SetRequesterFqdn("AUSF", "AUSF.ericsson.com")

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01&requester-nf-type=AUSF&requester-nf-instance-fqdn=AUSF.ericsson.com&snssais=[{\"sst\":2,\"sd\":\"A00000\"},{\"sst\":4,\"sd\":\"A00000\"}]&snssais={\"sst\":4,\"sd\":\"A00001\"}", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	nfDiscoveryRequestHandler(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestDiscoverRequestHandler: nfDiscoveryRequestHandler response code %d check fail", resp.Code)
	}

	cacheManager.Flush("AUSF")
	cacheManager.DeleteTargetNf("AUSF")
}

func TestResponseCarryDataHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	//log.SetLevel(log.DebugLevel)
	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusInternalServerError

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/nf-instances?target-nf-type=UDM&service-names=nudm-auth-01&target-nf-type=udm", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	body := []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
        "plmn": {
            "mcc": "460",
            "mnc": "000"
        }
		}]
	}`)

	sequenceID := utils.GetSequenceId()
	logcontent := &log.LogStruct{SequenceId: sequenceID}
	handleDiscoverySuccess(resp, req, logcontent, http.StatusOK, 86400, string(body))
	if resp.Code != http.StatusOK {
		t.Errorf("TesthandleNfDiscoveryResponse: handleNfDiscoveryResponse response code check fail")
	}
}

func stubNotifIPEndPoint() bool {
	ipEndPoint := structs.StatusNotifIPEndPoint{
		Ipv4Address: "192.168.110.112",
		Port:        12345,
	}

	ipEndPointData, err := json.Marshal(ipEndPoint)
	if err != nil {
		log.Errorf("stubNotifIPEndPoint: failed to Marshal StatusNotifIPEndPoint, %s", err.Error())
		return false
	}

	return structs.UpdateStatusNotifIPEndPoint(ipEndPointData)
}

/*
func TestSubscribeByNfType(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	client.InitHttpClient()

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("POST", http.StatusCreated)

	//	var cacheManager *cache.CacheManager
	//	cacheManager.SetCasheConfigUT("../../build/config/cache-index.json")
	//	cacheManager = cache.Instance()

	//	subscribeURL := "http://127.0.1.1:3212/nnrf-nfm/v1/nf-instances/udm-5g-01/"

	//	backupCmGetNrfBaseURL := cmGetNrfBaseURL
	//	defer func() {
	//		cmGetNrfBaseURL = backupCmGetNrfBaseURL
	//	}()
	//	cmGetNrfBaseURLStub(subscribeURL)

	cmNrfAgentConfHandler(testEvent, testCfgName, testFormat, cmDataNrfServicePort3212)

	//cm.Opts.SubsCallback = "http://10.111.14.85:80/nrf-notify-ntf/v1/notify/ausf"
	event := "testEvent"
	cfgName := "cfgName"
	format := cmproxy.NtfFormatFull
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	targetNfs, ok := common.CmGetTargetNfProfile()
	if !ok {
		t.Errorf("TestSubscribeByNfType: CmGetTargetNfProfile format check failure.")
	}

	subscriptionTimer = timer.NewTimer()
	go func() {
		for s := range subscriptionTimer.TimerChan() {
			t.Logf("TestSubscribeByNfType: %s", s)
		}
	}()

	t.Run("TestSubscribeByNfType", func(t *testing.T) {
		ok = subscribeTargetNf(&targetNfs[0])
		if !ok {
			t.Errorf("TestSubscribeByNfType: SubscribeByNfType failure.")
		}

		subscriptionID, _ := getSubscriptionID(targetNfs[0].RequesterNfType, targetNfs[0].TargetNfType, targetNfs[0].TargetServiceNames[0])
		if subscriptionID == "" {
			t.Errorf("TestSubscribeByNfType: SubscribeByNfType failure.")
		}

		unsubscribeByNfType(targetNfs[0].RequesterNfType)
	})
	t.Run("TestSubscribeByNfType RequesterNfDeregisteredHandler", func(t *testing.T) {
		ok = subscribeTargetNf(&targetNfs[0])
		if !ok {
			t.Errorf("TestRequesterNfDeregisteredHandler: SubscribeByNfType failure.")
		}
		subscriptionID, _ := getSubscriptionID(targetNfs[0].RequesterNfType, targetNfs[0].TargetNfType, targetNfs[0].TargetServiceNames[0])
		if subscriptionID == "" {
			t.Errorf("TestSubscribeByNfType: SubscribeByNfType failure.")
		}

		subscribeMessageBusHandler(consts.MsgbusTopicNamePrefix+"subscibe", msgRegEvent)
		var messageData structs.RegMsgBus
		err := json.Unmarshal(msgDeRegEvent, &messageData)
		if err != nil {
			t.Errorf("TestRequesterNfDeregisteredHandler: msgRegEvent unmarshal failure. %s", err.Error())
		}
		ok = nfSubscribeDeregisterEventHandler(&messageData)
		if !ok {
			t.Errorf("TestRequesterNfDeregisteredHandler: deregister event handle failure.")
		}

		unsubscribeByNfType(targetNfs[0].RequesterNfType)
	})
}
*/

func StubTargetNfProfilesHandler(Event, ConfigurationName string, RawData []byte) {
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	fmt.Println("StubTargetNfProfilesHandler:", Event, ConfigurationName, string(RawData))
	cmTargetNfProfilesHandler(Event, ConfigurationName, cmproxy.NtfFormatFull, RawData)
}

func TestCreateFsNotify(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	dirPath := "./test/"
	err := os.Mkdir(dirPath, 0777)
	if err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	defer func() { _ = os.Remove(dirPath) }()

	filePath := filepath.Join(dirPath, "test.json")
	defer func() { _ = os.Remove(filePath) }()

	go func() {
		configmapMonitor(dirPath, 30*time.Second, StubTargetNfProfilesHandler)
	}()

	time.Sleep(2 * time.Second)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	wr, _ := file.WriteString("test data")
	fmt.Println("wrote bytes: ", wr)

	time.Sleep(2 * time.Second)
}

func TestCreateFsTimer(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	dirPath := "./test/"
	err := os.Mkdir(dirPath, 0777)
	if err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	defer func() { _ = os.Remove(dirPath) }()

	filePath := filepath.Join(dirPath, "test.json")
	defer func() { _ = os.Remove(filePath) }()

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	wr, _ := file.Write(cmTargetNfProfile)
	fmt.Println("wrote bytes: ", wr)

	go func() {
		configmapMonitor(dirPath, 1*time.Second, StubTargetNfProfilesHandler)
	}()
	time.Sleep(2 * time.Second)

}

/*
func stupSubscriptionTimer() {
	//spt := timer.NewTimer()
	//setSubscriptionTimer(spt)
}
*/
/*
func TestTargetNfProfilesHandler(t *testing.T) {
	client.InitHttpClient()

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("POST", http.StatusCreated)

	var cacheManager *cache.CacheManager
	cacheManager.SetCasheConfigUT("../../build/config/cache-index.json")
	cacheManager = cache.Instance()

	//	subscribeURL := "http://127.0.1.1:3215/nnrf-nfm/v1/nf-instances/udm-5g-01/"

	//	backupCmGetNrfBaseURL := cmGetNrfBaseURL
	//	defer func() {
	//		cmGetNrfBaseURL = backupCmGetNrfBaseURL
	//	}()
	//	cmGetNrfBaseURLStub(subscribeURL)

	cmNrfAgentConfHandler(testEvent, testCfgName, testFormat, cmDataNrfService)

	stupSubscriptionTimer()

	//cm.Opts.SubsCallback = "http://10.111.14.85:80/nfname/nrf-notify-ntf/v1/notify"
	event := "CREATE"
	cfgName := "cfgName"
	configmapTargetNfProfilesHandler(event, cfgName, cmTargetNfProfile)
	targetNfs, ok := common.CmGetTargetNfProfile()
	if !ok {
		t.Errorf("TestTargetNfProfilesHandler: CmGetTargetNfProfile format check failure.")
	}
	for _, targetNf := range targetNfs {
		ok = subscribeTargetNf(&targetNf)
		if !ok {
			t.Errorf("TestTargetNfProfilesHandler: subscribeTargetNf failure.")
		}
		subscriptionID, _ := getSubscriptionID(targetNf.RequesterNfType, targetNf.TargetNfType, targetNf.TargetServiceNames[0])
		if subscriptionID == "" {
			t.Errorf("TestSubscribeByNfType: SubscribeByNfType failure.")
		}
	}
	ok, _ = cacheManager.GetSubscriptionIDURLs("AUSF")
	//if !ok || len(URLSet) != 2 {
	if !ok {
		t.Errorf("TestTargetNfProfilesHandler: GetSubscriptionIDURLs check URL failure.")
	}

	for _, targetNf := range targetNfs {
		unsubscribeByNfType(targetNf.RequesterNfType)
	}
}
*/

func TestFetchRequesterNfInfo(t *testing.T) {
	fHTTPDo := client.HTTPDo
	defer func() { client.HTTPDo = fHTTPDo }()

	cache.Instance().DeleteRequesterFqdn("AUSF")

	t.Run("TestFetchRequesterNfInfo01", func(t *testing.T) {
		var resp httpclient.HttpRespData
		resp.StatusCode = http.StatusOK
		resp.Body = []byte(`[{"nfInstanceId":"nef_instance","requesterNfFqdn":"fqdn.orig.ericsson.com","requesterNfType":"NEF"},{"nfInstanceId":"ausf_instance","requesterNfFqdn":"fqdn.orig.ericsson.com","requesterNfType":"AUSF"}]`)
		StubHTTPDo(&resp, nil)

		msg, ok := fetchRequesterNfInfo()
		//fmt.Printf("msg:%+v\n", msg)
		if msg == nil || ok == false {
			t.Errorf("TestFetchRequesterNfInfo: TestFetchRequesterNfInfo01 failure")
		}
	})

	t.Run("TestFetchRequesterNfInfo02", func(t *testing.T) {
		var resp httpclient.HttpRespData
		resp.StatusCode = http.StatusOK
		StubHTTPDo(&resp, nil)

		msg, ok := fetchRequesterNfInfo()
		if msg != nil || ok == false {
			t.Errorf("TestFetchRequesterNfInfo: TestFetchRequesterNfInfo02 failure")
		}
	})
}

func TestLoopFetchRequesterNfInfo(t *testing.T) {
	fHTTPDo := client.HTTPDo
	defer func() { client.HTTPDo = fHTTPDo }()

	t.Run("TestLoopFetchRequesterNfInfo", func(t *testing.T) {
		var resp httpclient.HttpRespData
		resp.StatusCode = http.StatusOK
		resp.Body = []byte(`[{"nfInstanceId":"nef_instance","requesterNfFqdn":"fqdn.orig.ericsson.com","requesterNfType":"NEF"},{"nfInstanceId":"ausf_instance","requesterNfFqdn":"fqdn.orig.ericsson.com","requesterNfType":"AUSF"}]`)
		StubHTTPDo(&resp, nil)

		fqdnInfoSlice := loopFetchRequesterNfInfo()
		fmt.Printf("fqdnInfoSlice is %+v\n", fqdnInfoSlice)
		if fqdnInfoSlice == nil {
			t.Errorf("TestLoopFetchRequesterNfInfo failure")
		}
	})

	t.Run("TestLoopFetchRequesterNfInfo02", func(t *testing.T) {
		var resp httpclient.HttpRespData
		resp.StatusCode = http.StatusOK
		StubHTTPDo(&resp, nil)

		fqdnInfoSlice := loopFetchRequesterNfInfo()
		fmt.Printf("fqdnInfoSlice is %+v\n", fqdnInfoSlice)
		if fqdnInfoSlice != nil {
			t.Errorf("TestLoopFetchRequesterNfInfo failure")
		}
	})
}

//func TestSubscriptionTrigger(t *testing.T) {
//	cache.Instance().SetRequesterFqdn("AUSF", "seliius03690.seli.gic.ericsson.se")
//	var messageData *structs.RegMsgBus = new(structs.RegMsgBus)
//	messageData.NfType = "AUSF"
//	messageData.NfInstanceID = "ausf-01"
//	messageData.FQDN = "seliius03690.seli.gic.ericsson.se"

//	ok := subscriptionTrigger()
//	if ok != true {
//		t.Errorf("TestSubscriptionTrigger: check failure")
//	}
//	cache.Instance().DeleteRequesterFqdn("AUSF")
//}

/*
func TestGetSubscriptionID(t *testing.T) {
	var subInfo structs.SubscriptionInfo
	subInfo.RequesterNfType = "AUSF"
	subInfo.TargetNfType = "UDM"
	subInfo.TargetServiceName = "udm-01"
	subInfo.SubscriptionID = "123-456-789"
	subInfo.ValidityTime = time.Time{}

	subscriptionInfoMap["AUSF"] = subInfo
	sid, _ := getSubscriptionID("AUSF", "UDM", "udm-01")
	if sid != "123-456-789" {
		t.Errorf("TestSubscriptionTrigger: check failure")
	}
	delete(subscriptionInfoMap, "AUSF")
}
*/
/*
func TestDoSubscriptionToNRF(t *testing.T) {
	cache.Instance().DeleteRequesterFqdn("AUSF")
	var oneSubData structs.OneSubscriptionData
	oneSubData.RequesterNfType = "AUSF"
	oneSubData.TargetNfType = "UDM"
	oneSubData.TargetServiceName = "udm-01"
	subscriptionID, _, _ := doSubscriptionToNRF(&oneSubData)
	if subscriptionID != "" {
		t.Errorf("TestDoSubscriptionToNRF: check failure")
	}
}
*/
/*
func TestForwardNfDiscoveryRequest(t *testing.T) {
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() { client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc }()

	StubSetup()

	targetNf := &structs.TargetNf{}
	targetNf.RequesterNfType = "AUSF"
	targetNf.TargetNfType = "UDM"
	targetNf.TargetServiceNames = []string{"nudm-ausf-01"}
	cache.Instance().SetTargetNf("AUSF", *targetNf)
	cache.Instance().SetReqFqdn("AUSF", "seliius19953")

	resp := &httpclient.HttpRespData{}
	resp.StatusCode = http.StatusOK
	resp.Body = searchResultUDM

	t.Run("TestForwardNfDiscoveryRequest", func(t *testing.T) {
		req := &http.Request{}
		req.URL, _ = url.Parse("http://127.0.0.1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			if url_suffix != "nf-instances?service-names=nudm-ausf-01&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius19953" {
				t.Errorf("fail to forward")
				return nil, fmt.Errorf("fail to forward")
			}
			return resp, nil
		}

		_, err := forwardNfDiscoveryRequest(targetNf, req)
		if err != nil {
			t.Errorf("fail to forward")
		}
	})
	t.Run("TestForwardNfDiscoveryRequest", func(t *testing.T) {
		req := &http.Request{}
		req.URL, _ = url.Parse("http://127.0.0.1/nf-instances?service-names=nudm-sdm&service-names=nudm-test100&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UDM&requester-nf-type=UDR&supi=imsi-600<imsi>9999&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%22000000%22%7D&snssais=%7B%22sst%22%3A+1,%22sd%22%3A+%22111111%22%7D")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			if url_suffix != "nf-instances?service-names=nudm-ausf-01&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius19953&snssais=%7B%22sst%22%3A+0%2C%22sd%22%3A+%22000000%22%7D&snssais=%7B%22sst%22%3A+1%2C%22sd%22%3A+%22111111%22%7D&supi=imsi-600%3Cimsi%3E9999" {
				t.Errorf("fail to forward")
				return nil, fmt.Errorf("fail to forward")
			}
			return resp, nil
		}

		_, err := forwardNfDiscoveryRequest(targetNf, req)
		if err != nil {
			t.Errorf("fail to forward")
		}
	})
}
*/
/*
func TestSubscriptionTimerHandler(t *testing.T) {
	fHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	defer func() { client.HTTPDoToNrfMgmt = fHTTPDoToNrfMgmt }()
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() { client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc }()

	StubSetup()

	subscriptionID := "20bd0bb9-edc1-4c74-8ec5-74e4fed79ac8"

	s := &structs.SubscriptionInfo{}
	s.RequesterNfType = "AUSF"
	s.TargetNfType = "UDM"
	s.TargetServiceName = "nudm-ausf-01"
	s.SubscriptionID = subscriptionID
	s.ValidityTime = time.Time{}

	resp := &httpclient.HttpRespData{}
	resp.StatusCode = http.StatusCreated
	resp.Location = "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/" + subscriptionID

	t.Run("TestSubscriptionTimerHandler", func(t *testing.T) {
		subscriptionTimerHandler(subscriptionID)
	})

	t.Run("TestSubscriptionTimerHandler", func(t *testing.T) {
		s.SubscriptionID = subscriptionID
		s.ValidityTime = time.Now()
		updateConfigmapStorage(s, subscriptionID)
		defer func() {
			s.SubscriptionID = ""
			updateConfigmapStorage(s, subscriptionID)
		}()

		client.HTTPDoToNrfMgmt = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			if method == "PUT" {
				return nil, errors.New("test error")
			}
			return resp, nil
		}

		subscriptionTimerHandler(subscriptionID)
		_, existed := getSubscriptionInfo(subscriptionID)
		if existed {
			t.Errorf("")
			return
		}

		updateConfigmapStorage(s, subscriptionID)
		resp.StatusCode = http.StatusNotFound
		client.HTTPDoToNrfMgmt = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			if method == "PUT" {
				resp.StatusCode = http.StatusNotFound
				return resp, nil
			}
			resp.StatusCode = http.StatusCreated
			return resp, nil
		}

		subscriptionTimerHandler(subscriptionID)
		_, existed = getSubscriptionInfo(subscriptionID)
		if existed {
			t.Errorf("")
		}
	})

	t.Run("TestSubscriptionTimerHandler", func(t *testing.T) {
		s.SubscriptionID = subscriptionID
		s.ValidityTime = time.Now()
		updateConfigmapStorage(s, subscriptionID)
		defer func() {
			s.SubscriptionID = ""
			updateConfigmapStorage(s, subscriptionID)
		}()

		client.HTTPDoToNrfMgmt = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			return resp, nil
		}

		subscriptionTimerHandler(subscriptionID)
		_, existed := getSubscriptionInfo(subscriptionID)
		if existed {
			t.Errorf("")
		}
	})
}
*/
func TestFixSearchResultFromNrf(t *testing.T) {
	t.Run("filterOneSvc", func(t *testing.T) {
		var searchParameter = &cache.SearchParameter{}
		searchParameter.SetServiceNames([]string{"nudm-ausf-01"})
		searchParameter.SetSupportedFeatures("")

		data, e := applyFilter(searchResultUDM, searchParameter)
		if e != nil {
			t.Errorf("%s", e.Error())
		} else {
			if !bytes.Contains(data, []byte("nudm-ausf-01")) ||
				bytes.Contains(data, []byte("nudm-auth-01")) {
				t.Errorf("%s", string(data))
			}
		}
	})
	t.Run("filterAllSvc", func(t *testing.T) {
		var searchParameter = &cache.SearchParameter{}
		searchParameter.SetServiceNames([]string{"nudm-auth-x"})
		searchParameter.SetSupportedFeatures("")

		data, e := applyFilter(searchResultUDM, searchParameter)
		if e == nil {
			t.Errorf("%s", string(data))
		}
	})
	t.Run("filterBySF", func(t *testing.T) {
		var searchParameter = &cache.SearchParameter{}
		searchParameter.SetServiceNames(nil)
		searchParameter.SetSupportedFeatures("01")

		data, e := applyFilter(searchResultUDM, searchParameter)
		if e != nil {
			t.Errorf("%s", e.Error())
		} else {
			if bytes.Contains(data, []byte("nudm-ausf-01")) ||
				!bytes.Contains(data, []byte("nudm-auth-01")) {
				t.Errorf("%s", string(data))
			}
		}
	})
	t.Run("filterAllBySF", func(t *testing.T) {
		var searchParameter = &cache.SearchParameter{}
		searchParameter.SetServiceNames(nil)
		searchParameter.SetSupportedFeatures("00")

		data, e := applyFilter(searchResultUDM, searchParameter)
		if e == nil {
			t.Errorf("%s", string(data))
		}
	})
}

/*
func TestBuildSubscriptionPatchData(t *testing.T) {
	validityTime := 86400
	patchData := buildSubscriptionPatchData(validityTime)
	if patchData == nil {
		t.Errorf("Expect patchData is not nil, but failure")
	}

	t.Logf("patchData : %+v\n", string(patchData))
}
*/
func TestBuildSubscriptionData(t *testing.T) {
	/*
		serviceNameCond := structs.ServiceNameCond{
			ServiceName: "nudr-dr",
		}
	*/
	oneSubsData := structs.OneSubscriptionData{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "namf-comm",
		NotifCondition:    nil,
	}

	fmt.Printf("%+v\n", oneSubsData)
}

func TestWaitDISCAgentReady(t *testing.T) {
	go func() {
		nfIsReady <- true
	}()
	waitDISCAgentReady()
}

/*
func TestHandleDiscoveryRequestForReadiness(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestHandleDiscoveryRequestForReadiness: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestHandleDiscoveryRequestForReadiness: Cached fail")
		}
	}

	event := "testEvent"
	cfgName := "cfgName"

	client.InitHttpClient()

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)

	//test format value
	format := cmproxy.NtfFormatFull
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	targetNfSet, ok := common.CmGetTargetNfProfile()
	if !ok {
		t.Errorf("TestHandleDiscoveryRequestForReadiness: CmTargetNfProfilesHandler format check failure.")
	}
	cacheManager.SetTargetNf("AUSF", targetNfSet[0])
	cacheManager.SetRequesterFqdn("AUSF", "request_fqdn")
	cmNrfAgentConfHandler(event, cfgName, format, cmDataNrfService)

	handleDiscoveryRequestForReadiness(&targetNfSet[0], "udm-5g-01")
}
*/
func TestGetReadinessCheckFlag(t *testing.T) {
	flag := getReadinessCheckFlag()
	if flag != true {
		t.Errorf("TestHandleDiscoveryRequestForReadiness:check failure.")
	}
}

func TestSetReadinessCheckFlag(t *testing.T) {
	setReadinessCheckFlag(false)
	if needReadinessCheck != false {
		t.Errorf("TestSetReadinessCheckFlag:check failure.")
	}
	setReadinessCheckFlag(true)
	if getReadinessCheckFlag() != true {
		t.Errorf("TestSetReadinessCheckFlag:check failure.")
	}
}

func TestHandleDiscovery(t *testing.T) {
	log.SetLevel(log.LevelUint("DEBUG"))

	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusOK

	req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/memcache/AUSF-roam", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")
	logcontent := &log.LogStruct{SequenceId: "pod-1"}
	logcontent.ResponseDescription = "Handle cache"

	t.Run("TestHandleCacheOperationFailure", func(t *testing.T) {
		handleDiscoveryFailure(resp, req, logcontent, 404, "")
	})

}

func TestSetup(t *testing.T) {
	log.SetLevel(log.LevelUint("DEBUG"))
	Setup()
}

func TestAgentRoleMonitor(t *testing.T) {
	log.SetLevel(log.LevelUint("DEBUG"))
	startAgentRoleMonitor(10)
	agentRoleMonitorStarted = true
	StopAgentRoleMonitor()
}

func TestGetAgentRole(t *testing.T) {
	log.SetLevel(log.LevelUint("DEBUG"))
	getAgentRole()
}

func TestConfigmapTargetNfProfilesHandler(t *testing.T) {
	log.SetLevel(log.LevelUint("DEBUG"))

	event, configurationName := "CREATE", "ausf"
	ausfTargetProfile := []byte(`{
      "targetNfProfiles": [
        {
          "requesterNfType": "AUSF",
          "targetNfType": "UDM",
          "targetServiceNames": [
            "nudm-uecm"
          ]
        }
      ]
    } `)
	configmapTargetNfProfilesHandler(event, configurationName, ausfTargetProfile)
}

func TestProxy(t *testing.T) {
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() { client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc }()

	log.SetLevel(log.LevelUint("DEBUG"))

	t.Run("TestProxy_404", func(t *testing.T) {
		resp := &httpclient.HttpRespData{}
		resp.StatusCode = http.StatusNotFound
		resp.Body = []byte("")

		var nobody = []byte("")
		req := httptest.NewRequest("GET", "http://10.107.152.30:3202/nrf-discovery-agent/v1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			return resp, nil
		}
		rw := httptest.NewRecorder()
		//rw.Code = http.StatusNotFound

		rest := proxy(rw, req) // just forward
		if rest != false {
			t.Errorf("TestProxy_404: test proxy return 404 failed")
		}
	})

	t.Run("TestProxy_Wrong_Query", func(t *testing.T) {
		resp := &httpclient.HttpRespData{}
		resp.StatusCode = http.StatusNotFound
		resp.Body = []byte("")

		var nobody = []byte("")
		req := httptest.NewRequest("GET", "http://10.107.152.30:3202/nrf-discovery-agent/v1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			return resp, errors.New("request error")
		}
		rw := httptest.NewRecorder()

		rest := proxy(rw, req) // just forward
		if rest != false {
			t.Errorf("TestProxy_Wrong_Query: test proxy return 404 failed")
		}
	})

	t.Run("TestProxy_200_Wrong_body", func(t *testing.T) {
		resp := &httpclient.HttpRespData{}
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("")

		var nobody = []byte("")
		req := httptest.NewRequest("GET", "http://10.107.152.30:3202/nrf-discovery-agent/v1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			return resp, nil
		}
		rw := httptest.NewRecorder()

		rest := proxy(rw, req) // just forward
		if rest != false {
			t.Errorf("TestProxy_200_Wrong_body: test proxy return 404 failed")
		}
	})

	t.Run("TestProxy_200_Correct_body", func(t *testing.T) {
		resp := &httpclient.HttpRespData{}
		resp.StatusCode = http.StatusOK
		resp.Body = searchResultUDM

		var nobody = []byte("")
		req := httptest.NewRequest("GET", "http://10.107.152.30:3202/nrf-discovery-agent/v1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
			return resp, nil
		}
		rw := httptest.NewRecorder()

		rest := proxy(rw, req) // just forward
		if rest != true {
			t.Errorf("TestProxy_200_Correct_body: test proxy return 404 failed")
		}
	})
}
func TestClose(t *testing.T) {
	var nobody = []byte("")
	req := httptest.NewRequest("GET", "http://10.107.152.30:3202/nrf-discovery-agent/v1/nf-instances?requester-nf-type=AUSF&target-nf-type=UDM", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")
	close(req)
}
