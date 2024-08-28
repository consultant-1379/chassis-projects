package cache

import (
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/timer"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

//ttlMonitor ttl monitor
type ttlMonitor struct {
	requesterNfType string
	targetNfType    string
	defaultTtl      uint
	timeManager     *timer.Timer
	mcacheType      cacheType
}

//init init ttl monitor
func (tm *ttlMonitor) init(requesterNfType, targetNfType string, mcacheType cacheType) {
	tm.requesterNfType = requesterNfType
	tm.targetNfType = targetNfType
	tm.defaultTtl = 7200
	tm.timeManager = timer.NewTimer()
	tm.mcacheType = mcacheType
}

//startMonitorWorker start ttl monitor worker
func (tm *ttlMonitor) startMonitorWorker() {
	go func() {
		for nfInstanceID := range tm.timeManager.TimerChan() {
			go tm.timeoutHandler(nfInstanceID)
		}
	}()
	tm.timeManager.StartTimer()

	log.Infof("Start nfType[%s,%s] profile TtlMonitor work thread", tm.requesterNfType, tm.targetNfType)
}

//stopMonitorWorker stop monitor worker thread
func (tm *ttlMonitor) stopMonitorWorker() {
	tm.timeManager.StopTimer()

	log.Infof("Stop nfType[%s,%s] profile TtlMonitor work thread", tm.requesterNfType, tm.targetNfType)
}

//destroyMonitor destroy monitor instance
func (tm *ttlMonitor) destroyMonitor() {
	tm.timeManager.DestroyTimer()
	tm.timeManager = nil

	log.Infof("Destroy nfType[%s,%s] TtlMonitor", tm.requesterNfType, tm.targetNfType)
}

//superviseDefaultTTL monitor nfProfile with default ttl
func (tm *ttlMonitor) superviseDefaultTTL(nfInstanceID string) {
	now := time.Now()
	validityTime := now.Add(time.Duration(tm.defaultTtl) * time.Second)

	tm.timeManager.AddTimePoint(&validityTime, nfInstanceID)

	log.Infof("Monitor with default ttl for nfInstanceID[%s], lifetime:%d(seconds)", nfInstanceID, tm.defaultTtl)
}

//supervise monitor nfProfile ttl
func (tm *ttlMonitor) supervise(nfInstanceID string, ttl uint) {
	now := time.Now()
	validityTime := now.Add(time.Duration(ttl) * time.Second)
	tm.delete(nfInstanceID)
	tm.timeManager.AddTimePoint(&validityTime, nfInstanceID)
	//tm.timeManager.ResetTimer()

	log.Infof("Monitor ttl for nfInstanceID[%s], lifetime:%d(seconds)", nfInstanceID, ttl)
}

//supervise monitor nfProfile ttl
func (tm *ttlMonitor) superviseAll(nfInstanceIDs []string, ttl uint) {
	now := time.Now()
	validityTime := now.Add(time.Duration(ttl) * time.Second)
	tm.deleteAll()

	for _, nfInstanceID := range nfInstanceIDs {
		tm.timeManager.AddTimePoint(&validityTime, nfInstanceID)
	}
	log.Infof("Monitor ttl for nfInstanceID[%v], lifetime:%d(seconds)", nfInstanceIDs, ttl)
}

func (tm *ttlMonitor) superviseTimestamp(nfInstanceID string, ttl time.Time) {
	tm.delete(nfInstanceID)
	tm.timeManager.AddTimePoint(&ttl, nfInstanceID)

	log.Infof("Monitor ttl for nfInstanceID[%s], timestamp:%v", nfInstanceID, ttl)
}

//stop stop nfProfile ttl timer
func (tm *ttlMonitor) stop(nfInstanceID string) {
	log.Infof("Stop TtlMonitor for nfInstanceID[%s]", nfInstanceID)
	tm.timeManager.StopTimePointTag(nfInstanceID)
}

//stopAll stop all nfProfile timer
func (tm *ttlMonitor) stopAll() {
	log.Infof("Stop TtlMonitor for nfType[%s,%s]", tm.requesterNfType, tm.targetNfType)
	tm.timeManager.StopTimer() // total stop all timer, after reset, the timer will start again
}

//delete delete nfProfile ttl timer
func (tm *ttlMonitor) delete(nfInstanceID string) {
	log.Infof("Delete timer for nfType[%s,%s]:%s", tm.requesterNfType, tm.targetNfType, nfInstanceID)
	tm.timeManager.DelTimePointTag(nfInstanceID)
}

//deleteAll delete all nfProfile ttl timer
func (tm *ttlMonitor) deleteAll() {
	log.Infof("Delete TtlMonitor for nfType[%s,%s]", tm.requesterNfType, tm.targetNfType)
	tm.timeManager.DelTimePointAll()
}

//leftLive get the leftLive ttl
func (tm *ttlMonitor) leftLive(nfInstanceID string) (uint, bool) {
	timeStamp, err := tm.timeManager.GetTimePoint(nfInstanceID)
	if err != nil {
		log.Errorf("%s", err.Error())
		return 0, false
	}

	now := time.Now().Unix()
	end := timeStamp.Unix()
	liveTtl := end - now
	if liveTtl < 0 {
		log.Errorf("The timer of %s have been expired, response ttl set to 0", nfInstanceID)
		return 0, false
	}

	log.Infof("Left ttl for nfType[%s,%s]:%s is %d", tm.requesterNfType, tm.targetNfType, nfInstanceID, liveTtl)
	return uint(liveTtl), true
}

func (tm *ttlMonitor) reset(nfInstanceID string, ttl uint) {
	//tm.Stop(nfInstanceID)
	tm.delete(nfInstanceID)
	tm.supervise(nfInstanceID, ttl)
}

func (tm *ttlMonitor) getTimePoint(nfInstanceID string) (time.Time, error) {
	timePoint, err := tm.timeManager.GetTimePoint(nfInstanceID)
	if err != nil {
		log.Errorf("Get monitor timePoint for profile[%s] fail, err:%s", nfInstanceID, err.Error())
		return time.Time{}, err
	}

	return timePoint, nil
}

func (tm *ttlMonitor) timeoutHandler(nfInstanceID string) {
	log.Infof("NfType[%s,%s]:profile[%s] ttl expired", tm.requesterNfType, tm.targetNfType, nfInstanceID)

	if tm.mcacheType == homeCache {
		ttl, ok := proberNfProfile(tm.requesterNfType, tm.targetNfType, nfInstanceID)
		if !ok {
			Instance().DeCached(tm.requesterNfType, tm.targetNfType, nfInstanceID, false)
			tm.delete(nfInstanceID)
			if election.IsActiveLeader("3201", consts.DiscoveryAgentReadinessProbe) {
				util.PushMessageToMSB(tm.requesterNfType, tm.targetNfType, nfInstanceID, consts.NFDeRegister, nil)
			}
		} else {
			tm.reset(nfInstanceID, ttl)
		}
	} else if tm.mcacheType == roamingCache {
		Instance().DeCached(tm.requesterNfType, tm.targetNfType, nfInstanceID, true)
		tm.delete(nfInstanceID)
		if election.IsActiveLeader("3201", consts.DiscoveryAgentReadinessProbe) {
			err := unsubscribeByNfInstanceID(tm.requesterNfType, tm.targetNfType, nfInstanceID)
			if err != nil {
				log.Warnf("Do unscribe by nfInstanceID fail, err:%s", err.Error())
			}
			util.PushMessageToMSB(tm.requesterNfType, tm.targetNfType, nfInstanceID, consts.NFDeRegister, nil)
		}
	}
}

/*
func (tm *ttlMonitor) buildInit() {
	tm.defaultTtl = 7200
	tm.timeManager = timer.NewTimer()
}
*/
/*
func (tm *ttlMonitor) setNfType(nfType string) {
	tm.nfType = nfType
}
*/
