package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
	"github.com/buger/jsonparser"
)

type Worker struct {
	name     string
	stopFlag bool
}

func (w *Worker) Stop() {
	w.stopFlag = true
}

func (w *Worker) IsRunning() bool {
	return w.stopFlag
}

type SubscribeWorker struct {
	interval int
	nfType   string
	callBack func(nfType string) bool
	Worker
}

func (sw *SubscribeWorker) Start() {
	log.Infof("Launch worker thread : %s", sw.name)
	go func() {
		for !sw.stopFlag {
			time.Sleep(time.Second * time.Duration(sw.interval))
			ret := sw.callBack(sw.nfType)
			if ret {
				break
			}
		}
		if sw.stopFlag {
			log.Infof("Stop worker thread : %s", sw.name)
		}
	}()
}

type FetchProfileWorker struct {
	interval int
	nfType   string
	callBack func(string) bool
	Worker
}

func (fw *FetchProfileWorker) Start() {
	log.Infof("Launch worker thread : %s", fw.name)
	go func() {
		for !fw.stopFlag {
			time.Sleep(time.Second * time.Duration(fw.interval))
			ret := fw.callBack(fw.nfType)
			if ret {
				break
			}
		}
		if fw.stopFlag {
			log.Infof("Stop worker thread : %s", fw.name)
		}
	}()
}

type DumpCacheWorker struct {
	interval int
	nfType   string
	callBack func(string) bool
	Worker
}

func (dw *DumpCacheWorker) Start() {
	log.Infof("Launch worker thread : %s", dw.name)
	go func() {
		for !dw.stopFlag {
			ret := dw.callBack(dw.nfType)
			if ret {
				break
			}
			time.Sleep(time.Second * time.Duration(dw.interval))
		}
		if dw.stopFlag {
			log.Infof("Stop worker thread : %s", dw.name)
		}
	}()
}

var (
	workerManager *WorkerManager
)

const (
	httpContentTypeJSON        = "application/json"
	httpHeaderJSONPatchJSON    = "application/json-patch+json"
	httpContentTypeProblemJSON = "application/problem+json"
	httpResponseFormat         = `{"title": "%s"}`
)

const (
	defaultValidityTime      = 876576 * time.Hour
	defaultTimeDelta         = 5 * time.Second
	defaultTimeDeltaForSlave = 2 * time.Second
)

func init() {
	workerManager = Instance()
}

func loopFetchTargetNfs(nfType string) []structs.TargetNf {
	emptyTargetnf := make([]structs.TargetNf, 0)
	if nfType == "NSSF" {
		//TODO: replace by configmap config
		log.Warningf("need not get targetNfProfiles for NSSF")
		return emptyTargetnf
	}
	for {
		targetNfs, ok := cache.Instance().GetTargetNfs(nfType)
		if !ok {
			log.Warnf("Get targetNfProfiles fail, wait deploy %s targetNf configmap", nfType)
			time.Sleep(time.Second * 5)
		} else {
			return targetNfs
		}
	}
}

func subscribe(nfType string) bool {
	targetNfs, ok := cache.Instance().GetTargetNfs(nfType)
	if !ok {
		log.Warnf("Get targetNfProfiles fail, wait deploy %s targetNf configmap", nfType)
		return false
	}

	ret := false
	doSubscibeFromNrf(targetNfs)
	ret = workerManager.checkSubscribeAllTaskSuccess(nfType)

	return ret
}

func doSubscibeFromNrf(targetNfs []structs.TargetNf) {
	for _, targetNf := range targetNfs {
		subscribeTargetNf(&targetNf)
	}
}

func subscribeTargetNf(targetNf *structs.TargetNf) {
	if len(targetNf.TargetServiceNames) == 0 {
		log.Errorf("subscribe NF targetServiceNames is empty")
		return
	}

	oneSubsData := structs.OneSubscriptionData{
		RequesterNfType: targetNf.RequesterNfType,
		TargetNfType:    targetNf.TargetNfType,
		NotifCondition:  targetNf.NotifCondition,
	}

	for _, serviceName := range targetNf.TargetServiceNames {
		oneSubsData.TargetServiceName = serviceName

		nfType := oneSubsData.RequesterNfType
		subscribeKey := fmt.Sprintf("%s-%s", oneSubsData.TargetNfType, oneSubsData.TargetServiceName)
		ret := workerManager.checkSubscribeTaskSuccess(nfType, subscribeKey)
		if ret {
			continue
		}

		err := subscribeServiceName(&oneSubsData)
		if err != nil {
			log.Errorf("POST subscribe to NRF fail, err:%s", err.Error())
			workerManager.setSubscribeTaskStatus(nfType, subscribeKey, FailureStatus)
		} else {
			workerManager.setSubscribeTaskStatus(nfType, subscribeKey, SuccessStatus)
		}
	}
}

func subscribeServiceName(oneSubsData *structs.OneSubscriptionData) error {
	subscriptionID, validityTime, err := subscribeExecutor(oneSubsData)
	if err != nil {
		log.Errorf("Subscribe to NRF fail, err:%s", err.Error())
		return err
	}

	//consider the exception
	if subscriptionID == "" {
		return nil
	}

	subscriptionInfo := structs.SubscriptionInfo{
		RequesterNfType:   oneSubsData.RequesterNfType,
		TargetNfType:      oneSubsData.TargetNfType,
		TargetServiceName: oneSubsData.TargetServiceName,
		NotifCondition:    oneSubsData.NotifCondition,
		SubscriptionID:    subscriptionID,
		ValidityTime:      validityTime,
	}

	log.Debugf("Discovery Agent subscribe from NRF success, ID:%s", subscriptionID)

	requesterNfType := oneSubsData.RequesterNfType
	targetNfType := oneSubsData.TargetNfType

	cache.Instance().AddSubscriptionInfo(requesterNfType, targetNfType, subscriptionInfo)
	cache.Instance().SuperviseSubscription(requesterNfType, targetNfType, subscriptionID, validityTime)
	cache.Instance().UpdateSubscriptionStorage()

	subscribeKey := fmt.Sprintf("%s-%s", oneSubsData.TargetNfType, oneSubsData.TargetServiceName)
	status := workerManager.fetchSubscribeTaskStatus(oneSubsData.RequesterNfType, subscribeKey)
	if status == FailureStatus {
		common.DispatchSubscrInfoToMessageBus(subscriptionInfo)
	}

	return nil
}

func subscribeExecutor(oneSubsData *structs.OneSubscriptionData) (string, time.Time, error) {
	if IsKeepCacheMode() {
		log.Info("NRF Disc Agent is Master, Woke Mode is KeepCache Mode, So no need to do subscription")
		return "", time.Time{}, nil
	}

	fqdn, exists := cache.Instance().GetRequesterFqdn(oneSubsData.RequesterNfType)
	if !exists {
		log.Warnf("%s nf instance was deregistered", oneSubsData.RequesterNfType)
		return "", time.Time{}, nil
	}

	var subscribeData []byte
	subscribeData = util.BuildSubscriptionPostData(oneSubsData, fqdn)

	if subscribeData == nil {
		return "", time.Time{}, fmt.Errorf("Build subscription POST data for %s fail", oneSubsData.RequesterNfType)
	}

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfMgmt("h2", "POST", "subscriptions", hdr, bytes.NewBuffer(subscribeData))
	if err != nil {
		log.Errorf("Failed to send subscription request to NRF, %s", err.Error())
		return "", time.Time{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		log.Errorf("Subscribe from NRF by POST method fail, statusCode:%d, body:%s", resp.StatusCode, string(resp.Body))
		return "", time.Time{}, fmt.Errorf("Subscribe from NRF for %s fail", oneSubsData.RequesterNfType)
	}

	subscriptionID, timeStamp, err := subscribeResponseParser(resp)
	if err != nil {
		return "", time.Time{}, err
	}

	return subscriptionID, timeStamp, nil
}

func subscribeResponseParser(resp *httpclient.HttpRespData) (string, time.Time, error) {
	location := resp.Location
	log.Debugf("subscribe response location:%s", location)

	subURL := strings.Split(location, "//")
	if len(subURL) < 2 {
		log.Errorf("subscriptionID in location is error")
		return "", time.Time{}, fmt.Errorf("subscriptionID in location is error")
	}
	subURLSuffix := strings.Split(subURL[1], "/")
	if len(subURLSuffix) < 5 {
		log.Errorf("subscriptionID in location is error")
		return "", time.Time{}, fmt.Errorf("subscriptionID in location is error")
	}
	subscriptionID := subURLSuffix[3] + "/" + subURLSuffix[4]

	defaultTime := time.Now().Add(defaultValidityTime)
	vt, err := jsonparser.GetString(resp.Body, "validityTime")
	if err != nil {
		log.Warnf("no validityTime in subscription boby")
		return subscriptionID, defaultTime, nil
	}
	validityTime, err := time.Parse(time.RFC3339, vt)
	if err != nil {
		log.Warnf("validityTime %s is invalid when parse by %s", vt, time.RFC3339)
		return subscriptionID, defaultTime, nil
	}

	return subscriptionID, validityTime.Add(-defaultTimeDelta), nil
}

//////////////////fetch profile data/////////////

func fetchProfile(nfType string) bool {
	targetNfs, ok := cache.Instance().GetTargetNfs(nfType)
	if !ok {
		log.Warnf("Get targetNfProfiles fail, wait deploy %s targetNf configmap", nfType)
		return false
	}

	ret := false
	doFetchProfileFromNrf(targetNfs)
	ret = workerManager.checkFetchProfileAllTaskSuccess(nfType)

	return ret
}

func doFetchProfileFromNrf(targetNfs []structs.TargetNf) {
	for _, targetNf := range targetNfs {
		nfType := targetNf.RequesterNfType
		fetchProfileKey := fmt.Sprintf("%s", targetNf.TargetNfType)
		ret := workerManager.checkFetchProfileTaskSuccess(nfType, fetchProfileKey)
		if ret {
			continue
		}

		err := fetchProfileTargetNf(targetNf)
		if err != nil {
			log.Errorf("Fetch profile for targetNf:%+v fail, err:%s", targetNf, err.Error())
			workerManager.setFetchProfileTaskStatus(nfType, fetchProfileKey, FailureStatus)
		} else {
			workerManager.setFetchProfileTaskStatus(nfType, fetchProfileKey, SuccessStatus)
		}
	}
}

func fetchProfileTargetNf(targetNf structs.TargetNf) error {
	fqdn, exists := cache.Instance().GetRequesterFqdn(targetNf.RequesterNfType)
	if !exists {
		log.Warnf("%s nf instance was deregistered", targetNf.RequesterNfType)
		return nil
	}

	query := util.GetDiscoveryRequestURL(&targetNf, "", fqdn)
	if query == "" {
		return fmt.Errorf("Failed to get Discovery request URL")
	}

	log.Debugf("Master fetch profile url:%s", query)
	err := fetchProfileExecutor(targetNf.RequesterNfType, targetNf.TargetNfType, query)

	return err
}

func fetchProfileExecutor(requesterNfType string, targetNfType string, query string) error {
	if IsKeepCacheMode() {
		log.Info("NRF Disc Agent is Master, Woke Mode is KeepCache Mode, So no need to fetch profile from NRF")
		return nil
	}

	if _, exists := cache.Instance().GetRequesterFqdn(requesterNfType); !exists {
		log.Warnf("%s nf instance was deregistered", requesterNfType)
		return nil
	}

	//pm.Inc(consts.NrfDiscoveryRequestsTotal)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		log.Errorf("Send discovery request to NRF fail, %s", err.Error())
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= http.StatusUnavailableForLegalReasons {
		log.Infof("NRF response code is: %+v and Master Agent will be ready", resp.StatusCode)
		return nil
	}

	err = handleDiscoveryResponse(requesterNfType, targetNfType, resp)
	if err != nil {
		return err
	}
	/*
		if err != nil {
			fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
			return err
		}
		fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	*/

	fetchProfileKey := fmt.Sprintf("%s", targetNfType)
	status := workerManager.fetchFetchProfileTaskStatus(requesterNfType, fetchProfileKey)
	if status == FailureStatus {
		util.PushMessageToMSB(requesterNfType, targetNfType, "", consts.NFEventDiscResult, resp.Body)
	}

	return nil
}

func handleDiscoveryResponse(nfType string, targetNfType string, resp *httpclient.HttpRespData) error {
	//pmNrfDiscoveryResponses(resp.StatusCode)
	log.Infof("NRF response: %s", resp.SimpleString())
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
	nfInstances, validityPeriod, ok := cache.SpliteSeachResult(resp.Body)
	if !ok {
		return fmt.Errorf("invalid SearchResult in NRF response, %+v", string(resp.Body))
	}
	for _, nfProfile := range nfInstances {
		cache.Instance().CachedWithTTL(nfType, targetNfType, nfProfile, validityPeriod, false)
	}
	cache.Instance().SetCacheStatus(nfType, targetNfType, true)

	return nil
}

func dumpCacheFromMaster(nfType string) bool {
	_, ok := cache.Instance().GetTargetNfs(nfType)
	if !ok {
		log.Warnf("Get targetNfProfiles fail, wait deploy %s targetNf configmap", nfType)
		return false
	}

	if _, exists := cache.Instance().GetRequesterFqdn(nfType); !exists {
		log.Warnf("%s nf instance was deregistered", nfType)
		return true
	}

	log.Infof("Fetch %s cache profile from master NRF Discovery Agent", nfType)

	masterURL := util.GetLeaderDiscURL()
	if masterURL == "" {
		log.Error("Get master URL fail")
		return false
	}
	query := "synccache/" + nfType

	queryUrl := masterURL + query
	log.Debugf("Master agent url:%s", queryUrl)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDo("h2", "GET", queryUrl, hdr, nil)
	if err != nil {
		log.Errorf("Send dump request to master discovery agent fail,err:%s", err.Error())
		return false
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Master agent response dump query fail, StatusCode:%d, body:%s", resp.StatusCode, string(resp.Body))
		return false
	}

	log.Debugf("Master agent response:%s", resp.SimpleString())
	err = handleDumpResponse(nfType, resp)
	if err != nil {
		log.Warnf("Parse dump response from mater fail, err:%s", err.Error())
		return false
	}

	return true
}

func handleDumpResponse(nfType string, resp *httpclient.HttpRespData) error {
	dumpContent := resp.Body

	cacheData := structs.CacheSyncData{
		RequestNfType: nfType,
	}

	err := json.Unmarshal(dumpContent, &cacheData)
	if err != nil {
		log.Errorf("Marshal dump cache data to CacheSyncData fail: err:%s", err.Error())
		return err
	}

	for _, targetCache := range cacheData.CacheInfos {
		targetNfType := targetCache.TargetNfType

		if len(targetCache.NfProfiles) >= 0 {
			fetchProfileKey := fmt.Sprintf("%s", targetNfType)
			workerManager.setFetchProfileTaskStatus(nfType, fetchProfileKey, SuccessStatus)
		}

		for _, profile := range targetCache.NfProfiles {
			nfInstanceID := util.GetNfInstanceID(profile)
			if nfInstanceID == "" {
				continue
			}

			validityTime := getTtlTimeStamp(targetCache.TtlInfos, nfInstanceID)
			cache.Instance().CachedWithTtlTimestamp(nfType, targetNfType, profile, validityTime, false)
		}

		for _, etagInfo := range targetCache.EtagInfos {
			cache.Instance().SaveEtag(nfType, targetNfType, etagInfo.NfInstanceID, etagInfo.FingerPrint)
		}

		for _, subscriptionInfo := range targetCache.SubscriptionInfos {
			cache.Instance().AddSubscriptionInfo(nfType, targetNfType, subscriptionInfo)
			subscribeKey := fmt.Sprintf("%s-%s", subscriptionInfo.TargetNfType, subscriptionInfo.TargetServiceName)
			workerManager.setSubscribeTaskStatus(nfType, subscribeKey, SuccessStatus)

			requesterNfType := subscriptionInfo.RequesterNfType
			targetNfType := subscriptionInfo.TargetNfType
			subscriptionID := subscriptionInfo.SubscriptionID
			timepoint := subscriptionInfo.ValidityTime
			cache.Instance().SuperviseSubscription(requesterNfType, targetNfType, subscriptionID, timepoint)
		}
	}

	for _, targetCache := range cacheData.RoamingCacheInfos {
		targetNfType := targetCache.TargetNfType

		if len(targetCache.NfProfiles) >= 0 {
			fetchProfileKey := fmt.Sprintf("%s", targetNfType)
			workerManager.setFetchProfileTaskStatus(nfType, fetchProfileKey, SuccessStatus)
		}

		for _, profile := range targetCache.NfProfiles {
			nfInstanceID := util.GetNfInstanceID(profile)
			if nfInstanceID == "" {
				continue
			}

			validityTime := getTtlTimeStamp(targetCache.TtlInfos, nfInstanceID)
			cache.Instance().CachedWithTtlTimestamp(nfType, targetNfType, profile, validityTime, true)
		}

		for _, etagInfo := range targetCache.EtagInfos {
			cache.Instance().SaveRoamingEtag(nfType, targetNfType, etagInfo.NfInstanceID, etagInfo.FingerPrint)
		}

		for _, subscriptionInfo := range targetCache.SubscriptionInfos {
			cache.Instance().AddRoamingSubscriptionInfo(nfType, targetNfType, subscriptionInfo)

			requesterNfType := subscriptionInfo.RequesterNfType
			targetNfType := subscriptionInfo.TargetNfType
			subscriptionID := subscriptionInfo.SubscriptionID
			timepoint := subscriptionInfo.ValidityTime
			cache.Instance().SuperviseRoamingSubscription(requesterNfType, targetNfType, subscriptionID, timepoint)
		}
	}

	return nil
}

func getTtlTimeStamp(ttlInfo []structs.TtlInfo, nfInstanceID string) time.Time {
	for _, ttl := range ttlInfo {
		if ttl.NfInstanceID == nfInstanceID {
			return ttl.ValidityTime
		}
	}

	return time.Time{}
}
