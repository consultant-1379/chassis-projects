package nfdiscservice

import (
	"testing"
)

func TestNfSupportAcrossRegion(t *testing.T) {
	if !nfSupportAcrossRegion("AUSF") {
		t.Fatalf("AUSF should support across region, but failed")
	}
	if nfSupportAcrossRegion("NRF") {
		t.Fatalf("NRF should not support accross region, but support")
	}
}

