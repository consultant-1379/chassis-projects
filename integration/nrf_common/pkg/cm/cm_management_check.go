package cm

import (
	//"time"
	"os"

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

	if GetNRFRole() == REGIONNRF {
		if !checkDiscoveryNrfService() {
			return false
		}
	}

	return true
}

func checkNrfCommon(serverName string) bool {
	nrfCommon := GetNRFCommon()
	if GetNRFRole() == UNKNOWNROLE {
		log.Errorf("one of the role choice common.region-nrf or common.plmn-nrf shall be configured")
		return false
	}

	if GetNRFRole() == REGIONNRF {

		if nrfCommon.RegionNrf.NextHop == nil || len(nrfCommon.RegionNrf.NextHop.Site) == 0 {
			log.Debugf("There's no next-hop NRF configured")
			return true
		}

		for _, site := range nrfCommon.RegionNrf.NextHop.Site {
			if len(site.Profile) == 0 {
				log.Errorf("common.region-nrf.next-hop.site.profile in CM need to be configured for site %s", site.ID)
				return false
			}
			for _, profile := range site.Profile {
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
	}

	return true
}

func checkNfProfile() bool {

	if len(GetNRFNFProfile().PlmnID) == 0 {
		log.Errorf("nf-profile.plmn-id in CM need to be configured correctly, current is empty")
		return false
	}

	if GetNRFNFProfile().InstanceID == "" {
		log.Errorf("nf-profile.instance-id in CM need to be configured,it can't be empty")
		return false
	}

	return true
}

func checkManagementNrfService() bool {

	nfProfilePoint := false
	nfProfile := GetNRFNFProfile()
	if nfProfile.Fqdn != "" || len(nfProfile.Ipv4Address) > 0 || len(nfProfile.Ipv6Address) > 0 {
		nfProfilePoint = true
	}
	nrfService := GetNRFNFServices()
	if len(nrfService.ManagementNfServices) == 0 {
		log.Errorf("at least one management service in CM need to be configured")
		return false
	}

	for _, managementNfService := range nrfService.ManagementNfServices {
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
	nfProfile := GetNRFNFProfile()
	if nfProfile.Fqdn != "" || len(nfProfile.Ipv4Address) > 0 || len(nfProfile.Ipv6Address) > 0 {
		nfProfilePoint = true
	}

	if len(GetNRFNFServices().DiscoveryNfServices) == 0 {
		log.Errorf("at least one discovery service in CM need to be configured")
		return false
	}
	for _, discoveryNfService := range GetNRFNFServices().DiscoveryNfServices {
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

	IPStackMode := os.Getenv("IP_STACK_MODE")
	for _, vIPEndpoint := range IPEndpoints {

		if IPStackMode == constvalue.IPStackv4 && vIPEndpoint.Ipv4Address == "" {
			log.Errorf("ipv4-address of ip-endpoint[%d] in CM need to be configured", vIPEndpoint.ID)
			return false
		}

		if IPStackMode == constvalue.IPStackv6 && vIPEndpoint.Ipv6Address == "" {
			log.Errorf("ipv6-address of ip-endpoint[%d] in CM need to be configured", vIPEndpoint.ID)
			return false
		}
	}
	return true
}
