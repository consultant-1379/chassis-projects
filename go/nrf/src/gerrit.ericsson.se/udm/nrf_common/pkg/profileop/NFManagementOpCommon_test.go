package profileop

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
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

func TestGetAmfSetId(t *testing.T) {
	amfInfo := []byte(`{
                             "amfSetId": "amfset1",
                             "amfRegionId": "amfRegion1",
                             "guamiList": [
                                 {
                                     "plmnId": {"mcc": "460", "mnc": "000"},
                                     "amfId": "123456"
                                 },
                                 {
                                     "plmnId": {"mcc": "461", "mnc": "01"},
                                     "amfId": "654321"
                                 }
                             ],
                             "taiList": [
                                 {
                                     "plmnId": {"mcc": "462", "mnc": "002"},
                                     "tac": "234567"
                                 },
                                 {
                                     "plmnId": {"mcc": "463", "mnc": "03"},
                                     "tac": "765432"
                                 }
                             ]
                         }`)
	amfSetId, problemDetails := getAmfSetId(amfInfo)

	if problemDetails != nil {
		t.Fatalf("should not return error, but did.")
	}

	if amfSetId != "amfset1" {
		t.Fatalf("getAmfProperties didn't return right amfSetId")
	}

	amfInfo = []byte(`{
                             "amfRegionId": "amfRegion1",
                             "guamiList": [
                                 {
                                     "plmnId": {"mcc": "460", "mnc": "000"},
                                     "amfId": "123456"
                                 },
                                 {
                                     "plmnId": {"mcc": "461", "mnc": "01"},
                                     "amfId": "654321"
                                 }
                             ],
                             "taiList": [
                                 {
                                     "plmnId": {"mcc": "462", "mnc": "002"},
                                     "tac": "234567"
                                 },
                                 {
                                     "plmnId": {"mcc": "463", "mnc": "03"},
                                     "tac": "765432"
                                 }
                             ]
                         }`)
	amfSetId, problemDetails = getAmfSetId(amfInfo)

	if problemDetails == nil {
		t.Fatalf("should return error, but not.")
	}

}

func TestGetAmfRegionId(t *testing.T) {
	amfInfo := []byte(`{
                             "amfSetId": "amfset1",
                             "amfRegionId": "amfRegion1",
                             "guamiList": [
                                 {
                                     "plmnId": {"mcc": "460", "mnc": "000"},
                                     "amfId": "123456"
                                 },
                                 {
                                     "plmnId": {"mcc": "461", "mnc": "01"},
                                     "amfId": "654321"
                                 }
                             ],
                             "taiList": [
                                 {
                                     "plmnId": {"mcc": "462", "mnc": "002"},
                                     "tac": "234567"
                                 },
                                 {
                                     "plmnId": {"mcc": "463", "mnc": "03"},
                                     "tac": "765432"
                                 }
                             ]
                         }`)
	amfRegionId, problemDetails := getAmfRegionId(amfInfo)

	if problemDetails != nil {
		t.Fatalf("should not return error, but did.")
	}

	if amfRegionId != "amfRegion1" {
		t.Fatalf("getAmfProperties didn't return right amfRegionId")
	}

	amfInfo = []byte(`{
                             "amfSetId": "amfset1",
                             "guamiList": [
                                 {
                                     "plmnId": {"mcc": "460", "mnc": "000"},
                                     "amfId": "123456"
                                 },
                                 {
                                     "plmnId": {"mcc": "461", "mnc": "01"},
                                     "amfId": "654321"
                                 }
                             ],
                             "taiList": [
                                 {
                                     "plmnId": {"mcc": "462", "mnc": "002"},
                                     "tac": "234567"
                                 },
                                 {
                                     "plmnId": {"mcc": "463", "mnc": "03"},
                                     "tac": "765432"
                                 }
                             ]
                         }`)
	amfRegionId, problemDetails = getAmfRegionId(amfInfo)

	if problemDetails == nil {
		t.Fatalf("should return error, but not.")
	}
}

func TestGetAusfProperties(t *testing.T) {
	nfProfile := []byte(`{
	                         "nfInstanceID":"ausf",
						    "nfType":"AUSF",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0
					    }`)
	groupId, routingIndicator, problemDetails := getAusfProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("should not return error, but did.")
	}

	if groupId != "" || routingIndicator != "" {
		t.Fatalf("getAusfProperties didn't return right value.")
	}

	nfProfile = []byte(`{
	                         "nfInstanceID":"ausf",
						    "nfType":"AUSF",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
							"ausfInfo": {
    						        "groupId": "group1",
    						        "routingIndicator": "1234",
								"supiRanges": {
									"start": "123",
									"end": "456"
								}
  						    }
					    }`)
	groupId, routingIndicator, problemDetails = getAusfProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("should not return error, but did.")
	}

	if groupId != "group1" || routingIndicator != "1234" {
		t.Fatalf("getAusfProperties didn't return right value.")
	}
}

func TestGetAusfGroupId(t *testing.T) {
	ausfInfo := []byte(`{
    						    "groupId": "group1",
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	groupId, _ := getAusfGroupId(ausfInfo)

	if groupId != "group1" {
		t.Fatalf("getAusfGroupId didn't return right groupId")
	}

	ausfInfo = []byte(`{
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	groupId, _ = getAusfGroupId(ausfInfo)

	if groupId != "" {
		t.Fatalf("getAusfGroupId didn't return right groupId")
	}
}

func TestGetAusfRoutingIndicator(t *testing.T) {
	ausfInfo := []byte(`{
    						    "groupId": "group1",
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	routingIndicator, _ := getAusfRoutingIndicator(ausfInfo)
	if routingIndicator != "1234" {
		t.Fatalf("getAusfRoutingIndicator didn't return right routingIndicator")
	}

	ausfInfo = []byte(`{
    						    "groupId": "group1",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	routingIndicator, _ = getAusfRoutingIndicator(ausfInfo)
	if routingIndicator != "" {
		t.Fatalf("getAusfRoutingIndicator didn't return right routingIndicator")
	}
}

func TestGetSmfPgwFqdn(t *testing.T) {
	smfInfo := []byte(`{
    							"pgwFqdn": "111111",
    							"dnnList": [
      							"dnn1",
      							"dnn2"
    							],
    							"taiList": [
      							{
        								"plmnId": {"mcc": "460", "mnc": "000"},
        								"tac": "123456"
      							},
      							{
        								"plmnId": {"mcc": "461", "mnc": "01"},
        								"tac": "654321"
      							}
    							]
  						}`)
	pgwFqdn, _ := getSmfPgwFqdn(smfInfo)
	if pgwFqdn != "111111" {
		t.Fatalf("getSmfPgwFqdn didn't return right pgwFqdn. It should be 111111")
	}

	smfInfo = []byte(`{
    						  "dnnList": [
      						  "dnn1",
      					      "dnn2"
    						  ],
    						  "taiList": [
      						  {
        							  "plmnId": {"mcc": "460", "mnc": "000"},
        						      "tac": "123456"
      						  },
      						  {
        							  "plmnId": {"mcc": "461", "mnc": "01"},
        							  "tac": "654321"
      						  }
    						  ]
  				      }`)
	pgwFqdn, _ = getSmfPgwFqdn(smfInfo)
	if pgwFqdn != "" {
		t.Fatalf("getSmfPgwFqdn didn't return right pgwFqdn. It should be empty")
	}
}

func TestGetUdmProperties(t *testing.T) {
	nfProfile := []byte(`{
	                         "nfInstanceID":"udm",
						    "nfType":"UDM",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0
					    }`)
	groupId, routingIndicator, problemDetails := getUdmProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("getUdmProperties should not return error, but did.")
	}

	if groupId != "" || routingIndicator != "" {
		t.Fatalf("getUdmProperties didn't return right value.")
	}

	nfProfile = []byte(`{
	                         "nfInstanceID":"udm",
						    "nfType":"UDM",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
							"udmInfo": {
    						        "groupId": "group1",
    						        "routingIndicator": "1234",
								"supiRanges": {
									"start": "123",
									"end": "456"
								}
  						    }
					    }`)
	groupId, routingIndicator, problemDetails = getUdmProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("getUdmProperties should not return error, but did.")
	}

	if groupId != "group1" || routingIndicator != "1234" {
		t.Fatalf("getUdmProperties didn't return right value.")
	}
}

func TestGetUdmGroupId(t *testing.T) {
	udmInfo := []byte(`{
    						    "groupId": "group1",
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	groupId, _ := getUdmGroupId(udmInfo)

	if groupId != "group1" {
		t.Fatalf("getUdmGroupId didn't return right groupId. It should be group1.")
	}

	udmInfo = []byte(`{
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	groupId, _ = getUdmGroupId(udmInfo)

	if groupId != "" {
		t.Fatalf("getUdmGroupId didn't return right groupId. It should be empty.")
	}
}

func TestGetUdmRoutingIndicator(t *testing.T) {
	udmInfo := []byte(`{
    						    "groupId": "group1",
    						    "routingIndicator": "1234",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	routingIndicator, _ := getUdmRoutingIndicator(udmInfo)

	if routingIndicator != "1234" {
		t.Fatalf("getUdmRoutingIndicator didn't return right routingIndicator. It should be 1234.")
	}

	udmInfo = []byte(`{
    						    "groupId": "group1",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	routingIndicator, _ = getUdmRoutingIndicator(udmInfo)

	if routingIndicator != "" {
		t.Fatalf("getUdmRoutingIndicator didn't return right routingIndicator. It should be empty.")
	}
}

func TestGetUdrProperties(t *testing.T) {
	nfProfile := []byte(`{
	                         "nfInstanceID":"udr",
						    "nfType":"UDR",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0
					    }`)
	groupId, problemDetails := getUdrProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("getUdrProperties should not return error, but did.")
	}

	if groupId != "" {
		t.Fatalf("getUdrProperties didn't return right value.")
	}

	nfProfile = []byte(`{
	                         "nfInstanceID":"udr",
						    "nfType":"UDR",
						    "plmn": "46000", 
						    "sNssai": {"sst": "0","sd": "0"},
						    "fqdn": "amf.mnc001.mcc460.5g",
						    "ipAddress": ["10.0.0.1"],
						    "capacity": 0,
							"udrInfo": {
    						        "groupId": "group1",
								"supiRanges": {
									"start": "123",
									"end": "456"
								}
  						    }
					    }`)
	groupId, problemDetails = getUdrProperties(nfProfile)

	if problemDetails != nil {
		t.Fatalf("getUdrProperties should not return error, but did.")
	}

	if groupId != "group1" {
		t.Fatalf("getUdrProperties didn't return right value.")
	}
}

func TestGetUdrGroupId(t *testing.T) {
	udrInfo := []byte(`{
    						    "groupId": "group1",
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  						}`)

	groupId, _ := getUdrGroupId(udrInfo)

	if groupId != "group1" {
		t.Fatalf("getUdrGroupId didn't return right groupId. It should be group1.")
	}

	udrInfo = []byte(`{
							"supiRanges": {
								"start": "123",
								"end": "456"
							}
  					  }`)

	groupId, _ = getUdrGroupId(udrInfo)

	if groupId != "" {
		t.Fatalf("getUdrGroupId didn't return right groupId. It should be empty.")
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

func TestLastSubString(t *testing.T) {}

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

func TestrequestFromNRFProv(t *testing.T) {
	flag := RequestFromNRFProv("eric-nrf-management:3004", 3004)
	if !flag {
		t.Fatalf("requestFromNRFProv failed")
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
