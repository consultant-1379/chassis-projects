package cm

import (
	"fmt"
)

var (
	// ProvisionService is configuration of provision service profile
	ProvisionService TProvisionService
)

// ParseConf is to parse configuration of provision service profile
func (conf *TProvisionService) ParseConf() {
	ProvisionService = *conf
}

// Show print discovery service profile
func (conf *TProvisionService) Show() {
	fmt.Printf("TProvisionService value is : %+v\n", ProvisionService)
}
