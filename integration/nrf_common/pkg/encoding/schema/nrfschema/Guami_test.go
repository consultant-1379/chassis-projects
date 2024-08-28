package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGenerateGrpcKeyForGuami(t *testing.T) {
	body := []byte(`{
		"plmnId": {
			"mcc": "460",
			"mnc": "00"
		},
		"amfId": "123456"
	}`)

	guami := &TGuami{}

	err := json.Unmarshal(body, guami)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	subKey := guami.GenerateGrpcKey()
	if subKey.SubKey1 != "46000" || subKey.SubKey2 != "123456" {
		t.Fatalf("TGuami.GenerateGrpcKey didn't return right value!")
	}
}
