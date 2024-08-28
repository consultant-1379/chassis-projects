package cache

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	pkgcm "gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/k8sapiclient"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

type cacheType int

const (
	homeCache    cacheType = 0
	roamingCache cacheType = 1
)

var (
	cacheConfig = "../../cache-index.json"

	instance *CacheManager
	nfTypes  []string
)

type CacheManager struct {
	mcache       map[string]*cacheAdapter //key: requestNfType
	roamingCache map[string]*cacheAdapter //key: requestNfType

	targetNfMutex       sync.Mutex
	requesterFqdnMutex  sync.Mutex
	requesterPlmnsMutex sync.Mutex

	targetNf       map[string][]structs.TargetNf //key: requestNfType
	requesterFqdn  map[string]string             //key: requestNfType
	requesterPlmns map[string][]structs.PlmnID

	//TODO: indexGroup change to indexGroup map[string][]string, key: targetNfType
	indexGroup  []string //index
	indexMapper searchIndexMapper
}

//Instance get the cacheManager instance
func Instance() *CacheManager {
	if instance == nil {
		instance = new(CacheManager)
		nfTypes = pkgcm.GetDefaultRequesterNfList()
		ok := instance.init()
		if !ok {
			instance = nil
		}
	}

	return instance
}

//InitCache init cache
func (cm *CacheManager) InitCache(requesterNfType string, targetNfType string) {
	cacheAdapter := cm.mcache[requesterNfType]
	roamingCacheAdapter := cm.roamingCache[requesterNfType]

	if cacheAdapter == nil && roamingCacheAdapter == nil {
		log.Warnf("No such cache adapter from requesterNfType:%s", requesterNfType)
		return
	}

	cacheAdapter.initCache(requesterNfType, targetNfType, cm.indexGroup, nil, homeCache)
	roamingCacheAdapter.initCache(requesterNfType, targetNfType, cm.indexGroup, nil, roamingCache)
}

//GetCacheStatus for getting cache status
func (cm *CacheManager) GetCacheStatus(requesterNfType string, targetNfType string) bool {
	homeAdapter := cm.mcache[requesterNfType]
	if homeAdapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return false
	}

	roamingAdapter := cm.roamingCache[requesterNfType]
	if roamingAdapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return homeAdapter.getCacheStatus(targetNfType) && roamingAdapter.getCacheStatus(targetNfType)
}

//SetCacheStatus for setting cache status
func (cm *CacheManager) SetCacheStatus(requesterNfType string, targetNfType string, status bool) {
	homeAdapter := cm.mcache[requesterNfType]
	if homeAdapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return
	}

	roamingAdapter := cm.roamingCache[requesterNfType]
	if roamingAdapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return
	}

	log.Debugf("set nftype[%s,%s] cache status to %t", requesterNfType, targetNfType, status)

	homeAdapter.setCacheStatus(targetNfType, status)
	roamingAdapter.setCacheStatus(targetNfType, status)
}

//EnterKeepCacheWorkMode enter keepCache work mode
func (cm *CacheManager) EnterKeepCacheWorkMode() bool {
	for _, requesterNfType := range nfTypes {
		if cm.mcache[requesterNfType] != nil {
			cm.mcache[requesterNfType].stopAllMonitorWorker()
			cm.mcache[requesterNfType].stopAllSubscrMonitorWorker()
		}

		if cm.roamingCache[requesterNfType] != nil {
			cm.roamingCache[requesterNfType].stopAllMonitorWorker()
			cm.roamingCache[requesterNfType].stopAllSubscrMonitorWorker()
		}
	}

	return true
}

//EnterNormalWorkMode enter normal work mode
func (cm *CacheManager) EnterNormalWorkMode() bool {
	//standby delay all profiles timeout(second)
	var standbyDelayTimeout uint = 3600
	for _, nfType := range nfTypes {
		targetNfs, ok := cm.GetTargetNfs(nfType)

		if !ok {
			log.Infof("Get targetNf for nfType[%s] fail, skip build cache for it", nfType)
			continue
		}
		for _, targetNf := range targetNfs {
			//(1)skip build new cache for offline nf
			if cm.GetCacheStatus(nfType, targetNf.TargetNfType) == false {
				log.Infof("Nf[%s,%s] no instance online, skip build new cache for it", nfType, targetNf.TargetNfType)
				continue
			}

			if election.IsActiveLeader(strconv.Itoa(pkgcm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
				//(2)build new cache
				//subscriptionMonitor := cm.getSubscriptionMonitor(nfType, targetNf.TargetNfType)
				//newCache := cm.buildCache(nfType, targetNf.TargetNfType, subscriptionMonitor)

				//(2)build new cache
				subscriptionMonitor := cm.getSubscriptionMonitor(nfType, targetNf.TargetNfType)
				newCache := cm.buildCache(nfType, targetNf.TargetNfType, subscriptionMonitor, homeCache)
				roamingSubscriptionMonitor := cm.getRoamingSubscriptionMonitor(nfType, targetNf.TargetNfType)
				newRoamingCache := cm.buildCache(nfType, targetNf.TargetNfType, roamingSubscriptionMonitor, roamingCache)

				//(3) build and reset ttlMonitor
				//newTtlMonitor := cm.buildTtlMonitor(nfType, targetNf.TargetNfType)
				//newTtlMonitor.startMonitorWorker()
				//newCache.setTtlMonitor(newTtlMonitor)

				newTtlMonitor := cm.buildTtlMonitor(nfType, targetNf.TargetNfType, homeCache)
				newTtlMonitor.startMonitorWorker()
				newCache.setTtlMonitor(newTtlMonitor)

				newRoamingTtlMonitor := cm.buildTtlMonitor(nfType, targetNf.TargetNfType, roamingCache)
				newRoamingTtlMonitor.startMonitorWorker()
				newRoamingCache.setTtlMonitor(newRoamingTtlMonitor)

				//(4)start fetch nfProfiles from NRF to newCache
				//(5)push data to MSB
				ok := SyncNrfData(&targetNf, true, newCache)
				if ok {
					log.Infof("Active pod sync profile data for targetNF[%v] from NRF-Disc success", targetNf)
				} else {
					log.Infof("Active pod sync profile data for targetNF[%v] from NRF-Disc failure", targetNf)
				}

				//(6)send DEREGISTERED message for decrease Nf profile
				cm.cacheProfileDiffHandler(nfType, targetNf.TargetNfType, newCache)

				//(7)switch to new cache
				cm.setCache(nfType, targetNf.TargetNfType, newCache)
				cm.setRoamingCache(nfType, targetNf.TargetNfType, newRoamingCache)
				cm.SetCacheStatus(nfType, targetNf.TargetNfType, true)

				//(8)start subscription ttlMonitor worker
				if cm.mcache[nfType] != nil {
					cm.mcache[nfType].startSubscrMonitorWorker(targetNf.TargetNfType)
				}
				if cm.roamingCache[nfType] != nil {
					cm.roamingCache[nfType].startSubscrMonitorWorker(targetNf.TargetNfType)
				}
			} else {
				//For slave, keep current cache by delay timeout,
				//and waiting for NFEventDiscResult and NFDeRegister from master to update cache

				//(2)delay all profile timeout
				cm.setSameTimeout(nfType, targetNf.TargetNfType, standbyDelayTimeout)

				//(3)start all monitor
				if cm.mcache[nfType] != nil {
					cm.mcache[nfType].startMonitorWorker(targetNf.TargetNfType)
					cm.mcache[nfType].startSubscrMonitorWorker(targetNf.TargetNfType)
				}
			}
		}
	}
	return true
}

func (cm *CacheManager) SuperviseSubscription(requesterNfType string, targetNfType string, subscriptionID string, timepoint time.Time) {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.superviseSubscription(targetNfType, subscriptionID, timepoint)
}

func (cm *CacheManager) SuperviseRoamingSubscription(requesterNfType string, targetNfType string, subscriptionID string, timepoint time.Time) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.superviseSubscription(targetNfType, subscriptionID, timepoint)
}

func (cm *CacheManager) getSubscriptionMonitor(requesterNfType string, targetNfType string) *subscriptionCache {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return nil
	}

	return adapter.getSubscriptionMonitor(targetNfType)
}

func (cm *CacheManager) getRoamingSubscriptionMonitor(requesterNfType string, targetNfType string) *subscriptionCache {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return nil
	}

	return adapter.getSubscriptionMonitor(targetNfType)
}

//DelSubscriptionMonitor delete subscriptionID ttl monitor
func (cm *CacheManager) DelSubscriptionMonitor(requesterNfType string, targetNfType string, subscriptionID string) {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.deleteSubscriptionMonitor(targetNfType, subscriptionID)
}

func (cm *CacheManager) DelRoamingSubscriptionMonitor(requesterNfType string, targetNfType string, subscriptionID string) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.deleteSubscriptionMonitor(targetNfType, subscriptionID)
}

//Probe cachemanager probe cache whether have cached the nfintanceid profile
func (cm *CacheManager) Probe(requesterNfType string, targetNfType string, nfInstanceID string) bool {
	adapter := cm.getCacheAdapter(requesterNfType, false)
	if adapter == nil {
		return false
	}
	exist := adapter.probe(targetNfType, nfInstanceID)
	if !exist {
		log.Infof("NRF-Disc-Agent did not cached the profile nfInstanceID[%s]", nfInstanceID)
	}
	return exist
}

//Probe cachemanager probe cache whether have cached the nfintanceid profile
func (cm *CacheManager) ProbeRoam(requesterNfType string, targetNfType string, nfInstanceID string) bool {
	adapter := cm.getCacheAdapter(requesterNfType, true)
	if adapter == nil {
		return false
	}
	exist := adapter.probe(targetNfType, nfInstanceID)
	if !exist {
		log.Infof("NRF-Disc-Agent did not roam cached the profile nfInstanceID[%s]", nfInstanceID)
	}
	return exist
}

//ProbeAllCache cachemanager probe normal/roaming cache whether have cached the nfintanceid profile
//ProbeAllCache return value (exist, isRoam)
func (cm *CacheManager) ProbeAllCache(requesterNfType string, targetNfType string, nfInstanceID string) (bool, bool) {
	normalExist := cm.Probe(requesterNfType, targetNfType, nfInstanceID)
	if normalExist {
		return true, false
	}
	roamExist := cm.ProbeRoam(requesterNfType, targetNfType, nfInstanceID)
	if roamExist {
		return true, true
	}
	return false, false
}

//Probe cachemanager probe cache whether have cached the nfintanceid profile
func (cm *CacheManager) GetProfileByID(requesterNfType string, targetNfType string, nfInstanceID string, isRoam bool) []byte {
	adapter := cm.getCacheAdapter(requesterNfType, isRoam)
	if adapter == nil {
		return nil
	}

	profile := adapter.getProfileByID(targetNfType, nfInstanceID)
	if profile == nil {
		log.Infof("No such profile by nfProfileID:%s from requesterNfType:%s targetNfType:%s", nfInstanceID, requesterNfType, targetNfType)
	}

	return profile
}

//Cached cache the profile, have discard this function
func (cm *CacheManager) Cached(requesterNfType string, targetNfType string, content []byte, isRoaming bool) {
	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return
	}

	_, ok := adapter.indexed(targetNfType, content, cm.indexGroup)
	if !ok {
		log.Errorf("no index item in NF profile %s", string(content))

	} else {
		adapter.showIndexContent(targetNfType)
		adapter.cached(targetNfType, content)
	}
}

//CachedWithTTL cache the profile with ttl
func (cm *CacheManager) CachedWithTTL(requesterNfType string, targetNfType string, content []byte, ttl uint, isRoaming bool) {
	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return
	}

	nfInstanceID, ok := adapter.indexed(targetNfType, content, cm.indexGroup)
	if !ok {
		log.Errorf("no index item in NF profile %+v", string(content))
	} else {
		adapter.showIndexContent(targetNfType)
		adapter.cached(targetNfType, content)
		adapter.supervise(targetNfType, nfInstanceID, ttl)

		if !cm.GetCacheStatus(requesterNfType, targetNfType) {
			cm.SetCacheStatus(requesterNfType, targetNfType, true)
		}
	}
	log.Debugf("CachedWithTTL nfType[%s,%s], ttl=%d, isRoaming=%t", requesterNfType, targetNfType, ttl, isRoaming)
}

//CachedWithTtlTimeStamp cache the profile with ttl timestamp
func (cm *CacheManager) CachedWithTtlTimestamp(requesterNfType string, targetNfType string, content []byte, ttl time.Time, isRoaming bool) {
	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return
	}

	nfInstanceID, ok := adapter.indexed(targetNfType, content, cm.indexGroup)
	if !ok {
		log.Errorf("no index item in NF profile %+v", string(content))
	} else {
		adapter.showIndexContent(targetNfType)
		adapter.cached(targetNfType, content)
		adapter.superviseTimestamp(targetNfType, nfInstanceID, ttl)
		if !cm.GetCacheStatus(requesterNfType, targetNfType) {
			cm.SetCacheStatus(requesterNfType, targetNfType, true)
		}
	}
}

//InjectionCachedWithTtl inject content to cache
func (cm *CacheManager) InjectionCachedWithTtl(cache *cache, content []byte, ttl uint) {
	if cache == nil {
		log.Errorf("cache ptr is nil")
		return
	}
	nfInstanceID, ok := cache.indexed(content, cm.indexGroup)
	if !ok {
		log.Errorf("no index item in NF profile %+v for cache", string(content))
	} else {
		cache.cached(content)
		ttlMonitor := cache.getTtlMonitor()
		if ttlMonitor == nil {
			log.Warn("Cache ttlMonitor is nil")
			return
		}
		ttlMonitor.supervise(nfInstanceID, ttl)
	}
}

//setSameTimeout is for set same timeout value for all profiles in cache
func (cm *CacheManager) setSameTimeout(requesterNfType string, targetNfType string, ttl uint) {
	if cm.mcache[requesterNfType] != nil {
		cm.mcache[requesterNfType].superviseAll(targetNfType, ttl)
	}
}

//DeCached delete cache by nfInstanceID
func (cm *CacheManager) DeCached(requesterNfType string, targetNfType string, nfInstanceID string, isRoaming bool) {
	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return
	}

	adapter.deIndex(targetNfType, nfInstanceID)
	adapter.deCached(targetNfType, nfInstanceID)
	adapter.stopByID(targetNfType, nfInstanceID)
	adapter.deleteEtag(targetNfType, nfInstanceID)
}

/*
//DeRoamingCached delete roaming cache by nfInstanceID
func (cm *CacheManager) DeRoamingCached(requesterNfType string, targetNfType string, nfInstanceID string) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.deIndex(targetNfType, nfInstanceID)
	adapter.deCached(targetNfType, nfInstanceID)
	adapter.stopByID(targetNfType, nfInstanceID)
	adapter.deleteEtag(targetNfType, nfInstanceID)
}
*/

func (cm *CacheManager) DeCachedByNfType(requesterNfType string) {
	targetNfs, ok := cm.GetTargetNfs(requesterNfType)
	if !ok {
		log.Warnf("Get targetNf for nfType[%s] fail, skip build cache for it", requesterNfType)
		return
	}

	adapter := cm.getCacheAdapter(requesterNfType, false)
	if adapter == nil {
		return
	}

	for _, targetNf := range targetNfs {
		adapter.deIndexByNfType(targetNf.TargetNfType)
		adapter.deCachedByNfType(targetNf.TargetNfType)
		adapter.deleteAllEtag(targetNf.TargetNfType)
	}
	adapter.stopAllTtlMonitor()

	adapter = cm.getCacheAdapter(requesterNfType, true)
	if adapter == nil {
		return
	}

	for _, targetNf := range targetNfs {
		adapter.deIndexByNfType(targetNf.TargetNfType)
		adapter.deCachedByNfType(targetNf.TargetNfType)
		adapter.deleteAllEtag(targetNf.TargetNfType)
	}
	adapter.stopAllTtlMonitor()
}

/*
func (cm *CacheManager) deCachedByNfType(requesterNfType string) {
	targetNfs, ok := cm.GetTargetNfs(requesterNfType)
	if !ok {
		log.Warnf("Get targetNf for nfType[%s] fail, skip build cache for it", requesterNfType)
		return
	}

	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	for _, targetNf := range targetNfs {
		adapter.deIndexByNfType(targetNf.TargetNfType)
		adapter.deCachedByNfType(targetNf.TargetNfType)
		adapter.deleteAllEtag(targetNf.TargetNfType)
	}
	adapter.stopAllTtlMonitor()
}

func (cm *CacheManager) deRoamingCachedByNfType(requesterNfType string) {
	targetNfs, ok := cm.GetTargetNfs(requesterNfType)
	if !ok {
		log.Warnf("Get targetNf for nfType[%s] fail, skip build cache for it", requesterNfType)
		return
	}

	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	for _, targetNf := range targetNfs {
		adapter.deIndexByNfType(targetNf.TargetNfType)
		adapter.deCachedByNfType(targetNf.TargetNfType)
		adapter.deleteAllEtag(targetNf.TargetNfType)
	}
	adapter.stopAllTtlMonitor()
}
*/
////DeCachedAll delete cache by nfInstanceID
//func (cm *CacheManager) DeCachedAll(nfInstanceID string) {
//	for _, requesterNfType := range cm.mcache.getNfTypes() {
//		cm.mcache.deIndex(requesterNfType, nfInstanceID)
//		cm.mcache.deCached(requesterNfType, nfInstanceID)
//		cm.ttlMonitor.stopByID(requesterNfType, nfInstanceID)
//	}
//}

//ReCached delete and then save in cache
func (cm *CacheManager) ReCached(requesterNfType string, targetNfType string, nfInstanceID string, content []byte, isRoaming bool) {
	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return
	}

	adapter.deIndex(targetNfType, nfInstanceID)
	adapter.deCached(targetNfType, nfInstanceID)
	_, ok := adapter.indexed(targetNfType, content, cm.indexGroup)
	if !ok {
		log.Errorf("no index item in NF profile %s", string(content))
	} else {
		adapter.showIndexContent(targetNfType)
		adapter.cached(targetNfType, content)
	}
}

//Search search in cache
func (cm *CacheManager) Search(requesterNfType string, targetNfType string, searchParameter *SearchParameter, isKeepCacheMode bool) ([]byte, bool) {
	cache := cm.getCache(requesterNfType, targetNfType)
	if cache == nil {
		log.Errorf("CacheManager have not create %s requesterNfType cache", requesterNfType)
		return nil, false
	}
	indexSearcher := indexSearcher{
		cache,
		searchParameter,
		cm.indexMapper,
	}

	ids, result := indexSearcher.search()
	if !result {
		log.Debugf("cache index search miss")
		return nil, false
	}

	result = searchParameter.ProfileSearchNecessary()
	if result {
		profileSearcher := profileSearcher{
			cache,
			ids,
			searchParameter,
		}
		ids, result = profileSearcher.search()
		if !result {
			log.Debugf("cache profile search miss")
			return nil, false
		}
	}

	/////////assemble response contents//////////
	var nfProfileList []string
	ids.Each(func(item interface{}) bool {
		switch t := item.(type) {
		case string:
			nfProfileList = append(nfProfileList, t)
		default:
			log.Errorf("unsupported interface type")
		}
		return false
	})

	contents := cm.assembleResponseContents(requesterNfType, targetNfType, nfProfileList, searchParameter, isKeepCacheMode, false)
	if contents == nil {
		log.Errorf("cache assemble search result error")
		return nil, false
	}

	return contents, true
}

//SearchRoamingCache search in roamingCache
func (cm *CacheManager) SearchRoamingCache(requesterNfType string, targetNfType string, searchParameter *SearchParameter, isKeepCacheMode bool) ([]byte, bool) {
	roamingCache := cm.getRoamingCache(requesterNfType, targetNfType)
	if roamingCache == nil {
		log.Errorf("CacheManager have not create %s requesterNfType roamingCache", requesterNfType)
		return nil, false
	}
	indexSearcher := indexSearcher{
		roamingCache,
		searchParameter,
		cm.indexMapper,
	}

	ids, result := indexSearcher.search()
	if !result {
		log.Debugf("cache index search miss")
		return nil, false
	}

	result = searchParameter.ProfileSearchNecessary()
	if result {
		profileSearcher := profileSearcher{
			roamingCache,
			ids,
			searchParameter,
		}
		ids, result = profileSearcher.search()
		if !result {
			log.Debugf("cache profile search miss")
			return nil, false
		}
	}

	/////////assemble response contents//////////
	var nfProfileList []string
	ids.Each(func(item interface{}) bool {
		switch t := item.(type) {
		case string:
			nfProfileList = append(nfProfileList, t)
		default:
			log.Errorf("unsupported interface type")
		}
		return false
	})

	contents := cm.assembleResponseContents(requesterNfType, targetNfType, nfProfileList, searchParameter, isKeepCacheMode, true)
	if contents == nil {
		log.Errorf("cache assemble search result error")
		return nil, false
	}

	return contents, true
}

//Flush flush the mcache content by requesterNfType
func (cm *CacheManager) Flush(requesterNfType string) {
	homeAdapter := cm.mcache[requesterNfType]
	if homeAdapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
	} else {
		for _, targetNfType := range homeAdapter.getTargetNfTypes() {
			homeAdapter.flush(targetNfType)
		}
		homeAdapter.stopAllTtlMonitor()
	}
}

//Flush flush the roamcache content by requesterNfType
func (cm *CacheManager) FlushRoam(requesterNfType string) {
	roamingAdapter := cm.roamingCache[requesterNfType]
	if roamingAdapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
	} else {
		for _, targetNfType := range roamingAdapter.getTargetNfTypes() {
			roamingAdapter.flush(targetNfType)
		}
		roamingAdapter.stopAllTtlMonitor()
	}
}

//FlushAll flush all the cache content
func (cm *CacheManager) FlushAll() {
	requesterNfTypes := cm.getRequesterNfTypes()
	for _, requesterNfType := range requesterNfTypes {
		cm.Flush(requesterNfType)
		//cm.mcache[requesterNfType].flushAll()
		//cm.mcache[requesterNfType].stopAllTtlMonitor()
	}

	requesterRoamNfTypes := cm.getRequesterRoamNfTypes()
	for _, requesterRoamNfType := range requesterRoamNfTypes {
		cm.FlushRoam(requesterRoamNfType)
	}
}

/*
func (cm *CacheManager) requesterNfTypeCacheExist(requesterNftype string) bool {
	requesterNfTypes := cm.getRequesterNfTypes()
	for _, nfType := range requesterNfTypes {
		if requesterNftype == nfType {
			return true
		}
	}
	return false
}
*/

//Sync is for sync cache content by requesterNfType
func (cm *CacheManager) Sync(requesterNfType string, syncData *structs.CacheSyncData) {
	syncData.RequestNfType = requesterNfType

	homeAdapter := cm.mcache[requesterNfType]
	if homeAdapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
	} else {
		homeAdapter.sync(syncData)
	}

	roamingAdapter := cm.roamingCache[requesterNfType]
	if roamingAdapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
	} else {
		roamingAdapter.sync(syncData)
	}
}

//Dump dump cache content by requesterNfType
func (cm *CacheManager) Dump(requesterNfType string, dumpData *structs.CacheDumpData) {
	dumpData.RequestNfType = requesterNfType

	homeAdapter := cm.mcache[requesterNfType]
	if homeAdapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
	} else {
		homeAdapter.dump(dumpData)
	}

	roamingAdapter := cm.roamingCache[requesterNfType]
	if roamingAdapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
	} else {
		roamingAdapter.dump(dumpData)
	}
}

//DumpAll dump all the cache content
func (cm *CacheManager) DumpAll(dumpDatas *[]structs.CacheDumpData) {
	for _, requesterNfType := range cm.getRequesterNfTypes() {
		dumpData := structs.CacheDumpData{
			RequestNfType: requesterNfType,
		}
		cm.Dump(requesterNfType, &dumpData)
		if len(dumpData.CacheInfos) == 0 && len(dumpData.RoamingCacheInfos) == 0 {
			continue
		}
		//dumpDatas = append(dumpDatas, dumpData)
		*dumpDatas = append(*dumpDatas, dumpData)
	}
}

//GetNfProfile fetch nfProfile by requesterNfType and targetNfType
func (cm *CacheManager) GetNfProfile(requesterNfType string, targetNfType string, nfInstanceID string) []byte {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil
	}

	cache := adapter.getCache(targetNfType)
	if cache == nil {
		log.Infof("CacheManager have not create home cache targetCache:%s for nfType:%s", targetNfType, requesterNfType)
		return nil
	}

	nfProfile := cache.fetchProfileByID(nfInstanceID)
	return nfProfile

}

//GetRoamingNfProfile fetch nfProfile by requesterNfType and targetNfType
func (cm *CacheManager) GetRoamingNfProfile(requesterNfType string, targetNfType string, nfInstanceID string) []byte {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil
	}

	cache := adapter.getCache(targetNfType)
	if cache == nil {
		log.Infof("CacheManager have not create roaming cache targetCache:%s for nfType:%s", targetNfType, requesterNfType)
		return nil
	}

	nfProfile := cache.fetchProfileByID(nfInstanceID)
	return nfProfile

}

//FetchNfProfile fetch nfProfile by requesterNfType
func (cm *CacheManager) FetchNfProfile(requesterNfType string, nfInstanceID string) []byte {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil
	}

	for _, targetNfType := range adapter.getTargetNfTypes() {
		cache := adapter.getCache(targetNfType)
		if cache == nil {
			log.Infof("CacheManager have not create home cache targetCache:%s for nfType:%s", targetNfType, requesterNfType)
			return nil
		}
		nfProfile := cache.fetchProfileByID(nfInstanceID)
		if nfProfile != nil {
			return nfProfile
		}
	}

	return nil
}

//FetchRoamingNfProfile fetch nfProfile by requesterNfType
func (cm *CacheManager) FetchRoamingNfProfile(requesterNfType string, nfInstanceID string) []byte {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil
	}

	for _, targetNfType := range adapter.getTargetNfTypes() {
		cache := adapter.getCache(targetNfType)
		if cache == nil {
			log.Infof("CacheManager have not create roaming cache targetCache:%s for nfType:%s", targetNfType, requesterNfType)
			return nil
		}
		nfProfile := cache.fetchProfileByID(nfInstanceID)
		if nfProfile != nil {
			return nfProfile
		}
	}

	return nil
}

//FetchProfileIDs get all the cache profile IDs
func (cm *CacheManager) FetchProfileIDs(requesterNfType string) []string {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.fetchAllProfileIDs()
}

//FetchProfileIDs get all the cache profile IDs
func (cm *CacheManager) FetchRoamingProfileIDs(requesterNfType string) []string {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.fetchAllProfileIDs()
}

//FetchAllProfileIDs get all the cache profile IDs
func (cm *CacheManager) FetchAllProfileIDs() []string {
	ids := make([]string, 0)
	for _, requesterNfType := range cm.getRequesterNfTypes() {
		adapter := cm.mcache[requesterNfType]
		if adapter == nil {
			log.Warnf("No such home cache for nfType:%s", requesterNfType)
			continue
		}
		id := adapter.fetchAllProfileIDs()
		ids = append(ids, id...)
	}
	return ids
}

//FetchAllRoamingProfileIDs get all the cache profile IDs
func (cm *CacheManager) FetchAllRoamingProfileIDs() []string {
	ids := make([]string, 0)
	for _, requesterNfType := range cm.getRequesterNfTypes() {
		adapter := cm.roamingCache[requesterNfType]
		if adapter == nil {
			log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
			continue
		}
		id := adapter.fetchAllProfileIDs()
		ids = append(ids, id...)
	}
	return ids
}

//DumpByID dump cache content by requesterNfType and nfInstanceID
func (cm *CacheManager) DumpByID(requesterNfType string, targetNfType string, nfInstanceID string) []byte {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.dumpByID(targetNfType, nfInstanceID)
}

//DumpRoamingByID dump cache content by requesterNfType and nfInstanceID
func (cm *CacheManager) DumpRoamingByID(requesterNfType string, targetNfType string, nfInstanceID string) []byte {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.dumpByID(targetNfType, nfInstanceID)
}

func (cm *CacheManager) HaveEtag(requesterNfType string, targetNfType string, nfInstanceID string) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.haveEtag(targetNfType, nfInstanceID)
}

func (cm *CacheManager) HaveRoamingEtag(requesterNfType string, targetNfType string, nfInstanceID string) bool {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.haveEtag(targetNfType, nfInstanceID)
}

func (cm *CacheManager) FetchEtag(requesterNfType string, targetNfType string, nfInstanceID string) string {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return ""
	}

	return adapter.fetchEtag(targetNfType, nfInstanceID)
}

func (cm *CacheManager) FetchRoamingEtag(requesterNfType string, targetNfType string, nfInstanceID string) string {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return ""
	}

	return adapter.fetchEtag(targetNfType, nfInstanceID)
}

//SaveEtag save roaming etag
func (cm *CacheManager) SaveEtag(requesterNfType string, targetNfType string, nfInstanceID string, value string) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.saveEtag(targetNfType, nfInstanceID, value)
}

//SaveRoamingEtag save roaming etag
func (cm *CacheManager) SaveRoamingEtag(requesterNfType string, targetNfType string, nfInstanceID string, value string) bool {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.saveEtag(targetNfType, nfInstanceID, value)
}

//AddSubsIDWithNfType is for adding subscribeId URL cache with requesterNfType and targetNfType
/*
func (cm *CacheManager) AddSubscriptionID(requesterNfType string, targetNfType string, subsIdURL string) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.addSubscriptionID(targetNfType, subsIdURL)
}

//AddRoamingSubscriptionID is for adding subscribeId URL cache with requesterNfType and targetNfType
func (cm *CacheManager) AddRoamingSubscriptionID(requesterNfType string, targetNfType string, subsIdURL string) bool {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.addSubscriptionID(targetNfType, subsIdURL)
}
*/

func (cm *CacheManager) AddSubscriptionInfo(requesterNfType string, targetNfType string, subscriptionInfo structs.SubscriptionInfo) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return false
	}

	return adapter.addSubscriptionInfo(targetNfType, subscriptionInfo)
}

func (cm *CacheManager) AddRoamingSubscriptionInfo(requesterNfType string, targetNfType string, subscriptionInfo structs.SubscriptionInfo) bool {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return false
	}

	return adapter.addSubscriptionInfo(targetNfType, subscriptionInfo)
}

/*
//DelSubsIDWithNfType is for delete subscribeId URL cache with requesterNfType
func (cm *CacheManager) DelSubscriptionID(requesterNfType string, targetNfType string, subsIdURL string) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.delSubscriptionID(targetNfType, subsIdURL)
}
*/

/*
//DelRoamingSubscriptionID is for delete subscribeId URL cache with requesterNfType
func (cm *CacheManager) DelRoamingSubscriptionID(requesterNfType string, targetNfType string, subsIdURL string) bool {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.delSubscriptionID(targetNfType, subsIdURL)
}
*/

func (cm *CacheManager) DelSubscriptionInfo(requesterNfType string, targetNfType string, subscriptionID string) {
	adapter := cm.getCacheAdapter(requesterNfType, false)
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.delSubscriptionInfo(targetNfType, subscriptionID)
}

func (cm *CacheManager) DelRoamingSubscriptionInfo(requesterNfType string, targetNfType string, subscriptionID string) {
	adapter := cm.getCacheAdapter(requesterNfType, true)
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.delSubscriptionInfo(targetNfType, subscriptionID)
}

//GetSubscriptionIDURLs is for get subscriptionID URL set for requester-NfType
func (cm *CacheManager) GetSubscriptionIDs(requesterNfType string, targetNfType string) ([]string, bool) {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil, false
	}

	return adapter.getSubscriptionIDs(targetNfType)
}

func (cm *CacheManager) GetNfProfileSubscriptionID(requesterNfType, targetNfType, nfInstanceID string) (string, bool) {
	adapter := cm.getCacheAdapter(requesterNfType, true)
	if adapter == nil {
		return "", false
	}

	return adapter.getNfProfileSubscriptionID(targetNfType, nfInstanceID)
}

//GetRoamingSubscriptionIDs is for get subscriptionID URL set for requester-NfType
func (cm *CacheManager) GetRoamingSubscriptionIDs(requesterNfType string, targetNfType string) ([]string, bool) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil, false
	}

	return adapter.getSubscriptionIDs(targetNfType)
}

//GetServiceSubscriptionInfo get subscriptionInfo via service
func (cm *CacheManager) GetServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName string) (structs.SubscriptionInfo, bool) {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return structs.SubscriptionInfo{}, false
	}

	return adapter.getServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName)
}

func (cm *CacheManager) GetNfProfileSubscriptionInfo(requesterNfType, targetNfType, nfInstanceID string) (structs.SubscriptionInfo, bool) {
	adapter := cm.getCacheAdapter(requesterNfType, true)
	if adapter == nil {
		return structs.SubscriptionInfo{}, false
	}

	return adapter.getNfProfileSubscriptionInfo(targetNfType, nfInstanceID)
}

//GetRoamingServiceSubscriptionInfo get subscriptionInfo via service
func (cm *CacheManager) GetRoamingServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName string) (structs.SubscriptionInfo, bool) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return structs.SubscriptionInfo{}, false
	}

	return adapter.getServiceSubscriptionInfo(requesterNfType, targetNfType, serviceName)
}

//UpdateSubscriptionStorage update subscriptionInfo to storage configmap
func (cm *CacheManager) UpdateSubscriptionStorage() bool {
	subscriptionInfoMap := cm.getAllSubscriptionInfo()
	rest := updateConfigmapStorage(subscriptionInfoMap)

	return rest
}

//ProbeSubscriptionInfo probe subscribe info
func (cm *CacheManager) ProbeSubscriptionInfo(requesterNfType, targetNfType, serviceName string) bool {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return false
	}

	return adapter.probeSubscriptionInfo(targetNfType, serviceName)
}

//SubscriptionInfoProvision provision subscriptionInfo from storage configmap to cacheManager
func (cm *CacheManager) SubscriptionInfoProvision(subscriptionInfoData []byte) bool {
	if len(subscriptionInfoData) == 0 {
		return false
	}

	updated := false
	subscriptionInfoContainer := make(map[string]structs.SubscriptionInfo, 0)
	err := json.Unmarshal(subscriptionInfoData, &subscriptionInfoContainer)
	if err == nil {
		for k, v := range subscriptionInfoContainer {
			adapter := cm.getCacheAdapter(v.RequesterNfType, false)
			if adapter == nil {
				return false
			}

			if v.ValidityTime.Before(time.Now()) {
				log.Warnf("delete expired subscription(%s) from configmap storage, %+v", k, v)
				adapter.delSubscriptionInfo(v.TargetNfType, v.SubscriptionID)
				delete(subscriptionInfoContainer, v.SubscriptionID)
				updated = true
				continue
			}

			adapter.addSubscriptionInfo(v.TargetNfType, v)
			adapter.superviseSubscription(v.TargetNfType, v.SubscriptionID, v.ValidityTime)
		}
	} else {
		updated = true
	}

	if updated && election.IsActiveLeader(strconv.Itoa(pkgcm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		subscriptionInfoDataNew, err := json.Marshal(subscriptionInfoContainer)
		if err != nil {
			log.Errorf("Marshal subscriptionInfo[%+v] fail, err:%s", subscriptionInfoContainer, err.Error())
			return false
		}
		err = k8sapiclient.SetConfigMapData(consts.ConfigMapStorage, consts.ConfigMapKeySubsInfo, subscriptionInfoDataNew)
		if err != nil {
			log.Errorf("Write subscriptionInfo to configmap[%s] fail, err:%s", consts.ConfigMapStorage, err.Error())
			return false
		}

		log.Debugf("Provision subscriptionInfo[%s] to configmap storage success", string(subscriptionInfoDataNew))
	}

	return true
}

//SetRequesterFqdn for set requester fqdn
func (cm *CacheManager) SetRequesterFqdn(nfType string, fqdn string) bool {
	cm.requesterFqdnMutex.Lock()
	defer cm.requesterFqdnMutex.Unlock()

	fqdnOld, ok := cm.requesterFqdn[nfType]
	if ok {
		if fqdn == fqdnOld {
			log.Debugf("CacheManager exist the same fqdn for nfType[%s], skip save action", nfType)
			return true
		} else {
			log.Debugf("The nfType[%s] fqdn need to be updated from [%s] to [%s]", nfType, fqdnOld, fqdn)
		}
	}

	cm.requesterFqdn[nfType] = fqdn

	_, ok = cm.requesterFqdn[nfType]
	if !ok {
		log.Debugf("Save nfType[%s],fqdn[%s] fail", nfType, fqdn)
		return false
	}

	log.Debugf("Save nfType[%s],fqdn[%s] success", nfType, fqdn)

	return true
}

//GetRequesterFqdn for get requester fqdn
func (cm *CacheManager) GetRequesterFqdn(nfType string) (string, bool) {
	cm.requesterFqdnMutex.Lock()
	defer cm.requesterFqdnMutex.Unlock()

	fqdn, ok := cm.requesterFqdn[nfType]
	if !ok {
		log.Debugf("Get fqdn by nfType[%s] fail", nfType)
		return "", false
	}
	log.Debugf("Get fqdn by nfType[%s] success, fqdn[%s]", nfType, fqdn)

	return fqdn, true
}

//DeleteRequesterFqdn for delete requester fqdn
func (cm *CacheManager) DeleteRequesterFqdn(nfType string) bool {
	cm.requesterFqdnMutex.Lock()
	defer cm.requesterFqdnMutex.Unlock()

	_, ok := cm.requesterFqdn[nfType]
	if !ok {
		log.Debugf("CacheManager not such fqdn for nfType[%s], skip delete action", nfType)
		return true
	}

	delete(cm.requesterFqdn, nfType)

	_, ok = cm.requesterFqdn[nfType]
	if ok {
		log.Debugf("Del fqdn by nfType[%s] fail", nfType)
		return false
	}

	log.Debugf("Del fqdn by nfType[%s] success", nfType)
	return true
}

//SetRequesterPlmns for set requester plmns
func (cm *CacheManager) SetRequesterPlmns(nfType string, plmns []structs.PlmnID) bool {
	cm.requesterPlmnsMutex.Lock()
	defer cm.requesterPlmnsMutex.Unlock()
	/*
		plmnsOld, ok := cm.requesterPlmns[nfType]
		if ok {
			if fqdn == fqdnOld {
				log.Debugf("CacheManager exist the same fqdn for nfType[%s], skip save action", nfType)
				return true
			} else {
				log.Debugf("The nfType[%s] fqdn need to be updated from [%s] to [%s]", nfType, fqdnOld, fqdn)
			}
		}
	*/
	cm.requesterPlmns[nfType] = plmns

	_, ok := cm.requesterPlmns[nfType]
	if !ok {
		log.Debugf("Set plmns[%v] by nfType[%s] fail", plmns, nfType)
		return false
	}
	log.Debugf("Set plmns[%s] by nfType[%s] success", plmns, nfType)

	return true
}

//GetRequesterPlmns for get requester plmns
func (cm *CacheManager) GetRequesterPlmns(nfType string) ([]structs.PlmnID, bool) {
	cm.requesterPlmnsMutex.Lock()
	defer cm.requesterPlmnsMutex.Unlock()

	plmns, ok := cm.requesterPlmns[nfType]
	if !ok {
		log.Debugf("Get plmns by nfType[%s] fail", nfType)
		return nil, false
	}
	log.Debugf("Get plmns by nfType[%s] success, plmns[%s]", nfType, plmns)

	return plmns, true
}

//DelRequesterPlmns for delete requester plmns
func (cm *CacheManager) DelRequesterPlmns(nfType string) bool {
	cm.requesterPlmnsMutex.Lock()
	defer cm.requesterPlmnsMutex.Unlock()

	_, ok := cm.requesterPlmns[nfType]
	if !ok {
		log.Debugf("No such plmns for nfType[%s], skip delete action", nfType)
		return true
	}

	delete(cm.requesterPlmns, nfType)

	_, ok = cm.requesterPlmns[nfType]
	if ok {
		log.Debugf("Del plmns by nfType[%s] fail", nfType)
		return false
	}
	log.Debugf("Del plmns by nfType[%s] success", nfType)

	return true
}

//SetTargetNf for cache targetNF from ConfigMap
func (cm *CacheManager) SetTargetNf(nfType string, targetNf structs.TargetNf) bool {
	cm.targetNfMutex.Lock()
	defer cm.targetNfMutex.Unlock()

	exist := false
	index := 0
	for i, loopTargetNf := range cm.targetNf[nfType] {
		if loopTargetNf.RequesterNfType == targetNf.RequesterNfType && loopTargetNf.TargetNfType == targetNf.TargetNfType {
			exist = true
			index = i
		}
	}
	if !exist {
		cm.targetNf[nfType] = append(cm.targetNf[nfType], targetNf)
	} else {
		cm.targetNf[nfType][index].TargetServiceNames = targetNf.TargetServiceNames
		cm.targetNf[nfType][index].NotifCondition = targetNf.NotifCondition
		cm.targetNf[nfType][index].SubscriptionValidityTime = targetNf.SubscriptionValidityTime
	}
	log.Debugf("Set targetNf for nfType[%s] success, targetNf[%s]", nfType, targetNf.Info())

	return true
}

//GetTargetNfs for get all targetNFs from cache
func (cm *CacheManager) GetTargetNfs(nfType string) ([]structs.TargetNf, bool) {
	cm.targetNfMutex.Lock()
	defer cm.targetNfMutex.Unlock()

	targetNFs, ok := cm.targetNf[nfType]
	if !ok {
		log.Debugf("Get targetNfs for nfType[%s] fail", nfType)
		return nil, false
	}

	log.Debugf("Get targetNfs by nfType[%s] success, targetNfs[%v]", nfType, targetNFs)
	return targetNFs, true
}

//GetTargetNf for get targetNFs from cache
func (cm *CacheManager) GetTargetNf(nfType, targetNfType string) (structs.TargetNf, bool) {
	cm.targetNfMutex.Lock()
	defer cm.targetNfMutex.Unlock()

	targetNFs, ok := cm.targetNf[nfType]
	if ok {
		for _, targetNF := range targetNFs {
			if targetNF.TargetNfType == targetNfType {
				log.Debugf("Get targetNf by nfType[%s] and targetNfType[%s] success, targetNf[%s]", nfType, targetNfType, targetNF.Info())
				return targetNF, true
			}
		}
	}

	log.Debugf("Get targetNf by nfType[%s] and targetNfType[%s] fail", nfType, targetNfType)
	return structs.TargetNf{}, false
}

//GetAllTargetNf for get all targetNFs
func (cm *CacheManager) GetAllTargetNf() map[string][]structs.TargetNf {
	cm.targetNfMutex.Lock()
	defer cm.targetNfMutex.Unlock()

	return cm.targetNf
}

//DeleteTargetNf for delete targetNF in cache
func (cm *CacheManager) DeleteTargetNf(nfType string) bool {
	cm.targetNfMutex.Lock()
	defer cm.targetNfMutex.Unlock()

	delete(cm.targetNf, nfType)
	log.Debugf("Delete targetNfs by nfType[%s] success", nfType)

	return true
}

//////////////////////private function////////////////////////

func (cm *CacheManager) init() bool {
	cm.indexGroup = make([]string, 0)
	cm.mcache = make(map[string]*cacheAdapter)
	cm.roamingCache = make(map[string]*cacheAdapter)
	cm.indexMapper = searchIndexMapper{}
	ok := cm.loadConfig()
	if !ok {
		log.Errorf("CacheManager load config file failure")
		return false
	}
	for _, requestNfType := range nfTypes {
		cm.mcache[requestNfType] = new(cacheAdapter)
		cm.mcache[requestNfType].init(homeCache)
		cm.roamingCache[requestNfType] = new(cacheAdapter)
		cm.roamingCache[requestNfType].init(roamingCache)
	}
	cm.targetNf = make(map[string][]structs.TargetNf)
	cm.requesterFqdn = make(map[string]string)
	cm.requesterPlmns = make(map[string][]structs.PlmnID)

	return true
}

func (cm *CacheManager) getCacheAdapter(requesterNfType string, isRoaming bool) *cacheAdapter {
	var adapter *cacheAdapter
	if isRoaming {
		adapter = cm.roamingCache[requesterNfType]
	} else {
		adapter = cm.mcache[requesterNfType]
	}

	if adapter == nil {
		if isRoaming {
			log.Warnf("No such roaming cacheAdapter for requesterNfType:%s", requesterNfType)
		} else {
			log.Warnf("No such home cacheAdapter for requesterNfType:%s", requesterNfType)
		}
	}

	return adapter
}

func (cm *CacheManager) getCache(requesterNfType string, targetNfType string) *cache {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such home cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.getCache(targetNfType)
}

func (cm *CacheManager) getRoamingCache(requesterNfType string, targetNfType string) *cache {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such roaming cache for nfType:%s", requesterNfType)
		return nil
	}

	return adapter.getCache(targetNfType)
}

func (cm *CacheManager) buildCache(requesterNfType, targetNfType string, subscriptionMonitor *subscriptionCache, mcacheType cacheType) *cache {
	cache := new(cache)
	cache.init(requesterNfType, targetNfType, cm.indexGroup, subscriptionMonitor, mcacheType)
	log.Info("Build cache success")

	return cache
}

func (cm *CacheManager) buildTtlMonitor(requesterNfType, targetNfType string, mcacheType cacheType) *ttlMonitor {
	ttlMonitor := new(ttlMonitor)
	ttlMonitor.init(requesterNfType, targetNfType, mcacheType)
	log.Info("Build ttlMonitor success")

	return ttlMonitor
}

func (cm *CacheManager) setCache(requesterNfType string, targetNfType string, cache *cache) {
	adapter := cm.mcache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.cacheMap[targetNfType] = cache
	log.Infof("Set home profile cache for req nfType[%s], target nfType[%s] success", requesterNfType, targetNfType)
}

func (cm *CacheManager) setRoamingCache(requesterNfType string, targetNfType string, cache *cache) {
	adapter := cm.roamingCache[requesterNfType]
	if adapter == nil {
		log.Warnf("No such cache from requesterNfType:%s", requesterNfType)
		return
	}

	adapter.cacheMap[targetNfType] = cache
	log.Infof("Set roaming profile cache for req nfType[%s], target nfType[%s] success", requesterNfType, targetNfType)
}

/*
func (cm *CacheManager) setTtlMonitor(requesterNfType string, targetNfType string, ttlMonitor *ttlMonitor) {
	cm.mcache[requesterNfType].setTtlMonitor(targetNfType, ttlMonitor)
	log.Infof("Set home profile ttlMonitor for nfType[%s] success", requesterNfType)
}
*/
/*
func (cm *CacheManager) setRoamingTtlMonitor(requesterNfType string, targetNfType string, ttlMonitor *ttlMonitor) {
	//adpter should offer the api and lock
	cm.roamingCache[requesterNfType].setTtlMonitor(targetNfType, ttlMonitor)
	log.Infof("Set roaming profile ttlMonitor for nfType[%s] success", requesterNfType)
}
*/
/*
func (cm *CacheManager) buildSubscriptionCache() *subscriptionCache {
	subscriptionCache := new(subscriptionCache)
	subscriptionCache.init()
	log.Info("Build subscription cache success")

	return subscriptionCache
}
*/
/*
func (cm *CacheManager) setSubscriptionCache(subscriptionCache *subscriptionCache) {
	cm.subscriptionCache = subscriptionCache
	log.Info("Set subscription cache success")
}
*/
func (cm *CacheManager) loadConfig() bool {
	content, err := ioutil.ReadFile(cacheConfig)
	if err != nil {
		log.Errorf("CacheManager read config file failure")
		return false
	}
	err = json.Unmarshal(content, &cm.indexMapper)
	if err != nil {
		log.Errorf("Unmarshal cache index config failure, please check, Error: %s", err.Error())
		return false
	}

	log.Debugf("CacheManager SearchDataIndexMapper : %v\n", cm.indexMapper)

	/*
		_, err = jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			index := make([]string, 0)
			jsonparser.ArrayEach(value, func(indexItem []byte, dataType jsonparser.ValueType, offset int, err error) {
				log.Debugf("index item : %s", string(indexItem))
				index = append(index, string(indexItem))
			})
			cm.indexGroup = append(cm.indexGroup, index)
		}, "indexItem")

		if err != nil {
			log.Errorf("Json parser cache config file error")
			return false
		}
	*/

	cm.injectionCacheIndexGroup()
	log.Infof("CacheManager indexgroup : %v", cm.indexGroup)

	return true
}

func (cm *CacheManager) injectionCacheIndexGroup() {
	var index string
	if cm.indexMapper.ServiceName != "" {
		index = cm.indexMapper.ServiceName
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.TargetNfType != "" {
		index = cm.indexMapper.TargetNfType
		cm.indexGroup = append(cm.indexGroup, index)
	}
	//	plmnId := PlmnId{}
	//	if cm.indexMapper.TargetPlmn != plmnId {
	//		index := make([]string, 0)
	//		index = append(index, cm.indexMapper.TargetPlmn.Mcc)
	//		index = append(index, cm.indexMapper.TargetPlmn.Mnc)
	//		cm.indexGroup = append(cm.indexGroup, index)
	//	}
	if len(cm.indexMapper.TargetPlmnList) != 0 {
		index = "mcc:mnc"
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.Dnn != "" {
		index = cm.indexMapper.Dnn
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.SmfServingArea != "" {
		index = cm.indexMapper.SmfServingArea
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.RoutingIndicator != "" {
		index = cm.indexMapper.RoutingIndicator
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.NsiList != "" {
		index = cm.indexMapper.NsiList
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.GroupIDList != "" {
		index = cm.indexMapper.GroupIDList
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.IPDomain != "" {
		index = cm.indexMapper.IPDomain
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.DnaiList != "" {
		index = cm.indexMapper.DnaiList
		cm.indexGroup = append(cm.indexGroup, index)
	}
	if cm.indexMapper.UpfIwkEpsInd != "" {
		index = cm.indexMapper.UpfIwkEpsInd
		cm.indexGroup = append(cm.indexGroup, index)
	}
	//log.Infof("injectionCacheIndexGroup cm.indexGroup:%+v", cm.indexGroup)
}

func (cm *CacheManager) assembleResponseContents(requesterNfType string, targetNfType string, ids []string, searchParameter *SearchParameter, isKeepCacheMode bool, isRoaming bool) []byte {
	if len(ids) == 0 {
		return nil
	}

	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return nil
	}

	var searchResult structs.SearchResult
	for _, id := range ids {
		log.Infof("Apply NfServices filter in profile:%s", id)
		content := adapter.fetchProfileByID(targetNfType, id)
		var nfProfile structs.SearchResultNFProfile
		err := json.Unmarshal(content, &nfProfile)
		if err != nil {
			log.Infof("Unmarshal searchResultNFProfile failure, err:%s", err.Error())
			continue
		}

		if !NfServiceFilter(&nfProfile, searchParameter) ||
			!snssaisFilter(&nfProfile, searchParameter) {
			continue
		}

		searchResult.NfInstances = append(searchResult.NfInstances, nfProfile)
	}

	if len(searchResult.NfInstances) != 0 {
		if isKeepCacheMode {
			searchResult.ValidityPeriod = 86400
		} else {
			searchResult.ValidityPeriod = int32(cm.calculateTTL(requesterNfType, targetNfType, ids, isRoaming))
		}
		contents, err := json.Marshal(searchResult)
		if err != nil {
			log.Errorf("Marshal searchResult failure, err:%s", err.Error())
			return nil
		}
		return contents
	} else {
		return nil
	}
}

func (cm *CacheManager) calculateTTL(requesterNfType string, targetNfType string, profileIds []string, isRoaming bool) uint {
	if len(profileIds) <= 0 {
		return 0
	}

	adapter := cm.getCacheAdapter(requesterNfType, isRoaming)
	if adapter == nil {
		return 0
	}

	ttl, _ := adapter.leftLive(targetNfType, profileIds[0])
	for _, id := range profileIds[1:] {
		tmp, _ := adapter.leftLive(targetNfType, id)
		if ttl > tmp {
			ttl = tmp
		}
	}

	return ttl
}

func (ca *CacheManager) getRequesterNfTypes() []string {
	nfTypes := make([]string, 0)
	for nfType, _ := range ca.mcache {
		nfTypes = append(nfTypes, nfType)
	}

	return nfTypes
}

func (ca *CacheManager) getRequesterRoamNfTypes() []string {
	nfTypes := make([]string, 0)
	for nfType, _ := range ca.roamingCache {
		nfTypes = append(nfTypes, nfType)
	}

	return nfTypes
}

//getAllSubscriptionInfo is to get all subscriptionInfo
func (cm *CacheManager) getAllSubscriptionInfo() map[string]structs.SubscriptionInfo {
	subscriptionInfo := make(map[string]structs.SubscriptionInfo, 0)
	for _, pCacheAdapter := range cm.mcache {
		if pCacheAdapter == nil {
			continue
		}
		for subId, subInfo := range pCacheAdapter.getAllSubscriptionInfo() {
			subscriptionInfo[subId] = subInfo
		}
	}

	return subscriptionInfo
}

func (cm *CacheManager) cacheProfileDiffHandler(requesterNfType string, targetNfType string, cache *cache) bool {
	if cache == nil {
		log.Errorf("cache ptr is nil")
		return false
	}

	adapter := cm.getCacheAdapter(requesterNfType, false)
	if adapter == nil {
		return false
	}

	oldIDs := adapter.fetchProfileIDs(targetNfType)
	log.Debugf("nftype[%s,%s], oldIDs: %v", requesterNfType, targetNfType, oldIDs)
	for _, id := range oldIDs {
		if !cache.probe(id) {
			if election.IsActiveLeader("3201", consts.DiscoveryAgentReadinessProbe) {
				util.PushMessageToMSB(requesterNfType, targetNfType, id, consts.NFDeRegister, nil)
			}
		}
	}

	return true
}

//////////////bellow code should extract to util///////////////////

//SetCacheConfig set cache config
func SetCacheConfig(cacheCfg string) {
	cacheConfig = cacheCfg
}

//NfServiceFilter is filter for serviceNames and supportedFeatures
func NfServiceFilter(searchResultNfProfile *structs.SearchResultNFProfile, searchParameter *SearchParameter) bool {
	log.Infof("NfServiceFilter: serviceNames[%+v], supportedFeatures[%s]",
		searchParameter.serviceNames, searchParameter.supportedFeatures)
	for i := 0; i < len(searchResultNfProfile.NfSrvList); {
		if !serviceNamesFilter(searchParameter.serviceNames, searchResultNfProfile.NfSrvList[i].SrvName) ||
			!supportedFeaturesFilter(searchParameter.supportedFeatures, searchResultNfProfile.NfSrvList[i].SupportedFeatures) {
			searchResultNfProfile.NfSrvList = append(searchResultNfProfile.NfSrvList[:i],
				searchResultNfProfile.NfSrvList[i+1:]...)
		} else {
			i++
		}
	}

	if (len(searchParameter.serviceNames) != 0 ||
		len(searchParameter.supportedFeatures) != 0) &&
		len(searchResultNfProfile.NfSrvList) == 0 {
		log.Errorf("NfServiceFilter: requested NF services not found in NF profile")
		return false
	}
	log.Debugf("NfServiceFilter: searchResultNfProfile.NfSrvList: %+v", searchResultNfProfile.NfSrvList)
	return true
}

func serviceNamesFilter(spServiceNames []string, serviceName string) bool {
	if len(spServiceNames) == 0 {
		return true
	}

	matched := false
	for _, svcName := range spServiceNames {
		if serviceName == svcName {
			matched = true
			break
		}
	}
	log.Debugf("serviceNamesFilter: %s is (%v) in %v", serviceName, matched, spServiceNames)
	return matched
}

func supportedFeaturesFilter(spSupportedFeatures, supportedFeatures string) bool {
	if len(spSupportedFeatures) == 0 {
		return true
	}

	if len(spSupportedFeatures)%2 != 0 {
		spSupportedFeatures = "0" + spSupportedFeatures
	}
	decode1, e1 := hex.DecodeString(spSupportedFeatures)
	if e1 != nil {
		return false
	}
	l1 := len(decode1)
	if len(supportedFeatures)%2 != 0 {
		supportedFeatures = "0" + supportedFeatures
	}
	decode2, e2 := hex.DecodeString(supportedFeatures)
	if e2 != nil {
		return false
	}
	l2 := len(decode2)

	matched := false
	for l1 > 0 && l2 > 0 {
		if (decode1[l1-1] & decode2[l2-1]) != 0 {
			matched = true
			break
		}
		l1--
		l2--
	}
	return matched
}

func snssaisFilter(nfProfile *structs.SearchResultNFProfile, searchParameter *SearchParameter) bool {
	log.Infof("snssaisFilter: nfProfile.SNSSAI %+v, searchParameter.snssai %+v",
		nfProfile.SNSSAI, searchParameter.snssai)

	if !searchParameter.searchSnssai() {
		return true
	}

	result := false
	if len(nfProfile.SNSSAI) == 0 {
		return false
	} else {
		for _, nfProfileSnssai := range nfProfile.SNSSAI {
			for _, sNssai := range searchParameter.snssai {
				if nfProfileSnssai.SST == sNssai.Sst &&
					nfProfileSnssai.SD == sNssai.Sd {
					result = true
					break
				}
			}
			if result == true {
				break
			}
		}
	}
	log.Debugf("snssaisFilter result %+v", result)
	return result
}
