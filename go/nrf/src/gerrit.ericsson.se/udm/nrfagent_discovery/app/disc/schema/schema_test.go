package schema

import (
	"fmt"
	//"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	nfSchemaSuffix    = "src/gerrit.ericsson.se/udm/nrfagent_common/helm/eric-nrfagent-common/config/schema/nfProfileInSearchResult.json"
	patchSchemaSuffix = "src/gerrit.ericsson.se/udm/nrfagent_common/helm/eric-nrfagent-common/config/schema/patchDocument.json"

	nfSchema    string
	patchSchema string

	nfInstance = `{
  "nfInstanceId": "5g-ausf-01",
  "nfType": "AUSF",
  "nfStatus": "REGISTERED",
  "plmnList":[
    {
	   "mcc": "460",
	   "mnc": "000"
    }
  ],
  "sNssais": [
    {
      "sst": 0,
      "sd": "abAB01"
    }
  ],
  "fqdn": "seliius03696.seli.gic.ericsson.se",
  "ipv4Addresses": [
    "172.16.208.1"
  ],
  "ipv6Addresses": [
    "1001:da8::36"
	],
  "capacity": 100,
  "load" : 100,
  "nfServices": [
    {
      "serviceInstanceId": "nausf-auth-01",
      "nfServiceStatus": "REGISTERED",
      "serviceName": "nausf-auth",
      "versions": [{
        "apiVersionInUri":"v1",
        "apiFullVersion": "1.R15.1.1 " ,
        "expiry":"2020-07-06T02:54:32Z"}],
      "scheme": "http",
      "fqdn": "seliius03696.seli.gic.ericsson.se",
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
      ]
    }
  ]
}`

	nfInstanceBad = `{
  "nfInstanceId": "5g-ausf-01",
  "nfType": "AUSF",
  "nfStatus": "REGISTERED",
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
  "plmn": {
    "mcc": "460",
    "mnc": "000"
  },
  "nfServices": [
    {
      "serviceInstanceId": "nausf-auth-01",
      "nfServiceStatus": "REGISTERED",
      "serviceName": "nausf-auth",
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
      "allowedNfDomains":["172.17.0.20:3000/nnrf-nfm/v1/nf-instances"],
      "supportedFeatures":"A0A0",
      "capacity": 100,
      "load" : 100
    }
  ]
}`

	validPatchBody = `[
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
    ]`

	inValidPatchBody = `[
        {
            "op": "replace",
            "path": "/ipAddress",
            "value": ["10.0.0.3"]
        },
        {
            "oop": "add",
            "path": "/nfServiceList/0/allowedNfTypes",
            "value": ["nrf", "amf", "ausf"]
        }
    ]`
)

func init() {
	goPath := os.Getenv("GOPATH")
	fmt.Println("goPath", goPath)
	nfSchema = goPath + "/" + nfSchemaSuffix
	patchSchema = goPath + "/" + patchSchemaSuffix
	//	fmt.Println("nfSchema: ", nfSchema)
}

func TestValidateNfProfile(t *testing.T) {

	nfSchemaContent, err := ioutil.ReadFile(nfSchema)
	if err != nil {
		t.Fatalf("Load nfschema file failure, err:%v\n", err)
	}

	patchSchemaContent, err := ioutil.ReadFile(patchSchema)
	if err != nil {
		t.Fatalf("Load patchschema file failure, err:%v\n", err)
	}

	mapSchemaFile := make(map[string][]byte)
	mapSchemaFile["nfProfileInSearchResult.json"] = nfSchemaContent
	mapSchemaFile["patchDocument.json"] = patchSchemaContent

	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	for key := range mapSchemaFile {
		AbsFileName := currDir + "/" + key
		err := ioutil.WriteFile(AbsFileName, mapSchemaFile[key], 0666)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	os.Setenv("SCHEMA_DIR", currDir)
	os.Setenv("SCHEMA_NF_PROFILE", "nfProfileInSearchResult.json")
	os.Setenv("SCHEMA_PATCH_DOCUMENT", "patchDocument.json")

	errInfo := LoadDiscoverSchema()
	//fmt.Println("0-->errInfo:", errInfo)

	if errInfo != nil {
		t.Fatalf("LoadDiscoverSchema error,error info is %+v", err.Error())
	}

	errInfo = ValidateNfProfile(nfInstance)
	//fmt.Println("1-->errInfo:", errInfo)
	if errInfo != nil {
		t.Fatalf("This is a valid nf profile, but validate failed !")
	}

	errInfo = ValidateNfProfile(nfInstanceBad)
	//fmt.Println("2-->errInfo:", errInfo)
	if errInfo == nil {
		t.Fatalf("This is a invalid nf profile, but validate ok !")
	}

	errInfo = ValidatePatchDocument(validPatchBody)
	//fmt.Println("1-->errInfo:", errInfo)
	if errInfo != nil {
		t.Fatalf("This is a valid patch body, but validate failed !")
	}

	errInfo = ValidatePatchDocument(inValidPatchBody)
	//fmt.Println("2-->errInfo:", errInfo)
	if errInfo == nil {
		t.Fatalf("This is a invalid patch body, but validate ok !")
	}

}

func TestSetSchemaNfProfile(t *testing.T) {
	ok := SetSchemaNfProfile(nil)
	if ok == true {
		t.Error("TestSetSchemaNfProfile failed")
	}
}
