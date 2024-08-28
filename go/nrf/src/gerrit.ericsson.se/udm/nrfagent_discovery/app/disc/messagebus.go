package disc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/buger/jsonparser"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/discutil"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc/schema"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"

	jsonpatch "github.com/evanphx/json-patch"
)

var (
	requesterNfList = make(map[string]bool)
	discAgentMutex  sync.Mutex
)

//InitMessageBus init messagebus
func initMessageBus() error {
	discMsgbus := msgbus.NewMessageBus(os.Getenv("MESSAGE_BUS_KAFKA"))
	if discMsgbus == nil {
		log.Error("Initialize message bus failure")
		return errors.New("Initialize message bus failure")
	}

	log.Infof("initialize message bus success")
	common.SetDiscMsgbus(discMsgbus)

	//nrf-agent-nfType
	CreateNfTypeTopic()

	var ok bool
	ok = CreateRegDiscInnerTopic()
	if ok {
		log.Infof("initMessageBus: create topic nrf-agent-regdiscinner success")
	} else {
		log.Error("initMessageBus: create topic nrf-agent-regdiscinner failure")
		return errors.New("Create topic nrf-agent-regdiscinner failure")
	}

	ok = CreateNtfDiscInnerTopic()
	if ok {
		log.Infof("initMessageBus: create topic nrf-agent-ntfdiscinner success")
	} else {
		log.Error("initMessageBus: create topic nrf-agent-ntfdiscinner failure")
		return errors.New("Create topic nrf-agent-ntfdiscinner failure")
	}

	ok = CreateDiscDiscInnerTopic()
	if ok {
		log.Infof("initMessageBus: create topic nrf-agent-discdiscinner success")
	} else {
		log.Error("initMessageBus: create topic nrf-agent-discdiscinner failure")
		return errors.New("Create topic nrf-agent-discdiscinner failure")
	}

	return nil
}

//CreateNfTypeTopic create nfType topic
func CreateNfTypeTopic() bool {
	discAgentMutex.Lock()
	defer discAgentMutex.Unlock()

	targetNfProfile := structs.GetTargetNfProfiles()
	if targetNfProfile == nil {
		return false
	}

	for _, nfProfile := range targetNfProfile {
		if nfProfile.RequesterNfType == "" {
			continue
		}
		_, existed := requesterNfList[strings.ToLower(nfProfile.RequesterNfType)]
		if !existed {
			requesterNfList[strings.ToLower(nfProfile.RequesterNfType)] = false
		}
	}

	if discMsgbus := common.GetDiscMsgbus(); discMsgbus != nil {
		for nfType, consumed := range requesterNfList {
			if nfType != "" && !consumed {
				topic := consts.MsgbusTopicNamePrefix + strings.ToLower(nfType)
				err := discMsgbus.ConsumeTopicPlusTopicName(topic, notificationMessageHandler)
				if err != nil {
					log.Errorf("Fail to create topic %s for %s Service", topic, nfType)
					//continue
					return false
				} else {
					requesterNfList[nfType] = true
					log.Infof("Success to create topic %s for %s Service", topic, nfType)
				}
			}
		}
	}

	return true
}

//CreateRegDiscInnerTopic create inner topic between reg and disc
func CreateRegDiscInnerTopic() bool {
	discMsgbus := common.GetDiscMsgbus()
	if discMsgbus == nil {
		log.Errorf("Failed to get messageBus for NRF Discovery Agent")
		return false
	}

	regDiscInnerTopicName := consts.MsgbusTopicNamePrefix + consts.RegDiscInner
	err := discMsgbus.ConsumeTopicPlusTopicName(regDiscInnerTopicName, regDiscInnerMessageHandler)
	if err != nil {
		log.Errorf("Failed to ConsumeTopicPlusTopicName %s for NRF Discovery Agent", regDiscInnerTopicName)
		return false
	}

	return true
}

//CreateNtfDiscInnerTopic create inner topic between ntf and disc
func CreateNtfDiscInnerTopic() bool {
	discMsgbus := common.GetDiscMsgbus()
	if discMsgbus == nil {
		log.Errorf("failed to get messageBus for NRF Discovery Agent")
		return false
	}

	ntfDiscInnerTopicName := consts.MsgbusTopicNamePrefix + consts.NtfDiscInner
	err := discMsgbus.ConsumeTopicPlusTopicName(ntfDiscInnerTopicName, ntfDiscInnerMessageHandler)
	if err != nil {
		log.Errorf("failed to ConsumeTopicPlusTopicName %s for NRF Discovery Agent", ntfDiscInnerTopicName)
		return false
	}

	return true
}

//CreateDiscDiscInnerTopic create inner topic between master disc and slave disc
func CreateDiscDiscInnerTopic() bool {
	discMsgbus := common.GetDiscMsgbus()
	if discMsgbus == nil {
		log.Errorf("Failed to get messageBus for NRF Discovery Agent")
		return false
	}

	discDiscInnerTopicName := consts.MsgbusTopicNamePrefix + consts.DiscDiscInner
	err := discMsgbus.ConsumeTopicPlusTopicName(discDiscInnerTopicName, discDiscInnerMessageHandler)
	if err != nil {
		log.Errorf("Failed to ConsumeTopicPlusTopicName %s for NRF Discovery Agent", discDiscInnerTopicName)
		return false
	}

	return true
}

////////////notification message handler////////////
func notificationMessageHandler(topic string, msg []byte) {
	var messageData structs.NotificationMsg
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		log.Errorf("notificationMessageHandler: Unmarshal fail, err:%s", err.Error())
		return
	}

	if common.GetSelfUUID() == messageData.AgentProducerID {
		log.Infof(":notificationMessageHandler: ignore the message send by agent self")
		return
	}

	nrfNtfAgent := "NRF-Notify-Agent"
	log.Infof("Message from messagebus consumer:%s,topic:%s,event:%s,content:%s", nrfNtfAgent, topic, messageData.NfEvent, string(msg))

	requesterNfType := strings.ToUpper(strings.Replace(topic, consts.MsgbusTopicNamePrefix, "", -1))

	_, exists := cache.Instance().GetRequesterFqdn(requesterNfType)
	if !exists {
		log.Warnf("RequesterNftype %s is not registered, ignore the related notify event", requesterNfType)
		return
	}

	switch messageData.NfEvent {
	case consts.NFRegister:
		ntfRegisterEventHandler(requesterNfType, &messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.NfEvent, topic)
	case consts.NFProfileChg:
		ntfChangeEventHandler(requesterNfType, &messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.NfEvent, topic)
	case consts.NFDeRegister:
		ntfDeregisterEventHandler(requesterNfType, &messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.NfEvent, topic)
	case consts.NFEventDiscResult:
		ntfDiscResultEventHandler(requesterNfType, &messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.NfEvent, topic)
	default:
		log.Errorf("Message event(%s) unknown", messageData.NfEvent)
	}
}

func ntfRegisterEventHandler(requesterNfType string, messageData *structs.NotificationMsg) {
	nfInstanceID := messageData.NfInstanceID
	targetNfType := messageData.NfType

	exist := cache.Instance().Probe(requesterNfType, targetNfType, nfInstanceID)
	if exist {
		log.Infof("NF Instance(%s) already exist in Cache", nfInstanceID)
		return
	}

	searchResult := messageData.MessageBody
	if searchResult == nil ||
		len(searchResult.NfInstances) != 1 {
		log.Error("The message body should contain one and at most one NF profile")
		return
	}
	if searchResult.NfInstances[0].NfInstanceID != nfInstanceID {
		log.Infof("NfInstanceID %s is not the same with it in NF profile", nfInstanceID)
		return
	}
	nfinstance, err := json.Marshal(searchResult.NfInstances[0])
	if err != nil {
		log.Errorf("Failed to Marshal nfInstances, %s", err.Error())
		return
	}

	cache.Instance().CachedWithTTL(requesterNfType, targetNfType, nfinstance, uint(cm.GetDefaultValidityPeriod()), false)
	cache.Instance().SetCacheStatus(requesterNfType, targetNfType, true)

	fmClearNoAvailableDestination(requesterNfType, targetNfType)
}

func ntfChangeEventHandler(requesterNfType string, messageData *structs.NotificationMsg) {
	nfInstanceID := messageData.NfInstanceID
	targetNfType := messageData.NfType

	if targetNfType == "" {
		log.Error("the message body do not contain targetNfType")
		return
	}
	searchResult := messageData.MessageBody
	var isRoam bool
	if searchResult == nil || len(searchResult.NfInstances) != 1 {
		//Error NFProfileChg message, not treat it as roaming
		isRoam = false
	} else {
		isRoam = discutil.IsRoamMessage(requesterNfType, searchResult.NfInstances)
	}
	if isRoam {
		exist := cacheManager.ProbeRoam(requesterNfType, targetNfType, nfInstanceID)
		if !exist {
			log.Infof("NF Instance(%s) do not exist in Roaming cache Cache", nfInstanceID)
			return
		}
	} else {
		exist := cacheManager.Probe(requesterNfType, targetNfType, nfInstanceID)
		if !exist {
			log.Infof("NF Instance(%s) do not exist in Cache", nfInstanceID)
			targetNf, ok := cacheManager.GetTargetNf(requesterNfType, targetNfType)
			if !ok {
				log.Errorf("GetTargetNfs failed")
				return
			}
			if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
				handleDiscoveryRequest(&targetNf, nfInstanceID)
			}
			return
		}
	}

	if searchResult == nil ||
		len(searchResult.NfInstances) != 1 {
		log.Error("The message body should contain one and at most one NF profile")
		return
	}
	if searchResult.NfInstances[0].NfInstanceID != nfInstanceID {
		log.Infof("nfInstanceID %s is not the same with it in NF profile", nfInstanceID)
		return
	}

	nfinstance, err := json.Marshal(searchResult.NfInstances[0])
	if err != nil {
		log.Errorf("Failed to Marshal nfInstances, %s", err.Error())
		return
	}

	cacheManager.ReCached(requesterNfType, targetNfType, nfInstanceID, nfinstance, isRoam)
}

func ntfDeregisterEventHandler(requesterNfType string, messageData *structs.NotificationMsg) {
	nfInstanceID := messageData.NfInstanceID
	targetNfType := messageData.NfType

	if targetNfType == "" {
		log.Errorf("messageData.NfType is empty")
		return
	}

	exist, isRoam := cache.Instance().ProbeAllCache(requesterNfType, targetNfType, nfInstanceID)
	if !exist {
		log.Infof("NF Instance(%s) do not exist in Cache", nfInstanceID)
		return
	}

	cacheManager.DeCached(requesterNfType, targetNfType, nfInstanceID, isRoam)

	fmRaiseNoAvailableDestination(requesterNfType, targetNfType)
}

func ntfDiscResultEventHandler(requesterNfType string, messageData *structs.NotificationMsg) {
	targetNfType := messageData.NfType
	results, err := json.Marshal(messageData.MessageBody)
	if err != nil {
		log.Errorf("Failed to marshal message, %s", err.Error())
		return
	}

	nfInstances, validityPeriod, ok := cache.SpliteSeachResult(results)
	if !ok {
		log.Errorf("SpliteSeachResult error")
		return
	}

	isRoam := discutil.IsRoamMessage(requesterNfType, messageData.MessageBody.NfInstances)
	for _, nfProfile := range nfInstances {
		cache.Instance().CachedWithTTL(requesterNfType, targetNfType, nfProfile, validityPeriod, isRoam)
	}

	fmClearNoAvailableDestination(requesterNfType, targetNfType)
}

////////////reg disc inner message handler////////////
func regDiscInnerMessageHandler(topic string, msg []byte) {
	log.Infof("ENTRY FROM %s: Message comes from MESSAGEBUS %s, %+v", "NRF Register Agent", topic, string(msg))

	var messageData structs.RegDiscInnerMsg
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		log.Errorf("regDiscInnerMessageHandler: unmarshal failure, Error:%s", err.Error())
		return
	}

	switch messageData.EventType {
	case consts.EventTypeRegister:
		regRegisterEventHandler(&messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.EventType, topic)
	case consts.EventTypeDeregister:
		regDeregisterEventHandler(&messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.EventType, topic)
	case consts.EventTypeFQDNChanged:
		regFqdnChangeEventHandler(&messageData)
		log.Infof("Handle event:%s from topic:%s finished", messageData.EventType, topic)
	default:
		log.Errorf("Message event(%s) unknown", messageData.EventType)
	}
}

func regRegisterEventHandler(messageData *structs.RegDiscInnerMsg) bool {
	requesterNfType := messageData.NfType
	nfInstanceID := messageData.NfInstanceID
	requesterNfFqdn := messageData.FQDN
	requsterPlmns := messageData.Plmns

	_, exists := cache.Instance().GetRequesterFqdn(requesterNfType)
	if exists {
		log.Warnf("RequesterNftype:%s, Instance:%s %s was already registered, and the fqdn:", requesterNfType, nfInstanceID, requesterNfFqdn)
		return true
	}

	cache.Instance().SetRequesterFqdn(requesterNfType, requesterNfFqdn)
	if len(messageData.Plmns) > 0 {
		cache.Instance().SetRequesterPlmns(requesterNfType, requsterPlmns)
	}
	workerManager.PrepareNfRegister(requesterNfType)

	return true
}

func regDeregisterEventHandler(messageData *structs.RegDiscInnerMsg) bool {
	requesterNfType := messageData.NfType

	_, exists := cache.Instance().GetRequesterFqdn(requesterNfType)
	if !exists {
		log.Warnf("nfSubscribeDeregisterEventHandler: %s was not registered to Register Agent", requesterNfType)
		return true
	}

	cache.Instance().DeleteRequesterFqdn(requesterNfType)
	cache.Instance().DeCachedByNfType(requesterNfType)
	cache.Instance().DelRequesterPlmns(requesterNfType)

	discutil.UnsubscribeNfDeregister(requesterNfType)

	fmClearNoAvailableDestination(requesterNfType, "")

	return true
}

func regFqdnChangeEventHandler(messageData *structs.RegDiscInnerMsg) {
	changedFqdn := messageData.FQDN
	requesterNfType := messageData.NfType

	if requesterNfType == "" {
		log.Errorf("RegDiscInnerMsg event FQDN Changed with wrong requesterNfType(%s)", requesterNfType)
		return
	}

	cacheManager.SetRequesterFqdn(requesterNfType, changedFqdn)

	//clean cache and need refetch nfprofile
	cacheManager.Flush(requesterNfType)
	cacheManager.FlushRoam(requesterNfType)
	log.Debugf("Clean cache by nfType and will fetch NF profile by new fqdn")

	targetNfs, ok := cacheManager.GetTargetNfs(requesterNfType)
	if !ok {
		log.Errorf("%s targetNf in cache is nil, So can not fetch profile without target-nf-type and service-names", requesterNfType)
		return
	}
	for _, targetNf := range targetNfs {
		handleDiscoveryRequest(&targetNf, "")
		log.Debugf("Finish handleDiscoveryRequestToNrf according to new fqdn %s", changedFqdn)
	}
}

func ntfDiscInnerMessageHandler(topic string, msg []byte) {
	log.Infof("Message from:%s, topic:%s, message:%+v", "NRF-Notify-Agent", topic, string(msg))
	var messageData structs.NtfDiscInnerMsg
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		log.Errorf("Unmarshal NtfDiscInnerMsg fail, eror:%s", err.Error())
		return
	}
	if common.GetSelfUUID() == messageData.AgentProducerID {
		log.Infof("Ignore the message send by self")
		return
	}

	//requesterNfType := strings.ToUpper(strings.Replace(topic, consts.MsgbusTopicNamePrefix, "", -1))
	requesterNfType := strings.ToUpper(messageData.ReqNfType)
	targetNfType := strings.ToUpper(messageData.NfType)
	if requesterNfType == "" || targetNfType == "" {
		log.Errorf("RequesterNfType and targetNfType can not be empty, requesterNfType:%s, targetNfType:%s", requesterNfType, targetNfType)
		return
	}
	nfInstanceID := messageData.NfInstanceID

	log.Infof("RequesterNfType is %s, targetNfType is %s", requesterNfType, targetNfType)

	if !election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Debugf("Message bus nrf-agent-ntfdiscinner leader will process notification message")
		return
	}
	if common.IsRoamNotifcation(messageData.ReqNfType) {
		roamNtfMessageHandler(&messageData)
	} else {
		if messageData.NfEvent == consts.NotifEventProfileChg {

			exist := cache.Instance().Probe(requesterNfType, targetNfType, nfInstanceID)
			if !exist {
				log.Errorf("NF Instance(%s) do not exist in Cache", nfInstanceID)
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			nfProfileInCache := cache.Instance().FetchNfProfile(requesterNfType, nfInstanceID)
			if nfProfileInCache == nil {
				log.Errorf("%s nfProfile in cache is empty, cache Probe failed", nfInstanceID)
				return
			}
			log.Debugf("Fetch cache nfProfile by requesterNfType and nfInstanceID successfully")

			var PatchApplyDataArray = make([]structs.NfProfilePatchApplyData, 0)
			for _, v := range messageData.MessageBody {
				PatchApplyData := structs.NfProfilePatchApplyData{}
				PatchApplyData.Op = strings.ToLower(v.Op)
				PatchApplyData.Path = v.Path
				PatchApplyData.From = v.From
				PatchApplyData.Value = v.NewValue
				PatchApplyDataArray = append(PatchApplyDataArray, PatchApplyData)
			}
			log.Debugf("PatchApplyDataArray are %+v", PatchApplyDataArray)
			PatchApplyDataArraySlice, err := json.Marshal(PatchApplyDataArray)
			if err != nil {
				log.Errorf("Json Marshal PatchApplyDataArray failed")
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			//use patch schema to validate patch body,
			err = schema.ValidatePatchDocument(string(PatchApplyDataArraySlice[:]))
			if err != nil {
				log.Errorf("Validate PatchApplyDataArraySlice failed, error is %s", err.Error())
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			updatedNfProfile, ok := applyPatchItems(nfProfileInCache, PatchApplyDataArraySlice)
			if ok != true {
				log.Errorf("Apply Patch Items failed")
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			sNFProfileByte, err := convertToSearchResultNFProfile(updatedNfProfile)
			if err != nil {
				log.Errorf("Convert updated nfProfile to search result Profile failed, error is %s", err.Error())
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			//validate the nfProfile after Patching using schema, if failed, will trigger refetch nfProfile to NRF
			err = schema.ValidateNfProfile(string(sNFProfileByte[:]))
			if err != nil {
				log.Errorf("Validate updated nfProfile failed, error is %s", err.Error())
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}

			cache.Instance().ReCached(requesterNfType, targetNfType, nfInstanceID, sNFProfileByte, false)
			log.Debugf("New %s nfProfile is: %+v", nfInstanceID, string(sNFProfileByte))

			//convert nfProfile to seatchResult structure
			searchResultBody, err := convertToSearchResultBody(sNFProfileByte)
			if err != nil {
				log.Errorf("convertToSearchResultBody failed, error is %s", err.Error())
				log.Infof("NRF Agent will fetch nfProfile %s to NRF", nfInstanceID)
				reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID)
				return
			}
			util.PushMessageToMSB(requesterNfType, targetNfType, nfInstanceID, consts.NFProfileChg, searchResultBody)
			log.Info("NRF Discovery Agent apply patch successfully")
		} else {
			log.Info("Message bus nrf-agent-ntfdiscinner only receive patch message and leader will process it")
		}
	}
}

func roamNtfMessageHandler(messageData *structs.NtfDiscInnerMsg) {
	if messageData == nil {
		return
	}

	//remove roaming flag and get the real request nfType
	reqNfType := common.GetReqNfTypeForRoam(messageData.ReqNfType)
	requesterNfType := strings.ToUpper(reqNfType)
	targetNfType := strings.ToUpper(messageData.NfType)

	nfInstanceID := messageData.NfInstanceID

	nfProfileInCache := cache.Instance().GetRoamingNfProfile(requesterNfType, targetNfType, nfInstanceID)
	if nfProfileInCache == nil {
		log.Warningf("%s nfProfile in cache is empty, ignore roaming notification event", nfInstanceID)
		return
	}
	log.Debugf("Fetch cache nfProfile by requesterNfType and nfInstanceID successfully")
	if messageData.NfEvent == consts.NotifEventProfileChg {
		if len(messageData.MessageBody) > 0 {
			roamPatchNtfMessageHandler(messageData, requesterNfType, targetNfType, nfInstanceID, &nfProfileInCache)
		} else if messageData.NfProfile != nil {
			roamUpdateNtfMessageHandler(messageData, requesterNfType, targetNfType, nfInstanceID, &nfProfileInCache)
		}
	} else if messageData.NfEvent == consts.NotifEventDeregister {
		if !cache.Instance().ProbeRoam(requesterNfType, targetNfType, messageData.NfInstanceID) {
			log.Warningf("roamNtfMessageHandler instance(%s) not exist for nftype[%s,%s]", messageData.NfInstanceID, requesterNfType, targetNfType)
			return
		}
		discutil.UnsubscribeRoamNfProfile(requesterNfType, targetNfType, nfInstanceID)
		cache.Instance().DeCached(requesterNfType, targetNfType, messageData.NfInstanceID, true)
		messageData.ReqNfType = reqNfType
		body, err := json.Marshal(messageData)
		if err != nil {
			log.Errorf("roamNtfMessageHandler marshal messageData fail, %s", err.Error())
			return
		}
		util.PushMessageToMSB(requesterNfType, targetNfType, nfInstanceID, consts.NFDeRegister, body)
	} else {
		log.Info("Message bus nrf-agent-ntfdiscinner will not handler roaming NotifEventRegister event")
	}

}

func roamPatchNtfMessageHandler(messageData *structs.NtfDiscInnerMsg, requesterNfType string, targetNfType string, nfInstanceID string, nfProfileInCache *[]byte) {
	if messageData == nil || nfProfileInCache == nil {
		return
	}
	var PatchApplyDataArray = make([]structs.NfProfilePatchApplyData, 0)
	for _, v := range messageData.MessageBody {
		PatchApplyData := structs.NfProfilePatchApplyData{}
		PatchApplyData.Op = strings.ToLower(v.Op)
		PatchApplyData.Path = v.Path
		PatchApplyData.From = v.From
		PatchApplyData.Value = v.NewValue
		PatchApplyDataArray = append(PatchApplyDataArray, PatchApplyData)

	}
	log.Debugf("PatchApplyDataArray are %+v", PatchApplyDataArray)
	PatchApplyDataArraySlice, err := json.Marshal(PatchApplyDataArray)
	if err != nil {
		log.Errorf("Json Marshal PatchApplyDataArray failed")
		deregisterInstance(requesterNfType, targetNfType, nfInstanceID)
		return
	}

	//use patch schema to validate patch body,
	err = schema.ValidatePatchDocument(string(PatchApplyDataArraySlice[:]))
	if err != nil {
		log.Errorf("Validate PatchApplyDataArraySlice failed, error is %s", err.Error())
		deregisterInstance(requesterNfType, targetNfType, nfInstanceID)
		return
	}

	updatedNfProfile, ok := applyPatchItems(*nfProfileInCache, PatchApplyDataArraySlice)
	if ok != true {
		log.Errorf("Apply Patch Items failed")
		deregisterInstance(requesterNfType, targetNfType, nfInstanceID)
		return
	}

	sendProfieChgToMSB(updatedNfProfile, requesterNfType, targetNfType, nfInstanceID)
}

func roamUpdateNtfMessageHandler(messageData *structs.NtfDiscInnerMsg, requesterNfType string, targetNfType string, nfInstanceID string, nfProfileInCache *[]byte) {
	if messageData == nil || nfProfileInCache == nil {
		return
	}
	if len(messageData.NfProfile.PlmnList) == 0 {
		//restore plmnList from cached profile
		plmnListData, _, _, err := jsonparser.Get(*nfProfileInCache, "plmnList")
		if err == nil {
			plmnIDs := make([]structs.PlmnID, 0)
			err := json.Unmarshal(plmnListData, &plmnIDs)
			if err != nil {
				log.Errorf("Unmarshal plmnList fail, %s", err.Error())
				return
			}
			messageData.NfProfile.PlmnList = plmnIDs
		}
	}

	updatedNfProfile, err := json.Marshal(messageData.NfProfile)
	if err != nil {
		log.Errorf("Marshal updatedNfProfile fail, %s", err.Error())
		return
	}
	sendProfieChgToMSB(updatedNfProfile, requesterNfType, targetNfType, nfInstanceID)
}

func sendProfieChgToMSB(nfProfile []byte, requesterNfType string, targetNfType string, nfInstanceID string) {
	sNFProfileByte, err := convertToSearchResultNFProfile(nfProfile)
	if err != nil {
		log.Errorf("sendProfieChgToMSB convert updated nfProfile to search result Profile failed, error is %s", err.Error())
		return
	}

	//validate the nfProfile after Patching using schema, if failed, will trigger refetch nfProfile to NRF
	err = schema.ValidateNfProfile(string(sNFProfileByte[:]))
	if err != nil {
		log.Errorf("sendProfieChgToMSB validate updated nfProfile failed, error is %s", err.Error())
		return
	}

	cache.Instance().ReCached(requesterNfType, targetNfType, nfInstanceID, sNFProfileByte, true)
	log.Debugf("sendProfieChgToMSB new %s nfProfile is: %+v", nfInstanceID, string(sNFProfileByte))

	//convert nfProfile to seatchResult structure
	searchResultBody, err := convertToSearchResultBody(sNFProfileByte)
	if err != nil {
		log.Errorf("convertToSearchResultBody failed, error is %s", err.Error())
		return
	}
	util.PushMessageToMSB(requesterNfType, targetNfType, nfInstanceID, consts.NFProfileChg, searchResultBody)
	log.Info("NRF Discovery Agent send profile change successfully")
}

func deregisterInstance(requesterNfType string, targetNfType string, nfInstanceID string) {
	log.Info("NRF Discovery Agent send deregister instance")
	discutil.UnsubscribeRoamNfProfile(requesterNfType, targetNfType, nfInstanceID)
	cache.Instance().DeCached(requesterNfType, targetNfType, nfInstanceID, true)
	util.PushMessageToMSB(requesterNfType, targetNfType, nfInstanceID, consts.NFDeRegister, nil)
}

////////////master-slave disc inner message handler////////////
func discDiscInnerMessageHandler(topic string, msg []byte) {
	log.Infof("ENTRY FROM %s: Message comes from MESSAGEBUS %s, %+v", "NRF Register Agent", topic, string(msg))

	var messageData structs.DiscDiscInnerMsg
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		log.Errorf("discDiscInnerMessageHandler: unmarshal failure, Error:%s", err.Error())
		return
	}

	if common.GetSelfUUID() == messageData.AgentProducerID {
		log.Infof("Ignore the message send by self")
		return
	}

	switch messageData.EventType {
	case consts.EventTypeSyncSubscrInfo:
		syncSubscrInfoEventHandler(&messageData)
	default:
		log.Errorf("Message event(%s) unknown", messageData.EventType)
	}
}

func syncSubscrInfoEventHandler(messageData *structs.DiscDiscInnerMsg) bool {
	if messageData == nil {
		log.Warnf("MessageData from topic:%s is nil, skip the message", consts.MsgbusTopicNamePrefix+consts.DiscDiscInner)
		return true
	}

	subscrInfo := messageData.SubscrInfo
	var isRoamSubscription bool
	if subscrInfo.NfInstanceID != "" {
		isRoamSubscription = true
	} else {
		isRoamSubscription = false
	}

	if !isRoamSubscription {
		//below func have been abandoned
		ok := cacheManager.AddSubscriptionInfo(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo)
		if !ok {
			log.Errorf("add subscrInfo event hander AddSubscriptionInfo failure: %s, %s, %s", subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID)
			return false
		}
		cacheManager.DelSubscriptionMonitor(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID)
		cacheManager.SuperviseSubscription(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID, subscrInfo.ValidityTime)
		log.Debugf("Add %s subscrInfo event hander finished", subscrInfo.SubscriptionID)

		requsterNfType := subscrInfo.RequesterNfType
		subscribeKey := fmt.Sprintf("%s-%s", subscrInfo.TargetNfType, subscrInfo.TargetServiceName)
		workerManager.InjectSuccessSubscribeTask(requsterNfType, subscribeKey)
	} else {
		ok := cacheManager.AddRoamingSubscriptionInfo(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo)
		if !ok {
			log.Errorf("add subscrInfo event hander AddRoamingSubscriptionInfo failure: %s, %s, %s", subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID)
			return false
		}
		cacheManager.DelRoamingSubscriptionMonitor(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID)
		cacheManager.SuperviseRoamingSubscription(subscrInfo.RequesterNfType, subscrInfo.TargetNfType, subscrInfo.SubscriptionID, subscrInfo.ValidityTime)
		log.Debugf("Add %s subscrInfo event hander finished", subscrInfo.SubscriptionID)
	}

	return true
}

/*
func sendSubscrInfoToSlave(subscriptionInfo structs.SubscriptionInfo) bool {
	//master aad subscriptionInfo to slave validityTime delay 3 seconds compare with master
	validityTime := subscriptionInfo.ValidityTime.Add(-defaultTimeDeltaForSlave)
	subscriptionInfo.ValidityTime = validityTime
	syncSubscrInfoMsg := structs.DiscDiscInnerMsg{
		EventType:  consts.EventTypeSyncSubscrInfo,
		SubscrInfo: subscriptionInfo,
	}
	jsonBuf, err := json.Marshal(syncSubscrInfoMsg)
	if err != nil {
		log.Errorf("sendSyncSubscrInfoToSlave: Failed to Marshal Disc message, Error: %s", err.Error())
		return false
	}

	innerTopicName := consts.MsgbusTopicNamePrefix + consts.DiscDiscInner
	return sendToMessageBus(innerTopicName, string(jsonBuf))
}
*/

///////////common code////////////////

func reFetchNFProfile(requesterNfType, targetNfType, nfInstanceID string) {
	targetNf, ok := cacheManager.GetTargetNf(requesterNfType, targetNfType)
	if !ok {
		log.Errorf("reFetchNFProfile: target NFProfile does not exist in cache")
		return
	}
	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		log.Infof("reFetchNFProfile: NRF Agent role is leader, will fetch nfProfile %s to NRF", nfInstanceID)
		handleDiscoveryRequest(&targetNf, nfInstanceID)
	}
}

func convertToSearchResultNFProfile(searchResultNFProfile []byte) ([]byte, error) {
	var sNfProfile structs.SearchResultNFProfile
	err := json.Unmarshal(searchResultNFProfile, &sNfProfile)
	if err != nil {
		return nil, err
	}
	sNFProfileByte, err := json.Marshal(sNfProfile)
	if err != nil {
		return nil, err
	}
	return sNFProfileByte, nil
}

func convertToSearchResultBody(searchResultNFProfile []byte) ([]byte, error) {
	//convert Ipv6Address to Ipv4Address in SearchResult
	searchResult := structs.SearchResult{
		ValidityPeriod: 86400,
	}
	var sNfProfile structs.SearchResultNFProfile
	err := json.Unmarshal(searchResultNFProfile, &sNfProfile)
	if err != nil {
		return nil, err
	}
	searchResult.NfInstances = append(searchResult.NfInstances, sNfProfile)
	jsonSearchResult, err := json.Marshal(searchResult)
	if err != nil {
		return nil, err
	}
	jsonSearchResult, err = common.ConvertIpv6ToIpv4InSearchResult(jsonSearchResult, cm.IsEnableConvertIpv6ToIpv4())
	if err != nil {
		return nil, err
	}
	return jsonSearchResult, nil
}

func applyPatchItems(oldNfProfile []byte, patchData []byte) ([]byte, bool) {
	log.Debug("Apply patch item to nfProfile")
	var updatedNfProfile []byte
	p, err := jsonpatch.DecodePatch(patchData)
	if err == nil {
		if updatedNfProfile, err = p.Apply(oldNfProfile); err != nil {
			log.Errorf("Apply patch fail, err:%s", err.Error())
			return oldNfProfile, false
		}
		return updatedNfProfile, true
	} else {
		log.Errorf("Decode patch fail, err:%s", err.Error())
		return oldNfProfile, false
	}

}

/*
func sendToMessageBus(topicName string, messageData string) bool {
	log.Debugf("Send message to message bus :Topic: %s message: %s", topicName, messageData)
	if discMsgbus := common.GetDiscMsgbus(); discMsgbus != nil {
		err := discMsgbus.SendMessage(topicName, messageData)
		if err != nil {
			log.Errorf("%s:Failed to send message to message bus, %s", topicName, err.Error())
			return false
		}
	} else {
		log.Errorf("message bus is not initialized")
		return false
	}
	return true
}
*/
