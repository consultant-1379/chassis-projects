package nfdisccache

import (
	"sync"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
)

//CacheItem is item stored in cache
type CacheItem struct {
	Key               string
	BodyCommon        string
	NfInfo            string
	NfServices        string
	MD5Sum            string
	Value             interface{}
	ProfileUpdateTime int
	ExpiredTime       int
}

//Cache is mapping of cache
type Cache struct {
	*cache
}

//cache is struct of cache
type cache struct {
	items          sync.Map
	updateChannel  chan string
	AddDataChannel chan []string
	stopChannel    chan bool
	//mu                sync.RWMutex
}

//newCache create a cache
func newCache() *Cache {
	var itemMap sync.Map
	c := &cache{items:itemMap}
	c.updateChannel = make(chan string, 100)
	c.AddDataChannel = make(chan []string, 500)
	c.stopChannel = make(chan bool, 1)
	C := &Cache{c}
	return C
}

//set the cache item into Cache
func (c *cache) set(item CacheItem) {
	c.items.Store(item.Key, item)
}

//get value from cache by key
func (c *cache) get(key string, profileUpdateTime int) (CacheItem, bool) {
	value, found := c.items.Load(key)
	if !found {
		return CacheItem{}, false
	}
	if value.(CacheItem).ProfileUpdateTime == profileUpdateTime {
		return value.(CacheItem), true
	} else if profileUpdateTime - value.(CacheItem).ProfileUpdateTime <= internalconf.DiscCacheTimeThreshold {
		//data is not latest, need to fetch new data
		c.updateChannel <- key
		return value.(CacheItem), true
	}
	return CacheItem{}, false
}

//startCache start a goroutine for cache function
func (c *Cache) startCache() {
	go c.start()
}

//start process chan message
func (c *Cache) start() {
	for {
		select {
		case key := <-c.updateChannel:
			c.updateCacheItem(key)
		case profileList := <-c.AddDataChannel:
			c.addCacheItems(profileList)
		case <-c.stopChannel:
			return
		}
	}
}

//stop add stop channel
func (c *Cache) stop() {
	c.stopChannel <- true
}

//updateCacheItem update cache item from db
func (c *cache) updateCacheItem(key string) {
	newProfile := getNFProfileByInstID(key)
	if newProfile != "" {
		_, found := c.items.Load(key)
		if found {
			log.Debugf("%v update from db, will update the cache data", key)
			c.addCacheItems([]string{newProfile})
		}
	}
}

func (c *cache) addCacheItems(profileList []string) {
	itemList := SplitNFProfileList(profileList)
	for _, item := range itemList {
		c.set(item)
	}
}