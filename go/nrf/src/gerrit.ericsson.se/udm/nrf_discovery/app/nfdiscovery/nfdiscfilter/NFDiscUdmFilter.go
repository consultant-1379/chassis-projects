package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFUDMInfoFilter to process udminfo filter in nfprofile
type NFUDMInfoFilter struct {
}

func (a *NFUDMInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
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

func (a *NFUDMInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFSMFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression

	routingIndicator := queryForm.GetNRFDiscRoutingIndicator()
	if "" != routingIndicator {
		path := getParamSearchPath(constvalue.NfTypeUDM, constvalue.SearchDataRoutingIndic)
		routingIndicatorExpression := buildStringSearchParameter(path, routingIndicator, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, routingIndicatorExpression)
	}

	gpsi := queryForm.GetNRFDiscGspi()
	if "" != gpsi {
		groupidExpression := createGroupIDInstanceIDExpression(gpsi, constvalue.NfTypeUDM, constvalue.SearchDataGpsi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			gpsiExpression := createGpsiExpression(constvalue.NfTypeUDM, constvalue.SearchDataGpsi, gpsi)
			gpsiExpressionAbsence := createGpsiExpressionForAbsence(constvalue.NfTypeUDM, constvalue.SearchDataGpsi)
			var gpsiExpressionList []*MetaExpression
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpression)
			gpsiExpressionList = append(gpsiExpressionList, gpsiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(gpsiExpressionList))
		}
	}

	externalGroupID := queryForm.GetNRFDiscExterGroupID()
	if "" != externalGroupID {
		rangeExpression := createRangeExpression(constvalue.NfTypeUDM, constvalue.SearchDataExterGroupID, externalGroupID, externalGroupID)
		rangeExpressionAbsence := createRangeAbsenceExpression(constvalue.NfTypeUDM, constvalue.SearchDataExterGroupID)
		var rangeExperssionList []*MetaExpression
		rangeExperssionList = append(rangeExperssionList, rangeExpression)
		rangeExperssionList = append(rangeExperssionList, rangeExpressionAbsence)
		metaExpressionList = append(metaExpressionList, buildORExpression(rangeExperssionList))
	}

	supi := queryForm.GetNRFDiscSupiValue()
	if "" != supi {
		groupidExpression := createGroupIDInstanceIDExpression(supi, constvalue.NfTypeUDM, constvalue.SearchDataSupi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			supiExpression := createSupiExpression(constvalue.NfTypeUDM, constvalue.SearchDataSupi, supi)
			supiExpressionAbsence := createSupiExpressionForAbsence(constvalue.NfTypeUDM, constvalue.SearchDataSupi)
			var supiExpressionList []*MetaExpression
			supiExpressionList = append(supiExpressionList, supiExpression)
			supiExpressionList = append(supiExpressionList, supiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(supiExpressionList))
		}
	}

	groupIDList := queryForm.GetNRFDiscGroupIDList()
	if groupIDList != nil && len(groupIDList) > 0 {
		var groupIDExpressionList []*MetaExpression
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(constvalue.NfTypeUDM, constvalue.SearchDataGroupIDList, groupIDList))

		var absenceExpressionList []*MetaExpression
		absenceExpressionList = append(absenceExpressionList, createSupiExpressionForAbsence(constvalue.NfTypeUDM, constvalue.SearchDataSupi))
		absenceExpressionList = append(absenceExpressionList, createGpsiExpressionForAbsence(constvalue.NfTypeUDM, constvalue.SearchDataGpsi))
		groupIDExpressionList = append(groupIDExpressionList, buildAndExpression(absenceExpressionList))

		metaExpressionList = append(metaExpressionList, buildORExpression(groupIDExpressionList))
	}

	return metaExpressionList
}
