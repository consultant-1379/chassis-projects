package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetPlmnID(t *testing.T) {
	body := []byte(`{
		"mcc": "460",
		"mnc": "00"
	}`)

	plmnID := &TPlmnId{}
	err := json.Unmarshal(body, plmnID)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}

	if plmnID.GetPlmnID() != "46000" {
		t.Fatalf("SubscriptionData.GetTargetPlmnID didn't return right plmnID!")
	}

	body = []byte(`{
		"mcc": "460",
		"mnc": "000"
	}`)

	plmnID = &TPlmnId{}
	err = json.Unmarshal(body, plmnID)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}

	if plmnID.GetPlmnID() != "460000" {
		t.Fatalf("SubscriptionData.GetTargetPlmnID didn't return right plmnID!")
	}
}
