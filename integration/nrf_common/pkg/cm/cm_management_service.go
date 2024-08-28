package cm

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

const (
	defaultHeartbeatTimer               = 120
	defaultHeartbeatTimerGracePeriod    = 5
	defaultValidityPeriodOfSubscription = 604800
	defautMinHeartbeat                  = 5
	defaultTrafficRateLimitPerInstance  = 10
)

var (
	// ValidityPeriodOfSubscription records the expired time of subscription
	ValidityPeriodOfSubscription = defaultValidityPeriodOfSubscription

	// TrafficRateLimitPerInstance records the traffic rate limit of per nf instance
	TrafficRateLimitPerInstance = defaultTrafficRateLimitPerInstance

	// ManagementService is configuration of management service
	ManagementService *TManagementService
)

func (conf *TManagementService) atomicSetMgmtService() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ManagementService)), unsafe.Pointer(conf))
}

//GetMgmtService to get mgmt service
func GetMgmtService() *TManagementService {
	return (*TManagementService)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ManagementService))))
}

// ParseConf is to parse management service
func (conf *TManagementService) ParseConf() {

	if conf.Heartbeat == nil {
		conf.Heartbeat = &THeartbeat{Default: defaultHeartbeatTimer, GracePeriod: defaultHeartbeatTimerGracePeriod}
	} else {

		if conf.Heartbeat.Default < defautMinHeartbeat {
			conf.Heartbeat.Default = defaultHeartbeatTimer
		}

		if conf.Heartbeat.GracePeriod <= 0 {
			conf.Heartbeat.GracePeriod = defaultHeartbeatTimerGracePeriod
		}
	}

	conf.toUpper()

	if conf.SubscriptionExpiredTime <= 0 {
		conf.SubscriptionExpiredTime = defaultValidityPeriodOfSubscription

	}
	ValidityPeriodOfSubscription = conf.SubscriptionExpiredTime

	if conf.TrafficRateLimitPerInstance <= 0 {
		conf.TrafficRateLimitPerInstance = defaultTrafficRateLimitPerInstance
	}
	TrafficRateLimitPerInstance = conf.TrafficRateLimitPerInstance

	conf.atomicSetMgmtService()
}

// toUpper UPPER some 3gpp value, e.g. NF type
func (conf *TManagementService) toUpper() {
	if conf.Heartbeat != nil {
		for index := range conf.Heartbeat.DefaultPerNfType {
			conf.Heartbeat.DefaultPerNfType[index].NfType = strings.ToUpper(conf.Heartbeat.DefaultPerNfType[index].NfType)
		}
	}
}

// SetDefaultForManagementService to set default value when it is not configured
func SetDefaultForManagementService() {
	var mgmtService TManagementService
	mgmtService.Heartbeat = &THeartbeat{Default: defaultHeartbeatTimer, GracePeriod: defaultHeartbeatTimerGracePeriod}
	mgmtService.TrafficRateLimitPerInstance = defaultTrafficRateLimitPerInstance
	mgmtService.SubscriptionExpiredTime = defaultValidityPeriodOfSubscription
	mgmtService.atomicSetMgmtService()
	ValidityPeriodOfSubscription = defaultValidityPeriodOfSubscription
}

// Show management service profile
func (conf *TManagementService) Show() {
	fmt.Printf("The default heartbeat is %d\n", GetMgmtService().Heartbeat.Default)
	for _, defaultPerNftype := range GetMgmtService().Heartbeat.DefaultPerNfType {
		fmt.Printf("The default heartbeat for NFType %s is %d\n", defaultPerNftype.NfType, defaultPerNftype.Value)
	}
	fmt.Printf("heartbeat grace period is : %d\n", GetMgmtService().Heartbeat.GracePeriod)
	fmt.Printf("subscription-expired-time : %d\n", ValidityPeriodOfSubscription)
}
