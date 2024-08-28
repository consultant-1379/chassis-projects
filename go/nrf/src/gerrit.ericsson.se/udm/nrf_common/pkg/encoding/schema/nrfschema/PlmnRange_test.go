package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsPlmnRangeValid(t *testing.T) {
	//right PlmnRange
	body := []byte(`{
	    "start": "1111",
	    "end": "2222"
    }`)

	PlmnRange := &TPlmnRange{}
	err := json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}

	//right PlmnRange
	body = []byte(`{
	    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	PlmnRange = &TPlmnRange{}
	err = json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}

	//wrong PlmnRange
	body = []byte(`{
    }`)

	PlmnRange = &TPlmnRange{}
	err = json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}

	//wrong PlmnRange
	body = []byte(`{
		"start": "1111",
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	PlmnRange = &TPlmnRange{}
	err = json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}

	//wrong PlmnRange
	body = []byte(`{
		"start": "1111",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	PlmnRange = &TPlmnRange{}
	err = json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}

	//wrong PlmnRange
	body = []byte(`{
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	PlmnRange = &TPlmnRange{}
	err = json.Unmarshal(body, PlmnRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if PlmnRange.IsValid() {
		t.Fatalf("TPlmnRange.IsValid didn't return right value!")
	}
}
