package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"testing"
	"net/http"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

func TestCalcProfilesEtag(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	filter := &NFDiscFilter{DiscNFPreFilterAction:&NFPreFilter{},
		DiscNFCommonFilter:&NFCommonFilter{},
		DiscNFServiceFiler:&NFServiceFilter{},
		DiscNFPostFilterAction:&NFPostFilter{}, }
	filter.Init(&DiscPara)
	data := (`[{
  "nfInstanceId": "12345678-abcd-ef12-1000-000000000002",
  "nfType": "AMF",
  "plmn": {
    "mcc": "240",
    "mnc": "240"
  },
  "sNssais": [
        {
      "sst": 2,
      "sd": "2"
    },
        {
      "sst": 4,
      "sd": "4"
    }
  ],
  "fqdn": "seliius03696.seli.gic.ericsson.se",
  "capacity": 100,
  "nfServices": [{
      "serviceInstanceId": "namf-comm-01",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "ipEndPoints": [
        {
          "ipv4Address": "172.16.208.1",
          "port": 30088
        }
      ],
      "capacity": 100 ,
      "supportedFeatures":"F000"
    ,"Priority":100},{
      "serviceInstanceId": "namf-comm-02",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "supportedFeatures":"0f00"
    ,"Priority":100}]
,"Priority":100}
,{
  "nfInstanceId": "12345678-abcd-ef12-1000-000000000001",
  "nfType": "AMF",
  "plmn": {
    "mcc": "240",
    "mnc": "240"
  },
  "sNssais": [
    {
      "sst": 1,
      "sd": "1"
    },
        {
      "sst": 2,
      "sd": "2"
    },
        {
      "sst": 3,
      "sd": "3"
    }
  ],
  "fqdn": "seliius03696.seli.gic.ericsson.se",
  "capacity": 100,
  "nfServices": [{
      "serviceInstanceId": "namf-comm-01",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "ipEndPoints": [
        {
          "ipv4Address": "172.16.208.1",
          "port": 30088
        }
      ],
      "capacity": 100 ,
      "supportedFeatures":"F000"
    ,"Priority":100},{
      "serviceInstanceId": "namf-comm-02",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "supportedFeatures":"0f00"
    ,"Priority":100}]
,"Priority":100}
]`)
	//test mutli-profiles with mutli-services
	filter.filterInfo.newProfiles=data
	filter.calcProfilesEtag()
	if "" == filter.GetFilterInfoEtag() {
		t.Fatalf("Etag should not empty, but failed")
	}
	etag := filter.GetFilterInfoEtag()
	filter.calcProfilesEtag()
	if etag != filter.GetFilterInfoEtag() {
		t.Fatalf("the same profile, every time calc etag by md5, the etag should be equal, but not equal")
	}

	data2 := (`[{
  "nfInstanceId": "12345678-abcd-ef12-1000-000000000002",
  "nfType": "AMF",
  "plmn": {
    "mcc": "240",
    "mnc": "240"
  },
  "sNssais": [
        {
      "sst": 2,
      "sd": "2"
    },
        {
      "sst": 4,
      "sd": "4"
    }
  ],
  "fqdn": "seliius03696.seli.gic.ericsson.se",
  "capacity": 100,
  "nfServices": [{
      "serviceInstanceId": "namf-comm-01",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "ipEndPoints": [
        {
          "ipv4Address": "172.16.208.1",
          "port": 30088
        }
      ],
      "capacity": 100 ,
      "supportedFeatures":"F000"
    ,"Priority":100},{
      "serviceInstanceId": "namf-comm-02",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "supportedFeatures":"0f00"
    ,"Priority":100}]
,"Priority":100}]`)
	//test one-profile with multi servces
	filter.filterInfo.newProfiles=data2
	filter.calcProfilesEtag()
	etag = filter.GetFilterInfoEtag()
	if "" == etag {
		t.Fatalf("Etag should not empty, but failed")
	}
	filter.calcProfilesEtag()
	if etag != filter.GetFilterInfoEtag() {
		t.Fatalf("the same profile, every time calc etag by md5, the etag should be equal, but not equal")
	}

	data3 := (`[{
  "nfInstanceId": "12345678-abcd-ef12-1000-000000000002",
  "nfType": "AMF",
  "plmn": {
    "mcc": "240",
    "mnc": "240"
  },
  "sNssais": [
        {
      "sst": 2,
      "sd": "2"
    },
        {
      "sst": 4,
      "sd": "4"
    }
  ],
  "fqdn": "seliius03696.seli.gic.ericsson.se",
  "capacity": 100,
  "nfServices": [{
      "serviceInstanceId": "namf-comm-01",
      "serviceName": "namf-comm",
      "version": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "schema": "https",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
      "ipEndPoints": [
        {
          "ipv4Address": "172.16.208.1",
          "port": 30088
        }
      ],
      "capacity": 100 ,
      "supportedFeatures":"F000"
    ,"Priority":100}
    ]
,"Priority":100}]`)
	//test one profile with one service
	filter.filterInfo.newProfiles= data3
	filter.calcProfilesEtag()
	etag = filter.GetFilterInfoEtag()
	if "" == etag {
		t.Fatalf("Etag should not empty, but failed")
	}
	filter.calcProfilesEtag()
	if etag != filter.GetFilterInfoEtag() {
		t.Fatalf("the same profile, every time calc etag by md5, the etag should be equal, but not equal")
	}
}

func TestGeneratorErrorInfo(t *testing.T) {
	var DiscPara nfdiscrequest.DiscGetPara
	filter := &NFDiscFilter{
		filterInfo:&FilterInfo{logcontent:&log.LogStruct{}}}
	filter.Init(&DiscPara)
	filter.generatorErrorInfo(false)
	if filter.filterInfo.statusCode != http.StatusNotFound {
		t.Fatal("statuscode should be 404, but not")
	}
	if filter.filterInfo.errorInfo != "requested NF profile not found" {
		t.Fatal("errorInfo should match, but not")
	}

	filter2 := &NFDiscFilter{
		filterInfo:&FilterInfo{logcontent:&log.LogStruct{}}}
	filter2.Init(&DiscPara)
	filter2.filterInfo.nfTypeForbiddenInProfile = true
	filter2.filterInfo.plmnForbiddenInProfile = true
	filter2.generatorErrorInfo(false)
	if filter2.filterInfo.statusCode != http.StatusForbidden {
		t.Fatal("statuscode should be 403, but not")
	}
	if filter2.filterInfo.errorInfo != "not allowed requester-plmn in nfprofile or requester-nf-type in nfprofile" {
		t.Fatal("errorInfo should match, but not")
	}
}

func TestParseResultForCache(t *testing.T)  {
	var nfResponse []string
	nfResponse = append(nfResponse, "struct(nfInstanceId:12345678-9abc-def0-1000-100000000021,profileUpdateTime:1)")
	nfResponse = append(nfResponse, "struct(nfInstanceId:12345678-9abc-def0-1000-100000000022,profileUpdateTime:2)")
	cacheResponse := parseResultForCache(nfResponse)
	//fmt.Printf("%v\n", cacheResponse)
	for index, value := range cacheResponse {
		if index==0{
			if value.ProfileUpdateTime !=1 || value.NfInstanceID!="12345678-9abc-def0-1000-100000000021" {
				t.Fatal("response should match, but not")
			}
		}else {
			if value.ProfileUpdateTime !=2 || value.NfInstanceID!="12345678-9abc-def0-1000-100000000022" {

				t.Fatal("response should match, but not")
			}
		}
	}
}