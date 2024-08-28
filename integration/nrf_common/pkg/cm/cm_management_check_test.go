package cm

import (
	"os"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestcheckNrfCommon(t *testing.T) {
	// case1: role choice not configured, the checking not pass

	nrfCommon := &TCommon{}

	nrfCommon.ParseConf()

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("the role choice common.plmn-nrf or common.region-nrf is not configured, but checkNrfCommon() pass")
	}

	// case2: role choice is plmn-nrf, the checking pass
	nrfCommon = &TCommon{
		PlmnNrf: &TPlmnNrf{
			Mode: "forward",
		},
	}

	nrfCommon.ParseConf()

	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("the role choice common.plmn-nrf is configured, but checkNrfCommon() didn't pass")
	}

	// case3: role choice is region-nrf, but plmn-nrf.profile is not configured, the checking not pass
	nrfCommon = &TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{},
		},
	}

	nrfCommon.ParseConf()

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("common.region-nrf.plmn-nrf.profile is not configured, but checkNrfCommon() pass")
	}

	// case4: role choice is region-nrf, but plmn-nrf.profile.service is not configured, the checking not pass
	nrfCommon = &TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{
				Site: []TSiteNRF{
					TSiteNRF{Profile: []TNrfProfile{}},
				},
			},
		},
	}

	nrfCommon.ParseConf()

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("common.region-nrf.plmn-nrf.profile.service is not configured, but checkNrfCommon() pass")
	}

	// case5: role choice is region-nrf, but plmn-nrf.profile.service.scheme is wrong, the checking not pass
	nrfCommon = &TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{
				Site: []TSiteNRF{
					TSiteNRF{
						Profile: []TNrfProfile{
							TNrfProfile{
								Service: []TNrfService{
									TNrfService{
										Scheme: "tcp",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	nrfCommon.ParseConf()

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("common.region-nrf.plmn-nrf.profile.service.scheme is wrong, but checkNrfCommon() pass")
	}

	// case6: role choice is region-nrf, but no plmn nrf address is configured, the checking not pass
	nrfCommon = &TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{
				Site: []TSiteNRF{
					TSiteNRF{
						Profile: []TNrfProfile{
							TNrfProfile{
								Service: []TNrfService{
									TNrfService{
										Scheme: "http",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	nrfCommon.ParseConf()

	if checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("plmn nrf address is not configured, but checkNrfCommon() pass")
	}

	// case7: role choice is region-nrf, and plmn nrf address is configured rightly, the checking pass
	nrfCommon = &TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{
				Site: []TSiteNRF{
					TSiteNRF{
						Profile: []TNrfProfile{
							TNrfProfile{

								Service: []TNrfService{
									TNrfService{
										Scheme: "http",
										IPEndpoint: []TIPEndpoint{
											TIPEndpoint{
												ID:          1,
												Transport:   "tcp",
												Ipv4Address: "10.10.10.10",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	nrfCommon.ParseConf()

	if !checkNrfCommon(constvalue.NNRFNFM) {
		t.Fatal("plmn nrf address is configured rightly, but checkNrfCommon() didn't pass")
	}
}

func TestCheckNfProfile(t *testing.T) {

	// case1: plmn-id is not configured, the checking not pass
	nfprofile := &TNfProfile{
		PlmnID:     []TPLMN{},
		InstanceID: "nrf01",
	}

	nfprofile.ParseConf()

	if checkNfProfile() {
		t.Fatal("nf-profile.plmn-id is not configured, but checkNfProfile pass")
	}

	// case2: instance-id is not configured, the checking not pass
	nfprofile = &TNfProfile{
		PlmnID: []TPLMN{
			TPLMN{
				Mcc: "460",
				Mnc: "00",
			},
		},
	}

	nfprofile.ParseConf()

	if checkNfProfile() {
		t.Fatal("nf-profile.instance-id is not configured, but checkNfProfile pass")
	}

	// case3: plmn-id and instance-id are configured, the checking pass
	nfprofile = &TNfProfile{
		PlmnID: []TPLMN{
			TPLMN{
				Mcc: "460",
				Mnc: "00",
			},
		},
		InstanceID: "nrf01",
	}

	nfprofile.ParseConf()

	if !checkNfProfile() {
		t.Fatal("nf-profile.instance-id and nf-profile.plmn-id are configured, but checkNfProfile didn't pass")
	}

}

func TestCheckManagementNrfService(t *testing.T) {

	var nfprofile0 TNfProfile
	nfprofile0.atomicSetNFProfile()
	var nrfService0 TNRFNFServices
	nrfService0.init()
	nrfService0.atomicSetNRFNFServices()
	if checkManagementNrfService() {
		t.Fatal("check management services fail")
	}

	nfService := TNfService{
		Scheme: "http-x",
	}
	var nrfService TNRFNFServices
	nrfService.init()
	nrfService.ManagementNfServices = append(nrfService.ManagementNfServices, nfService)
	nrfService.atomicSetNRFNFServices()

	if checkManagementNrfService() {
		t.Fatal("check management service scheme fail")
	}

	var nrfService1 TNRFNFServices
	nrfService1.init()
	nfService.Scheme = "http"
	nrfService1.ManagementNfServices = append(nrfService1.ManagementNfServices, nfService)
	nrfService1.atomicSetNRFNFServices()
	if checkManagementNrfService() {
		t.Fatal("check management service fqdn fail")
	}

	// if management service ip-endpoint and fqdn are not configured, nf-profile.fqdn or ipv4-address is used.

	var nfProfile TNfProfile
	nfProfile.Fqdn = "mcc.mnc.se"
	nfProfile.atomicSetNFProfile()

	var nrfService2 TNRFNFServices
	nrfService2.init()
	nfService.InstanceID = "nrf-mgmt-01"
	nfService.Scheme = "http"
	nrfService2.ManagementNfServices = append(nrfService2.ManagementNfServices, nfService)
	nrfService2.atomicSetNRFNFServices()
	if !checkManagementNrfService() {
		t.Fatal("check nf-profile.fqdn fail")
	}

	vIPEndpoint := TIPEndpoint{
		Port: 80,
	}

	os.Setenv("IP_STACK_MODE", "ipv4")

	var nfProfile1 TNfProfile
	nfProfile1.Fqdn = ""
	nfProfile1.atomicSetNFProfile()

	var nrfService3 TNRFNFServices
	nrfService3.init()
	nrfService3.ManagementNfServices = append(nrfService3.ManagementNfServices, nfService)
	nrfService3.ManagementNfServices[0].IPEndpoint = append(nrfService3.ManagementNfServices[0].IPEndpoint, vIPEndpoint)
	nrfService3.atomicSetNRFNFServices()
	if checkManagementNrfService() {
		t.Fatal("check management service ip-endpoint fail")
	}

	var nrfService4 TNRFNFServices
	nrfService4.init()
	nrfService4.ManagementNfServices = append(nrfService4.ManagementNfServices, TNfService{Scheme: "http", InstanceID: "nrf-mgmt-01"})
	nrfService4.ManagementNfServices[0].IPEndpoint = append(nrfService4.ManagementNfServices[0].IPEndpoint, vIPEndpoint)
	nrfService4.ManagementNfServices[0].IPEndpoint[0].Ipv4Address = "127.0.0.1"
	nrfService4.atomicSetNRFNFServices()
	if !checkManagementNrfService() {
		t.Fatal("check management service ip-endpoint fail")
	}

}

func TestCheckDiscoveryNrfService(t *testing.T) {

	var nfprofile0 TNfProfile
	nfprofile0.atomicSetNFProfile()
	var nrfService0 TNRFNFServices
	nrfService0.init()
	nrfService0.atomicSetNRFNFServices()

	if checkDiscoveryNrfService() {
		t.Fatal("check discovery services fail")
	}

	nfService := TNfService{
		Scheme: "http-x",
	}

	var nrfService TNRFNFServices
	nrfService.init()
	nrfService.DiscoveryNfServices = append(nrfService.DiscoveryNfServices, nfService)
	nrfService.atomicSetNRFNFServices()

	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service scheme fail")
	}

	var nrfService1 TNRFNFServices
	nrfService1.init()
	nfService.Scheme = "http"
	nrfService1.DiscoveryNfServices = append(nrfService1.DiscoveryNfServices, nfService)
	nrfService1.atomicSetNRFNFServices()
	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service fqdn fail")
	}

	// if discovery service ip-endpoint and fqdn are not configured, nf-profie.fqdn or ipv4-address is used.
	var nfProfile TNfProfile
	nfProfile.Fqdn = "mcc.mnc.se"
	nfProfile.atomicSetNFProfile()

	var nrfService2 TNRFNFServices
	nrfService2.init()
	nfService.InstanceID = "nrf-mgmt-01"
	nfService.Scheme = "http"
	nrfService2.DiscoveryNfServices = append(nrfService2.DiscoveryNfServices, nfService)
	nrfService2.atomicSetNRFNFServices()
	if !checkDiscoveryNrfService() {
		t.Fatal("check nf-profie.fqdn fail")
	}

	vIPEndpoint := TIPEndpoint{
		Port: 80,
	}

	var nfProfile1 TNfProfile
	nfProfile1.Fqdn = ""
	nfProfile1.atomicSetNFProfile()

	var nrfService3 TNRFNFServices
	nrfService3.init()
	nrfService3.DiscoveryNfServices = append(nrfService3.DiscoveryNfServices, nfService)
	nrfService3.DiscoveryNfServices[0].IPEndpoint = append(nrfService3.DiscoveryNfServices[0].IPEndpoint, vIPEndpoint)
	nrfService3.atomicSetNRFNFServices()
	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}

	var nrfService4 TNRFNFServices
	nrfService4.init()
	nrfService4.DiscoveryNfServices = append(nrfService4.DiscoveryNfServices, TNfService{Scheme: "http", InstanceID: "nrf-mgmt-01"})
	nrfService4.DiscoveryNfServices[0].IPEndpoint = append(nrfService4.DiscoveryNfServices[0].IPEndpoint, TIPEndpoint{Port: 80, Ipv4Address: "127.0.0.1"})
	nrfService4.atomicSetNRFNFServices()

	if !checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}

	os.Setenv("IP_STACK_MODE", "ipv6")
	if checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}
	var nrfService5 TNRFNFServices
	nrfService5.init()
	nrfService5.DiscoveryNfServices = append(nrfService5.DiscoveryNfServices, TNfService{Scheme: "http", InstanceID: "nrf-mgmt-01"})
	nrfService5.DiscoveryNfServices[0].IPEndpoint = append(nrfService5.DiscoveryNfServices[0].IPEndpoint, TIPEndpoint{Port: 80, Ipv4Address: "127.0.0.1", Ipv6Address: "1080::8:800:200C:417A"})
	nrfService5.atomicSetNRFNFServices()
	if !checkDiscoveryNrfService() {
		t.Fatal("check discovery service ip-endpoint fail")
	}

}
