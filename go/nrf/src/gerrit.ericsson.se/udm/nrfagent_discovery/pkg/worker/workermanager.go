package worker

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
)

type TaskStatus int

const (
	InitStatus    TaskStatus = 0
	SuccessStatus TaskStatus = 1
	FailureStatus TaskStatus = 2
	UnknowStatus  TaskStatus = 3
)

var (
	instance     *WorkerManager
	cacheManager *cache.CacheManager
)

var (
	isKeepCacheStatus bool                 //true: keepCache Mode, false: normal Mode
	monitorStarted    bool                 = false
	preWorkMode       client.NRFConnStatus = client.NRFConnUnknown
	curWorkMode       client.NRFConnStatus = client.NRFConnUnknown
	quitMonitor                            = make(chan bool)
)

//WorkerManager worker manager
type WorkerManager struct {
	subscribeContainerMutex    sync.Mutex
	fetchProfileContainerMutex sync.Mutex

	subscribeWorkers      map[string]*SubscribeWorker
	fetchProfileWorkers   map[string]*FetchProfileWorker
	dumpCacheWorkers      map[string]*DumpCacheWorker
	subscribeContainer    map[string]map[string]TaskStatus //key:requesterNfType, value:targetNfType-serviceName
	fetchProfileContainer map[string]map[string]TaskStatus //key:requesterNfType, value:targetNfType

	isKeepCacheStatus bool
}

//Instance get the workerManager instance
func Instance() *WorkerManager {
	if instance == nil {
		instance = new(WorkerManager)
		instance.init()
	}

	return instance
}

func (wm *WorkerManager) PrepareDiscoveryAgent(nfInfo structs.NfInfoForRegDisc) {
	nfType := nfInfo.RequesterNfType
	requesterNfFqdn := nfInfo.RequesterNfFqdn
	requesterPlmns := nfInfo.RequesterPlmns

	targetNfs := loopFetchTargetNfs(nfType)
	if len(targetNfs) == 0 {
		return
	}

	//set firstly
	cache.Instance().SetRequesterFqdn(nfType, requesterNfFqdn)
	if len(requesterPlmns) > 0 {
		cache.Instance().SetRequesterPlmns(nfType, requesterPlmns)
	}

	//master and slave have the same task on backlog
	wm.injectSubscribeBacklogTask(nfType, targetNfs)
	wm.injectFetchProfileBacklogTask(nfType, targetNfs)

	wm.syncCacheSubscribeInfo(nfType, targetNfs)

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		subscribeRet := subscribe(nfType)
		if !subscribeRet {
			log.Warnf("Nrf discovery master agent will launch subscribe worker to do the left taskes for nfType:%s", nfType)
			wm.launchSubscribeWorker(nfType)
		}

		fetchProfileRet := fetchProfile(nfType)
		if !fetchProfileRet {
			log.Warnf("Nrf discovery master agent will launch fetchProfile worker to do the left taskes for nfType:%s", nfType)
			wm.launchFetchProfileWorker(nfType)
		}
	} else {
		wm.launchDumpCacheWorker(nfType)
	}
}

func (wm *WorkerManager) PrepareNfRegister(nfType string) {
	if wm.haveBacklog(nfType) {
		log.Infof("Disc agent have plan the nfType:%s register tasks", nfType)
		return
	}

	targetNfs := loopFetchTargetNfs(nfType)
	if len(targetNfs) == 0 {
		return
	}

	//master and slave have the same task on backlog
	workerManager.injectSubscribeBacklogTask(nfType, targetNfs)
	workerManager.injectFetchProfileBacklogTask(nfType, targetNfs)

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		subscribeRet := subscribe(nfType)
		if !subscribeRet {
			log.Warnf("Nrf discovery master agent will launch subscribe worker to do the left taskes for nfType:%s", nfType)
			wm.launchSubscribeWorker(nfType)
		}

		fetchProfileRet := fetchProfile(nfType)
		if !fetchProfileRet {
			log.Warnf("Nrf discovery master agent will launch fetchProfile worker to do the left taskes for nfType:%s", nfType)
			wm.launchFetchProfileWorker(nfType)
		}
	} else {
		wm.launchDumpCacheWorker(nfType)
	}
}

func (wm *WorkerManager) WaitAgentReady() bool {
	ret := false
	for {
		ret = wm.agentReady()
		if ret {
			return true
		} else {
			log.Infof("Agent is not ready, will check after 5 seconds")
			time.Sleep(time.Second * 5)
		}
	}

	return false
}

func (wm *WorkerManager) DumpReady(nfType string) bool {
	ret1 := wm.checkSubscribeAllTaskDone(nfType)
	ret2 := wm.checkFetchProfileAllTaskDone(nfType)

	return ret1 && ret2
}

func (wm *WorkerManager) LaunchAllLeftTask() {
	for nfType, _ := range wm.subscribeContainer {
		wm.launchLeftTask(nfType)
	}
}

func (wm *WorkerManager) StopSubscribeWorker(nfType string) {
	subscribeWorker, ok := wm.subscribeWorkers[nfType]
	if !ok {
		log.Printf("WorkerManager no subscribe worker for %s", nfType)
		return
	}

	if subscribeWorker.IsRunning() {
		subscribeWorker.Stop()
	}
}

func (wm *WorkerManager) StopFetchProfileWorker(nfType string) {
	fetchProfileWorker, ok := wm.fetchProfileWorkers[nfType]
	if !ok {
		log.Printf("WorkerManager no fetchProfile worker for %s", nfType)
		return
	}

	if fetchProfileWorker.IsRunning() {
		fetchProfileWorker.Stop()
	}
}

func (wm *WorkerManager) StopDumpCacheWorker(nfType string) {
	dumpCacheWorker, ok := wm.dumpCacheWorkers[nfType]
	if !ok {
		log.Printf("WorkerManager no dumpCache worker for %s", nfType)
		return
	}

	if dumpCacheWorker.IsRunning() {
		dumpCacheWorker.Stop()
	}
}

func (wm *WorkerManager) StopAllWorker() {
	for _, subWorker := range wm.subscribeWorkers {
		if subWorker.IsRunning() {
			subWorker.Stop()
		}
	}

	for _, fetchWorker := range wm.fetchProfileWorkers {
		if fetchWorker.IsRunning() {
			fetchWorker.Stop()
		}
	}

	for _, dumpWorker := range wm.dumpCacheWorkers {
		if dumpWorker.IsRunning() {
			dumpWorker.Start()
		}
	}
}

func (wm *WorkerManager) InjectSuccessSubscribeTask(nfType string, key string) {
	wm.setSubscribeTaskStatus(nfType, key, SuccessStatus)
}

////////////////private/////////////////

func (wm *WorkerManager) init() {
	wm.subscribeWorkers = make(map[string]*SubscribeWorker)
	wm.fetchProfileWorkers = make(map[string]*FetchProfileWorker)
	wm.dumpCacheWorkers = make(map[string]*DumpCacheWorker)

	wm.subscribeContainer = make(map[string]map[string]TaskStatus)
	wm.fetchProfileContainer = make(map[string]map[string]TaskStatus)
	cacheManager = cache.Instance()
}

func (wm *WorkerManager) agentReady() bool {
	for _, fetchProfileTasks := range wm.fetchProfileContainer {
		for _, status := range fetchProfileTasks {
			if status == SuccessStatus {
				return true
			}
		}
	}

	return false
}

func (wm *WorkerManager) launchSubscribeWorker(nfType string) {
	workerName := fmt.Sprintf("%s-subscribe-worker", nfType)

	subscribeWorkerThread := SubscribeWorker{
		interval: 5,
		nfType:   nfType,
		callBack: subscribe,
		Worker: Worker{
			name:     workerName,
			stopFlag: false,
		},
	}

	subscribeWorkerThread.Start()

	wm.subscribeWorkers[nfType] = &subscribeWorkerThread
}

func (wm *WorkerManager) launchFetchProfileWorker(nfType string) {
	workerName := fmt.Sprintf("%s-fetchProfile-worker", nfType)

	fetchProfileWorkerThread := FetchProfileWorker{
		interval: 5,
		nfType:   nfType,
		callBack: fetchProfile,
		Worker: Worker{
			name:     workerName,
			stopFlag: false,
		},
	}

	fetchProfileWorkerThread.Start()

	wm.fetchProfileWorkers[nfType] = &fetchProfileWorkerThread
}

func (wm *WorkerManager) launchDumpCacheWorker(nfType string) {
	workerName := fmt.Sprintf("%s-dumpCache-worker", nfType)

	dumpCacheWorkerThread := DumpCacheWorker{
		interval: 5,
		nfType:   nfType,
		callBack: dumpCacheFromMaster,
		Worker: Worker{
			name:     workerName,
			stopFlag: false,
		},
	}

	dumpCacheWorkerThread.Start()

	wm.dumpCacheWorkers[nfType] = &dumpCacheWorkerThread
}

func (wm *WorkerManager) injectSubscribeBacklogTask(nfType string, targetNfs []structs.TargetNf) bool {
	if targetNfs == nil {
		return false
	}

	subscribeTask := make(map[string]TaskStatus)

	for _, targetNf := range targetNfs {
		targetNfType := targetNf.TargetNfType

		for _, serviceName := range targetNf.TargetServiceNames {
			subscribeKey := fmt.Sprintf("%s-%s", targetNfType, serviceName)
			subscribeTask[subscribeKey] = InitStatus
		}
	}

	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()
	wm.subscribeContainer[nfType] = subscribeTask

	log.Debugf("Set subscribe backlog tasks for nfType:%s success, tasks:%v", nfType, subscribeTask)
	return true
}

func (wm *WorkerManager) injectFetchProfileBacklogTask(nfType string, targetNfs []structs.TargetNf) bool {
	if targetNfs == nil {
		return false
	}

	fetchProfileTask := make(map[string]TaskStatus)

	for _, targetNf := range targetNfs {
		targetNfType := targetNf.TargetNfType

		fetchProfileKey := fmt.Sprintf("%s", targetNfType)
		fetchProfileTask[fetchProfileKey] = InitStatus
	}

	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()
	wm.fetchProfileContainer[nfType] = fetchProfileTask

	log.Debugf("Set fetchProfile backlog tasks for nfType:%s success, tasks:%v", nfType, fetchProfileTask)
	return true
}

func (wm *WorkerManager) resetSubscribeBacklogTask(nfType string) bool {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return true
	}

	for key, _ := range backLogTask {
		backLogTask[key] = InitStatus
	}

	return true
}

func (wm *WorkerManager) resetFetchProfileBacklogTask(nfType string) bool {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return true
	}

	for key, _ := range backLogTask {
		backLogTask[key] = InitStatus
	}

	return true
}

func (wm *WorkerManager) setSubscribeTaskStatus(nfType string, key string, status TaskStatus) {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return
	}

	_, ok = backLogTask[key]
	if !ok {
		log.Warnf("Do not need do such subscribe for : %s", key)
		return
	}

	log.Infof("Set NF:%s subscribe backlog task:%s status:%d", nfType, key, status)
	backLogTask[key] = status
}

func (wm *WorkerManager) setFetchProfileTaskStatus(nfType string, key string, status TaskStatus) {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return
	}

	_, ok = backLogTask[key]
	if !ok {
		log.Warnf("Do not need fetch profile for : %s", key)
		return
	}

	log.Infof("Set NF:%s fetchProfile backlog task:%s status:%d", nfType, key, status)
	backLogTask[key] = status
}

func (wm *WorkerManager) checkSubscribeAllTaskSuccess(nfType string) bool {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return true
	}

	for key, ret := range backLogTask {
		if ret != SuccessStatus {
			log.Infof("NF:%s subscribe backlog task:%s status:%d, not success", nfType, key, ret)
			return false
		}
	}

	log.Infof("NF:%s subscribe all backlog task success", nfType)
	return true
}

func (wm *WorkerManager) checkFetchProfileAllTaskSuccess(nfType string) bool {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return true
	}

	for key, ret := range backLogTask {
		if ret != SuccessStatus {
			log.Infof("NF:%s fetchProfile backlog task:%s status:%d, not success", nfType, key, ret)
			return false
		}
	}

	log.Infof("NF:%s fetchProfile all backlog task success", nfType)
	return true
}

func (wm *WorkerManager) checkSubscribeAllTaskDone(nfType string) bool {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return true
	}

	for key, ret := range backLogTask {
		if ret == InitStatus {
			log.Infof("NF:%s subscribe backlog task:%s status:InitStauts, not Done", nfType, key)
			return false
		}
	}

	log.Infof("NF:%s subscribe all backlog task done", nfType)
	return true
}

func (wm *WorkerManager) checkFetchProfileAllTaskDone(nfType string) bool {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return true
	}

	for key, ret := range backLogTask {
		if ret == InitStatus {
			log.Infof("NF:%s fetchProfile backlog task:%s status:InitStauts, not Done", nfType, key)
			return false
		}
	}

	log.Infof("NF:%s fetchProfile all backlog task done", nfType)
	return true
}

func (wm *WorkerManager) checkSubscribeTaskSuccess(nfType string, key string) bool {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return true
	}

	log.Infof("Check subscribe task:%s status", key)
	ret, ok := backLogTask[key]
	if !ok {
		log.Warnf("Do not need do such subscribe for : %s", key)
		return true
	}
	log.Infof("RequesterNfType:%s subscribe task[%s] status:%v", nfType, key, ret)

	return ret == SuccessStatus
}

func (wm *WorkerManager) checkFetchProfileTaskSuccess(nfType string, key string) bool {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return true
	}

	log.Infof("Check fetchProfile task:%s status", key)
	ret, ok := backLogTask[key]
	if !ok {
		log.Warnf("Do not need fetch profile for : %s", key)
		return true
	}
	log.Infof("RequesterNfType:%s fetchProfile task[%s] status:%v", nfType, key, ret)

	return ret == SuccessStatus
}

func (wm *WorkerManager) fetchSubscribeTaskStatus(nfType string, key string) TaskStatus {
	wm.subscribeContainerMutex.Lock()
	defer wm.subscribeContainerMutex.Unlock()

	backLogTask, ok := wm.subscribeContainer[nfType]
	if !ok {
		log.Warnf("No such subscribe backlog task for %s", nfType)
		return UnknowStatus
	}

	log.Infof("Check subscribe task:%s status", key)
	ret, ok := backLogTask[key]
	if !ok {
		log.Warnf("Do not need do such subscribe for : %s", key)
		return UnknowStatus
	}
	log.Infof("RequesterNfType:%s subscribe task[%s] status:%v", nfType, key, ret)

	return ret
}

func (wm *WorkerManager) fetchFetchProfileTaskStatus(nfType string, key string) TaskStatus {
	wm.fetchProfileContainerMutex.Lock()
	defer wm.fetchProfileContainerMutex.Unlock()

	backLogTask, ok := wm.fetchProfileContainer[nfType]
	if !ok {
		log.Warnf("No such fetchProfile backlog task for %s", nfType)
		return UnknowStatus
	}

	log.Infof("Check fetchProfile task:%s status", key)
	ret, ok := backLogTask[key]
	if !ok {
		log.Warnf("Do not need fetch profile for : %s", key)
		return UnknowStatus
	}
	log.Infof("RequesterNfType:%s fetchProfile task[%s] status:%v", nfType, key, ret)

	return ret
}

func (wm *WorkerManager) launchLeftTask(nfType string) {
	subscribeRet := subscribe(nfType)
	if !subscribeRet {
		wm.launchSubscribeWorker(nfType)
	}

	fetchProfileRet := fetchProfile(nfType)
	if !fetchProfileRet {
		wm.launchFetchProfileWorker(nfType)
	}
}

func (wm *WorkerManager) haveBacklog(nfType string) bool {
	_, subOk := wm.subscribeContainer[nfType]
	_, fetchOk := wm.fetchProfileContainer[nfType]

	return subOk && fetchOk
}

func (wm *WorkerManager) syncCacheSubscribeInfo(nfType string, targetNfs []structs.TargetNf) {
	for _, targetNf := range targetNfs {
		targetNfType := targetNf.TargetNfType
		for _, serviceName := range targetNf.TargetServiceNames {
			subscribeKey := fmt.Sprintf("%s-%s", targetNfType, serviceName)
			exist := cache.Instance().ProbeSubscriptionInfo(nfType, targetNfType, serviceName)
			if exist {
				wm.setSubscribeTaskStatus(nfType, subscribeKey, SuccessStatus)
			}
		}
	}
}

//StartWorkModeMonitor is for start workMode monitor
func StartWorkModeMonitor(monitorTimer int) {
	log.Debugf("startWorkModeMonitor: Start work mode monitor, timer: %d second", monitorTimer)
	monitorStarted = true
	ticker := time.NewTicker(time.Second * time.Duration(monitorTimer))
	go func() {
		keepCacheCounter := 0
		for {
			select {
			case <-ticker.C:
				cmKeepCacheRetryCount := cm.GetKeepCacheRetryCount()
				if cmKeepCacheRetryCount <= 0 {
					log.Debugf("keep cache function is switch off, configmap KeepCacheRetryCount is %d", cmKeepCacheRetryCount)
					continue
				}
				preWorkMode = curWorkMode
				curWorkMode = client.GetNRFConnStatus()
				if preWorkMode == client.NRFConnLost && curWorkMode == client.NRFConnNormal {
					keepCacheCounter = 0
					if isKeepCacheStatus {
						log.Warningf("Leave keep cache mode, previous discWorkMode: %d", preWorkMode)
						setKeepCacheMode(false)
						cacheManager.EnterNormalWorkMode()
					}
				} else {
					if !isKeepCacheStatus && curWorkMode == client.NRFConnLost {
						if preWorkMode == client.NRFConnLost {
							keepCacheCounter++
						} else if preWorkMode == client.NRFConnNormal {
							keepCacheCounter = 1
						}
						if keepCacheCounter >= cmKeepCacheRetryCount {
							log.Warningf("Enter keep cache mode, keepCacheCounter: %d", keepCacheCounter)
							setKeepCacheMode(true)
							cacheManager.EnterKeepCacheWorkMode()
						}
					}
				}
			case <-quitMonitor:
				ticker.Stop()
				monitorStarted = false
				log.Debugf("Ticker which is used to monitor nf-profile version stops")
				return
			}
		}
	}()
}

// StopWorkModeMonitor is to stop work mode monitor
func StopWorkModeMonitor() {
	if monitorStarted {
		log.Debugf("Stop work mode monitor.")
		quitMonitor <- true
	}

	for {
		if !monitorStarted {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500))
	}
}

//IsKeepCacheMode is for get keep cache status
func IsKeepCacheMode() bool {
	return isKeepCacheStatus
}

func setKeepCacheMode(status bool) {
	isKeepCacheStatus = status
}
