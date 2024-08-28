package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsNotifConditionValid(t *testing.T) {
	//NotifCondition with only monitoredAttributes is valid
	body := []byte(`{
		"monitoredAttributes": ["nfStatus", "plmnList"]
	}`)

	notifCondition := &TNotifCondition{}

	err := json.Unmarshal(body, notifCondition)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !notifCondition.IsValid() {
		t.Fatalf("TNotifCondition.IsValid didn't return right value!")
	}

	//NotifCondition with only unmonitoredAttributes is valid
	body = []byte(`{
		"unmonitoredAttributes": ["nfStatus", "plmnList"]
	}`)

	notifCondition = &TNotifCondition{}

	err = json.Unmarshal(body, notifCondition)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !notifCondition.IsValid() {
		t.Fatalf("TNotifCondition.IsValid didn't return right value!")
	}

	//NotifCondition with both monitoredAttributes and unmonitoredAttributes is valid
	body = []byte(`{
		"monitoredAttributes": ["nfStatus", "plmnList"],
		"unmonitoredAttributes": ["sNssais", "nsiList"]
	}`)

	notifCondition = &TNotifCondition{}

	err = json.Unmarshal(body, notifCondition)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if notifCondition.IsValid() {
		t.Fatalf("TNotifCondition.IsValid didn't return right value!")
	}
}
