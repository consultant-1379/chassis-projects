package schema

import (
	"fmt"
	"os"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

func init() {
	log.SetLevel(log.FatalLevel)

	goPath := os.Getenv("GOPATH")
	os.Setenv("SCHEMA_DIR", fmt.Sprintf("%s/src/gerrit.ericsson.se/udm/nrf_common/helm/eric-nrf-common/config/schema", goPath))
	os.Setenv("SCHEMA_NF_PROFILE", "nfProfile.json")
	os.Setenv("SCHEMA_PATCH_DOCUMENT", "patchDocument.json")
	os.Setenv("SCHEMA_SUBSCRIPTIONDATA", "subscriptionData.json")
	os.Setenv("SCHEMA_SUBSCRIPTIONPATCH", "subscriptionPatch.json")
}

func TestLoadManagementSchema(t *testing.T) {
	err := LoadManagementSchema()
	if err != nil {
		t.Fatalf("LoadManagementSchema error, %v", err)
	}

	if schemaNfProfile == nil || schemaPatchDocument == nil || schemaSubscriptionData == nil || schemaSubscriptionPatch == nil {
		t.Fatalf("LoadManagementSchema error, %v", err)
	}
}

func TestValidateNfProfile(t *testing.T) {

	LoadManagementSchema()

	problemDetails := ValidateNfProfile(`{
        "nfInstanceId": "udm1",
        "nfType": "UDM",
        "nfStatus": "REGISTERED",
		"plmnList": [
		    {
				"mcc": "460",
				"mnc": "00"
			},
			{
				"mcc": "460",
				"mnc": "00"
			},
			{
				"mcc": "460",
				"mnc": "00"
			}
		],
        "sNssais": [
            {
                "sst": 1,
                "sd": "123456"
            },
            {
                "sst": 21,
                "sd": "123456"
            }
        ],
        "fqdn": "testing",
        "interPlmnFqdn": "testing",
        "ipv4Addresses": [
            "10.10.10.100"
        ],
        "ipv6Addresses": [
            "2001:db8:85a3::8a2e:370:7334"
        ],
        "capacity": 0,
        "nrfInfo": {
            "servedUdrInfo": {
                "udr1": {
                    "supiRanges": [
                        {
                            "start": "001",
                            "end": "002",
                            "pattern": "string"
                        }
                    ]
                }
            }
        },
        "udrInfo": {
            "supiRanges": [
                {
                    "start": "001",
                    "end": "002",
                    "pattern": "string"
                }
            ]
        },
        "amfInfo": {
            "amfSetId": "string",
			"amfRegionId": "string",
			"guamiList": [
			    {
					"plmnId": {
					    "mcc": "460",
						"mnc": "00"
					},
					"amfId": "123456"
				}
			]
        },
        "smfInfo": {
            "sNssaiSmfInfoList": [
			    {
				    "sNssai": {
						"sst": 1,
                         "sd": "123456"
					},
					"dnnSmfInfoList": [
					    {
							"dnn": "123"
						}	
					]
			    }
			]
        },
        "upfInfo": {
            "sNssaiUpfInfoList": [
                {
                    "sNssai": {
                        "sst": 3,
                        "sd": "123456"
                    },
                    "dnnUpfInfoList": [
                        {
                            "dnn": "dnn3"
                        }
                    ]
                }
            ]
        },
		"bsfInfo": {
			"ipv6PrefixRanges": [
			    {
					"start": "2001:db8:abcd:12::0/64",
					"end": "2001:db8:abcd:12::0/64"
				}
			]
		},
        "nfServices": [
            {
                "serviceInstanceId": "srv1",
                "serviceName": "srv1",
                "versions": [
				    {
						"apiVersionInUri": "http://test6",
						"apiFullVersion": "0.1"
					}
				],
                "scheme": "http",
				"nfServiceStatus": "SUSPENDED",
                "fqdn": "fqdn1",
                "interPlmnFqdn": "fqdn1",
                "ipEndPoints": [
                    {
                       "ipv4Address": "255.255.255.0",
                       "ipv6Address": "fda:5cc1:23:4::1f",
                       "port": 0
                    }
                ],
                "apiPrefix": "string",
                "defaultNotificationSubscriptions": [
                    {
                        "callbackUri": "string",
                        "notificationType": "N1_MESSAGES"
                    }
                ],
                "allowedPlmns": [
                    {
                        "mcc": "000",
                        "mnc": "00"
                    },
                    {
                        "mcc": "001",
                        "mnc": "01"
                    },
                    {
                        "mcc": "002",
                        "mnc": "02"
                    }
                ],
                "allowedNfTypes": [
                    "AUSF",
                    "UDM"
                ],
                "allowedNssais": [
                    {
                        "sst": 0,
                        "sd": "123456"
                    }
                ],
                "capacity": 0
            }
        ]
    }`)

	if problemDetails != nil {
		t.Fatalf("This is a valid nf profile, but validate failed! %+v", problemDetails.InvalidParams[0])
	}

	problemDetails = ValidateNfProfile(`{
        "nnfInstanceId": "udm1",
        "nfType": "UDM",
        "nfStatus": "REGISTERED",
        "sNssais": [
            {
                "sst": 1,
                "sd": "sd1"
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
            "10.10.10.10.10"
        ],
        "ipv6Prefixes": [
            "ipv6"
        ],
        "capacity": 0,
        "udrInfo": {
            "supiRanges": [
                {
                    "start": "string",
                    "end": "string",
                    "pattern": "string"
                }
            ]
        },
        "amfInfo": {
            "amfSetId": "string"
        },
        "smfInfo": {
            "dnnList": [
                "dnn1",
                "dnn2"
            ],
            "servingArea": [
                "string"
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
                            "dnn": "dnn3"
                        }
                    ]
                }
            ]
        },
        "nfServices": [
            {
                "serviceInstanceId": "srv1",
                "serviceName": "srv1",
                "version": "v1",
                "schema": "schema1",
                "fqdn": "fqdn1",
                "interPlmnFqdn": "fqdn1",
                "ipEndPoints": [
                    {
                       "ipv4Address": "string",
                       "ipv6Address": "string",
                       "ipv6Prefix": "string",
                       "port": 0
                    }
                ],
                "apiPrefix": "string",
                "defaultNotificationSubscriptions": [
                    {
                        "callbackUri": "string",
                        "notificationType": "UDM"
                    }
                ],
                "allowedPlmns": [
                    {
                        "mcc": "000",
                        "mnc": "00"
                    },
                    {
                        "mcc": "001",
                        "mnc": "01"
                    },
                    {
                        "mcc": "002",
                        "mnc": "02"
                    }
                ],
                "allowedNfTypes": [
                    "AUSF",
                    "UDM"
                ],
                "allowedNssais": [
                    {
                        "sst": 0,
                        "sd": "string"
                    }
                ],
                "capacity": 0
            }
        ]
    }`)

	if problemDetails == nil {
		t.Fatalf("This is a invalid nf profile, but validate ok!")
	}
}

func TestValidatePatchDocument(t *testing.T) {
	LoadManagementSchema()

	problemDetails := ValidatePatchDocument(`[
        {
            "op": "replace",
            "path": "/ipAddress",
            "value": ["10.0.0.3"]
        },
        {
            "op": "add",
            "path": "/nfServiceList/0/allowedNfTypes",
            "value": ["nrf", "amf", "ausf"]
        }
    ]`)

	if problemDetails != nil {
		t.Fatalf("This is a valid nf patchDocument, but validate failed!")
	}

	problemDetails = ValidatePatchDocument(`[
        {
            "oop": "replace",
            "path": "/ipAddress",
            "value": ["10.0.0.3"]
        },
        {
            "op": "add",
            "path": "/nfServiceList/0/allowedNfTypes",
            "value": ["nrf", "amf", "ausf"]
        }
    ]`)

	if problemDetails == nil {
		t.Fatalf("This is a invalid nf patchDocument, but validate ok!")
	}

	problemDetails = ValidatePatchDocument(`[]`)

	if problemDetails == nil {
		t.Fatalf("This is a invalid nf patchDocument, but validate ok!")
	}
}

func TestValidateSubscriptionData(t *testing.T) {
	LoadManagementSchema()

	problemDetails := ValidateSubscriptionData(`{
		"nfStatusNotificationUri": "http://seliius04099:20001",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED", "NF_PROFILE_CHANGED"],
		"plmnId": {
			"mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
		    "nfType": "UDM",
            "serviceName": "nudm-ueau",
			"amfSetId": "amfSet01",
			"amfRegionId": "amfRegion01",
			"guamiList": [
			    {
					"plmnId": {
						"mcc": "460",
						"mnc": "00"
					},
					"amfId": "123456"
				},
				{
					"plmnId": {
						"mcc": "460",
						"mnc": "00"
					},
					"amfId": "123456"					
				}
			],
			"snssaiList": [
			    {
					"sst": 1,
					"sd": "123456"
				},
				{
					"sst": 2,
					"sd": "123456"
				}
			],
			"nsiList": ["nsi01", "nsi02"],
			"nfGroupId": "group01"
        },
		"validityTime": "2019-01-09T02:31:25Z",
		"notifCondition": {
			"monitoredAttributes": ["fqdn", "ipv4Addresses"],
			"unmonitoredAttributes": ["load", "priority"]
		},
		"reqNfType": "UDM",
		"reqNfFqdn": "xxxx"
    }`)

	if problemDetails != nil {
		t.Fatalf("This is a valid nf subscriptionData, but validate failed!")
	}

	problemDetails = ValidateSubscriptionData(`{
        "nfType": "UDMM",
        "serviceName": "nudm-ueau",
        "callbackUri": "http://seliius04099:20001"
    }`)

	if problemDetails == nil {
		t.Fatalf("This is a invalid nf subscriptionData, but validate ok!")
	}

}

func TestValidateSubscriptionPatch(t *testing.T) {
	LoadManagementSchema()

	// patch body must include at least one item
	patchBody := `[
	]`

	problemDetails := ValidateSubscriptionPatch(patchBody)

	if problemDetails == nil {
		t.Fatalf("This is a invalid subscription patch body, but validate ok!")
	}

	// op should only be replace
	patchBody = `[
	    {
			"op": "add",
			"path": "/validityTime",
			"value": "2019-01-18T02:02:51Z"
		}
	]`

	problemDetails = ValidateSubscriptionPatch(patchBody)

	if problemDetails == nil {
		t.Fatalf("This is a invalid subscription patch body, but validate ok!")
	}

	// path should only be /validityTime
	patchBody = `[
	    {
			"op": "replace",
			"path": "/reqNfType",
			"value": "AMF"
		}
	]`

	problemDetails = ValidateSubscriptionPatch(patchBody)

	if problemDetails == nil {
		t.Fatalf("This is a invalid subscription patch body, but validate ok!")
	}

	// invalid value format
	patchBody = `[
	    {
			"op": "replace",
			"path": "/validityTime",
			"value": "20191212"
		}
	]`

	problemDetails = ValidateSubscriptionPatch(patchBody)

	if problemDetails == nil {
		t.Fatalf("This is a invalid subscription patch body, but validate ok!")
	}

	// patch body must include at most one item
	patchBody = `[
	    {
			"op": "replace",
			"path": "/validityTime",
			"value": "2019-01-18T02:02:51Z"
		},
		{
			"op": "replace",
			"path": "/validityTime",
			"value": "2019-01-18T02:02:51Z"
		}
	]`

	problemDetails = ValidateSubscriptionPatch(patchBody)

	if problemDetails == nil {
		t.Fatalf("This is a invalid subscription patch body, but validate ok!")
	}

	// a valid patch body
	patchBody = `[
	    {
			"op": "replace",
			"path": "/validityTime",
			"value": "2019-01-18T02:02:51Z"
		}
	]`

	problemDetails = ValidateSubscriptionPatch(patchBody)

	if problemDetails != nil {
		t.Fatalf("This is a valid subscription patch body, but validate failed!")
	}
}
