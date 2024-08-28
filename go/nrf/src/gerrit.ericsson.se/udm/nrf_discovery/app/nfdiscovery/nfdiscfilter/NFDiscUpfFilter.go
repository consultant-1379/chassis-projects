package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"github.com/buger/jsonparser"
)

//NFUPFInfoFilter to process upfinfo filter in nfprofile
type NFUPFInfoFilter struct {

}

func (a *NFUPFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	upfIwkEpsInd, err := queryForm.GetNRFDiscBoolValue(constvalue.SearchDataUpfIwkEpsInd)
	if err == nil {
		log.Debugf("Search nfProfile with upf-iwk-eps-ind = %v", upfIwkEpsInd)
		if !a.isMatchedUpfIwkEpsInd(upfIwkEpsInd, nfInfo) {
			log.Debugf("No Matched nfProfile with upf-iwk-eps-ind = %v", upfIwkEpsInd)
			return false
		}
	}
	return true
}

//isMatchedUpfIwkEpsInd is to match upf-iwk-eps-ind in upfInfo
func (a *NFUPFInfoFilter)isMatchedUpfIwkEpsInd(upfIwkEpsInd bool, nfInfo []byte) bool {
	matched := false
	upfIwkEpsIndInProfile, err := jsonparser.GetBoolean(nfInfo, constvalue.IwkEpsInd)
	if err == nil {
		if upfIwkEpsInd == upfIwkEpsIndInProfile {
			matched = true
		}
	}
	if err != nil && !upfIwkEpsInd {
		matched = true
		log.Debugf("%s is not exist in nfprofile, default is false, param %s is %v", constvalue.IwkEpsInd, constvalue.SearchDataUpfIwkEpsInd, upfIwkEpsInd)
	}

	return matched
}

func (a *NFUPFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFUPFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	var snssaiTotalExpressionList []*MetaExpression
	var snssaiUpfInfoExpressionList []*MetaExpression

	smfServingArea := queryForm.GetNRFDiscSmfServingArea()
	if "" != smfServingArea {
		smfServingAreaPath := getParamSearchPath(constvalue.NfTypeUPF, constvalue.SearchDataSmfServingArea)
		smfServingAreaExpression := buildStringSearchParameter(smfServingAreaPath, smfServingArea, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, smfServingAreaExpression)
	}

	if isSnssaisParaOnly(queryForm) {
		//snssais search in nfProfile
		snssaiExpression := createSnssaiFilter(queryForm, constvalue.Common, constvalue.SearchDataSnssais)
		if snssaiExpression != nil {
			snssaiTotalExpressionList = append(snssaiTotalExpressionList, snssaiExpression)
		}
		//snssais search in upfInfo
		snssaiUpfInfoExpression := createSnssaiFilter(queryForm, constvalue.NfTypeUPF, constvalue.SearchDataSnssais)
		absenceUpfInfoSnssaiExpression := createAbsenceInfoSnssaiExpression(constvalue.NfTypeUPF, constvalue.SearchDataSnssais)
		var snssaiExpressionList []*MetaExpression
		snssaiExpressionList = append(snssaiExpressionList, snssaiUpfInfoExpression, absenceUpfInfoSnssaiExpression)
		snssaiTotalExpression := buildORExpression(snssaiExpressionList)
		if snssaiUpfInfoExpression != nil {
			snssaiTotalExpressionList = append(snssaiTotalExpressionList, snssaiTotalExpression)
		}
		metaExpressionList = append(metaExpressionList, buildORExpression(snssaiTotalExpressionList))
	} else {
		dnn := queryForm.GetNRFDiscDnnValue()
		if "" != dnn {
			dnnPath := getParamSearchPath(constvalue.NfTypeUPF, constvalue.SearchDataDnn)
			dnnExpression := buildStringSearchParameter(dnnPath, dnn, constvalue.EQ)
			snssaiUpfInfoExpressionList = append(snssaiUpfInfoExpressionList, dnnExpression)
		}

		dnaiList := queryForm.GetNRFDiscDnaiList()
		if len(dnaiList) > 0 {
			dnaiListPath := getParamSearchPath(constvalue.NfTypeUPF, constvalue.SearchDataDnaiList)
			dnaiListExpression := createDnaiListFilter(dnaiList, dnaiListPath)
			absenceDnaiListExpression := buildStringSearchParameter(dnaiListPath, constvalue.EmptyDnai, constvalue.EQ)
			var dnaiExpressionList []*MetaExpression
			dnaiExpressionList = append(dnaiExpressionList, dnaiListExpression, absenceDnaiListExpression)
			dnaiTotalExpression := buildORExpression(dnaiExpressionList)
			if dnaiListExpression != nil {
				snssaiUpfInfoExpressionList = append(snssaiUpfInfoExpressionList, dnaiTotalExpression)
			}
		}

		sNssais := queryForm.GetNRFDiscListSnssais(constvalue.SearchDataSnssais)
		if sNssais != "" {
			snssaiExpression := createSnssaiFilter(queryForm, constvalue.NfTypeUPF, constvalue.SearchDataSnssais)
			absenceUpfInfoSnssaiExpression := createAbsenceInfoSnssaiExpression(constvalue.NfTypeUPF, constvalue.SearchDataSnssais)
			var snssaiExpressionList []*MetaExpression
			snssaiExpressionList = append(snssaiExpressionList, snssaiExpression, absenceUpfInfoSnssaiExpression)
			snssaiTotalExpression := buildORExpression(snssaiExpressionList)
			if snssaiTotalExpression != nil {
				snssaiUpfInfoExpressionList = append(snssaiUpfInfoExpressionList, snssaiTotalExpression)
			}
		}

		if buildAndExpression(snssaiUpfInfoExpressionList) != nil {
			metaExpressionList = append(metaExpressionList, buildAndExpression(snssaiUpfInfoExpressionList))
		}
	}
	return metaExpressionList
}
