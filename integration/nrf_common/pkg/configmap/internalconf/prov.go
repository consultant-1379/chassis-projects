package internalconf

import (
	"fmt"
)

const (
	defaultSyncNFProfileTimer = 900
)

var (
	//SyncNFProfileTimer is the cyclic timer value for sync nfProfile to group profile
	SyncNFProfileTimer = defaultSyncNFProfileTimer
)

// InternalProvConf is format of internalConf.json in NRF Provision
type InternalProvConf struct {
	HttpServer HttpServerConf `json:"httpServer"`
	Provision  ProvisionConf  `json:"provision, omitempty"`
}

// ProvisionConf is used to store provision service info
type ProvisionConf struct {
	SyncNFProfileTimer int `json:"syncNFProfileTimer, omitempty"`
}

// ParseConf sets the value of global variable from struct
func (conf *InternalProvConf) ParseConf() {
	HTTPWithXVersion = conf.HttpServer.HTTPWithXVersion

	SyncNFProfileTimer = conf.Provision.SyncNFProfileTimer
	if SyncNFProfileTimer <= 0 {
		SyncNFProfileTimer = defaultSyncNFProfileTimer
	}
}

// ShowInternalConf shows the value of global variable
func (conf *InternalProvConf) ShowInternalConf() {
	fmt.Printf("http header with x-version : %v\n", HTTPWithXVersion)
	fmt.Printf("Provision sync NFProfile timer : %d\n", SyncNFProfileTimer)
}
