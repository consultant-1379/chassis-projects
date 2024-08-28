package nfdiscservice

import (
	"fmt"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscfilter"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"net/http"
)

//NFDiscLocalSearchService to discover nfprofile in local
type NFDiscLocalSearchService struct {
}

//Execute to execute discover nfprofile in local
func (l *NFDiscLocalSearchService) Execute(httpInfo *HTTPInfo) {
	if httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) == constvalue.NfTypeNRF {
		l.regionNRFDiscNRFHandler(httpInfo)
		httpInfo.acrossRegionSearch = false
	} else {
		l.regionNRFHandler(httpInfo)
		if httpInfo.statusCode == http.StatusNotFound {
			getCacheProfile(httpInfo)
		}

	}
}

func (l *NFDiscLocalSearchService) regionNRFHandler(httpInfo *HTTPInfo) {
	filter := &nfdiscfilter.NFDiscFilter{DiscNFPreFilterAction: &nfdiscfilter.NFPreFilter{},
		DiscNFCommonFilter:     &nfdiscfilter.NFCommonFilter{},
		DiscNFServiceFiler:     &nfdiscfilter.NFServiceFilter{},
		DiscNFPostFilterAction: &nfdiscfilter.NFPostFilter{}}
	switch httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
	case "AMF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFAMFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFAMFInfoFilter{}
		}
	case "PCF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFPCFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFPCFInfoFilter{}
		}
	case "SMF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFSMFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFSMFInfoFilter{}
		}
	case "UDM":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFUDMInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFUDMInfoFilter{}
		}
	case "UDR":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFUDRInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFUDRInfoFilter{}
		}
	case "AUSF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFAUSFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFAUSFInfoFilter{}
		}
	case "UPF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFUPFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFUPFInfoFilter{}
		}
	case "BSF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFBSFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFBSFInfoFilter{}
		}
	case "CHF":
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFCHFInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFCHFInfoFilter{}
		}
	default:
		{
			filter.DiscNFPreFilterAction.DiscNFInfoFilter = &nfdiscfilter.NFOtherInfoFilter{}
			filter.DiscNFInfoFilter = &nfdiscfilter.NFOtherInfoFilter{}
		}
	}

	filter.Init(&httpInfo.queryForm)
	ret := filter.Filter()
	if ret == 1 {
		log.Warningf("Filter NFProfile in KVDB fail: %s", filter.GetFilterInfoProblem().ToString())
		httpInfo.statusCode = filter.GetFilterInfoStatusCode()
		httpInfo.logcontent = filter.GetFilterInfoLogcontent()
		httpInfo.problem = filter.GetFilterInfoProblem()
		httpInfo.body = filter.GetFilterInfoProblem().ToString()
		return
	} else if ret == 2 {
		log.Warningf("Filter NFProfile fail: %s", filter.GetFilterInfoErrorInfo())
		httpInfo.statusCode = filter.GetFilterInfoStatusCode()
		httpInfo.problem = filter.GetFilterInfoProblem()
		httpInfo.logcontent = filter.GetFilterInfoLogcontent()
		httpInfo.body = filter.GetFilterInfoProblem().ToString()
		return
	}

	etag := filter.GetFilterInfoEtag()
	httpInfo.body = filter.GetFilterInfoProfiles()
	httpInfo.statusCode = filter.GetFilterInfoStatusCode()

	httpInfo.body, httpInfo.statusCode = l.setEtagOnResponse(httpInfo.queryForm, httpInfo.body, etag, httpInfo.req)
	httpInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
	httpInfo.logcontent.ResponseDescription = fmt.Sprintf(`"successful"`)

	if httpInfo.statusCode == http.StatusNotFound {
		httpInfo.body = "NFProfile Not Found"
		httpInfo.problem.Title = "NFProfile Not Found"
		return
	}
	httpInfo.rw.Header().Set("ETag", etag)
	cacheControl := fmt.Sprintf("public, max-age=%d", cm.ValidityPeriodOfSearchResult)
	httpInfo.rw.Header().Set("Cache-Control", cacheControl)
	//handleDiscoverySuccess(rw, req, logcontent, statusCode, newProfiles)
	return

}

func (l *NFDiscLocalSearchService) setEtagOnResponse(queryForm nfdiscrequest.DiscGetPara, newProfiles string, eTag string, req *http.Request) (string, int) {
	validityPeriod := cm.ValidityPeriodOfSearchResult
	statusCode := http.StatusOK
	if eTag != "" {
		if nil != nfdiscrequest.GetNRFDiscIfNoneMatch(req) {
			eTagList := nfdiscrequest.GetNRFDiscIfNoneMatch(req)
			if eTagList[0] == "*" {
				statusCode = http.StatusNotModified
			} else {
				for _, i := range eTagList {
					log.Debugf("eTag List: %s", i)
					if i == eTag {
						statusCode = http.StatusNotModified
					}
				}

			}
		}

		profiles := fmt.Sprintf(constvalue.SearchResult, validityPeriod, newProfiles)
		if statusCode == http.StatusNotModified {
			profiles = ""
		}
		return profiles, statusCode
	}
	return "", http.StatusNotFound
}

func (l *NFDiscLocalSearchService) supportParaForNRFProfile(queryForm nfdiscrequest.DiscGetPara) bool {
	ret := true
	for key := range queryForm.GetValue() {
		if constvalue.NRFParaMap[key] != true {
			return false
		}
	}
	return ret
}

//regionNRFDiscNRFHandler is to handle region nrf disc nrfprofile
func (l *NFDiscLocalSearchService) regionNRFDiscNRFHandler(httpInfo *HTTPInfo) {
	if !l.supportParaForNRFProfile(httpInfo.queryForm) {
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.problem.Title = "NRFProfile Not Found"
		httpInfo.body = httpInfo.problem.ToString()
		return
	}
	filter := &nfdiscfilter.NFDiscFilter{DiscNFPreFilterAction: &nfdiscfilter.NFPreFilter{},
		DiscNFCommonFilter:     &nfdiscfilter.NFCommonFilter{},
		DiscNFInfoFilter:       &nfdiscfilter.NFNRFInfoFilter{},
		DiscNFServiceFiler:     &nfdiscfilter.NFServiceFilter{},
		DiscNFPostFilterAction: &nfdiscfilter.NFPostFilter{}}
	filter.Init(&httpInfo.queryForm)
	ret := filter.NrfFilter()
	if ret == 1 {
		log.Warningf("Filter NRFProfile in KVDB fail: %s", filter.GetFilterInfoProblem().ToString())
		httpInfo.statusCode = filter.GetFilterInfoStatusCode()
		httpInfo.logcontent = filter.GetFilterInfoLogcontent()
		httpInfo.problem = filter.GetFilterInfoProblem()
		httpInfo.body = filter.GetFilterInfoProblem().ToString()
		return
	} else if ret == 2 {
		log.Warningf("Filter NRFProfile fail: %s", filter.GetFilterInfoErrorInfo())
		httpInfo.statusCode = filter.GetFilterInfoStatusCode()
		httpInfo.logcontent = filter.GetFilterInfoLogcontent()
		httpInfo.problem = filter.GetFilterInfoProblem()
		httpInfo.body = filter.GetFilterInfoProblem().ToString()
		return
	}

	etag := filter.GetFilterInfoEtag()
	httpInfo.body = filter.GetFilterInfoProfiles()
	//statusCode := filter.GetFilterInfoStatusCode()

	httpInfo.body, httpInfo.statusCode = l.setEtagOnResponse(httpInfo.queryForm, httpInfo.body, etag, httpInfo.req)
	httpInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), httpInfo.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
	httpInfo.logcontent.ResponseDescription = fmt.Sprintf(`"successful"`)

	if httpInfo.statusCode == http.StatusNotFound {

		httpInfo.statusCode = http.StatusNotFound
		httpInfo.problem.Title = "NRFProfile Not Found"
		httpInfo.body = httpInfo.problem.ToString()
		return
	}
	httpInfo.rw.Header().Set("ETag", etag)
	cacheControl := fmt.Sprintf("public, max-age=%d", cm.ValidityPeriodOfSearchResult)
	httpInfo.rw.Header().Set("Cache-Control", cacheControl)

	return
}
