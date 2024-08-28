package nfdiscservice

import (
	"com/dbproxy/nfmessage/cachenfprofile"
	"com/dbproxy"
	"fmt"
	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/kvdbclient"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"github.com/buger/jsonparser"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	//MarkDiscLocalCacheStatus when cm.DiscLocalCacheEnable from true to false, trigger to clear all cache
	MarkDiscLocalCacheStatus = false
	//MarkDiscLocalCacheCapacity to store local cache capacit, when cm change, if diff, alter capacity
	MarkDiscLocalCacheCapacity = 100
	//ClearCacheProfilesCommand is command to clear local cache
	ClearCacheProfilesCommand = "remove --region=/ericsson-nrf-cachenfprofiles --all"
	//QueryCacheProfilesCountCommand is command to check cache items number in cache
	QueryCacheProfilesCountCommand = "query --query=\"select count(*) from /ericsson-nrf-cachenfprofiles\""
	//QueryCacheProrilesMatchResult is to check whether clear cache success
	QueryCacheProrilesMatchResult = "Result\n------\n0\n"
	//RemoteCachePutTime is field for cachenfprofiles's entry, when insert this entry
	RemoteCachePutTime              = "put_time"
	//RemoteCacheExpiryTime is field for cachenfprofiles's entry, when the entry timeout
	RemoteCacheExpiryTime           = "expiry_time"
	//RemoteCacheFromRedirect is flag value to indicate that the cached nfprofiles got by rediect
	RemoteCacheFromRedirect         = "FromRedirect"
	//RemoteCacheFROM is to indicate the entry from which nrf
	RemoteCacheFROM                 = "from"
	//RemoteCacheRawProfile is field for cachenfprofiels to store rawprofile
	RemoteCacheRawProfile           = "rawProfile"
	//PlmnNRFInstanceID to store plmn NRF instanceID, if some instance remove from cm, need delete cache entrys that from this instanceid
	PlmnNRFInstanceID []string
)

//InitMarkDiscLocalCacheCapacity is sync cache capacity when discovery start
func InitMarkDiscLocalCacheCapacity() {
	MarkDiscLocalCacheCapacity = cm.DiscLocalCacheCapacity
	PlmnNRFInstanceID = cm.PlmnNRFInstanceID
	transferParameter(cm.DiscLocalCacheCapacity, true)
}

func transferParameter(capacity int, init bool){
	log.Debugf("local-cache-capacity : %v", capacity)
     go func (capa int, init bool) {
	     for {
		     if init {
			     time.Sleep(time.Duration(5)*time.Second)
		     }
		     paraReq := &dbproxy.ParaRequest{ParameterName:"local-cache-capacity", ParameterValue:strconv.Itoa(capa)}
		     paraResp, err := dbmgmt.TransferParameter(paraReq)
		     if err == nil && paraResp.Code == dbmgmt.DbPutSuccess {
			     break;
		     } else {
			     time.Sleep(time.Duration(1)*time.Second)
		     }
	     }
     }(capacity, init)
}

func responseAllowCache(resp *httpclient.HttpRespData) bool {
	if !cm.DiscLocalCacheEnable || resp.StatusCode != http.StatusOK {
		return false
	}

	valList := nfdiscrequest.GetNRFDiscRespCacheControl(resp)

	for _, v := range valList {
		if v == constvalue.SearchDataCacheControlNoStore || v == constvalue.SearchDataCacheControlNoCache || v == constvalue.SearchDataCacheControlPrivate || v == constvalue.SearchDataCacheControlMaxAge0 {
			return false
		}
	}
	return true
}

func putCacheProfile(req *http.Request, resp *httpclient.HttpRespData, from string) {
	//Only localcache enable and response statuscode is 200, discovery cache nfprofile
	if !responseAllowCache(resp) {
		return
	}
	log.Debugf("Query parameters: %s", req.URL.RawQuery)
	oriqueryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Debugf("ParseQuery parameters fail : %s, parameters: %s", err.Error(), req.URL.RawQuery)
	}
	var queryForm nfdiscrequest.DiscGetPara
	//queryForm.value = oriqueryForm
	queryForm.InitMember(oriqueryForm)
	problem := queryForm.ValidateNRFDiscovery()
	if problem != nil {
		log.Debugf("generator localcache key fail : %s", problem.ToString())
	}
	key := queryForm.GetLocalCacheKey()
	period, err := jsonparser.GetInt(resp.Body, constvalue.SearchResultValidityPeriod)
	if err != nil || period <= 0 {
		return
	}

	expiredTime := int64(cm.DiscLocalCacheTimeout)
	if period < int64(cm.DiscLocalCacheTimeout) {
		expiredTime = period
	}

	var etag string

	for _, v := range (*resp.Header)[constvalue.HTTPHeaderEtag] {
		if etag != "" {
			etag = etag + "," + v

		} else {
			etag = etag + v
		}
	}
	rawProfile := fmt.Sprintf(`{"put_time": %v, "expiry_time": %v, "from": "%v", "%s" :["%s"],"rawProfile": %s}`, (time.Now().Unix())*1000, expiredTime, from, constvalue.HTTPHeaderEtag, etag, string(resp.Body))
	log.Debugf("cachenrfprofile profile: %s key: %s", rawProfile, key)
	cacheNFProfilePutReq := &cachenfprofile.CacheNFProfilePutRequest{
		CacheNfInstanceId: key,
		RawCacheNfProfile: rawProfile,
	}

	cacheNFprofileResp, err := dbmgmt.PutCacheNFProfile(cacheNFProfilePutReq)

	if err != nil || cacheNFprofileResp.Code != dbmgmt.DbPutSuccess {
		log.Warningf("put nfprofile into ericsson-nrf-cachenfprofiles fail request: %s, nfprofile: %s", req.RequestURI, string(resp.Body))
		return
	}

	log.Debugf("put nfprofile into ericsson-nrf-cachenfprofiles success request: %s, nfprofile: %s", req.RequestURI, string(resp.Body))
}

func etagMatched(req *http.Request, body []byte) (bool, []string) {
	etagList := nfdiscrequest.GetNRFDiscIfNoneMatch(req)

	var etagInCache []string

	_, err := jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		etagInCache = append(etagInCache, string(value))
	}, constvalue.HTTPHeaderEtag)

	if etagList == nil {
		return false, etagInCache
	}

	if etagList[0] == "*" {
		return true, etagInCache
	}

	if err != nil {
		return false, etagInCache
	}
	for _, v := range etagList {
		for _, vv := range etagInCache {
			if v == vv {
				return true, etagInCache
			}
		}
	}
	return false, etagInCache
}

func getMaxAgeFromRequest(req *http.Request) (int64, bool) {
	valList := nfdiscrequest.GetNRFDiscReqCacheControl(req)
	for _, v := range valList {
		if v != constvalue.SearchDataCacheControlMaxAge0 && strings.HasPrefix(v, "max-age=") {
			tmpList := strings.Split(v, "=")
			if len(tmpList) == 2 {
				value, err := strconv.ParseInt(tmpList[1], 10, 64)
				if err != nil || value < 0 {
					return 0, false
				}
				return value, true
			}
			return 0, false
		}
	}
	return 0, false
}

func getExpiredTime(resp *cachenfprofile.CacheNFProfileGetResponse, req *http.Request) int {
	putTime, err1 := jsonparser.GetInt([]byte(resp.CacheNfProfile), RemoteCachePutTime)
	expiredTime, err2 := jsonparser.GetInt([]byte(resp.CacheNfProfile), RemoteCacheExpiryTime)
	if err1 != nil || err2 != nil {
		log.Warnf("get puttime expirytime from cachenfprofile fail")
	}
	currentTime := time.Now().Unix()
	log.Debugf("nfprofile put time: %d expiredTime: %d, currentTime: %d", putTime, expiredTime, currentTime)
	if currentTime*1000 >= (putTime + expiredTime*1000) {
		return 0
	}

	maxage, isExist := getMaxAgeFromRequest(req)
	log.Debugf("Request Cache-Control Header max-age: %d", maxage)
	if isExist {
		if (maxage < 0) || (maxage*1000 < (currentTime*1000 - putTime)) {
			return 0
		}
	}

	return int(putTime/1000 + expiredTime - currentTime)
}

func requestAllowCache(req *http.Request) bool {
	if !cm.DiscLocalCacheEnable {
		return false
	}

	valList := nfdiscrequest.GetNRFDiscReqCacheControl(req)
	for _, v := range valList {
		if v == constvalue.SearchDataCacheControlNoCache || v == constvalue.SearchDataCacheControlNoStore || v == constvalue.SearchDataCacheControlMaxAge0 {
			return false
		}
	}

	return true
}

func getCacheProfile(httpInfo *HTTPInfo) {

	if !requestAllowCache(httpInfo.req) {
		httpInfo.body = ""
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.acrossRegionSearch = true
		return //"", false, http.StatusNotFound, rw
	}

	key := httpInfo.queryForm.GetLocalCacheKey()

	cacheNFProfileGetRequest := &cachenfprofile.CacheNFProfileGetRequest{
		CacheNfInstanceId: key,
	}

	cacheNFProfileResponse, err := dbmgmt.GetCacheNFProfile(cacheNFProfileGetRequest)

	if err != nil {
		log.Warningf("get nfprofile from ericcson-nrf-cachenfprofiles fail, key: %s", key)
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.body = ""
		httpInfo.acrossRegionSearch = true
		return //"", false, http.StatusNotFound, rw
	}

	if cacheNFProfileResponse.Code != dbmgmt.DbGetSuccess {
		log.Debugf("not found nfprofile from ericcson-nrf-cachenfprofiles by key : %s", key)
		httpInfo.body = ""
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.acrossRegionSearch = true
		return //"", false, http.StatusNotFound, rw
	}
	timeOut := getExpiredTime(cacheNFProfileResponse, httpInfo.req)
	matched, etagList := etagMatched(httpInfo.req, []byte(cacheNFProfileResponse.CacheNfProfile))
	if matched && timeOut != 0 {
		log.Debugf("get nfprofiles from ericsson-nrf-cachenfprofiles success, 304 not modify")
		for _, v := range etagList {
			httpInfo.rw.Header().Add("Etag", v)
		}
		cacheControl := fmt.Sprintf("public, max-age=%d", timeOut)
		httpInfo.rw.Header().Set("Cache-Control", cacheControl)
		httpInfo.body = ""
		httpInfo.statusCode = http.StatusNotModified
		return //"", true, http.StatusNotModified, rw
	}


	profiles, _, _, err := jsonparser.Get([]byte(cacheNFProfileResponse.CacheNfProfile), RemoteCacheRawProfile, constvalue.SearchResultNFInstances)
	if err != nil {
		log.Debugf("parser cachenfprofile fail error : %s : %s", err.Error(), cacheNFProfileResponse.CacheNfProfile)
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.body = ""
		httpInfo.acrossRegionSearch = true
		return //"", false , http.StatusNotFound, rw
	}

	isRedirect := false
        from,  err := jsonparser.GetString([]byte(cacheNFProfileResponse.CacheNfProfile), RemoteCacheFROM)
	if err == nil && from == RemoteCacheFromRedirect {
		isRedirect = true
	}

        if timeOut == 0 && isRedirect {
		log.Debugf("Get timeout cache nfprofiles that from redirect, no need return timeout cache items")
		httpInfo.statusCode = http.StatusNotFound
		httpInfo.body = ""
		httpInfo.acrossRegionSearch = true;
		return
	}

	bodyInfo := fmt.Sprintf(constvalue.SearchResult, timeOut, string(profiles))
	log.Debugf("get nrfprofiles from ericsson-nrf-cachenfprofiels success, nfprofile :%s, timeout: %v", bodyInfo, timeOut)
        if timeOut != 0 {
		for _, v := range etagList {
			httpInfo.rw.Header().Add("Etag", v)
		}
		cacheControl := fmt.Sprintf("public, max-age=%d", timeOut)
		httpInfo.rw.Header().Set("Cache-Control", cacheControl)
	} else {
		for _, v := range etagList {
			httpInfo.header.Add("Etag", v)
		}
		cacheControl := fmt.Sprintf("public, max-age=%d", timeOut)
		httpInfo.header.Set("Cache-Control", cacheControl)
	}
	httpInfo.body = bodyInfo
	httpInfo.statusCode = http.StatusOK
	if timeOut == 0 {
		httpInfo.cacheTimeout = true
		httpInfo.acrossRegionSearch = true
	}
	return //bodyInfo, true, http.StatusOK, rw
}

//LocalCacheCMUpdate is for when cache configurations is changed, to clear cache or alter cache capacity
func LocalCacheCMUpdate() bool {

	if cm.DiscLocalCacheEnable {
		MarkDiscLocalCacheStatus = cm.DiscLocalCacheEnable
	} else {
		if cm.DiscLocalCacheEnable != MarkDiscLocalCacheStatus {
			clearCacheNFProfile()
			MarkDiscLocalCacheStatus = cm.DiscLocalCacheEnable

		}
	}

	if cm.DiscLocalCacheCapacity != MarkDiscLocalCacheCapacity {
		transferParameter(cm.DiscLocalCacheCapacity, false)
		MarkDiscLocalCacheCapacity = cm.DiscLocalCacheCapacity
	}
        var removedPlmn []string
	for _, v := range PlmnNRFInstanceID{
		matched := false
		for _, vv := range cm.PlmnNRFInstanceID{
			if v == vv{
				matched = true
				break
			}
		}

		if !matched {
			removedPlmn = append(removedPlmn, v)
		}
	}
	PlmnNRFInstanceID = PlmnNRFInstanceID[0:0]
	PlmnNRFInstanceID = cm.PlmnNRFInstanceID

	if len(removedPlmn) > 0 {
		err := dbmgmt.Remove("ericsson-nrf-cachenfprofiles", removedPlmn)
		if err != nil{
			log.Warnf("Delete cache entrys from %v fail", removedPlmn)
		}

	}
	return true
}

//when cm update, localcacheenable  from enable to disable, clear all cache nfprofiles
func clearCacheNFProfile() {
	for retry := 1; retry <= 3; retry++ {
		_, err := kvdbclient.GetInstance().SendGFSHCommand(ClearCacheProfilesCommand)
		if err == nil {
			commandID, err := kvdbclient.GetInstance().SendGFSHCommand(QueryCacheProfilesCountCommand)
			result, err2 := kvdbclient.GetInstance().GetGFSHCommandResult(commandID)
			if err == nil && err2 == nil {
				log.Debugf("Local Cache info: %d, %s", result.StatusCode, result.Output)
				if strings.Contains(result.Output, QueryCacheProrilesMatchResult) {
					log.Debugf("Clear Cache Success")
					return
				}
			}
		}
		sleep := rand.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

	log.Warningf("When local-cache-enable become from true to false, clear Cache fail")
}

//when cm update localcachecapcity, need disable local cache
/*func alterLocalCacheCapacity() bool {
	for retry := 1; retry <= 3; retry++ {
		commandID, err := kvdbclient.GetInstance().SendGFSHCommand(DescribeCacheRegionInfoCommand)
		if err == nil {
			sleep := rand.Intn(100)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			result, err := kvdbclient.GetInstance().GetGFSHCommandResult(commandID)
			if err == nil {
				if result.ExecutionStatus == "EXECUTED" {
					log.Debugf("Local Cache Info : %d %s %s %s", result.StatusCode, result.ExecutionStatus, result.Command, result.Output)
					subStr := fmt.Sprintf(CacheRegionMaxEntryMatcheResult, cm.DiscLocalCacheCapacity)
					if strings.Contains(result.Output, subStr) {
						return true
					}
					command := fmt.Sprintf(AlterCacheRegionCapacityCommand, cm.DiscLocalCacheCapacity)
					_, err = kvdbclient.GetInstance().SendGFSHCommand(command)
					if err != nil {
						log.Warningf("Alter local cache capacity fail")
					}
				}
			}
		}
		sleep := rand.Intn(1000)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
	log.Warningf("modify local cache capacity fail")
	return false
}*/
