package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/k8sapiclient"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/subscribe"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
	"github.com/buger/jsonparser"
)

//SpliteSeachResult split the search result from NRF
func SpliteSeachResult(content []byte) ([][]byte, uint, bool) {
	var nfinstanceArray [][]byte
	validityPeriod, err := jsonparser.GetInt(content, "validityPeriod")
	if err != nil {
		log.Errorf("The search result does not contain validityPeriod")
		return nil, 0, false
	}
	if validityPeriod < 0 {
		log.Errorf("The search result validityPeriod < 0")
		return nil, 0, false
	}

	nfinstanceArray = make([][]byte, 0)
	_, err = jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		nfinstanceArray = append(nfinstanceArray, value)
	}, "nfInstances")
	if err != nil {
		log.Errorf("Get nfInstances from search result failure")
		return nil, 0, false
	}

	return nfinstanceArray, uint(validityPeriod), true
}

func getNfType(nfType, nfInstanceID string) (string, error) {
	content := Instance().FetchNfProfile(nfType, nfInstanceID)
	if content == nil {
		return "", fmt.Errorf("Get nfProfile[%s] from cache failure", nfInstanceID)
	}

	nfTypeRaw, err := jsonparser.GetString(content, "nfType")
	if err != nil {
		return "", err
	}

	return nfTypeRaw, nil
}

func proberNfProfile(requesterNfType, targetNfType, nfInstanceID string) (uint, bool) {
	targetNf, exists := Instance().GetTargetNf(requesterNfType, targetNfType)
	if !exists {
		log.Errorf("requsterNfType:%s targetNfType:%s do not deploy in configmap", requesterNfType, targetNfType)
		return 0, false
	}

	targetNfServices := ""
	for _, service := range targetNf.TargetServiceNames {
		targetNfServices = targetNfServices + "service-names=" + service + "&"
	}
	queryURL := "nf-instances?" +
		"requester-nf-type=" + requesterNfType +
		"&target-nf-type=" + targetNfType +
		"&" + targetNfServices +
		"target-nf-instance-id=" + nfInstanceID

	requesterFqdn, exists := Instance().GetRequesterFqdn(requesterNfType)
	if !exists {
		log.Errorf("%s was not registered to Register Agent", requesterNfType)
		return 0, false
	}
	if requesterFqdn != "" {
		queryURL += "&requester-nf-instance-fqdn=" + requesterFqdn
	}
	log.Infof("query url:%s", queryURL)

	pm.Inc(consts.NrfDiscoveryRequestsTotal)
	hdr := make(map[string]string)
	hdr["Content-Type"] = "application/json"
	if Instance().HaveEtag(requesterNfType, targetNfType, nfInstanceID) {
		hdr["If-None-Match"] = Instance().FetchEtag(requesterNfType, targetNfType, nfInstanceID)
	}
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", queryURL, hdr, nil)
	if err != nil {
		log.Errorf("Failed to send Discover request to NRF, Error: \"%s\"", err.Error())
		return 0, false
	}
	nrfDiscoveryResponsesPmHandler(resp.StatusCode)

	log.Infof("NRF resp:%s", resp.SimpleString())
	if resp.StatusCode == http.StatusNotModified {
		log.Infof("Discover etag statusNotModified, StatusCode: \"%d\"", resp.StatusCode)
		ttlValue := 86400
		cacheControl := resp.Header.Get("Cache-Control")
		if cacheControl != "" {
			ttls := strings.Split(cacheControl, "=")
			ttl, err := strconv.Atoi(ttls[1])
			if err != nil {
				log.Warnf("Nrf response header less Cache-Control item")
			} else {
				ttlValue = ttl
			}
		}
		log.Infof("CacheMonitor will update [%s] ttl[%d]\n", nfInstanceID, ttlValue)
		return uint(ttlValue), true
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Discover failed, StatusCode: \"%d\"", resp.StatusCode)
		return 0, false
	}

	results, err := common.ConvertIpv6ToIpv4InSearchResult(resp.Body, cm.IsEnableConvertIpv6ToIpv4())
	if err != nil {
		return 0, false
	}
	log.Infof("NRF response: %+v", string(results))

	nfInstances, validityPeriod, ok := SpliteSeachResult(results)
	if !ok {
		log.Error("Invalid response from NRF Discovery")
		return 0, false
	}
	if len(nfInstances) != 1 {
		log.Error("NRF response should contain one and at most one NF profile")
		return 0, false
	}

	//save etag value
	value := resp.Header.Get("ETag")
	if value != "" {
		value = strings.Replace(value, "W/", "", -1)
		log.Infof("CacheManager will save etag[%s] for %s", value, nfInstanceID)
		Instance().SaveEtag(requesterNfType, targetNfType, nfInstanceID, value)
	}
	for _, nfProfile := range nfInstances {
		Instance().Cached(requesterNfType, targetNfType, nfProfile, false)
	}

	if election.IsActiveLeader("3201", consts.DiscoveryAgentReadinessProbe) {
		util.PushMessageToMSB(requesterNfType, targetNfType, nfInstanceID, consts.NFProfileChg, results)
	}

	log.Infof("NfinstanceID:%s validityPeriod:%d", nfInstanceID, validityPeriod)
	return validityPeriod, true
}

func nrfDiscoveryResponsesPmHandler(code int) {
	if code >= 200 && code <= 299 {
		pm.Inc(consts.NrfDiscoveryResponses2xx)
	} else if code >= 300 && code <= 399 {
		pm.Inc(consts.NrfDiscoveryResponses3xx)
	} else if code >= 400 && code <= 499 {
		pm.Inc(consts.NrfDiscoveryResponses4xx)
	} else if code >= 500 && code <= 599 {
		pm.Inc(consts.NrfDiscoveryResponses5xx)
	}
	pm.Inc(consts.NrfDiscoveryResponsesTotal)
}

/////////////subscription//////////////

var (
	defaultValidityTime      = 876576 * time.Hour
	defaultTimeDelta         = 5 * time.Second
	defaultTimeDeltaForSlave = 2 * time.Second
	//defaultMasterTimeDelta = 5 * time.Second
	//defaultSlaveTimeDelta  = 3 * time.Second // slave do need delay

	httpContentTypeJSON        = "application/json"
	httpHeaderJSONPatchJSON    = "application/json-patch+json"
	httpContentTypeProblemJSON = "application/problem+json"
)

func handleSubscriptionResponse(resp *httpclient.HttpRespData) (string, *time.Time) {
	location := resp.Location
	log.Debugf("Subscription response location is %s", location)

	subURL := strings.Split(location, "//")
	if len(subURL) < 2 {
		log.Errorf("SubscriptionID in location is error")
		return "", nil
	}
	subURLSuffix := strings.Split(subURL[1], "/")
	if len(subURLSuffix) < 5 {
		log.Errorf("SubscriptionID in location is error")
		return "", nil
	}
	subscriptionID := subURLSuffix[3] + "/" + subURLSuffix[4]

	timestamp, err := jsonparser.GetString(resp.Body, "validityTime")
	if err != nil {
		log.Warnf("Get validityTime in subscription boby fail, err=%s", err.Error())
		return "", nil
	}
	validityTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Warnf("Parse timestamp[%s] by RFC3339 fail, err=%s", timestamp, err.Error())
		return "", nil
	}
	validityTime = validityTime.Add(-defaultTimeDelta)

	return subscriptionID, &validityTime
}

func doSubscriptionToNrf(oneSubsData *structs.OneSubscriptionData) (string, *time.Time) {
	requesterNfFqdn, exists := Instance().GetRequesterFqdn(oneSubsData.RequesterNfType)
	if !exists {
		log.Warnf("Disc agent don't receive %s FQDN from Register Agent", oneSubsData.RequesterNfType)
		return "", nil
	}

	subscriptionData := util.BuildSubscriptionPostData(oneSubsData, requesterNfFqdn)
	if subscriptionData == nil {
		return "", nil
	}

	hdr := make(map[string]string)
	var resp *httpclient.HttpRespData
	hdr["Content-Type"] = httpContentTypeJSON

	log.Infof("Do subscribe from NRF, POST data:[%s]", string(subscriptionData))
	resp, err := client.HTTPDoToNrfMgmt("h2", "POST", "subscriptions", hdr, bytes.NewBuffer(subscriptionData))
	if err != nil {
		log.Errorf("Do subscribe by POST from NRF fail,err:%s", err.Error())
		return "", nil
	}
	log.Infof("Subscribe POST response from NRF : %s", string(resp.Body))
	if resp.StatusCode != http.StatusCreated {
		log.Errorf("Do subscribe by POST from NRF fail:code[%d], body[%s]", resp.StatusCode, string(resp.Body))
		return "", nil
	}

	return handleSubscriptionResponse(resp)
}

func doRoamSubscriptionToNrf(subscriptionInfo structs.SubscriptionInfo) (string, time.Time, error) {
	var errInfo string

	callbackCluster := util.GetStatusNotifURLs()
	if callbackCluster == "" {
		errInfo = "CallbackCluster can not be empty"
		return "", time.Time{}, fmt.Errorf("%s", errInfo)
	}
	/*
		plmnID, ok := Instance().GetRequesterPlmns(subscriptionInfo.RequesterNfType)
		if !ok {
			errInfo = fmt.Sprintf("Get plmnID for nfType:%s fail", subscriptionInfo.RequesterNfType)
			return "", time.Time{}, fmt.Errorf("%s", errInfo)
		}
	*/
	nfInstanceID := subscriptionInfo.NfInstanceID
	requesterNfType := subscriptionInfo.RequesterNfType
	targetNfType := subscriptionInfo.TargetNfType
	targetPlmnID := subscriptionInfo.TargetPlmnID

	var subscriptionData structs.SubscriptionData

	subscriptionData.NfStatusNotificationURI = callbackCluster + "/nrf-notify-agent/v1/notify/" + requesterNfType + "-roam" + "/" + targetNfType
	nfInstanceIdCond := structs.NfInstanceIDCond{
		NfInstanceID: nfInstanceID,
	}
	subscriptionData.SubscrCond = nfInstanceIdCond
	subscriptionData.PlmnID = &targetPlmnID

	resp, err := subscribe.SubscribePostExecutor(&subscriptionData)
	if err != nil {
		log.Errorf("Subscribe roam nfProfile:%s fail, err:%s", nfInstanceID, err.Error())
		return "", time.Time{}, err
	}

	subscriptionID, validityTime, err := subscribe.SubscribePostRespParser(resp)
	if err != nil {
		log.Errorf("subscribe nfProfile:%s by POST metthod fail, err:%s", nfInstanceID, err.Error())
		return "", time.Time{}, err
	}
	log.Debugf("Subscribe nfProfile[%s] by POST method success, subscriptionID[%s], validateTime[%v]", nfInstanceID, subscriptionID, validityTime)

	return subscriptionID, validityTime, nil
	/*
		subscriptionInfoNew := structs.SubscriptionInfo{
			RequesterNfType: requesterNfType,
			TargetNfType:    targetNfType,
			NfInstanceID:    nfInstanceID,
			TargetPlmnID:    targetPlmnID,
			SubscriptionID:  subscriptionID,
			ValidityTime:    validityTime,
		}

		Instance().AddRoamingSubscriptionInfo(requesterNfType, targetNfType, subscriptionInfoNew)
		Instance().SuperviseRoamingSubscription(requesterNfType, targetNfType, subscriptionID, validityTime)

		common.DispatchSubscrInfoToMessageBus(subscriptionInfo)
	*/
}

func prolongSubscriptionFromNrf(oneSubsData *structs.OneSubscriptionData, subscriptionID string) *time.Time {
	requesterNfType := oneSubsData.RequesterNfType
	targetNfType := oneSubsData.TargetNfType
	//targetServiceName := oneSubsData.TargetServiceName

	targetNf, ok := Instance().GetTargetNf(requesterNfType, targetNfType)
	if !ok {
		log.Errorf("Failed to getTargetNf for requestNfType[%s] and targetNfType[%s]", requesterNfType, targetNfType)
		return nil
	}

	subscriptionPatchBody := util.BuildSubscriptionPatchData(targetNf.SubscriptionValidityTime)
	if subscriptionPatchBody == nil {
		log.Error("Build subscribe patch data fail")
		return nil
	}
	log.Infof("Prolong subscribe ttl from NRF, PATCH data:[%s]", string(subscriptionPatchBody))

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpHeaderJSONPatchJSON
	var resp *httpclient.HttpRespData
	var err error

	//resp, err = client.HTTPDoToNrfMgmt("h2", "PATCH", "subscriptions/"+subscriptionID, hdr, bytes.NewBuffer(subscriptionPatchBody))
	resp, err = client.HTTPDoToNrfMgmt("h2", "PATCH", subscriptionID, hdr, bytes.NewBuffer(subscriptionPatchBody))
	if err != nil {
		log.Errorf("Send subscribe request to NRF fail,err:%s", err.Error())
		return nil
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Errorf("Prolong subscription[%s] by PATCH fail, code[%d], body[%s]", subscriptionID, resp.StatusCode, string(resp.Body))
		return nil
	}
	log.Infof("Subscribe PATCH response body from NRF:[%s]", string(resp.Body))

	var validityTime time.Time
	if resp.StatusCode == http.StatusNoContent {
		now := time.Now()
		prolongValue := time.Duration(targetNf.SubscriptionValidityTime) * time.Second
		nextTimeStamp := now.Add(prolongValue)
		validityTime = nextTimeStamp
	} else if resp.StatusCode == http.StatusOK {
		timestamp, err := jsonparser.GetString(resp.Body, "validityTime")
		if err != nil {
			log.Warnf("Get validityTime in subscription boby fail, err=%s", err.Error())
			return nil
		}
		validityTime, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			log.Warnf("Parse timestamp[%s] by RFC3339 fail, err=%s", timestamp, err.Error())
			return nil
		}
	}

	validityTime = validityTime.Add(-defaultTimeDelta)
	return &validityTime
}

func prolongRoamSubscriptionFromNrf(subscriptionID string) *time.Time {
	var defaultTtl int = 2592000

	now := time.Now()
	prolongValue := time.Duration(defaultTtl) * time.Second
	nextTimeStamp := now.Add(prolongValue)
	nextTimeStampStr := nextTimeStamp.Format(time.RFC3339)

	item := structs.PatchItem{
		Op:    "replace",
		Path:  "/validityTime",
		Value: nextTimeStampStr,
	}
	patchItems := make([]structs.PatchItem, 0)
	patchItems = append(patchItems, item)

	log.Infof("Prolong subscribe ttl from NRF, PATCH data:[%v]", patchItems)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpHeaderJSONPatchJSON

	resp, err := subscribe.SubscribePatchExecutor(subscriptionID, patchItems)
	if err != nil {
		log.Errorf("Subscribe subscriptionID:%s by PATCH method fail, err:%s", subscriptionID, err.Error())
		return nil
	}

	validityTime, err := subscribe.SubscribePatchRespParser(resp, nextTimeStamp)
	if err != nil {
		log.Errorf("Parse subscribe PATCH repsonse fail, err:%s", err.Error())
		return nil
	}

	return &validityTime
}

func unsubscribeByNfInstanceID(requesterNfType, targetNfType, nfInstanceID string) error {
	subscriptionID, ok := Instance().GetNfProfileSubscriptionID(requesterNfType, targetNfType, nfInstanceID)
	if !ok {
		errInfo := fmt.Sprintf("Get subscriptionID by requesterNfType:%s targetNfType:%s nfInstanceID:%s failure", requesterNfType, targetNfType, nfInstanceID)
		return fmt.Errorf("%s", errInfo)
	}

	return subscribe.UnSubscribeExecutor(subscriptionID)
}

func fetchSubscriptionInfoFromMaster(subscriptionInfo *structs.SubscriptionInfo, isRoam bool) (string, *time.Time) {
	if subscriptionInfo == nil {
		return "", nil
	}

	subIDURLPrefix := util.GetLeaderDiscURL()
	if subIDURLPrefix == "" {
		return "", nil
	}

	var subIDURL string
	if isRoam {
		subIDURL = subIDURLPrefix + "subscriptions?" +
			consts.SearchDataRequesterNfType + "=" + subscriptionInfo.RequesterNfType + "&" +
			consts.SearchDataTargetNfType + "=" + subscriptionInfo.TargetNfType + "&" +
			consts.SearchDataTargetInstID + "=" + subscriptionInfo.NfInstanceID
	} else {
		subIDURL = subIDURLPrefix + "roam-subscriptions?" +
			consts.SearchDataRequesterNfType + "=" + subscriptionInfo.RequesterNfType + "&" +
			consts.SearchDataTargetNfType + "=" + subscriptionInfo.TargetNfType + "&" +
			consts.SearchDataServiceName + "=" + subscriptionInfo.TargetServiceName
	}
	log.Infof("Subscription request to master url %s", subIDURL)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDo("h2", "GET", subIDURL, hdr, nil)
	if err != nil {
		log.Errorf("Failed to send request to master agent, err:%s", err.Error())
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Master agent response:%s", resp.SimpleString())
		return "", nil
	}

	log.Infof("Master agent response body:%s", string(resp.Body))
	var subsInfo structs.SubscriptionInfo
	err = json.Unmarshal(resp.Body, &subsInfo)
	if err != nil {
		log.Errorf("Unmarsh response body to SubscriptionInfo fail, err=%s", err.Error())
		return "", nil
	}
	validityTime := subsInfo.ValidityTime.Add(-defaultTimeDeltaForSlave)

	return subsInfo.SubscriptionID, &validityTime
}

func updateConfigmapStorage(subscriptionInfoMap map[string]structs.SubscriptionInfo) bool {
	if subscriptionInfoMap == nil {
		return false
	}

	// ONLY master agent should write data to configmap storage
	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		jsonBuf, err := json.Marshal(subscriptionInfoMap)
		if err != nil {
			log.Errorf("failed to marshal subscription info %+v, %s", subscriptionInfoMap, err.Error())
			return false
		}

		err = k8sapiclient.SetConfigMapData(consts.ConfigMapStorage, consts.ConfigMapKeySubsInfo, jsonBuf)
		if err != nil {
			log.Errorf("failed to write FQDN to configmap %s, %s", consts.ConfigMapStorage, err.Error())
			return false
		}
	}
	log.Debugf("updateConfigmapStorage done")

	return true
}

func SyncNrfData(targetNf *structs.TargetNf, isSwitchCache bool, cache *cache) bool {
	requestOptions := ""
	/*
		if nfInstanceID != "" {
			requestOptions = consts.SearchDataTargetInstID + "=" + nfInstanceID
		}
	*/
	requesterNfFqdn, exists := Instance().GetRequesterFqdn(targetNf.RequesterNfType)
	if !exists {
		log.Warnf("Disc agent don't receive %s FQDN from Register Agent", targetNf.RequesterNfType)
		return false
	}

	query := util.GetDiscoveryRequestURL(targetNf, requestOptions, requesterNfFqdn)
	if query == "" {
		log.Error("Get Discovery request URL fail")
		return false
	}

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		if _, err := handleDiscoveryRequestToNrf(targetNf, query, isSwitchCache, cache); err != nil {
			log.Errorf("Sync data from NRF-Disc by query[%v] fail, err:%s", *targetNf, err.Error())
			return false
		}
	} else {
		log.Debug("NRF-Agent-Disc slave node sync data by messageBus from active node")
	}

	return true
}

func handleDiscoveryRequestToNrf(targetNf *structs.TargetNf, query string, isSwitchCache bool, cache *cache) (*httpclient.HttpRespData, error) {
	//pm.Inc(consts.NrfDiscoveryRequestsTotal)

	log.Infof("Try to fetch %s %+v profile from NRF Discovery", targetNf.TargetNfType, targetNf.TargetServiceNames)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		log.Errorf("failed to send discovery request to NRF, %s", err.Error())
		return nil, err
	}
	err = handleDiscoveryResponse(targetNf, resp, isSwitchCache, cache)
	if err != nil {
		//fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return nil, err
	}
	//fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	util.PushMessageToMSB(targetNf.RequesterNfType, targetNf.TargetNfType, "", consts.NFEventDiscResult, resp.Body)
	return resp, nil
}

func handleDiscoveryResponse(targetNf *structs.TargetNf, resp *httpclient.HttpRespData, isSwitchCache bool, cache *cache) error {
	//pmNrfDiscoveryResponses(resp.StatusCode)

	log.Infof("NRF response %+v", string(resp.Body))
	if resp.StatusCode != http.StatusOK {
		log.Errorf("NRF response: %+v", string(resp.Body))
		return errors.New(string(resp.Body))
	}
	var err error
	resp.Body, err = common.ConvertIpv6ToIpv4InSearchResult(resp.Body, cm.IsEnableConvertIpv6ToIpv4())
	if err != nil {
		log.Errorf("Failed to convert Ipv6Address to Ipv4Address in NF profile, %s", err.Error())
		return err
	}

	//log.Infof("NRF response: %+v", string(resp.Body))
	nfInstances, validityPeriod, ok := SpliteSeachResult(resp.Body)
	if !ok {
		return fmt.Errorf("invalid SearchResult in NRF response, %+v", string(resp.Body))
	}
	for _, nfProfile := range nfInstances {
		//build cache with input cache ptr when isSwitchCache flag set
		if isSwitchCache {
			if cache == nil {
				log.Error("Discovery response : isSwitchCache is true, but cache is nil ")
				return errors.New("Discovery response : isSwitchCache is true, but cache is nil")
			}
			Instance().InjectionCachedWithTtl(cache, nfProfile, validityPeriod)
		} else {
			Instance().CachedWithTTL(targetNf.RequesterNfType, targetNf.TargetNfType, nfProfile, validityPeriod, false)
		}
	}

	return nil
}
