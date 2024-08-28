package discutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/subscribe"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/worker"
	"github.com/buger/jsonparser"
)

//MatchResult is for match function result value
type MatchResult int32

const (
	//ResultError is for MatchResult error
	ResultError MatchResult = 0
	//ResultFoundMatch is for MatchResult found and match
	ResultFoundMatch MatchResult = 1
	//ResultFoundNotMatch is for MatchResult found and not match
	ResultFoundNotMatch MatchResult = 2
)

const (
	delayForMaster = 3 * time.Second
)

//getNRFDiscReqCacheControl to get  request header cache-control value
func getNRFDiscReqCacheControl(req *http.Request) []string {
	value := req.Header.Get(consts.SearchDataCacheControl)
	log.Debugf("Request Cache-Control Header Field-value: %s", value)
	valList := strings.Split(value, ",")
	var retList []string
	for _, v := range valList {
		vv := strings.Replace(v, " ", "", -1)
		retList = append(retList, vv)
	}
	return retList
}

//getNRFDiscRespCacheControl to get response header cache-control value
func getNRFDiscRespCacheControl(resp *httpclient.HttpRespData) []string {
	value := resp.Header.Get(consts.SearchDataCacheControl)
	log.Debugf("Response Cache-Control Header Field-value: %s", value)
	valList := strings.Split(value, ",")
	var retList []string
	for _, v := range valList {
		vv := strings.Replace(v, " ", "", -1)
		retList = append(retList, vv)
	}
	return retList
}

//RequestAllowCache query allow nrf-agent search in cache
func RequestAllowCache(req *http.Request) bool {
	valList := getNRFDiscReqCacheControl(req)
	for _, v := range valList {
		if v == consts.SearchDataCacheControlNoCache || v == consts.SearchDataCacheControlNoStore || v == consts.SearchDataCacheControlMaxAge0 {
			return false
		}
	}

	return true
}

//ResponseAllowCache response allow been cached
func ResponseAllowCache(resp *httpclient.HttpRespData) bool {
	if resp.StatusCode != http.StatusOK {
		return false
	}

	valList := getNRFDiscRespCacheControl(resp)
	for _, v := range valList {
		if v == consts.SearchDataCacheControlNoCache || v == consts.SearchDataCacheControlNoStore || v == consts.SearchDataCacheControlPrivate || v == consts.SearchDataCacheControlMaxAge0 {
			return false
		}
	}
	return true
}

//GetMaxAgeFromRequest get max-age from request
func GetMaxAgeFromRequest(req *http.Request) (int64, bool) {
	valList := getNRFDiscReqCacheControl(req)
	for _, v := range valList {
		if v != consts.SearchDataCacheControlMaxAge0 && strings.HasPrefix(v, "max-age=") {
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

/*
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
*/

//ServiceNameScopeVerify verify serviceName search parameter scope
func ServiceNameScopeVerify(serviceNames []string, targetNf *structs.TargetNf) error {
	if targetNf == nil {
		return fmt.Errorf("TargetNf is nil")
	}
	targetServices := targetNf.TargetServiceNames

	if len(serviceNames) == 0 {
		return nil
	}

	for _, service := range serviceNames {
		found := false
		for _, targetService := range targetServices {
			if service == targetService {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("SearchParameter service-names:%s, do not configured in configmap targetNfProfile", service)
		}
	}

	return nil
}

func RejectVerify(oriqueryForm map[string][]string, plmns []structs.PlmnID) error {
	targetPlmnListValue := oriqueryForm[consts.SearchDataTargetPlmnList]
	requesterPlmnListValue := oriqueryForm[consts.SearchDataRequesterPlmnList]
	//without requester-plmn-list, with target-plmn-list
	if len(targetPlmnListValue) > 0 && len(requesterPlmnListValue) == 0 {
		requesterNfType := oriqueryForm[consts.SearchDataRequesterNfType][0]
		nfPlmns, ok := cache.Instance().GetRequesterPlmns(requesterNfType)
		if !ok {
			return nil
		}
		ok = isSubset(nfPlmns, plmns)
		if !ok {
			return fmt.Errorf("NF:%s config plmns:%v, query target-plmn-list(%v). URL need add requester-plmn-list parameter", requesterNfType, nfPlmns, plmns)
		}
	}

	return nil
}

func ForwardVerify(oriqueryForm map[string][]string) bool {
	targetPlmnListValue := oriqueryForm[consts.SearchDataTargetPlmnList]
	requesterPlmnListValue := oriqueryForm[consts.SearchDataRequesterPlmnList]
	if len(targetPlmnListValue) == 0 && len(requesterPlmnListValue) > 0 {
		return true
	}

	requesterNfType := oriqueryForm[consts.SearchDataRequesterNfType][0]
	plmns, ok := cache.Instance().GetRequesterPlmns(requesterNfType)
	isRoam := RoamingVerify(oriqueryForm)
	if (!ok || len(plmns) == 0) && isRoam {
		return true
	}

	return false
}

func RoamingVerify(oriqueryForm map[string][]string) bool {
	targetPlmnListValue := oriqueryForm[consts.SearchDataTargetPlmnList]
	requesterPlmnListValue := oriqueryForm[consts.SearchDataRequesterPlmnList]

	if len(targetPlmnListValue) > 0 && len(requesterPlmnListValue) > 0 {
		return true
	}

	return false
}

func BuildRoamMessage(requesterNfType, targetNfType string, rawResp []byte, roamTargetPlmnID structs.PlmnID) []byte {
	validityPeriod, err := util.GetValidityPeriod(rawResp)
	if err != nil {
		log.Errorf("Get validityPeriod failed, err:%s", err.Error())
		return nil
	}

	nfInstances, err := util.GetNfInstances(rawResp)
	if err != nil {
		log.Errorf("Get nfInstances failed, err:%s", err.Error())
		return nil
	}

	searchResult := structs.SearchResult{
		ValidityPeriod: int32(validityPeriod),
	}
	nfProfiles := make([]structs.SearchResultNFProfile, 0)
	for _, nfProfile := range nfInstances {
		var oneNfProfile structs.SearchResultNFProfile
		err := json.Unmarshal(nfProfile, &oneNfProfile)
		if err != nil {
			log.Warnf("Unmarsh nfProfile fail, will abandon this nfProfile, err:%s", err.Error())
			continue
		}

		ok := cache.Instance().ProbeRoam(requesterNfType, targetNfType, oneNfProfile.NfInstanceID)
		if ok {
			oldProfile := cache.Instance().GetProfileByID(requesterNfType, targetNfType, oneNfProfile.NfInstanceID, true)
			if oldProfile == nil {
				continue
			}
			if appendPlmnInfo(&oneNfProfile, oldProfile, roamTargetPlmnID) {
				nfProfiles = append(nfProfiles, oneNfProfile)
			}
		} else {
			missPlmn := plmnMissProber(nfProfile)
			if missPlmn {
				addPlmnInfo(&oneNfProfile, roamTargetPlmnID)
			}
			nfProfiles = append(nfProfiles, oneNfProfile)
		}

	}

	if len(nfProfiles) == 0 {
		return nil
	}

	searchResult.NfInstances = nfProfiles
	newRespData, err := json.Marshal(searchResult)
	if err != nil {
		log.Warnf("Marsh searchResult after add plmnid fail, err:%s", err.Error())
		return nil
	}

	return newRespData
}

func IsRoamMessage(nfType string, nfProfiles []structs.SearchResultNFProfile) bool {
	nfPlmns, ok := cache.Instance().GetRequesterPlmns(nfType)
	if !ok {
		return false
	}

	for _, nfProfile := range nfProfiles {
		ok = haveInstersection(nfProfile.PLMN, nfPlmns)
		if ok {
			return false
		}
	}

	return true
}

func SubscribeRoamNfProfiles(nfProfiles []byte, requesterNfType, targetNfType string, plmnID *structs.PlmnID, validityPeriod int64) {
	now := time.Now()
	prolongValue := time.Duration(validityPeriod) * time.Second
	timeStamp := now.Add(prolongValue)
	timeStampStr := timeStamp.Format(time.RFC3339)

	nfInstances, err := util.GetNfInstances(nfProfiles)
	if err != nil {
		log.Errorf("Get nfInstances failed, err:%s", err.Error())
		return
	}

	for _, nfProfile := range nfInstances {
		nfInstanceID := util.GetNfInstanceID(nfProfile)
		if nfInstanceID == "" {
			continue
		}
		subscribeRoamNfProfile(requesterNfType, targetNfType, nfInstanceID, plmnID, timeStampStr)
	}
}

func UnsubscribeNfDeregister(requesterNfType string) {
	if worker.IsKeepCacheMode() {
		log.Info("keep cache mode, not send message to NRF.")
		return
	}

	targetNfs, ok := cache.Instance().GetTargetNfs(requesterNfType)
	if !ok {
		log.Errorf("Get targetNf for nfType[%s], need check configmap deployemnt", requesterNfType)
		return
	}

	for _, targetNf := range targetNfs {
		targetNfType := targetNf.TargetNfType

		subscriptionIDs, ok := cache.Instance().GetSubscriptionIDs(requesterNfType, targetNfType)
		if ok {
			for _, subscriptionID := range subscriptionIDs {
				if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
					err := subscribe.UnSubscribeExecutor(subscriptionID)
					if err != nil {
						log.Errorf("UnSubscribe subscriptionID:%s fail, err:%s", subscriptionID, err.Error())
					}
				}
				cache.Instance().DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
				cache.Instance().DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
			}
		}

		roamSubscriptionIDs, ok := cache.Instance().GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
		if ok {
			for _, subscriptionID := range roamSubscriptionIDs {
				if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
					err := subscribe.UnSubscribeExecutor(subscriptionID)
					if err != nil {
						log.Errorf("UnSubscribe roam subscriptionID:%s fail, err:%s", subscriptionID, err.Error())
					}
				}
				cache.Instance().DelRoamingSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
				cache.Instance().DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
			}
		}
	}

	cache.Instance().UpdateSubscriptionStorage()
}

func UnsubscribeRoamNfProfile(requesterNfType, targetNfType, nfInstanceID string) {
	if worker.IsKeepCacheMode() {
		log.Info("keep cache mode, not send message to NRF.")
		return
	}

	subscriptionID, ok := cache.Instance().GetNfProfileSubscriptionID(requesterNfType, targetNfType, nfInstanceID)
	if ok {
		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			err := subscribe.UnSubscribeExecutor(subscriptionID)
			if err != nil {
				log.Errorf("UnSubscribe roam subscriptionID:%s fail, err:%s", subscriptionID, err.Error())
			}
		}
		cache.Instance().DelRoamingSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
		cache.Instance().DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
	}
}

/////////////////private///////////////

func subscribeRoamNfProfile(requesterNfType string, targetNfType string, nfInstanceID string, plmnID *structs.PlmnID, validityTimeStr string) {
	callbackCluster := util.GetStatusNotifURLs()
	if callbackCluster == "" {
		log.Error("CallbackCluster can not be empty")
		return
	}

	var subscriptionData structs.SubscriptionData

	subscriptionData.NfStatusNotificationURI = callbackCluster + "/nrf-notify-agent/v1/notify/" + requesterNfType + "-roam" + "/" + targetNfType
	subscriptionData.ValidityTime = validityTimeStr
	nfInstanceIdCond := structs.NfInstanceIDCond{
		NfInstanceID: nfInstanceID,
	}
	subscriptionData.SubscrCond = nfInstanceIdCond
	subscriptionData.PlmnID = plmnID

	resp, err := subscribe.SubscribePostExecutor(&subscriptionData)
	if err != nil {
		log.Errorf("Subscribe roam nfProfile:%s fail, err:%s", nfInstanceID, err.Error())
		return
	}

	//the validityTime less 5 seconds compare with NRF-Mgmt
	subscriptionID, validityTime, err := subscribe.SubscribePostRespParser(resp)
	if err != nil {
		log.Errorf("subscribe nfProfile:%s by POST metthod fail, err:%s", nfInstanceID, err.Error())
		return
	}
	log.Debugf("Subscribe nfProfile[%s] by POST method success, subscriptionID[%s], validateTime[%v]", nfInstanceID, subscriptionID, validityTime)

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		validityTime = validityTime.Add(-delayForMaster)
	}

	subscriptionInfo := structs.SubscriptionInfo{
		RequesterNfType: requesterNfType,
		TargetNfType:    targetNfType,
		NfInstanceID:    nfInstanceID,
		TargetPlmnID:    *plmnID,
		SubscriptionID:  subscriptionID,
		ValidityTime:    validityTime,
	}

	cache.Instance().AddRoamingSubscriptionInfo(requesterNfType, targetNfType, subscriptionInfo)
	cache.Instance().SuperviseRoamingSubscription(requesterNfType, targetNfType, subscriptionID, validityTime)

	common.DispatchSubscrInfoToMessageBus(subscriptionInfo)

}

func haveInstersection(nfPlmns []structs.PlmnID, targetPlmns []structs.PlmnID) bool {
	for _, plmnID := range targetPlmns {
		ok := isContain(nfPlmns, plmnID)
		if ok {
			return true
		}
	}

	return false
}

func isSubset(nfPlmns []structs.PlmnID, targetPlmns []structs.PlmnID) bool {
	for _, plmnID := range targetPlmns {
		ok := isContain(nfPlmns, plmnID)
		if !ok {
			return false
		}
	}

	return true
}

func plmnMissProber(nfProfile []byte) bool {
	if nfProfile == nil {
		log.Info("Profile is nil")
		return false
	}

	_, _, _, err := jsonparser.Get(nfProfile, "plmnList")
	if err == nil {
		return false
	} else {
		log.Info("Profile less plmnlist info")
	}

	return true
}

func addPlmnInfo(nfProfile *structs.SearchResultNFProfile, plmnID structs.PlmnID) bool {
	plmns := make([]structs.PlmnID, 0)
	plmns = append(plmns, plmnID)

	nfProfile.PLMN = plmns

	return true
}

func appendPlmnInfo(nfProfile *structs.SearchResultNFProfile, oldNfProfile []byte, plmnID structs.PlmnID) bool {
	var oldProfile structs.SearchResultNFProfile
	err := json.Unmarshal(oldNfProfile, &oldProfile)
	if err != nil {
		log.Errorf("Unmarsh nfProfile fail, err:%s", err.Error())
		return false
	}

	plmns := mergePlmns(nfProfile.PLMN, oldProfile.PLMN, plmnID)
	nfProfile.PLMN = plmns

	return true
}

func isContain(plmns []structs.PlmnID, plmn structs.PlmnID) bool {
	for _, plmnItem := range plmns {
		if plmnItem.Mcc == plmn.Mcc && plmnItem.Mnc == plmn.Mnc {
			return true
		}
	}

	return false
}

func mergePlmns(profilePlmns []structs.PlmnID, oldProfilePlmns []structs.PlmnID, plmnID structs.PlmnID) []structs.PlmnID {
	plmns := make([]structs.PlmnID, 0)
	plmns = append(plmns, profilePlmns...)

	for _, plmn := range oldProfilePlmns {
		ok := isContain(plmns, plmn)
		if !ok {
			plmns = append(plmns, plmn)
		}
	}

	plmnItem := structs.PlmnID{
		Mcc: plmnID.Mcc,
		Mnc: plmnID.Mnc,
	}

	ok := isContain(plmns, plmnItem)
	if !ok {
		plmns = append(plmns, plmnItem)
	}

	return plmns
}
