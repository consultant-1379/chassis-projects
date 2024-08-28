package nfdiscfilter

import (
	"regexp"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"github.com/buger/jsonparser"
)

//NFCHFInfoFilter to process upfinfo filter in nfprofile
type NFCHFInfoFilter struct {
}

//filter is to filter chfInfo in nfprofile
func (a *NFCHFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	plmnID, _ := queryForm.GetNRFDiscPlmnValue(constvalue.SearchDataChfSupportedPlmn)
	if plmnID != "" {
		log.Debugf("Search nfProfile with chf-supported-plmn %s", plmnID)
		if !a.isMatchedChfSupportedPlmn(plmnID, nfInfo) {
			log.Debugf("No Matched nfProfile with chf-supported-plmn %s", plmnID)
			return false
		}
	}

	if !filterInfo.KVDBSearch && queryForm.GetNRFDiscGspi() != "" {
		log.Debugf("Search nfProfile with gpsi: %s", queryForm.GetNRFDiscGspi())
		if !isMatchedGpsi(queryForm, nfInfo) {
			log.Debugf("No Matched nfProfile with gpsi: %s", queryForm.GetNRFDiscGspi())
			return false
		}

	}

	if !filterInfo.KVDBSearch && queryForm.GetNRFDiscSupiValue() != "" {
		log.Debugf("Search nfProfile with supi %s", queryForm.GetNRFDiscSupiValue())
		matchResult := isMatchedGroupID(queryForm, filterInfo.groupID, nfInfo)
		if !(matchResult == nfdiscutil.ResultFoundMatch) || isMatchedSupi(queryForm, nfInfo) {
			log.Debugf("No Matched nfProfile with supi %s", queryForm.GetNRFDiscSupiValue())
			return false
		}
	}

	return true
}

//filterByKVDB is to generate kvdb search parameters
func (a *NFCHFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFCHFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression

	gpsi := queryForm.GetNRFDiscGspi()
	if "" != gpsi {
		groupidExpression := createGroupIDInstanceIDExpression(gpsi, constvalue.NfTypeCHF, constvalue.SearchDataGpsi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			gpsiExpression := createGpsiExpression(constvalue.NfTypeCHF, constvalue.SearchDataGpsi, gpsi)
			gpsiExpressionAbsence := createGpsiExpressionForAbsence(constvalue.NfTypeCHF, constvalue.SearchDataGpsi)
			var gpsiExpressionList []*MetaExpression
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpression)
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(gpsiExpressionList))
		}
	}

	supi := queryForm.GetNRFDiscSupiValue()
	if "" != supi {
		groupidExpression := createGroupIDInstanceIDExpression(gpsi, constvalue.NfTypeCHF, constvalue.SearchDataSupi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			supiExpression := createSupiExpression(constvalue.NfTypeCHF, constvalue.SearchDataSupi, supi)
			supiExpressionAbsence := createSupiExpressionForAbsence(constvalue.NfTypeCHF, constvalue.SearchDataSupi)
			var supiExpressionList []*MetaExpression
			supiExpressionList = append(supiExpressionList, supiExpression)
			supiExpressionList = append(supiExpressionList, supiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(supiExpressionList))
		}
	}

	//plmnID, _ := queryForm.GetNRFDiscPlmnValue(constvalue.SearchDataChfSupportedPlmn)
	//if plmnID != "" && queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) == constvalue.NfTypeCHF {
	//	chfSupportedPlmnExpression := createRangeExpression(constvalue.NfTypeCHF, constvalue.SearchDataChfSupportedPlmn, plmnID, plmnID)
	//	plmnRangeExpressionAbsence := createPlmnExpressionForAbsence(constvalue.NfTypeCHF, constvalue.SearchDataChfSupportedPlmn)
	//	var rangeExperssionList []*MetaExpression
	//	rangeExperssionList = append(rangeExperssionList, chfSupportedPlmnExpression, plmnRangeExpressionAbsence)
	//	metaExpressionList = append(metaExpressionList, buildORExpression(rangeExperssionList))
	//}

	return metaExpressionList
}

//mcc match
func (a *NFCHFInfoFilter) isMccMatched(start, end, plmnid string )bool {
	s, err1 := strconv.ParseInt(start[0:3], 10, 64)
	if err1 != nil {
		log.Debugf("mcc start parseint error, err=%v", err1)
		return false
	}
	e, err2 := strconv.ParseInt(end[0:3], 10, 64)
	if err2 != nil {
		log.Debugf("mcc end parseint error, err=%v", err2)
		return false
	}
	mcc, err3 := strconv.ParseInt(plmnid[0:3], 10, 64)
	if err3 != nil {
		log.Debugf("mcc plmnid parseint error, err=%v", err3)
		return false
	}

	if mcc >= s && mcc <= e {
		return true
	}

	return false
}
//mnc match
func (a *NFCHFInfoFilter) isMncMatched(start, end, plmnid string) bool {
	if len(start) > len(plmnid) || len(end) < len(plmnid){
		return false
	}
	s, err1 := strconv.ParseInt(start[3:], 10, 64)
	if err1 != nil {
		log.Debugf("mnc start parseint error, err=%v", err1)
		return false
	}
	e, err2 := strconv.ParseInt(end[3:], 10, 64)
	if err2 != nil {
		log.Debugf("mnc end parseint error, err=%v", err2)
		return false
	}
	mnc, err3 := strconv.ParseInt(plmnid[3:], 10, 64)
	if err3 != nil {
		log.Debugf("mnc plmnid parseint error, err=%v", err3)
		return false
	}

	if mnc >= s && mnc <= e {
		return true
	}

	return false
}
//isMatchedChfSupportedPlmn is to match chf-supported-plmn in chfInfo
func (a *NFCHFInfoFilter) isMatchedChfSupportedPlmn(plmnID string, nfInfo []byte) bool {
	ret := false
	num := 0
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		num = num + 1
		if ret {
			return
		}
		pattern, err1 := jsonparser.GetString(value, "pattern")
		s, err2 := jsonparser.GetString(value, "start")
		e, err3 := jsonparser.GetString(value, "end")
		if err1 != nil && err2 != nil && err3 != nil {
			log.Debugf("plmnRangeList not have start & end & pattern, match each supi")
			ret = true
			return
		}

		if err1 == nil {
			matched, matchErr := regexp.MatchString(pattern, plmnID)
			if matchErr != nil {
				log.Debugf("plmnId regexp match error, err=%v", matchErr)
			}
			log.Debugf("The chf-supported-plmn: %s, pattern : %s, matched result: %v", plmnID, pattern, matched)
			if matched {
				ret = true
				return
			}
		}

		if err2 != nil || err3 != nil {
			return
		}

                if a.isMccMatched(s, e, plmnID) && a. isMncMatched(s, e, plmnID) {
			ret = true
			log.Debugf("the chf-supported-plmn is %s, range is %s-%s matched result is %v", plmnID, s, e, ret)
			return
		}
		log.Debugf("the chf-supported-plmn is %s, range is %s-%s matched result is %v", plmnID, s, e, ret)

	}, constvalue.PlmnRangeList)

	if err != nil || ret == true {
		return true
	}

	if num == 0 && err == nil {
		log.Debugf("plmnRangeList is [], this will match all chf-supported-plmn")
		return true
	}
	return ret
}
