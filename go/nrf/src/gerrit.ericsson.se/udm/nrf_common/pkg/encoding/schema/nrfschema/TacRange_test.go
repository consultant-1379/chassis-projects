package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsTacRangeValid(t *testing.T) {
	//right TacRange
	body := []byte(`{
	    "start": "1111",
	    "end": "2222"
    }`)

	tacRange := &TTacRange{}
	err := json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}

	//right TacRange
	body = []byte(`{
	    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	tacRange = &TTacRange{}
	err = json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}

	//wrong TacRange
	body = []byte(`{
    }`)

	tacRange = &TTacRange{}
	err = json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}

	//wrong TacRange
	body = []byte(`{
		"start": "1111",
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	tacRange = &TTacRange{}
	err = json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}

	//wrong TacRange
	body = []byte(`{
		"start": "1111",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	tacRange = &TTacRange{}
	err = json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}

	//wrong TacRange
	body = []byte(`{
	    "end": "2222",
		"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
    }`)

	tacRange = &TTacRange{}
	err = json.Unmarshal(body, tacRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if tacRange.IsValid() {
		t.Fatalf("TTacRange.IsValid didn't return right value!")
	}
}
