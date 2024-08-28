package cache

import (
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/timer"
)

//subscriptionCache is for unsubscribe URL with subscriptionID cache
type subscriptionCache struct {
	requesterNfType string
	targetNfType    string
	timeManager     *timer.Timer

	mcacheType cacheType

	//subscriptionMutex       sync.Mutex
	nfInstanceIDContainer   map[string]string
	subscriptionIDContainer []string // subscriptionID list

	subscriptionInfoMutex     sync.Mutex
	subscriptionInfoContainer map[string]structs.SubscriptionInfo //key: subscriptionID, value: subscriptionInfo
}

func (sp *subscriptionCache) init(requesterNfType, targetNfType string, mcacheType cacheType) {
	sp.requesterNfType = requesterNfType
	sp.targetNfType = targetNfType
	sp.nfInstanceIDContainer = make(map[string]string)
	sp.subscriptionIDContainer = make([]string, 0)
	sp.subscriptionInfoContainer = make(map[string]structs.SubscriptionInfo)
	sp.timeManager = timer.NewTimer()
	sp.mcacheType = mcacheType

	sp.startMonitorWorker()
}

func (sp *subscriptionCache) startMonitorWorker() {
	go func() {
		for subscriptionID := range sp.timeManager.TimerChan() {
			sp.timeoutHandler(subscriptionID)
		}
	}()
	sp.timeManager.StartTimer()

	log.Infof("Start nfType[%s,%s] subscription TtlMonitor work thread", sp.requesterNfType, sp.targetNfType)
}

func (sp *subscriptionCache) stopMonitorWorker() {
	sp.timeManager.StopTimer()

	log.Infof("Stop nfType[%s,%s] subscription TtlMonitor work thread", sp.requesterNfType, sp.targetNfType)
}

func (sp *subscriptionCache) supervise(subscriptionID string, timepoint time.Time) {
	sp.timeManager.AddTimePoint(&timepoint, subscriptionID)
	log.Infof("Monitor nfType[%s,%s] subscriptionID[%s], timestamp:%v", sp.requesterNfType, sp.targetNfType, subscriptionID, timepoint)
}

func (sp *subscriptionCache) deleteSubscriptionMonitor(subscriptionID string) {
	log.Infof("Delete nfType[%s,%s] subscriptionID[%s] ttlMonitor", sp.requesterNfType, sp.targetNfType, subscriptionID)
	sp.timeManager.DelTimePointTag(subscriptionID)
}

func (sp *subscriptionCache) getNfProfileSubscriptionID(nfInstanceID string) (string, bool) {
	sp.subscriptionInfoMutex.Lock()
	defer sp.subscriptionInfoMutex.Unlock()

	subscriptionID, ok := sp.nfInstanceIDContainer[nfInstanceID]
	if !ok {
		return "", false
	}

	return subscriptionID, true
}

func (sp *subscriptionCache) getSubscriptionInfo(subscriptionID string) (structs.SubscriptionInfo, bool) {
	sp.subscriptionInfoMutex.Lock()
	defer sp.subscriptionInfoMutex.Unlock()

	subscriptionInfo, ok := sp.subscriptionInfoContainer[subscriptionID]
	if !ok {
		return structs.SubscriptionInfo{}, false
	}

	return subscriptionInfo, true
}

func (sp *subscriptionCache) addSubscriptionInfo(subscriptionInfo structs.SubscriptionInfo) {
	sp.subscriptionInfoMutex.Lock()
	defer sp.subscriptionInfoMutex.Unlock()

	subscriptionID := subscriptionInfo.SubscriptionID
	if subscriptionID == "" {
		log.Warn("subscriptionInfo less subscriptionID item, will skip add to subscriptionCache")
		return
	}

	ok, _ := getIndexForSet(sp.subscriptionIDContainer, subscriptionID)
	if !ok {
		sp.subscriptionIDContainer = append(sp.subscriptionIDContainer, subscriptionID)
	}

	if subscriptionInfo.NfInstanceID != "" {
		sp.nfInstanceIDContainer[subscriptionInfo.NfInstanceID] = subscriptionInfo.SubscriptionID
	}

	sp.subscriptionInfoContainer[subscriptionID] = subscriptionInfo
}

func (sp *subscriptionCache) delSubscriptionInfo(subscriptionID string) {
	sp.subscriptionInfoMutex.Lock()
	defer sp.subscriptionInfoMutex.Unlock()

	ok, index := getIndexForSet(sp.subscriptionIDContainer, subscriptionID)
	if !ok {
		log.Warnf("No such subscriptionID[%s] in cache for requesterNfType:%s targetNfType:%s", subscriptionID, sp.requesterNfType, sp.targetNfType)
	} else {
		sp.subscriptionIDContainer = append(sp.subscriptionIDContainer[:index], sp.subscriptionIDContainer[index+1:]...)
	}

	subscriptionInfo, ok := sp.subscriptionInfoContainer[subscriptionID]
	if ok {
		if subscriptionInfo.NfInstanceID != "" {
			delete(sp.nfInstanceIDContainer, subscriptionInfo.NfInstanceID)
		}
	}

	delete(sp.subscriptionInfoContainer, subscriptionID)
}

func (sp *subscriptionCache) probeSubscriptionInfo(serviceName string) bool {
	sp.subscriptionInfoMutex.Lock()
	defer sp.subscriptionInfoMutex.Unlock()

	for _, subscriptionInfo := range sp.subscriptionInfoContainer {
		if subscriptionInfo.TargetServiceName == serviceName {
			return true
		}
	}

	return false
}

/*
func (sp *subscriptionCache) addSubscriptionID(subscriptionID string) bool {
	if len(subscriptionID) == 0 {
		log.Warnf("addSubscriptionID: nfType[%s,%s] subsId is empty", sp.requesterNfType, sp.targetNfType)
		return false
	}

	sp.subscriptionMutex.Lock()
	defer sp.subscriptionMutex.Unlock()

	ok, _ := getIndexForSet(sp.subscriptionIDContainer, subscriptionID)
	if ok {
		log.Warnf("addSubscriptionID: nfType[%s,%s] subIdURL(%s) already added", sp.requesterNfType, sp.targetNfType, subscriptionID)
		return false
	}
	sp.subscriptionIDContainer = append(sp.subscriptionIDContainer, subscriptionID)

	return true
}
*/
/*
func (sp *subscriptionCache) delSubscriptionID(subscriptionID string) bool {
	if len(subscriptionID) == 0 {
		log.Warnf("DelWithNfType: nfType[%s,%s] subsId is empty", sp.requesterNfType, sp.targetNfType)
		return false
	}

	sp.subscriptionMutex.Lock()
	defer sp.subscriptionMutex.Unlock()

	ok, index := getIndexForSet(sp.subscriptionIDContainer, subscriptionID)
	if !ok {
		log.Warnf("DelWithNfType: nfType[%s,%s] subIdURL(%s) do not exist", sp.requesterNfType, sp.targetNfType, subscriptionID)
		return false
	}
	sp.subscriptionIDContainer = append(sp.subscriptionIDContainer[:index], sp.subscriptionIDContainer[index+1:]...)
	return true
}
*/

func (sp *subscriptionCache) getSubscriptionIDs() ([]string, bool) {
	if len(sp.subscriptionIDContainer) == 0 {
		return nil, false
	}

	return sp.subscriptionIDContainer, true
}

func (sp *subscriptionCache) getServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName string) (structs.SubscriptionInfo, bool) {
	if requesterNfType == "" || targetNfType == "" || serviceName == "" {
		return structs.SubscriptionInfo{}, false
	}

	for _, subscriptionInfo := range sp.subscriptionInfoContainer {
		if subscriptionInfo.RequesterNfType == requesterNfType &&
			subscriptionInfo.TargetNfType == targetNfType &&
			subscriptionInfo.TargetServiceName == serviceName {
			return subscriptionInfo, true
		}
	}

	return structs.SubscriptionInfo{}, false
}

func (sp *subscriptionCache) getNfProfileSubscriptionInfo(nfInstanceID string) (structs.SubscriptionInfo, bool) {
	if nfInstanceID == "" {
		return structs.SubscriptionInfo{}, false
	}

	for _, subscriptionInfo := range sp.subscriptionInfoContainer {
		if subscriptionInfo.NfInstanceID == nfInstanceID {
			return subscriptionInfo, true
		}
	}

	return structs.SubscriptionInfo{}, false
}

func (sp *subscriptionCache) getAllSubscriptionInfo() []structs.SubscriptionInfo {
	var subscriptionInfos = []structs.SubscriptionInfo{}
	for _, subscriptionInfo := range sp.subscriptionInfoContainer {
		subscriptionInfos = append(subscriptionInfos, subscriptionInfo)
	}
	return subscriptionInfos
}

func getIndexForSet(array []string, element string) (bool, int) {
	for i, e := range array {
		if e == element {
			return true, i
		}
	}
	return false, -1
}

////////////time handler////////////////

func (sp *subscriptionCache) timeoutHandler(subscriptionID string) {
	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		if sp.mcacheType == homeCache {
			go sp.masterTimeoutHandler(subscriptionID)
		} else {
			go sp.masterRoamTimeoutHandler(subscriptionID)
		}
	} else {
		if sp.mcacheType == homeCache {
			go sp.slaveTimeoutHandler(subscriptionID)
		} else {
			go sp.slaveRoamTimeoutHandler(subscriptionID)
		}
	}
}

func (sp *subscriptionCache) masterTimeoutHandler(subscriptionID string) {
	subscriptionInfo, existed := sp.getSubscriptionInfo(subscriptionID)
	if !existed {
		log.Warnf("Subscription(%s) does not exist in cache, skip call timeoutHandler", subscriptionID)
		return
	}
	log.Infof("nfType[%s,%s] Subscription(%s) expired, %+v", sp.requesterNfType, sp.targetNfType, subscriptionID, subscriptionInfo)

	oneSubsData := structs.OneSubscriptionData{
		RequesterNfType:   subscriptionInfo.RequesterNfType,
		TargetNfType:      subscriptionInfo.TargetNfType,
		TargetServiceName: subscriptionInfo.TargetServiceName,
		//NotifCondition:    subscriptionInfo.NotifCondition, //skip the parameter
	}

	newValidityTime := prolongSubscriptionFromNrf(&oneSubsData, subscriptionID) // failure return "" and time.Time{}
	if newValidityTime != nil {
		log.Warnf("Prolong subscription[%s] from NRF-MGMT to timestamp[%v] success", subscriptionID, *newValidityTime)
		subscriptionInfo.ValidityTime = *newValidityTime
		sp.addSubscriptionInfo(subscriptionInfo)

		updateConfigmapStorage(sp.subscriptionInfoContainer)
		sp.supervise(subscriptionID, *newValidityTime)

		return
	}

	/////do new subscribe to NRF/////
	log.Warnf("Prolong subscription[%s] from NRF-MGMT fail, will do a new subscribe", subscriptionID)

	//sp.delSubscriptionID(subscriptionID)
	sp.delSubscriptionInfo(subscriptionID)
	updateConfigmapStorage(sp.subscriptionInfoContainer)

	go func() {
		var newSubscriptionID string
		var newValidityTime *time.Time
		for {
			newSubscriptionID, newValidityTime = doSubscriptionToNrf(&oneSubsData)
			if newSubscriptionID == "" && newValidityTime == nil {
				log.Infof("Do subscribe from NRF-MGMT fail, will retry after 5 seconds")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Infof("Do subscription from NRF-MGMT success, subscriptionID[%s], validityTime[%v]", newSubscriptionID, *newValidityTime)
		subscriptionInfo.SubscriptionID = newSubscriptionID
		subscriptionInfo.ValidityTime = *newValidityTime

		sp.addSubscriptionInfo(subscriptionInfo)

		updateConfigmapStorage(sp.subscriptionInfoContainer)

		sp.supervise(newSubscriptionID, *newValidityTime)

		//subscribeTargetNfWorker(&targetNf)

		// why need sync data again, because if the subscription timeout, some ntf message maybe skip
		//handleDiscoveryRequest(&targetNf, "")

		requesterNfType := subscriptionInfo.RequesterNfType
		targetNfType := subscriptionInfo.TargetNfType
		targetNf, ok := Instance().GetTargetNf(requesterNfType, targetNfType)
		if !ok {
			log.Errorf("%s targetNf in cache is nil, So can not fetch profile without target-nf-type and service-names", requesterNfType)
			return
		}
		ok = SyncNrfData(&targetNf, false, nil)
		if !ok {
			log.Infof("Sync profile data for targetNF[%v] from NRF-Disc failure", targetNf)
		}
	}()
}

func (sp *subscriptionCache) masterRoamTimeoutHandler(subscriptionID string) {
	subscriptionInfo, existed := sp.getSubscriptionInfo(subscriptionID)
	if !existed {
		log.Warnf("Subscription(%s) does not exist in cache, skip call timeoutHandler", subscriptionID)
		return
	}
	log.Infof("Roam nfProfile[%s] Subscription(%s) expired, %+v", subscriptionInfo.NfInstanceID, subscriptionID, subscriptionInfo)

	newValidityTime := prolongRoamSubscriptionFromNrf(subscriptionID)
	if newValidityTime != nil {
		log.Warnf("Prolong subscription[%s] from NRF-MGMT to timestamp[%v] success", subscriptionID, *newValidityTime)
		subscriptionInfo.ValidityTime = *newValidityTime
		sp.addSubscriptionInfo(subscriptionInfo)

		//updateConfigmapStorage(sp.subscriptionInfoContainer)
		sp.supervise(subscriptionID, *newValidityTime)

		return
	}

	/////do new subscribe to NRF/////
	log.Warnf("Prolong subscription[%s] from NRF-MGMT fail, will do a new subscribe", subscriptionID)

	sp.delSubscriptionInfo(subscriptionID)
	//updateConfigmapStorage(sp.subscriptionInfoContainer)

	go func() {
		var newSubscriptionID string
		var newValidityTime time.Time
		var err error

		for {
			newSubscriptionID, newValidityTime, err = doRoamSubscriptionToNrf(subscriptionInfo)
			if err != nil {
				log.Infof("Do subscribe from NRF-MGMT fail, err:%s, will retry after 5 seconds", err.Error())
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Infof("Do subscription from NRF-MGMT success, subscriptionID[%s], validityTime[%v]", newSubscriptionID, newValidityTime)
		subscriptionInfo.SubscriptionID = newSubscriptionID
		subscriptionInfo.ValidityTime = newValidityTime

		sp.addSubscriptionInfo(subscriptionInfo)
		//updateConfigmapStorage(sp.subscriptionInfoContainer)
		sp.supervise(newSubscriptionID, newValidityTime)
		//subscribeTargetNfWorker(&targetNf)
	}()
}

func (sp *subscriptionCache) slaveTimeoutHandler(subscriptionID string) {
	subscriptionInfo, existed := sp.getSubscriptionInfo(subscriptionID)
	if !existed {
		log.Warnf("Subscription(%s) does not exist in cache, skip call timeoutHandler", subscriptionID)
		return
	}
	log.Infof("nfType[%s,%s] Subscription(%s) expired, %+v", sp.requesterNfType, sp.targetNfType, subscriptionID, subscriptionInfo)

	sp.delSubscriptionInfo(subscriptionID)

	go func() {
		var newSubscriptionID string
		var validityTime *time.Time
		for {
			newSubscriptionID, validityTime = fetchSubscriptionInfoFromMaster(&subscriptionInfo, false)
			if newSubscriptionID == "" && validityTime == nil {
				log.Infof("Fetch subscriptionInfo from master fail, will retry after 5 seconds")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Infof("Fetch subscriptionInfo from master success, subscriptionID[%s], validityTime[%v]", subscriptionID, *validityTime)

		subscriptionInfo.SubscriptionID = newSubscriptionID
		subscriptionInfo.ValidityTime = *validityTime

		if subscriptionID != newSubscriptionID {
			sp.addSubscriptionInfo(subscriptionInfo)
		}

		sp.supervise(newSubscriptionID, *validityTime)
	}()
}

func (sp *subscriptionCache) slaveRoamTimeoutHandler(subscriptionID string) {
	subscriptionInfo, existed := sp.getSubscriptionInfo(subscriptionID)
	if !existed {
		log.Warnf("Subscription(%s) does not exist in cache, skip call timeoutHandler", subscriptionID)
		return
	}
	log.Infof("nfType[%s,%s] Subscription(%s) expired, %+v", sp.requesterNfType, sp.targetNfType, subscriptionID, subscriptionInfo)

	sp.delSubscriptionInfo(subscriptionID)

	go func() {
		var newSubscriptionID string
		var validityTime *time.Time
		for {
			newSubscriptionID, validityTime = fetchSubscriptionInfoFromMaster(&subscriptionInfo, true)
			if newSubscriptionID == "" && validityTime == nil {
				log.Infof("Fetch subscriptionInfo from master fail, will retry after 5 seconds")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
		log.Infof("Fetch subscriptionInfo from master success, subscriptionID[%s], validityTime[%v]", subscriptionID, *validityTime)

		subscriptionInfo.SubscriptionID = newSubscriptionID
		subscriptionInfo.ValidityTime = *validityTime

		if subscriptionID != newSubscriptionID {
			sp.addSubscriptionInfo(subscriptionInfo)
		}

		sp.supervise(newSubscriptionID, *validityTime)
	}()
}
