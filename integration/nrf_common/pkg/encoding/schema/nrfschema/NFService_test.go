package nrfschema

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestGetInvalidIPEndPointIndexs(t *testing.T) {
	//NFService without ipEndPoints
	body := []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED"
	}`)

	nfService := &TNFService{}
	err := json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidIPEndPointIndexs() != nil {
		t.Fatalf("TNFService.GetInvalidIPEndPointIndexs didn't return value!")
	}

	//NFService with right ipEndPoints
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED",
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
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidIPEndPointIndexs() != nil {
		t.Fatalf("TNFService.GetInvalidIPEndPointIndexs didn't return value!")
	}

	//NFService with wrong ipEndPoints
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED",
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
			},
			{
				"port": 80,
				"ipv4Address": "10.10.10.10",
				"ipv6Address": "1030::C9B4:FF12:48AA:1A2B"
			}
		]
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidIPEndPointIndexs() == nil {
		t.Fatalf("TNFService.GetInvalidIPEndPointIndexs didn't return value!")
	}
}

func TestGetInvalidChfServiceInfoIndex(t *testing.T) {
	//NFService without chfServiceInfo is valid
	body := []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED"
	}`)

	nfService := &TNFService{}
	err := json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidChfServiceInfoIndex() != "" {
		t.Fatalf("TNFService.GetInvalidChfServiceInfoIndex didn't return value!")
	}

	//NFService with right chfServiceInfo is valid
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED",
		"chfServiceInfo": {
		    "primaryChfServiceInstance": "serv01"
		}
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidChfServiceInfoIndex() != "" {
		t.Fatalf("TNFService.GetInvalidChfServiceInfoIndex didn't return value!")
	}

	//NFService with invalid chfServiceInfo is invalid
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "srv1",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED",
		"chfServiceInfo": {
		    "primaryChfServiceInstance": "serv01",
			"secondaryChfServiceInstance": "serv02"
		}
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidChfServiceInfoIndex() != constvalue.NFServiceChfServiceInfo {
		t.Fatalf("TNFService.GetInvalidChfServiceInfoIndex didn't return value!")
	}
}

func TestGetInvalidServiceNameIndex(t *testing.T) {
	//NFService with custom serviceName is valid
	body := []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "namf-xxxx",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED"
	}`)

	nfService := &TNFService{}
	err := json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidServiceNameIndex(constvalue.NfTypeAMF) != "" {
		t.Fatalf("TNFService.GetInvalidServiceNameIndex didn't return value!")
	}

	//NFService with matched serviceName is valid
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "namf-comm",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED"
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidServiceNameIndex(constvalue.NfTypeAMF) != "" {
		t.Fatalf("TNFService.GetInvalidServiceNameIndex didn't return value!")
	}

	//NFService with unmatched serviceName is invalid
	body = []byte(`{
        "serviceInstanceId": "srvId01",
        "serviceName": "nausf-auth",
        "versions": [],
        "scheme": "http",
		"nfServiceStatus": "SUSPENDED"
	}`)

	nfService = &TNFService{}
	err = json.Unmarshal(body, nfService)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService.GetInvalidServiceNameIndex(constvalue.NfTypeAMF) != constvalue.NFServiceName {
		t.Fatalf("TNFService.GetInvalidServiceNameIndex didn't return value!")
	}
}

func TestGenerateMd5ForNFService(t *testing.T) {
	//two NFServices between which only the attribute(interPlmnFqdn or allowedPlmns
	//or allowedNfTypes or allowedNfDomains or allowedNssais) is different, their md5 shall be the same
	body1 := []byte(`{
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
				"interPlmnFqdn": "http://test1",
				"allowedPlmns": [
				    {
						"mcc": "460",
						"mnc": "01"
					}
				],
				"allowedNfTypes": ["AUSF"],
				"allowedNfDomains": ["domain01"],	
			    "allowedNssais": [
				    {
						"sst": 100,
						"sd": "111111"
					}
				]	
			}`)

	body2 := []byte(`{
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
				"interPlmnFqdn": "http://test2",
				"allowedPlmns": [
				    {
						"mcc": "460",
						"mnc": "02"
					}
				],
				"allowedNfTypes": ["AUSF"],
				"allowedNfDomains": ["domain02"],	
			    "allowedNssais": [
				    {
						"sst": 200,
						"sd": "22222"
					}
				]				
			}`)

	nfService1 := &TNFService{}
	err := json.Unmarshal(body1, nfService1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	nfService2 := &TNFService{}
	err = json.Unmarshal(body2, nfService2)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService1.GenerateMd5() != nfService2.GenerateMd5() {
		t.Fatalf("TNFService.GenerateMd5 didn't return right value!")
	}

	//two NFServices between which attribute except for (interPlmnFqdn or allowedPlmns
	//or allowedNfTypes or allowedNfDomains or allowedNssais) is different, their md5 shall be different
	body1 = []byte(`{
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
				"interPlmnFqdn": "http://test1",
				"allowedPlmns": [
				    {
						"mcc": "460",
						"mnc": "01"
					}
				],
				"allowedNfTypes": ["AUSF"],
				"allowedNfDomains": ["domain01"],	
			    "allowedNssais": [
				    {
						"sst": 100,
						"sd": "111111"
					}
				]	
			}`)

	body2 = []byte(`{
				"serviceInstanceId": "namf-01",
				"serviceName": "namf-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "https",
				"nfServiceStatus": "REGISTERED",
				"interPlmnFqdn": "http://test2",
				"allowedPlmns": [
				    {
						"mcc": "460",
						"mnc": "02"
					}
				],
				"allowedNfTypes": ["AUSF"],
				"allowedNfDomains": ["domain02"],	
			    "allowedNssais": [
				    {
						"sst": 200,
						"sd": "22222"
					}
				]				
			}`)

	nfService1 = &TNFService{}
	err = json.Unmarshal(body1, nfService1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	nfService2 = &TNFService{}
	err = json.Unmarshal(body2, nfService2)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if nfService1.GenerateMd5() == nfService2.GenerateMd5() {
		t.Fatalf("TNFService.GenerateMd5 didn't return right value!")
	}
}

func TestAllowedParametersInService(t *testing.T) {
	body1 := []byte(`{
				"serviceInstanceId": "nudm-01",
				"serviceName": "nudm-01",
				"versions": [
				    {
						"apiVersionInUri": "http://test",
						"apiFullVersion": "0.1"
					}
				],
				"scheme": "http",
				"nfServiceStatus": "REGISTERED",
				"interPlmnFqdn": "http://test1",
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
				]	
			}`)

	nfService1 := &TNFService{}
	err := json.Unmarshal(body1, nfService1)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !nfService1.IsAllowedNfType("AUSF") || !nfService1.IsAllowedNfType("AMF") {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfService1.IsAllowedNfType("UDR") {
		t.Fatalf("Should not allow, but YES!")
	}

	if !nfService1.IsAllowedNfDomain("seliius12345.seli.gic.ericsson.se") || !nfService1.IsAllowedNfDomain("seliius1234.seli.gic.ericsson.se") {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfService1.IsAllowedNfDomain("seliius123.seli.gic.ericsson.se") {
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

	if !nfService1.IsAllowedPlmn(plmnID1) || !nfService1.IsAllowedPlmn(plmnID2) {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfService1.IsAllowedPlmn(plmnID3) {
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

	if !nfService1.IsAllowedNssai(snssai1) || !nfService1.IsAllowedNssai(snssai3) {
		t.Fatalf("Should allow, but NOT!")
	}

	if nfService1.IsAllowedNssai(snssai2) {
		t.Fatalf("Should not allow, but YES!")
	}
}
