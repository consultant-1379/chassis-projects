package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
	"fmt"

	"bytes"
	"encoding/json"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidSupiRangeIndexs return invalid supiRanges index
func (a *TAusfInfo) GetInvalidSupiRangeIndexs() []string {
	var invalidSupiRangeIndexs []string
	if a.SupiRanges != nil {
		index := 0
		for _, item := range a.SupiRanges {
			if !item.IsValid() {
				invalidSupiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.SupiRanges, index)
				invalidSupiRangeIndexs = append(invalidSupiRangeIndexs, invalidSupiRangeIndex)
			}
			index++
		}
	}

	return invalidSupiRangeIndexs
}

func (a *TAusfInfo) createNfInfo() string {
	var supiRanges string
	var groupID string
	if a.SupiRanges != nil && len(a.SupiRanges) > 0 {
		supiRangeList := ""
		for _, v := range a.SupiRanges {

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

	if a.GroupId != "" {
		groupID = fmt.Sprintf(`"groupId":"%s"`, a.GroupId)
	} else {
		groupID = fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
	}

	var supiMatchAll string
	if a.GroupId != "" || (a.SupiRanges != nil && len(a.SupiRanges) > 0) {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
	}

	var routingIndicators string
	if a.RoutingIndicators != nil && len(a.RoutingIndicators) > 0 {
		for _, v := range a.RoutingIndicators {
			if routingIndicators == "" {
				routingIndicators = fmt.Sprintf(`"%v"`, v)
			} else {
				routingIndicators = fmt.Sprintf(`%s,"%v"`, routingIndicators, v)
			}
		}
		routingIndicators = fmt.Sprintf(`"routingIndicators":[%s]`, routingIndicators)
	}
	if routingIndicators != "" {
		return fmt.Sprintf(`"ausfInfo":{%s, %s, %s, %s}`, groupID, supiRanges, supiMatchAll, routingIndicators)
	}
	return fmt.Sprintf(`"ausfInfo":{%s, %s, %s}`, groupID, supiRanges, supiMatchAll)
}

//GenerateNfGroupCond generate NfGroupCond for subscription
func (a *TAusfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	if a.GroupId != "" {
		return &subscription.SubKeyStruct{
			SubKey1: a.GroupId,
			SubKey2: constvalue.NfTypeAUSF,
		}
	}

	return nil
}
