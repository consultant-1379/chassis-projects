package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsChfServiceInfoValid(t *testing.T) {
	//ChfServiceInfo without primaryChfServiceInstance and secondaryChfServiceInstance
	body := []byte(`{
	}`)

	chfServiceInfo := &TChfServiceInfo{}
	err := json.Unmarshal(body, chfServiceInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !chfServiceInfo.IsValid() {
		t.Fatalf("TChfServiceInfo.IsValid didn't return right value!")
	}

	//ChfServiceInfo with primaryChfServiceInstance
	body = []byte(`{
		"primaryChfServiceInstance": "serv01"
	}`)

	chfServiceInfo = &TChfServiceInfo{}
	err = json.Unmarshal(body, chfServiceInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !chfServiceInfo.IsValid() {
		t.Fatalf("TChfServiceInfo.IsValid didn't return right value!")
	}

	//ChfServiceInfo with secondaryChfServiceInstance
	body = []byte(`{
		"secondaryChfServiceInstance": "serv01"
	}`)

	chfServiceInfo = &TChfServiceInfo{}
	err = json.Unmarshal(body, chfServiceInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !chfServiceInfo.IsValid() {
		t.Fatalf("TChfServiceInfo.IsValid didn't return right value!")
	}

	//ChfServiceInfo with primaryChfServiceInstance and secondaryChfServiceInstance
	body = []byte(`{
		"primaryChfServiceInstance": "serv01",
		"secondaryChfServiceInstance": "serv02"
	}`)

	chfServiceInfo = &TChfServiceInfo{}
	err = json.Unmarshal(body, chfServiceInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if chfServiceInfo.IsValid() {
		t.Fatalf("TChfServiceInfo.IsValid didn't return right value!")
	}
}
