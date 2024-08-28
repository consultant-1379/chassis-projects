package nrfschema

import (
	"bytes"
	"com/dbproxy/nfmessage/subscription"
	"encoding/json"
	"fmt"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidSupiRangeIndexs return invalid supiRanges index
func (u *TUdmInfo) GetInvalidSupiRangeIndexs() []string {
	var invalidSupiRangeIndexs []string
	if u.SupiRanges != nil {
		index := 0
		for _, item := range u.SupiRanges {
			if !item.IsValid() {
				invalidSupiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.SupiRanges, index)
				invalidSupiRangeIndexs = append(invalidSupiRangeIndexs, invalidSupiRangeIndex)
			}
			index++
		}
	}

	return invalidSupiRangeIndexs
}

// GetInvalidGpsiRangeIndexs return invalid gpsiRanges index
func (u *TUdmInfo) GetInvalidGpsiRangeIndexs() []string {
	var invalidGpsiRangeIndexs []string
	if u.GpsiRanges != nil {
		index := 0
		for _, item := range u.GpsiRanges {
			if !item.IsValid() {
				invalidGpsiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.GpsiRanges, index)
				invalidGpsiRangeIndexs = append(invalidGpsiRangeIndexs, invalidGpsiRangeIndex)
			}
			index++
		}
	}

	return invalidGpsiRangeIndexs
}

// GetInvalidEGIRangeIndexs return invalid externalGroupIdentifiersRanges index
func (u *TUdmInfo) GetInvalidEGIRangeIndexs() []string {
	var invalidEGIRangeIndexs []string
	if u.ExternalGroupIdentifiersRanges != nil {
		index := 0
		for _, item := range u.ExternalGroupIdentifiersRanges {
			if !item.IsValid() {
				invalidEGIRangeIndex := fmt.Sprintf("%s[%d]", constvalue.ExternalGroupIdentityfiersRanges, index)
				invalidEGIRangeIndexs = append(invalidEGIRangeIndexs, invalidEGIRangeIndex)
			}
			index++
		}
	}

	return invalidEGIRangeIndexs
}

func (u *TUdmInfo) createNfInfo() string {
	var supiRanges string
	var gpsiRanges string
	var externalID string
	var groupID string
	var supiMatchAll string
	var gpsiMatchAll string
	if u.ExternalGroupIdentifiersRanges != nil && len(u.ExternalGroupIdentifiersRanges) > 0 {
		externalIDList := ""
		for _, v := range u.ExternalGroupIdentifiersRanges {

			externalIDRangeItem := ""
			if v.Start != "" {
				start := fmt.Sprintf(`"start":"%s"`, v.Start)
				if externalIDRangeItem != "" {
					externalIDRangeItem += ","
				}
				externalIDRangeItem += start
			}

			if v.End != "" {
				end := fmt.Sprintf(`"end":"%s"`, v.End)
				if externalIDRangeItem != "" {
					externalIDRangeItem += ","
				}
				externalIDRangeItem += end
			}

			if v.Pattern != "" {
				buffer := &bytes.Buffer{}
				encoder := json.NewEncoder(buffer)
				err := encoder.Encode(v.Pattern)
				if nil == err {
					pattern := fmt.Sprintf(`"pattern":%s`, string(buffer.Bytes()))
					if externalIDRangeItem != "" {
						externalIDRangeItem += ","
					}
					externalIDRangeItem += pattern
				}
			}

			if externalIDRangeItem != "" {
				if externalIDList != "" {
					externalIDList += ","
				}
				externalIDList += "{" + externalIDRangeItem + "}"
			}
		}

		if externalIDList == "" {
			externalID = fmt.Sprintf(`"externalGroupIdentifiersRanges":[{"pattern":"%s"}]`, constvalue.EmptyExternalIDPattern)
		} else {
			externalID = fmt.Sprintf(`"externalGroupIdentifiersRanges":[%s]`, externalIDList)
		}
	} else {
		externalID = fmt.Sprintf(`"externalGroupIdentifiersRanges":[{"pattern":"%s"}]`, constvalue.EmptyExternalIDPattern)
	}

	if u.SupiRanges != nil && len(u.SupiRanges) > 0 {
		supiRangeList := ""
		for _, v := range u.SupiRanges {

			supiRangeItem := ""
			if v.Start != "" {
				start := fmt.Sprintf(`"start":"%s"`, v.Start)
				if supiRangeItem != "" {
					supiRangeItem += ","
				}
				supiRangeItem += start
			}

			if v.End != "" {
				end := fmt.Sprintf(`"end":"%s"`, v.End)
				if supiRangeItem != "" {
					supiRangeItem += ","
				}
				supiRangeItem += end
			}

			if v.Pattern != "" {
				buffer := &bytes.Buffer{}
				encoder := json.NewEncoder(buffer)
				err := encoder.Encode(v.Pattern)
				if nil == err {
					pattern := fmt.Sprintf(`"pattern":%s`, string(buffer.Bytes()))
					if supiRangeItem != "" {
						supiRangeItem += ","
					}
					supiRangeItem += pattern
				}
			}

			if supiRangeItem != "" {
				if supiRangeList != "" {
					supiRangeList += ","
				}
				supiRangeList += "{" + supiRangeItem + "}"
			}
		}

		if supiRangeList == "" {
			supiRanges = fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		} else {
			supiRanges = fmt.Sprintf(`"supiRanges":[%s]`, supiRangeList)
		}
	} else {
		supiRanges = fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
	}

	if u.GpsiRanges != nil && len(u.GpsiRanges) > 0 {
		gpsiRangeList := ""
		for _, v := range u.GpsiRanges {

			gpsiRangeItem := ""
			if v.Start != "" {
				start := fmt.Sprintf(`"start":"%s"`, v.Start)
				if gpsiRangeItem != "" {
					gpsiRangeItem += ","
				}
				gpsiRangeItem += start
			}

			if v.End != "" {
				end := fmt.Sprintf(`"end":"%s"`, v.End)
				if gpsiRangeItem != "" {
					gpsiRangeItem += ","
				}
				gpsiRangeItem += end
			}

			if v.Pattern != "" {
				buffer := &bytes.Buffer{}
				encoder := json.NewEncoder(buffer)
				err := encoder.Encode(v.Pattern)
				if nil == err {
					pattern := fmt.Sprintf(`"pattern":%s`, string(buffer.Bytes()))
					if gpsiRangeItem != "" {
						gpsiRangeItem += ","
					}
					gpsiRangeItem += pattern
				}
			}

			if gpsiRangeItem != "" {
				if gpsiRangeList != "" {
					gpsiRangeList += ","
				}
				gpsiRangeList += "{" + gpsiRangeItem + "}"
			}
		}

		if gpsiRangeList == "" {
			gpsiRanges = fmt.Sprintf(`"gpsiRanges":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		} else {
			gpsiRanges = fmt.Sprintf(`"gpsiRanges":[%s]`, gpsiRangeList)
		}
	} else {
		gpsiRanges = fmt.Sprintf(`"gpsiRanges":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
	}

	if u.GroupId != "" {
		groupID = fmt.Sprintf(`"groupId":"%s"`, u.GroupId)
	}

	if u.GroupId != "" || (u.SupiRanges != nil && len(u.SupiRanges) > 0) {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
	}

	if u.GroupId != "" || (u.GpsiRanges != nil && len(u.GpsiRanges) > 0) {
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
	}

	var routingIndicators string
	if u.RoutingIndicators != nil && len(u.RoutingIndicators) > 0 {
		for _, v := range u.RoutingIndicators {
			if routingIndicators == "" {
				routingIndicators = fmt.Sprintf(`"%v"`, v)
			} else {
				routingIndicators = fmt.Sprintf(`%s,"%v"`, routingIndicators, v)
			}
		}
		routingIndicators = fmt.Sprintf(`"routingIndicators":[%s]`, routingIndicators)
	} else {
		routingIndicators = fmt.Sprintf(`"routingIndicators":["%s"]`, constvalue.EmptyRoutingIndicator)
	}
	if groupID != "" {
		return fmt.Sprintf(`"udmInfo":{%s,%s,%s,%s,%s,%s,%s}`, groupID, supiRanges, gpsiRanges, externalID, supiMatchAll, gpsiMatchAll, routingIndicators)
	}
	return fmt.Sprintf(`"udmInfo":{%s,%s,%s,%s,%s,%s}`, supiRanges, gpsiRanges, externalID, supiMatchAll, gpsiMatchAll, routingIndicators)
}

//GenerateNfGroupCond generate NfGroupCond for subscription
func (u *TUdmInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	if u.GroupId != "" {
		return &subscription.SubKeyStruct{
			SubKey1: u.GroupId,
			SubKey2: constvalue.NfTypeUDM,
		}
	}

	return nil
}

// IsEqual is to check if NFInfo is equal
func (u *TUdmInfo) IsEqual(c TNfInfo) bool {

	b := c.(*TUdmInfo)

	if u.GroupId != b.GroupId || len(u.ExternalGroupIdentifiersRanges) != len(b.ExternalGroupIdentifiersRanges) {
		return false
	}

	if len(u.GpsiRanges) != len(b.GpsiRanges) || len(u.SupiRanges) != len(b.SupiRanges) {
		return false
	}

	if len(u.RoutingIndicators) != len(b.RoutingIndicators) {
		return false
	}

	for k, item := range u.RoutingIndicators {
		if item != b.RoutingIndicators[k] {
			return false
		}
	}

	for k, item := range u.ExternalGroupIdentifiersRanges {
		bb := b.ExternalGroupIdentifiersRanges[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	for k, item := range u.GpsiRanges {
		bb := b.GpsiRanges[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	for k, item := range u.SupiRanges {
		bb := b.SupiRanges[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	return true
}
