package nrfschema

import (
	"bytes"
	"encoding/json"
	"fmt"

	"com/dbproxy/nfmessage/subscription"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidSupiRangeIndexs return invalid supiRangeList index
func (p *TPcfInfo) GetInvalidSupiRangeIndexs() []string {
	var invalidSupiRangeIndexs []string
	if p.SupiRanges != nil {
		index := 0
		for _, item := range p.SupiRanges {
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
func (p *TPcfInfo) GetInvalidGpsiRangeIndexs() []string {
	var invalidGpsiRangeIndexs []string
	if p.GpsiRanges != nil {
		index := 0
		for _, item := range p.GpsiRanges {
			if !item.IsValid() {
				invalidGpsiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.GpsiRanges, index)
				invalidGpsiRangeIndexs = append(invalidGpsiRangeIndexs, invalidGpsiRangeIndex)
			}
			index++
		}
	}

	return invalidGpsiRangeIndexs
}

func (p *TPcfInfo) createNfInfo() string {

	var supiRanges string
	var gpsiRanges string
	var groupID string
	var supiMatchAll string
	var gpsiMatchAll string

	if p.SupiRanges != nil && len(p.SupiRanges) > 0 {
		supiRangeList := ""
		for _, v := range p.SupiRanges {

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

	if p.GpsiRanges != nil && len(p.GpsiRanges) > 0 {
		gpsiRangeList := ""
		for _, v := range p.GpsiRanges {

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

	if p.GroupId != "" {
		groupID = fmt.Sprintf(`"groupId":"%s"`, p.GroupId)
	}

	if p.GroupId != "" || (p.SupiRanges != nil && len(p.SupiRanges) > 0) {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
	}

	if p.GroupId != "" || (p.GpsiRanges != nil && len(p.GpsiRanges) > 0) {
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
	}

	var dnnList string
	if p.DnnList != nil && len(p.DnnList) > 0 {
		for _, v := range p.DnnList {
			if dnnList == "" {
				dnnList = fmt.Sprintf(`"%v"`, v)
			} else {
				dnnList = fmt.Sprintf(`%s,"%v"`, dnnList, v)
			}
		}
		dnnList = fmt.Sprintf(`"dnnList":[%s]`, dnnList)
	} else {
		dnnList = fmt.Sprintf(`"dnnList":["%s"]`, constvalue.EmptyDnn)
	}
	if groupID != "" {
		return fmt.Sprintf(`"pcfInfo":{%s,%s,%s,%s,%s,%s}`, groupID, supiRanges, gpsiRanges, supiMatchAll, gpsiMatchAll, dnnList)
	}
	return fmt.Sprintf(`"pcfInfo":{%s,%s,%s,%s,%s}`, supiRanges, gpsiRanges, supiMatchAll, gpsiMatchAll, dnnList)
}

// GenerateNfGroupCond generate NfGroupCond
func (p *TPcfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (p *TPcfInfo) IsEqual(c TNfInfo) bool {

	b := c.(*TPcfInfo)

	if p.GroupId != b.GroupId || p.RxDiamHost != b.RxDiamHost || p.RxDiamRealm != b.RxDiamRealm {
		return false
	}

	if len(p.DnnList) != len(b.DnnList) {
		return false
	}

	if len(p.SupiRanges) != len(b.SupiRanges) {
		return false
	}

	if len(p.GpsiRanges) != len(b.GpsiRanges) {
		return false
	}

	for k, item := range p.DnnList {
		if item != b.DnnList[k] {
			return false
		}
	}

	for k, item := range p.SupiRanges {
		bb := b.SupiRanges[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	for k, item := range p.GpsiRanges {
		bb := b.GpsiRanges[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	return true
}
