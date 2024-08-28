package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"github.com/buger/jsonparser"
)

//NFUDRInfoFilter to process udrinf filter in nfprofile
type NFUDRInfoFilter struct {
}

func (a *NFUDRInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	if !filterInfo.KVDBSearch && "" != queryForm.GetNRFDiscDataSet() {
		log.Debugf("Search nfProfile with DataSet: %s", queryForm.GetNRFDiscDataSet())
		if !a.isMatchedDataSet(queryForm.GetNRFDiscDataSet(), nfInfo) {
			log.Debugf("No Matched nfProfile with DataSet: %s", queryForm.GetNRFDiscDataSet())
			return false
		}
	}

	if !filterInfo.KVDBSearch && queryForm.GetNRFDiscExterGroupID() != "" {
		log.Debugf("Search nfProfile with externalGroupID: %s", queryForm.GetNRFDiscExterGroupID())
		if !isMatchedExternalGroupID(queryForm, nfInfo) {
			log.Debugf("No Matched nfProfile with externalGroupID: %s", queryForm.GetNRFDiscExterGroupID())
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

func (a *NFUDRInfoFilter) isMatchedDataSet(dataSet string, nfInfo []byte) bool {
	matched := false
	log.Debugf("nfProfile: %s", string(nfInfo))
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if matched {
			return
		}
		log.Debugf("dataset value: %s", string(value[:]))
		dataSetInProfile := string(value[:])
		if dataSet == dataSetInProfile {
			matched = true
			return
		}
	}, constvalue.SupportedDataSets)

	if err != nil {
		matched = false
	}

	return matched
}

func (a *NFUDRInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFUDRInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression

	dataSet := queryForm.GetNRFDiscDataSet()
	if "" != dataSet {
		dataSetPath := getParamSearchPath(constvalue.NfTypeUDR, constvalue.SearchDataDataSet)
		dataSetExpression := buildStringSearchParameter(dataSetPath, dataSet, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, dataSetExpression)
	}

	gpsi := queryForm.GetNRFDiscGspi()
	if "" != gpsi {
		groupidExpression := createGroupIDInstanceIDExpression(gpsi, constvalue.NfTypeUDR, constvalue.SearchDataGpsi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			gpsiExpression := createGpsiExpression(constvalue.NfTypeUDR, constvalue.SearchDataGpsi, gpsi)
			gpsiExpressionAbsence := createGpsiExpressionForAbsence(constvalue.NfTypeUDR, constvalue.SearchDataGpsi)
			var gpsiExpressionList []*MetaExpression
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpression)
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(gpsiExpressionList))
		}
	}

	externalGroupID := queryForm.GetNRFDiscExterGroupID()
	if "" != externalGroupID {
		rangeExpression := createRangeExpression(constvalue.NfTypeUDR, constvalue.SearchDataExterGroupID, externalGroupID, externalGroupID)
		rangeExpressionAbsence := createRangeAbsenceExpression(constvalue.NfTypeUDR, constvalue.SearchDataExterGroupID)
		var rangeExperssionList []*MetaExpression
		rangeExperssionList = append(rangeExperssionList, rangeExpression)
		rangeExperssionList = append(rangeExperssionList, rangeExpressionAbsence)
		metaExpressionList = append(metaExpressionList, buildORExpression(rangeExperssionList))
	}

	supi := queryForm.GetNRFDiscSupiValue()
	if "" != supi {
		groupidExpression := createGroupIDInstanceIDExpression(supi, constvalue.NfTypeUDR, constvalue.SearchDataSupi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		}else {
			supiExpression := createSupiExpression(constvalue.NfTypeUDR, constvalue.SearchDataSupi, supi)
			supiExpressionAbsence := createSupiExpressionForAbsence(constvalue.NfTypeUDR, constvalue.SearchDataSupi)
			var supiExpressionList []*MetaExpression
			supiExpressionList = append(supiExpressionList, supiExpression)
			supiExpressionList = append(supiExpressionList, supiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(supiExpressionList))
		}
	}

	groupIDList := queryForm.GetNRFDiscGroupIDList()
	if groupIDList != nil && len(groupIDList) > 0 {
		var groupIDExpressionList []*MetaExpression
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(constvalue.NfTypeUDR, constvalue.SearchDataGroupIDList, groupIDList))

		var absenceExpressionList []*MetaExpression
		absenceExpressionList = append(absenceExpressionList, createSupiExpressionForAbsence(constvalue.NfTypeUDR, constvalue.SearchDataSupi))
		absenceExpressionList = append(absenceExpressionList, createGpsiExpressionForAbsence(constvalue.NfTypeUDR, constvalue.SearchDataGpsi))
                groupIDExpressionList = append(groupIDExpressionList, buildAndExpression(absenceExpressionList))

		metaExpressionList = append(metaExpressionList, buildORExpression(groupIDExpressionList))
	}

	return metaExpressionList
}
