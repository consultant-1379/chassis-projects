package nfdiscfilter

import "testing"

func TestIsMatchedSnssais(t *testing.T) {
	filter := &NFNRFInfoFilter{}
	rawNrfProfile := []byte(`{
	    "nsiList": [
		"12345"
		     ],
	    "nfSetId": "5112345",
	    "nfStatus": "REGISTERED",
	    "interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
	    "nfType": "NRF",
	    "fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
	    "locality": "",
	    "plmn": {
		"mcc": "460",
		"mnc": "000"
		},
	    "priority": 10,
	    "nfInstanceId": "0c765084-9cc5-49c6-9876-ae2f5fa2a63f",
	    "ipv6Addresses": [],
	    "capacity": 100,
	    "ipv4Addresses": [ "127.0.0.1"],
	    "sNssais": [
		     {
				"sst": 1,
				"sd": "0"
			},
		     {
				"sst": 0,
				"sd": "1"
			},
			{
				"sst": 0,
				"sd": "Ab1"
			}
			],
	    "ipv6Prefixes": [],
	    "nfServices": [
		{
		"interPlmnFqdn": "seliius03696.seli.gic.ericsson.se",
		"scheme": "https",
		"fqdn": "nrf.5gc.mnc000.mcc460.3gppnetwork.org",
		"serviceInstanceId": "nnrf-nfm-01",
		"supportedFeatures": "1F",
		"ipEndPoints": [
		    {
				    "transport": "TCP",
				    "ipv4Address": "127.0.0.1",
				    "ipv6Address": "192.168.0.1",
				    "port": 443
				}
			    ],
		"apiPrefix": "",
		"priority": 5,
		"version": [
				{
				    "apiVersionInUri": "v1",
				    "apiFullVersion": "1.R15.1.1",
				    "expiry": "2020-07-06T02:54:32Z"
				}
			    ],
		"capacity": 100,
		"serviceName": "nnrf-disc",
		"allowedNfDomains": [],
		"allowedPlmns": [
				{
				    "mcc": "460",
				    "mnc": "00"
				}
			    ],
		"allowedNfTypes": [ ],
		"allowedNssais": [
				{
				    "sst": 1,
				    "sd": "0"
				}
			    ]
			}
	    ],
	    "nrfInfo":{
		    "udmInfoSum": {
			"supiRanges": [
				{
					"start": "",
					"end": "",
					"pattern": ""
				}
			],
			"externalGroupIdentifiersRanges": [
				{
					"start": "123000000",
					"end": "123999999",
					"pattern": "^imsi-12345\\d{4}$"
				}
			],
			"gpsiRanges": [
				{
					"start": "",
					"end": "",
					"pattern": ""
				}
			],
			"groupIdList": [
				"",
				""
			],
			"routingIndicatorList": [
				"",
				""
			]
		},
	    }
	}`)
	snssais := "[{\"sst\": 1,\"sd\": \"1\"},{\"sst\": 2,\"sd\": \"2\"}]"
	matched := filter.isMatchedSnssais(snssais, rawNrfProfile)
	if matched {
		t.Fatal("snssais should not be matched, but matched")
	}
	snssais2 := "[{\"sst\": 1,\"sd\": \"0\"},{\"sst\": 2,\"sd\": \"2\"}]"
	matched2 := filter.isMatchedSnssais(snssais2, rawNrfProfile)
	if !matched2 {
		t.Fatal("snssais should be matched, but not matched")
	}
	snssais3 := "[{\"sst\": 1,\"sd\": \"0\"}]"
	matched3 := filter.isMatchedSnssais(snssais3, rawNrfProfile)
	if !matched3 {
		t.Fatal("snssais should be matched, but not matched")
	}
	snssais4 := "[{\"sst\": 0,\"sd\": \"aB1\"}]"
	matched4 := filter.isMatchedSnssais(snssais4, rawNrfProfile)
	if !matched4 {
		t.Fatal("snssais should be matched, but not matched")
	}
}
