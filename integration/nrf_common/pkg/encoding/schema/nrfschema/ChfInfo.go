package nrfschema

import (
	"fmt"

	"bytes"
	"com/dbproxy/nfmessage/subscription"
	"encoding/json"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidSupiRangeIndexs return invalid supiRangeList index
func (c *TChfInfo) GetInvalidSupiRangeIndexs() []string {
	var invalidSupiRangeIndexs []string
	if c.SupiRangeList != nil {
		index := 0
		for _, item := range c.SupiRangeList {
			if !item.IsValid() {
				invalidSupiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.SupiRangeList, index)
				invalidSupiRangeIndexs = append(invalidSupiRangeIndexs, invalidSupiRangeIndex)
			}
			index++
		}
	}

	return invalidSupiRangeIndexs
}

// GetInvalidGpsiRangeIndexs return invalid gpsiRangeList index
func (c *TChfInfo) GetInvalidGpsiRangeIndexs() []string {
	var invalidGpsiRangeIndexs []string
	if c.GpsiRangeList != nil {
		index := 0
		for _, item := range c.GpsiRangeList {
			if !item.IsValid() {
				invalidGpsiRangeIndex := fmt.Sprintf("%s[%d]", constvalue.GpsiRangeList, index)
				invalidGpsiRangeIndexs = append(invalidGpsiRangeIndexs, invalidGpsiRangeIndex)
			}
			index++
		}
	}

	return invalidGpsiRangeIndexs
}

// GetInvalidPlmnRangeIndexs return invalid plmnRangeList index
func (c *TChfInfo) GetInvalidPlmnRangeIndexs() []string {
	var invalidPlmnRangeIndexs []string
	if c.PlmnRangeList != nil {
		index := 0
		for _, item := range c.PlmnRangeList {
			if !item.IsValid() {
				invalidPlmnRangeIndex := fmt.Sprintf("%s[%d]", constvalue.PlmnRangeList, index)
				invalidPlmnRangeIndexs = append(invalidPlmnRangeIndexs, invalidPlmnRangeIndex)
			}
			index++
		}
	}

	return invalidPlmnRangeIndexs
}

func (c *TChfInfo) createNfInfo() string {
	var supiRanges string
	var gpsiRanges string
	var plmnRanges string
	var supiMatchAll string
	var gpsiMatchAll string
	if c.SupiRangeList != nil && len(c.SupiRangeList) > 0 {
		supiRangeList := ""
		for _, v := range c.SupiRangeList {

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
			supiRanges = fmt.Sprintf(`"supiRangeList":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		} else {
			supiRanges = fmt.Sprintf(`"supiRangeList":[%s]`, supiRangeList)
		}
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		supiRanges = fmt.Sprintf(`"supiRangeList":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		supiMatchAll = fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
	}

	if c.GpsiRangeList != nil && len(c.GpsiRangeList) > 0 {
		gpsiRangeList := ""
		for _, v := range c.GpsiRangeList {

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
			gpsiRanges = fmt.Sprintf(`"gpsiRangeList":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		} else {
			gpsiRanges = fmt.Sprintf(`"gpsiRangeList":[%s]`, gpsiRangeList)
		}
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.NoMatchAll)
	} else {
		gpsiRanges = fmt.Sprintf(`"gpsiRangeList":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		gpsiMatchAll = fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
	}

	if c.PlmnRangeList != nil && len(c.PlmnRangeList) > 0 {

		plmnRangeList := ""
		for _, v := range c.PlmnRangeList {

			plmnRangeItem := ""
			if v.Start != "" {
				start := fmt.Sprintf(`"start":"%s"`, v.Start)
				if plmnRangeItem != "" {
					plmnRangeItem += ","
				}
				plmnRangeItem += start
			}

			if v.End != "" {
				end := fmt.Sprintf(`"end":"%s"`, v.End)
				if plmnRangeItem != "" {
					plmnRangeItem += ","
				}
				plmnRangeItem += end
			}

			if v.Pattern != "" {
				buffer := &bytes.Buffer{}
				encoder := json.NewEncoder(buffer)
				err := encoder.Encode(v.Pattern)
				if nil == err {
					pattern := fmt.Sprintf(`"pattern":%s`, string(buffer.Bytes()))
					if plmnRangeItem != "" {
						plmnRangeItem += ","
					}
					plmnRangeItem += pattern
				}
			}

			if plmnRangeItem != "" {
				if plmnRangeList != "" {
					plmnRangeList += ","
				}
				plmnRangeList += "{" + plmnRangeItem + "}"
			}
		}

		if plmnRangeList == "" {
			plmnRanges = fmt.Sprintf(`"plmnRangeList":[{"pattern":"%s"}]`, constvalue.EmptyPlmnRangePattern)
		} else {
			plmnRanges = fmt.Sprintf(`"plmnRangeList":[%s]`, plmnRangeList)
		}

	} else {
		plmnRanges = fmt.Sprintf(`"plmnRangeList":[{"pattern":"%s"}]`, constvalue.EmptyPlmnRangePattern)
	}

	return fmt.Sprintf(`"chfInfo":{%s,%s,%s,%s,%s}`, supiRanges, gpsiRanges, plmnRanges, supiMatchAll, gpsiMatchAll)
}

// GenerateNfGroupCond generate NfGroupCond
func (c *TChfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (c *TChfInfo) IsEqual(d TNfInfo) bool {

	b := d.(*TChfInfo)

	if len(c.PlmnRangeList) != len(b.PlmnRangeList) {
		return false
	}

	if len(c.GpsiRangeList) != len(b.GpsiRangeList) {
		return false
	}

	if len(c.SupiRangeList) != len(b.SupiRangeList) {
		return false
	}

	for k, item := range c.PlmnRangeList {
		bb := b.PlmnRangeList[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	for k, item := range c.GpsiRangeList {
		bb := b.GpsiRangeList[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	for k, item := range c.SupiRangeList {
		bb := b.SupiRangeList[k]
		if !item.IsEqual(&bb) {
			return false
		}
	}

	return true
}
