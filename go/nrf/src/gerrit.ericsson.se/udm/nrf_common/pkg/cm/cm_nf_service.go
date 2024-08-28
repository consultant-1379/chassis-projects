package cm

import (
	"fmt"
	"os"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

const (
	defaultIngressIP   = "https://127.0.0.1:443"
	defaultIngressFqdn = "https://nrf.ericsson.se:443"
)

var (
	// ManagementNfServices is nf service of management
	ManagementNfServices []TNfService
	// DiscoveryNfServices is nf service of discovery
	DiscoveryNfServices []TNfService
	// ServiceVersion is the version of nf service, got from helm chart version.
	ServiceVersion string
)

// SetServiceVersion is to set the service version
func SetServiceVersion() {
	ServiceVersion = os.Getenv("SERVICE_VERSION")
}

// GetMgmtIngressAddress return the Ingress Address of NRF-management, which is configured in CM
func GetMgmtIngressAddress() string {

	IngressAddress := constructIngressAddress(ManagementNfServices)

	if IngressAddress == defaultIngressFqdn || IngressAddress == defaultIngressIP {
		log.Debugf("The value of ingress-address is still default one, please configure it in CM !")
	}
	return IngressAddress
}

// GetDiscIngressAddress return the Ingress Address of NRF-discovery, which is configured in CM
func GetDiscIngressAddress() string {

	IngressAddress := constructIngressAddress(DiscoveryNfServices)

	if IngressAddress == defaultIngressFqdn || IngressAddress == defaultIngressIP {
		log.Debugf("The value of ingress-address is still default one, please configure it in CM !")
	}
	return IngressAddress
}

func constructIngressAddress(NfServices []TNfService) string {
	port := 0
	for _, NfService := range NfServices {

		for _, ipEndpoint := range NfService.IPEndpoint {
			port = ipEndpoint.Port
			ipAddress := ""
			if ipEndpoint.Ipv4Address != "" {
				ipAddress = ipEndpoint.Ipv4Address
			} else if ipEndpoint.Ipv4Address == "" && ipEndpoint.Ipv6Address != "" {
				ipAddress = "[" + ipEndpoint.Ipv6Address + "]"
			}

			if NfService.Scheme != "" && ipAddress != "" {
				ingressIPURI := constructURI(NfService.Scheme, ipEndpoint.Ipv4Address, port)
				if !isDefaultIngressIP(ingressIPURI) {
					return ingressIPURI
				}
			}
		}

		if NfService.Scheme != "" && NfService.Fqdn != "" {
			ingressFqdn := constructURI(NfService.Scheme, NfService.Fqdn, port)
			if !isDefaultIngressFqdn(ingressFqdn) {
				return ingressFqdn
			}
		}
	}

	return ""
}

func isDefaultIngressFqdn(fqdn string) bool {
	flag := false
	if fqdn == defaultIngressFqdn {
		flag = true
	}
	return flag
}

func isDefaultIngressIP(ip string) bool {
	flag := false
	if ip == defaultIngressIP {
		flag = true
	}
	return flag
}

func constructURI(scheme string, host string, port int) string {
	if port == 0 {
		if "https" == strings.ToLower(scheme) {
			port = 443
		} else if "http" == strings.ToLower(scheme) {
			port = 80

		} else {
			log.Errorf("NfService.Scheme is incorrect value:%s!", scheme)
		}
	}

	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

// Show is to print nf service info
func (*TNfService) Show() {

}
