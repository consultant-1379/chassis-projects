package nfdiscfilter

import (
	"com/dbproxy/nfmessage/nrfprofile"
	"math"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"com/dbproxy"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
)

//PreFilterInterface as interface for different nftype to genenator search parameter
type PreFilterInterface interface {
	filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression
}

//NFPreFilter to use kvdb filter nfprofile first
type NFPreFilter struct {
	DiscNFInfoFilter PreFilterInterface
}

//generatorGRPCRequst is to generate nfprofile request
func (p *NFPreFilter) generatorGRPCRequst(nfProfileGetRequest *dbproxy.QueryRequest, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) {
	log.Debugf("Enter NFPreFilter generatorGRPCRequst")

	groupIDList := queryForm.GetNRFDiscGroupIDList()
	if groupIDList != nil {
		filterInfo.groupID = append(filterInfo.groupID, groupIDList...)
	}

	if queryForm.GetNRFDiscNFInstIDValue() != "" {

		if queryForm.GetNRFDiscSupiValue() != "" {
			groupID, _ := nfdiscutil.GetGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscSupiValue())
			filterInfo.groupID = append(filterInfo.groupID, groupID...)
		}

		if queryForm.GetNRFDiscGspi() != "" {
			groupID, _ := nfdiscutil.GetGpsiGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscGspi())
			filterInfo.groupID = append(filterInfo.groupID, groupID...)
		}
		p.filterByInstanceID(nfProfileGetRequest, queryForm)
		filterInfo.isInstanceIDSearch = true
		return
	} else if queryForm.GetNRFDiscTargetNFFQDN() != "" {
		if queryForm.GetNRFDiscSupiValue() != "" {
			groupID, _ := nfdiscutil.GetGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscSupiValue())
			filterInfo.groupID = append(filterInfo.groupID, groupID...)
		}

		if queryForm.GetNRFDiscGspi() != "" {
			groupID, _ := nfdiscutil.GetGpsiGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscGspi())
			filterInfo.groupID = append(filterInfo.groupID, groupID...)
		}

		p.filterByTargetNFFQDN(nfProfileGetRequest, queryForm)
		return
	}
	filterInfo.KVDBSearch = true
	p.filterByNFTypeInfo(nfProfileGetRequest, queryForm)

}

//filterByInstanceID is to filter nfprofile by instance id
func (p *NFPreFilter) filterByInstanceID(nfProfileGetRequest *dbproxy.QueryRequest, queryForm *nfdiscrequest.DiscGetPara) {
	nfProfileGetRequest.RegionName = configmap.DBNfprofileRegionName
	nfProfileGetRequest.Query = []string{queryForm.GetNRFDiscNFInstIDValue()}
}

//filterByTargetNFFQDN is to filter nfprofile by fqdn
func (p *NFPreFilter) filterByTargetNFFQDN(nfProfileGetRequest *dbproxy.QueryRequest, queryForm *nfdiscrequest.DiscGetPara) {
	var metaExpressionList []*MetaExpression
	targetFqdn := queryForm.GetNRFDiscTargetNFFQDN()
	fqdnPath := getParamSearchPath(constvalue.Common, constvalue.SearchDataTargetNFFQDN)
	targetFqdnExpression := buildStringSearchParameter(fqdnPath, targetFqdn, constvalue.EQ)
	metaExpressionList = append(metaExpressionList, targetFqdnExpression)
	andMetaExpression := buildAndExpression(metaExpressionList)

	var metaExpressionString string
	andMetaExpression.metaExpressionToString(&metaExpressionString)
	log.Debugf("build expression={%s}", metaExpressionString)

	var searchOql string
	if internalconf.DiscCacheEnable {
		buildInstIDOql(andMetaExpression, &searchOql)
	} else {
		buildOql(andMetaExpression, &searchOql)
	}
	log.Debugf("filter searchOql={%s}", searchOql)

	nfProfileGetRequest.RegionName = configmap.DBNfprofileRegionName
	nfProfileGetRequest.Query = []string{searchOql}
}

//filterByNFTypeInfo is to filter nfprofile by info
func (p *NFPreFilter) filterByNFTypeInfo(nfProfileGetRequest *dbproxy.QueryRequest, queryForm *nfdiscrequest.DiscGetPara) {
	log.Debugf("Enter NFPreFilter filterByNFTypeInfo")
	var metaExpressionList []*MetaExpression

	//metaExpressionList = append(metaExpressionList, p.DiscNFInfoFilter.filterByTypeInfo()...)
	metaExpressionList = append(metaExpressionList, p.DiscNFInfoFilter.filterByKVDB(queryForm)...)

	nfType := queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)
	if nfType != constvalue.NfTypeUPF && nfType != constvalue.NfTypeSMF {
		snssaiExpression := createSnssaiFilter(queryForm, constvalue.Common, constvalue.SearchDataSnssais)
		if snssaiExpression != nil {
			metaExpressionList = append(metaExpressionList, snssaiExpression)
		}
	}

	targetNfType := queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)
	nftypePath := getParamSearchPath(constvalue.Common, constvalue.SearchDataTargetNfType)
	targetNfTypeExpression := buildStringSearchParameter(nftypePath, targetNfType, constvalue.EQ)
	metaExpressionList = append(metaExpressionList, targetNfTypeExpression)
	andMetaExpression := buildAndExpression(metaExpressionList)

	var metaExpressionString string
	andMetaExpression.metaExpressionToString(&metaExpressionString)
	log.Debugf("build expression={%s}", metaExpressionString)

	var searchOql string
	if internalconf.DiscCacheEnable {
		buildInstIDOql(andMetaExpression, &searchOql)
	} else {
		buildOql(andMetaExpression, &searchOql)
	}
	log.Debugf("filter searchOql={%s}", searchOql)

	nfProfileGetRequest.RegionName = configmap.DBNfprofileRegionName
	nfProfileGetRequest.Query = []string{searchOql}
}

//generatorNRFGRPCRequst is to generate nrfprofile filter request
func (p *NFPreFilter) generatorNRFGRPCRequst(nrfProfileGetRequest *nrfprofile.NRFProfileGetRequest, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) {
	log.Debugf("Enter NRFDiscPreFilter generatorGRPCRequst")
	if queryForm.GetNRFDiscNFInstIDValue() != "" {
		filterInfo.KVDBSearch = true
		p.filterNRFByInstanceID(nrfProfileGetRequest, queryForm)
		return
	}
	filterInfo.KVDBSearch = false
	p.filterNRFByExpiredTime(nrfProfileGetRequest, queryForm)
}

//filterNRFByInstanceID is to filter nrfprofile by instanceId
func (p *NFPreFilter) filterNRFByInstanceID(nrfProfileGetRequest *nrfprofile.NRFProfileGetRequest, queryForm *nfdiscrequest.DiscGetPara) {
	instanceIDRequest := &nrfprofile.NRFProfileGetRequest_NrfInstanceId{
		NrfInstanceId: queryForm.GetNRFDiscNFInstIDValue(),
	}
	nrfProfileGetRequest.Data = instanceIDRequest
}

//filterNRFByExpiredTime is to filter nrfprofile by expiredTime
func (p *NFPreFilter) filterNRFByExpiredTime(nrfProfileGetRequest *nrfprofile.NRFProfileGetRequest, queryForm *nfdiscrequest.DiscGetPara) {
	nrfProfileIndex := &nrfprofile.NRFProfileIndex{
		Key1: uint64(time.Now().Unix()) * 1000,
		Key2: math.MaxInt64,
	}
	nrfProfileFilter := &nrfprofile.NRFProfileFilter{
		AndOperation: true,
		Index:        nrfProfileIndex,
	}
	nrfProfileFilterData := &nrfprofile.NRFProfileGetRequest_Filter{
		Filter: nrfProfileFilter,
	}
	nrfProfileGetRequest.Data = nrfProfileFilterData
}
