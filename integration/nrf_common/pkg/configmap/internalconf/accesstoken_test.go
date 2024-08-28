package internalconf

import (
	"encoding/json"
	"testing"
)

func TestParseAccessTokenInternalConf(t *testing.T) {
	jsonContent := []byte(`{
  	"httpServer": {
	  "idleTimeout": 1,
	  "activeTimeout": 2,
	  "httpWithXVersion": false
	},
	"dbProxy": {
	  "connectionNum": 100,
	  "grpcContextTimeout": 3
	},
	"accessToken": {
	  "expiredTime" : 3600
	},
	"nfServices": {
	  "NRF": {
		"nnrf-nfm": {
		  "allowed-nf-type": [
			"AMF",
			"SMF",
			"UDM",
			"AUSF",
			"NEF",
			"PCF",
			"SMSF",
			"NSSF",
			"UPF",
			"BSF",
			"CHF",
			"NRF"
		  ]
		},
		"nnrf-disc": {
		  "allowed-nf-type": [
			"AMF",
			"SMF",
			"PCF",
			"NEF",
			"NSSF",
			"SMSF",
			"AUSF",
			"CHF",
			"NRF"
		  ]
		}
	  },
	  "UDM": {
		"nudm-sdm": {
		  "allowed-nf-type": [
			"AMF",
			"SMF",
			"AUSF",
			"SMSF"
		  ]
		},
		"nudm-abc": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		},
		"nudm-uecm": {
		  "allowed-nf-type": [
			"AMF",
			"SMF",
			"SMSF",
			"NEF",
			"GMLC"
		  ]
		},
		"nudm-ueau": {
		  "allowed-nf-type": [
			"AUSF"
		  ]
		},
		"nudm-ee": {
		  "allowed-nf-type": [
			"NEF"
		  ]
		},
		"nudm-pp": {
		  "allowed-nf-type": [
			"NEF"
		  ]
		},
		"nudm-niddau": {
		  "allowed-nf-type": [
			"*"
		  ]
		}
	  },
	  "AMF": {
		"namf-comm": {
		  "allowed-nf-type": [
			"SMF",
			"SMSF",
			"PCF",
			"LMF",
			"NEF",
			"UDM"
		  ]
		},
		"namf-evts": {
		  "allowed-nf-type": [
			"NEF",
			"SMF",
			"UDM"
		  ]
		},
		"namf-mt": {
		  "allowed-nf-type": [
			"UDM",
			"SMSF"
		  ]
		},
		"namf-loc": {
		  "allowed-nf-type": [
			"UDM",
			"GMLC"
		  ]
		}
	  },
	  "SMF": {
		"nsmf-pdusession": {
		  "allowed-nf-type": [
			"SMF",
			"AMF"
		  ]
		},
		"nsmf-event-exposure": {
		  "allowed-nf-type": [
			"NEF",
			"AMF"
		  ]
		}
	  },
	  "AUSF": {
		"nausf-auth": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		},
		"nausf-sorprotection": {
		  "allowed-nf-type": [
			"UDM"
		  ]
		},
		"nausf-upuprotection": {
		  "allowed-nf-type": [
			"*"
		  ]
		}
	  },
	  "NEF": {
		"nnef-pfdmanagement": {
		  "allowed-nf-type": [
			"SMF",
			"AF"
		  ]
		}
	  },
	  "PCF": {
		"npcf-am-policy-control": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		},
		"npcf-smpolicycontrol": {
		  "allowed-nf-type": [
			"SMF"
		  ]
		},
		"npcf-policyauthorization": {
		  "allowed-nf-type": [
			"AF",
			"NEF"
		  ]
		},
		"npcf-bdtpolicycontrol": {
		  "allowed-nf-type": [
			"NEF"
		  ]
		},
		"npcf-eventexposure": {
		  "allowed-nf-type": [
			"NEF"
		  ]
		},
		"npcf-ue-policy-control": {
		  "allowed-nf-type": [
			"AMF",
			"PCF"
		  ]
		}
	  },
	  "SMSF": {
		"nsmsf-sms": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		}
	  },
	  "NSSF": {
		"nnssf-nsselection": {
		  "allowed-nf-type": [
			"AMF",
			"NSSF"
		  ]
		},
		"nnssf-nssaiavailability": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		}
	  },
	  "UDR": {
		"nudr-dr": {
		  "allowed-nf-type": [
			"UDM",
			"PCF",
			"NEF"
		  ]
		}
	  },
	  "LMF": {
		"nlmf-loc": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		}
	  },
	  "5G-EIR": {
		"n5g-eir-eic": {
		  "allowed-nf-type": [
			"AMF"
		  ]
		}
	  },
	  "BSF": {
		"nbsf-management": {
		  "allowed-nf-type": [
			"PCF",
			"NEF",
			"AF"
		  ]
		}
	  },
	  "CHF": {
		"nchf-spendinglimitcontrol": {
		  "allowed-nf-type": [
			"PCF"
		  ]
		},
		"nchf-convergedcharging": {
		  "allowed-nf-type": [
			"*"
		  ]
		}
	  },
	  "NWDAF": {
		"nnwdaf-eventssubscription": {
		  "allowed-nf-type": [
			"PCF",
			"NSSF"
		  ]
		},
		"nnwdaf-analyticsinfo": {
		  "allowed-nf-type": [
			"PCF",
			"NSSF"
		  ]
		}
	  }
	}
  }`)

	internalConf := &InternalAccessTokenConf{}

	err := json.Unmarshal(jsonContent, internalConf)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	internalConf.ParseConf()

	//NfTypeToNfservices
	nfservices := NfTypeToNfservices["NRF"]
	if len(nfservices) != 2 {
		t.Fatalf("NRF services parsed error")
	}

	isNnrfNfmExist, isNnrfDiscExist := false, false

	for _, nfservice := range nfservices {

		switch nfservice {
		case "nnrf-nfm":
			isNnrfNfmExist = true

		case "nnrf-disc":
			isNnrfDiscExist = true
		}
	}

	if !(isNnrfNfmExist && isNnrfDiscExist) {
		t.Fatalf("NRF services parsed error")
	}

	//Nfservcies
	nfType := NfServiceToNfType["nnrf-nfm"]
	if nfType != "NRF" {
		t.Fatalf("NRF services parsed error")
	}

	//NfServiceAllowedNfType
	allowedNfTypes := NfServiceAllowedNfType["nnrf-disc"]
	isAusfAllowed := false
	for _, allowedNfType := range allowedNfTypes {
		if allowedNfType == "AUSF" {
			isAusfAllowed = true
		}

		if allowedNfType == "UDR" {
			t.Fatalf("NRF services parsed error")
		}

	}
	if !isAusfAllowed {
		t.Fatalf("NRF services parsed error")
	}

	if TokenExpiredTime != 3600 {
		t.Fatalf("expired time parsed error!")
	}
}
