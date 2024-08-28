package nrfschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"com/dbproxy/nfmessage/subscription"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidTaiRangeIndexs return invalid taiRangeList index
func (s *TSmfInfo) GetInvalidTaiRangeIndexs() []string {
	var invalidTaiRangeIndexs []string

	if s.TaiRangeList != nil {
		index := 0
		for _, taiRange := range s.TaiRangeList {
			invalidTacRangeIndexs := taiRange.GetInvalidTacRangeIndexs()
			if invalidTacRangeIndexs != nil {
				for _, invalidTacRangeIndex := range invalidTacRangeIndexs {
					invalidTaiRangeIndex := fmt.Sprintf("%s[%d].%s", constvalue.TaiRangeList, index, invalidTacRangeIndex)
					invalidTaiRangeIndexs = append(invalidTaiRangeIndexs, invalidTaiRangeIndex)
				}
			}
			index++
		}
	}

	return invalidTaiRangeIndexs
}

func (s *TSmfInfo) createNfInfo() string {
	var taiList string
	var taiRangeList string

	if nil != s.TaiList && len(s.TaiList) > 0 {

		taiListStr := ""
		for _, tai := range s.TaiList {
			plmnId := tai.PlmnId
			tac := strings.ToLower(tai.Tac)
			if tac == "" {
				tac = constvalue.EmptyTac
			}

			plmnIdStr, _ := json.Marshal(plmnId)
			taiStr := fmt.Sprintf(`{"plmnId":%s, "tac": "%s"}`, plmnIdStr, tac)
			if taiListStr != "" {
				taiListStr += ","
			}
			taiListStr += taiStr
		}
		taiList = fmt.Sprintf(`"taiList":[%s]`, taiListStr)

	} else {

		plmnid := fmt.Sprintf(`{"mcc":"%s","mnc":"%s"}`, constvalue.EmptyMcc, constvalue.EmptyMnc)
		taiList = fmt.Sprintf(`"taiList":[{"plmnId":%s, "tac": "%s"}]`, plmnid, constvalue.EmptyTac)
	}

	if nil != s.TaiRangeList && len(s.TaiRangeList) > 0 {

		taiRangeListStr := ""
		for _, taiRange := range s.TaiRangeList {
			tacRangeListStr := ""
			if nil != taiRange.TacRangeList && len(taiRange.TacRangeList) > 0 {
				for _, tacRange := range taiRange.TacRangeList {
					tacRangeStr := ""
					if tacRange.Start != "" && tacRange.End != "" {
						tacRangeStr = fmt.Sprintf(`{"start": "%s", "end": "%s"}`, strings.ToLower(tacRange.Start), strings.ToLower(tacRange.End))
					} else if tacRange.Pattern != "" {
						tacRangeStr = fmt.Sprintf(`{"pattern": "%s"}`, tacRange.Pattern)
					} else {
						tacRangeStr = fmt.Sprintf(`{"pattern": "%s"}`, constvalue.EmptyTacRangePattern)
					}

					if tacRangeListStr != "" {
						tacRangeListStr += ","
					}
					tacRangeListStr += tacRangeStr
				}
			} else {
				tacRangeListStr = fmt.Sprintf(`{"pattern": "%s"}`, constvalue.EmptyTac)
			}
			plmnIdStr, _ := json.Marshal(taiRange.PlmnID)

			taiRangeStr := fmt.Sprintf(`{"plmnId": %s, "tacRangeList":[%s]}`, plmnIdStr, tacRangeListStr)

			if taiRangeListStr != "" {
				taiRangeListStr += ","
			}
			taiRangeListStr += taiRangeStr
		}
		taiRangeList = fmt.Sprintf(`"taiRangeList":[%s]`, taiRangeListStr)

	} else {

		plmnid := fmt.Sprintf(`{"mcc":"%s","mnc":"%s"}`, constvalue.EmptyMcc, constvalue.EmptyMnc)
		tacRange := fmt.Sprintf(`{"pattern":"%s"}`, constvalue.EmptyTacRangePattern)
		taiRangeList = fmt.Sprintf(`"taiRangeList":[{"plmnId":%s, "tacRangeList": [%s]}]`, plmnid, tacRange)
	}

	var sNssaiSmfInfoList string
	if s.SNssaiSmfInfoList != nil && len(s.SNssaiSmfInfoList) > 0 {
		sNssaiDnnList := ""
		for _, item := range s.SNssaiSmfInfoList {
			sst := item.SNssai.Sst
			sd := strings.ToLower(item.SNssai.Sd)
			if sd == "" {
				sd = constvalue.EmptySd
			}
			sNssaisList := fmt.Sprintf(`"sNssai":{"sst":%d,"sd":"%s"}`, sst, sd)
			dnnSmfInfoList, _ := json.Marshal(item.DnnSmfInfoList)
			if sNssaiDnnList != "" {
				sNssaiDnnList += fmt.Sprintf(`,{%s,"dnnSmfInfoList":%s}`, sNssaisList, dnnSmfInfoList)
			} else {
				sNssaiDnnList += fmt.Sprintf(`{%s,"dnnSmfInfoList":%s}`, sNssaisList, dnnSmfInfoList)
			}
		}
		sNssaiSmfInfoList = fmt.Sprintf(`"sNssaiSmfInfoList":[%s]`, sNssaiDnnList)
	}
	if s.PgwFqdn != "" {
		return fmt.Sprintf(`"smfInfo":{%s,%s, %s,"pgwFqdn":"%s"}`, taiList, taiRangeList, sNssaiSmfInfoList, s.PgwFqdn)
	}
	return fmt.Sprintf(`"smfInfo":{%s,%s, %s}`, taiList, taiRangeList, sNssaiSmfInfoList)
}

// GenerateNfGroupCond generate NfGroupCond
func (s *TSmfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (s *TSmfInfo) IsEqual(c TNfInfo) bool {

	b := c.(*TSmfInfo)

	if len(s.PgwFqdn) != len(b.PgwFqdn) {
		return false
	}

	if len(s.AccessType) != len(b.AccessType) {
		return false
	}

	if len(s.TaiList) != len(b.TaiList) {
		return false
	}

	if len(s.SNssaiSmfInfoList) != len(b.SNssaiSmfInfoList) {
		return false
	}

	if len(s.TaiRangeList) != len(b.TaiRangeList) {
		return false
	}

	for k, item := range s.AccessType {
		if item != b.AccessType[k] {
			return false
		}
	}

	for k, item := range s.TaiList {
		bb := b.TaiList[k]
		if item.Tac != bb.Tac || item.PlmnId.Mcc != bb.PlmnId.Mcc || item.PlmnId.Mnc != bb.PlmnId.Mnc {
			return false
		}
	}

	for k, item := range s.SNssaiSmfInfoList {
		bb := b.SNssaiSmfInfoList[k]
		if item.SNssai.Sst != bb.SNssai.Sst || item.SNssai.Sd != bb.SNssai.Sd {
			return false
		}

		if len(item.DnnSmfInfoList) != len(bb.DnnSmfInfoList) {
			return false
		}
		for j, inItem := range item.DnnSmfInfoList {
			if inItem.Dnn != bb.DnnSmfInfoList[j].Dnn {
				return false
			}
		}

	}

	for k, item := range s.TaiRangeList {
		bb := b.TaiRangeList[k]

		if (item.PlmnID != nil && bb.PlmnID == nil) || (item.PlmnID == nil && bb.PlmnID != nil) {
			return false
		} else if item.PlmnID != nil && bb.PlmnID != nil {

			if item.PlmnID.Mcc != bb.PlmnID.Mcc || item.PlmnID.Mnc != bb.PlmnID.Mnc {
				return false
			}
		}

		if len(item.TacRangeList) != len(bb.TacRangeList) {
			return false
		}
		for j, inItem := range item.TacRangeList {
			bbb := bb.TacRangeList[j]
			if !inItem.IsEqual(&bbb) {
				return false
			}
		}

	}

	return true
}
