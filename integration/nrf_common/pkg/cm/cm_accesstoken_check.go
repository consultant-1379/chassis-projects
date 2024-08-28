package cm

import (
	//"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

// CheckConfForAccessTokenService is to check basic configuration of CM  for nrf-accesstoken serivice
func CheckConfForAccessTokenService() bool {

	if GetNRFNFProfile().InstanceID == "" {
		log.Errorf("nf-profile.instance-id in CM need to be configured,it can't be empty")
		return false
	}

	return true
}
