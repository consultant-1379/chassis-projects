package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetUdrInvalidSupiRangeIndexs(t *testing.T) {
	//UdrInfo without supiRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udrInfo := &TUdrInfo{}
	err := json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//UdrInfo with right supiRanges
	body = []byte(`{
	    "groupId": "001",
		"supiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//UdrInfo with wrong supiRanges
	body = []byte(`{
	    "groupId": "001",
		"supiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			},
			{
				"pattern": "string"
			},
			{
				"start": "1111"
			},
			{
				"end": "1111"
			},
			{
				"start": "1111",
				"pattern": "string"
			},
			{
				"end": "1111",
				"pattern": "string"
			},
			{
				"start": "1111",
				"end": "2222",
				"pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidSupiRangeIndexs() == nil {
		t.Fatalf("TUdrInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}
}

func TestGetUdrInvalidGpsiRangeIndexs(t *testing.T) {
	//UdrInfo without gpsiRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udrInfo := &TUdrInfo{}
	err := json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidGpsiRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//UdrInfo with right gpsiRanges
	body = []byte(`{
	    "groupId": "001",
		"gpsiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidGpsiRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//UdrInfo with wrong gpsiRanges
	body = []byte(`{
	    "groupId": "001",
		"gpsiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			},
			{
				"pattern": "string"
			},
			{
				"start": "1111"
			},
			{
				"end": "1111"
			},
			{
				"start": "1111",
				"pattern": "string"
			},
			{
				"end": "1111",
				"pattern": "string"
			},
			{
				"start": "1111",
				"end": "2222",
				"pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidGpsiRangeIndexs() == nil {
		t.Fatalf("TUdrInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}
}

func TestGetUdrInvalidEGIRangeIndexs(t *testing.T) {
	//UdrInfo without externalGroupIdentifiersRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udrInfo := &TUdrInfo{}
	err := json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidEGIRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}

	//UdrInfo with right externalGroupIdentifiersRanges
	body = []byte(`{
	    "groupId": "001",
		"externalGroupIdentifiersRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidEGIRangeIndexs() != nil {
		t.Fatalf("TUdrInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}

	//UdrInfo with wrong externalGroupIdentifiersRanges
	body = []byte(`{
	    "groupId": "001",
		"externalGroupIdentifiersRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
			},
			{
				"pattern": "string"
			},
			{
				"start": "1111"
			},
			{
				"end": "1111"
			},
			{
				"start": "1111",
				"pattern": "string"
			},
			{
				"end": "1111",
				"pattern": "string"
			},
			{
				"start": "1111",
				"end": "2222",
				"pattern": "string"
			}
		]
    }`)

	udrInfo = &TUdrInfo{}
	err = json.Unmarshal(body, udrInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udrInfo.GetInvalidEGIRangeIndexs() == nil {
		t.Fatalf("TUdrInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}
}
