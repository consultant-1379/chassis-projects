package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"strings"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/buger/jsonparser"
)
//NFSMFInfoFilter to process smfinfo filter in nfprofile
type NFSMFInfoFilter struct {

}

func (a *NFSMFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	if !a.isMatchedByPGWInd(nfInfo, queryForm) {
		return false
	}
	accessType := queryForm.GetNRFDiscAccessType();
	if accessType != "" {
		if !a.isMatchedAccessType(nfInfo, accessType) {
			return false
		}
	}

	return true
}

func (a *NFSMFInfoFilter) isMatchedByPGWInd(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara) bool {
	pgwind, err := queryForm.GetNRFDiscBoolValue(constvalue.SearchDataPGWInd)
	if err == nil {
		pgwFQDN, err := jsonparser.GetString(nfInfo, constvalue.PgwFqdn)
		if pgwind {
			if err != nil {
				return false //pgwind is true, pgwFqdn should be exist, not match
			}
			pgwFQDN = strings.Replace(pgwFQDN, " ", "", -1)
			if "" == pgwFQDN {
				return false //pgwind is true, pgwFqdn should be exist, and value should not empty, not match
			}
			return true
		}
		if err != nil {
			return true //pgwInd is false, pgwFqdn should not exist, match
		}
		pgwFQDN = strings.Replace(pgwFQDN, " ", "", -1)
		if "" == pgwFQDN {
			return true //pgwind is false, but pgwFqdn value is empty, match
		}
		return false
	}
	return true
}

//isMatchedAccessType is to match access-type in smfInfo
func (a *NFSMFInfoFilter) isMatchedAccessType(nfInfo []byte, accessType string) bool {
	matched := false
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		accessTypeInProfile := string(value[:])
		if accessType == accessTypeInProfile {
			matched = true
			return
		}
	}, constvalue.AccessType)
	if err != nil {
		matched = true
	}
	return matched
}

func (a *NFSMFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFSMFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	var snssaiTotalExpressionList []*MetaExpression
	var snssaiSmfInfoExpressionList []*MetaExpression
	pgw := queryForm.GetNRFDiscPGW()
	if "" != pgw {
		pgwPath := getParamSearchPath(constvalue.NfTypeSMF, constvalue.SearchDataPGW)
		pgwExpression := buildStringSearchParameter(pgwPath, pgw, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, pgwExpression)
	}

	plmnid, tac := queryForm.GetNRFDiscTaiType()
	if len(plmnid) >= 5 && tac != "" {
		tac = strings.ToLower(tac)
		taiNormal := createTaiExpression(constvalue.NfTypeSMF, constvalue.SearchDataTai, plmnid, tac)
		taiAbsence := createTaiExpressionForAbsence(constvalue.NfTypeSMF, constvalue.SearchDataTai)

		var taiExpressionList []*MetaExpression
		taiExpressionList = append(taiExpressionList, taiNormal)
		taiExpressionList = append(taiExpressionList, taiAbsence)

		metaExpressionList = append(metaExpressionList, buildORExpression(taiExpressionList))

	}

	if isSnssaisParaOnly(queryForm) {
		//snssais search in nfProfile
		snssaiExpression := createSnssaiFilter(queryForm, constvalue.Common, constvalue.SearchDataSnssais)
		if snssaiExpression != nil {
			snssaiTotalExpressionList = append(snssaiTotalExpressionList, snssaiExpression)
		}
		//snssais search in smfInfo
		snssaiUpfInfoExpression := createSnssaiFilter(queryForm, constvalue.NfTypeSMF, constvalue.SearchDataSnssais)
		absenceUpfInfoSnssaiExpression := createAbsenceInfoSnssaiExpression(constvalue.NfTypeSMF, constvalue.SearchDataSnssais)
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
			dnnPath := getParamSearchPath(constvalue.NfTypeSMF, constvalue.SearchDataDnn)
			dnnExpression := buildStringSearchParameter(dnnPath, dnn, constvalue.EQ)
			if dnnExpression != nil {
				snssaiSmfInfoExpressionList = append(snssaiSmfInfoExpressionList, dnnExpression)
			}
		}

		sNssais := queryForm.GetNRFDiscListSnssais(constvalue.SearchDataSnssais)
		if sNssais != "" {
			snssaiExpression := createSnssaiFilter(queryForm, constvalue.NfTypeSMF, constvalue.SearchDataSnssais)
			absenceSnssaiExpression := createAbsenceInfoSnssaiExpression(constvalue.NfTypeSMF, constvalue.SearchDataSnssais)
			var snssaiExpressionList []*MetaExpression
			snssaiExpressionList = append(snssaiExpressionList, snssaiExpression, absenceSnssaiExpression)
			snssaiTotalExpression := buildORExpression(snssaiExpressionList)
			if snssaiTotalExpression != nil {
				snssaiSmfInfoExpressionList = append(snssaiSmfInfoExpressionList, snssaiTotalExpression)
			}
		}

		if buildAndExpression(snssaiSmfInfoExpressionList) != nil {
			metaExpressionList = append(metaExpressionList, buildAndExpression(snssaiSmfInfoExpressionList))
		}
	}

	return metaExpressionList
}
