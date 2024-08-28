package cm

import (
	"fmt"
	"strings"
)

var (
	// NrfPolicy is configuration of nrf policy
	NrfPolicy TNrfPolicy
)

// ParseConf is to parse nf profile
func (conf *TNrfPolicy) ParseConf() {
	NrfPolicy = *conf
	NrfPolicy.toUpper()

}

// toUpper UPPER some 3gpp value, e.g. NF type
func (conf *TNrfPolicy) toUpper() {
	if conf.ManagementService != nil && conf.ManagementService.Subscription != nil {
		for index := range conf.ManagementService.Subscription.AllowedSubscriptionAllNFs {
			conf.ManagementService.Subscription.AllowedSubscriptionAllNFs[index].AllowedNfType = strings.ToUpper(conf.ManagementService.Subscription.AllowedSubscriptionAllNFs[index].AllowedNfType)
		}
	}
}

// Show is to print nf profile info
func (conf *TNrfPolicy) Show() {
	allowedSubToAllNFs := conf.ManagementService.Subscription.AllowedSubscriptionAllNFs
	if len(allowedSubToAllNFs) > 0 {
		for index := range allowedSubToAllNFs {
			fmt.Printf("allowedSubToAllNFs[%d].allowedNFType is %s\n", index, allowedSubToAllNFs[index].AllowedNfType)
			fmt.Printf("allowedSubToAllNFs[%d].AllowedNfDomains is %s\n", index, allowedSubToAllNFs[index].AllowedNfDomains)
		}
	}
}
