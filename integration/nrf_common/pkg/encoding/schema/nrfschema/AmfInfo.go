package nrfschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"com/dbproxy/nfmessage/subscription"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidateN2InterfaceAmfInfoIndex return invalid n2InterfaceAmfInfo index
func (a *TAmfInfo) GetInvalidN2InterfaceAmfInfoIndex() string {
	if a.N2InterfaceAmfInfo != nil && !a.N2InterfaceAmfInfo.IsValid() {
		return constvalue.N2InterfaceAmfInfo
	}

	return ""
}

// GetInvalidTaiRangeIndexs return invalid taiRangeList index
func (a *TAmfInfo) GetInvalidTaiRangeIndexs() []string {
	var invalidTaiRangeIndexs []string

	if a.TaiRangeList != nil {
		index := 0
		for _, taiRange := range a.TaiRangeList {
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

func (a *TAmfInfo) createNfInfo() string {
	var taiList string
	var taiRangeList string
	var guamiList string
	var backupInfoAmfFailure string
	var backupInfoAmfRemoval string
	if nil != a.TaiList && len(a.TaiList) > 0 {

		taiListStr := ""
		for _, tai := range a.TaiList {
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

	if nil != a.TaiRangeList && len(a.TaiRangeList) > 0 {

		taiRangeListStr := ""
		for _, taiRange := range a.TaiRangeList {
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

	if nil != a.GuamiList && len(a.GuamiList) > 0 {
		guamiListStr := ""
		for _, guami := range a.GuamiList {
			plmnId, _ := json.Marshal(guami.PlmnId)
			amfId := strings.ToLower(guami.AmfId)

			guamiStr := fmt.Sprintf(`{"plmnId": %s, "amfId": "%s"}`, plmnId, amfId)
			if guamiListStr != "" {
				guamiListStr += ","
			}
			guamiListStr += guamiStr
		}

		guamiList = fmt.Sprintf(`"guamiList": [%s]`, guamiListStr)
	}

	if nil != a.BackupInfoAmfFailure && len(a.BackupInfoAmfFailure) > 0 {
		guamiListStr := ""
		for _, guami := range a.BackupInfoAmfFailure {
			plmnID, _ := json.Marshal(guami.PlmnId)
			amfID := strings.ToLower(guami.AmfId)

			guamiStr := fmt.Sprintf(`{"plmnId": %s, "amfId": "%s"}`, plmnID, amfID)
			if guamiListStr != "" {
				guamiListStr += ","
			}
			guamiListStr += guamiStr
		}

		backupInfoAmfFailure = fmt.Sprintf(`"backupInfoAmfFailure": [%s]`, guamiListStr)
	} else {
		backupInfoAmfFailure = fmt.Sprintf(`"backupInfoAmfFailure": [{"plmnId": {"mcc": "%s"}}]`, constvalue.EmptyMcc)
	}

	if nil != a.BackupInfoAmfRemoval && len(a.BackupInfoAmfRemoval) > 0 {
		guamiListStr := ""
		for _, guami := range a.BackupInfoAmfRemoval {
			plmnID, _ := json.Marshal(guami.PlmnId)
			amfID := strings.ToLower(guami.AmfId)

			guamiStr := fmt.Sprintf(`{"plmnId": %s, "amfId": "%s"}`, plmnID, amfID)
			if guamiListStr != "" {
				guamiListStr += ","
			}
			guamiListStr += guamiStr
		}

		backupInfoAmfRemoval = fmt.Sprintf(`"backupInfoAmfRemoval": [%s]`, guamiListStr)
	} else {
		backupInfoAmfRemoval = fmt.Sprintf(`"backupInfoAmfRemoval": [{"plmnId": {"mcc": "%s"}}]`, constvalue.EmptyMcc)
	}

	if guamiList != "" {
		if backupInfoAmfFailure != "" {
			guamiList = guamiList + "," + backupInfoAmfFailure
		}

		if backupInfoAmfRemoval != "" {
			guamiList = guamiList + "," + backupInfoAmfRemoval
		}

	} else {
		if backupInfoAmfFailure != "" {
			guamiList = backupInfoAmfFailure
		}

		if backupInfoAmfRemoval != "" {
			if guamiList != "" {
				guamiList = guamiList + "," + backupInfoAmfRemoval
			} else {
				guamiList = backupInfoAmfRemoval
			}
		}
	}
	var amfInfo string
	if guamiList != "" {
		amfInfo = fmt.Sprintf(`%s,%s, %s`, taiList, taiRangeList, guamiList)
	} else {
		amfInfo = fmt.Sprintf(`%s,%s`, taiList, taiRangeList)
	}
	if a.AmfRegionId != "" {
		amfInfo = fmt.Sprintf(`%s,"amfRegionId":"%s"`, amfInfo, a.AmfRegionId)
	}
	if a.AmfSetId != "" {
		amfInfo = fmt.Sprintf(`%s,"amfSetId":"%s"`, amfInfo, a.AmfSetId)
	}
	return fmt.Sprintf(`"amfInfo":{%s}`, amfInfo)
}

//GenerateAmfCond generate AmfCond for subscription
func (a *TAmfInfo) GenerateAmfCond() []*subscription.SubKeyStruct {
	var subKeys []*subscription.SubKeyStruct

	if a.AmfSetId != "" {
		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: a.AmfSetId,
			SubKey2: constvalue.Wildcard,
		})
	}

	if a.AmfRegionId != "" {
		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: constvalue.Wildcard,
			SubKey2: a.AmfRegionId,
		})
	}

	if a.AmfSetId != "" && a.AmfRegionId != "" {
		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: a.AmfSetId,
			SubKey2: a.AmfRegionId,
		})
	}

	return subKeys
}

// GenerateGuamiListCond generate GuamiListCond for subscription
func (a *TAmfInfo) GenerateGuamiListCond() []*subscription.SubKeyStruct {
	var subKeys []*subscription.SubKeyStruct

	if a.GuamiList != nil {
		for _, item := range a.GuamiList {
			subKeys = append(subKeys, item.GenerateGrpcKey())
		}
	}

	return subKeys
}

// GenerateNfGroupCond generate NfGroupCond
func (a *TAmfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (a *TAmfInfo) IsEqual(c TNfInfo) bool {

	b := c.(*TAmfInfo)

	if a.AmfRegionId != b.AmfRegionId {
		return false
	}

	if a.AmfSetId != b.AmfSetId {
		return false
	}

	if (a.N2InterfaceAmfInfo != nil && b.N2InterfaceAmfInfo == nil) || (a.N2InterfaceAmfInfo == nil && b.N2InterfaceAmfInfo != nil) {
		return false
	}

	if len(a.BackupInfoAmfFailure) != len(b.BackupInfoAmfFailure) {
		return false
	}

	if len(a.BackupInfoAmfRemoval) != len(b.BackupInfoAmfRemoval) {
		return false
	}

	if len(a.GuamiList) != len(b.GuamiList) {
		return false
	}

	if len(a.TaiList) != len(b.TaiList) {
		return false
	}

	if len(a.TaiRangeList) != len(b.TaiRangeList) {
		return false
	}

	if a.N2InterfaceAmfInfo != nil && b.N2InterfaceAmfInfo != nil {

		if len(a.N2InterfaceAmfInfo.Ipv4EndpointAddress) != len(b.N2InterfaceAmfInfo.Ipv4EndpointAddress) {
			return false
		}

		if len(a.N2InterfaceAmfInfo.Ipv6EndpointAddress) != len(b.N2InterfaceAmfInfo.Ipv6EndpointAddress) {
			return false
		}

		if a.N2InterfaceAmfInfo.AmfName != b.N2InterfaceAmfInfo.AmfName {
			return false
		}

		for k, ipv4 := range a.N2InterfaceAmfInfo.Ipv4EndpointAddress {
			if ipv4 != b.N2InterfaceAmfInfo.Ipv4EndpointAddress[k] {
				return false
			}
		}
		for k, ipv6 := range a.N2InterfaceAmfInfo.Ipv6EndpointAddress {
			if ipv6 != b.N2InterfaceAmfInfo.Ipv6EndpointAddress[k] {
				return false
			}
		}
	}

	for k, item := range a.BackupInfoAmfFailure {
		bb := b.BackupInfoAmfFailure[k]
		if item.AmfId != bb.AmfId || item.PlmnId.Mcc != bb.PlmnId.Mcc || item.PlmnId.Mnc != bb.PlmnId.Mnc {
			return false
		}
	}

	for k, item := range a.BackupInfoAmfRemoval {
		bb := b.BackupInfoAmfRemoval[k]
		if item.AmfId != bb.AmfId || item.PlmnId.Mcc != bb.PlmnId.Mcc || item.PlmnId.Mnc != bb.PlmnId.Mnc {
			return false
		}
	}

	for k, item := range a.GuamiList {
		bb := b.GuamiList[k]
		if item.AmfId != bb.AmfId || item.PlmnId.Mcc != bb.PlmnId.Mcc || item.PlmnId.Mnc != bb.PlmnId.Mnc {
			return false
		}
	}

	for k, item := range a.TaiList {
		bb := b.TaiList[k]
		if item.Tac != bb.Tac || item.PlmnId.Mcc != bb.PlmnId.Mcc || item.PlmnId.Mnc != bb.PlmnId.Mnc {
			return false
		}
	}

	for k, item := range a.TaiRangeList {
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
			if inItem.Pattern != bbb.Pattern || inItem.Start != bbb.Start || inItem.End != bbb.End {
				return false
			}
		}
	}

	return true
}
