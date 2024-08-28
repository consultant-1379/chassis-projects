package nfdiscfilter

import (
	"fmt"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/buger/jsonparser"
	"sort"
	"crypto/md5"
	"net/http"
	"time"
	"com/dbproxy/nfmessage/nrfprofile"
	"com/dbproxy"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdisccache"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"strings"
	"strconv"
)
//FilterInterface interface to prcess different nftype discovery request
type FilterInterface interface {
	filter(nfInfo []byte, queryFrom *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool
}

//FilterInfo to store filter result in different phase fitler
type FilterInfo struct {
	nrfProfileList           []*nrfprofile.NRFProfileInfo
	customNFProfileList      []nfdisccache.CacheItem
	customNRFProfileList     []nfdisccache.CacheItem
	logcontent               *log.LogStruct
	statusCode               int
	problemDetails           *problemdetails.ProblemDetails
	KVDBSearch               bool
	originProfile            nfdisccache.CacheItem
	newProfiles              string
	backupNewProfiles        string
	groupID                  []string
	isInstanceIDSearch       bool

	etagKeys                 []string
	backuupEtagKeys          []string
	nfProfilesMd5Sum         map[string]string
	backupNFProfileMd5Sum    map[string]string
	etagExist                bool
	etagStr                  string

	plmnForbiddenInProfile   bool
	plmnForbiddenInService   bool
	nfTypeForbiddenInProfile bool
	nfTypeForbiddenInService bool
	domainForbiddenInService bool
	domainForbiddenInProfile bool

	errorInfo                string
}

//NFDiscFilter as interface to process discovery request
type NFDiscFilter struct {
	nrfProfileGetRequest   *nrfprofile.NRFProfileGetRequest
	nfProfileGetRequest    *dbproxy.QueryRequest
	queryForm              *nfdiscrequest.DiscGetPara
	filterInfo             *FilterInfo

	DiscNFPreFilterAction  *NFPreFilter
	DiscNFCommonFilter     *NFCommonFilter
	DiscNFInfoFilter       FilterInterface
	DiscNFServiceFiler     *NFServiceFilter
	DiscNFPostFilterAction *NFPostFilter
}

//GetFilterInfoErrorInfo to get filter result's error info
func (f *NFDiscFilter) GetFilterInfoErrorInfo() string {
	return f.filterInfo.errorInfo
}

//GetFilterInfoEtag to get nfprofiles's etag
func (f *NFDiscFilter) GetFilterInfoEtag() string {
	return f.filterInfo.etagStr
}

//GetFilterInfoStatusCode to get the statuscode for response
func (f *NFDiscFilter) GetFilterInfoStatusCode() int {
	return f.filterInfo.statusCode
}

//GetFilterInfoLogcontent to get logcontent
func (f *NFDiscFilter) GetFilterInfoLogcontent() *log.LogStruct {
	return f.filterInfo.logcontent
}

//GetFilterInfoProblem to get problemdetails when filter fail
func (f *NFDiscFilter) GetFilterInfoProblem() *problemdetails.ProblemDetails {
	return f.filterInfo.problemDetails
}

//GetFilterInfoProfiles to get profiles for response
func (f *NFDiscFilter) GetFilterInfoProfiles() string {
	if f.filterInfo.newProfiles != "" {
		return "[" + f.filterInfo.newProfiles + "]"
	}

	return "[" + f.filterInfo.backupNewProfiles + "]"
}

//Init  to initial NFDiscFilter to process discovery request
func (f *NFDiscFilter) Init(queryForm *nfdiscrequest.DiscGetPara) {
	f.queryForm = queryForm

	f.filterInfo = &FilterInfo{nfProfilesMd5Sum:make(map[string]string, 1), backupNFProfileMd5Sum:make(map[string]string, 1)}
	f.filterInfo.logcontent = &log.LogStruct{}
	f.filterInfo.problemDetails = &problemdetails.ProblemDetails{}
	f.nfProfileGetRequest = &dbproxy.QueryRequest{}
	f.nrfProfileGetRequest = &nrfprofile.NRFProfileGetRequest{}
}

//Filter to filter nfprofile for discovery request
func (f *NFDiscFilter) Filter() int {
	f.DiscNFPreFilterAction.generatorGRPCRequst(f.nfProfileGetRequest, f.queryForm, f.filterInfo)

	if !f.getNFProfileFromKVDB() {
		return 1
	}

	log.Debugf("Begin to Filter")
	startTime := time.Now().UnixNano()/1000000
	for _, item := range f.filterInfo.customNFProfileList{

		log.Debugf("nfprofile : %s", item.Value)
		//log.Debugf("nfService: %s", string(item.NfServices))
		//log.Debugf("nfInfo: %s", string(item.NfInfo))
		//log.Debugf("bodyCommon: %s", string(item.BodyCommon))
		f.filterInfo.originProfile = item

		if !f.DiscNFCommonFilter.filter([]byte(f.filterInfo.originProfile.BodyCommon), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFInfoFilter.filter([]byte(f.filterInfo.originProfile.NfInfo), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFServiceFiler.filter([]byte(f.filterInfo.originProfile.NfServices), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFPostFilterAction.filter([]byte(f.filterInfo.originProfile.BodyCommon), f.queryForm, f.filterInfo) {
			continue
		}

	}

	if f.filterInfo.newProfiles == "" {
		if f.filterInfo.backupNewProfiles == "" {
			log.Debugf("not found any nfprofile")
			f.generatorErrorInfo(false)
			return 2
		}
	}

	log.Debugf("Found nfprofiles : %s", f.filterInfo.newProfiles)
	f.calcEtagSum()
	if internalconf.EnableTimeStatistics {
		entTime := time.Now().UnixNano() / 1000000
		dbmgmt.DBLatency.InnerFilterChannel <- dbmgmt.Latency{FilterStartTime:startTime, FilterEndTime:entTime}
	}
	return 0
}

//NrfFilter is to get nrfprofile and filter nrfprofile
func (f *NFDiscFilter) NrfFilter() int {
	f.DiscNFPreFilterAction.generatorNRFGRPCRequst(f.nrfProfileGetRequest, f.queryForm, f.filterInfo)

	if !f.getNRFProfileFromKVDB() {
		return 1
	}

	log.Debugf("Begin to Filter")
	for _, item := range f.filterInfo.customNRFProfileList {
		if f.filterInfo.KVDBSearch {
			expiredTime := item.ExpiredTime
			if expiredTime < int(time.Now().Unix() * 1000) {
				continue
			}
		}
		log.Debugf("nrfprofile : %s", item.Value)

		f.filterInfo.originProfile = item

		if !f.DiscNFCommonFilter.filter([]byte(f.filterInfo.originProfile.BodyCommon), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFInfoFilter.filter([]byte(f.filterInfo.originProfile.BodyCommon), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFServiceFiler.filter([]byte(f.filterInfo.originProfile.NfServices), f.queryForm, f.filterInfo) {
			continue
		}

		if !f.DiscNFPostFilterAction.filter([]byte(f.filterInfo.originProfile.BodyCommon), f.queryForm, f.filterInfo) {
			continue
		}

	}

	if f.filterInfo.newProfiles == "" {
		if f.filterInfo.backupNewProfiles == "" {
			log.Debugf("not found any nrfprofile")
			f.generatorErrorInfo(true)
			return 2
		}
	}
	log.Debugf("Found nrfprofiles : %s", f.filterInfo.newProfiles)
	f.calcProfilesEtag()
	return 0
}

func (f *NFDiscFilter) generatorErrorInfo(isNrfProfile bool) {

	if f.filterInfo.plmnForbiddenInProfile {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-plmn in nfprofile"
		} else {
			f.filterInfo.errorInfo += " or requester-plmn in nfprofile"
		}
	}
	if f.filterInfo.nfTypeForbiddenInProfile {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-nf-type in nfprofile"
		} else {
			f.filterInfo.errorInfo += " or requester-nf-type in nfprofile"
		}
	}
	if f.filterInfo.domainForbiddenInProfile {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-nf-instance-fqdn in nfprofile"
		} else {
			f.filterInfo.errorInfo += " or requester-nf-instance-fqdn in nfprofile"
		}
	}
	if f.filterInfo.plmnForbiddenInService {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-plmn in nfservice"
		} else {
			f.filterInfo.errorInfo += " or requester-plmn in nfservice"
		}
	}
	if f.filterInfo.nfTypeForbiddenInService {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-nf-type in nfservice"
		} else {
			f.filterInfo.errorInfo += " or requester-nf-type in nfservice"
		}
	}
	if f.filterInfo.domainForbiddenInService {
		if f.filterInfo.errorInfo == "" {
			f.filterInfo.errorInfo += "not allowed requester-nf-instance-fqdn in nfservice"
		} else {
			f.filterInfo.errorInfo += " or requester-nf-instance-fqdn in nfservice"
		}
	}
	if f.filterInfo.errorInfo == "" {
		if isNrfProfile {
			f.filterInfo.errorInfo = "requested NRF profile not found "
		} else {
			f.filterInfo.errorInfo = "requested NF profile not found"
		}
		f.filterInfo.statusCode = http.StatusNotFound
	} else {
		f.filterInfo.statusCode = http.StatusForbidden
	}

	f.filterInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
	f.filterInfo.logcontent.ResponseDescription = f.filterInfo.errorInfo
	f.filterInfo.problemDetails.Title = f.filterInfo.errorInfo
}

func (f *NFDiscFilter) getNFProfileFromKVDB() bool {
	log.Debugf("Enter NFDiscFilter GetNFProfileFromKVDB")
	//rw http.ResponseWriter, req *http.Request, problemDetails *problemdetails.ProblemDetails, logcontent *log.LogStruct)(string, error, int, string){
	var nfProfileResponse *dbproxy.QueryResponse
	var err error
	if f.filterInfo.isInstanceIDSearch {
		nfProfileResponse, err = dbmgmt.QueryWithKey(f.nfProfileGetRequest)
	} else {
		nfProfileResponse, err = dbmgmt.QueryWithFilter(f.nfProfileGetRequest)
	}
	if err != nil {
		log.Debugf("Discover NF profile failed. DB error, %v", err)
		errorInfo := fmt.Sprintf("Discover NF profile failed. DB error, %v", err)
		f.filterInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
		f.filterInfo.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		f.filterInfo.problemDetails.Title = errorInfo
		//handleDiscoveryFailure(rw, req, logcontent, http.StatusInternalServerError, problemDetails.ToString())
		//return rw, logcontent, http.StatusInternalServerError, problemDetails.ToString()
		f.filterInfo.statusCode = http.StatusInternalServerError
		return false
	}
	if nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		log.Debugf("requested NF profile not found from DB")
		errorInfo := "requested NF profile not found"
		f.filterInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
		f.filterInfo.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		f.filterInfo.problemDetails.Title = errorInfo
		//handleDiscoveryFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())
		//return rw, logcontent, http.StatusNotFound, problemDetails.ToString()
		f.filterInfo.statusCode = http.StatusNotFound
		return false
	}

	if internalconf.DiscCacheEnable {
		log.Debugf("nfprofile=%v", nfProfileResponse.GetValue())
		if f.filterInfo.isInstanceIDSearch {
			f.filterInfo.customNFProfileList = nfdisccache.SplitNFProfileList(nfProfileResponse.GetValue())
			nfdisccache.NfProfileCache.AddDataChannel <- nfProfileResponse.GetValue()
		} else {
			cacheStructList := nfProfileResponse.GetValue()
			if len(cacheStructList) > 0 {
				nfProfiles, notFoundKeys := nfdisccache.GetNFProfileFromCache(parseResultForCache(cacheStructList))
				f.filterInfo.customNFProfileList = nfProfiles
				if len(notFoundKeys) > 0 {
					log.Debugf("Not find these keys in cache,{%v}", notFoundKeys)
					profilesInDB := getNFProfileByInstIDList(notFoundKeys)
					f.filterInfo.customNFProfileList = append(f.filterInfo.customNFProfileList, nfdisccache.SplitNFProfileList(profilesInDB)...)
				}
			}
		}
	} else {
		f.filterInfo.customNFProfileList = nfdisccache.SplitNFProfileList(nfProfileResponse.GetValue())
	}


	return true
}

//getNFProfileByInstID use instanceId list to get nfprofile
func getNFProfileByInstIDList(keys []string) []string {
	nfProfileGetRequest := &dbproxy.QueryRequest{
		RegionName: configmap.DBNfprofileRegionName,
		Query: keys,
	}
	nfProfileResponse, err := dbmgmt.QueryWithKey(nfProfileGetRequest)
	if err != nil {
		log.Debugf("Discover NF profile failed. DB error, %v", err)
		return []string{}
	}
	if nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		log.Debugf("requested NF profile not found from DB")
		return []string{}
	}
	nfdisccache.NfProfileCache.AddDataChannel <- nfProfileResponse.GetValue()
	return nfProfileResponse.GetValue()
}

//parseResultForCache is to parse [struct(nfInstanceId:12345678-9abc-def0-1000-100000000021,profileUpdateTime:1)]
func parseResultForCache(nfResponse []string) []nfdisccache.CacheNFResponse {
	var result []nfdisccache.CacheNFResponse
	for _, value := range nfResponse {
		cacheNFResponse := nfdisccache.CacheNFResponse{}
		index1 := strings.Index(value, "nfInstanceId:")
		index2 := strings.Index(value, ",")
		cacheNFResponse.NfInstanceID = value[index1+13:index2]
		index3 := strings.Index(value, "profileUpdateTime:")
		index4 := strings.Index(value, ")")
		var err error
		cacheNFResponse.ProfileUpdateTime, err = strconv.Atoi(value[index3+18:index4])
		if err != nil {
			log.Debugf("profileUpdateTime parseInt error, profileUpdateTime=%v, err=%v", value[index3+18:index4], err)
		}
		result = append(result, cacheNFResponse)
	}
	return result
}

func (f *NFDiscFilter) calcEtagSum() {
	var md5sumstr string
	if f.filterInfo.newProfiles != "" {
		sort.Strings(f.filterInfo.etagKeys)
		for _, k := range f.filterInfo.etagKeys {
			md5sumstr = md5sumstr + f.filterInfo.nfProfilesMd5Sum[k]
		}
		eTag := md5.Sum([]byte(md5sumstr))
		f.filterInfo.etagStr = fmt.Sprintf("%x", eTag)
	} else {
		sort.Strings(f.filterInfo.backuupEtagKeys)
		for _, k := range f.filterInfo.backuupEtagKeys{
			 md5sumstr = md5sumstr + f.filterInfo.backupNFProfileMd5Sum[k]
		}
		eTag := md5.Sum([]byte(md5sumstr))
		f.filterInfo.etagStr = fmt.Sprintf("%x", eTag)
	}
}

//getNRFProfileFromKVDB is to get nrfprofile from kvdb
func (f *NFDiscFilter) getNRFProfileFromKVDB() bool {
	log.Debugf("Enter NRFDiscFilter GetNRFProfileFromKVDB")
	//rw http.ResponseWriter, req *http.Request, problemDetails *problemdetails.ProblemDetails, logcontent *log.LogStruct)(string, error, int, string){
	nrfProfileResponse, err := dbmgmt.GetNRFProfile(f.nrfProfileGetRequest)
	if err != nil {
		log.Debugf("Discover NRF profile failed. DB error, %v", err)
		errorInfo := fmt.Sprintf("Discover NRF profile failed. DB error, %v", err)
		f.filterInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
		f.filterInfo.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		f.filterInfo.problemDetails.Title = errorInfo
		f.filterInfo.statusCode = http.StatusInternalServerError
		return false
	}
	if nrfProfileResponse.Code != dbmgmt.DbGetSuccess {
		log.Debugf("requested NRF profile not found from DB")
		errorInfo := "requested NRF profile not found"
		f.filterInfo.logcontent.RequestDescription = fmt.Sprintf(`{"target-nf-type":"%s", "requester-nf-type":"%s"}`, f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), f.queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType))
		f.filterInfo.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		f.filterInfo.problemDetails.Title = errorInfo
		f.filterInfo.statusCode = http.StatusNotFound
		return false
	}

	f.filterInfo.nrfProfileList = nrfProfileResponse.GetNrfProfile()
	FramTotalNumber := nrfProfileResponse.GetFragmentNrfprofileInfo().GetTotalNumber()
	FramTranBumber := nrfProfileResponse.GetFragmentNrfprofileInfo().GetTransmittedNumber()
	FramSessionID := nrfProfileResponse.GetFragmentNrfprofileInfo().GetFragmentSessionId()
	for FramTranBumber != FramTotalNumber {
		var nrfProfileGetRequest *nrfprofile.NRFProfileGetRequest
		nrfProfileFilterData := &nrfprofile.NRFProfileGetRequest_FragmentSessionId{
			FragmentSessionId: FramSessionID,
		}
		nrfProfileGetRequest = &nrfprofile.NRFProfileGetRequest{
			Data: nrfProfileFilterData,
		}

		response, err := dbmgmt.GetNRFProfile(nrfProfileGetRequest)
		if err != nil || response.Code != dbmgmt.DbGetSuccess {
			break
		} else {
			nrfProfilesInfo := response.GetNrfProfile()
			for _, item := range nrfProfilesInfo {
				f.filterInfo.nrfProfileList = append(f.filterInfo.nrfProfileList, item)
			}
			FramTranBumber = response.GetFragmentNrfprofileInfo().GetTransmittedNumber()
			FramSessionID = response.GetFragmentNrfprofileInfo().GetFragmentSessionId()
		}
	}
	f.filterInfo.customNRFProfileList = nfdisccache.SplitNRFProfileList(f.filterInfo.nrfProfileList)
	return true
}

//calcProfilesEtag is to calculate nrfprofile etag value
func (f *NFDiscFilter) calcProfilesEtag() {
	sortedProfileJSON := make(map[string]string, 1)
	_, err := jsonparser.ArrayEach([]byte(f.GetFilterInfoProfiles()), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		profileInstID, err1 := jsonparser.GetString(value, constvalue.NfInstanceId)
		if err1 != nil {
			return
		}
		sortedServiceJSON := make(map[string][]byte, 1)
		_, err2 := jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			serviceInstID, err3 := jsonparser.GetString(value1, constvalue.NFServiceInstanceId)
			if err3 != nil {
				return
			}
			sortedServiceJSON[serviceInstID] = value1
		}, constvalue.NfServices)
		if err2 == nil {
			profileBody := jsonparser.Delete(value, constvalue.NfServices)
			var keys []string
			for k := range sortedServiceJSON {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			tempNFProfile := string(profileBody)
			for _, k := range keys {
				tempNFProfile = tempNFProfile + string(sortedServiceJSON[k])
			}
			sortedProfileJSON[profileInstID] = tempNFProfile
		}
	})

	if err == nil {
		var keys []string
		for k := range sortedProfileJSON {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var temp string
		for _, k := range keys {
			temp = temp + sortedProfileJSON[k]
		}
		md5value := md5.Sum([]byte(temp))
		f.filterInfo.etagStr = fmt.Sprintf("%x", md5value)
	} else {
		f.filterInfo.etagStr = ""
		log.Errorf("%v", err)
	}
}