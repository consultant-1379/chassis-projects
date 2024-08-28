package cm

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

const (
	// REGIONNRF is the role name of region nrf
	REGIONNRF = "region-nrf"
	// PLMNNRF is the role name of plmn nrf
	PLMNNRF = "plmn-nrf"
	// UNKNOWNROLE is the role name when neither region nrf nor plmn nrf
	UNKNOWNROLE          = "unknown"
	defaultHierarchyMode = "forward"

	// SAMELAYER means other nrf is in the same layer with me
	SAMELAYER = "same-layer"

	// UPPERLAYER means other nrf is in upper layer of me, e.g. PLMN NRF
	UPPERLAYER = "upper-layer"
)

var (
	// NrfCommon is configuration of nrf common
	NrfCommon *TCommon
)

func (conf *TCommon) atomicSetCommon() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&NrfCommon)), unsafe.Pointer(conf))
}

//GetNRFCommon to get nrfcommon configuration in cm
func GetNRFCommon() *TCommon {
	return (*TCommon)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&NrfCommon))))
}

// ParseConf is to parse nrf common
func (conf *TCommon) ParseConf() {

	if conf.GeoRed == nil {
		witnessNf := &TWitnessNF{
			IdentityType:  "",
			IdentityValue: "",
		}
		geoRed := &TGeoRed{
			KeepManagementService: true,
			KeepDiscoveryService:  true,
			WitnessNF:             witnessNf,
		}
		conf.GeoRed = geoRed
	} else {
		if conf.GeoRed.WitnessNF == nil {
			witnessNf := &TWitnessNF{
				IdentityType:  "",
				IdentityValue: "",
			}
			conf.GeoRed.WitnessNF = witnessNf
		}
	}

	if conf.PlmnNrf != nil && conf.PlmnNrf.Mode == "" {
		conf.PlmnNrf.Mode = defaultHierarchyMode
	}

	conf.atomicSetCommon()
}

// Show is to print nrf common
func (conf *TCommon) Show() {
	fmt.Printf("the nrfcommon is %v", GetNRFCommon())
}

// GetNRFRole returns the nrf role
func GetNRFRole() string {
	if GetNRFCommon().RegionNrf != nil {
		return REGIONNRF
	}

	if GetNRFCommon().PlmnNrf != nil {
		return PLMNNRF
	}

	return UNKNOWNROLE
}

// GetNextHopLayerType returns the next hop is the same layer or up layer
func GetNextHopLayerType() string {
	if GetNRFCommon().RegionNrf == nil {
		return ""
	}

	peerNRFInfoMap := GetPeerNRFInfoMap()
	if peerNRFInfoMap == nil {
		return ""
	}

	for _, peerNRFInfo := range *peerNRFInfoMap {
		if peerNRFInfo.Layer == "upper-layer" {
			return UPPERLAYER
		}
	}

	return SAMELAYER
}
