package util

import (
	"log"
	"testing"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var udmProfile = []byte(`{
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
}
`)

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

func TestGetLeaderDiscURL(t *testing.T) {
	t.Log("Execute case TestGetLeaderDiscURL")
	discAgentLeaderURL := GetLeaderDiscURL()
	t.Logf("discAgentLeaderURL:%s", discAgentLeaderURL)
	if discAgentLeaderURL != "" {
		t.Fatal("Expect discURL is empty, but not")
	}
}

func TestGetDiscoveryRequestURL(t *testing.T) {
	t.Log("Execute case TestGetDiscoveryRequestURL")

	serviceNames := make([]string, 0)
	serviceNames = append(serviceNames, "udm-server1")
	serviceNames = append(serviceNames, "udm-server2")
	targetNf := &structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       serviceNames,
		NotifCondition:           nil,
		SubscriptionValidityTime: 86400,
	}

	fqdn := "seliius03696.seli.gic.ericsson.se"
	discURL := GetDiscoveryRequestURL(targetNf, "", fqdn)
	t.Logf("discURL:%s", discURL)

	expectURL := "nf-instances?service-names=udm-server1&service-names=udm-server2&target-nf-type=UDM&requester-nf-type=AUSF&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se"
	if discURL != expectURL {
		log.Fatalf("Expect discURL is %s, but not", expectURL)
	}
}

func TestGetStatusNotifURLs(t *testing.T) {
	t.Log("Execute case TestGetStatusNotifURLs")
	notifyURL := GetStatusNotifURLs()
	t.Logf("notifyURL:%s", notifyURL)
	if notifyURL != "" {
		t.Fatal("Expect notifyURL is empty, but not")
	}
}

func TestBuildSubscriptionPostData(t *testing.T) {
	t.Log("Execute case TestBuildSubscriptionPostData")

	subsData := &structs.OneSubscriptionData{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "UDM-server1",
		NotifCondition:    nil,
	}
	fqdn := "seliius03696.seli.gic.ericsson.se"

	subscriptionPostData := BuildSubscriptionPostData(subsData, fqdn)
	t.Logf("subscriptionPostData : %s", subscriptionPostData)
	if subscriptionPostData != nil {
		t.Fatal("Expect subscriptionPostData is nil, but not")
	}
}

func TestBuildSubscriptionPostRoamData(t *testing.T) {
	t.Log("Execute case TestBuildSubscriptionPostRoamData")

	subsData := &structs.OneSubscriptionData{
		RequesterNfType: "AUSF",
		TargetNfType:    "UDM",
		NfInstanceID:    "udm-5g-01",
		NotifCondition:  nil,
	}
	fqdn := "seliius03696.seli.gic.ericsson.se"
	plmnID := structs.PlmnID{
		Mcc: "450",
		Mnc: "000",
	}
	validityTime := "2019-04-02T17:11:28+08:00"

	subscriptionPostRoamData := BuildSubscriptionPostRoamData(subsData, fqdn, &plmnID, validityTime)
	t.Logf("subscriptionPostRoamData : %s", subscriptionPostRoamData)
	if subscriptionPostRoamData != nil {
		t.Fatal("Expect subscriptionPostRoamData is nil, but not")
	}
}

func TestBuildSubscriptionPatchData(t *testing.T) {
	t.Log("Execute case TestBuildSubscriptionPatchData")

	validityTime := 86400
	subscriptionPatchData := BuildSubscriptionPatchData(validityTime)
	t.Logf("subscriptionPatchData : %s", subscriptionPatchData)
	if subscriptionPatchData == nil {
		t.Fatal("Expect subscriptionPatchData is not nil, but not")
	}
}

func TestGetNfInstanceID(t *testing.T) {
	nfInstanceID := GetNfInstanceID(udmProfile)
	t.Logf("nfInstanceID : %s", nfInstanceID)
	if nfInstanceID == "" {
		t.Fatal("Expect get nfInstanceID success, but not")
	}
}

func TestGetNfInstances(t *testing.T) {
	empBody := make([]byte, 0)
	nfInstances, _ := GetNfInstances(empBody)
	t.Logf("NRF response body is nil")
	if nfInstances != nil {
		t.Fatal("Expect NRF response is empty, but not")
	}
	nfInstances, _ = GetNfInstances(searchResultUDM)
	t.Logf("nfInstances : %s", nfInstances)
	if nfInstances == nil {
		t.Fatal("Expect get nfInstances success, but not")
	}
}

func TestGetValidityPeriod(t *testing.T) {
	empBody := make([]byte, 0)
	validityPeriod, _ := GetValidityPeriod(empBody)
	t.Logf("NRF response body is nil")
	if validityPeriod != 0 {
		t.Fatal("Expect NRF response is empty, but not")
	}
	validityPeriod, _ = GetValidityPeriod(searchResultUDM)
	t.Logf("validityPeriod: %d", validityPeriod)
	if validityPeriod <= 0 {
		t.Fatal("Except get validityPeriod success, but not")
	}
}
