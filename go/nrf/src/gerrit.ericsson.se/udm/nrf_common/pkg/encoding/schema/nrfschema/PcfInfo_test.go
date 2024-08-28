package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetPcfInvalidSupiRangeIndexs(t *testing.T) {
	//PcfInfo without supiRanges
	body := []byte(`{
	    "dnnList": ["dnn1", "dnn2"]
    }`)

	pcfInfo := &TPcfInfo{}
	err := json.Unmarshal(body, pcfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if pcfInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TPcfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//PcfInfo with right supiRanges
	body = []byte(`{
	    "dnnList": ["dnn1", "dnn2"],
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

	pcfInfo = &TPcfInfo{}
	err = json.Unmarshal(body, pcfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if pcfInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TPcfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//PcfInfo with wrong supiRanges
	body = []byte(`{
	    "dnnList": ["dnn1", "dnn2"],
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

	pcfInfo = &TPcfInfo{}
	err = json.Unmarshal(body, pcfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if pcfInfo.GetInvalidSupiRangeIndexs() == nil {
		t.Fatalf("TPcfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}
}
