package cm

import (
	"testing"
)

func TestSetDefaultForDiscoveryService(t *testing.T) {

	SetDefaultForDiscoveryService()
	discService := GetDiscService()
	if discService.ResponseCacheTime != defaultValidityPeriodOfSearchResult {
		t.Fatal("default value of DiscoveryService.ResponseCacheTime is incorrect !")
	}

	if discService.LocalCache.Enabled {
		t.Fatal("default value of discService.LocalCache.Enabled is incorrect !")
	}

	if discService.LocalCache.Capacity != defaultLocalCacheCapacity {
		t.Fatal("default value of discService.LocalCache.Capacity is incorrect !")
	}

	if discService.LocalCache.Timeout != defaultLocalCacheTimeout {
		t.Fatal("default value of discService.LocalCache.Timeout is incorrect !")
	}

}

func TestDiscoveryParse(t *testing.T) {

	// case1: all the values are appropriate
	discoveryProfileIns := &TDiscoveryService{
		ResponseCacheTime: 20000,
		LocalCache: &TLocalCache{
			Enabled:  true,
			Timeout:  100,
			Capacity: 200,
		},
		OverloadRedirection: &TOverloadRedirection{
			Enabled: true,
			Mode:    "manual",
			RedirectionURL: []string{
				"http://www.example1.com",
			},
		},
	}

	discoveryProfileIns.ParseConf()

	if ValidityPeriodOfSearchResult != 20000 {
		t.Fatal("ValidityPeriodOfSearchResult should be 20000, but not !")
	}

	if !DiscLocalCacheEnable {
		t.Fatal("DiscLocalCacheEnable should be true, but not !")
	}

	if DiscLocalCacheTimeout != 100 {
		t.Fatal("DiscLocalCacheTimeout should be 200, but not !")
	}

	if DiscLocalCacheCapacity != 200 {
		t.Fatal("DiscLocalCacheCapacity should be 200, but not !")
	}

	if !OverloadRedirectionEnabled {
		t.Fatal("OverloadRedirectionEnabled should be true, but not !")
	}

	if GetOverloadRedirectionURL() != "http://www.example1.com" {
		t.Fatal("GetOverloadRedirectionURL() didn't return right value !")
	}

	// case2: all the values are not appropriate, use default
	discoveryProfileIns = &TDiscoveryService{
		ResponseCacheTime: -20000,
		LocalCache: &TLocalCache{
			Enabled:  true,
			Timeout:  -100,
			Capacity: -200,
		},
		OverloadRedirection: &TOverloadRedirection{
			Enabled:        false,
			Mode:           "auto",
			RedirectionURL: []string{},
		},
	}

	discoveryProfileIns.ParseConf()

	if ValidityPeriodOfSearchResult != defaultValidityPeriodOfSearchResult {
		t.Fatal("ValidityPeriodOfSearchResult should be 86400, but not !")
	}

	if !DiscLocalCacheEnable {
		t.Fatal("DiscLocalCacheEnable should be true, but not !")
	}

	if DiscLocalCacheTimeout != defaultLocalCacheTimeout {
		t.Fatal("DiscLocalCacheTimeout should be 3000, but not !")
	}

	if DiscLocalCacheCapacity != defaultLocalCacheCapacity {
		t.Fatal("DiscLocalCacheCapacity should be 100, but not !")
	}

	if OverloadRedirectionEnabled {
		t.Fatal("OverloadRedirectionEnabled should be false, but not !")
	}

	if GetOverloadRedirectionURL() != "http://www.example.com" {
		t.Fatal("GetOverloadRedirectionURL() didn't return right value !")
	}
}
