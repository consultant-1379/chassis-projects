package cm

import (
	"testing"
)

func TestExpectedDiscoveryCacheTime(t *testing.T) {

	discoveryProfileIns := &TDiscoveryService{ResponseCacheTime: 20000}

	discoveryProfileIns.ParseConf()

	if ValidityPeriodOfSearchResult != 20000 {
		t.Fatal("ValidityPeriodOfSearchResult should be 20000, but not !")
	}
}

func TestUnExpectedDiscoveryCacheTime(t *testing.T) {
	discoveryProfileIns := &TDiscoveryService{ResponseCacheTime: -20000}
	discoveryProfileIns.ParseConf()

	if ValidityPeriodOfSearchResult != 86400 {
		t.Fatal("ValidityPeriodOfSearchResult should be 86400, but not !")
	}
}
