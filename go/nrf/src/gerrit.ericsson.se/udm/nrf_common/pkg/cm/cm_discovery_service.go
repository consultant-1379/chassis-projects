package cm

import (
	"fmt"
)

const (
	defaultValidityPeriodOfSearchResult = 86400
	defaultLocalCacheTimeout            = 3000
	defaultLocalCacheCapacity           = 100
)

var (
	// ValidityPeriodOfSearchResult is discovery cache time
	ValidityPeriodOfSearchResult = defaultValidityPeriodOfSearchResult
	// DiscoveryService is configuration of discovery service
	DiscoveryService TDiscoveryService
	//DiscLocalCacheEnable is whether local cache enable
	DiscLocalCacheEnable = false
	//DiscLocalCacheTimeout is local cache timeout
	DiscLocalCacheTimeout = defaultLocalCacheTimeout
	//DiscLocalCacheCapacity is local cache capacity
	DiscLocalCacheCapacity = defaultLocalCacheCapacity
)

// ParseConf is to parse configuration of discovery service
func (conf *TDiscoveryService) ParseConf() {

	DiscoveryService = *conf
	ValidityPeriodOfSearchResult = conf.ResponseCacheTime
	if ValidityPeriodOfSearchResult <= 0 {
		//NrfConfig.DiscoveryServiceProfile.DiscoveryCacheTime = defaultValidityPeriodOfSearchResult
		ValidityPeriodOfSearchResult = defaultValidityPeriodOfSearchResult
		DiscoveryService.ResponseCacheTime = defaultValidityPeriodOfSearchResult
	}

	DiscLocalCacheEnable = conf.LocalCacheEnable

	DiscLocalCacheTimeout = conf.LocalCacheTimeout
	if DiscLocalCacheTimeout <= 0 {
		DiscLocalCacheTimeout = defaultLocalCacheTimeout
		DiscoveryService.LocalCacheTimeout = defaultLocalCacheTimeout
	}

	DiscLocalCacheCapacity = conf.LocalCacheCapacity
	if DiscLocalCacheCapacity <= 0 || DiscLocalCacheCapacity >= 20000{
		DiscLocalCacheCapacity = defaultLocalCacheCapacity
		DiscoveryService.LocalCacheCapacity = defaultLocalCacheCapacity
	}
}

// Show print discovery service profile
func (conf *TDiscoveryService) Show() {
	fmt.Printf("discovery-response-cache-time : %d\n", defaultValidityPeriodOfSearchResult)
}
