package cache

import (
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

type cacheAdapter struct {
	adapterType cacheType
	cacheMap    map[string]*cache //key: targetNfType
}

func (ca *cacheAdapter) initCache(requesterNfType string, targetNfType string, indexGroup []string, subscription *subscriptionCache, mcacheType cacheType) {
	if ca.cacheMap == nil {
		log.Warnf("Not init cacheAdapter for requeesterNfType:%s", requesterNfType)
		return
	}

	mcache := new(cache)
	mcache.init(requesterNfType, targetNfType, indexGroup, subscription, mcacheType)

	ca.cacheMap[targetNfType] = mcache
}

func (ca *cacheAdapter) init(adapterType cacheType) {
	ca.adapterType = adapterType
	ca.cacheMap = make(map[string]*cache)
}

func (ca *cacheAdapter) cached(targetNfType string, data []byte) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.cached(data)
}

func (ca *cacheAdapter) getCacheStatus(nfType string) bool {
	cache := ca.getCache(nfType)
	if cache == nil {
		return false
	}

	return cache.getCacheStatus()
}

func (ca *cacheAdapter) setCacheStatus(nfType string, status bool) {
	cache := ca.getCache(nfType)
	if cache == nil {
		return
	}

	cache.setCacheStatus(status)
}

func (ca *cacheAdapter) indexed(targetNfType string, content []byte, indexGroup []string) (string, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return "", false
	}

	return cache.indexed(content, indexGroup)
}

func (ca *cacheAdapter) deCached(targetNfType string, nfInstanceID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.deCached(nfInstanceID)
}

func (ca *cacheAdapter) deCachedByNfType(targetNfType string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.deCachedAll()
}

func (ca *cacheAdapter) deleteAllEtag(targetNfType string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.deleteAllEtag()
}

func (ca *cacheAdapter) deIndex(targetNfType string, nfInstanceID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.deIndex(nfInstanceID)
}

func (ca *cacheAdapter) deIndexByNfType(targetNfType string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}

	cache.deIndexAll()
}

func (ca *cacheAdapter) probe(targetNfType string, nfInstanceID string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}

	return cache.probe(nfInstanceID)
}

func (ca *cacheAdapter) getTargetNfTypes() []string {
	nfTypes := make([]string, 0)
	for nfType, _ := range ca.cacheMap {
		nfTypes = append(nfTypes, nfType)
	}

	return nfTypes
}

func (ca *cacheAdapter) fetchProfileByID(nfType string, nfInstanceID string) []byte {
	cache := ca.getCache(nfType)
	if cache == nil {
		return nil
	}

	return cache.fetchProfileByID(nfInstanceID)
}

func (ca *cacheAdapter) showIndexContent(nfType string) {
	cache := ca.getCache(nfType)
	if cache == nil {
		return
	}
	cache.showIndexContent()
}

func (ca *cacheAdapter) flush(nfType string) {
	cache := ca.getCache(nfType)
	if cache == nil {
		return
	}

	cache.flush()
}

func (ca *cacheAdapter) flushAll() {
	for _, cache := range ca.cacheMap {
		cache.flush()
	}
}

func (ca *cacheAdapter) sync(syncData *structs.CacheSyncData) {
	if syncData == nil {
		log.Error("syncData is nil")
		return
	}

	for targetNfType, cache := range ca.cacheMap {
		cacheInfo := structs.CacheSyncInfo{
			TargetNfType: targetNfType,
		}
		if cache == nil {
			log.Errorf("targetNfType:%s cache is nil", targetNfType)
			continue
		}

		cache.sync(&cacheInfo)

		if ca.adapterType == homeCache {
			syncData.CacheInfos = append(syncData.CacheInfos, cacheInfo)
		} else if ca.adapterType == roamingCache {
			syncData.RoamingCacheInfos = append(syncData.RoamingCacheInfos, cacheInfo)
		}

	}

	/*
		for targetNfType, cache := range ca.cacheMap {
			cacheInfo := structs.CacheSyncInfo{
				TargetNfType: targetNfType,
			}
			if cache == nil {
				log.Errorf("targetNfType:%s cache is nil", targetNfType)
				continue
			}

			//	if ca.cacheType == homeCache {
			//		cache.sync(&cacheInfo, false)
			//	} else if ca.cacheType == roamingCache {
			//		cache.sync(&cacheInfo, true)
			//	}

			cache.sync(&cacheInfo)
			//syncData.CacheInfos = append(syncData.CacheInfos, cacheInfo)
			cacheSyncInfos = append(cacheSyncInfos, cacheInfo)
		}
	*/
}

func (ca *cacheAdapter) dump(dumpData *structs.CacheDumpData) {
	//func (ca *cacheAdapter) dump(cacheDumpInfos []structs.CacheDumpInfo) {
	if dumpData == nil {
		log.Errorf("dumpData is nil")
		return
	}

	for targetNfType, cache := range ca.cacheMap {
		cacheInfo := structs.CacheDumpInfo{
			TargetNfType: targetNfType,
		}
		if cache == nil {
			log.Errorf("targetNfType:%s cache is nil", targetNfType)
			continue
		}

		cache.dump(&cacheInfo)
		if ca.adapterType == homeCache {
			dumpData.CacheInfos = append(dumpData.CacheInfos, cacheInfo)
		} else if ca.adapterType == roamingCache {
			dumpData.RoamingCacheInfos = append(dumpData.RoamingCacheInfos, cacheInfo)
		}
	}
}

func (ca *cacheAdapter) fetchProfileIDs(targetNfType string) []string {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return nil
	}

	return cache.fetchIDs()
}

func (ca *cacheAdapter) fetchAllProfileIDs() []string {
	ids := make([]string, 0)
	targetNfTypes := ca.getTargetNfTypes()
	for _, nfType := range targetNfTypes {
		id := ca.fetchProfileIDs(nfType)
		ids = append(ids, id...)
	}

	return ids
}

func (ca *cacheAdapter) dumpByID(targetNfType string, nfInstanceID string) []byte {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return nil
	}

	return cache.dumpByID(nfInstanceID)
}

func (ca *cacheAdapter) getCache(nfType string) *cache {
	cache := ca.cacheMap[nfType]
	if cache == nil {
		log.Warnf("CacheAdpter doesn't build cache for nfType[%s]", nfType)
		return nil
	}

	return cache
}

////////////////////ttlMonitor related function//////////////////////////
//Supervise supervise with ttl
func (ca *cacheAdapter) supervise(targetNfType, nfInstanceID string, second uint) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}
	ttlMonitor.supervise(nfInstanceID, second)
}

//SuperviseAll supervise all instances with ttl
func (ca *cacheAdapter) superviseAll(targetNfType string, second uint) {
	cache := ca.cacheMap[targetNfType]
	if cache == nil {
		return
	}
	ids := cache.fetchIDs()
	if len(ids) == 0 {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}

	ttlMonitor.superviseAll(ids, second)
}

func (ca *cacheAdapter) superviseTimestamp(targetNfType, nfInstanceID string, ttl time.Time) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}
	ttlMonitor.superviseTimestamp(nfInstanceID, ttl)
}

func (ca *cacheAdapter) superviseDefaultTTL(targetNfType, nfInstanceID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}
	ttlMonitor.superviseDefaultTTL(nfInstanceID)
}

func (ca *cacheAdapter) stopByID(targetNfType string, nfInstanceID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}
	ttlMonitor.stop(nfInstanceID)
	ttlMonitor.delete(nfInstanceID)
}

func (ca *cacheAdapter) stopAllMonitorWorker() {
	for _, cache := range ca.cacheMap {
		if cache != nil {
			//stop ttlMonitor Worker
			ttlMonitor := cache.getTtlMonitor()
			if ttlMonitor != nil {
				ttlMonitor.stopMonitorWorker()
			}
		}
	}
}

func (ca *cacheAdapter) startAllMonitorWorker() {
	for _, cache := range ca.cacheMap {
		if cache != nil {
			//start ttlMonitor Worker
			ttlMonitor := cache.getTtlMonitor()
			if ttlMonitor != nil {
				ttlMonitor.startMonitorWorker()
			}
		}
	}
}

func (ca *cacheAdapter) startMonitorWorker(targetNfType string) {
	cache := ca.cacheMap[targetNfType]
	if cache == nil {
		return
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return
	}
	ttlMonitor.startMonitorWorker()
}

func (ca *cacheAdapter) setTtlMonitor(targetNfType string, ttlMonitor *ttlMonitor) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	cache.setTtlMonitor(ttlMonitor)
}

func (ca *cacheAdapter) stopAllTtlMonitor() {
	for _, cache := range ca.cacheMap {
		if cache != nil {
			ttlMonitor := cache.getTtlMonitor()
			if ttlMonitor == nil {
				continue
			}
			ttlMonitor.stopAll()
			ttlMonitor.deleteAll()
		}
	}
}

func (ca *cacheAdapter) leftLive(targetNfType, nfInstanceID string) (uint, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return 0, false
	}
	ttlMonitor := cache.getTtlMonitor()
	if ttlMonitor == nil {
		return 0, false
	}

	return ttlMonitor.leftLive(nfInstanceID)
}

func (ca *cacheAdapter) haveEtag(targetNfType, nfInstanceID string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}

	return cache.haveEtag(nfInstanceID)
}

func (ca *cacheAdapter) fetchEtag(targetNfType, nfInstanceID string) string {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return ""
	}
	return cache.fetchEtag(nfInstanceID)
}

func (ca *cacheAdapter) saveEtag(targetNfType, nfInstanceID string, value string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}
	return cache.saveEtag(nfInstanceID, value)
}

func (ca *cacheAdapter) deleteEtag(targetNfType, nfInstanceID string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}
	return cache.deleteEtag(nfInstanceID)
}

////////////////////subscriptionCache related function//////////////////////////
func (ca *cacheAdapter) stopAllSubscrMonitorWorker() {
	for _, cache := range ca.cacheMap {
		if cache != nil {
			//stop subscriptionCache Worker
			subscription := cache.getSubscriptionCache()
			if subscription != nil {
				subscription.stopMonitorWorker()

			}
		}
	}
}

func (ca *cacheAdapter) startSubscrMonitorWorker(targetNfType string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return
	}
	subscription.startMonitorWorker()
}

func (ca *cacheAdapter) superviseSubscription(targetNfType, subscriptionID string, timepoint time.Time) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return
	}
	subscription.supervise(subscriptionID, timepoint)
}

func (ca *cacheAdapter) getSubscriptionMonitor(targetNfType string) *subscriptionCache {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return nil
	}
	return cache.getSubscriptionCache()
}

func (ca *cacheAdapter) getSubscriptionInfo(targetNfType, subscriptionID string) (structs.SubscriptionInfo, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return structs.SubscriptionInfo{}, false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return structs.SubscriptionInfo{}, false
	}

	return subscription.getSubscriptionInfo(subscriptionID)
}

func (ca *cacheAdapter) addSubscriptionInfo(targetNfType string, subscriptionInfo structs.SubscriptionInfo) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return false
	}

	subscription.addSubscriptionInfo(subscriptionInfo)
	return true
}

func (ca *cacheAdapter) delSubscriptionInfo(targetNfType, subscriptionID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return
	}
	subscription.delSubscriptionInfo(subscriptionID)
}

func (ca *cacheAdapter) probeSubscriptionInfo(targetNfType, serviceName string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}

	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return false
	}

	return subscription.probeSubscriptionInfo(serviceName)
}

func (ca *cacheAdapter) deleteSubscriptionMonitor(targetNfType, subscriptionID string) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return
	}
	subscription.deleteSubscriptionMonitor(subscriptionID)
}

/*
func (ca *cacheAdapter) addSubscriptionID(targetNfType, subsIdURL string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return false
	}
	return subscription.addSubscriptionID(subsIdURL)
}
*/
/*
func (ca *cacheAdapter) delSubscriptionID(targetNfType, subsIdURL string) bool {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return false
	}
	return subscription.delSubscriptionID(subsIdURL)
}
*/
func (ca *cacheAdapter) getProfileByID(targetNfType string, nfInstanceID string) []byte {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return nil
	}

	return cache.fetchProfileByID(nfInstanceID)
}

func (ca *cacheAdapter) getServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName string) (structs.SubscriptionInfo, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return structs.SubscriptionInfo{}, false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return structs.SubscriptionInfo{}, false
	}
	return subscription.getServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName)
}

func (ca *cacheAdapter) getNfProfileSubscriptionInfo(targetNfType, nfInstanceID string) (structs.SubscriptionInfo, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return structs.SubscriptionInfo{}, false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return structs.SubscriptionInfo{}, false
	}

	return subscription.getNfProfileSubscriptionInfo(nfInstanceID)
}

func (ca *cacheAdapter) getNfProfileSubscriptionID(targetNfType, nfInstanceID string) (string, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		return "", false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		return "", false
	}

	return subscription.getNfProfileSubscriptionID(nfInstanceID)
}

func (ca *cacheAdapter) getSubscriptionIDs(targetNfType string) ([]string, bool) {
	cache := ca.getCache(targetNfType)
	if cache == nil {
		log.Warnf("No cache for targetNfType:%s", targetNfType)
		return nil, false
	}
	subscription := cache.getSubscriptionCache()
	if subscription == nil {
		log.Warnf("No subscriptionCache for targetNfType:%s", targetNfType)
		return nil, false
	}

	return subscription.getSubscriptionIDs()
}

func (ca *cacheAdapter) getAllSubscriptionInfo() map[string]structs.SubscriptionInfo {
	subscriptionInfo := make(map[string]structs.SubscriptionInfo, 0)
	for _, cache := range ca.cacheMap {
		if cache == nil {
			continue
		}
		subscription := cache.getSubscriptionCache()
		if subscription == nil {
			continue
		}
		for subId, subInfo := range subscription.subscriptionInfoContainer {
			subscriptionInfo[subId] = subInfo
		}
	}
	return subscriptionInfo
}
