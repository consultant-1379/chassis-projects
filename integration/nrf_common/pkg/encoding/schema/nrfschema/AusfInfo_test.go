package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetAusfInvalidSupiRangeIndexs(t *testing.T) {
	//ausfInfo without supiRanges
	body := []byte(`{
	    "groupId": "1234",
		"routingIndicators": [
		    "1111", "2222", "3333"
		]
    }`)

	ausfInfo := &TAusfInfo{}
	err := json.Unmarshal(body, ausfInfo)
	if err != nil {
		t.Fatalf("Unmarshal erro, %v", err)
	}

	if ausfInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TAusfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//ausfInfo with right supiRanges
	body = []byte(`{
	    "groupId": "1234",
		"supiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
				"pattern": "string"
			}
		],
		"routingIndicators": [
		    "1111", "2222", "3333"
		]
    }`)

	ausfInfo = &TAusfInfo{}
	err = json.Unmarshal(body, ausfInfo)
	if err != nil {
		t.Fatalf("Unmarshal erro, %v", err)
	}

	if ausfInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TAusfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//ausfInfo with wrong supiRanges
	body = []byte(`{
	    "groupId": "1234",
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
		],
		"routingIndicators": [
		    "1111", "2222", "3333"
		]
    }`)

	ausfInfo = &TAusfInfo{}
	err = json.Unmarshal(body, ausfInfo)
	if err != nil {
		t.Fatalf("Unmarshal erro, %v", err)
	}

	if ausfInfo.GetInvalidSupiRangeIndexs() == nil {
		t.Fatalf("TAusfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}
}

func TestGenerateNfGroupCond(t *testing.T) {
	//ausfInfo without groupId
	body := []byte(`{
		"supiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
				"pattern": "string"
			}
		],
		"routingIndicators": [
		    "1111", "2222", "3333"
		]
    }`)

	ausfInfo := &TAusfInfo{}
	err := json.Unmarshal(body, ausfInfo)
	if err != nil {
		t.Fatalf("Unmarshal erro, %v", err)
	}

	if ausfInfo.GenerateNfGroupCond() != nil {
		t.Fatalf("TAusfInfo.GenerateNfGroupCond didn't return right value!")
	}

	//ausfInfo with groupId
	body = []byte(`{
		"groupId": "group01",
		"supiRanges": [
		    {
				"start": "1111",
				"end": "2222"
			},
			{
				"pattern": "string"
			}
		],
		"routingIndicators": [
		    "1111", "2222", "3333"
		]
    }`)

	ausfInfo = &TAusfInfo{}
	err = json.Unmarshal(body, ausfInfo)
	if err != nil {
		t.Fatalf("Unmarshal erro, %v", err)
	}

	nfGroupCond := ausfInfo.GenerateNfGroupCond()

	if nfGroupCond == nil {
		t.Fatalf("TAusfInfo.GenerateNfGroupCond didn't return right value!")
	}
}
