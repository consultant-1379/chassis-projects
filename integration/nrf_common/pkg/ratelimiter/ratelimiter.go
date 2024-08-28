package ratelimiter

import (
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"golang.org/x/time/rate"
)

type LimiterWrapper struct {
	rateLimiter *rate.Limiter
	lastUpdate  time.Time
}

var rateLimiterMap = make(map[string]*LimiterWrapper)
var mtx sync.RWMutex

var rateLimit float64
var bucketSize int
var rateLimitMtx sync.Mutex

func Init(limit float64, size int) {
	rateLimit = limit
	bucketSize = size

	go startMonitorExpiredLimiter(3600, 3600)
}

// SetRateLimit reset rateLimit and bucketSize after init
func SetRateLimit(limit float64, size int) {
	rateLimitMtx.Lock()

	rateLimit = limit
	bucketSize = size

	mtx.Lock()
	for instance := range rateLimiterMap {
		limiter := rate.NewLimiter(rate.Limit(rateLimit), bucketSize)
		rateLimiterMap[instance].rateLimiter = limiter
	}
	mtx.Unlock()

	rateLimitMtx.Unlock()
}

// GetRateLimiter return rateLimit and bucketSize
func GetRateLimiter() (limit float64, size int) {
	rateLimitMtx.Lock()
	defer rateLimitMtx.Unlock()
	return rateLimit, bucketSize
}

// Allow returns whether the traffic is allowd or not
func Allow(source string) bool {
	return getRateLimiter(source).Allow()
}

func getRateLimiter(source string) *rate.Limiter {
	mtx.RLock()
	limiterWrapper, exist := rateLimiterMap[source]

	if !exist {
		mtx.RUnlock()
		return createRateLimiter(source)
	}

	defer mtx.RUnlock()
	limiterWrapper.lastUpdate = time.Now()
	return limiterWrapper.rateLimiter
}

func createRateLimiter(source string) *rate.Limiter {
	rateLimitMtx.Lock()
	limiter := rate.NewLimiter(rate.Limit(rateLimit), bucketSize)

	mtx.Lock()
	rateLimiterMap[source] = &LimiterWrapper{limiter, time.Now()}
	mtx.Unlock()

	rateLimitMtx.Unlock()
	return limiter
}

func startMonitorExpiredLimiter(monitorInterval int, expiredTime int) {
	if monitorInterval <= 0 {
		monitorInterval = 3600
	}
	for {
		time.Sleep(time.Duration(monitorInterval) * time.Second)
		mtx.Lock()
		for instance, limiterWrapperPointer := range rateLimiterMap {
			if t := time.Since(limiterWrapperPointer.lastUpdate); t > time.Duration(expiredTime)*time.Second {
				log.Debugf("Delete the rateLimiter for instance %v", instance)
				delete(rateLimiterMap, instance)
			}
		}
		mtx.Unlock()
	}

}
