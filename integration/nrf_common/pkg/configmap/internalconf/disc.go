package internalconf

import (
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

const (
	defaultRegionNrfForwardRetryTime     = 3
	defaultRegionNrfForwardRetryWaitTime = 1

	defaultRegionNrfRedirectRetryTime     = 3
	defaultRegionNrfRedirectRetryWaitTime = 1

	defaultPlmnNrfForwardRetryTime     = 3
	defaultPlmnNrfForwardRetryWaitTime = 1

	defaultDiscCacheEnable        = false
	defaultDiscCacheTimeThreshold = 0

	defaultEnableTimeStatistics = false
	defaultStatisticsNum        = 100
)

var (
	//identify how many times  to retry after forwarding to PLMN NRF failed
	RegionNrfForwardRetryTime = defaultRegionNrfForwardRetryTime

	//identify how long to wait to retry
	RegionNrfForwardRetryWaitTime = defaultRegionNrfForwardRetryWaitTime

	//identify how many times  to retry after redirect to Region NRF failed
	RegionNrfRedirectRetryTime = defaultRegionNrfRedirectRetryTime

	//identify how long to wait to retry
	RegionNrfRedirectRetryWaitTime = defaultRegionNrfRedirectRetryWaitTime

	//identify how many times to retry after forwarding to target region nrf failed
	PlmnNrfForwardRetryTime = defaultPlmnNrfForwardRetryTime

	//identify  how long to wait to retry
	PlmnNrfForwardRetryWaitTime = defaultPlmnNrfForwardRetryWaitTime

	//DiscCacheTimeThreshold is the threshold cache item lastUpdateTime compare to lastUpdateTime in db
	DiscCacheTimeThreshold = defaultDiscCacheTimeThreshold

	//EnableTimeStatistics is a switch to decide whether enable time statistics in discovery
	EnableTimeStatistics = defaultEnableTimeStatistics
	//StatisticsNum is a request num limit to do avg time statistics
	StatisticsNum = defaultStatisticsNum
)

// InternalDiscConf is format of internalConf.json in NRF Discovery
type InternalDiscConf struct {
	HomeNrf            HomeNrfConf            `json:"homeNrf"`
	HttpServer         HttpServerConf         `json:"httpServer"`
	Dbproxy            DbproxyConf            `json:"dbProxy"`
	RegionNrf          DiscRegionNrfConf      `json:"regionNrf, omitempty"`
	PlmnNrf            PlmnNrfConf            `json:"plmnNrf, omitempty"`
	OverloadProtection OverloadProtectionConf `json:"OverloadProtection"`
	DiscCache          DiscCacheConf          `json:"discCache"`
	TimeStatistics     TimeStatisticsConf     `json:"timestatistics"`
}

// DiscRegionNrfConf is format of obejct "regionNrf" in internalConf.json
type DiscRegionNrfConf struct {
	Forward  ForwardConf `json:"forward, omitempty"`
	Redirect ForwardConf `json:"redirect, omitempty"`
}

// PlmnNrfConf is format of object "plmnNrf" in internalConf.json
type PlmnNrfConf struct {
	Forward *ForwardConf `json:"forward"`
}

//DiscCacheConf is config about the discovery cache in memory
type DiscCacheConf struct {
	DiscCacheTimeThreshold int `json:"discCacheTimeThreshold"`
}

//TimeStatisticsConf is config about the discovery do time statistics
type TimeStatisticsConf struct {
	EnableTimestatistics bool `json:"enableTimeStatistics"`
	StatisticsNum        int  `json:"statisticsNum"`
}

// ParseConf sets the value of global variable from struct
func (conf *InternalDiscConf) ParseConf() {

	HomeNrfForwardRetryTime = conf.HomeNrf.Forward.RetryTime
	if HomeNrfForwardRetryTime <= 0 {
		HomeNrfForwardRetryTime = defaultHomeNrfForwardRetryTime
	}

	HomeNrfForwardRetryWaitTime = conf.HomeNrf.Forward.RetryWaitTime
	if HomeNrfForwardRetryWaitTime <= 0 {
		HomeNrfForwardRetryWaitTime = defaultHomeNrfForwardRetryWaitTime
	}

	httpServerIdleTimeout := conf.HttpServer.IdleTimeout
	if httpServerIdleTimeout <= 0 {
		httpServerIdleTimeout = defaultHttpServerIdleTimeout
	}
	httpserver.SetIdleTimeout(httpServerIdleTimeout)

	httpServerActiveTimeout := conf.HttpServer.ActiveTimeout
	if httpServerActiveTimeout <= 0 {
		httpServerActiveTimeout = defaultHttpServerActiveTimeout
	}
	httpserver.SetActiveTimeout(httpServerActiveTimeout)

	HTTPWithXVersion = conf.HttpServer.HTTPWithXVersion

	DbproxyConnectionNum = conf.Dbproxy.ConnectionNum
	if DbproxyConnectionNum <= 0 {
		DbproxyConnectionNum = defaultDbproxyConnectionNum
	}

	DbproxyGrpcCtxTimeout = conf.Dbproxy.GrpcCtxTimeout
	if DbproxyGrpcCtxTimeout <= 0 {
		DbproxyGrpcCtxTimeout = defaultDbproxyGrpcCtxTimeout
	}

	RegionNrfForwardRetryTime = conf.RegionNrf.Forward.RetryTime
	if RegionNrfForwardRetryTime <= 0 {
		RegionNrfForwardRetryTime = defaultRegionNrfForwardRetryTime
	}

	RegionNrfForwardRetryWaitTime = conf.RegionNrf.Forward.RetryWaitTime
	if RegionNrfForwardRetryWaitTime <= 0 {
		RegionNrfForwardRetryWaitTime = defaultRegionNrfForwardRetryWaitTime
	}

	RegionNrfRedirectRetryTime = conf.RegionNrf.Redirect.RetryTime
	if RegionNrfRedirectRetryTime <= 0 {
		RegionNrfRedirectRetryTime = defaultRegionNrfRedirectRetryTime
	}

	RegionNrfRedirectRetryWaitTime = conf.RegionNrf.Redirect.RetryWaitTime
	if RegionNrfRedirectRetryWaitTime <= 0 {
		RegionNrfRedirectRetryWaitTime = defaultRegionNrfRedirectRetryWaitTime
	}

	PlmnNrfForwardRetryTime = conf.PlmnNrf.Forward.RetryTime
	if PlmnNrfForwardRetryTime <= 0 {
		PlmnNrfForwardRetryTime = defaultPlmnNrfForwardRetryTime
	}

	PlmnNrfForwardRetryWaitTime = conf.PlmnNrf.Forward.RetryWaitTime
	if PlmnNrfForwardRetryWaitTime <= 0 {
		PlmnNrfForwardRetryWaitTime = defaultPlmnNrfForwardRetryWaitTime
	}

	DiscCacheTimeThreshold = conf.DiscCache.DiscCacheTimeThreshold
	if DiscCacheTimeThreshold < 0 {
		DiscCacheTimeThreshold = defaultDiscCacheTimeThreshold
	}

	EnableTimeStatistics = conf.TimeStatistics.EnableTimestatistics
	StatisticsNum = conf.TimeStatistics.StatisticsNum
	if StatisticsNum <= 0 {
		StatisticsNum = defaultStatisticsNum
	}

	OverloadProtection = conf.OverloadProtection
}

// ShowInternalConf shows the value of global variable
func (conf *InternalDiscConf) ShowInternalConf() {
	fmt.Printf("home NRF forward retry time : %d\n", HomeNrfForwardRetryTime)
	fmt.Printf("home NRF forward retry wait time: %d\n", HomeNrfForwardRetryWaitTime)
	fmt.Printf("http server idleTimeout : %v\n", httpserver.GetIdleTimeout())
	fmt.Printf("http server activeTimeout : %v\n", httpserver.GetActiveTimeout())
	fmt.Printf("http header with x-version : %v\n", HTTPWithXVersion)
	fmt.Printf("DbproxyConnectionNum  : %d\n", DbproxyConnectionNum)
	fmt.Printf("DbproxyGrpcCtxTimeout : %d\n", DbproxyGrpcCtxTimeout)
	fmt.Printf("region NRF forward retry time : %d\n", RegionNrfForwardRetryTime)
	fmt.Printf("region NRF forward retry wait time : %d\n", RegionNrfForwardRetryWaitTime)
	fmt.Printf("region NRF redirect retry time : %d\n", RegionNrfRedirectRetryTime)
	fmt.Printf("region NRF redirect retry wait time : %d\n", RegionNrfRedirectRetryWaitTime)
	fmt.Printf("plmn NRF forward retry time : %d\n", PlmnNrfForwardRetryTime)
	fmt.Printf("plmn NRF forward retry wait time : %d\n", PlmnNrfForwardRetryWaitTime)
	fmt.Printf("disc cache time threshold : %v\n", DiscCacheTimeThreshold)
	fmt.Printf("disc enable time statistics : %v\n", EnableTimeStatistics)
	fmt.Printf("disc time statistics num : %v\n", StatisticsNum)
	fmt.Printf("OverloadControlEnabled = %v, ", OverloadProtection.Enabled)
	fmt.Printf("OverloadControlLevel = %v, OverloadTriggerLatencyThreshold = %v ms, OverloadControlLatencyThreshold = %v ms, OverloadTriggerSampleWindow = %v, OverloadControlSampleWindow = %v, WorkerNumber = %v, WorkerQueueCapacity = %v, DeniedWorkerNumber = %v, DeniedWorkerQueueCapacity = %v, IdleInterval = %v ms, IdleRecoverRatio = %v, CounterReportInterval = %v ms, OverloadAlarmClearWindow = %v\n", OverloadProtection.OverloadControlLevel, OverloadProtection.OverloadTriggerLatencyThreshold, OverloadProtection.OverloadControlLatencyThreshold, OverloadProtection.OverloadTriggerSampleWindow, OverloadProtection.OverloadControlSampleWindow, OverloadProtection.WorkerNumber, OverloadProtection.WorkerQueueCapacity, OverloadProtection.DeniedWorkerNumber, OverloadProtection.DeniedWorkerQueueCapacity, OverloadProtection.IdleInterval, OverloadProtection.IdleRecoverRatio, OverloadProtection.CounterReportInterval, OverloadProtection.OverloadAlarmClearWindow)
	for _, engine := range OverloadProtection.WorkEngines {
		fmt.Printf("GroupPriority = %v, ", engine.GroupPriority)
	}
}
