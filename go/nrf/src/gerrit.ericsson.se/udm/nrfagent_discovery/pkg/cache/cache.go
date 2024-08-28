package cache

import (
	"sync"

	"encoding/json"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache/provider"
	"github.com/buger/jsonparser"
	"github.com/deckarep/golang-set"
)

type cache struct {
	profileMutex     sync.Mutex
	cacheIndexMutex  sync.Mutex
	cacheStatusMutex sync.Mutex
	cacheStatus      bool
	profiles         map[string][]byte                //key: nfInstanceID
	cacheIndex       map[string]map[string]mapset.Set //key: indexCategory, indexValue
	indexGroup       []string
	ttlMonitor       *ttlMonitor
	subscription     *subscriptionCache
	etags            map[string]string //key: nfInstanceID
	mcacheType       cacheType
}

func (c *cache) init(requesterNfType string, targetNfType string, indexGroup []string, subscription *subscriptionCache, mcacheType cacheType) {
	c.profiles = make(map[string][]byte)
	c.cacheIndex = make(map[string]map[string]mapset.Set)
	c.etags = make(map[string]string)

	c.mcacheType = mcacheType

	c.indexGroup = indexGroup
	for _, index := range indexGroup {
		categoryMap := make(map[string]mapset.Set)
		c.cacheIndex[index] = categoryMap
	}
	c.ttlMonitor = new(ttlMonitor)
	c.ttlMonitor.init(requesterNfType, targetNfType, mcacheType)
	c.ttlMonitor.startMonitorWorker()

	if subscription == nil {
		c.subscription = new(subscriptionCache)
		c.subscription.init(requesterNfType, targetNfType, mcacheType)
	} else {
		c.subscription = subscription
	}
}

func (c *cache) cached(data []byte) {
	nfInstanceID, err := jsonparser.GetString(data, provider.NfInstanceId)
	if err != nil {
		log.Warnf("The nfPrpfile miss of nfInstanceID, can not been cached")
		return
	}

	c.profileMutex.Lock()
	defer c.profileMutex.Unlock()

	c.profiles[nfInstanceID] = make([]byte, 0)
	c.profiles[nfInstanceID] = append(c.profiles[nfInstanceID], data...)

	log.Debugf("Cached nfProfile[%s] success", nfInstanceID)
}

func (c *cache) indexed(content []byte, indexGroup []string) (string, bool) {
	result := false
	var sNfProfile structs.SearchResultNFProfile
	err := json.Unmarshal(content, &sNfProfile)
	if err != nil {
		log.Errorf("Index: Unmarshal SearchResultNFProfile failure. err(%s)", err.Error())
		return "", false
	}

	for _, indexCategory := range indexGroup {

		var indexContent string
		switch indexCategory {
		case "serviceName":
			exist := false
			for _, nfService := range sNfProfile.NfSrvList {
				if nfService.SrvName != "" {
					c.index(indexCategory, nfService.SrvName, sNfProfile.NfInstanceID)
					result = true
					exist = true
				}
			}
			if !exist {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		case "mcc:mnc":
			exist := false
			for _, plmn := range sNfProfile.PLMN {
				mccMnc := plmn.Mcc + ":" + plmn.Mnc
				c.index(indexCategory, mccMnc, sNfProfile.NfInstanceID)
				result = true
				exist = true
			}
			if !exist {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		case "dnn":
			exist := false
			var dnnList = make([]string, 0)
			if sNfProfile.NfType == "PCF" && sNfProfile.PcfInfo != nil {
				dnnList = sNfProfile.PcfInfo.Dnnlist
			} else if sNfProfile.NfType == "BSF" && sNfProfile.BsfInfo != nil {
				dnnList = sNfProfile.BsfInfo.Dnnlist
			} else if sNfProfile.NfType == "UPF" && sNfProfile.UpfInfo != nil {
				for _, sNssaiUpfInfo := range sNfProfile.UpfInfo.SNssaiUpfInfoList {
					for _, dnnUpfInfo := range sNssaiUpfInfo.DNNUpfInfoList {
						if dnnUpfInfo.DNN != "" {
							dnnList = append(dnnList, dnnUpfInfo.DNN)
						}
					}
				}
			}
			for _, dnn := range dnnList {
				c.index(indexCategory, dnn, sNfProfile.NfInstanceID)
				result = true
				exist = true
			}
			if !exist {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		case "routingIndicators":
			exist := false
			var routingIndicators = make([]string, 0)
			if sNfProfile.NfType == "UDM" && sNfProfile.UdmInfo != nil {
				routingIndicators = sNfProfile.UdmInfo.RoutingIndicators
			} else if sNfProfile.NfType == "AUSF" && sNfProfile.AusfInfo != nil {
				routingIndicators = sNfProfile.AusfInfo.RoutingIndicators
			}
			for _, routingIndicator := range routingIndicators {
				c.index(indexCategory, routingIndicator, sNfProfile.NfInstanceID)
				result = true
				exist = true
			}
			if !exist {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		case "smfServingArea":
			if sNfProfile.NfType == "UPF" && sNfProfile.UpfInfo != nil {
				exist := false
				for _, areaValue := range sNfProfile.UpfInfo.SmfServingArea {
					c.index(indexCategory, areaValue, sNfProfile.NfInstanceID)
					result = true
					exist = true
				}
				if !exist {
					log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
				}
			}
		case "nsiList":
			exist := false
			for _, nsiItem := range sNfProfile.NsiList {
				if nsiItem != "" {
					c.index(indexCategory, nsiItem, sNfProfile.NfInstanceID)
					result = true
					exist = true
				}
			}
			if !exist {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		case "groupId":
			var groupID string
			var supiRangeList = make([]structs.SupiRange, 0)
			var gpsiRangeList = make([]structs.IdentityRange, 0)
			if sNfProfile.NfType == "UDR" && sNfProfile.UdrInfo != nil {
				groupID = sNfProfile.UdrInfo.GroupID
				supiRangeList = sNfProfile.UdrInfo.SupiRanges
				gpsiRangeList = sNfProfile.UdrInfo.GpsiRanges
			} else if sNfProfile.NfType == "UDM" && sNfProfile.UdmInfo != nil {
				groupID = sNfProfile.UdmInfo.GroupID
				supiRangeList = sNfProfile.UdmInfo.SupiRanges
				gpsiRangeList = sNfProfile.UdmInfo.GpsiRanges
			} else if sNfProfile.NfType == "AUSF" && sNfProfile.AusfInfo != nil {
				groupID = sNfProfile.AusfInfo.GroupID
				supiRangeList = sNfProfile.AusfInfo.SupiRanges
			} else {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
				break
			}

			if groupID != "" {
				c.index(indexCategory, groupID, sNfProfile.NfInstanceID)
				result = true
			} else {
				//NFInfo not include groupID/supiRange/GpsiRange, it match all groupID
				if len(supiRangeList) == 0 && len(gpsiRangeList) == 0 {
					c.index(indexCategory, provider.MatchAllGroupID, sNfProfile.NfInstanceID)
					result = true
				} else {
					log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
				}
			}
		case "ipDomainList":
			exist := false
			if sNfProfile.BsfInfo != nil {
				for _, domainItem := range sNfProfile.BsfInfo.IPDomainList {
					c.index(indexCategory, domainItem, sNfProfile.NfInstanceID)
					result = true
					exist = true
				}
			}
			if !exist {
				log.Infof("cacheIndex %s:%s does not exist in nfInstanceId[%s]", provider.BsfInfo, indexCategory, sNfProfile.NfInstanceID)
			}
		case "dnaiList":
			exist := false
			if sNfProfile.UpfInfo != nil {
				for _, sNssaiUpfInfo := range sNfProfile.UpfInfo.SNssaiUpfInfoList {
					for _, dnnUpfInfo := range sNssaiUpfInfo.DNNUpfInfoList {
						for _, dnaiItem := range dnnUpfInfo.DnaiList {
							c.index(indexCategory, dnaiItem, sNfProfile.NfInstanceID)
							result = true
							exist = true
						}
					}
				}
			}
			if !exist {
				log.Infof("cacheIndex %s:%s does not exist in nfInstanceId[%s]", provider.UpfInfo, indexCategory, sNfProfile.NfInstanceID)
			}
		case "iwkEpsInd":
			exist := false
			if sNfProfile.UpfInfo != nil {
				if sNfProfile.UpfInfo.IwkEpsInd != nil {
					if *(sNfProfile.UpfInfo.IwkEpsInd) {
						c.index(indexCategory, "true", sNfProfile.NfInstanceID)
					} else {
						c.index(indexCategory, "false", sNfProfile.NfInstanceID)
					}
					result = true
					exist = true
				}
			}
			if !exist {
				log.Infof("cacheIndex %s:%s does not exist in nfInstanceId[%s]", provider.UpfInfo, indexCategory, sNfProfile.NfInstanceID)
			}
		case "nfType":
			if sNfProfile.NfType != "" {
				c.index(indexCategory, sNfProfile.NfType, sNfProfile.NfInstanceID)
				result = true
			} else {
				log.Warningf("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
			}
		default:
			indexContent, err = jsonparser.GetString(content, indexCategory)
			if err != nil {
				log.Infof("cacheIndex %s does not exist in nfInstanceId[%s]", indexCategory, sNfProfile.NfInstanceID)
				break
			}
			c.index(indexCategory, indexContent, sNfProfile.NfInstanceID)
			result = true
		}
	}

	return sNfProfile.NfInstanceID, result
}

func (c *cache) deIndex(id string) {
	c.cacheIndexMutex.Lock()
	defer c.cacheIndexMutex.Unlock()

	for categoryKey := range c.cacheIndex {
		for valueKey := range c.cacheIndex[categoryKey] {
			if c.cacheIndex[categoryKey][valueKey] == nil {
				continue
			}
			// Remove immediately?
			// if !c.cacheIndex[categoryKey][valueKey].Contains(id) {
			// 	continue
			// }
			c.cacheIndex[categoryKey][valueKey].Remove(id)
		}
	}
}

func (c *cache) deIndexAll() {
	c.cacheIndexMutex.Lock()
	defer c.cacheIndexMutex.Unlock()

	for categoryKey := range c.cacheIndex {
		for valueKey := range c.cacheIndex[categoryKey] {
			if c.cacheIndex[categoryKey][valueKey] != nil {
				c.cacheIndex[categoryKey][valueKey].Clear()
			}
			delete(c.cacheIndex[categoryKey], valueKey)
		}
	}
	log.Infof("DeIndexAll: decache all the cacheIndex")
}

func (c *cache) deCached(id string) {
	c.profileMutex.Lock()
	defer c.profileMutex.Unlock()

	_, ok := c.profiles[id]
	if ok == false {
		log.Infof("Cache no such profile for nfInstanceID[%s]\n", id)
	} else {
		delete(c.profiles, id)
	}
}

func (c *cache) deCachedAll() {
	c.profileMutex.Lock()
	defer c.profileMutex.Unlock()

	for k, _ := range c.profiles {
		delete(c.profiles, k)
	}
	log.Infof("DeCachedAll: delete all cache content")
}

func (c *cache) probe(id string) bool {
	c.profileMutex.Lock()
	defer c.profileMutex.Unlock()

	_, ok := c.profiles[id]

	return ok
}

func (c *cache) fetchProfileByID(id string) []byte {
	c.profileMutex.Lock()
	defer c.profileMutex.Unlock()

	content, ok := c.profiles[id]
	if !ok {
		log.Infof("There is no such nfProfile for NfInstanceId[%s]", id)
		return nil
	}

	return content
}

func (c *cache) showIndexContent() {
	for k, v := range c.cacheIndex {
		log.Debugf("key:%s value:%+v", k, v)
	}
}

func (c *cache) flush() {
	c.deIndexAll()
	c.deCachedAll()
	c.etags = make(map[string]string)
}

func (c *cache) sync(cacheInfo *structs.CacheSyncInfo) {
	if cacheInfo == nil {
		log.Errorf("sync: cacheInfo PTR is nil")
		return
	}
	for instanceID, profile := range c.profiles {
		cacheInfo.NfProfiles = append(cacheInfo.NfProfiles, profile)
		if c.ttlMonitor != nil {
			timePoint, err := c.ttlMonitor.getTimePoint(instanceID)
			if err == nil {
				ttlInfo := structs.TtlInfo{
					NfInstanceID: instanceID,
					ValidityTime: timePoint,
				}
				cacheInfo.TtlInfos = append(cacheInfo.TtlInfos, ttlInfo)
			}
		}
	}
	if c.subscription != nil {
		cacheInfo.SubscriptionInfos = c.subscription.getAllSubscriptionInfo()
	}
	for instanceID, etag := range c.etags {
		etagInfo := structs.EtagInfo{
			NfInstanceID: instanceID,
			FingerPrint:  etag,
		}
		cacheInfo.EtagInfos = append(cacheInfo.EtagInfos, etagInfo)
	}
}

func (c *cache) dump(cacheInfo *structs.CacheDumpInfo) {
	if cacheInfo == nil {
		log.Errorf("Dump: cacheInfo PTR is nil")
		return
	}

	for instanceID, profile := range c.profiles {
		cacheInfo.NfProfiles = append(cacheInfo.NfProfiles, string(profile))

		if c.ttlMonitor != nil {
			timePoint, err := c.ttlMonitor.getTimePoint(instanceID)
			if err == nil {
				ttlInfo := structs.TtlInfo{
					NfInstanceID: instanceID,
					ValidityTime: timePoint,
				}
				cacheInfo.TtlInfos = append(cacheInfo.TtlInfos, ttlInfo)
			}
		}
	}
	if c.subscription != nil {
		cacheInfo.SubscriptionInfos = c.subscription.getAllSubscriptionInfo()
	}
	for instanceID, etag := range c.etags {
		etagInfo := structs.EtagInfo{
			NfInstanceID: instanceID,
			FingerPrint:  etag,
		}
		cacheInfo.EtagInfos = append(cacheInfo.EtagInfos, etagInfo)
	}
}

func (c *cache) fetchIDs() []string {
	ids := make([]string, 0)
	for k, _ := range c.profiles {
		ids = append(ids, k)
	}

	return ids
}

func (c *cache) dumpByID(id string) []byte {
	profile, ok := c.profiles[id]
	if !ok {
		log.Errorf("There is no such cached profile for nfinstanceid[%s]", id)
	}

	return profile
}

func (c *cache) getCacheStatus() bool {
	c.cacheStatusMutex.Lock()
	defer c.cacheStatusMutex.Unlock()

	return c.cacheStatus
}

func (c *cache) setCacheStatus(status bool) {
	c.cacheStatusMutex.Lock()
	defer c.cacheStatusMutex.Unlock()

	c.cacheStatus = status
}

func (c *cache) index(categoryKey string, valueKey string, nfInstanceId string) {
	c.cacheIndexMutex.Lock()
	defer c.cacheIndexMutex.Unlock()

	if _, ok := c.cacheIndex[categoryKey]; !ok {
		c.cacheIndex[categoryKey] = make(map[string]mapset.Set)
	}
	if _, ok := c.cacheIndex[categoryKey][valueKey]; !ok {
		c.cacheIndex[categoryKey][valueKey] = mapset.NewSet()
	}

	if c.cacheIndex[categoryKey][valueKey].Contains(nfInstanceId) {
		log.Debugf("cacheIndex %s[%s] has already created for nfInstanceId[%s]", categoryKey, valueKey, nfInstanceId)
		return
	}
	log.Infof("Create cacheIndex %s[%s] for nfInstanceId[%s]", categoryKey, valueKey, nfInstanceId)
	c.cacheIndex[categoryKey][valueKey].Add(nfInstanceId)
}

func (c *cache) haveEtag(nfInstanceID string) bool {
	_, ok := c.etags[nfInstanceID]
	if !ok {
		return false
	}

	return true
}

func (c *cache) fetchEtag(nfInstanceID string) string {
	return c.etags[nfInstanceID]
}

func (c *cache) saveEtag(nfInstanceID string, value string) bool {
	if c.etags == nil {
		return false
	}
	c.etags[nfInstanceID] = value
	return true
}

func (c *cache) deleteEtag(nfInstanceID string) bool {
	if c.etags == nil {
		return false
	}

	if c.haveEtag(nfInstanceID) {
		delete(c.etags, nfInstanceID)
	}

	return true
}

func (c *cache) deleteAllEtag() {
	c.etags = make(map[string]string)
}

////////////////////ttlMonitor related function//////////////////////////
func (c *cache) getTtlMonitor() *ttlMonitor {
	return c.ttlMonitor
}

func (c *cache) setTtlMonitor(ttlMonitor *ttlMonitor) {
	c.ttlMonitor = ttlMonitor
}

func (c *cache) getSubscriptionCache() *subscriptionCache {
	return c.subscription
}

/*
func (c *cache) checkCacheMapExist(key string) bool {
	for _, cacheItem := range c.cacheContent {
		_, ok := cacheItem[key]
		if ok {
			return true
		}
	}

	return false
}*/

/*
func (c *cache) fetchCacheMap(key string) map[string][]string {
	ok := c.checkCacheMapExist(key)
	if !ok {
		cacheMap := make(map[string][]string)
		cacheMap[key] = make([]string, 0)
		c.cacheContent = append(c.cacheContent, cacheMap)
	}

	for _, cacheMap := range c.cacheContent {
		_, ok := cacheMap[key]
		if ok {
			return cacheMap
		}
	}

	return nil
}*/
/*
func (c *cache) index(key string, nfInstanceId string) {
	cacheMap := c.fetchCacheMap(key)
	if cacheMap == nil {
		return
	}

	cacheMap[key] = append(cacheMap[key], nfInstanceId)
}*/

// func (c *cache) getIndex(array []string, element string) int {
// 	for i, e := range array {
// 		if e == element {
// 			return i
// 		}
// 	}

// 	return -1
// }
