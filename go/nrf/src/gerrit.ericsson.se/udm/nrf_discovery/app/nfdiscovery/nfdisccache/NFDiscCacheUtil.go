package nfdisccache

import (
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"com/dbproxy"
	"com/dbproxy/nfmessage/nrfprofile"
)

//CacheNFResponse is struct of result from db
type CacheNFResponse struct {
	NfInstanceID      string
	ProfileUpdateTime int
}

//NfProfileCache is nfprofile cache
var NfProfileCache *Cache

//InitNfProfileCache init nfprofile cache
func InitNfProfileCache() {
	NfProfileCache = newCache()
	NfProfileCache.startCache()
}

// GetNFProfileFromCache get nfprofile from cache by instId list
func GetNFProfileFromCache(keys []CacheNFResponse) ([]CacheItem, []string) {
	var items []CacheItem
	var foundKeys []string
	var notFoundKeys []string
	for _, value := range keys {
		item, exist := NfProfileCache.get(value.NfInstanceID, value.ProfileUpdateTime)
		if exist {
			foundKeys = append(foundKeys, value.NfInstanceID)
			items = append(items, item)
		} else {
			notFoundKeys = append(notFoundKeys, value.NfInstanceID)
		}
	}
	if len(foundKeys) > 0 {
		log.Debugf("Find keys in cache, {%v}", foundKeys)
	}
	return items, notFoundKeys
}

//getNFProfileByInstID get nfprofile by instId from db
func getNFProfileByInstID(key string) (string) {
	nfProfileGetRequest := &dbproxy.QueryRequest{
		RegionName: configmap.DBNfprofileRegionName,
		Query: []string{key},
	}
	nfProfileResponse, err := dbmgmt.QueryWithKey(nfProfileGetRequest)
	if err != nil {
		log.Debugf("Discover NF profile failed. DB error, %v", err)
		return ""
	}
	if nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		log.Debugf("requested NF profile not found from DB")
		return ""
	}
	if len(nfProfileResponse.GetValue()) <= 0 {
		return ""
	}
	nfProfile := nfProfileResponse.GetValue()[0]
	if err == nil {
		return nfProfile
	}
	return ""
}

//SplitNFProfileList is to split nfprofile with bodycommon, nfservices, nfinfo, md5sum
func SplitNFProfileList(nfProfiles []string) []CacheItem {
	var itemList []CacheItem
	for _, item := range (nfProfiles) {
		nfType, err := jsonparser.GetString([]byte(item), constvalue.BODY, constvalue.NfType)
		if err != nil {
			log.Error("nfType not exist in profile")
			continue
		}
		infoName := constvalue.TargetNFInfo[nfType]
		var bodyCommon []byte
		body, err0 := dbmgmt.GetBody([]byte(item))
		bodyStr := string(body)
		profileUpdateTime, err1 := jsonparser.GetInt([]byte(item), constvalue.ProfileUpdateTime)
		if err0 != nil || err1 != nil {
			log.Error("No body or profileUpdateTime in custom nf profile")
			continue
		}
		md5sum, err2 := dbmgmt.GetMd5Sum([]byte(item))
		if err2 != nil {
			log.Errorf("No md5sum in custom nf profile")
		}
		instanceID, err3 := jsonparser.GetString(body, constvalue.NfInstanceId)
		if err3 != nil {
			log.Error("No InstanceId in body")
			continue
		}
		nfServices, _, _, err4 := jsonparser.Get(body, constvalue.NfServices)
		if err4 != nil {
			log.Debugf("nfservices not exist in nfProfile, err=%v", err4)
		}

		bodyCommon = jsonparser.Delete([]byte(bodyStr), constvalue.NfServices)
		bodyCommonStr := string(bodyCommon)
		var nfInfo []byte
		var exist jsonparser.ValueType
		if infoName != "" {
			nfInfo, exist, _, err = jsonparser.Get([]byte(bodyCommonStr), infoName)
			if exist != jsonparser.NotExist && err == nil {
				bodyCommon = jsonparser.Delete([]byte(bodyCommonStr), infoName)
			}
		}
		itemList = append(itemList, CacheItem{Key:instanceID, Value:bodyStr, BodyCommon:string(bodyCommon), NfInfo:string(nfInfo), NfServices:string(nfServices), ProfileUpdateTime:int(profileUpdateTime), MD5Sum:string(md5sum)})
	}
	return itemList
}

//SplitNRFProfileList is to split nrfprofile
func SplitNRFProfileList(nrfProfiles []*nrfprofile.NRFProfileInfo) []CacheItem {
	var itemList []CacheItem
	for _, item := range (nrfProfiles) {
		body := item.RawNrfProfile
		nfType, err := jsonparser.GetString(body, constvalue.NfType)
		if err != nil {
			log.Errorf("nfType not exist in profile, error=%v", err)
			continue
		}
		infoName := constvalue.TargetNFInfo[nfType]
		var bodyCommon []byte
		bodyStr := string(body)
		expiredTime := item.ExpiredTime
		instanceID, err3 := jsonparser.GetString(body, constvalue.NfInstanceId)
		if err3 != nil {
			log.Error("No InstanceId in body")
			continue
		}
		nfServices, _, _, err4 := jsonparser.Get(body, constvalue.NfServices)
		if err4 != nil {
			log.Debugf("nfservices not exist in nfProfile, err=%v", err4)
		}

		bodyCommon = jsonparser.Delete([]byte(bodyStr), constvalue.NfServices)
		bodyCommonStr := string(bodyCommon)
		var nfInfo []byte
		var exist jsonparser.ValueType
		if infoName != "" {
			nfInfo, exist, _, err = jsonparser.Get([]byte(bodyCommonStr), infoName)
			if exist != jsonparser.NotExist && err == nil {
				bodyCommon = jsonparser.Delete([]byte(bodyCommonStr), infoName)
			}
		}
		itemList = append(itemList, CacheItem{Key:instanceID, Value:bodyStr, BodyCommon:string(bodyCommon), NfInfo:string(nfInfo), NfServices:string(nfServices), ExpiredTime:int(expiredTime)})
	}
	return itemList
}