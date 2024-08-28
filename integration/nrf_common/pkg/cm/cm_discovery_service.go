package cm

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

const (
	defaultValidityPeriodOfSearchResult = 86400
	defaultLocalCacheTimeout            = 3000
	defaultLocalCacheCapacity           = 100
	defaultOverloadRedirectionURL       = "http://www.example.com"
)

var (
	// DiscoveryService is configuration of discovery service
	DiscoveryService *TDiscoveryService
	// ValidityPeriodOfSearchResult is discovery cache time
	ValidityPeriodOfSearchResult = defaultValidityPeriodOfSearchResult
	//DiscLocalCacheEnable is whether local cache enable
	DiscLocalCacheEnable = false
	//DiscLocalCacheTimeout is local cache timeout
	DiscLocalCacheTimeout = defaultLocalCacheTimeout
	//DiscLocalCacheCapacity is local cache capacity
	DiscLocalCacheCapacity = defaultLocalCacheCapacity

	// OverloadRedirectionEnabled is whether overload redirection feature enabled
	OverloadRedirectionEnabled = false
)

func (conf *TDiscoveryService) atomicSetDiscService() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&DiscoveryService)), unsafe.Pointer(conf))
}

// GetDiscService to get disc service
func GetDiscService() *TDiscoveryService {
	return (*TDiscoveryService)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&DiscoveryService))))
}

// ParseConf is to parse configuration of discovery service
func (conf *TDiscoveryService) ParseConf() {

	ValidityPeriodOfSearchResult = conf.ResponseCacheTime
	if ValidityPeriodOfSearchResult <= 0 {
		//NrfConfig.DiscoveryServiceProfile.DiscoveryCacheTime = defaultValidityPeriodOfSearchResult
		ValidityPeriodOfSearchResult = defaultValidityPeriodOfSearchResult
		conf.ResponseCacheTime = defaultValidityPeriodOfSearchResult
	}

	if conf.LocalCache == nil {
		conf.LocalCache = &TLocalCache{
			Enabled:  false,
			Timeout:  defaultLocalCacheTimeout,
			Capacity: defaultLocalCacheCapacity,
		}
	}

	DiscLocalCacheEnable = conf.LocalCache.Enabled

	DiscLocalCacheTimeout = conf.LocalCache.Timeout
	if DiscLocalCacheTimeout <= 0 {
		DiscLocalCacheTimeout = defaultLocalCacheTimeout
		conf.LocalCache.Timeout = defaultLocalCacheTimeout
	}

	DiscLocalCacheCapacity = conf.LocalCache.Capacity
	if DiscLocalCacheCapacity <= 0 || DiscLocalCacheCapacity >= 20000 {
		DiscLocalCacheCapacity = defaultLocalCacheCapacity
		conf.LocalCache.Capacity = defaultLocalCacheCapacity
	}

	if conf.OverloadRedirection == nil {
		conf.OverloadRedirection = &TOverloadRedirection{
			Enabled:        false,
			Mode:           "auto",
			RedirectionURL: nil,
		}
	}

	OverloadRedirectionEnabled = conf.OverloadRedirection.Enabled
	if conf.OverloadRedirection.Mode == "" || conf.OverloadRedirection.Mode == "auto" {
		OverloadRedirectionEnabled = false
	}

	conf.atomicSetDiscService()
}

// SetDefaultForDiscoveryService to set default value when it is not configured
func SetDefaultForDiscoveryService() {
	var discService TDiscoveryService

	discService.ResponseCacheTime = defaultValidityPeriodOfSearchResult
	discService.LocalCache = &TLocalCache{
		Enabled:  false,
		Timeout:  defaultLocalCacheTimeout,
		Capacity: defaultLocalCacheCapacity,
	}

	discService.atomicSetDiscService()
}

// Show print discovery service profile
func (conf *TDiscoveryService) Show() {
	fmt.Printf("discovery-response-cache-time : %d\n", defaultValidityPeriodOfSearchResult)
}

// GetOverloadRedirectionURL return the overload redirection URL
func GetOverloadRedirectionURL() string {
	discService := GetDiscService()
	if discService != nil && discService.OverloadRedirection != nil && len(discService.OverloadRedirection.RedirectionURL) > 0 {
		return discService.OverloadRedirection.RedirectionURL[0]
	}

	return defaultOverloadRedirectionURL
}
