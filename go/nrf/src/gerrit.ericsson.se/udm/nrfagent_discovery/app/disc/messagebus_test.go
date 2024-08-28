package disc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/schema"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"github.com/buger/jsonparser"
)

var nfProfileUDMWithInvalidInterFqdn = []byte(`
{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
		"nfStatus": "REGISTERED",
        "plmnList": [{
            "mcc": "466",
            "mnc": "000"
        }],
        "sNssais": [{
                "sst": 2,
                "sd": "A00000"
            },
            {
                "sst": 4,
                "sd": "A00000"
            }
        ],
        "nsiList": ["069"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
		"interPlmnFqdn" : "seliius03695.seli.gic.ericsson.se",
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

var nfProfileUDM = []byte(`
{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
		"nfStatus": "REGISTERED",
        "plmnList": [{
            "mcc": "466",
            "mnc": "000"
        }],
        "sNssais": [{
                "sst": 2,
                "sd": "A00000"
            },
            {
                "sst": 4,
                "sd": "A00000"
            }
        ],
        "nsiList": ["069"],
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

var searchRNFProfile = []byte(`
{
	"validityPeriod": 86400,
     "nfInstances":[{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
		"nfStatus": "REGISTERED",
        "plmnList": [{
            "mcc": "466",
            "mnc": "000"
        }],
        "sNssais": [{
                "sst": 2,
                "sd": "A00000"
            },
            {
                "sst": 4,
                "sd": "A00000"
            }
        ],
        "nsiList": ["069"],
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
}`)

var msgChange = []byte(`
{"event": "NF_PROFILE_CHANGED", "nfType":"UDM", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "messageBody": {
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "udm-5g-01",
        "nfType": "UDM",
		"nfStatus": "REGISTERED",
        "plmnList": [{
            "mcc": "460",
            "mnc": "000"
        }],
        "sNssais": [{
                "sst": 2,
                "sd": "A00002"
            },
            {
                "sst": 4,
                "sd": "A00004"
            }
        ],
        "nsiList": ["069"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
        "ipv4Addresses": ["172.16.208.1"],
        "ipv6Addresses": ["FF01::1101"],
        "ipv6Prefixes": ["2001:db8:abcd:12::0/64"],
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
        "amfInfo": {
            "amfSetId": "amfSet-01",
            "amfRegionId": "amfRegion-01",
            "guamiList": [{
                "plmnId": {
                    "mcc": "460",
                    "mnc": "000"
                },
                "amfId": "a00001"
            }]
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
        "pcfInfo": {
            "dnnlist": ["pcf-dnn1","pcf-dnn2"]
        },
        "bsfInfo": {
            "ipv4AdddressRanges": [{
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
			"nfServiceStatus" : "REGISTED",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "FF01::1101",
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
}
 `)

var msgChangePatchBody = []byte(`
{"event": "NF_PROFILE_CHANGED", "nfType":"UDM", "reqNfType": "AUSF", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "messageBody":[
 { "op": "REPLACE","path": "/capacity", "newValue": 50 },
 { "op": "remove","path": "/ipv6Addresses"}
]}`)

var msgChangeRoamPatchBody = []byte(`
{"event": "NF_PROFILE_CHANGED", "nfType":"UDM", "reqNfType": "AUSF-roam", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "messageBody":[
 { "op": "REPLACE","path": "/capacity", "newValue": 50 },
 { "op": "remove","path": "/ipv6Addresses"}
]}`)

var msgChangePatchWithInterPlmnFqdn = []byte(`
{"event": "NF_PROFILE_CHANGED", "nfType":"UDM", "reqNfType": "AUSF", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "messageBody":[
 { "op": "add","path": "/interPlmnFqdn", "newValue": "seliius03695.seli.gic.ericsson.se" },
 { "op": "replace","path": "/capacity", "newValue": 30 }
]}`)

var msgChangeRoamPatchWithInterPlmnFqdn = []byte(`
{"event": "NF_PROFILE_CHANGED", "nfType":"UDM", "reqNfType": "AUSF-roam", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "messageBody":[
 { "op": "add","path": "/interPlmnFqdn", "newValue": "seliius03695.seli.gic.ericsson.se" },
 { "op": "replace","path": "/capacity", "newValue": 30 }
]}`)

var msgDeregRoam = []byte(`
{"event": "NF_DEREGISTERED", "nfType":"UDM", "reqNfType": "AUSF-roam", "nfInstanceId": "udm-5g-01", "agentProducerId": "string"}`)

var msgChangePatchApplyBody = []byte(`
[
 { "op": "replace","path": "/capacity", "value": 50 },
 { "op": "remove","path": "/ipv6Addresses"}
]
`)

var msgDeReg = []byte(`
{"event": "NF_DEREGISTERED", "nfType": "UDM", "nfInstanceId": "udm-5g-01", "agentProducerId": "string", "MessageBody": {
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
                "sd": "A00001"
            },
            {
                "sst": 4,
                "sd": "A00001"
            }
        ],
        "nsiList": ["069"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
        "ipv4Addresses": ["172.16.208.1"],
        "ipv6Addresses": ["FF01::1101"],
        "ipv6Prefixes": ["2001:db8:abcd:12::0/64"],
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
                "amfId": "amf001"
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
        "pcfInfo": {
            "dnnlist":  ["pcf-dnn1","pcf-dnn2"]
        },
        "bsfInfo": {
            "ipv4AdddressRanges": [{
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
            "version": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "schema": "https://",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "FF01::1101",
                "ipv6Prefix": "2001:db8:abcd:12::0/64",
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
}        `)

var msgAgentEvent = []byte(`
{"event": "AGENT_EVENT_WITHOUT_BODY", "nfInstanceId": "\"udm-5g-01\"", "agentProducerId": "string"}`)

var msgRegEvent = []byte(`{"eventType": "REGISTER", "nfInstanceId": "\"udm-5g-01\"", "nfType": "AUSF","fqdn": "ausf_fqdn"}`)
var msgDeRegEvent = []byte(`{"eventType": "DEREGISTER", "nfInstanceId": "\"udm-5g-01\"", "nfType": "AUSF"}`)
var msgFqdnChgEvent = []byte(`{"eventType": "FQDN_CHANGED", "nfInstanceId": "\"udm-5g-01\"", "nfType": "AUSF","fqdn": "ausf_new_fqdn"}`)

//var randnum = 1

func getnfProfileSchamaLoader() {
	var nfSchemaSuffix = "src/gerrit.ericsson.se/udm/nrfagent_common/helm/eric-nrfagent-common/config/schema/nfProfileInSearchResult.json"
	var patchSchemaSuffix = "src/gerrit.ericsson.se/udm/nrfagent_common/helm/eric-nrfagent-common/config/schema/patchDocument.json"
	var nfSchema string
	var patchSchema string

	goPath := os.Getenv("GOPATH")
	//fmt.Println("goPath", goPath)
	nfSchema = goPath + "/" + nfSchemaSuffix
	patchSchema = goPath + "/" + patchSchemaSuffix
	//fmt.Println("nfSchema: ", nfSchema)

	nfSchemaContent, err := ioutil.ReadFile(nfSchema)
	if err != nil {
		fmt.Errorf("Load schema file failure, err:%v\n", err)
	}

	patchSchemaContent, err := ioutil.ReadFile(patchSchema)
	if err != nil {
		fmt.Printf("Load patchschema file failure, err:%v\n", err.Error())
	}

	mapSchemaFile := make(map[string][]byte)
	mapSchemaFile["nfProfileInSearchResult.json"] = nfSchemaContent
	mapSchemaFile["patchDocument.json"] = patchSchemaContent

	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	for key := range mapSchemaFile {
		AbsFileName := currDir + "/" + key
		err := ioutil.WriteFile(AbsFileName, mapSchemaFile[key], 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	os.Setenv("SCHEMA_DIR", currDir)
	os.Setenv("SCHEMA_NF_PROFILE", "nfProfileInSearchResult.json")
	os.Setenv("SCHEMA_PATCH_DOCUMENT", "patchDocument.json")

	errInfo := schema.LoadDiscoverSchema()
	if errInfo != nil {
		fmt.Printf("LoadDiscoverSchema error,error info is %+v", err.Error())
	}
}

func TestConvertToSearchResultNFProfile(t *testing.T) {
	var sNFProfileConvertData map[string]interface{}
	var sNFProfileStandardData map[string]interface{}
	sNFProfile, err := convertToSearchResultNFProfile(nfProfileUDMWithInvalidInterFqdn)
	err = json.Unmarshal(sNFProfile, &sNFProfileConvertData)
	if err != nil {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult Unmarshal sNFProfileConvertData failed, error is %s", err.Error())
	}
	err = json.Unmarshal(nfProfileUDM, &sNFProfileStandardData)
	if err != nil {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult Unmarshal sNFProfileStandardData failed, error is %s", err.Error())
	}
	ok := reflect.DeepEqual(sNFProfileConvertData, sNFProfileStandardData)
	if !ok {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult failed, error is %s", err.Error())
	}
}

func TestConvertToSearchResultBody(t *testing.T) {
	var sNFProfileConvertBody map[string]interface{}
	var sNFProfileStandardBody map[string]interface{}
	sNFProfile, err := convertToSearchResultBody(nfProfileUDM)
	err = json.Unmarshal(sNFProfile, &sNFProfileConvertBody)
	if err != nil {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult Unmarshal sNFProfileConvertData failed, error is %s", err.Error())
	}
	err = json.Unmarshal(searchRNFProfile, &sNFProfileStandardBody)
	if err != nil {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult Unmarshal sNFProfileStandardData failed, error is %s", err.Error())
	}
	ok := reflect.DeepEqual(sNFProfileConvertBody, sNFProfileStandardBody)
	if !ok {
		t.Fatalf("TestConvertToSearchResult: convertToSearchResult failed, error is %s", err.Error())
	}
}

func TestApplyPatchItems(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestApplyPatchItems: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestApplyPatchItems: Cached fail")
		}
	}

	nfProfileInCache := cacheManager.FetchNfProfile("AUSF", "udm-5g-01")
	updatedNfProfile, ok := applyPatchItems(nfProfileInCache, msgChangePatchApplyBody)
	if ok != true {
		log.Error("TestApplyPatchItems: apply Patch Items error and return")
		return
	}

	capacity, err := jsonparser.GetInt(updatedNfProfile, "capacity")
	if err != nil || capacity != 50 {
		t.Errorf("TestApplyPatchItems: NF_PROFILE_CHANGED fail")
	}

	ipv6Addresses, _, _, err := jsonparser.Get(updatedNfProfile, "ipv6Addresses")
	if err == nil || ipv6Addresses != nil {
		t.Errorf("TestApplyPatchItems: NF_PROFILE_CHANGED fail")
	}

}

func TestNtfDiscInnerMessageHandler(t *testing.T) {
	getnfProfileSchamaLoader()
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestNtfDiscInnerHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestNtfDiscInnerHandler: Cached fail")
		}
	}

	t.Run("testApplyPatchNormally", func(t *testing.T) {
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", msgChangePatchBody)
		content := cacheManager.DumpByID("AUSF", "UDM", "udm-5g-01")
		capacity, err := jsonparser.GetInt(content, "capacity")
		if err != nil || capacity != 50 {
			t.Errorf("TestNtfDiscInnerHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}

		ipv6Addresses, _, _, err := jsonparser.Get(content, "ipv6Addresses")
		if err == nil || ipv6Addresses != nil {
			t.Errorf("TestNtfDiscInnerHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}
	})

	t.Run("testApplyPatchWithInvalidField_InterPlmnFqdn", func(t *testing.T) {
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", msgChangePatchWithInterPlmnFqdn)
		content := cacheManager.DumpByID("AUSF", "UDM", "udm-5g-01")
		_, err := jsonparser.GetString(content, "interPlmnFqdn")
		if err == nil {
			t.Errorf("TestNtfDiscInnerHandler: testApplyPatchWithInvalidField_InterPlmnFqdn NF_PROFILE_CHANGED failed")
		}
		capacity, err := jsonparser.GetInt(content, "capacity")
		if err != nil || capacity != 30 {
			t.Errorf("TestNtfDiscInnerHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}
	})
	/*
		notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", msgDeReg)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if ok {
			t.Errorf("TestNtfDiscInnerHandler: NF_DEREGISTERED fail")
		}
	*/
	cacheManager.Flush("AUSF")
}

func TestNtfDiscInnerRoamMessageHandler(t *testing.T) {
	getnfProfileSchamaLoader()
	log.SetLevel(log.ErrorLevel)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestNtfDiscInnerRoamMessageHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: Cached fail")
		}
	}

	t.Run("testApplyPatchNormally", func(t *testing.T) {
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", msgChangeRoamPatchBody)
		content := cacheManager.DumpRoamingByID("AUSF", "UDM", "udm-5g-01")
		capacity, err := jsonparser.GetInt(content, "capacity")
		if err != nil || capacity != 50 {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}

		ipv6Addresses, _, _, err := jsonparser.Get(content, "ipv6Addresses")
		if err == nil || ipv6Addresses != nil {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}
	})

	t.Run("testApplyPatchWithInvalidField_InterPlmnFqdn", func(t *testing.T) {
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", msgChangeRoamPatchWithInterPlmnFqdn)
		content := cacheManager.DumpRoamingByID("AUSF", "UDM", "udm-5g-01")
		_, err := jsonparser.GetString(content, "interPlmnFqdn")
		if err == nil {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testApplyPatchWithInvalidField_InterPlmnFqdn NF_PROFILE_CHANGED failed")
		}
		capacity, err := jsonparser.GetInt(content, "capacity")
		if err != nil || capacity != 30 {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testApplyPatchNormally NF_PROFILE_CHANGED fail")
		}
	})
	t.Run("testCompeteUpdateRoamProfile", func(t *testing.T) {
		var ntfDiscMsg structs.NtfDiscInnerMsg
		err := json.Unmarshal(msgChangeRoamPatchBody, &ntfDiscMsg)
		if err != nil {
			t.Errorf("Unmarshal msgChangeRoamPatchBody failure, %s", err.Error())
		}
		var nfProfile structs.NfProfile
		err = json.Unmarshal(nfProfileUDM, &nfProfile)
		if err != nil {
			t.Errorf("Unmarshal nfProfileUDM failure, %s", err.Error())
		}
		ntfDiscMsg.MessageBody = make([]structs.NfProfilePatchData, 0)
		ntfDiscMsg.NfProfile = &nfProfile
		changeMsgBody, err := json.Marshal(ntfDiscMsg)
		if err != nil {
			t.Errorf("Marshal faliure, %s", err.Error())
		}
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", changeMsgBody)
		content := cacheManager.DumpRoamingByID("AUSF", "UDM", "udm-5g-01")
		capacity, err := jsonparser.GetInt(content, "capacity")
		if err != nil || capacity != 100 {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testCompeteUpdateRoamProfile fail")
		}
	})

	t.Run("testDeleteRoamProfile", func(t *testing.T) {
		ntfDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+"ntfdiscinner", msgDeregRoam)
		content := cacheManager.DumpRoamingByID("AUSF", "UDM", "udm-5g-01")
		if len(content) != 0 {
			t.Errorf("TestNtfDiscInnerRoamMessageHandler: testDeleteRoamProfile fail")
		}
	})

	cacheManager.FlushRoam("AUSF")
}

func TestNtfMessageHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	cacheManager.SetRequesterFqdn("AUSF", "ausf_fqdn")
	var ntfMsg structs.NotificationMsg
	err := json.Unmarshal(msgChange, &ntfMsg)
	if err != nil {
		t.Errorf("Unmarshal msgChange failure, %s", err.Error())
	}
	ntfMsg.NfEvent = consts.NFRegister
	ntfMsgBody, err := json.Marshal(ntfMsg)
	if err != nil {
		t.Errorf("Marshal ntfMsg faliure, %s", err.Error())
	}
	t.Run("testChangeProfile", func(t *testing.T) {

		notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", ntfMsgBody)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TesNnotificationMessageHandler: Cached fail")
		}
	})

	t.Run("testChangeProfile", func(t *testing.T) {
		notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", msgChange)
		content := cacheManager.DumpByID("AUSF", "UDM", "udm-5g-01")
		indexContent, err := jsonparser.GetString(content, "nfType")
		if err != nil || indexContent != "UDM" {
			t.Errorf("TesNnotificationMessageHandler: NF_PROFILE_CHANGED fail")
		}
	})

	t.Run("testDeregProfile", func(t *testing.T) {
		notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", msgDeReg)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if ok {
			t.Errorf("TesNnotificationMessageHandler: NF_DEREGISTERED fail")
		}
	})

	t.Run("testDiscResultEventHandler", func(t *testing.T) {
		ntfMsg.NfEvent = consts.NFEventDiscResult
		discResultBody, err := json.Marshal(ntfMsg)
		if err != nil {
			t.Errorf("Marshal ntfMsg faliure, %s", err.Error())
		}
		notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", discResultBody)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TesNnotificationMessageHandler: testDiscResultEventHandler fail")
		}
	})
	cacheManager.Flush("AUSF")
}

func StubHTTPServerHandler200(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("location", "http://127.0.1.1:3214/nrf-nfm/v1/nf-instances/ausf-5g-01/subscribeId01")
	rw.WriteHeader(http.StatusOK)
	rw.Write(searchResultUDM)
	fmt.Println("Recieved Request.")
}

func TestNotificationMessageHandlerWithoutBody(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	h := httpserver.InitHTTPServer(
		httpserver.Trace(true),
		httpserver.HostPort("", "3211"),
		httpserver.HTTP2(true),
		httpserver.MaxConcurrentStreams(5),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.SetRoute(),
		httpserver.PathFunc("/nnrf-disc/v1/nf-instances", "GET", StubHTTPServerHandler200),
	)
	h.Run()

	defer func() { h.Stop() }()
	client.InitHttpClient()

	event := "testEvent"
	cfgName := "cfgName"

	//test format value
	format := cmproxy.NtfFormatFull
	cmTargetNfProfilesHandler(event, cfgName, format, cmTargetNfProfile)
	_, ok := common.CmGetTargetNfProfile()
	if !ok {
		t.Errorf("TestNotificationMessageHandlerWithoutBody: CmTargetNfProfilesHandler format check failure.")
	}

	cmNrfAgentConfHandler(event, cfgName, format, cmDataNrfService)

	nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
	if nfinstanceByte == nil {
		t.Errorf("TestNotificationMessageHandlerWithoutBody: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestNotificationMessageHandlerWithoutBody: Cached fail")
		}
	}

	notificationMessageHandler(consts.MsgbusTopicNamePrefix+"AUSF", msgAgentEvent)
	content := cacheManager.DumpByID("AUSF", "UDM", "udm-5g-01")
	indexContent, err := jsonparser.GetString(content, "nfType")
	if err != nil || indexContent != "UDM" {
		t.Errorf("TestNotificationMessageHandlerWithoutBody: nfType check failure.")
	}

	cacheManager.Flush("AUSF")
}

func TestRegDiscInnerMessageHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	client.InitHttpClient()
	activeLeaderMock(true)
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusOK)
	StubHTTPDoToNrf("POST", http.StatusCreated)

	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	})

	t.Run("TestEventTypeRegister", func(t *testing.T) {
		regDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.RegDiscInner, msgRegEvent)
		fqdn, ok := cacheManager.GetRequesterFqdn("AUSF")

		if !ok || fqdn != "ausf_fqdn" {
			t.Errorf("TesRregDiscInnerMessageHandler  register check failure")
		}
	})

	t.Run("TestEventTypeDeregister", func(t *testing.T) {
		regDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.RegDiscInner, msgFqdnChgEvent)
		fqdn, ok := cacheManager.GetRequesterFqdn("AUSF")

		if !ok || fqdn != "ausf_new_fqdn" {
			t.Errorf("TesRregDiscInnerMessageHandler  deregister check failure")
		}
	})
	t.Run("TestEventTypeDeregister", func(t *testing.T) {
		regDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.RegDiscInner, msgDeRegEvent)
		_, ok := cacheManager.GetRequesterFqdn("AUSF")
		if ok {
			t.Errorf("TesRregDiscInnerMessageHandler  deregister check failure")
		}
	})
	cacheManager.Flush("AUSF")
}

/*
func TestSendSubscrInfoToSlave(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	subInfo := structs.SubscriptionInfo{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "udm-01",
		SubscriptionID:    "123-456-789",
		ValidityTime:      time.Time{},
	}

	ok := sendSubscrInfoToSlave(subInfo)
	if ok {
		t.Errorf("TestSendSubscrInfoToSlave: check failure")
	}
}
*/
func TestDiscDiscInnerMessageHandler(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	ok := syncSubscrInfoEventHandler(nil)
	if !ok {
		t.Errorf("TestSyncSubscrInfoEventHandler: check active disc action failure")
	}

	activeLeaderMock(false)

	t.Run("TestSyncSubscrInfoWrong", func(t *testing.T) {
		subInfoWrong := structs.SubscriptionInfo{
			RequesterNfType:   "UDR",
			TargetNfType:      "UDM",
			TargetServiceName: "udm-01",
			SubscriptionID:    "123-456-789",
			ValidityTime:      time.Time{},
		}
		syncSubscrInfoWrongMsg := structs.DiscDiscInnerMsg{
			EventType:       consts.EventTypeSyncSubscrInfo,
			AgentProducerID: "string",
			SubscrInfo:      subInfoWrong,
		}

		discDiscWrongBody, err := json.Marshal(syncSubscrInfoWrongMsg)
		if err != nil {
			t.Errorf("Marshal faliure, %s", err.Error())
		}
		discDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.DiscDiscInner, discDiscWrongBody)
		ids, ok := cacheManager.GetSubscriptionIDs("UDR", "UDM")
		if len(ids) != 0 || ok {
			t.Errorf("TestSyncSubscrInfoEventHandler: check slave disc message handle check failure")
		}
	})

	t.Run("TestSyncSubscrInfoMsg", func(t *testing.T) {
		subInfo := structs.SubscriptionInfo{
			RequesterNfType:   "AUSF",
			TargetNfType:      "UDM",
			TargetServiceName: "udm-01",
			SubscriptionID:    "123-456-789",
			ValidityTime:      time.Time{},
		}
		syncSubscrInfoMsg := structs.DiscDiscInnerMsg{
			EventType:       consts.EventTypeSyncSubscrInfo,
			AgentProducerID: "string",
			SubscrInfo:      subInfo,
		}

		discDiscBody, err := json.Marshal(syncSubscrInfoMsg)
		if err != nil {
			t.Errorf("Marshal faliure, %s", err.Error())
		}

		discDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.DiscDiscInner, discDiscBody)
		ids, ok := cacheManager.GetSubscriptionIDs("AUSF", "UDM")

		if !ok || len(ids) != 1 || ids[0] != "123-456-789" {
			t.Errorf("TestSyncSubscrInfoEventHandler: check slave disc message handle check failure")
		}
		cacheManager.DelSubscriptionInfo("AUSF", "UDM", "123-456-789")
	})

	t.Run("TestSyncRoamSubscrInfoMsg", func(t *testing.T) {
		subInfo := structs.SubscriptionInfo{
			RequesterNfType: "AUSF",
			TargetNfType:    "UDM",
			NfInstanceID:    "udm-5g-01",
			SubscriptionID:  "111-222-333",
			ValidityTime:    time.Time{},
		}
		syncSubscrInfoMsg := structs.DiscDiscInnerMsg{
			EventType:       consts.EventTypeSyncSubscrInfo,
			AgentProducerID: "string",
			SubscrInfo:      subInfo,
		}

		discDiscBody, err := json.Marshal(syncSubscrInfoMsg)
		if err != nil {
			t.Errorf("Marshal faliure, %s", err.Error())
		}

		discDiscInnerMessageHandler(consts.MsgbusTopicNamePrefix+consts.DiscDiscInner, discDiscBody)
		ids, ok := cacheManager.GetRoamingSubscriptionIDs("AUSF", "UDM")

		if !ok || len(ids) != 1 || ids[0] != "111-222-333" {
			t.Errorf("TestSyncSubscrInfoEventHandler: check slave disc message handle check failure")
		}
		cacheManager.DelRoamingSubscriptionInfo("AUSF", "UDM", "123-222-333")
	})
	activeLeaderMock(true)
}
