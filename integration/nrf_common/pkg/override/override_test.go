package override

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
	"gerrit.ericsson.se/udm/nrf_common/pkg/schema"
	"github.com/buger/jsonparser"
)

var patchBody1 = []byte(`[
	{
			"op": "replace",
			"path": "/nfServices/0/priority",
			"value": 15
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/nfServices/0/priority"
	}
]`)

var patchBody2 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/priority",
			"value": 50
	}
]`)

var patchBody3 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/allowedPlmns",
			"value": [{ "mcc": "460",  "mnc": "00" }]
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/nfServices/0/allowedPlmns"
	}
]`)

var patchBody4 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/allowedPlmns",
			"value":[{ "mcc": "450",  "mnc": "30" }]
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/capacity"
	}
]`)

var patchBody5 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/allowedNfTypes",
			"value":["AUSF","UDM","AMF"]
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/nfServices/0/allowedNfTypes"
	}
]`)

var patchBody6 = []byte(`[
	{
			"op": "remove",
			"path": "/provisionInfo/overrideAttrList/0"
	}
]`)

var patchBody7 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/allowedNssais",
			"value":[{ "sst": 123,  "sd": "abcdef" },{  "sst": 456,  "sd": "a12345" }]
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/nfServices/0/allowedNssais"
	}
]`)
var patchBody8 = []byte(`[
	{
			"op": "replace",
			"path": "/nfServices/0/priority",
			"value": 15
	},
    {
			"op": "replace",
			"path": "capacity",
			"value": 15
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "/nfServices/0/priority"
	},
	{
			"op": "add",
			"path": "/provisionInfo/overrideAttrList/-",
			"value": "capacity"
	}
]`)

//var patchBody9 = []byte(`[
//	{
//			"op": "add",
//			"path": "/nfServices/0/allowedNfTypes/-",
//			"value": "SMF"
//	},
//    {
//			"op": "replace",
//			"path": "capacity",
//			"value": 15
//	},
//	{
//			"op": "add",
//			"path": "/provisionInfo/overrideAttrList/-",
//			"value": "/nfServices/0/allowedNfTypes"
//	},
//	{
//			"op": "add",
//			"path": "/provisionInfo/overrideAttrList/-",
//			"value": "capacity"
//	}
//]`)

var patchBody10 = []byte(`[
    {
        "op": "add",
        "path": "/nfServices/0/recoveryTime",
        "value": "2018-12-11T23:20:50Z"
    },
    {
        "op": "add",
        "path": "/provisionInfo/overrideAttrList/-",
        "value": "/nfServices/0/recoveryTime"
    }
]
`)

var patchBody11 = []byte(`[
    {
        "op": "replace",
        "path": "/priority",
        "value": 11
    },
    {
        "op": "replace",
        "path": "/nfService/0/capacity",
        "value": 80
    },     
    {
        "op": "add",
        "path": "/provisionInfo/overrideAttrList",
        "value": ["/priority","/nfService/0/capacity"]
    }
] 
`)

var nfPatchBody1 = []byte(`[
	{
			"op": "add",
			"path": "/nfServices/0/priority",
			"value": 2
	},
    {
			"op": "add",
			"path": "/nfServices/0/capacity",
			"value": 50
	}
]`)

var nfPatchRemoveOverride = []byte(`[
	{
			"op": "remove",
			"path": "/provisionInfo/overrideAttrList"
	}
]`)
var patchChangeMcc = []byte(`[
	{
			"op": "replace",
			"path": "/nfServices/0/allowedPlmns/0/mcc",
			"value": "460"
	}
]`)

var nfProfileBody = []byte(`{
	                     "nfInstanceId":"udm01",
						"nfType":"",
						    "plmn": {
        							"mcc": "460",
						        "mnc": "000"
        						},
						"sNssai": {"sst": "0","sd": "1"},
						"fqdn": "udm.mnc001.mcc460.5g",
						"ipAddress": ["10.0.0.2"],
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

var profileFull = []byte(`{
	"nfInstanceId": "5g-udm-01",
	"nfType": "UDM",
	"nfStatus": "REGISTERED",
	"plmn": {
	  "mcc": "460",
	  "mnc": "66"
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
	"capacity": 120,
	"priority": 12,
	"load" : 100,
	"udmInfo": {
	  "groupId": "gid01",
	   "gpsiRanges": [
	   {
		 "start": "12300000",
		 "end": "12399999"
	   }
	],
	"supiRanges": [
	   {
		 "start": "12300000",
		 "end": "12399999"
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
		"priority": 10,
		"load" : 100
	  }
	]
}`)

var ProfileWithoutPriority = []byte(`{
	"nfInstanceId": "5g-udm-01",
	"nfType": "UDM",
	"nfStatus": "REGISTERED",
	"plmn": {
	  "mcc": "460",
	  "mnc": "66"
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
	"capacity": 120,
	"load" : 100,
	"udmInfo": {
	  "groupId": "gid01",
	   "gpsiRanges": [
	   {
		 "start": "12300000",
		 "end": "12399999"
	   }
	],
	"supiRanges": [
	   {
		 "start": "12300000",
		 "end": "12399999"
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

func init() {
	log.SetLevel(log.FatalLevel)
	loadSchema()
}

func loadSchema() {
	goPath := os.Getenv("GOPATH")
	os.Setenv("SCHEMA_DIR", fmt.Sprintf("%s/src/gerrit.ericsson.se/udm/nrf_common/helm/eric-nrf-common/config/schema", goPath))
	os.Setenv("SCHEMA_NF_PROFILE", "nfProfile.json")
	os.Setenv("SCHEMA_PATCH_DOCUMENT", "patchDocument.json")
	os.Setenv("SCHEMA_SUBSCRIPTIONDATA", "subscriptionData.json")
	os.Setenv("SCHEMA_SUBSCRIPTIONPATCH", "subscriptionPatch.json")

	err := schema.LoadManagementSchema()
	if err != nil {
		log.Fatalf("LoadManagementSchema error, %v", err)
	}
}

func TestApplyOverrideAttributes(t *testing.T) {
	//////////2//////////////////
	var overrideInfo []nrfschema.OverrideInfo
	oneReplaceInfo1 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/capacity",
		Value:  "50",
	}

	oneAddInfo1 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/priority",
		Value:  "30",
	}
	overrideInfo = append(overrideInfo, oneReplaceInfo1)
	overrideInfo = append(overrideInfo, oneAddInfo1)

	updateProfileFull, updateOverrideInfo, err := applyOverrideAttributes(profileFull, overrideInfo)
	if err != nil {
		t.Fatalf("Expect apply overrideAttributes success, but failure, err:%s", err.Error())
	}
	if len(updateOverrideInfo) != 2 {
		t.Fatalf("Expect new overrideAttributes is 2, but no")
	}
	_, err = jsonparser.ArrayEach(updateProfileFull, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		capacity, err := jsonparser.GetInt(value, "capacity")
		if err != nil {
			t.Fatalf("Get capacity from nfService failure, err:%s", err.Error())
		}
		if capacity != 50 {
			t.Fatalf("Expect capacity is 50, but no")
		}

		priority, err := jsonparser.GetInt(value, "priority")
		if err != nil {
			t.Fatalf("Get priority from nfService failure, err:%s", err.Error())
		}
		if priority != 30 {
			t.Fatalf("Expect priority is 30, but no")
		}
	}, "nfServices")

	//////////2/////////////
	var overrideInfo2 []nrfschema.OverrideInfo
	oneReplaceInfo2 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/capacity",
		Value:  "60",
	}
	oneAddInfo2 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/priority",
		Value:  "50",
	}
	overrideInfo2 = append(overrideInfo2, oneReplaceInfo2)
	overrideInfo2 = append(overrideInfo2, oneAddInfo2)

	updateProfileFull2, updateOverrideInfo2, err := applyOverrideAttributes(ProfileWithoutPriority, overrideInfo2)
	if err != nil {
		t.Fatalf("Expect apply overrideAttributes success, but failure")
	}
	if len(updateOverrideInfo2) != 2 {
		t.Fatalf("Expect new overrideAttributes is 2, but no")
	}
	_, err = jsonparser.ArrayEach(updateProfileFull2, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		capacity, err := jsonparser.GetInt(value, "capacity")
		if err != nil {
			t.Fatalf("Get capacity from nfService failure, err:%s", err.Error())
		}
		if capacity != 60 {
			t.Fatalf("Expect capacity is 50, but no")
		}

		priority, err := jsonparser.GetInt(value, "priority")
		if err != nil {
			t.Fatalf("Get priority from nfService failure, err:%s", err.Error())
		}
		if priority != 50 {
			t.Fatalf("Expect priority is 30, but no")
		}
	}, "nfServices")

	//////////3//////////
	var overrideInfo3 []nrfschema.OverrideInfo
	oneReplaceInfo3 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/capacity",
		Value:  "80",
	}
	oneAddInfo3 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/1/priority",
		Value:  "60",
	}
	overrideInfo3 = append(overrideInfo3, oneReplaceInfo3)
	overrideInfo3 = append(overrideInfo3, oneAddInfo3)

	updateProfileFull3, updateOverrideInfo3, err := applyOverrideAttributes(ProfileWithoutPriority, overrideInfo3)
	if err != nil {
		t.Fatalf("Expect apply overrideAttributes success, but failure")
	}
	if len(updateOverrideInfo3) != 1 {
		t.Fatalf("Expect new overrideAttributes is 1, but no")
	}
	_, err = jsonparser.ArrayEach(updateProfileFull3, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		capacity, err := jsonparser.GetInt(value, "capacity")
		if err != nil {
			t.Fatalf("Get capacity from nfService failure, err:%s", err.Error())
		}
		if capacity != 80 {
			t.Fatalf("Expect capacity is 80, but no")
		}

		_, err = jsonparser.GetInt(value, "priority")
		if err == nil {
			t.Fatalf("Expect without priority from nfService, but not")
		}
	}, "nfServices")

}

func TestConvertPatchItem(t *testing.T) {
	oneInfo1 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/capacity",
		Value:  "50",
	}

	patchContent := getPatchContent(oneInfo1)
	if len(patchContent) == 0 {
		t.Fatalf("Expect convert success, but failure")
	}

	oneInfo2 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/capacity",
		Value:  "no-num",
	}

	patchContent = getPatchContent(oneInfo2)
	if len(patchContent) == 0 {
		t.Fatalf("Expect convert failure, but success")
	}
	oneInfo3 := nrfschema.OverrideInfo{
		Action: "replace",
		Path:   "/nfServices/0/allowedPlmns",
		Value:  "[{\"mcc\":\"450\",\"mnc\":\"30\"}]",
	}

	patchContent = getPatchContent(oneInfo3)
	if len(patchContent) == 0 {
		t.Fatalf("Expect convert failure, but success")
	}
}

func TestConstructOverrideInfo(t *testing.T) {
	//test add override info /nfServices/0/priority
	patchData1 := make([]nrfschema.TPatchItem, 0)
	overrideList := make([]nrfschema.OverrideInfo, 0)
	err := json.Unmarshal(patchBody1, &patchData1)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ := constructOverrideInfo(true, patchData1, overrideList)

	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 1 || overrideList[0].Value != "15" || len(overrideAttrList) != 1 {
		t.Fatalf("TestConstructOverrideInfo: patchData1 check result failure.")
	}

	//test change already override item /nfServices/0/priority value
	patchData2 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody2, &patchData2)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData2, overrideList)
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 1 || overrideList[0].Value != "50" || len(overrideAttrList) != 1 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}

	//test add override info /nfServices/0/allowedPlmns
	patchData3 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody3, &patchData3)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData3, overrideList)
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 2 || overrideList[0].Path != "/nfServices/0/priority" || overrideList[1].Path != "/nfServices/0/allowedPlmns" || len(overrideAttrList) != 2 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}

	//test change already override item /nfServices/0/allowedPlmns value, and ignore overrideAttrList /capacity  without patch value
	patchData4 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody4, &patchData4)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, _, _ = constructOverrideInfo(true, patchData4, overrideList)
	//t.Logf("overrideStr: %s, %d", overrideStr, len(overrideStr))
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 2 || overrideList[0].Path != "/nfServices/0/priority" || overrideList[1].Value != "[{\"mcc\":\"450\",\"mnc\":\"30\"}]" || len(overrideAttrList) != 2 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}

	//test add override item /nfServices/0/allowedNfTypes value
	patchData5 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody5, &patchData5)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData5, overrideList)
	//t.Logf("overrideStr: %s, %d", overrideStr, len(overrideStr))
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 3 || overrideList[0].Path != "/nfServices/0/priority" || overrideList[2].Value != "[\"AUSF\",\"UDM\",\"AMF\"]" || len(overrideAttrList) != 3 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
	//test remove override item allowedNfTypes
	patchData6 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody6, &patchData6)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData6, overrideList)
	//t.Logf("overrideStr: %s, %d", overrideStr, len(overrideStr))
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	if err != nil || len(overrideList) != 2 || overrideList[0].Path != "/nfServices/0/allowedPlmns" || overrideList[1].Value != "[\"AUSF\",\"UDM\",\"AMF\"]" || len(overrideAttrList) != 2 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
	patchData7 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody7, &patchData7)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData7, overrideList)
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	//t.Logf("overrideStr: %s, %v", overrideStr, overrideList)
	if err != nil || len(overrideList) != 3 || overrideList[0].Path != "/nfServices/0/allowedPlmns" || overrideList[1].Value != "[\"AUSF\",\"UDM\",\"AMF\"]" || overrideList[2].Path != "/nfServices/0/allowedNssais" || len(overrideAttrList) != 3 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
	patchData10 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody10, &patchData10)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData10, overrideList)
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	//t.Logf("overrideStr: %s, %v", overrideStr, overrideList)
	if err != nil || len(overrideList) != 4 || overrideList[0].Path != "/nfServices/0/allowedPlmns" || overrideList[1].Value != "[\"AUSF\",\"UDM\",\"AMF\"]" || overrideList[3].Path != "/nfServices/0/recoveryTime" || len(overrideAttrList) != 4 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
	patchData11 := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(patchBody11, &patchData11)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, _ = constructOverrideInfo(true, patchData11, overrideList)
	err = json.Unmarshal([]byte(overrideStr), &overrideList)
	//t.Logf("overrideStr: %s", overrideStr)
	//t.Logf("overrideStr: %s, %v", overrideStr, overrideList)
	if err != nil || len(overrideList) != 6 || overrideList[4].Path != "/priority" || overrideList[5].Value != "80" || len(overrideAttrList) != 6 {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
	patchDataRemove := make([]nrfschema.TPatchItem, 0)
	err = json.Unmarshal(nfPatchRemoveOverride, &patchDataRemove)
	if err != nil {
		t.Fatalf("Unmarshal patch data failed!")
	}
	overrideStr, overrideAttrList, err = constructOverrideInfo(true, patchDataRemove, overrideList)
	//t.Logf("overrideStr: %s", overrideStr)
	//t.Logf("overrideStr: %s, %v", overrideStr, overrideList)
	if err != nil || overrideStr != "" {
		t.Fatalf("TestConstructOverrideInfo:  check result failure.")
	}
}

func TestUpdateProvOverrideInfo(t *testing.T) {
	nfProfile := &nrfschema.TNFProfile{}
	err := json.Unmarshal(nfProfileBody, nfProfile)
	if err != nil {
		t.Fatalf("TestUpdateProvOverrideInfo: Unmarshal failed, %s!", err.Error())
	}
	var overrideList = make([]string, 0)
	overrideList = append(overrideList, "/nfServices/0/priority")
	overrideList = append(overrideList, "capacity")

	UpdateProvOverrideInfo(nfProfile, constvalue.Cmode_Provisioned, overrideList)
	if nfProfile.ProvisionInfo.CreateMode != constvalue.CMODE_PROVISIONED || len(nfProfile.ProvisionInfo.OverrideAttrList) != 0 {
		t.Fatalf("TestUpdateProvOverrideInfo: patchData2 remove overrideList override item check failure.")
	}

	UpdateProvOverrideInfo(nfProfile, constvalue.Cmode_NFRegistered, overrideList)
	//t.Logf("TestUpdateProvOverrideInfo nfProfile:%v, overrideList: %v", nfProfile.ProvisionInfo, overrideList)
	if nfProfile.ProvisionInfo.CreateMode != constvalue.CMODE_NF_REGISTERED || nfProfile.ProvisionInfo.OverrideAttrList[0] != "/nfServices/0/priority" || nfProfile.ProvisionInfo.OverrideAttrList[1] != "capacity" {
		t.Fatalf("TestUpdateProvOverrideInfo: patchData2 remove overrideList override item check failure.")
	}

}

func TestAppendProvisionInfo(t *testing.T) {
	overrideInfo := []byte(`[{"path":"/nfServices/0/allowedPlmns","action":"replace","value":"{\"mcc\":\"450\",\"mnc\":\"30\"}"},{"path":"/nfServices/0/allowedNfTypes","action":"replace","value":"[\"AUSF\",\"UDM\",\"AMF\"]"},{"path":"/nfServices/0/allowedNssais","action":"replace","value":"[{\"sd\":\"abcdef\",\"sst\":123},{\"sd\":\"a12345\",\"sst\":456}]"}]`)
	var overrideExist bool = true
	var provFlag int64 = constvalue.Cmode_NFRegistered
	succ := AppendProvisionInfo(&nfProfileBody, overrideExist, &overrideInfo, provFlag)
	if !succ {
		t.Fatalf("TestAppendProvisionInfo: AppendProvisionInfo NF_REGISTERED value check failure")
	}
	provFlag = constvalue.Cmode_Provisioned
	succ = AppendProvisionInfo(&nfProfileBody, overrideExist, &overrideInfo, provFlag)
	//t.Logf("TestAppendProvisionInfo ret:%t, profile: %s", succ, string(nfProfileBody))
	if !succ {
		t.Fatalf("TestAppendProvisionInfo: AppendProvisionInfo PROVISIONED value check failure")
	}
}
