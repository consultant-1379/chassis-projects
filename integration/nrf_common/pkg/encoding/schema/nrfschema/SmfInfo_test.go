package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetSmfInvalidTaiRangeIndexs(t *testing.T) {
	//SmfInfo without taiRangeList
	body := []byte(`{
	    "dnnList": ["dnn1", "dnn2"]
    }`)

	smfInfo := &TSmfInfo{}
	err := json.Unmarshal(body, smfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if smfInfo.GetInvalidTaiRangeIndexs() != nil {
		t.Fatalf("TSmfInfo.GetInvalidTaiRangeIndexs didn't right value!")
	}

	//SmfInfo with right taiRangeList
	body = []byte(`{
	    "dnnList": ["dnn1", "dnn2"],
		"taiRangeList": [
	        {
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
	        },
        	    {
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
	        }
	    ]
    }`)

	smfInfo = &TSmfInfo{}
	err = json.Unmarshal(body, smfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if smfInfo.GetInvalidTaiRangeIndexs() != nil {
		t.Fatalf("TSmfInfo.GetInvalidTaiRangeIndexs didn't right value!")
	}

	//SmfInfo with wrong taiRangeList
	body = []byte(`{
	    "dnnList": ["dnn1", "dnn2"],
		"taiRangeList": [
	        {
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
						"end": "1234",
				        "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"	
					}
		        ]
	        },
        	    {
		        "plmnId": {
                    "mcc": "460",
                    "mnc": "05"
                },
		        "tacRangeList": [
					{},
					{
			            "start": "1234",
			            "end": "1234",
				        "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"						
					},
					{
						"end": "1234",
				        "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"	
					},
		            {
			            "start": "1234",
			            "end": "1234"
		            },
		            {
			            "pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		            }
		        ]
	        }
	    ]
    }`)

	smfInfo = &TSmfInfo{}
	err = json.Unmarshal(body, smfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if smfInfo.GetInvalidTaiRangeIndexs() == nil {
		t.Fatalf("TSmfInfo.GetInvalidTaiRangeIndexs didn't right value!")
	}

}
