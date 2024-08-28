package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestGetInvalidN2InterfaceAmfInfoIndex(t *testing.T) {
	//amfInfo without n2InterfaceAmfInfo
	body := []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "01"
        },
        "amfId": "800001"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      }
    ],
    "taiList": [
      {
        "plmnId": {
          "mcc": "240",
          "mnc": "81"
        },
        "tac": "00000C"
      },
      {
        "plmnId": {
          "mcc": "240",
          "mnc": "81"
        },
        "tac": "00000B"
      }
    ]
    }`)

	amfInfo := &TAmfInfo{}
	err := json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidN2InterfaceAmfInfoIndex() != "" {
		t.Fatalf("TAmfInfo.GetInvalidN2InterfaceAmfInfoIndex didn't return right value!")
	}

	//amfInfo with right n2InterfaceAmfInfo
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
	"n2InterfaceAmfInfo": {
	  "ipv4EndpointAddress": [
	    "10.10.10.10"
	  ]
	}
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidN2InterfaceAmfInfoIndex() != "" {
		t.Fatalf("TAmfInfo.GetInvalidN2InterfaceAmfInfoIndex didn't return right value!")
	}

	//amfInfo with right n2InterfaceAmfInfo
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
	"n2InterfaceAmfInfo": {
	  "ipv6EndpointAddress": [
	    "1030::C9B4:FF12:48AA:1A2B"
	  ]
	}
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidN2InterfaceAmfInfoIndex() != "" {
		t.Fatalf("TAmfInfo.GetInvalidN2InterfaceAmfInfoIndex didn't return right value!")
	}

	//amfInfo with right n2InterfaceAmfInfo
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
	"n2InterfaceAmfInfo": {
	  "ipv4EndpointAddress": [
	    "10.10.10.10"
	  ],
	  "ipv6EndpointAddress": [
	    "1030::C9B4:FF12:48AA:1A2B"
	  ]
	}
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidN2InterfaceAmfInfoIndex() != "" {
		t.Fatalf("TAmfInfo.GetInvalidN2InterfaceAmfInfoIndex didn't return right value!")
	}

	//amfInfo with wrong n2InterfaceAmfInfo
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
	"n2InterfaceAmfInfo": {
      "amfName": "amf01"
	}
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidN2InterfaceAmfInfoIndex() == "" {
		t.Fatalf("TAmfInfo.GetInvalidN2InterfaceAmfInfoIndex didn't return right value!")
	}
}

func TestGetAmfInvalidTaiRangeIndexs(t *testing.T) {
	//amfInfo without taiRangeList
	body := []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
    "taiList": [
      {
        "plmnId": {
          "mcc": "240",
          "mnc": "81"
        },
        "tac": "00000C"
      },
      {
        "plmnId": {
          "mcc": "240",
          "mnc": "81"
        },
        "tac": "00000B"
      }
    ]
    }`)

	amfInfo := &TAmfInfo{}
	err := json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidTaiRangeIndexs() != nil {
		t.Fatalf("TAmfInfo.GetInvalidTaiRangeIndexs didn't return right value!")
	}

	//amfInfo with right taiRangeList
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
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

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidTaiRangeIndexs() != nil {
		t.Fatalf("TAmfInfo.GetInvalidTaiRangeIndexs didn't return right value!")
	}

	//amfInfo with wrong taiRangeList
	body = []byte(`{
    "amfSetId": "amfSet01",
    "amfRegionId": "amfRegion01",
    "guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
    "taiRangeList": [
	  {
		"plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
		"tacRangeList": [
		  {
			"start": "1234",
			"end": "1234",
			"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		  },
		  {}
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
		  {},
		  {
			"end": "1234",
			"pattern": "^([A-Fa-f0-9]{4}|[A-Fa-f0-9]{6})$"
		  }
		]
	  }   
	]
	
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GetInvalidTaiRangeIndexs() == nil {
		t.Fatalf("TAmfInfo.GetInvalidTaiRangeIndexs didn't return right value!")
	}
}

func TestGenerateAmfCond(t *testing.T) {
	//amfInfo without amfSetId and amfRegionId
	body := []byte(`{
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
		  }
		]
	  }   
	]	
    }`)

	amfInfo := &TAmfInfo{}
	err := json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GenerateAmfCond() != nil {
		t.Fatalf("TAmfInfo.GenerateAmfCond didn't return right value!")
	}

	//amfInfo with amfSetId and amfRegionId
	body = []byte(`{
	"amfSetId": "amfSet01",
	"amfRegionId": "amfRegion01",
	"guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
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
		  }
		]
	  }   
	]	
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	amfConds := amfInfo.GenerateAmfCond()
	if amfConds == nil {
		t.Fatalf("TAmfInfo.GenerateAmfCond didn't return right value!")
	}

	if len(amfConds) != 3 {
		t.Fatalf("TAmfInfo.GenerateAmfCond didn't return right value!")
	}
}

func TestGenerateGuamiListCond(t *testing.T) {
	//amfInfo without guamiList
	body := []byte(`{
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
		  }
		]
	  }   
	]	
    }`)

	amfInfo := &TAmfInfo{}
	err := json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if amfInfo.GenerateGuamiListCond() != nil {
		t.Fatalf("TAmfInfo.GenerateGuamiListCond didn't return right value!")
	}

	//amfInfo with guamiList
	body = []byte(`{
	"amfSetId": "amfSet01",
	"amfRegionId": "amfRegion01",
	"guamiList": [
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "02"
        },
        "amfId": "800002"
      },
      {
        "plmnId": {
          "mcc": "460",
          "mnc": "05"
        },
        "amfId": "800005"
      }
    ],
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
		  }
		]
	  }   
	]	
    }`)

	amfInfo = &TAmfInfo{}
	err = json.Unmarshal(body, amfInfo)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	guamiListConds := amfInfo.GenerateGuamiListCond()
	if guamiListConds == nil {
		t.Fatalf("TAmfInfo.GenerateGuamiListCond didn't return right value!")
	}

	if len(guamiListConds) != 2 {
		t.Fatalf("TAmfInfo.GenerateGuamiListCond didn't return right value!")
	}
}
