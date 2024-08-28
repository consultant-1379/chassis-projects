package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetUdmInvalidSupiRangeIndexs(t *testing.T) {
	//UdmInfo without supiRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udmInfo := &TUdmInfo{}
	err := json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//UdmInfo with right supiRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//UdmInfo with wrong supiRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidSupiRangeIndexs() == nil {
		t.Fatalf("TUdmInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}
}

func TestGetUdmInvalidGpsiRangeIndexs(t *testing.T) {
	//UdmInfo without gpsiRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udmInfo := &TUdmInfo{}
	err := json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidGpsiRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//UdmInfo with right gpsiRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidGpsiRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//UdmInfo with wrong gpsiRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidGpsiRangeIndexs() == nil {
		t.Fatalf("TUdmInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}
}

func TestGetUdmInvalidEGIRangeIndexs(t *testing.T) {
	//UdmInfo without externalGroupIdentifiersRanges
	body := []byte(`{
	    "groupId": "001"
    }`)

	udmInfo := &TUdmInfo{}
	err := json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidEGIRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}

	//UdmInfo with right externalGroupIdentifiersRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidEGIRangeIndexs() != nil {
		t.Fatalf("TUdmInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}

	//UdmInfo with wrong externalGroupIdentifiersRanges
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

	udmInfo = &TUdmInfo{}
	err = json.Unmarshal(body, udmInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if udmInfo.GetInvalidEGIRangeIndexs() == nil {
		t.Fatalf("TUdmInfo.GetInvalidEGIRangeIndexs didn't return right value!")
	}
}
