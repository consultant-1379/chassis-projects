package cm

import (
	"fmt"
	"strings"

	"encoding/json"
	"sync/atomic"
	"unsafe"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

const (
	defaultNrfProfileType = "nrf"
	defaultNrfStatus      = "suspended"
)

var (
	// NfProfile is configuration of nf profile
	NfProfile *TNfProfile
	//PlmnListRaw is raw data of plmnList in NfProfile
	PlmnListRaw *TPlmnListRawData
)

// TPlmnListRawData to store homeplmn raw data
type TPlmnListRawData struct {
	// PlmnList is raw data of plmnList in NfProfile
	PlmnList []byte
}

func (p *TPlmnListRawData) atomicSetPlmnListRaw() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&PlmnListRaw)), unsafe.Pointer(p))
}

func (p *TPlmnListRawData) init() {
	p.PlmnList = make([]byte, 1)
}

//GetPlmnListRaw to get plmnlist raw
func GetPlmnListRaw() *TPlmnListRawData {
	return (*TPlmnListRawData)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&PlmnListRaw))))
}

func (conf *TNfProfile) atomicSetNFProfile() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&NfProfile)), unsafe.Pointer(conf))
}

//GetNRFNFProfile to get nrf nfprofile
func GetNRFNFProfile() *TNfProfile {
	return (*TNfProfile)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&NfProfile))))
}

// ParseConf is to parse nf profile
func (conf *TNfProfile) ParseConf() {

	if conf.Type == "" {
		conf.Type = defaultNrfProfileType
	}

	if conf.RequestedStatus == "" {
		conf.RequestedStatus = defaultNrfStatus
	}

	conf.toUpper()

	var nrfNFServices TNRFNFServices
	nrfNFServices.init()
	for _, service := range conf.Service {
		if service.Name == constvalue.NNRFNFM {
			nrfNFServices.ManagementNfServices = append(nrfNFServices.ManagementNfServices, service)
		} else if service.Name == constvalue.NNRFDISC {
			nrfNFServices.DiscoveryNfServices = append(nrfNFServices.DiscoveryNfServices, service)
		}
	}

	nrfNFServices.atomicSetNRFNFServices()

	conf.plmnList()
	conf.atomicSetNFProfile()
}

func (conf *TNfProfile) plmnList() {
	var plmnListRaw TPlmnListRawData
	plmnListRaw.init()
	plmnList, err := json.Marshal(conf.PlmnID)
	if err != nil {
		fmt.Printf("Marshal PlmnList fail, %v", err)
		plmnListRaw.PlmnList = []byte("")
	} else {
		plmnListRaw.PlmnList = plmnList
	}
	plmnListRaw.atomicSetPlmnListRaw()
}

// toUpper UPPER some 3gpp value, e.g. NF type and NF status
func (conf *TNfProfile) toUpper() {
	conf.Type = strings.ToUpper(conf.Type)
	conf.RequestedStatus = strings.ToUpper(conf.RequestedStatus)

	for index := range conf.AllowedNfType {
		conf.AllowedNfType[index] = strings.ToUpper(conf.AllowedNfType[index])
	}

	for index := range conf.Service {
		conf.Service[index].Status = strings.ToUpper(conf.Service[index].Status)

		for subIndex := range conf.Service[index].IPEndpoint {
			conf.Service[index].IPEndpoint[subIndex].Transport = strings.ToUpper(conf.Service[index].IPEndpoint[subIndex].Transport)
		}

		for subIndex := range conf.Service[index].AllowedNfType {
			conf.Service[index].AllowedNfType[subIndex] = strings.ToUpper(conf.Service[index].AllowedNfType[subIndex])
		}

		for subIndex := range conf.Service[index].DefaultNotificationSubscription {
			conf.Service[index].DefaultNotificationSubscription[subIndex].NotificationType = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].NotificationType)
			conf.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass)
			conf.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass = strings.ToUpper(conf.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass)
		}
	}
}

// Show is to print nf profile info
func (conf *TNfProfile) Show() {
	for index, plmn := range NfProfile.PlmnID {
		fmt.Printf("Plmnlist[%d] mcc : %s\n", index, plmn.Mcc)
		fmt.Printf("Plmnlist[%d] mnc : %s\n", index, plmn.Mnc)
	}
}
