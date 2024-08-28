package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetInvalidInterfaceUpfInfoIndexs(t *testing.T) {
	//UpfInfo without interfaceUpfInfoList
	body := []byte(`{
	}`)

	upfInfo := &TUpfInfo{}
	err := json.Unmarshal(body, upfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if upfInfo.GetInvalidInterfaceUpfInfoIndexs() != nil {
		t.Fatalf("TUpfInfo.GetInvalidInterfaceUpfInfoIndexs didn't return value!")
	}

	//UpfInfo with right interfaceUpfInfoList
	body = []byte(`{
		"interfaceUpfInfoList": [
		    {
			    "interfaceType": "N2",
				"ipv4EndpointAddresses": [
		            "10.10.10.10",
			        "10.10.10.11"
		        ]
			},
			{
			    "interfaceType": "N2",
				"ipv6EndpointAddresses": [
		            "1030::C9B4:FF12:48AA:1A2B",
			        "1030::C9B4:FF12:48AA:1A2B"
		        ]	
			},
			{
				"interfaceType": "N2",
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
			    "interfaceType": "N2",
				"endpointFqdn": "http://test"
			}
		]
		
	}`)

	upfInfo = &TUpfInfo{}
	err = json.Unmarshal(body, upfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if upfInfo.GetInvalidInterfaceUpfInfoIndexs() != nil {
		t.Fatalf("TUpfInfo.GetInvalidInterfaceUpfInfoIndexs didn't return value!")
	}

	//UpfInfo with wrong interfaceUpfInfoList
	body = []byte(`{
		"interfaceUpfInfoList": [
		    {
			    "interfaceType": "N2",
				"ipv4EndpointAddresses": [
		            "10.10.10.10",
			        "10.10.10.11"
		        ]
			},
			{
			    "interfaceType": "N2",
				"ipv6EndpointAddresses": [
		            "1030::C9B4:FF12:48AA:1A2B",
			        "1030::C9B4:FF12:48AA:1A2B"
		        ]	
			},
			{
			    "interfaceType": "N2",
				"endpointFqdn": "http://test"
			},
			{
				"interfaceType": "N2",
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
			    "interfaceType": "N2"
			}
		]
		
	}`)

	upfInfo = &TUpfInfo{}
	err = json.Unmarshal(body, upfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if upfInfo.GetInvalidInterfaceUpfInfoIndexs() == nil {
		t.Fatalf("TUpfInfo.GetInvalidInterfaceUpfInfoIndexs didn't return value!")
	}
}

func TestCreateNfInfo(t *testing.T) {
	body := []byte(`{
		"sNssaiUpfInfoList": [
			{
				"sNssai": {
					"sst": 111,
					"sd": "123"
				},
				"dnnUpfInfoList": [
					{
						"dnn": "222"
					},
					{
						"dnn": "333"
					}
				]
			},
			{
				"sNssai": {
					"sst": 222,
					"sd": "222"
				},
				"dnnUpfInfoList": [
					{
						"dnn": "444"
					},
					{
						"dnn": "555",
						"dnaiList":["111","222"]
					}
				]
			}
		]
	}`)

	upfInfo := &TUpfInfo{}
	err := json.Unmarshal(body, upfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}
	upfInfoJson := upfInfo.createNfInfo()
	if upfInfoJson != `"upfInfo":{"sNssaiUpfInfoList":[{"sNssai":{"sst":111,"sd":"123"},"dnnUpfInfoList":[{"dnn":"222","dnaiList":["RESERVED_EMPTY_DNAI"]},{"dnn":"333","dnaiList":["RESERVED_EMPTY_DNAI"]}]},{"sNssai":{"sst":222,"sd":"222"},"dnnUpfInfoList":[{"dnn":"444","dnaiList":["RESERVED_EMPTY_DNAI"]},{"dnn":"555","dnaiList":["111","222"]}]}], "smfServingArea":["RESERVED_EMPTY_SMFSERVINGAREA"]}` {
		t.Fatal("upfInfo helper should matched, but fail")
	}
}
