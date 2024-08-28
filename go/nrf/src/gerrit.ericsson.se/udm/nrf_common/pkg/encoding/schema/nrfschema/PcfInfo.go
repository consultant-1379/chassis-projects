package nrfschema

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (p *TPcfInfo) createNfInfo() string {

	var supiRanges string
	var groupID string

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

	if p.GroupId != "" {
		groupID = fmt.Sprintf(`"groupId":"%s"`, p.GroupId)
	} else {
		groupID = fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
	}

	var supiMatchAll string
	if p.GroupId != "" || (p.SupiRanges != nil && len(p.SupiRanges) > 0) {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
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
	}
	if dnnList != "" {
		return fmt.Sprintf(`"pcfInfo":{%s,%s,%s,%s}`, groupID, supiRanges, supiMatchAll, dnnList)
	}
	return fmt.Sprintf(`"pcfInfo":{%s,%s,%s}`, groupID, supiRanges, supiMatchAll)
}
