package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsInterfaceUpfInfoItemValid(t *testing.T) {
	//right InterfaceUpfInfoItem
	body := []byte(`{
		"networkInstance": "network1",
		"ipv4EndpointAddresses": [
		    "10.10.10.10",
			"10.10.10.11"
		]
	}`)

	interfaceUpfInfoItem := &TInterfaceUpfInfoItem{}
	err := json.Unmarshal(body, interfaceUpfInfoItem)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !interfaceUpfInfoItem.IsValid() {
		t.Fatalf("TInterfaceUpfInfoItem.IsValid didn't return right value!")
	}

	//right InterfaceUpfInfoItem
	body = []byte(`{
		"networkInstance": "network1",
		"ipv6EndpointAddresses": [
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		]
	}`)

	interfaceUpfInfoItem = &TInterfaceUpfInfoItem{}
	err = json.Unmarshal(body, interfaceUpfInfoItem)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !interfaceUpfInfoItem.IsValid() {
		t.Fatalf("TInterfaceUpfInfoItem.IsValid didn't return right value!")
	}

	//right InterfaceUpfInfoItem
	body = []byte(`{
		"networkInstance": "network1",
        "endpointFqdn": "http://test"
	}`)

	interfaceUpfInfoItem = &TInterfaceUpfInfoItem{}
	err = json.Unmarshal(body, interfaceUpfInfoItem)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !interfaceUpfInfoItem.IsValid() {
		t.Fatalf("TInterfaceUpfInfoItem.IsValid didn't return right value!")
	}

	//right InterfaceUpfInfoItem
	body = []byte(`{
		"networkInstance": "network1",
        "endpointFqdn": "http://test",
		"ipv4EndpointAddresses": [
		    "10.10.10.10",
			"10.10.10.11"
		],
		"ipv6EndpointAddresses": [
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		]
	}`)

	interfaceUpfInfoItem = &TInterfaceUpfInfoItem{}
	err = json.Unmarshal(body, interfaceUpfInfoItem)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !interfaceUpfInfoItem.IsValid() {
		t.Fatalf("TInterfaceUpfInfoItem.IsValid didn't return right value!")
	}

	//wrong InterfaceUpfInfoItem
	body = []byte(`{
		"networkInstance": "network1"
	}`)

	interfaceUpfInfoItem = &TInterfaceUpfInfoItem{}
	err = json.Unmarshal(body, interfaceUpfInfoItem)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if interfaceUpfInfoItem.IsValid() {
		t.Fatalf("TInterfaceUpfInfoItem.IsValid didn't return right value!")
	}
}
