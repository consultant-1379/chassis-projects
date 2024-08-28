package cache

import (
	"bytes"
	"strconv"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

func StubSetupBenchmark() {
	//ResetLogLevel
	log.SetLevel(log.ErrorLevel)

	//Setup Common Env
	//	var cacheManager *CacheManager
	//	cacheManager.SetCacheConfigUT("../../build/config/cache-index.json")
}

func newSearchParam(id string) *SearchParameter {
	searchParameter := &SearchParameter{}
	searchParameter.SetServiceNames([]string{"nudm-uecm"})
	searchParameter.SetTargetNfType("UDM")
	searchParameter.SetRequesterNfType("AUSF")
	if id != "" {
		searchParameter.SetTargetNfInstanceID(id)
	}

	return searchParameter
}

var searchResultNfProfileUDM = []byte(`{
    "capacity" : 100,
    "fqdn" : "seliius03696.seli.gic.ericsson.se",
    "ipv4Addresses" : [ "172.16.208.1", "172.16.208.2", "172.16.208.3", "172.16.208.4", "172.16.208.5", "172.16.208.6", "172.16.208.7", "172.16.208.8" ],
    "nfInstanceId" : "12345678-9udm-def0-1000-0000",
    "nfServices" : [ {
      "defaultNotificationSubscriptions" : [ {
        "callbackUri" : "/nnrf-nfm/v1/nf-instances/udm-5g-01",
        "n1MessageClass" : "5GMM",
        "n2InformationClass" : "SM",
        "notificationType" : "N1_MESSAGES"
      } ],
      "fqdn" : "seliius03690.seli.gic.ericsson.se",
      "ipEndPoints" : [ {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30088,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0002",
        "port" : 30090,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      } ],
      "nfServiceStatus" : "REGISTERED",
      "scheme" : "https",
      "serviceInstanceId" : "nudm-uecm-01",
      "serviceName" : "nudm-uecm",
      "supportedFeatures" : "A0A0",
      "versions" : [ {
        "apiFullVersion" : "1.R15.1.1 ",
        "apiVersionInUri" : "v1",
        "expiry" : "2020-07-06T02:54:32Z"
      } ]
    },
    {
      "fqdn" : "seliius03690.seli.gic.ericsson.se",
      "ipEndPoints" : [ {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30088,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0002",
        "port" : 30090,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      }, {
        "ipv6Address" : "fe80:1234:0000:0000:0000:0000:0000:0000",
        "port" : 30089,
        "transport" : "TCP"
      } ],
      "nfServiceStatus" : "REGISTERED",
      "scheme" : "https",
      "serviceInstanceId" : "nudm-uecm-02",
      "serviceName" : "nudm-uecm",
      "supportedFeatures" : "A0A0",
      "versions" : [ {
        "apiFullVersion" : "1.R15.1.1 ",
        "apiVersionInUri" : "v1",
        "expiry" : "2020-07-06T02:54:32Z"
      } ]
    } ],
    "nfStatus" : "REGISTERED",
    "nfType" : "UDM",
    "nsiList" : [ "111111", "222222", "333333" ],
    "plmn" : {
      "mcc" : "460",
      "mnc" : "00"
    },
    "sNssais" : [ {
      "sd" : "fff000",
      "sst" : 0
    }, {
      "sd" : "fff111",
      "sst" : 1
    } ],
    "udmInfo" : {
      "externalGroupIdentifiersRanges" : [ {
        "pattern" : "^msisdn-52345678905\\d{4}$"
      } ],
      "gpsiRanges" : [ {
        "pattern" : "^msisdn-42345678904\\d{4}$"
      } ],
      "routingIndicators" : [ "4321" ],
      "supiRanges" : [ {
        "pattern" : "^imsi-60000\\d{4}$"
      } ]
    }
  }
`)

//BenchmarkCached
// func BenchmarkCached(b *testing.B) {
// }

//BenchmarkCachedWithTTL
func BenchmarkCachedWithTTL(b *testing.B) {
	StubSetupBenchmark()
	Instance().Flush("AUSF")

	for n := 0; n < b.N; n++ {
		Instance().CachedWithTTL("AUSF", "UDM", searchResultNfProfileUDM, 86400, false)
	}
}

func TestSearchInCache(t *testing.T) {
	StubSetupBenchmark()
	Instance().Flush("AUSF")

	for i := 0; i < 610; i++ {
		s := bytes.Replace(searchResultNfProfileUDM,
			[]byte("12345678-9udm-def0-1000-0000"), []byte("12345678-9udm-def0-1000-"+strconv.Itoa(i)), -1)
		Instance().CachedWithTTL("AUSF", "UDM", s, 86400, false)
	}
	// log.SetLevel(log.DebugLevel)

	sp := newSearchParam("12345678-9udm-def0-1000-x")
	content, ok := Instance().Search("AUSF", "UDM", sp, false)
	t.Log(string(content), ok)

	sp = newSearchParam("12345678-9udm-def0-1000-310")
	content, ok = Instance().Search("AUSF", "UDM", sp, false)
	t.Log(string(content), ok)
}

func benchSearchInCache(id string, num int, b *testing.B) {
	StubSetupBenchmark()
	Instance().Flush("AUSF")

	for i := 0; i < num; i++ {
		s := bytes.Replace(searchResultNfProfileUDM,
			[]byte("12345678-9udm-def0-1000-0000"), []byte("12345678-9udm-def0-1000-"+strconv.Itoa(i)), -1)
		Instance().CachedWithTTL("AUSF", "UDM", s, 86400, false)
	}
	// log.SetLevel(log.DebugLevel)

	sp := newSearchParam(id)
	for n := 0; n < b.N; n++ {
		Instance().Search("AUSF", "UDM", sp, false)
	}
}

//BenchmarkSearchInCache
func BenchmarkSearchInCacheHit(b *testing.B) {
	benchSearchInCache("12345678-9udm-def0-1000-0000", 1, b)
}

func BenchmarkSearchInCacheMiss(b *testing.B) {
	benchSearchInCache("12345678-9udm-def0-1000-x", 1, b)
}

func BenchmarkSearchInCacheUDM610Hit(b *testing.B) {
	benchSearchInCache("12345678-9udm-def0-1000-310", 610, b)
}

func BenchmarkSearchInCacheUDM610Miss(b *testing.B) {
	benchSearchInCache("12345678-9udm-def0-1000-x", 610, b)
}
