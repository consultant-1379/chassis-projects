package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFPCFInfoFilter to process pcfinfo filter in nfprofile
type NFPCFInfoFilter struct {
}

func (a *NFPCFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
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

func (a *NFPCFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFAMFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	dnn := queryForm.GetNRFDiscDnnValue()
	if "" != dnn {
		dnnPath := getParamSearchPath(constvalue.NfTypePCF, constvalue.SearchDataDnn)
		dnnExpression := buildStringSearchParameter(dnnPath, dnn, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, dnnExpression)
	}

	supi := queryForm.GetNRFDiscSupiValue()
	if "" != supi {
		groupidExpression := createGroupIDInstanceIDExpression(supi, constvalue.NfTypePCF, constvalue.SearchDataSupi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			supiExpression := createSupiExpression(constvalue.NfTypePCF, constvalue.SearchDataSupi, supi)
			supiExpressionAbsence := createSupiExpressionForAbsence(constvalue.NfTypePCF, constvalue.SearchDataSupi)
			var supiExpressionList []*MetaExpression
			supiExpressionList = append(supiExpressionList, supiExpression)
			supiExpressionList = append(supiExpressionList, supiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(supiExpressionList))
		}
	}
        
	groupIDList := queryForm.GetNRFDiscGroupIDList()
	if groupIDList != nil && len(groupIDList) > 0 {
		var groupIDExpressionList []*MetaExpression
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(constvalue.NfTypePCF, constvalue.SearchDataGroupIDList, groupIDList))

		//var absenceExpressionList []*common.MetaExpression
		//absenceExpressionList = append(absenceExpressionList, createSupiExpressionForAbsence(constvalue.NfTypePCF, constvalue.SearchDataSupi))
		//absenceExpressionList = append(absenceExpressionList, createGpsiExpressionForAbsence(constvalue.NfTypePCF, constvalue.SearchDataGpsi))
		groupIDExpressionList = append(groupIDExpressionList, createSupiExpressionForAbsence(constvalue.NfTypePCF, constvalue.SearchDataSupi))


		metaExpressionList = append(metaExpressionList, buildORExpression(groupIDExpressionList))
	}

	return metaExpressionList
}
