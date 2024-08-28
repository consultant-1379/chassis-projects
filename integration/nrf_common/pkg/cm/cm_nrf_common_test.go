package cm

import (
	"testing"
)

func TestGetNRFRole(t *testing.T) {
	// case1: neither common.plmn-nrf nor common.region-nrf is configured, should return unknown
	nrfCommon := TCommon{}
	nrfCommon.atomicSetCommon()

	if GetNRFRole() != UNKNOWNROLE {
		t.Fatal("neither common.plmn-nrf nor common.region-nrf is configured, GetNRFRole should return unknown, but not !")
	}

	// case2: common.plmn-nrf is configured, should return plmn-nrf
	nrfCommon = TCommon{
		PlmnNrf: &TPlmnNrf{},
	}
	nrfCommon.atomicSetCommon()

	if GetNRFRole() != PLMNNRF {
		t.Fatal("common.plmn-nrf is configured, GetNRFRole should return plmn-nrf, but not !")
	}

	// case3: common.region-nrf is configured, should return plmn-nrf
	nrfCommon = TCommon{
		RegionNrf: &TRegionNrf{
			NextHop: &TNextHopNrf{
				Site: []TSiteNRF{
					TSiteNRF{Profile: []TNrfProfile{}},
				},
			},
		},
	}

	nrfCommon.atomicSetCommon()

	if GetNRFRole() != REGIONNRF {
		t.Fatal("common.region-nrf is configured, GetNRFRole should return region-nrf, but not !")
	}
}

func TestParseConfForCommon(t *testing.T) {
	NrfCommon.ParseConf()

	if NrfCommon.GeoRed == nil || NrfCommon.GeoRed.WitnessNF == nil {
		t.Fatalf("NrfCommon parse error!")
	}

	if NrfCommon.GeoRed.KeepDiscoveryService != true || NrfCommon.GeoRed.KeepManagementService != true {
		t.Fatalf("NrfCommon parse error!")
	}

	if NrfCommon.GeoRed.WitnessNF.IdentityType != "" || NrfCommon.GeoRed.WitnessNF.IdentityValue != "" {
		t.Fatalf("NrfCommon parse error!")
	}

}
