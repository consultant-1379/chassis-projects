package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsSupiRangeValid(t *testing.T) {
	//right SupiRange
	body := []byte(`{
	    "start": "1111",
	    "end": "2222"
    }`)

	supiRange := &TSupiRange{}
	err := json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}

	//right SupiRange
	body = []byte(`{
	    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	supiRange = &TSupiRange{}
	err = json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}

	//wrong SupiRange
	body = []byte(`{
    }`)

	supiRange = &TSupiRange{}
	err = json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}

	//wrong SupiRange
	body = []byte(`{
		"start": "1111",
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	supiRange = &TSupiRange{}
	err = json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}

	//wrong SupiRange
	body = []byte(`{
		"start": "1111",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	supiRange = &TSupiRange{}
	err = json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}

	//wrong SupiRange
	body = []byte(`{
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	supiRange = &TSupiRange{}
	err = json.Unmarshal(body, supiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if supiRange.IsValid() {
		t.Fatalf("TSupiRange.IsValid didn't return right value!")
	}
}
