package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestValidateCommon(t *testing.T) {
	//NF profile with fqdn is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateCommon() != nil {
		t.Fatalf("TNFProfile.ValidateCommon didn't return right value!")
	}

	//NF profile with ipv4Addresses is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"ipv4Addresses": [ 
		    "10.10.10.10",
			"10.10.10.11"
		],
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateCommon() != nil {
		t.Fatalf("TNFProfile.ValidateCommon didn't return right value!")
	}

	//NF profile with ipv6Addresses is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"ipv6Addresses": [ 
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		],
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateCommon() != nil {
		t.Fatalf("TNFProfile.ValidateCommon didn't return right value!")
	}

	//NF profile with fqdn and ipv4Addresses and ipv6Addresses is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"ipv4Addresses": [ 
		    "10.10.10.10",
			"10.10.10.11"
		],
		"ipv6Addresses": [ 
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		],
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateCommon() != nil {
		t.Fatalf("TNFProfile.ValidateCommon didn't return right value!")
	}

	//NF profile without fqdn and ipv4Addresses and ipv6Addresses is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateCommon() == nil {
		t.Fatalf("TNFProfile.ValidateCommon didn't return right value!")
	}
}

func TestValidateService(t *testing.T) {
	//NF profile without nfServices is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com"
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() != nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who doesn't include ipEndPoints and chfServiceInfo is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"			
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() != nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who include valid ipEndPoints and chfServiceInfo is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
				"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01"
				}				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
				"chfServiceInfo": {
					"secondaryChfServiceInstance": "serv02"
				}				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() != nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who include invalid ipEndPoints is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
				"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01"
				}				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10",
						"ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
				"chfServiceInfo": {
					"secondaryChfServiceInstance": "serv02"
				}				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() == nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who include invalid chfServiceInfo is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
				"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01",
					"secondaryChfServiceInstance": "serv02"
				}			
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
			    	"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01",
					"secondaryChfServiceInstance": "serv02"
				}	
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() == nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who include invalid ipEndPoints and chfServiceInfo is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
			    	"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01",
					"secondaryChfServiceInstance": "serv02"
				}				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10",
						"ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ],
			    	"chfServiceInfo": {
					"primaryChfServiceInstance": "serv01",
					"secondaryChfServiceInstance": "serv02"
				}				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() == nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who includ matched serviceName is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-comm",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-evts",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"	
			},
			{
				"serviceInstanceId": "namf-03",
				"serviceName": "namf-mt",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"	
			},
			{
				"serviceInstanceId": "namf-04",
				"serviceName": "namf-loc",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"	
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() != nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}

	//NF profile with nfServices who includ unmatched serviceName is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-xxxx",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-comm",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"	
			},
			{
				"serviceInstanceId": "namf-03",
				"serviceName": "nausf-auth",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"	
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateService() == nil {
		t.Fatalf("TNFProfile.ValidateService didn't return right value!")
	}
}

func TestValidateAmf(t *testing.T) {
	//NF profile without amfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAmf() != nil {
		t.Fatalf("TNFProfile.ValidateAmf didn't return right value!")
	}

	//NF profile with amfInfo who includes valid taiRangeList and n2InterfaceAmfInfo is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                     },
                     "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }  
	        ],
			"n2InterfaceAmfInfo": {
	            "ipv4EndpointAddress": [
	    	    	    	    "10.10.10.10"
	  	    	    ],
	  	    	    "ipv6EndpointAddress": [
	    	    	    	    "1030::C9B4:FF12:48AA:1A2B"
	  	    	    ]
		    	}
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAmf() != nil {
		t.Fatalf("TNFProfile.ValidateAmf didn't return right value!")
	}

	//NF profile with amfInfo who has invalid taiRangeList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                    },
                    "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
				{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  		    {
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }
	        ]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAmf() == nil {
		t.Fatalf("TNFProfile.ValidateAmf didn't return right value!")
	}

	//NF profile with amfInfo who includes invalid n2InterfaceAmfInfo is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                    },
                    "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
				{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }
	        ],
	        "n2InterfaceAmfInfo": {
	        }
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAmf() == nil {
		t.Fatalf("TNFProfile.ValidateAmf didn't return right value!")
	}

	//NF profile with amfInfo who includes invalid taiRangeList and n2InterfaceAmfInfo is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                    },
                    "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
				{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  		    {
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }
	        ],
	        "n2InterfaceAmfInfo": {
	        }
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}
}

func TestValidateAusf(t *testing.T) {
	//NF profile without ausfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AUSF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAusf() != nil {
		t.Fatalf("TNFProfile.ValidateAusf didn't return right value!")
	}

	//NF profile with ausfInfo who doesn't include supiRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AUSF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"ausfInfo": {
			"groupId": "001"
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAusf() != nil {
		t.Fatalf("TNFProfile.ValidateAusf didn't return right value!")
	}

	//NF profile with ausfInfo who include valid supiRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AUSF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"ausfInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAusf() != nil {
		t.Fatalf("TNFProfile.ValidateAusf didn't return right value!")
	}

	//NF profile with ausfInfo who include invalid supiRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AUSF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"ausfInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateAusf() == nil {
		t.Fatalf("TNFProfile.ValidateAusf didn't return right value!")
	}

}

func TestValidateChf(t *testing.T) {
	//NF profile without chfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "CHF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateChf() != nil {
		t.Fatalf("TNFProfile.ValidateChf didn't return right value!")
	}

	//NF profile with chfInfo who doesn't include supiRangeList, gpsiRangeList and plmnRangeList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "CHF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"chfInfo": {
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateChf() != nil {
		t.Fatalf("TNFProfile.ValidateChf didn't return right value!")
	}

	//NF profile with chfInfo who include valid supiRangeList, gpsiRangeList and plmnRangeList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "CHF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"chfInfo": {
            "gpsiRangeList": [
			    {
				    "start": "1111",
				    "end": "2222"
			    },
			    {
			        "pattern": "string"
			    }
		    ],
	        "supiRangeList": [
			    {
				    "start": "1111",
				    "end": "2222"
			    },
			    {
			        "pattern": "string"
			    }
		    ],
		    "plmnRangeList": [
			    {
				    "start": "1111",
				    "end": "2222"
			    },
			    {
			        "pattern": "string"
			    }
		    ]
		}
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateChf() != nil {
		t.Fatalf("TNFProfile.ValidateChf didn't return right value!")
	}

	//NF profile with chfInfo who include invalid supiRangeList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "CHF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"chfInfo": {
			"supiRangeList": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRangeList": [
			    {
				    "start": "1111",
				    "end": "2222"
			    },
			    {
			        "pattern": "string"
			    }
		    ],
		    "plmnRangeList": [
			    {
				    "start": "1111",
				    "end": "2222"
			    },
			    {
			        "pattern": "string"
			    }
		    ]
		}
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateChf() == nil {
		t.Fatalf("TNFProfile.ValidateChf didn't return right value!")
	}

	//NF profile with chfInfo who include invalid supiRangeList or gpsiRangeList or plmnRangeList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "CHF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"chfInfo": {
			"supiRangeList": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRangeList": [
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
		    ],
		    "plmnRangeList": [
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
		    ]
		}
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateChf() == nil {
		t.Fatalf("TNFProfile.ValidateChf didn't return right value!")
	}
}

func TestValidatePcf(t *testing.T) {
	//NF profile without pcfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "PCF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidatePcf() != nil {
		t.Fatalf("TNFProfile.ValidatePcf didn't return right value!")
	}

	//NF profile with pcfInfo who doesn't include supiRangeList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "PCF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"pcfInfo": {
			"dnnList": ["dnn1", "dnn2"]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidatePcf() != nil {
		t.Fatalf("TNFProfile.ValidatePcf didn't return right value!")
	}

	//NF profile with pcfInfo who include valid supiRangeList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "PCF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"pcfInfo": {
			"dnnList": ["dnn1", "dnn2"],
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidatePcf() != nil {
		t.Fatalf("TNFProfile.ValidatePcf didn't return right value!")
	}

	//NF profile with pcfInfo who include invalid supiRangeList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "PCF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"pcfInfo": {
			"dnnList": ["dnn1", "dnn2"],
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidatePcf() == nil {
		t.Fatalf("TNFProfile.ValidatePcf didn't return right value!")
	}
}

func TestValidateSmf(t *testing.T) {
	//NF profile without smfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "SMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateSmf() != nil {
		t.Fatalf("TNFProfile.ValidateSmf didn't return right value!")
	}

	//NF profile with smfInfo who doesn't include taiRangeList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "SMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"smfInfo": {
			"dnnList": ["dnn1", "dnn2"]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateSmf() != nil {
		t.Fatalf("TNFProfile.ValidateSmf didn't return right value!")
	}

	//NF profile with valid smfInfo is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "SMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"smfInfo": {
			"dnnList": ["dnn1", "dnn2"],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }  
	        ]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateSmf() != nil {
		t.Fatalf("TNFProfile.ValidateSmf didn't return right value!")
	}

	//NF profile with smfInfo who includes invalid taiRangeList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "SMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"smfInfo": {
			"dnnList": ["dnn1", "dnn2"],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
				{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  		    {
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }
	        ]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateSmf() == nil {
		t.Fatalf("TNFProfile.ValidateSmf didn't return right value!")
	}
}

func TestValidateUdm(t *testing.T) {
	//NF profile without udmInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() != nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who doesn't include supiRanges and gpsiRanges and externalGroupIdentifiersRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001"
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() != nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who include valid supiRanges and gpsiRanges and externalGroupIdentifiersRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"externalGroupIdentifiersRanges ": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() != nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who include invalid supiRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() == nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who include invalid gpsiRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001",
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() == nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who include invalid externalGroupIdentifiersRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001",
			"externalGroupIdentifiersRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() == nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}

	//NF profile with udmInfo who include invalid supiRanges and invalid gpsiRanges and invalid externalGroupIdentifiersRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDM",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udmInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"externalGroupIdentifiersRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdm() == nil {
		t.Fatalf("TNFProfile.ValidateUdm didn't return right value!")
	}
}

func TestValidateUdr(t *testing.T) {
	//NF profile without udrInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() != nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who doesn't include supiRanges and gpsiRanges and externalGroupIdentifiersRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001"
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() != nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who include valid supiRanges and gpsiRanges and externalGroupIdentifiersRanges is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"externalGroupIdentifiersRanges ": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() != nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who include invalid supiRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() == nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who includes invalid gpsiRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001",
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() == nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who includes invalid externalGroupIdentifiersRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001",
			"externalGroupIdentifiersRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() == nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}

	//NF profile with udrInfo who includes invalid supiRanges and invalid gpsiRanges and invalid externalGroupIdentifiersRanges is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "001",
			"supiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"gpsiRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			],
			"externalGroupIdentifiersRanges": [
			    {
					"start": "1111",
					"end": "2222"
				},
				{
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
				},
				{
					"start": "1111",
					"end": "2222",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    },
				{
					"end": "2222",
					"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
				},
				{
					"start": "1111",
				    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
			    }
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUdr() == nil {
		t.Fatalf("TNFProfile.ValidateUdr didn't return right value!")
	}
}

func TestValidateUpf(t *testing.T) {
	//NF profile without upfInfo is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UPF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUpf() != nil {
		t.Fatalf("TNFProfile.ValidateUpf didn't return right value!")
	}

	//NF profile with upfInfo who doesn't include interfaceUpfInfoList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UPF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"upfInfo": {
			"smfServingArea": ["000", "111"]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUpf() != nil {
		t.Fatalf("TNFProfile.ValidateUpf didn't return right value!")
	}

	//NF profile with upfInfo who includes valid interfaceUpfInfoList is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UPF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"upfInfo": {
			"smfServingArea": ["000", "111"],
			"interfaceUpfInfoList": [
			    {
					"interfaceType": "N3",
					"ipv4EndpointAddresses": [
					    "10.10.10.10",
						"10.10.10.11"
					]
				},
				{
					"interfaceType": "N3",
				    "ipv6EndpointAddresses": [
		                "1030::C9B4:FF12:48AA:1A2B",
			            "1030::C9B4:FF12:48AA:1A2B"
		            ]
				},
				{
					"interfaceType": "N3",
					"endpointFqdn": "http://test"
				}
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUpf() != nil {
		t.Fatalf("TNFProfile.ValidateUpf didn't return right value!")
	}

	//NF profile with upfInfo who includes invalid interfaceUpfInfoList is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UPF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"upfInfo": {
			"smfServingArea": ["000", "111"],
			"interfaceUpfInfoList": [
			    {
					"interfaceType": "N3",
					"ipv4EndpointAddresses": [
					    "10.10.10.10",
						"10.10.10.11"
					]
				},
				{
					"interfaceType": "N3",
				    "ipv6EndpointAddresses": [
		                "1030::C9B4:FF12:48AA:1A2B",
			            "1030::C9B4:FF12:48AA:1A2B"
		            ]
				},
				{
					"interfaceType": "N3",
					"endpointFqdn": "http://test"
				},
				{
					"interfaceType": "N3"
				},
				{
					"interfaceType": "N3",
					"ipv4EndpointAddresses": [
					    "10.10.10.10",
						"10.10.10.11"
					],
				    "ipv6EndpointAddresses": [
		                "1030::C9B4:FF12:48AA:1A2B",
			            "1030::C9B4:FF12:48AA:1A2B"
		            ]
				},
				{
					"interfaceType": "N3"
				}
			]
		},
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.ValidateUpf() == nil {
		t.Fatalf("TNFProfile.ValidateUpf didn't return right value!")
	}
}

func TestValidate(t *testing.T) {
	//NF profile with fqdn is valid
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "NRF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": []
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() != nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}

	//NF profile without fqdn and ipv4Addresses and ipv6Addresses is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "NRF",
		"nfStatus": "REGISTERED",
		"nfServices": []
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() == nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}

	//NF profile with nfServices who include valid ipEndPoints is valid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "NRF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() != nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}

	//NF profile with nfServices who include invalid ipEndPoints is invalid
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "NRF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
						"ipv4Address": "10.10.10.10",
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10",
						"ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() == nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}

	//a valid AMF profile
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                     },
                     "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }  
	        ],
			"n2InterfaceAmfInfo": {
	            "ipv4EndpointAddress": [
	    	    	    	    "10.10.10.10"
	  	    	    ],
	  	    	    "ipv6EndpointAddress": [
	    	    	    	    "1030::C9B4:FF12:48AA:1A2B"
	  	    	    ]
		    	}
		},
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() != nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}

	//a invalid AMF profile
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfSet01",
            "amfRegionId": "amfRegion01",
            "guamiList": [
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "02"
                    },
                    "amfId": "800002"
                },
                {
                    "plmnId": {
                        "mcc": "460",
                        "mnc": "05"
                    },
                    "amfId": "800005"
                }
            ],
			"taiRangeList": [
	  			{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
				{
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            },
	  		    {
		  			"plmnId": {
            				"mcc": "460",
            				"mnc": "05"
                     },
		             "tacRangeList": [
		  		        {
			                "start": "1234",
			                "end": "1234"
		                },
		                {
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                },
					    {
		                },
		                {
                             "start": "1234",
			                "end": "1234",
			                "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		                }
		            ]
	            }
	        ],
	        "n2InterfaceAmfInfo": {
	        }
		},
		"nfServices": [
		    {
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			},
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"ipEndPoints": [
		            {
				        "port": 80
			        },
		            {
				        "port": 80,
				        "ipv4Address": "10.10.10.10"
			        },
			        {
				        "port": 80,
				        "ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			        }
		        ]				
			}
		]
	}`)

	profile = &TNFProfile{}
	err = json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if profile.Validate() == nil {
		t.Fatalf("TNFProfile.Validate didn't return right value!")
	}
}

func TestGenerateMd5ForNFProfile(t *testing.T) {
	//two NF profiles between which only the NFServices is different, their md5 shall be the same
	body1 := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	body2 := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "https",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	nfProfile1 := &TNFProfile{}
	err := json.Unmarshal(body1, nfProfile1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	nfProfile2 := &TNFProfile{}
	err = json.Unmarshal(body2, nfProfile2)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile1.GenerateMd5() != nfProfile2.GenerateMd5() {
		t.Fatalf("TNFProfile.GenerateMd5 didn't return right value!")
	}

	//two NF profiles between which attributes except for the NFServices is different, their md5 shall be different
	body1 = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	body2 = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "DEREGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	nfProfile1 = &TNFProfile{}
	err = json.Unmarshal(body1, nfProfile1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	nfProfile2 = &TNFProfile{}
	err = json.Unmarshal(body2, nfProfile2)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile1.GenerateMd5() == nfProfile2.GenerateMd5() {
		t.Fatalf("TNFProfile.GenerateMd5 didn't return right value!")
	}
}

func TestCreateSnssaisHelperInfo(t *testing.T) {
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AUSF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"sNssais": [
			{
				"sst": 3,
				"sd": "111111"
			},
			{
				"sst": 2
			}
		]
	}`)

	profile := &TNFProfile{}
	err := json.Unmarshal(body, profile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}
	result := `"sNssais":[{"sst":3,"sd":"111111"},{"sst":2,"sd":"RESERVED_EMPTY_SD"}]`
	if result != profile.createSnssaisHelperInfo() {
		t.Fatal("func createSnssaisHelperInfo() create snssais helper generate wrong")
	}
}

func TestNFProfileEqual(t *testing.T) {
	//two NF profiles between which only the NFServices is different, their md5 shall be the same
	body1 := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"heartbeatTimer" : 120,
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	nfProfile1 := &TNFProfile{}
	err := json.Unmarshal(body1, nfProfile1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	//case1 : body2 is the same with body1
	body2 := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"heartbeatTimer" : 120,
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	nfProfile2 := &TNFProfile{}
	err = json.Unmarshal(body2, nfProfile2)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile1.Equal(nfProfile2) != true {
		t.Fatalf("Profile should be equal, but Not!")
	}

	//case 2 : body3 is different with body1
	body3 := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"nfServices": [
			{
				"serviceInstanceId": "namf-02",
				"serviceName": "namf-02",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED"				
			}
		]
	}`)

	nfProfile3 := &TNFProfile{}
	err = json.Unmarshal(body3, nfProfile3)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile1.Equal(nfProfile3) == true {
		t.Fatalf("Profile should Not be equal, but Not!")
	}
}

func TestIsNfInfoExist(t *testing.T) {
	//case 1 : NF type whose nfInfo doesn't belong to nrfInfo return false
	body := []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "UDR",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"udrInfo": {
			"groupId": "udm01",
			"supiRanges": [
			    {
					"start": "111111111",
					"end": "222222222"
				}
			]
		}
	}`)

	nfProfile := &TNFProfile{}
	err := json.Unmarshal(body, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile.IsNfInfoExist() {
		t.Fatalf("TNFProfile.IsNfInfoExist should return false, but Not!")
	}

	//case 2 : NF type whose nfInfo belongs to nrfInfo, but doesn't exist, return false
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com"
	}`)

	nfProfile = &TNFProfile{}
	err = json.Unmarshal(body, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfProfile.IsNfInfoExist() {
		t.Fatalf("TNFProfile.IsNfInfoExist should return false, but Not!")
	}

	//case 1 : NF type whose nfInfo belongs to nrfInfo, and exist, return true
	body = []byte(`{
		"nfInstanceId": "0000-1111-2222-3333",
		"nfType": "AMF",
		"nfStatus": "REGISTERED",
		"fqdn": "http://www.test.com",
		"amfInfo": {
			"amfSetId": "amfset01",
			"amfRegionId": "amfRegion01",
			"guamiList": [
			    {
					"amfId": "111111",
					"plmnId": {
						"mcc": "460",
						"mnc": "000"
					}
				}
			]
		}
	}`)

	nfProfile = &TNFProfile{}
	err = json.Unmarshal(body, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !nfProfile.IsNfInfoExist() {
		t.Fatalf("TNFProfile.IsNfInfoExist should return true, but Not!")
	}
}

func TestAllowedParametersInNFProfile(t *testing.T) {
	body1 := []byte(`
	{
	"nfInstanceId": "0000-1111-2222-3333",
	"nfType": "AMF",
	"heartbeatTimer" : 120,
	"nfStatus": "REGISTERED",
	"fqdn": "http://www.test.com",
	"allowedPlmns": [
				    {
						"mcc": "460",
						"mnc": "01"
					},
					 {
						"mcc": "460",
						"mnc": "00"
					 }
				],
	"allowedNfTypes": ["AUSF", "AMF"],
	"allowedNfDomains": ["^seliius\\d{5}.seli.gic.ericsson.se$","^seliius\\d{4}.seli.gic.ericsson.se$"],
	"allowedNssais": [
		{
			"sst": 100,
			"sd": "111111"
		},
		{
			"sst": 200,
			"sd": "22222"
		},
		{
			"sst": 300
		}
	],
	"nfServices": [
		{
			"serviceInstanceId": "namf-01",
			"serviceName": "namf-01",
			"versions": [
				{
					"apiVersionInUri": "http://test",
					"apiFullVersion": "0.1"
				}
			],
			"scheme": "http",
			"nfServiceStatus": "REGISTERED"
					
		}
	]
    }`)

	nfProfile1 := &TNFProfile{}
	err := json.Unmarshal(body1, nfProfile1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !nfProfile1.IsAllowedNfType("AUSF") || !nfProfile1.IsAllowedNfType("AMF") {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfProfile1.IsAllowedNfType("UDR") {
		t.Fatalf("Should not allow, but YES!")
	}

	if !nfProfile1.IsAllowedNfDomain("seliius2121.seli.gic.ericsson.se") || !nfProfile1.IsAllowedNfDomain("seliius22121.seli.gic.ericsson.se") {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfProfile1.IsAllowedNfDomain("seliius121.seli.gic.ericsson.se") {
		t.Fatalf("Should not allow, but YES!")
	}

	//IsAllowedPlmn
	plmnID1 := &TPlmnId{
		Mcc: "460",
		Mnc: "00",
	}
	plmnID2 := &TPlmnId{
		Mcc: "460",
		Mnc: "01",
	}
	plmnID3 := &TPlmnId{
		Mcc: "460",
		Mnc: "02",
	}

	if !nfProfile1.IsAllowedPlmn(plmnID1) || !nfProfile1.IsAllowedPlmn(plmnID2) {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfProfile1.IsAllowedPlmn(plmnID3) {
		t.Fatalf("Should not allow, but YES!")
	}

	// IsAllowedSNssi

	snssai1 := &TSnssai{
		Sst: 100,
		Sd:  "111111",
	}

	snssai2 := &TSnssai{
		Sst: 200,
		Sd:  "111111",
	}

	snssai3 := &TSnssai{
		Sst: 300,
		Sd:  "111111",
	}

	if !nfProfile1.IsAllowedNssai(snssai1) || !nfProfile1.IsAllowedNssai(snssai3) {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfProfile1.IsAllowedNssai(snssai2) {
		t.Fatalf("Should not allow, but YES!")
	}
}
