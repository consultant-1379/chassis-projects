package cm

import (
	"testing"
)

func TestConstructPlmnNrfAPIRoot(t *testing.T) {

	vPlmnNrf := &TPlmnNrf{}
	NrfCommon.Role = "region-nrf"
	NrfCommon.PlmnNrf = vPlmnNrf

	priority := 2
	vProfile := TNrfProfile{
		Priority: &priority,
		Fqdn:     "mcc.mnc.se",
	}

	vNrfService := TNrfService{
		Scheme: "http",
		Name:   "nnrf-nfm",
	}

	vProfile.Service = append(vProfile.Service, vNrfService)
	vPlmnNrf.Profile = append(vPlmnNrf.Profile, vProfile)

	// use fqdn of profile
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://mcc.mnc.se:80" {
		t.Fatalf("use fqdn of profile as plmn nrf address failed!")
	}

	// use Ipv4Address of profile
	vIPEndpoint := TIPEndpoint{
		Ipv4Address: "",
		Port:        0,
		Ipv6Address: "",
	}
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint = append(NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint, vIPEndpoint)
	vIPv4Address := "192.168.1.1"
	NrfCommon.PlmnNrf.Profile[0].Ipv4Address = append(NrfCommon.PlmnNrf.Profile[0].Ipv4Address, vIPv4Address)
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://192.168.1.1:80" {
		t.Fatalf("use Ipv4Address of profile as plmn nrf address failed!")
	}

	// use fqdn of service
	NrfCommon.PlmnNrf.Profile[0].Service[0].Fqdn = "mcc460.mnc000.se"
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://mcc460.mnc000.se:80" {
		t.Fatalf("use fqdn of service as plmn nrf address failed!")
	}

	// use ipEndpoint's ipv4 of service

	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint[0].Ipv4Address = "192.168.1.2"
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint[0].Port = 81
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://192.168.1.2:81" {
		t.Fatalf("use ipEndpoint's ipv4 of service as plmn nrf address failed!")
	}

	// use IPEndpoints's ipv6 of  NrfServiceEndpoint, and use specific port
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint[0].Ipv4Address = ""
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint[0].Port = 0
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint[0].Ipv6Address = "2001:470:c:1818::2"
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://[2001:470:c:1818::2]:80" {
		t.Fatalf("use ipEndpoint'ipv6 of nrf-service-endpoints as plmn nrf address failed!")
	}

	// add api-prefix for management
	NrfCommon.PlmnNrf.Profile[0].Service[0].APIPrefix = "mgmt"
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if PlmnNrfAPIRootMap[2][0] != "http://[2001:470:c:1818::2]:80/mgmt" {
		t.Fatalf("add api-prefix before plmn nrf address failed!")
	}

	// add api-prefix for discovery
	NrfCommon.PlmnNrf.Profile[0].Service[0].Name = "nnrf-disc"
	NrfCommon.PlmnNrf.Profile[0].Service[0].APIPrefix = "disc"
	ConstructPlmnNrfAPIRoot("nnrf-disc")
	if PlmnNrfAPIRootMap[2][0] != "http://[2001:470:c:1818::2]:80/disc" {
		t.Fatalf("add api-prefix before plmn nrf address failed!")
	}

	// management service ignores discovery addresses
	NrfCommon.PlmnNrf.Profile[0].Service[0].Name = "nnrf-disc"
	ConstructPlmnNrfAPIRoot("nnrf-nfm")
	if len(PlmnNrfAPIRootMap) != 0 {
		t.Fatalf("construct plmn nrf address for management failed!")
	}
}

func TestParsePlmnNrfHostPortByNrfProfile(t *testing.T) {
	var vIPEndpoints []TIPEndpoint
	var vIPv4Addresses []string
	var vIPv6Addresses []string

	//test IP is empty, fqdn is not empty,port uses default 80
	ipEndpoint := TIPEndpoint{
		ID:          1,
		Transport:   "TCP",
		Ipv4Address: "",
		Port:        0,
		Ipv6Address: "",
	}
	vIPEndpoints = append(vIPEndpoints, ipEndpoint)
	addrMap := parsePlmnNrfHostPortByNrfProfile(vIPv4Addresses, vIPv6Addresses, vIPEndpoints, "http", "ericsson.se")
	if len(addrMap) != 1 || addrMap["ericsson.se"] != 80 {
		t.Fatalf("parse IPEndpoints for fqdn failed!")
	}

	// test IPv6 and multiple address
	IPv6Address := "3ffe:2a00:100:7031::1"
	vIPv6Addresses = append(vIPv6Addresses, IPv6Address)
	addrMap = parsePlmnNrfHostPortByNrfProfile(vIPv4Addresses, vIPv6Addresses, vIPEndpoints, "https", "ericsson.se")
	if len(addrMap) != 1 || addrMap["[3ffe:2a00:100:7031::1]"] != 443 {
		t.Fatalf("parse IPEndpoints for IPv6 failed!")
	}

	// test IP is not empty, fqdn is empty, port uses default 443
	IPv4Address := "127.0.0.1"
	vIPv4Addresses = append(vIPv4Addresses, IPv4Address)
	addrMap = parsePlmnNrfHostPortByNrfProfile(vIPv4Addresses, vIPv6Addresses, vIPEndpoints, "https", "")
	if len(addrMap) != 1 || addrMap["127.0.0.1"] != 443 {
		t.Fatalf("parse IPEndpoints for IPv4 failed!")
	}
}

func TestParsePlmnNrfHostPortByNrfService(t *testing.T) {
	var vIPEndpoints []TIPEndpoint

	//test IP is empty, fqdn is not empty,port uses default 80
	ipEndpoint := TIPEndpoint{
		ID:          1,
		Transport:   "TCP",
		Ipv4Address: "",
		Port:        0,
		Ipv6Address: "",
	}
	vIPEndpoints = append(vIPEndpoints, ipEndpoint)
	addrMap := parsePlmnNrfHostPortByNrfService(vIPEndpoints, "http", "ericsson.se")
	if len(addrMap) != 1 || addrMap["ericsson.se"] != 80 {
		t.Fatalf("parse IPEndpoints for fqdn failed!")
	}

	// test IP is not empty, fqdn is empty, port uses default 443
	ipEndpoint1 := TIPEndpoint{
		ID:          1,
		Transport:   "TCP",
		Ipv4Address: "127.0.0.1",
		Port:        0,
		Ipv6Address: "",
	}
	vIPEndpoints = append(vIPEndpoints, ipEndpoint1)
	addrMap = parsePlmnNrfHostPortByNrfService(vIPEndpoints, "https", "")
	if len(addrMap) != 1 || addrMap["127.0.0.1"] != 443 {
		t.Fatalf("parse IPEndpoints for IPv4 failed!")
	}

	// test IPv6 and multiple address
	ipEndpoint3 := TIPEndpoint{
		ID:          1,
		Transport:   "TCP",
		Ipv4Address: "",
		Port:        80,
		Ipv6Address: "3ffe:2a00:100:7031::1",
	}
	vIPEndpoints = append(vIPEndpoints, ipEndpoint3)

	addrMap = parsePlmnNrfHostPortByNrfService(vIPEndpoints, "https", "ericsson.se")
	if len(addrMap) != 2 || addrMap["127.0.0.1"] != 443 || addrMap["[3ffe:2a00:100:7031::1]"] != 80 {
		t.Fatalf("parse IPEndpoints for IPv6 failed!")
	}
}
