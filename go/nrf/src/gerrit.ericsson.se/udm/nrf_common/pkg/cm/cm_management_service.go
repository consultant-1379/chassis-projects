package cm

import (
	"fmt"
	"strings"
)

const (
	defaultHeartbeatTimer               = 120
	defaultHeartbeatTimerGracePeriod    = 5
	defaultValidityPeriodOfSubscription = 604800
)

var (
	// ValidityPeriodOfSubscription records the expired time of subscription
	ValidityPeriodOfSubscription = defaultValidityPeriodOfSubscription
	// ManagementService is configuration of management service
	ManagementService TManagementService
)

// ParseConf is to parse management service
func (conf *TManagementService) ParseConf() {
	ManagementService = *conf

	if conf.Heartbeat == nil {
		ManagementService.Heartbeat = &THeartbeat{Default: defaultHeartbeatTimer, GracePeriod: defaultHeartbeatTimerGracePeriod}
	}
	ManagementService.toUpper()
	ValidityPeriodOfSubscription = conf.SubscriptionExpiredTime

	if ValidityPeriodOfSubscription <= 0 {
		ManagementService.SubscriptionExpiredTime = defaultValidityPeriodOfSubscription
		ValidityPeriodOfSubscription = defaultValidityPeriodOfSubscription
	}
}

// toUpper UPPER some 3gpp value, e.g. NF type
func (conf *TManagementService) toUpper() {
	if conf.Heartbeat != nil {
		for index := range conf.Heartbeat.DefaultPerNfType {
			conf.Heartbeat.DefaultPerNfType[index].NfType = strings.ToUpper(conf.Heartbeat.DefaultPerNfType[index].NfType)
		}
	}
}

// Show management service profile
func (conf *TManagementService) Show() {
	fmt.Printf("The default heartbeat is %d\n", ManagementService.Heartbeat.Default)
	for _, defaultPerNftype := range ManagementService.Heartbeat.DefaultPerNfType {
		fmt.Printf("The default heartbeat for NFType %s is %d\n", defaultPerNftype.NfType, defaultPerNftype.Value)
	}
	fmt.Printf("heartbeat grace period is : %d\n", ManagementService.Heartbeat.GracePeriod)
	fmt.Printf("subscription-expired-time : %d\n", ValidityPeriodOfSubscription)
}
