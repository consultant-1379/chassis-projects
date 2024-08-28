package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetInvalidTacRangeIndexs(t *testing.T) {
	//taiRange without tacRangeList
	body := []byte(`{
		"plmnId": {
          "mcc": "460",
          "mnc": "05"
        }
	}`)

	taiRange := &TTaiRange{}
	err := json.Unmarshal(body, taiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if taiRange.GetInvalidTacRangeIndexs() != nil {
		t.Fatalf("TTaiRange.GetInvalidTacRangeIndexs didn't return riht value!")
	}

	//taiRange with right tacRangeList
	body = []byte(`{
		"plmnId": {
            "mcc": "460",
            "mnc": "05"
        },
		"tacRangeList": [
		    {
			    "start": "1234",
			    "end": "1234"
		    },
		    {
			    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		    }
		]
	}`)

	taiRange = &TTaiRange{}
	err = json.Unmarshal(body, taiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if taiRange.GetInvalidTacRangeIndexs() != nil {
		t.Fatalf("TTaiRange.GetInvalidTacRangeIndexs didn't return riht value!")
	}

	//taiRange with wrong tacRangeList
	body = []byte(`{
		"plmnId": {
            "mcc": "460",
            "mnc": "05"
        },
		"tacRangeList": [
		    {
			    "start": "1234",
			    "end": "1234"
		    },
		    {
			    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		    },
			{},
		    {
			    "start": "1234",
			    "end": "1234",
				"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		    },
		    {
				"start": "1234",
			    "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		    },
			{
			    "end": "1234",
				"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		    }
		]
	}`)

	taiRange = &TTaiRange{}
	err = json.Unmarshal(body, taiRange)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if taiRange.GetInvalidTacRangeIndexs() == nil {
		t.Fatalf("TTaiRange.GetInvalidTacRangeIndexs didn't return riht value!")
	}

}
