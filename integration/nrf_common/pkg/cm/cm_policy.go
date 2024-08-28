package cm

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

const (
	defaultMediumStart     = 8
	defaultLowStart        = 16
	defaultMessagePriority = 24
)

var (
	// NrfPolicy is configuration of nrf policy
	NrfPolicy *TNrfPolicy
)

func (conf *TNrfPolicy) atomicSetNRFPolicy() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&NrfPolicy)), unsafe.Pointer(conf))
}

//GetNRFPolicy to get nrfpolicy in cm
func GetNRFPolicy() *TNrfPolicy {
	return (*TNrfPolicy)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&NrfPolicy))))
}

// ParseConf is to parse nf profile
func (conf *TNrfPolicy) ParseConf() {
	if conf.ManagementService == nil {
		subscription := &TSubscriptionPolicy{}
		managementService := &TNrfManagementServicePolicy{Subscription: subscription}
		conf.ManagementService = managementService
	}

	if conf.ManagementService.Subscription == nil {
		subscription := &TSubscriptionPolicy{}
		conf.ManagementService.Subscription = subscription
	}

	conf.toUpper()
	conf.atomicSetNRFPolicy()
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

// PriorityGroup is used to construct priority group
type PriorityGroup struct {
	Level int
	Start int
	End   int
}

func getDefaultPriorityGroup() []PriorityGroup {
	lowGroup := PriorityGroup{
		Level: 3,
		Start: 16,
		End:   31,
	}

	mediumGroup := PriorityGroup{
		Level: 2,
		Start: 8,
		End:   15,
	}

	highGroup := PriorityGroup{
		Level: 1,
		Start: 0,
		End:   7,
	}

	return []PriorityGroup{
		lowGroup,
		mediumGroup,
		highGroup,
	}
}

// GetPriorityGroup returns the message priority group
func GetPriorityGroup() []PriorityGroup {
	nrfPolicy := GetNRFPolicy()
	if nrfPolicy == nil {
		return getDefaultPriorityGroup()
	}

	priorityPolicy := nrfPolicy.MessagePriorityPolicy
	if priorityPolicy == nil {
		return getDefaultPriorityGroup()
	}

	var lowStart int
	var mediumStart int

	if priorityPolicy.LowStart == nil {
		lowStart = defaultLowStart
	} else {
		lowStart = *(priorityPolicy.LowStart)
	}

	if priorityPolicy.MediumStart == nil {
		mediumStart = defaultMediumStart
	} else {
		mediumStart = *(priorityPolicy.MediumStart)
	}

	if mediumStart <= 0 || mediumStart >= lowStart || lowStart > 31 {
		return getDefaultPriorityGroup()
	}

	lowGroup := PriorityGroup{
		Level: 3,
		Start: lowStart,
		End:   31,
	}

	mediumGroup := PriorityGroup{
		Level: 2,
		Start: mediumStart,
		End:   lowStart - 1,
	}

	highGroup := PriorityGroup{
		Level: 1,
		Start: 0,
		End:   mediumStart - 1,
	}

	return []PriorityGroup{
		lowGroup,
		mediumGroup,
		highGroup,
	}
}

// GetDefaultPriority returns the default message priority
func GetDefaultPriority() int {
	nrfPolicy := GetNRFPolicy()
	if nrfPolicy == nil || nrfPolicy.MessagePriorityPolicy == nil {
		return defaultMessagePriority
	}

	return *(nrfPolicy.MessagePriorityPolicy.DefaultPriority)
}
