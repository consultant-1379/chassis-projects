package cm

import (
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestCheckNrfCommon_nrfrole(t *testing.T) {

	nrfCommon := &TCommon{
		Role: "plmn-nrf",
	}

	nrfCommon.ParseConf()

	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check nrf role fail!")
	}

	NrfCommon.Role = "invalidrole"

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("the Nrf Role must region-nrf or plmn-nrf")
	}
}

func TestCheckNrfCommon_PlmnNrf(t *testing.T) {

	NrfCommon.Role = "region-nrf"
	plmnnrf := &TPlmnNrf{Mode: "load-balance"}
	NrfCommon.PlmnNrf = plmnnrf
	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile fail")
	}

	priority := 5
	profile := TNrfProfile{
		ID:       "plmn-nrf-profile1",
		Priority: &priority,
	}
	NrfCommon.PlmnNrf.Profile = append(NrfCommon.PlmnNrf.Profile, profile)
	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile.Fqdn fail")
	}

	//IPEndpoint
	nrfService := TNrfService{Scheme: "httpx", Name: constvalue.NNRFNFM}
	NrfCommon.PlmnNrf.Profile[0].Service = append(NrfCommon.PlmnNrf.Profile[0].Service, nrfService)
	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile[0].Service.Scheme fail")
	}

	NrfCommon.PlmnNrf.Profile[0].Service[0].Scheme = "http"

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile[0].Service.scheme and IPEndpoint fail")
	}

	NrfCommon.PlmnNrf.Profile[0].Fqdn = "mcc.mnc.se"
	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile[0].Fqdn fail")
	}

	NrfCommon.PlmnNrf.Profile[0].Fqdn = ""
	NrfCommon.PlmnNrf.Profile[0].Service[0].Fqdn = "mcc460.mnc00.se"
	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile[0].Service[0].Fqdn fail")
	}

	NrfCommon.PlmnNrf.Profile[0].Service[0].Fqdn = ""
	vIPEndpoint := TIPEndpoint{
		Ipv4Address: "127.0.0.1",
		Port:        0,
		Ipv6Address: "",
	}
	NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint = append(NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint, vIPEndpoint)
	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("check NrfCommon.PlmnNrf.Profile[0].Service[0].IPEndpoint fail")
	}

}

func TestCheckNfProfile(t *testing.T) {

	plmnList := []TPLMN{}

	nfprofile := &TNfProfile{
		Status: "REGISTERED-test",
		PlmnID: plmnList,
	}

	nfprofile.ParseConf()

	if checkNfProfile() {
		t.Fatal("check nf-profile.plmn-id fail")
	}

}

func TestCheckManagementNrfService(t *testing.T) {

	ManagementNfServices = []TNfService{}

	if checkManagementNrfService() {
		t.Fatal("check management services fail")
	}

	nfService := TNfService{
		Scheme: "http-x",
	}

	ManagementNfServices = append(ManagementNfServices, nfService)

	if checkManagementNrfService() {
		t.Fatal("check management service scheme fail")
	}
	ManagementNfServices[0].Scheme = "http"

	if checkManagementNrfService() {
		t.Fatal("check management service fqdn fail")
	}

	// if management service ip-endpoint and fqdn are not configured, nf-profile.fqdn or ipv4-address is used.
	NfProfile.Fqdn = "mcc.mnc.se"
	ManagementNfServices[0].InstanceID = "nrf-mgmt-01"
	if !checkManagementNrfService() {
		t.Fatal("check nf-profile.fqdn fail")
	}

	NfProfile.Fqdn = ""
	vIPEndpoint := TIPEndpoint{
		Port: 80,
	}
	ManagementNfServices[0].IPEndpoint = append(ManagementNfServices[0].IPEndpoint, vIPEndpoint)
	if checkManagementNrfService() {
		t.Fatal("check management service ip-endpoint fail")
	}
	ManagementNfServices[0].IPEndpoint[0].Ipv4Address = "127.0.0.1"

	if !checkManagementNrfService() {
		t.Fatal("check management service ip-endpoint fail")
	}

}

func TestCheckDiscoveryNrfService(t *testing.T) {

	DiscoveryNfServices = []TNfService{}

	if checkDiscoveryNrfService() {
		t.Fatal("check discovery services fail")
	}

	nfService := TNfService{
		Scheme: "http-x",
	}

	DiscoveryNfServices = append(DiscoveryNfServices, nfService)

	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service scheme fail")
	}
	DiscoveryNfServices[0].Scheme = "http"

	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service fqdn fail")
	}

	// if discovery service ip-endpoint and fqdn are not configured, nf-profie.fqdn or ipv4-address is used.
	NfProfile.Fqdn = "mcc.mnc.se"
	DiscoveryNfServices[0].InstanceID = "nrf-mgmt-01"
	if !checkDiscoveryNrfService() {
		t.Fatal("check nf-profie.fqdn fail")
	}

	NfProfile.Fqdn = ""
	vIPEndpoint := TIPEndpoint{
		Port: 80,
	}
	DiscoveryNfServices[0].IPEndpoint = append(DiscoveryNfServices[0].IPEndpoint, vIPEndpoint)
	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}
	DiscoveryNfServices[0].IPEndpoint[0].Ipv4Address = "127.0.0.1"

	if !checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}

}
