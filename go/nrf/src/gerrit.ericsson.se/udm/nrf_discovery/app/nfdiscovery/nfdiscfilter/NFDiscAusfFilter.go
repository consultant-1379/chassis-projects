package nfdiscfilter

import (

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFAUSFInfoFilter to process ausfinfo filter in nfprofile
type NFAUSFInfoFilter struct {
}

func (a *NFAUSFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
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

func (a *NFAUSFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFAUSFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression

	routingIndicator := queryForm.GetNRFDiscRoutingIndicator()
	if "" != routingIndicator {
		routingIndicatorPath := getParamSearchPath(constvalue.NfTypeAUSF, constvalue.SearchDataRoutingIndic)
		routingIndicatorExpression := buildStringSearchParameter(routingIndicatorPath, routingIndicator, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, routingIndicatorExpression)
	}

	supi := queryForm.GetNRFDiscSupiValue()
	if "" != supi {
		groupidExpression := createGroupIDInstanceIDExpression(supi, constvalue.NfTypeAUSF, constvalue.SearchDataSupi)

		if nil != groupidExpression {
			metaExpressionList = append(metaExpressionList, groupidExpression)

		} else {
			supiExpression := createSupiExpression(constvalue.NfTypeAUSF, constvalue.SearchDataSupi, supi)
			supiExpressionAbsence := createSupiExpressionForAbsence(constvalue.NfTypeAUSF, constvalue.SearchDataSupi)
			var supiExpressionList []*MetaExpression
			supiExpressionList = append(supiExpressionList, supiExpression)
			supiExpressionList = append(supiExpressionList, supiExpressionAbsence)
			metaExpressionList = append(metaExpressionList, buildORExpression(supiExpressionList))
		}
	}

	groupIDList := queryForm.GetNRFDiscGroupIDList()
	if groupIDList != nil && len(groupIDList) > 0 {
		var groupIDExpressionList []*MetaExpression
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(constvalue.NfTypeAUSF, constvalue.SearchDataGroupIDList, groupIDList))
		groupIDExpressionList = append(groupIDExpressionList, createSupiExpressionForAbsence(constvalue.NfTypeAUSF, constvalue.SearchDataSupi))

		metaExpressionList = append(metaExpressionList, buildORExpression(groupIDExpressionList))
	}
	return metaExpressionList
}
