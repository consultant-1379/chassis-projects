package nfdiscfilter

import (
	"testing"
)

func TestIsMatchedChfSupportedPlmn(t *testing.T) {
	filter := &NFCHFInfoFilter{}
	nfProfile := []byte(`{
	"plmnRangeList": [
      	{
        	"start": "46000",
        	"end": "46011"
      	},
        {
        	"pattern": "^46[3-4]{1}[0-9]{2,3}$"
      	},
	{
		"start": "460111"
		"end": "460222"
	}
    	]
      	}`)
	if !filter.isMatchedChfSupportedPlmn("46010", []byte(nfProfile)) {
		t.Fatalf("should be matched , but failed")
	}
	if !filter.isMatchedChfSupportedPlmn("460120", []byte(nfProfile)) {
		t.Fatalf("should be matched , but failed")
	}
	if !filter.isMatchedChfSupportedPlmn("463999", []byte(nfProfile)) {
		t.Fatalf("should be matched , but failed")
	}
	if filter.isMatchedChfSupportedPlmn("46100", []byte(nfProfile)) {
		t.Fatalf("should not be matched , but matched")
	}
	nfProfile = []byte(`{
	"plmnRangeList": [
      	{
        	"start": "001001",
        	"end": "003003"
      	}
    	]
      	}`)

	if filter.isMatchedChfSupportedPlmn("002018", nfProfile) {
		t.Fatalf("should not be matched , but matched")
	}

	nfProfile = []byte(`{
	"plmnRangeList": [
      	{
        	"start": "00100",
        	"end": "003999"
      	}
    	]
      	}`)

	if !filter.isMatchedChfSupportedPlmn("00201", nfProfile) {
		t.Fatalf("should be matched , but not matched")
	}

	if !filter.isMatchedChfSupportedPlmn("002001", nfProfile) {
		t.Fatalf("should be matched , but not matched")
	}
}
