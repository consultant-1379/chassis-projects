package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsIdentityRangeValid(t *testing.T) {
	//right IdentityRange
	body := []byte(`{
	    "start": "1111",
	    "end": "2222"
    }`)

	identityRange := &TIdentityRange{}

	err := json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

	//right IdentityRange
	body = []byte(`{
	    "pattern": "string"
    }`)

	identityRange = &TIdentityRange{}

	err = json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

	//wrong IdentityRange
	body = []byte(`{
    }`)

	identityRange = &TIdentityRange{}

	err = json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

	//wrong IdentityRange
	body = []byte(`{
		"start": "1111",
		"end": "2222",
		"pattern": "string"
    }`)

	identityRange = &TIdentityRange{}

	err = json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

	//wrong IdentityRange
	body = []byte(`{
		"start": "1111",
		"pattern": "string"
    }`)

	identityRange = &TIdentityRange{}

	err = json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

	//wrong IdentityRange
	body = []byte(`{
		"end": "2222",
		"pattern": "string"
    }`)

	identityRange = &TIdentityRange{}

	err = json.Unmarshal(body, identityRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if identityRange.IsValid() {
		t.Fatalf("TIdentityRange.IsValid didn't return right value!")
	}

}
