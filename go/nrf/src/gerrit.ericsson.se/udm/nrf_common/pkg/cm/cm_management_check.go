package cm

import (
	//"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// CheckConfForMgmt is to check basic configuration of CM  for nrf-management serivice
func CheckConfForMgmt() bool {

	if !checkNrfCommon(constvalue.NNRFNFM) {
		return false
	}

	if !checkNfProfile() {
		return false
	}

	if !checkManagementNrfService() {
		return false
	}

	return true
}

func checkNrfCommon(serverName string) bool {

	if NrfCommon.Role != "plmn-nrf" && NrfCommon.Role != "region-nrf" && NrfCommon.Role != "slice-nrf" {
		log.Errorf("common.role in CM need to be configured correctly,current value:%s", NrfCommon.Role)
		return false
	}

	if NrfCommon.Role == "region-nrf" {

		if len(NrfCommon.PlmnNrf.Profile) == 0 {
			log.Errorf("common.plmn-nrf.profile in CM need to be configured")
			return false
		}

		for _, profile := range NrfCommon.PlmnNrf.Profile {

			nrfLevelPoint := false
			if profile.Fqdn != "" || len(profile.Ipv4Address) > 0 || len(profile.Ipv6Address) > 0 {
				nrfLevelPoint = true
			}

			if len(profile.Service) == 0 {
				log.Errorf("common.plmn-nrf.profile.service in CM need to be configured")
				return false
			}

			for _, service := range profile.Service {

				if serverName != service.Name {
					continue
				}

				if service.Scheme != "http" && service.Scheme != "https" {
					log.Errorf("common.plmn-nrf.profile.service.scheme in CM need to be configured correctly,current value:%s", service.Scheme)
					return false
				}

				if !nrfLevelPoint {
					if service.Fqdn == "" && !checkIPEndPoint(service.IPEndpoint) {
						log.Errorf("common.plmn-nrf.profile.service.fqdn and ip-endpoint in CM is at least one occurrence")
						return false
					}
				}
			}
		}
	}

	return true
}

func checkNfProfile() bool {

	if len(NfProfile.PlmnID) == 0 {
		log.Errorf("nf-profile.plmn-id in CM need to be configured correctly, current is empty")
		return false
	}

	if NfProfile.InstanceID == "" {
		log.Errorf("nf-profile.instance-id in CM need to be configured,it can't be empty")
		return false
	}

	return true
}

func checkManagementNrfService() bool {

	nfProfilePoint := false
	if NfProfile.Fqdn != "" || len(NfProfile.Ipv4Address) > 0 || len(NfProfile.Ipv6Address) > 0 {
		nfProfilePoint = true
	}

	if len(ManagementNfServices) == 0 {
		log.Errorf("at least one management service in CM need to be configured")
		return false
	}

	for _, managementNfService := range ManagementNfServices {
		if managementNfService.Scheme != "http" && managementNfService.Scheme != "https" {
			log.Errorf("management service scheme in CM need to be configured correctly,current value:%s", managementNfService.Scheme)
			return false
		}

		if managementNfService.InstanceID == "" {
			log.Errorf("management service instance-id in CM need to be configured correctly,current value:%s", managementNfService.InstanceID)
			return false
		}

		if !nfProfilePoint {
			if managementNfService.Fqdn == "" && !checkIPEndPoint(managementNfService.IPEndpoint) {
				log.Errorf("management service fqdn and ip-endpoint in CM is at least one occurrence")
				return false
			}
		}
	}

	return true
}

func checkDiscoveryNrfService() bool {

	nfProfilePoint := false
	if NfProfile.Fqdn != "" || len(NfProfile.Ipv4Address) > 0 || len(NfProfile.Ipv6Address) > 0 {
		nfProfilePoint = true
	}

	if len(DiscoveryNfServices) == 0 {
		log.Errorf("at least one discovery service in CM need to be configured")
		return false
	}

	for _, discoveryNfService := range DiscoveryNfServices {
		if discoveryNfService.Scheme != "http" && discoveryNfService.Scheme != "https" {
			log.Errorf("discovery service scheme in CM need to be configured correctly,current value:%s", discoveryNfService.Scheme)
			return false
		}

		if discoveryNfService.InstanceID == "" {
			log.Errorf("discovery service instance-id in CM need to be configured correctly,current value:%s", discoveryNfService.InstanceID)
			return false
		}

		if !nfProfilePoint {
			if discoveryNfService.Fqdn == "" && !checkIPEndPoint(discoveryNfService.IPEndpoint) {
				log.Errorf("discovery service fqdn and ip-endpoint in CM is at least one occurrence")
				return false
			}
		}
	}
	return true
}

func checkIPEndPoint(IPEndpoints []TIPEndpoint) bool {

	if len(IPEndpoints) == 0 {
		return false
	}

	for _, vIPEndpoint := range IPEndpoints {
		if vIPEndpoint.Ipv4Address == "" && vIPEndpoint.Ipv6Address == "" {
			log.Errorf("ip-endpoint.ipv4-address and ip-endpoint.ipv6-address in CM is at least one occurrence")
			return false
		}
	}
	return true
}
