package cm

import (
	"fmt"
	"os"
	"strings"

	"sync/atomic"
	"unsafe"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

const (
	defaultIngressIP   = "https://127.0.0.1:443"
	defaultIngressFqdn = "https://nrf.ericsson.se:443"
)

var (
	// ServiceVersion is the version of nf service, got from helm chart version.
	ServiceVersion string

	//NRFNFServices to store mgmt&disc services
	NRFNFServices *TNRFNFServices
)

//TNRFNFServices to store mgmt&disc services
type TNRFNFServices struct {
	// ManagementNfServices is nf service of management
	ManagementNfServices []TNfService
	// DiscoveryNfServices is nf service of discovery
	DiscoveryNfServices []TNfService
	// ServiceVersion is the version of nf service, got from helm chart version.
}

func (conf *TNRFNFServices) atomicSetNRFNFServices() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&NRFNFServices)), unsafe.Pointer(conf))
}

// GetNRFNFServices is to get nrf nfservcies
func GetNRFNFServices() *TNRFNFServices {
	return (*TNRFNFServices)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&NRFNFServices))))
}

func (conf *TNRFNFServices) init() {
	conf.DiscoveryNfServices = make([]TNfService, 0)
	conf.ManagementNfServices = make([]TNfService, 0)
}

// SetServiceVersion is to set the service version
func SetServiceVersion() {
	ServiceVersion = os.Getenv("SERVICE_VERSION")
}

// GetMgmtIngressAddress return the Ingress Address of NRF-management, which is configured in CM
func GetMgmtIngressAddress() string {

	IngressAddress := constructIngressAddress(GetNRFNFServices().ManagementNfServices)

	if IngressAddress == defaultIngressFqdn || IngressAddress == defaultIngressIP {
		log.Debugf("The value of ingress-address is still default one, please configure it in CM !")
	}
	return IngressAddress
}

// GetDiscIngressAddress return the Ingress Address of NRF-discovery, which is configured in CM
func GetDiscIngressAddress() string {

	IngressAddress := constructIngressAddress(GetNRFNFServices().DiscoveryNfServices)

	if IngressAddress == defaultIngressFqdn || IngressAddress == defaultIngressIP {
		log.Debugf("The value of ingress-address is still default one, please configure it in CM !")
	}
	return IngressAddress
}

func constructIngressAddress(NfServices []TNfService) string {
	IPStackMode := os.Getenv("IP_STACK_MODE")
	scheme := "http"
	port := 0
	for _, NfService := range NfServices {
		scheme = NfService.Scheme
		for _, ipEndpoint := range NfService.IPEndpoint {
			port = ipEndpoint.Port
			ipAddress := ""
			if IPStackMode == constvalue.IPStackv4 && ipEndpoint.Ipv4Address != "" {
				ipAddress = ipEndpoint.Ipv4Address
			} else if IPStackMode == constvalue.IPStackv6 && ipEndpoint.Ipv6Address != "" {
				ipAddress = "[" + ipEndpoint.Ipv6Address + "]"
			}

			if scheme != "" && ipAddress != "" {
				ingressIPURI := constructURI(scheme, ipAddress, port)
				if !isDefaultIngressIP(ingressIPURI) {
					return ingressIPURI
				}
			}
		}

		if scheme != "" && NfService.Fqdn != "" {
			ingressFqdn := constructURI(scheme, NfService.Fqdn, port)
			if !isDefaultIngressFqdn(ingressFqdn) {
				return ingressFqdn
			}
		}
	}

	nfProfile := GetNRFNFProfile()
	if IPStackMode == constvalue.IPStackv4 && scheme != "" {
		for _, ipv4 := range nfProfile.Ipv4Address {
			ingressIPURI := constructURI(scheme, ipv4, port)
			if !isDefaultIngressIP(ingressIPURI) {
				return ingressIPURI
			}
		}
	}

	if IPStackMode == constvalue.IPStackv6 && scheme != "" {
		for _, ipv6 := range nfProfile.Ipv6Address {
			ipv6 = "[" + ipv6 + "]"
			ingressIPURI := constructURI(scheme, ipv6, port)
			if !isDefaultIngressIP(ingressIPURI) {
				return ingressIPURI
			}
		}
	}

	if scheme != "" {
		if nfProfile.Fqdn != "" {
			ingressFqdn := constructURI(scheme, nfProfile.Fqdn, port)
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
