package profileop

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
)

var (
	newNfProfileData = []byte(`{
		"nfInstanceId": "5g-udm-01",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"plmn": {
		  "mcc": "460",
		  "mnc": "00"
		},
		"sNssais": [
		  {
			"sst": 0,
			"sd": "abAB01"
		  }
		],
		"fqdn": "seliius03696.seli.gic.ericsson.se",
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"ipv4Addresses": [
		  "172.16.208.1"
		],
		"ipv6Addresses": [
		  "FE80:1234::0000"
		],
		"capacity": 100,
		"load" : 100,
		"udmInfo": {
		  "groupId": "gid01",
		   "gpsiRanges": [
		   {
			 "start": "12300000",
			 "end": "12399999"
		   },
		   {
			"start": "12400000",
			"end": "12499999"
		   }
		],
		"supiRanges": [
		   {
			 "start": "12300000",
			 "end": "12399999"
		   },
		   {
			"start": "12400000",
			"end": "12499999"
		  }
		]
		},
		"nfServices": [
		  {
			"serviceInstanceId": "nudm-auth-01",
			"nfServiceStatus": "REGISTERED",
			"serviceName": "nudm-auth",
			"versions": [{
			  "apiVersionInUri":"v1",
			  "apiFullVersion": "1.R15.1.1 " ,
			  "expiry":"2020-07-06T02:54:32Z"}],
			"scheme": "http",
			"fqdn": "seliius03696.seli.gic.ericsson.se",
			"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
			"ipEndPoints":[
			  {
				"ipv4Address": "172.16.208.1",
				"transport": "TCP",
				"port": 30088
			  }
			],
			"apiPrefix": "mytest/nausf-auth/v1",
			"defaultNotificationSubscriptions": [
			  {
				"notificationType": "N1_MESSAGES",
				"callbackUri": "/nnrf-nfm/v1/nf-instances/ausf-5g-01",
				"n1MessageClass": "5GMM",
				"n2InformationClass": "SM"
			  }
			],
			"allowedPlmns": [
			  {
				"mcc": "460",
				"mnc": "00"
			  }
			],
			"allowedNfTypes": [
			  "NEF", "PCF", "SMSF", "NSSF",
			  "UDR", "LMF", "5G_EIR", "SEPP", "UPF", "N3IWF", "AF", "UDSF"
			],
			"allowedNssais": [
			  {
				"sst": 0,
				"sd": "abAB01"
			  }
			],
			"supportedFeatures":"A0A0",
			"capacity": 100,
			"load" : 100
		  }
		]
	  }`)
)

func init() {
	log.SetLevel(log.FatalLevel)
}

func TestGenerateExpiredTime(t *testing.T) {
	expiredTime := GenerateExpiredTime(2)
	time.Sleep(1 * time.Second)
	if uint64(time.Now().Unix()*1000) > expiredTime {
		t.Fatalf("Should not be expired, but not!")
	}
	time.Sleep(2 * time.Second)
	if uint64(time.Now().Unix()*1000) <= expiredTime {
		t.Fatalf("Should be expired, but not!")
	}
}

func TestGenerateLastUpdateTime(t *testing.T) {
	currentTime := uint64(time.Now().Unix() * 1000)
	time.Sleep(1 * time.Second)
	lastUpdateTime := GenerateLastUpdateTime()
	if currentTime >= lastUpdateTime {
		t.Fatalf("lastUpdateTime should be larger than currentTime")
	}
}

func TestGenerateProfileUpdateTime(t *testing.T) {
	currentTime := uint64(time.Now().UnixNano() / 1000000)
	time.Sleep(1 * time.Second)
	profileUpdateTime := GenerateProfileUpdateTime()
	if currentTime >= profileUpdateTime {
		t.Fatalf("profileUpdateTime should be larger than currentTime")
	}
}

func TestCheckSupportedFields(t *testing.T) {
	var targetParameterArray []string
	targetParameterArray = append(targetParameterArray, "1")
	queryForm := make(map[string][]string)
	queryForm["nf-type"] = targetParameterArray
	queryForm["limit"] = targetParameterArray
	invalidParameters := CheckSupportedFields(queryForm, "nf-type", "limit")
	if invalidParameters != nil {
		t.Fatal("Should return nil, but not !")
	}

	queryForm["nfType"] = targetParameterArray
	invalidParameters = CheckSupportedFields(queryForm, "nf-type", "limit")
	if invalidParameters == nil {
		t.Fatal("Should not return true, but did !")
	}
}

func TestGetServiceName(t *testing.T) {

	nfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123",
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123",
							}
						]
					}`)

	serviceNames, problemDetails := GetServiceName(nfProfile)
	if problemDetails != nil {
		fmt.Printf(problemDetails.ToString())
		t.Fatalf("valid nfServices, but return no service names")
	}
	if serviceNames[0] != "udm-svc1" || serviceNames[1] != "udm-svc2" {
		t.Fatalf("Get invalid nfServiceNames")
	}

	nfProfileInvalid := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceNAME": "udm-svc1",
							   "version": [],
							   "schema": "abc123",
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version":[],
							   "schema": "abc123",
							}
						]
					}`)

	serviceNames, problemDetails = GetServiceName(nfProfileInvalid)
	if problemDetails == nil {
		t.Fatalf("invalid nfServices, but return no error")
	}
}

func TestValidateOtherRules(t *testing.T) {
	profile := `{
  "nfInstanceId": "udm1",
  "nfType": "UDM",
  "nfStatus": "REGISTERED",
  "sNssais": [
    {
      "sst": 1,
      "sd": "sd11"
    },
        {
      "sst": 21,
      "sd": "sd21"
    }
  ],
  "fqdn": "string",
  "interPlmnFqdn": "string",
  "ipv4Addresses": [
    "10.10.10.10"
  ],
  "ipv6Addresses": [
    "CDCD:910A:2222:5498:8475:1111:3900:2020"
  ],
  "nsiList": [],
  "capacity": 0,
  "udrInfo": {
    "supiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
    "gpsiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
	"externalGroupIdentifiersRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ]
  },
  "udmInfo": {
    "groupId": "001",
    "gpsiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
	"externalGroupIdentifiersRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
    "routingIndicator": "1234"
  },
  "nfServices": [
    {
      "serviceInstanceId": "srv1",
      "serviceName": "srv1",
      "versions": [{"apiVersionInUri":"","apiFullVersion":""}],
      "scheme": "http",
      "ipEndPoints": [
        {
          "ipv4Address": "10.10.10.10",
          "port": 0
        }
      ]
    },
    {
      "serviceInstanceId": "srv2",
      "serviceName": "srv2",
      "versions": [{"apiVersionInUri":"","apiFullVersion":""}],
      "scheme": "http",
      "ipEndPoints": [
        {
          "ipv6Address": "2000:0:0:0:0:0:0:1",
          "port": 0
        }
      ]
    },
    {
      "serviceInstanceId": "srv3",
      "serviceName": "srv3",
      "versions": [],
      "scheme": "http",
      "ipEndPoints": [
        {
          "port": 0
        }
      ]
    },
	{
      "serviceInstanceId": "srv4",
      "serviceName": "srv4",
      "versions": [],
      "scheme": "http"
    }
  ]
}`

	problemDetails := ValidateOtherRules([]byte(profile))
	if problemDetails != nil {
		t.Fatalf("%s", problemDetails.ToString())
		t.Fatalf("This is a valid nf profile, but validate failed !")
	}

	profile = `{
  "nfInstanceId": "udm1",
  "nfType": "UDM",
  "nfStatus": "REGISTERED",
  "sNssais": [
    {
      "sst": 1,
      "sd": "sd11"
    },
        {
      "sst": 21,
      "sd": "sd21"
    }
  ],
  "fqdn": "string",
  "interPlmnFqdn": "string",
  "ipv4Addresses": [
    "10.10.10.10"
  ],
  "ipv6Addresses": [
    "CDCD:910A:2222:5498:8475:1111:3900:2020"
  ],
  "nsiList": [],
  "capacity": 0,
  "udrInfo": {
    "supiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
    "gpsiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
	"externalGroupIdentifiersRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ]
  },
  "udmInfo": {
    "groupId": "001",
    "gpsiRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
	"externalGroupIdentifiersRanges": [
      {
        "start": "0",
        "end": "9"
      },
	  {
        "pattern": "string"
      },
	  {
        "pattern": "string"
      }
    ],
    "routingIndicator": "1234"
  },
  "nfServices": [
    {
      "serviceInstanceId": "srv1",
      "serviceName": "srv1",
      "versions": [{"apiVersionInUri":"","apiFullVersion":""}],
      "scheme": "http",
      "ipEndPoints": [
        {
          "ipv4Address": "10.10.10.10",
          "port": 0
        }
      ]
    },
    {
      "serviceInstanceId": "srv2",
      "serviceName": "srv2",
      "versions": [{"apiVersionInUri":"","apiFullVersion":""}],
      "scheme": "http",
      "ipEndPoints": [
        {
		  "ipv4Address": "10.10.10.10",
          "ipv6Address": "2000:0:0:0:0:0:0:1",
          "port": 0
        }
      ]
    },
    {
      "serviceInstanceId": "srv3",
      "serviceName": "srv3",
      "versions": [],
      "scheme": "http",
      "ipEndPoints": [
        {
          "port": 0
        }
      ]
    },
	{
      "serviceInstanceId": "srv4",
      "serviceName": "srv4",
      "versions": [],
      "scheme": "http"
    }
  ]
}`

	problemDetails = ValidateOtherRules([]byte(profile))
	if problemDetails == nil {
		t.Fatalf("This is a invalid nf profile, but validate OK !")
	}

}

func TestGetChangedServiceName(t *testing.T) {
	oldNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	changedServiceNames, _ := GetChangedServiceName(oldNfProfile, newNfProfile)

	if changedServiceNames != nil {
		t.Fatalf("nfServices in Old profile is the same as in new profile, but difference returned: %s,%s", changedServiceNames[0], changedServiceNames[1])
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc1234"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc1234"
							}
						]
					}`)

	changedServiceNames, _ = GetChangedServiceName(oldNfProfile, newNfProfile)

	if changedServiceNames == nil {
		t.Fatalf("nfServices in Old profile is different from in new profile, but no difference returned")
	}

	if len(changedServiceNames) != 2 {
		t.Fatalf("Count of changed services not match")
	}
	for _, item := range changedServiceNames {
		if item != "udm-svc1" && item != "udm-svc2" {
			t.Fatalf("Servicename of changed services not match")
		}
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc3",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	changedServiceNames, _ = GetChangedServiceName(oldNfProfile, newNfProfile)

	if changedServiceNames == nil {
		t.Fatalf("nfServices in Old profile is different from in new profile, but no difference returned")
	}

	if len(changedServiceNames) != 2 {
		t.Fatalf("Count of changed services not match")
	}
	for _, item := range changedServiceNames {
		if item != "udm-svc2" && item != "udm-svc3" {
			t.Fatalf("Servicename of changed services not match")
		}
	}

}

func TestGetUnionServiceName(t *testing.T) {
	oldNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm-svc1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc2",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm-svc2",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc3",
							   "serviceName": "udm-svc3",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc3",
							   "serviceName": "udm-svc3",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	unionServiceNames, _ := GetUnionServiceName(oldNfProfile, newNfProfile)

	length := len(unionServiceNames)

	if length != 3 {
		t.Fatalf("There should be 3 services, but not.")
	}

	mapServiceExist := map[string]bool{
		"udm-svc1": false,
		"udm-svc2": false,
		"udm-svc3": false,
	}

	for _, item := range unionServiceNames {
		mapServiceExist[item] = true
	}

	ok := true
	for _, v := range mapServiceExist {
		if !v {
			ok = false
			break
		}
	}

	if !ok {
		t.Fatalf("GetUnionServiceName not return the right service names")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm-svc1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc2",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc2",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0
					}`)

	unionServiceNames, _ = GetUnionServiceName(oldNfProfile, newNfProfile)

	length = len(unionServiceNames)

	if length != 2 {
		t.Fatalf("There should be 2 services, but not.")
	}

	mapServiceExist = map[string]bool{
		"udm-svc1": false,
		"udm-svc2": false,
	}

	for _, item := range unionServiceNames {
		mapServiceExist[item] = true
	}

	ok = true
	for _, v := range mapServiceExist {
		if !v {
			ok = false
			break
		}
	}

	if !ok {
		t.Fatalf("GetUnionServiceName not return the right service names")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"UDM",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm-svc1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc2",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm-svc1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	unionServiceNames, _ = GetUnionServiceName(oldNfProfile, newNfProfile)

	length = len(unionServiceNames)

	if length != 2 {
		t.Fatalf("There should be 2 services, but not.")
	}

	mapServiceExist = map[string]bool{
		"udm-svc1": false,
		"udm-svc2": false,
	}

	for _, item := range unionServiceNames {
		mapServiceExist[item] = true
	}

	ok = true
	for _, v := range mapServiceExist {
		if !v {
			ok = false
			break
		}
	}

	if !ok {
		t.Fatalf("GetUnionServiceName not return the right service names")
	}

}

func TestIsProfileChanged(t *testing.T) {
	oldNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)
	if IsProfileChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("IsProfileChanged should return false, but return true.")
	}

	oldNfProfile = []byte(`{
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						],
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)
	if IsProfileChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("IsProfileChanged should return false, but return true.")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "1"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)
	if !IsProfileChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("IsProfileChanged should return true, but return false.")
	}
}

func TestIsServiceChanged(t *testing.T) {
	oldNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile := []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	if isServiceChanged("udm-svc1", oldNfProfile, newNfProfile) {
		t.Fatalf("Service not change, but return true")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc1234"
							}
						]
					}`)

	if !isServiceChanged("udm-svc2", oldNfProfile, newNfProfile) {
		t.Fatalf("Service change, but return false")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc1234"
							}
						]
					}`)

	if isServiceChanged("udm-svc3", oldNfProfile, newNfProfile) {
		t.Fatalf("Service change, but return false")
	}

	oldNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
							{
							   "serviceInstanceId": "udm2",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	newNfProfile = []byte(`{
	                    "nfInstanceID":"udm",
						"nfType":"udm",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm1",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm2",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc1234"
							}
						]
					}`)

	if !isServiceChanged("udm-svc1", oldNfProfile, newNfProfile) {
		t.Fatalf("Service change, but return false")
	}
}

func TestGetNFType(t *testing.T) {
	nfProfile := []byte(`{
	                         "nfInstanceId":"udm",
						    "nfType":"UDM",
							"nfStatus": "REGISTERED",
						    "plmn": {"mcc":"460", "mnc":"01"},
						    "sNssais": [{"sst": 1,"sd": "s1"}, {"sst": 2,"sd": "s2"}],
							"nsiList": ["nsi1", "nsi2"],
						    "fqdn": "udm.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
						    "nfServices": [
							    {
							        "serviceName": "udm1"
						        },
								{
							        "serviceName": "udm2"
						        },
								{
							        "serviceName": "udm3"
						        }
					  		]
					    }`)
	nfType, problemDetails := GetNFType(nfProfile)
	if problemDetails != nil {
		t.Fatalf("getNFType should not return error, but did.")
	}

	if nfType != "UDM" {
		t.Fatalf("getNFType didn't return right nfType.")
	}

	nfProfile = []byte(`{
	                         "nfInstanceId":"udm",
							"nfStatus": "REGISTERED",
						    "plmn": {"mcc":"460", "mnc":"01"},
						    "sNssais": [{"sst": 1,"sd": "s1"}, {"sst": 2,"sd": "s2"}],
							"nsiList": ["nsi1", "nsi2"],
						    "fqdn": "udm.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
						    "nfServices": [
							    {
							        "serviceName": "udm1"
						        },
								{
							        "serviceName": "udm2"
						        },
								{
							        "serviceName": "udm3"
						        }
					  		]
					    }`)
	nfType, problemDetails = GetNFType(nfProfile)
	if problemDetails == nil {
		t.Fatalf("getNFType should return error, but not.")
	}
}

func TestGetNFStatus(t *testing.T) {
	nfProfile := []byte(`{
	                         "nfInstanceId":"udm",
						    "nfType":"UDM",
							"nfStatus": "REGISTERED",
						    "plmn": {"mcc":"460", "mnc":"01"},
						    "sNssais": [{"sst": 1,"sd": "s1"}, {"sst": 2,"sd": "s2"}],
							"nsiList": ["nsi1", "nsi2"],
						    "fqdn": "udm.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
						    "nfServices": [
							    {
							        "serviceName": "udm1"
						        },
								{
							        "serviceName": "udm2"
						        },
								{
							        "serviceName": "udm3"
						        }
					  		]
					    }`)
	nfStatus, problemDetails := GetNFStatus(nfProfile)
	if problemDetails != nil {
		t.Fatalf("getNFStatus should not return error, but did.")
	}

	if nfStatus != "REGISTERED" {
		t.Fatalf("getNFStatus didn't return right nfStatus.")
	}

	nfProfile = []byte(`{
	                         "nfInstanceId":"udm",
							"nfType":"UDM",
						    "plmn": {"mcc":"460", "mnc":"01"},
						    "sNssais": [{"sst": 1,"sd": "s1"}, {"sst": 2,"sd": "s2"}],
							"nsiList": ["nsi1", "nsi2"],
						    "fqdn": "udm.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
						    "nfServices": [
							    {
							        "serviceName": "udm1"
						        },
								{
							        "serviceName": "udm2"
						        },
								{
							        "serviceName": "udm3"
						        }
					  		]
					    }`)
	nfStatus, problemDetails = GetNFStatus(nfProfile)
	if problemDetails == nil {
		t.Fatalf("getNFStatus should return error, but not.")
	}
}

func TestGetNfInfoChangedCode(t *testing.T) {}

func TestGetNfInfo(t *testing.T) {}

func TestLastSubString(t *testing.T) {
	if LastSubString("http://10.10.10.10:8080", "/") != "10.10.10.10:8080" {
		t.Fatalf("LastSubString didn't return right value.")
	}
	if LastSubString("https://10.10.10.10:8080", "/") != "10.10.10.10:8080" {
		t.Fatalf("LastSubString didn't return right value.")
	}

	if LastSubString("http://www.example.com:8080", "/") != "www.example.com:8080" {
		t.Fatalf("LastSubString didn't return right value.")
	}

	if LastSubString("https://www.example.com:8080", "/") != "www.example.com:8080" {
		t.Fatalf("LastSubString didn't return right value.")
	}
}

func TestGetNFInstanceID(t *testing.T) {}

func TestContructNfInstanceURI(t *testing.T) {}

func TestConstructNRFAddressWithPlmnId(t *testing.T) {}

func TestIsProfileCommonPartChanged(t *testing.T) {
	oldNfProfile := []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)
	newNfProfile := []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	if IsProfileCommonPartChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("profile common part didn't change, but return true.")
	}

	oldNfProfile = []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc1234"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc1234"
							}
						]
					}`)
	newNfProfile = []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	if IsProfileCommonPartChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("profile common part didn't change, but return true.")
	}

	oldNfProfile = []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.1"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)
	newNfProfile = []byte(`{
	                     "nfInstanceId":"udm",
						"nfType":"UDM",
						"nfStatus": "REGISTERED",
						"plmn": "46000", 
						"sNssai": {"sst": "0","sd": "0"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.2"],
						"capacity": 0,
						"nfServices": [
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc1",
							   "version": [],
							   "schema": "abc123"
							},
						    {
							   "serviceInstanceId": "udm",
							   "serviceName": "udm-svc2",
							   "version": [],
							   "schema": "abc123"
							}
						]
					}`)

	if !IsProfileCommonPartChanged(oldNfProfile, newNfProfile) {
		t.Fatalf("profile common part changed, but return false.")
	}

}

func TestConstructPlmnID(t *testing.T) {
	if ConstructPlmnID("123", "45") != "12345" {
		t.Fatalf("constructPlmnID didn't construct right plmnID")
	}

	if ConstructPlmnID("123", "456") != "123456" {
		t.Fatalf("constructPlmnID didn't construct right plmnID")
	}
}

func TestRequestFromNRFProv(t *testing.T) {
	flag := RequestFromNRFProv("eric-nrf-management:3004", 3004)
	if !flag {
		t.Fatalf("requestFromNRFProv failed")
	}

	flag = RequestFromNRFProv("10.109.1.2:80", 109)
	if flag {
		t.Fatalf("requestFromNRFProv failed")
	}

	flag = RequestFromNRFProv("[3ffe:2a00:109:7031::1]:80", 109)
	if flag {
		t.Fatalf("requestFromNRFProv failed")
	}

	flag = RequestFromNRFProv("[3ffe:2a00:109:7031::1]:3004", 3004)
	if !flag {
		t.Fatalf("requestFromNRFProv failed")
	}
}

func TestRequestFromNRFMgmt(t *testing.T) {
	flag := RequestFromNRFMgmt("eric-nrf-management:3004", 3004)
	if !flag {
		t.Fatalf("requestFromNRFMgmt failed")
	}

	flag = RequestFromNRFMgmt("10.109.1.2:80", 109)
	if flag {
		t.Fatalf("requestFromNRFMgmt failed")
	}

	flag = RequestFromNRFMgmt("[3ffe:2a00:109:7031::1]:80", 109)
	if flag {
		t.Fatalf("requestFromNRFMgmt failed")
	}

	flag = RequestFromNRFMgmt("[3ffe:2a00:109:7031::1]:3004", 3004)
	if !flag {
		t.Fatalf("requestFromNRFMgmt failed")
	}
}

var overrideBody = []byte(`[
	{
			"action": "replace",
			"path": "/nfServices/0/priority",
			"value": "10"
	},
	{
			"action": "add",
			"path": "/nfServices/1/capacity",
			"value": "50"
	}
]`)

func TestRebuildNfServiceOverrideInfo(t *testing.T) {
	oldMapper := make(map[string]int, 0)
	newMapper := make(map[string]int, 0)

	oldMapper["nudm-auth-01"] = 0
	oldMapper["nudm-auth-02"] = 1

	newMapper["nudm-auth-01"] = 0
	newMapper["nudm-auth-02"] = 1
	newMapper["nudm-auth-03"] = 2

	var overrideData []nrfschema.OverrideInfo
	err := json.Unmarshal(overrideBody, &overrideData)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}

	newOverrideData := RebuildNfServiceOverrideInfo(oldMapper, newMapper, overrideData)
	if len(overrideData) != len(newOverrideData) {
		t.Fatalf("Expect new overrideInfo is the same, but not")
	}

	for i, overrideItem := range overrideData {
		newOverrideItem := newOverrideData[i]
		if overrideItem.Path != newOverrideItem.Path || overrideItem.Action != newOverrideItem.Action || overrideItem.Value != newOverrideItem.Value {
			t.Fatalf("Expect each overrideItem is the same, but not")
		}
	}

	///////2///////
	oldMapper2 := make(map[string]int, 0)
	newMapper2 := make(map[string]int, 0)

	oldMapper2["nudm-auth-01"] = 0
	oldMapper2["nudm-auth-02"] = 1

	newMapper2["nudm-auth-01"] = 0

	var overrideData2 []nrfschema.OverrideInfo
	err = json.Unmarshal(overrideBody, &overrideData2)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}

	newOverrideData2 := RebuildNfServiceOverrideInfo(oldMapper2, newMapper2, overrideData2)
	if len(newOverrideData2) != 1 {
		t.Fatalf("Expect new overrideInfo left one item, but not")
	}

	overrideItem := newOverrideData2[0]
	if overrideItem.Path != "/nfServices/0/priority" || overrideItem.Action != "replace" || overrideItem.Value != "10" {
		t.Fatalf("The left overrideItem is not correct")
	}
}

func TestSplitNrfInfo(t *testing.T) {
	nrfProfile := []byte(`{
	"capacity": 100,
	"fqdn": "seliius03696.seli.gic.ericsson.se",
	"nfInstanceId": "12345678-abcd-ef12-1000-000000000010",
	"nfServices": [{
			"capacity": 100,
			"fqdn": "seliius03690.seli.gic.ericsson.se",
			"ipEndPoints": [{
				"ipv4Address": "172.16.208.1",
				"port": 30088
			}],
			"nfServiceStatus": "REGISTERED",
			"priority": 100,
			"scheme": "https",
			"serviceInstanceId": "nudm-uecm-01",
			"serviceName": "nudm-uecm",
			"versions": [{
				"apiFullVersion": "1.R15.1.1",
				"apiVersionInUri": "v1",
				"expiry": "2020-07-06T02: 54: 32Z"
			}]
		},
		{
			"fqdn": "seliius03690.seli.gic.ericsson.se",
			"ipEndPoints": [{
				"ipv4Address": "172.16.208.2",
				"port": 30088
			}],
			"nfServiceStatus": "REGISTERED",
			"priority": 100,
			"scheme": "https",
			"serviceInstanceId": "nudm-uecm-02",
			"serviceName": "nudm-uecm",
			"versions": [{
				"apiFullVersion": "1.R15.1.1",
				"apiVersionInUri": "v1",
				"expiry": "2020-07-06T02: 54: 32Z"
			}]
		},
		{
			"fqdn": "seliius03690.seli.gic.ericsson.se",
			"ipEndPoints": [{
				"ipv4Address": "172.16.208.3",
				"port": 30088
			}],
			"nfServiceStatus": "REGISTERED",
			"priority": 100,
			"scheme": "https",
			"serviceInstanceId": "nudm-sdm-01",
			"serviceName": "nudm-sdm",
			"versions": [{
				"apiFullVersion": "1.R15.1.1",
				"apiVersionInUri": "v1",
				"expiry": "2020-07-06T02: 54: 32Z"
			}]
		}
	],
	"nfStatus": "REGISTERED",
	"nfType": "NRF",
	"plmn": {
		"mcc": "460",
		"mnc": "000"
	},
	"priority": 100,
	"sNssais": [{
		"sd": "222222",
		"sst": 2
	}],
	"nrfInfo": {
		"servedUdrInfo": {
			"instanceId-001": {
				"groupId": "gg",
				"supiRanges": [{
					"start": "111",
					"end": "222"
				}],
				"gpsiRanges": [{
					"start": "111",
					"end": "222"
				}],
				"externalGroupIdentityfiersRanges": [{
					"start": "111",
					"end": "222"
				}],
				"supportedDataSets": [
					"hh"
				]
			}
		},
		"servedUdmInfo": {
			"instanceId-002": {
				"groupId": "",
				"supiRanges": [{
					"start": "111",
					"end": "222"
				}],
				"gpsiRanges": [{
					"start": "111",
					"end": "222"
				}],
				"externalGroupIdentityfiersRanges": [{
					"start": "111",
					"end": "222"
				}],
				"routingIndicators": [
					"hh"
				]
			}
		},
		"servedAusfInfo": {
			"instanceId-003": {
				"groupId": "gg",
				"supiRanges": [{
					"start": "111",
					"end": "222"
				}],
				"routingIndicators": [
					"hh"
				]
			}
		},
		"servedAmfInfo": {
			"instanceId-004": {
				"amfSetId": "hh",
				"amfRegionId": "hh",
				"guamiList": [{
					"plmnId": {
						"mnc": "000",
						"mcc": "460"
					},
					"amfId": "hh"
				}],
				"taiList": [{
					"plmnId": {
						"mnc": "000",
						"mcc": "460"
					},
					"tac": "12312"
				}],
				"taiRangeList": [{
					"plmnId": {
						"mnc": "000",
						"mcc": "460"
					},
					"tacRangeList": [{
						"start": "111",
						"end": "222"
					}]
				}],
				"backupInfoAmfFailure": [{
					"plmnId": {
						"mnc": "460",
						"mcc": "00"
					},
					"amfId": "12"
				}],
				"backupInfoAmfRemoval": [{
					"plmnId": {
						"mnc": "460",
						"mcc": "000"
					},
					"amfId": "hh"
				}],
				"n2InterfaceAmfInfo": {
					"ipv4EndpointAddress": [],
					"ipv6EndpointAddress": [],
					"amfName": "12213"
				}
			}
		},
		"servedSmfInfo": {
			"instanceId-005": {
				"sNssaiSmfInfoList": [{
					"sNssai": {
						"sst": 111,
						"sd": "222"
					},
					"dnnSmfInfoList": [{
						"dnn": "333"
					}]
				}],
				"taiList": [{
					"plmnId": {
						"mnc": "000",
						"mcc": "460"
					},
					"tac": "123"
				}],
				"taiRangeList": [{
					"plmnId": {
						"mnc": "000",
						"mcc": "460"
					},
					"tacRangeList": [{
						"start": "1221",
						"end": "111"
					}]
				}],
				"pgwFqdn": "123",
				"accessType": [
					"3GPP_ACCESS"
				]
			}
		},
		"servedUpfInfo": {
			"instanceId-006": {
				"sNssaiUpfInfoList": [{
					"sNssai": {
						"sst": 111,
						"sd": "11"
					},
					"dnnUpfInfoList": [{
						"dnn": "12",
						"dnaiList": []
					}]
				}],
				"smfServingArea": [
					"fd"
				],
				"interfaceUpfInfoList": [{
					"interfaceType": [
						"fd"
					],
					"ipv4EndpointAddress": [
						""
					],
					"ipv6EndpointAddress": [
						""
					],
					"endpointFqdn": "",
					"networkInstance": ""
				}],
				"iwkEpsInd": true
			}
		},
		"servedPcfInfo": {
			"instanceId-007": {
				"dnnList": [
					"12"
				],
				"supiRanges": [{
					"start": "111",
					"end": "222"
				}]
			}
		},
		"servedBsfInfo": {
			"instanceId-008": {
				"dnnList": [
					"11"
				],
				"ipDomainList": [
					"22"
				],
				"ipv4AddressRanges": [{
					"start": "11",
					"end": "123"
				}]
			}
		},
		"servedChfInfo": {
			"instanceId-009": {
				"supiRangeList": [{
					"start": "111",
					"end": "222"
				}],
				"gpsiRangeList": [{
					"start": "111",
					"end": "222"
				}],
				"plmnRangeList": [{
					"start": "111",
					"end": "222"
				}]
			}
		}
	}
}`)
	nrfprofileScheme := &nrfschema.TNFProfile{}
	err := json.Unmarshal(nrfProfile, nrfprofileScheme)
	if err != nil {
		t.Fatalf("unmarshal fail, error=%v", err)
	}
	fmt.Println(SplitNrfInfo("12345678-abcd-ef12-1000-000000000010", nrfprofileScheme.NrfInfo))
}

func TestInsertNfInfoToNrfInfo(t *testing.T) {
	udrProfile := `{"nfInstanceId": "instId-udr",
		"nrfInstanceId": "instId-udr",
		"nfType": "UDR",
		"udrInfo": {
			"groupId": "gg",
				"supiRanges": [{
			"start": "111",
			"end": "222"
			}],
			"gpsiRanges": [{
			"start": "111",
			"end": "222"
			}],
			"externalGroupIdentityfiersRanges": [{
			"start": "111",
			"end": "222"
			}],
			"supportedDataSets": [
			"hh"
			]
		}
	}`
	udmProfile := `{"nfInstanceId": "instId-udm",
		"nrfInstanceId": "instId-udm",
		"nfType": "UDM",
		"udmInfo": {
			"groupId": "",
				"supiRanges": [{
			"start": "111",
			"end": "222"
			}],
			"gpsiRanges": [{
			"start": "111",
			"end": "222"
			}],
			"externalGroupIdentityfiersRanges": [{
			"start": "111",
			"end": "222"
			}],
			"routingIndicators": [
			"hh"
			]
		}
	}`

	ausfProfile := `{"nfInstanceId": "instId-ausf",
		"nrfInstanceId": "instId-ausf",
		"nfType": "AUSF",
		"ausfInfo": {
			"groupId": "gg",
				"supiRanges": [{
			"start": "111",
			"end": "222"
			}],
			"routingIndicators": [
			"hh"
			]
		}
	}`

	amfProfile := `{"nfInstanceId": "instId-amf",
		"nrfInstanceId": "instId-amf",
		"nfType": "AMF",
		"amfInfo": {
			"amfSetId": "hh",
				"amfRegionId": "hh",
				"guamiList": [{
			"plmnId": {
			"mnc": "000",
			"mcc": "460"
			},
			"amfId": "hh"
			}],
			"taiList": [{
			"plmnId": {
			"mnc": "000",
			"mcc": "460"
			},
			"tac": "12312"
			}],
			"taiRangeList": [{
			"plmnId": {
			"mnc": "000",
			"mcc": "460"
			},
			"tacRangeList": [{
			"start": "111",
			"end": "222"
			}]
			}],
			"backupInfoAmfFailure": [{
			"plmnId": {
			"mnc": "460",
			"mcc": "00"
			},
			"amfId": "12"
			}],
			"backupInfoAmfRemoval": [{
			"plmnId": {
			"mnc": "460",
			"mcc": "000"
			},
			"amfId": "hh"
			}],
			"n2InterfaceAmfInfo": {
				"ipv4EndpointAddress": [],
				"ipv6EndpointAddress": [],
				"amfName": "12213"
			}
		}
	}`

	smfProfile := `{"nfInstanceId": "instId-smf",
		"nrfInstanceId": "instId-smf",
		"nfType": "SMF",
		"smfInfo": {
			"sNssaiSmfInfoList": [{
			"sNssai": {
			"sst": 111,
			"sd": "222"
			},
			"dnnSmfInfoList": [{
			"dnn": "333"
			}]
			}],
			"taiList": [{
			"plmnId": {
			"mnc": "000",
			"mcc": "460"
			},
			"tac": "123"
			}],
			"taiRangeList": [{
			"plmnId": {
			"mnc": "000",
			"mcc": "460"
			},
			"tacRangeList": [{
			"start": "1221",
			"end": "111"
			}]
			}],
			"pgwFqdn": "123",
				"accessType": [
			"3GPP_ACCESS"
			]
		}
	}`

	upfProfile := `{"nfInstanceId": "instId-upf",
		"nrfInstanceId": "instId-upf",
		"nfType": "UPF",
		"upfInfo": {
			"sNssaiUpfInfoList": [{
			"sNssai": {
			"sst": 111,
			"sd": "11"
			},
			"dnnUpfInfoList": [{
			"dnn": "12",
			"dnaiList": []
			}]
			}],
			"smfServingArea": [
			"fd"
			],
			"interfaceUpfInfoList": [{
			"interfaceType": [
			"fd"
			],
			"ipv4EndpointAddress": [
			""
			],
			"ipv6EndpointAddress": [
			""
			],
			"endpointFqdn": "",
			"networkInstance": ""
			}],
			"iwkEpsInd": true
		}
	}`

	pcfProfile := `{"nfInstanceId": "instId-pcf",
		"nrfInstanceId": "instId-pcf",
		"nfType": "PCF",
		"pcfInfo": {
			"dnnList": [
			"12"
			],
			"supiRanges": [{
			"start": "111",
			"end": "222"
			}]
		}
	}`
	bsfProfile := `{"nfInstanceId": "instId-bsf",
		"nrfInstanceId": "instId-bsf",
		"nfType": "BSF",
		"bsfInfo": {
			"dnnList": [
			"11"
			],
			"ipDomainList": [
			"22"
			],
			"ipv4AddressRanges": [{
			"start": "11",
			"end": "123"
			}]
		}
	}`
	chfProfile := `{"nfInstanceId": "instId-chf",
		"nrfInstanceId": "instId-chf",
		"nfType": "CHF",
		"chfInfo": {
			"supiRangeList": [{
			"start": "111",
			"end": "222"
			}],
			"gpsiRangeList": [{
			"start": "111",
			"end": "222"
			}],
			"plmnRangeList": [{
			"start": "111",
			"end": "222"
			}]
		}
	}`
	nrfInfo := &nrfschema.TNrfInfo{}
	insertNfInfoToNrfInfo(nrfInfo, udrProfile)
	if nrfInfo.ServedUdrInfo["instId-udr"] == nil {
		t.Fatal("udr should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, udmProfile)
	if nrfInfo.ServedUdmInfo["instId-udm"] == nil {
		t.Fatal("udm should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, upfProfile)
	if nrfInfo.ServedUpfInfo["instId-upf"] == nil {
		t.Fatal("upf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, ausfProfile)
	if nrfInfo.ServedAusfInfo["instId-ausf"] == nil {
		t.Fatal("ausf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, amfProfile)
	if nrfInfo.ServedAmfInfo["instId-amf"] == nil {
		t.Fatal("amf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, bsfProfile)
	if nrfInfo.ServedBsfInfo["instId-bsf"] == nil {
		t.Fatal("bsf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, chfProfile)
	if nrfInfo.ServedChfInfo["instId-chf"] == nil {
		t.Fatal("chf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, pcfProfile)
	if nrfInfo.ServedPcfInfo["instId-pcf"] == nil {
		t.Fatal("pcf should not be null")
	}
	insertNfInfoToNrfInfo(nrfInfo, smfProfile)
	if nrfInfo.ServedSmfInfo["instId-smf"] == nil {
		t.Fatal("smf should not be null")
	}
}

func TestSupiRangeInfoExist(t *testing.T) {
	nfProfile := &nrfschema.TNFProfile{}
	err := json.Unmarshal(newNfProfileData, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal NF profile error, %s", err.Error())
	}
	exist := SupiRangeInfoExist(nfProfile)
	if !exist {
		t.Fatal("Expect udmInfo supiRanges exist, but not")
	}
}

func TestGpsiRangeInfoExist(t *testing.T) {
	nfProfile := &nrfschema.TNFProfile{}
	err := json.Unmarshal(newNfProfileData, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal NF profile error, %s", err.Error())
	}
	exist := GpsiRangeInfoExist(nfProfile)
	if !exist {
		t.Fatal("Expect udmInfo gpsiRanges is exist, but not")
	}
}
