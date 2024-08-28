package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"github.com/buger/jsonparser"
)

/*
func pushMessageToMSB(requesterNfType, targetNfType, nfInstanceID string, event string, resp []byte) bool {
	var msgBusDisc structs.NotificationMsgBus
	msgBusDisc.AgentProducerID = common.GetSelfUUID()
	msgBusDisc.NfEvent = event
	msgBusDisc.NfType = targetNfType
	msgBusDisc.NfInstanceID = nfInstanceID

	err := json.Unmarshal(resp, &msgBusDisc.MessageBody)
	if err != nil {
		log.Errorf("Decode MessageBody message Unmarshal fail, %s", err.Error())
		return false
	}
	log.Infof("RequesterNfType %s, msgBusDisc %+v", requesterNfType, msgBusDisc)
	jsonBuf, err := json.Marshal(msgBusDisc)
	if err != nil {
		log.Errorf("Failed to Marshal Disc message, %s", err.Error())
		return false
	}

	topic := consts.MsgbusTopicNamePrefix + strings.ToLower(requesterNfType)
	if discMsgbus := common.GetDiscMsgbus(); discMsgbus != nil {
		err := discMsgbus.SendMessage(topic, string(jsonBuf))
		if err != nil {
			log.Errorf("Failed to send notification to message bus, %s", err.Error())
		} else {
			log.Debugf("Push message to MSB succeed")
		}
	} else {
		log.Warnf("Message bus was not initialized")
	}

	return true
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

func syncNrfData(targetNf *structs.TargetNf, nfInstanceID string) bool {
	requestOptions := ""
	if nfInstanceID != "" {
		requestOptions = consts.SearchDataTargetInstID + "=" + nfInstanceID
	}
	query := getDiscoveryRequestURL(targetNf, requestOptions)
	if query == "" {
		log.Error("Get Discovery request URL fail")
		return false
	}

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		if _, err := handleDiscoveryRequestToNrf(targetNf, query); err != nil {
			log.Errorf("Sync data from NRF-Disc by query[%v] fail, err:%s", *targetNf, err.Error())
			return false
		}
	} else {
		log.Debug("NRF-Agent-Disc slave node sync data by messageBus from active node")
	}

	return true
}

func handleDiscoveryRequestToNrf(targetNf *structs.TargetNf, query string) (*httpclient.HttpRespData, error) {
	//pm.Inc(consts.NrfDiscoveryRequestsTotal)

	log.Infof("Try to fetch %s %+v profile from NRF Discovery", targetNf.TargetNfType, targetNf.TargetServiceNames)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		log.Errorf("failed to send discovery request to NRF, %s", err.Error())
		return nil, err
	}
	err = handleDiscoveryResponse(targetNf, resp)
	if err != nil {
		//fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return nil, err
	}
	//fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	pushSyncDataToMSB(targetNf.RequesterNfType, targetNf.TargetNfType, "", consts.NFEventDiscResult, resp.Body)
	return resp, nil
}

func handleDiscoveryResponse(targetNf *structs.TargetNf, resp *httpclient.HttpRespData) error {
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
		Instance().CachedWithTTL(targetNf.RequesterNfType, nfProfile, validityPeriod)
	}

	return nil
}

func pushSyncDataToMSB(requesterNfType, targetNfType, nfInstanceID string, event string, resp []byte) bool {
	var msgBusDisc structs.NotificationMsgBus
	msgBusDisc.AgentProducerID = common.GetSelfUUID()
	msgBusDisc.NfEvent = event
	msgBusDisc.NfType = targetNfType
	msgBusDisc.NfInstanceID = nfInstanceID

	err := json.Unmarshal(resp, &msgBusDisc.MessageBody)
	if err != nil {
		log.Errorf("decode MessageBody message Unmarshal fail, %s", err.Error())
		return false
	}
	log.Infof("requesterNfType %s, msgBusDisc %+v", requesterNfType, msgBusDisc)
	jsonBuf, err := json.Marshal(msgBusDisc)
	if err != nil {
		log.Errorf("failed to Marshal Disc message, %s", err.Error())
		return false
	}

	topic := consts.MsgbusTopicNamePrefix + strings.ToLower(requesterNfType)
	if discMsgbus := common.GetDiscMsgbus(); discMsgbus != nil {
		err := discMsgbus.SendMessage(topic, string(jsonBuf))
		if err != nil {
			log.Errorf("failed to send notification to message bus, %s", err.Error())
		} else {
			log.Debugf("push message to MSB succeed")
		}
	} else {
		log.Warnf("message bus was not initialized")
	}

	return true
}

func reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID string) {
	targetNf, ok := cache.Instance().GetTargetNf(requesterNfType, targetNfType)
	if !ok {
		log.Errorf("Target NFProfile does not exist in cache")
		return
	}
	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Infof("NRF Agent role is leader, will fetch nfProfile %s to NRF", nfInstanceID)
		handleDiscoveryRequest(&targetNf, nfInstanceID)
	}
}

func handleDiscoveryRequest(targetNf *structs.TargetNf, nfInstanceID string) {
	requestOptions := ""
	if nfInstanceID != "" {
		requestOptions = consts.SearchDataTargetInstID + "=" + nfInstanceID
	}
	query := getDiscoveryRequestURL(targetNf, requestOptions)
	if query == "" {
		log.Errorf("handleDiscoveryRequest: failed to get Discovery request URL")
		return
	}

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Debugf("handleDiscoveryRequest: Master Agent KeepCache Mode status %t", client.IsKeepCacheMode())
		if client.IsKeepCacheMode() {
			log.Info("handleDiscoveryRequest: Master Agent is keep cache mode, not send message to NRF.")
			return
		}
		if _, err := handleDiscoveryRequestToNrf(targetNf, query); err != nil {
			log.Errorf("handleDiscoveryRequest: failed to query %s,%+v from NRF, %s",
				targetNf.TargetNfType, targetNf.TargetServiceNames, err.Error())
		}
	} else {
		if _, err := handleDiscoveryRequestToMaster(targetNf, query); err != nil {
			log.Errorf("handleDiscoveryRequest: failed to query %s,%+v from master agent, %s",
				targetNf.TargetNfType, targetNf.TargetServiceNames, err.Error())
		}
	}
}

func handleDiscoveryRequestToNrf(targetNf *structs.TargetNf, query string) (*httpclient.HttpRespData, error) {
	pm.Inc(consts.NrfDiscoveryRequestsTotal)

	log.Infof("Try to fetch %s %+v profile from NRF Discovery", targetNf.TargetNfType, targetNf.TargetServiceNames)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfDisc("h2", "GET", query, hdr, nil)
	if err != nil {
		log.Errorf("failed to send discovery request to NRF, %s", err.Error())
		return nil, err
	}

	err = handleDiscoveryResponse(targetNf, resp)
	if err != nil {
		fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return nil, err
	}
	fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	pushMessageToMSB(targetNf.RequesterNfType, targetNf.TargetNfType, "", consts.NFEventDiscResult, resp.Body)
	return resp, nil
}

func handleDiscoveryRequestToMaster(targetNf *structs.TargetNf, query string) (*httpclient.HttpRespData, error) {
	log.Infof("handleDiscoveryRequestToMaster: Try to fetch %s %+v profile from master NRF Discovery Agent",
		targetNf.TargetNfType, targetNf.TargetServiceNames)

	masterURL := getLeaderDiscURL()
	if masterURL == "" {
		log.Errorf("handleDiscoveryRequestToMaster: failed to get master URL")
		return nil, errors.New("NRF Discovery master agent URI not found")
	}
	log.Debugf("handleDiscoveryRequestToMaster: master discovery agent url \"%s\"", masterURL+query)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDo("h2", "GET", masterURL+query, hdr, nil)
	if err != nil {
		log.Errorf("handleDiscoveryRequestToMaster: failed to send discovery request to master discovery agent, %s", err.Error())
		return nil, err
	}

	err = handleDiscoveryResponse(targetNf, resp)
	if err != nil {
		fmRaiseNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
		return resp, err
	}
	fmClearNoAvailableDestination(targetNf.RequesterNfType, targetNf.TargetNfType)
	return resp, nil
}
*/

/*
func close(req *http.Request) {
	err := req.Body.Close()
	if err != nil {
		log.Error("close http request Body failure")
	}
}
*/

///////////////latest//////////

var (
	//Compile to compile partern into memory
	Compile map[string]*regexp.Regexp
)

//PreComplieRegexp to compile pattern into memory
func PreComplieRegexp() {
	Compile = make(map[string]*regexp.Regexp)
	re, _ := regexp.Compile("^[A-Fa-f0-9]*$")
	Compile[consts.SearchDataSupportedFeatures] = re

	re1, _ := regexp.Compile("^((http|https)://).*$")
	Compile[consts.SearchDataHnrfURI] = re1

	re2, _ := regexp.Compile("^[A-Fa-f0-9]{8}-[0-9]{3}-[0-9]{2,3}-([A-Fa-f0-9][A-Fa-f0-9]){1,10}$")
	Compile[consts.SearchDataExterGroupID] = re2

	re3, _ := regexp.Compile("^(SUBSCRIPTION|POLICY|EXPOSURE|APPLICATION)$")
	Compile[consts.SearchDataDataSet] = re3

	re4, _ := regexp.Compile("^[0-9]{3}$")
	Compile[consts.SearchDataMcc] = re4

	re5, _ := regexp.Compile("^[0-9]{2,3}$")
	Compile[consts.SearchDataMnc] = re5

	re6, _ := regexp.Compile("^[A-Fa-f0-9]{6}$")
	Compile[consts.SearchDataAmfID] = re6

	re7, _ := regexp.Compile("(^[A-Fa-f0-9]{4}$)|(^[A-Fa-f0-9]{6}$)")
	Compile[consts.SearchDataTac] = re7

	re8, _ := regexp.Compile("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$")
	Compile[consts.SearchDataUEIPv4Addr] = re8

	re9, _ := regexp.Compile("^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$")
	Compile[consts.SearchDataGpsi] = re9

	re10, _ := regexp.Compile("^(imsi-[0-9]{5,15}|nai-.+|.+)$")
	Compile[consts.SearchDataSupi] = re10

	re11, _ := regexp.Compile("[0-9]{5,15}")
	Compile[consts.GpsiRanges] = re11

	re12, _ := regexp.Compile("[0-9]{5,15}")
	Compile[consts.SupiRanges] = re12

	re13, _ := regexp.Compile("imsi-[0-9]{5,15}|suci-[0-9]{5,15}")
	Compile[consts.SupiFormat] = re13

	re14, _ := regexp.Compile("^[A-Fa-f0-9]{6}$")
	Compile[consts.SearchDataSnssaiSd] = re14

	re15, _ := regexp.Compile("^[0-9]{1,4}$")
	Compile[consts.SearchDataRoutingIndic] = re15
}

func GetValidityPeriod(body []byte) (int64, error) {
	if body == nil {
		return 0, fmt.Errorf("%s", "NRF response body is nil")
	}

	validityPeriod, err := jsonparser.GetInt(body, "validityPeriod")
	if err != nil {
		return 0, fmt.Errorf("Get validityPeriod from searchResult failed, err=%s", err.Error())
	}
	if validityPeriod < 0 {
		return validityPeriod, fmt.Errorf("%s", "NRF response validityPeriod < 0")
	}

	return validityPeriod, nil
}

func GetNfInstances(body []byte) ([][]byte, error) {
	if body == nil {
		return nil, fmt.Errorf("%s", "NRF response body is nil")
	}

	nfInstances := make([][]byte, 0)
	_, err := jsonparser.ArrayEach(body, func(nfInstance []byte, dataType jsonparser.ValueType, offset int, err error) {
		nfInstances = append(nfInstances, nfInstance)
	}, "nfInstances")
	if err != nil {
		log.Errorf("Get nfInstances from searchResult failed, err=%s", err.Error())
		return nil, err
	}

	return nfInstances, nil
}

func GetNfInstanceID(profile []byte) string {
	nfInstanceID, err := jsonparser.GetString(profile, "nfInstanceId")
	if err != nil {
		log.Errorf("The profile get nfInstanceID fail, err:%s", err.Error())
		return ""
	}

	return nfInstanceID
}

func GetLeaderDiscURL() string {
	leader := election.GetLeader()
	port := cm.Opts.PortHTTP2WithoutTLS
	if leader == "" {
		return ""
	}

	return "http://" + leader + ":" + strconv.Itoa(port) + "/nrf-discovery-agent/v1/"
}

func GetDiscoveryRequestURL(targetNf *structs.TargetNf, opts string, requesterFqdn string) string {
	requesterNfType := targetNf.RequesterNfType
	targetNfType := targetNf.TargetNfType
	targetNfServices := ""
	for _, service := range targetNf.TargetServiceNames {
		targetNfServices += consts.SearchDataServiceName + "=" + service + "&"
	}

	queryURL := "nf-instances?" +
		targetNfServices +
		consts.SearchDataTargetNfType + "=" + targetNfType + "&" +
		consts.SearchDataRequesterNfType + "=" + requesterNfType
	/*
		requesterFqdn, ok := cache.Instance().GetRequesterFqdn(requesterNfType)
		if !ok {
			log.Errorf("NfType:%s was not registered to Register Agent", requesterNfType)
			return ""
		}
	*/
	if requesterFqdn != "" {
		queryURL += "&" + consts.SearchDataRequesterNFInstFQDN + "=" + requesterFqdn
	}
	if opts != "" {
		queryURL += "&" + opts
	}
	log.Debugf("NRF Discovery Request URL: %s", queryURL)

	return queryURL
}

func GetStatusNotifURLs() string {
	statusNotifIPEndpoint, ok := structs.GetStatusNotifIPEndPoint()
	if !ok {
		log.Errorf("Get statusNotifIPEndpoint URL failure")
		return ""
	}
	ipAddress := statusNotifIPEndpoint.Ipv4Address
	if cm.IsEnableIpv6() {
		ipAddress = statusNotifIPEndpoint.Ipv6Address
	}
	nrfURLs := "http://" + ipAddress + ":" + strconv.Itoa(statusNotifIPEndpoint.Port)
	return nrfURLs
}

func BuildSubscriptionPostData(oneSubsData *structs.OneSubscriptionData, requesterNfFqdn string) []byte {
	requesterNfType := oneSubsData.RequesterNfType
	targetNfType := oneSubsData.TargetNfType
	targetServiceName := oneSubsData.TargetServiceName
	if requesterNfType == "" || targetNfType == "" || targetServiceName == "" {
		log.Error("RequesterNfType,targetNfType,targetServiceName in nfTargetProfile is mandatary")
		return nil
	}
	/*
		requesterNfFqdn, exists := cache.Instance().GetRequesterFqdn(requesterNfType)
		if !exists {
			log.Warnf("Disc agent don't receive %s FQDN from Register Agent", requesterNfType)
			return nil
		}
	*/
	callbackCluster := GetStatusNotifURLs()
	if callbackCluster == "" {
		log.Error("CallbackCluster can not be empty")
		return nil
	}

	subscribeData := &structs.SubscriptionData{}

	subscribeData.NfStatusNotificationURI = callbackCluster + "/nrf-notify-agent/v1/notify/" + requesterNfType + "/" + targetNfType
	subscribeData.ReqNfType = requesterNfType
	subscribeData.ReqNfFqdn = requesterNfFqdn
	serviceNameCond := structs.ServiceNameCond{
		ServiceName: targetServiceName,
	}
	subscribeData.SubscrCond = serviceNameCond
	subscribeData.NotifCondition = oneSubsData.NotifCondition

	log.Debugf("Build POST subscriptionData:%+v", subscribeData)

	subscribeDataRaw, err := json.Marshal(subscribeData)
	if err != nil {
		log.Errorf("Marshal subscribeData fail, err:%s", err.Error())
		return nil
	}

	return subscribeDataRaw
}

func BuildSubscriptionPostRoamData(oneSubsData *structs.OneSubscriptionData, requesterNfFqdn string, plmnID *structs.PlmnID, validityTime string) []byte {
	requesterNfType := oneSubsData.RequesterNfType
	targetNfType := oneSubsData.TargetNfType
	nfInstanceID := oneSubsData.NfInstanceID
	if requesterNfType == "" || targetNfType == "" || nfInstanceID == "" {
		log.Error("RequesterNfType,targetNfType,nfInstanceID in nfTargetProfile is mandatary")
		return nil
	}

	callbackCluster := GetStatusNotifURLs()
	if callbackCluster == "" {
		log.Error("CallbackCluster can not be empty")
		return nil
	}

	subscribeData := &structs.SubscriptionData{}

	subscribeData.NfStatusNotificationURI = callbackCluster + "/nrf-notify-agent/v1/notify/" + requesterNfType + consts.RoamSuffix + "/" + targetNfType
	subscribeData.ReqNfType = requesterNfType
	subscribeData.ReqNfFqdn = requesterNfFqdn
	subscribeData.ValidityTime = validityTime
	nfInstanceIdCond := structs.NfInstanceIDCond{
		NfInstanceID: nfInstanceID,
	}
	subscribeData.SubscrCond = nfInstanceIdCond
	subscribeData.PlmnID = plmnID
	subscribeData.NotifCondition = oneSubsData.NotifCondition

	log.Debugf("Build POST roaming subscriptionData:%+v", subscribeData)

	subscribeDataRaw, err := json.Marshal(subscribeData)
	if err != nil {
		log.Errorf("Marshal subscribeData fail, err:%s", err.Error())
		return nil
	}

	return subscribeDataRaw
}

func BuildSubscriptionPatchData(validityTime int) []byte {
	now := time.Now()
	prolongValue := time.Duration(validityTime) * time.Second
	nextTimeStamp := now.Add(prolongValue)
	nextTimeStampStr := nextTimeStamp.Format(time.RFC3339)

	item := structs.PatchItem{
		Op:    "replace",
		Path:  "/validityTime",
		Value: nextTimeStampStr,
	}

	items := make([]structs.PatchItem, 0)
	items = append(items, item)

	log.Debugf("Build PATCH subscriptionData:%+v", items)

	patchData, err := json.Marshal(items)
	if err != nil {
		log.Errorf("marsh patch items failure, err:%s", err.Error())
		return nil
	}

	return patchData
}

var PushMessageToMSB = func(requesterNfType, targetNfType, nfInstanceID string, event string, resp []byte) bool {
	/*
		if !election.IsActiveLeader("3201", consts.DiscoveryAgentReadinessProbe) {
			return true
		}
	*/
	var msgBusDisc structs.NotificationMsg
	msgBusDisc.AgentProducerID = common.GetSelfUUID()
	msgBusDisc.NfEvent = event
	msgBusDisc.NfType = targetNfType
	msgBusDisc.NfInstanceID = nfInstanceID

	if resp != nil {
		err := json.Unmarshal(resp, &msgBusDisc.MessageBody)
		if err != nil {
			log.Errorf("Decode MessageBody message Unmarshal fail, %s", err.Error())
			return false
		}
	}
	log.Infof("RequesterNfType %s, msgBusDisc %+v", requesterNfType, msgBusDisc)

	jsonBuf, err := json.Marshal(msgBusDisc)
	if err != nil {
		log.Errorf("Failed to Marshal Disc message, %s", err.Error())
		return false
	}
	log.Infof("Dispatch message body : %s", string(jsonBuf))

	discMsgbus := common.GetDiscMsgbus()
	if discMsgbus == nil {
		log.Warnf("Message bus was not initialized")
		return false
	}

	topic := consts.MsgbusTopicNamePrefix + strings.ToLower(requesterNfType)
	err = discMsgbus.SendMessage(topic, string(jsonBuf))
	if err != nil {
		log.Errorf("Failed to send notification to message bus, %s", err.Error())
		return false
	}

	log.Debugf("Push message to MSB succeed")

	return true
}

/*
func unscribeByNfType() {
	if worker.IsKeepCacheMode() {
		log.Info("keep cache mode, not send message to NRF.")
		return
	}
	targetNfs, ok := cacheManager.GetTargetNfs(requesterNfType)
	if !ok {
		log.Errorf("Failed to get targetNfProfiles for nfType[%s], please check configmap status", requesterNfType)
		return
	}

	for _, targetNf := range targetNfs {
		unsubscribeSubscription(requesterNfType, targetNf.TargetNfType)
	}
}

func unSubscribeByNfinstanceID() {

}

func UnsubscribeByNfInstanceID(requesterNfType, targetNfType, nfInstanceID string) {
	subscriptionID, ok := cacheManager.GetNfProfileSubscriptionID(requesterNfType, targetNfType, nfInstanceID)
	if !ok {
		log.Warnf("No such subscription for nfProfile:%s, will skip do unsubscribe by nfInstanceID", nfInstanceID)
		return
	}
	log.Infof("nfProfile:%s subscriptionID:%s", nfInstanceID, subscriptionID)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
		if err != nil {
			log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
		} else {
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			} else {
				log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			}
		}
	} else {
		log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
	}

	cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
	cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
}

func unsubscribeSubscription(requesterNfType, targetNfType string) {
	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	ok, subscriptionIDURLs := cache.Instance().GetSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(subscriptionIDURLs) > 0 {
		subscriptionID := subscriptionIDURLs[0]
		log.Infof("subscriptionID of %s: %+v", requesterNfType, subscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
		cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
		cacheManager.UpdateSubscriptionStorage()

		ok, subscriptionIDURLs = cacheManager.GetSubscriptionIDs(requesterNfType, targetNfType)
	}

	ok, roamsubscriptionIDURLs := cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(roamsubscriptionIDURLs) > 0 {
		roamsubscriptionID := roamsubscriptionIDURLs[0]
		log.Infof("roaming subscriptionID of %s: %+v", requesterNfType, roamsubscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", roamsubscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE roaming subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelRoamingSubscriptionInfo(requesterNfType, targetNfType, roamsubscriptionID)
		cacheManager.DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, roamsubscriptionID)

		ok, roamsubscriptionIDURLs = cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	}
}
*/
