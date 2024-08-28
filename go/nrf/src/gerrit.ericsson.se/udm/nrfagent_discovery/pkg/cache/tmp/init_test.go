package cache

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
)

var contUdmReg = []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
        "plmn": {
            "mcc": "460",
            "mnc": "000"
        },
        "sNssais": [{
                "sst": 2,
                "sd": "2"
            },
            {
                "sst": 4,
                "sd": "4"
            }
        ],
        "nsiList": ["069"],
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
            "routingIndicators": ["1111", "1234", "5678"]
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
            "supportedFeatures": "A2"
        }]
    }]
}
`)

var contUdmRegRoam = []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-02",
        "nfType": "UDM",
        "plmn": {
            "mcc": "450",
            "mnc": "000"
        },
        "sNssais": [{
                "sst": 2,
                "sd": "2"
            },
            {
                "sst": 4,
                "sd": "4"
            }
        ],
        "nsiList": ["069"],
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
            "routingIndicators": ["1111", "1234", "5678"]
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
            "serviceInstanceId": "nudm-auth-02",
            "serviceName": "nudm-auth-02",
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
            "serviceInstanceId": "nudm-ausf-02",
            "serviceName": "nudm-ausf-02",
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
            "supportedFeatures": "A2"
        }]
    }]
}
`)

var (
	cacheManager *CacheManager
)

func TestMain(m *testing.M) {
	setupEnv()
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func setupEnv() {
	SetCacheConfig("../../build/config/cache-index.json")
	cacheManager = Instance()
	cacheManager.InitCache("AUSF", "UDM")
	activeLeaderMock(true)
}

func activeLeaderMock(isLeader bool) {
	election.IsActiveLeader = func(probePort, probeURL string) bool {
		return isLeader
	}
}

var nfProfileAUSF = []byte(`{"heartBeatTimer":120, "nfProfile":{
	"fqdn": "seliius03696.seli.gic.ericsson.se",
	"nfInstanceId": "0c765084-9cc5-49c6-9876-ae2f5fa2a63f",
	"nfServices": [
	  {
		"fqdn": "seliius03696.seli.gic.ericsson.se",
		"schema": "https://",
		"serviceInstanceId": "nausf-auth-01",
		"serviceName": "nausf-auth",
		"version": [
		  {
			"apiFullVersion": "1.R15.1.1 ",
			"apiVersionInUri": "v1",
			"expiry": "2020-07-06T02:54:32Z"
		  }
		]
	  }
	],
	"nfStatus": "REGISTERED",
	"nfType": "AUSF"
  }
  }`)

var subscrRsp = []byte(`{"validityTime":"2019-04-02T17:11:28+08:00", 
    "nfStatusNotificationUri": "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/subscriptionTest",
	"subscriptionId": "subscriptionTest"	
  }`)

var masterSubscrRsp = []byte(`{"validityTime":"2019-04-02T17:11:28+08:00", 
    "nfStatusNotificationUri": "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/subscriptionTest",
	"subscriptionID": "subscriptionTest",
	"requesterNfType": "AUSF",
	"targetNfType": "UDM"
  }`)

var searchResultUDM = []byte(`{
	  "validityPeriod": 43200,
	  "nfInstances": [{
		  "nfInstanceId": "udm-5g-01",
		  "nfType": "UDM",
		  "nfStatus": "REGISTERED",
		  "plmnList": [{
			  "mcc": "466",
			  "mnc": "000"
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
		  "nsiList": ["069","001","101"],
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
				"groupId": "udr-01"
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
			  "groupId": "udm-01"
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
			  "serviceInstanceId": "nudm-auth-01",
			  "serviceName": "nudm-auth-01",
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
			  "serviceInstanceId": "nudm-ausf-01",
			  "serviceName": "nudm-uecm",
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
		  },
		  {
			  "serviceInstanceId": "nudm-ausf-01",
			  "serviceName": "nudm-ausf-01",
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

func StubHTTPDoToNrf(method string, code int) {
	respMgmt := &httpclient.HttpRespData{}
	respMgmt.StatusCode = code
	if method == "POST" {
		subID, _ := utils.GetUUIDString()
		respMgmt.Location = "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/" + subID
		respMgmt.Body = subscrRsp
	} else if method == "PUT" {
		respMgmt.Body = nfProfileAUSF
	} else if method == "PATCHSubscr" {
		respMgmt.Body = subscrRsp
	}
	StubHTTPDoToNrfMgmt(respMgmt, nil)

	respDisc := &httpclient.HttpRespData{}
	respDisc.StatusCode = code
	if method == "GET" {
		respDisc.Body = searchResultUDM
		respDisc.Header = &http.Header{}
		respDisc.Header.Set("Etag", "etag_test")
	}

	StubHTTPDoToNrfDisc(respDisc, nil)
}

func StubHTTPDoToNrfMgmt(resp *httpclient.HttpRespData, err error) {
	client.HTTPDoToNrfMgmt = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
		return resp, err
	}
}

func StubHTTPDoToNrfDisc(resp *httpclient.HttpRespData, err error) {
	client.HTTPDoToNrfDisc = func(httpv, method, url_suffix string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
		return resp, err
	}
}

func StubHTTPDoToMaster(method string, code int) {
	respMaster := &httpclient.HttpRespData{}
	respMaster.StatusCode = code
	if method == "GETSubscr" {
		respMaster.Location = "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/SubscriptionTest"
		respMaster.Body = masterSubscrRsp
	}
	StubHTTPDo(respMaster, nil)
}

func StubHTTPDo(resp *httpclient.HttpRespData, err error) {
	client.HTTPDo = func(httpv, method, url string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
		return resp, err
	}
}
