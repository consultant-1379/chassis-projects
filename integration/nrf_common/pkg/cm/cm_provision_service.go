package cm

import (
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"unsafe"
	"sync/atomic"
)

var (
	// ProvisionService is configuration of provision service profile
	ProvisionService *TProvisionService
	ingressAddress   *TProvIngressAddress
)

//TProvIngressAddress provisionservcie's ingressaddress
type TProvIngressAddress struct {
	ingressAddress []string
}

func (p *TProvIngressAddress)init(){
	p.ingressAddress = make([]string, 0)
}

func (p *TProvIngressAddress)atomicSetIngressAddress(){
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ingressAddress)), unsafe.Pointer(p))
}
//SetPrivisionService use atomic to set provisionservice
func (conf *TProvisionService)atomicSetPrivisionService() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ProvisionService)), unsafe.Pointer(conf))
}

//GetPrivisionService to get provisionservice by atomic
func GetPrivisionService()*TProvisionService{
	return (*TProvisionService)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ProvisionService))))
}

// ParseConf is to parse configuration of provision service profile
func (conf *TProvisionService) ParseConf() {
	conf.atomicSetPrivisionService()
	ProvisionService.constructIngressAddress()
}

func (conf *TProvisionService) constructIngressAddress() {
	var ingressAddressTmp []string

	for _, item := range GetPrivisionService().ProvAddress {

		if item.ID < 0 {
			log.Warning("constructIngressAddress: ProvisionServProfile Addresses ID is empty, Please correct it")
			continue
		}

		if item.Scheme == "" {
			log.Warning("constructIngressAddress: ProvisionServProfile Addresses Scheme is empty, Please correct it")
			continue
		}

		port := item.Port
		if port == 0 {
			if item.Scheme == "https" {
				port = 443
			} else if item.Scheme == "http" {
				port = 80
			}
		}

		var address string

		if IPStackMode == constvalue.IPStackv4 && item.Ipv4Address != "" {
			address = item.Ipv4Address
		} else if IPStackMode == constvalue.IPStackv6 && item.Ipv6Address != "" {
			address = "[" + item.Ipv6Address + "]"
		} else {
			address = item.Fqdn
		}

		if address == "" {
			log.Warning("constructIngressAddress: ProvisionServProfile Addresses address is empty, Please correct it")
			continue
		}

		ingressAddressTmp = append(ingressAddressTmp, fmt.Sprintf("%s://%s:%d", item.Scheme, address, port))
	}

	var ingress TProvIngressAddress
	ingress.init()
	ingress.ingressAddress = ingressAddressTmp
	ingress.atomicSetIngressAddress()
}

// Show print discovery service profile
func (conf *TProvisionService) Show() {
	fmt.Printf("TProvisionService value is : %+v\n", GetPrivisionService())
}

// CheckProvisionService to check if provision service profile is configured
func CheckProvisionService() bool {

	if len(GetPrivisionService().ProvAddress) == 0 {
		log.Errorf("ProvisionService.ProvAddress in CM need to be configured")
		return false
	}

	for _, provAddress := range GetPrivisionService().ProvAddress {
		if provAddress.Scheme != "http" && provAddress.Scheme != "https" {
			log.Errorf("provision service scheme in CM need to be configured correctly,current value:%s", provAddress.Scheme)
			return false
		}

		if IPStackMode == constvalue.IPStackv4 {
			if provAddress.Fqdn == "" && provAddress.Ipv4Address == "" {
				log.Errorf("provision service fqdn and ipv4-address in CM is at least one occurrence")
				return false
			}
		} else if IPStackMode == constvalue.IPStackv6 {
			if provAddress.Fqdn == "" && provAddress.Ipv6Address == "" {
				log.Errorf("provision service fqdn and ipv6-address in CM is at least one occurrence")
				return false
			}
		}
	}
	return true
}

// GetProvIngressAddress for getting the location address prefix
func GetProvIngressAddress() string {
	ingress := (*TProvIngressAddress)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ingressAddress))))
	if len(ingress.ingressAddress) < 1 {
		return ""
	}

	return ingress.ingressAddress[0]
}
