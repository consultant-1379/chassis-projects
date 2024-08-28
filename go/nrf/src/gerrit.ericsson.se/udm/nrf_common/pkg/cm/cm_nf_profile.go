package cm

import (
	"fmt"
	"strings"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

var (
	// NfProfile is configuration of nf profile
	NfProfile TNfProfile
)

// ParseConf is to parse nf profile
func (conf *TNfProfile) ParseConf() {
	NfProfile = *conf
	NfProfile.toUpper()

	var ManagementNfServicesTmp []TNfService
	var DiscoveryNfServicesTmp []TNfService

	for _, service := range NfProfile.Service {
		if service.Name == constvalue.NNRFNFM {
			ManagementNfServicesTmp = append(ManagementNfServicesTmp, service)
		} else if service.Name == constvalue.NNRFDISC {
			DiscoveryNfServicesTmp = append(DiscoveryNfServicesTmp, service)
		}
	}

	ManagementNfServices = ManagementNfServicesTmp
	DiscoveryNfServices = DiscoveryNfServicesTmp

}

// toUpper UPPER some 3gpp value, e.g. NF type and NF status
func (conf *TNfProfile) toUpper() {
	conf.Type = strings.ToUpper(conf.Type)
	conf.Status = strings.ToUpper(conf.Status)

	for index := range conf.AllowedNfType {
		conf.AllowedNfType[index] = strings.ToUpper(conf.AllowedNfType[index])
	}

	for index := range conf.Service {
		conf.Service[index].Status = strings.ToUpper(conf.Service[index].Status)

		for subIndex := range conf.Service[index].IPEndpoint {
			conf.Service[index].IPEndpoint[subIndex].Transport = strings.ToUpper(conf.Service[index].IPEndpoint[subIndex].Transport)
		}

		for subIndex := range conf.Service[index].AllowedNfType {
			conf.Service[index].AllowedNfType[subIndex] = strings.ToUpper(conf.Service[index].AllowedNfType[subIndex])
		}

		for subIndex := range conf.Service[index].DefaultNotificationSubscription {
			conf.Service[index].DefaultNotificationSubscription[subIndex].NotificationType = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].NotificationType)
			conf.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass)
			conf.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass)
		}
	}
}

// Show is to print nf profile info
func (conf *TNfProfile) Show() {
	for index, plmn := range NfProfile.PlmnID {
		fmt.Printf("Plmnlist[%d] mcc : %s\n", index, plmn.Mcc)
		fmt.Printf("Plmnlist[%d] mnc : %s\n", index, plmn.Mnc)
	}
}
