package nrfschema

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestGenerateGrpcPutKeyForSnssai(t *testing.T) {
	//Snssai with sd
	body := []byte(`{
		"sst": 10,
		"sd": "123456"
	}`)

	snssai := &TSnssai{}
	err := json.Unmarshal(body, snssai)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	subKey := snssai.GenerateGrpcPutKey()
	if subKey.SubKey1 != "10" || subKey.SubKey2 != "123456" {
		t.Fatalf("TSnssai.GenerateGrpcPutKey didn't return right value!")
	}

	//Snssai without sd
	body = []byte(`{
		"sst": 10
	}`)

	snssai = &TSnssai{}
	err = json.Unmarshal(body, snssai)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	subKey = snssai.GenerateGrpcPutKey()
	if subKey.SubKey1 != "10" || subKey.SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSnssai.GenerateGrpcPutKey didn't return right value!")
	}
}

func TestGenerateGrpcGetKeyForSnssai(t *testing.T) {
	//Snssai with sd
	body := []byte(`{
		"sst": 10,
		"sd": "123456"
	}`)

	snssai := &TSnssai{}
	err := json.Unmarshal(body, snssai)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	subKeys := snssai.GenerateGrpcGetKey()
	if len(subKeys) != 2 {
		t.Fatalf("TSnssai.GenerateGrpcGetKey didn't return right value!")
	}

	//Snssai without sd
	body = []byte(`{
		"sst": 10
	}`)

	snssai = &TSnssai{}
	err = json.Unmarshal(body, snssai)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	subKeys = snssai.GenerateGrpcGetKey()
	if len(subKeys) != 1 {
		t.Fatalf("TSnssai.GenerateGrpcGetKey didn't return right value!")
	}

	if subKeys[0].SubKey1 != "10" || subKeys[0].SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSnssai.GenerateGrpcGetKey didn't return right value!")
	}
}
