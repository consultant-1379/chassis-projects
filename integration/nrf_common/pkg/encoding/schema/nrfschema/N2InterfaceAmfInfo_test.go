package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsN2InterfaceAmfInfoValid(t *testing.T) {
	//right N2InterfaceAmfInfo
	body := []byte(`{
		"ipv4EndpointAddress": [
		    "10.10.10.10",
			"10.10.10.11"
		]
	}`)

	n2InterfaceAmfInfo := &TN2InterfaceAmfInfo{}
	err := json.Unmarshal(body, n2InterfaceAmfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !n2InterfaceAmfInfo.IsValid() {
		t.Fatalf("TN2InterfaceAmfInfo.IsValid didn't return right value!")
	}

	//right N2InterfaceAmfInfo
	body = []byte(`{
		"ipv6EndpointAddress": [
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		]
	}`)

	n2InterfaceAmfInfo = &TN2InterfaceAmfInfo{}
	err = json.Unmarshal(body, n2InterfaceAmfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !n2InterfaceAmfInfo.IsValid() {
		t.Fatalf("TN2InterfaceAmfInfo.IsValid didn't return right value!")
	}

	//right N2InterfaceAmfInfo
	body = []byte(`{
		"ipv4EndpointAddress": [
		    "10.10.10.10",
			"10.10.10.11"
		],
		"ipv6EndpointAddress": [
		    "1030::C9B4:FF12:48AA:1A2B",
			"1030::C9B4:FF12:48AA:1A2B"
		]
	}`)

	n2InterfaceAmfInfo = &TN2InterfaceAmfInfo{}
	err = json.Unmarshal(body, n2InterfaceAmfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !n2InterfaceAmfInfo.IsValid() {
		t.Fatalf("TN2InterfaceAmfInfo.IsValid didn't return right value!")
	}

	//wrong N2InterfaceAmfInfo
	body = []byte(`{
	}`)

	n2InterfaceAmfInfo = &TN2InterfaceAmfInfo{}
	err = json.Unmarshal(body, n2InterfaceAmfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if n2InterfaceAmfInfo.IsValid() {
		t.Fatalf("TN2InterfaceAmfInfo.IsValid didn't return right value!")
	}
}
