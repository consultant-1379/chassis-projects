package cm

import (
	"os"
	"testing"
)

func TestExpectedIngressAddress(t *testing.T) {
	os.Setenv("IP_STACK_MODE", "ipv4")
	nfService := TNfService{
		IPEndpoint: []TIPEndpoint{
			TIPEndpoint{
				Transport:   "tcp",
				Ipv4Address: "127.0.0.2",
				Port:        443,
			},
		},
		Fqdn:          "test.ericsson.se",
		Scheme:        "https",
		InterPlmnFqdn: "int.ericsson.se",
	}

	var nrfServce TNRFNFServices
	nrfServce.init()
	nrfServce.ManagementNfServices = append(nrfServce.ManagementNfServices, nfService)
	nrfServce.atomicSetNRFNFServices()
	if GetMgmtIngressAddress() != "https://127.0.0.2:443" {
		t.Fatal("The ingress address should be https://127.0.0.2:443, but not !")
	}

	if GetNRFNFServices().ManagementNfServices[0].InterPlmnFqdn != "int.ericsson.se" {
		t.Fatal("The InterPlmnFqdn got from conf is incorrect!")
	}

}

func TestExpectedIngressAddress2(t *testing.T) {

	nfService := TNfService{
		Fqdn:          "test.ericsson.se",
		Scheme:        "https",
		InterPlmnFqdn: "int.ericsson.se",
	}

	var nrfServce TNRFNFServices
	nrfServce.init()
	nrfServce.ManagementNfServices = append(nrfServce.ManagementNfServices, nfService)
	nrfServce.atomicSetNRFNFServices()
	if GetMgmtIngressAddress() != "https://test.ericsson.se:443" {
		t.Fatal("The ingress address should be https://test.ericsson.se:443, but not !")
	}

	nfService = TNfService{
		Fqdn:          "test.ericsson.se",
		Scheme:        "http",
		InterPlmnFqdn: "int.ericsson.se",
	}
	var nrfServce1 TNRFNFServices
	nrfServce1.init()
	nrfServce1.ManagementNfServices = append(nrfServce1.ManagementNfServices, nfService)
	nrfServce1.atomicSetNRFNFServices()
	if GetMgmtIngressAddress() != "http://test.ericsson.se:80" {
		t.Fatal("The ingress address should be https://test.ericsson.se:80, but not!")
	}
}

func TestIsDefaultIngressIP(t *testing.T) {
	ipURI := "https://127.0.0.1:443"
	if !isDefaultIngressIP(ipURI) {
		t.Fatal("Should be default Ingress IP, but NOT")
	}

	ipURI = "http://150.236.12.123:30082"
	if isDefaultIngressIP(ipURI) {
		t.Fatal("Should NOT be default Ingress IP, but it is")
	}
}

func TestIsDefaultIngressFqdn(t *testing.T) {
	fqdnURI := "https://nrf.ericsson.se:443"
	if !isDefaultIngressFqdn(fqdnURI) {
		t.Fatal("Should be default Ingress fqdn, but NOT")
	}

	fqdnURI = "http://nrf.r20.ericsson.se:30082"
	if isDefaultIngressFqdn(fqdnURI) {
		t.Fatal("Should NOT be default Ingress fqdn, but it is")
	}
}
