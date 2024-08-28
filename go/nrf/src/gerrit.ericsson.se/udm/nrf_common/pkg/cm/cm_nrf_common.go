package cm

import (
	"fmt"
)

var (
	// NrfCommon is configuration of nrf common
	NrfCommon TCommon
)

// ParseConf is to parse nrf common
func (conf *TCommon) ParseConf() {
	NrfCommon = *conf

	if NrfCommon.GeoRed == nil {
		witnessNf := &TWitnessNF{
			IdentityType:  "",
			IdentityValue: "",
		}
		geoRed := &TGeoRed{
			KeepManagementService: true,
			KeepDiscoveryService:  true,
			WitnessNF:             witnessNf,
		}
		NrfCommon.GeoRed = geoRed
	} else {
		if NrfCommon.GeoRed.WitnessNF == nil {
			witnessNf := &TWitnessNF{
				IdentityType:  "",
				IdentityValue: "",
			}
			NrfCommon.GeoRed.WitnessNF = witnessNf
		}
	}

	if NrfCommon.RemoteDefaultSetting == nil {
		remoteDefaultSetting := &TRemoteDefaultSetting{
			Scheme: "https",
			Port:   443,
		}
		NrfCommon.RemoteDefaultSetting = remoteDefaultSetting
	}
}

// Show is to print nrf common
func (conf *TCommon) Show() {
	fmt.Printf("the nrfcommon is %v", NrfCommon)
}
