package nrfschema

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

var (
	//NfGroupCondNfType records the nfType who has nfGroupId
	NfGroupCondNfType = map[string]bool{
		constvalue.NfTypeAUSF: true,
		constvalue.NfTypeUDM:  true,
		constvalue.NfTypeUDR:  true,
	}
)

//Validate returns the reason of invalidation
func (s *TSubscrCond) Validate() string {
	countOfCOnds := 0

	if s.NfInstanceID != "" {
		countOfCOnds++
	}
	if s.ServiceName != "" {
		countOfCOnds++
	}
	if s.AmfSetID != "" || s.AmfRegionID != "" {
		countOfCOnds++
	}
	if s.GuamiList != nil {
		countOfCOnds++
	}
	if s.SnssaiList != nil || s.NsiList != nil {
		countOfCOnds++
	}
	if s.NfGroupID != "" || s.NfType != "" {
		countOfCOnds++
	}

	if countOfCOnds != 1 {
		return constvalue.StatusSubscribeRule1
	}

	if s.NsiList != nil && s.SnssaiList == nil {
		return constvalue.StatusSubscribeRule2
	}

	if s.NfGroupID != "" {
		if s.NfType == "" {
			return constvalue.StatusSubscribeRule3
		}

		if !NfGroupCondNfType[s.NfType] {
			return constvalue.StatusSubscribeRule4
		}
	}

	return ""

}
