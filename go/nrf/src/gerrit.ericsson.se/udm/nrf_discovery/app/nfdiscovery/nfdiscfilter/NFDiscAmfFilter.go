package nfdiscfilter

import (
	"strings"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
)
//NFAMFInfoFilter to process amfinfo filter in nfprofile
type NFAMFInfoFilter struct {

}

func (a *NFAMFInfoFilter) filter(nfInfo []byte, queryFrom *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	return true

}

func (a *NFAMFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFAMFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	amfSetID := queryForm.GetNRFDiscAMFSetID()
	if "" != amfSetID {
		amfSetIDPath := getParamSearchPath(constvalue.NfTypeAMF, constvalue.SearchDataAmfSetID)
		amfSetIDExpression := buildStringSearchParameter(amfSetIDPath, amfSetID, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, amfSetIDExpression)
	}

	amfRegionID := queryForm.GetNRFDiscAMFRegionID()
	if "" != amfRegionID {
		amfRegionIDPath := getParamSearchPath(constvalue.NfTypeAMF, constvalue.SearchDataAmfRegionID)
		amfRegionIDExpression := buildStringSearchParameter(amfRegionIDPath, amfRegionID, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, amfRegionIDExpression)
	}

	plmnid, amfid := queryForm.GetNRFDiscGuamiType()
	if len(plmnid) >= 5 && amfid != "" {
		amfid = strings.ToLower(amfid)
		metaExpressionList = append(metaExpressionList, createGuamiFilter(constvalue.NfTypeAMF, constvalue.SearchDataGuami, plmnid, amfid))
	}

	plmnid, tac := queryForm.GetNRFDiscTaiType()
	if len(plmnid) >= 5 && tac != "" {
		tac = strings.ToLower(tac)
		//metaExpressionList = append(metaExpressionList, createTaiExpression("body.amfInfo", plmnid, tac))
		taiNormal := createTaiExpression(constvalue.NfTypeAMF, constvalue.SearchDataTai, plmnid, tac)
		taiAbsence := createTaiExpressionForAbsence(constvalue.NfTypeAMF, constvalue.SearchDataTai)

		var taiExpressionList []*MetaExpression
		taiExpressionList = append(taiExpressionList, taiNormal)
		taiExpressionList = append(taiExpressionList, taiAbsence)

		metaExpressionList = append(metaExpressionList, buildORExpression(taiExpressionList))
	}

	return metaExpressionList
}