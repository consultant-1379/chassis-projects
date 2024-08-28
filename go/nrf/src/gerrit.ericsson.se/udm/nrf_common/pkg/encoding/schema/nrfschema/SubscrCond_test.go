package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestSubscrCondValidate(t *testing.T) {
	// empty object is invalid
	body := []byte(`{}`)

	subscrCond := &TSubscrCond{}

	err := json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() == "" {
		t.Fatalf("it is a invalid SubscrCond, but validate pass")
	}

	// multiple conditions is invalid
	body = []byte(`{
		"nfInstanceId": "amf01",
		"nfType": "AMF"
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() == "" {
		t.Fatalf("it is a invalid SubscrCond, but validate pass")
	}

	// nsiList present not along with snssaiList is invalid
	body = []byte(`{
		"nsiList": ["nsi1", "nsi2"]
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() == "" {
		t.Fatalf("it is a invalid SubscrCond, but validate pass")
	}

	// nsiList present along with snssaiList is valid
	body = []byte(`{
		"snssaiList": [
		    {
			    "sst": 1,
			    "sd": "123456"
			}
		],
		"nsiList": ["nsi1", "nsi2"]
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() != "" {
		t.Fatalf("it is a valid SubscrCond, but validate fail")
	}

	// nfGroupId present not along with nfType is invalid
	body = []byte(`{
		"nfGroupId": "group01"
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() == "" {
		t.Fatalf("it is a invalid SubscrCond, but validate pass")
	}

	// nfGroupId present along with nfType not equal to AUSF, UDM or UDR is invalid
	body = []byte(`{
		"nfGroupId": "group01",
		"nfType": "AMF"
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() == "" {
		t.Fatalf("it is a invalid SubscrCond, but validate pass")
	}

	// nfGroupId present along with nfType equal to AUSF, UDM or UDR is valid
	body = []byte(`{
		"nfGroupId": "group01",
		"nfType": "AUSF"
	}`)

	subscrCond = &TSubscrCond{}

	err = json.Unmarshal(body, subscrCond)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if subscrCond.Validate() != "" {
		t.Fatalf("it is a valid SubscrCond, but validate fail")
	}
}
