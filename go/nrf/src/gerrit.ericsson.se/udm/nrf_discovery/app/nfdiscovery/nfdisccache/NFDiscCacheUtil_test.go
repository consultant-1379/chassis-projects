package nfdisccache

import (
	"testing"
	"com/dbproxy/nfmessage/nrfprofile"
	"time"
)

func TestSplitNFProfileList(t *testing.T) {
	t.Fatal("split nfprofile return a wrong value")
}
func TestSplitNRFProfileList(t *testing.T) {
	rawNRFProfile := (`{
        "capacity": 100,
        "fqdn": "seliius03696.seli.gic.ericsson.se",
        "nfInstanceId": "12345678-abcd-ef12-1000-000000000010",
        "nfServices": [
            {
                "capacity": 100,
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.1",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-01",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            },
            {
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.2",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-02",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            },
            {
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.3",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-sdm-01",
                "serviceName": "nudm-sdm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            }
        ],
        "nfStatus": "REGISTERED",
        "nfType": "NRF",
        "plmn": {
            "mcc": "460",
            "mnc": "000"
        },
        "priority": 100,
        "sNssais": [
            {
                "sd": "222222",
                "sst": 2
            }
        ],
        "nrfInfo": {
            "externalGroupIdentifiersRanges": [
                {
                    "pattern": "^groupid-[A-Fa-f0-9]{8}-[0-9]{3}-[0-9]{2,3}-([A-Fa-f0-9][A-Fa-f0-9]){1}$"
                }
            ],
            "gpsiRanges": [
                {
                    "end": "423456789059999",
                    "start": "423456789040000"
                }
            ],
            "routingIndicators": [
                "1234"
            ],
            "supiRanges": [
                {
                    "end": "666669999",
                    "start": "600000000"
                }
            ]
        }
    }
	`)
	nrfProfileInfo := nrfprofile.NRFProfileInfo{RawNrfProfile:[]byte(rawNRFProfile), ExpiredTime:1545288737000}
	cacheItems := SplitNRFProfileList([]*nrfprofile.NRFProfileInfo{&nrfProfileInfo})
	if cacheItems[0].Key != "12345678-abcd-ef12-1000-000000000010" && cacheItems[0].ProfileUpdateTime != 1545288737000 {
		t.Fatal("split nfprofile return a wrong value")
	}
}

func TestGetNFProfileFromCache(t *testing.T) {
	customNFProfile := (`{
    "expiredTime": 1545288862000,
    "profileUpdateTime": 1545288737000,
    "provisioned": 1,
    "md5sum": {
        "nfProfile": "25d7476ff58fec0e0975a5d0edb9f0ab",
        "nudm-uecm-01": "740bca9170d7756e9c9d0ff96b9041fe",
        "nudm-uecm-02": "5f97289663520ec178dbb6247fb126ce",
        "nudm-sdm-01": "c91f05388da3c67c3d1f94559be06cdf"
    },
    "helper": {
        "nfServices": {
            "allowedNfDomains": [
                "RESERVED_EMPTY_DOMAIN"
            ],
            "allowedNfTypes": [
                "RESERVED_EMPTY_TYPE"
            ],
            "allowedPlmns": [
                {
                    "mcc": "XXX",
                    "mnc": "YYY"
                }
            ]
        },
        "xxxInfo": {
		}
        "sNssais": [
            {
                "sst": 2,
                "sd": "222222"
            }
        ]
    },
    "body": {
        "capacity": 100,
        "fqdn": "seliius03696.seli.gic.ericsson.se",
        "nfInstanceId": "12345678-abcd-ef12-1000-000000000010",
        "nfServices": [
            {
                "capacity": 100,
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.1",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-01",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            },
            {
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.2",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-uecm-02",
                "serviceName": "nudm-uecm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            },
            {
                "fqdn": "seliius03690.seli.gic.ericsson.se",
                "ipEndPoints": [
                    {
                        "ipv4Address": "172.16.208.3",
                        "port": 30088
                    }
                ],
                "nfServiceStatus": "REGISTERED",
                "priority": 100,
                "scheme": "https",
                "serviceInstanceId": "nudm-sdm-01",
                "serviceName": "nudm-sdm",
                "versions": [
                    {
                        "apiFullVersion": "1.R15.1.1",
                        "apiVersionInUri": "v1",
                        "expiry": "2020-07-06T02: 54: 32Z"
                    }
                ]
            }
        ],
        "nfStatus": "REGISTERED",
        "nfType": "UDM",
        "plmn": {
            "mcc": "460",
            "mnc": "000"
        },
        "priority": 100,
        "sNssais": [
            {
                "sd": "222222",
                "sst": 2
            }
        ],
        "udmInfo": {
            "externalGroupIdentifiersRanges": [
                {
                    "pattern": "^groupid-[A-Fa-f0-9]{8}-[0-9]{3}-[0-9]{2,3}-([A-Fa-f0-9][A-Fa-f0-9]){1}$"
                }
            ],
            "gpsiRanges": [
                {
                    "end": "423456789059999",
                    "start": "423456789040000"
                }
            ],
            "routingIndicators": [
                "1234"
            ],
            "supiRanges": [
                {
                    "end": "666669999",
                    "start": "600000000"
                }
            ]
        }
    }
}
	`)
	InitNfProfileCache()
	NfProfileCache.AddDataChannel <- []string{customNFProfile}
	time.Sleep(1 * time.Second)
	var keys []CacheNFResponse
	keys = append(keys, CacheNFResponse{NfInstanceID:"12345678-abcd-ef12-1000-000000000010", ProfileUpdateTime:1545288737000})
	keys = append(keys, CacheNFResponse{NfInstanceID:"12345678-abcd-ef12-1000-000000000011", ProfileUpdateTime:1545288737000})
	cacheItems, notfoundKeys := GetNFProfileFromCache(keys)
	if len(cacheItems) != 1 && len(notfoundKeys) != 1 && notfoundKeys[0] != "12345678-abcd-ef12-1000-000000000011" {
		t.Fatal("get cache nfprofile fail")
	}
}
