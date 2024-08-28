package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetChfInvalidSupiRangeIndexs(t *testing.T) {
	//ChfInfo without supiRangeList
	body := []byte(`{
	    "gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo := &TChfInfo{}
	err := json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if chfInfo.GetInvalidSupiRangeIndexs() != nil {
		t.Fatalf("TChfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//ChfInfo with right supiRangeList
	body = []byte(`{
		"supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
	    "gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidSupiRangeIndexs := chfInfo.GetInvalidSupiRangeIndexs()
	if invalidSupiRangeIndexs != nil {
		t.Fatalf("TChfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	//ChfInfo with wrong supiRangeList
	body = []byte(`{
		"supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
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
				"end": "1111",
			    "pattern": "string"
			},
			{}
		],
	    "gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidSupiRangeIndexs = chfInfo.GetInvalidSupiRangeIndexs()
	if invalidSupiRangeIndexs == nil {
		t.Fatalf("TChfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}

	if len(invalidSupiRangeIndexs) != 4 {
		t.Fatalf("TChfInfo.GetInvalidSupiRangeIndexs didn't return right value!")
	}
}

func TestGetChfInvalidGpsiRangeIndexs(t *testing.T) {
	//ChfInfo without gpsiRangeList
	body := []byte(`{
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo := &TChfInfo{}
	err := json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if chfInfo.GetInvalidGpsiRangeIndexs() != nil {
		t.Fatalf("TChfInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//ChfInfo with right gpsiRangeList
	body = []byte(`{
		"gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidGpsiRangeIndexs := chfInfo.GetInvalidGpsiRangeIndexs()
	if invalidGpsiRangeIndexs != nil {
		t.Fatalf("TChfInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	//ChfInfo with wrong gpsiRangeList
	body = []byte(`{
		"gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
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
				"end": "1111",
			    "pattern": "string"
			},
			{}
		],
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidGpsiRangeIndexs = chfInfo.GetInvalidGpsiRangeIndexs()
	if invalidGpsiRangeIndexs == nil {
		t.Fatalf("TChfInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}

	if len(invalidGpsiRangeIndexs) != 4 {
		t.Fatalf("TChfInfo.GetInvalidGpsiRangeIndexs didn't return right value!")
	}
}

func TestGetChfInvalidPlmnRangeIndexs(t *testing.T) {
	//ChfInfo without plmnRangeList
	body := []byte(`{
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo := &TChfInfo{}
	err := json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if chfInfo.GetInvalidPlmnRangeIndexs() != nil {
		t.Fatalf("TChfInfo.GetInvalidPlmnRangeIndexs didn't return right value!")
	}

	//ChfInfo with right plmnRangeList
	body = []byte(`{
		"gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidPlmnRangeIndexs := chfInfo.GetInvalidPlmnRangeIndexs()
	if invalidPlmnRangeIndexs != nil {
		t.Fatalf("TChfInfo.GetInvalidPlmnRangeIndexs didn't return right value!")
	}

	//ChfInfo with wrong plmnRangeList
	body = []byte(`{
		"plmnRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
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
				"end": "1111",
			    "pattern": "string"
			},
			{}
		],
	    "supiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		],
		"gpsiRangeList": [
			{
				"start": "1111",
				"end": "2222"
			},
			{
			    "pattern": "string"
			}
		]
    }`)

	chfInfo = &TChfInfo{}
	err = json.Unmarshal(body, chfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	invalidPlmnRangeIndexs = chfInfo.GetInvalidPlmnRangeIndexs()
	if invalidPlmnRangeIndexs == nil {
		t.Fatalf("TChfInfo.GetInvalidPlmnRangeIndexs didn't return right value!")
	}

	if len(invalidPlmnRangeIndexs) != 4 {
		t.Fatalf("TChfInfo.GetInvalidPlmnRangeIndexs didn't return right value!")
	}
}

func TestChfCreateNfInfo(t *testing.T) {
	//body := []byte(`{
	//	"plmnRangeList": [
	//		{
	//			"start": "1111",
	//			"end": "2222"
	//		},
	//		{
	//		    "pattern": "string"
	//		},
	//		{
	//			"start": "1111",
	//			"end": "1111",
	//		    "pattern": "string"
	//		},
	//		{}
	//	],
	//    	"supiRangeList": [
	//		{
	//			"start": "1111",
	//			"end": "2222"
	//		},
	//		{
	//		    "pattern": "string"
	//		}
	//	],
	//	"gpsiRangeList": [
	//		{
	//			"start": "1111",
	//			"end": "2222"
	//		},
	//		{
	//		    "pattern": "string"
	//		}
	//	]
	//}`)
	body2 := []byte(`{

	}`)
	//chfInfo := &TChfInfo{}
	//err := json.Unmarshal(body, chfInfo)
	//if err != nil {
	//	t.Fatalf("Unmarshal error, %v", err)
	//}
	//chfInfoJson := chfInfo.createNfInfo()

	chfInfo2 := &TChfInfo{}
	err2 := json.Unmarshal(body2, chfInfo2)
	if err2 != nil {
		t.Fatalf("Unmarshal error, %v", err2)
	}
	chfInfoJson2 := chfInfo2.createNfInfo()
	if chfInfoJson2 != `"chfInfo":{"supiRangeList":[{"pattern":"RESERVED_EMPTY_SUPI_RANGE_PATTERN"}],"gpsiRangeList":[{"pattern":"RESERVED_EMPTY_GPSI_RANGE_PATTERN"}],"plmnRangeList":[{"pattern":"RESERVED_EMPTY_PLMN_RANGE_PATTERN"}],"supiMatchAll":"MATCH_ALL","gpsiMatchAll":"MATCH_ALL"}` {
		t.Fatal("chfInfo helper should matched, but fail")
	}
}
