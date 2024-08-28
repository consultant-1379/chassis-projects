package internalconf

import (
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

const (
	defaultNrfInfoCheckInterval = 60
)

var (
	// NrfInfoCheckInterval is the interval to check whether the nrfInfo is changed or not
	NrfInfoCheckInterval = defaultNrfInfoCheckInterval

	//TrafficRateLimitPerNfInstance is the number of requests per second allowed for each NF instance
	TrafficRateLimitPerNfInstance = 10

	// RetryAfterRangeStart is min of RetryAfter Range
	// In case of service overload, the HTTP header field "Retry-After" will be added in the response
	// to indicate how many seconds the NF instance has to wait before making a new request.
	RetryAfterRangeStart = 3

	// RetryAfterRangeStart is max of RetryAfter Range
	RetryAfterRangeEnd = 10

	// HeartBeatTimerMin is the min allowed heartbeat timer
	HeartBeatTimerMin = 5

	// HeartBeatTimerOffset is the offset against the heartbeat default value
	HeartBeatTimerOffset = 10
)

// OverloadLogChangeHandler is function type to handle the case that log changes in configmap
type OverloadLogChangeHandler func(bool)

var overloadLogChangeHandler OverloadLogChangeHandler

// InternalMgmtConf is format of internalConf.json in NRF Management
type InternalMgmtConf struct {
	HomeNrf            HomeNrfConf            `json:"homeNrf"`
	HttpServer         HttpServerConf         `json:"httpServer"`
	Dbproxy            DbproxyConf            `json:"dbProxy"`
	RegionNrf          MgmtRegionNrfConf      `json:"regionNrf, omitempty"`
	NFStatusNotify     NFStatusNotifyConf     `json:"NFStatusNotify"`
	OverloadProtection OverloadProtectionConf `json:"OverloadProtection"`
	OverloadControl    *OverloadControlConf   `json:"overloadControl"`
	PriorityPolicy     *PriorityPolicyConf    `json:"PriorityPolicy"`
	Heartbeat          *HeartbeatConf         `json:"heartbeat"`
}

// MgmtRegionNrfConf is format of obejct "regionNrf" in internalConf.json
type MgmtRegionNrfConf struct {
	NrfInfoCheckInterval int `json:"nrfInfoCheckInterval, omitempty"`
}

// OverloadControlConf is for some configuration of overload control
type OverloadControlConf struct {
	EnableDebugLog                bool        `json:"enableDebugLog"`
	TrafficRateLimitPerNfInstance int         `json:"trafficRateLimitPerNfInstance"`
	RetryAfterRange               NumberRange `json:"retryAfterRange"`
}

// HeartbeatConf is internal conf for heartbeat. min is the min allowed heartbeat timer.
type HeartbeatConf struct {
	Min    int `json:"min"`
	Offset int `json:"offset"`
}

// NumberRange is a range of number including start and end
type NumberRange struct {
	Start int `json:"start`
	End   int `json:"end"`
}

// ParseConf sets the value of global variable from struct
func (conf *InternalMgmtConf) ParseConf() {

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

	NrfInfoCheckInterval = conf.RegionNrf.NrfInfoCheckInterval
	if NrfInfoCheckInterval <= 0 {
		NrfInfoCheckInterval = defaultNrfInfoCheckInterval
	}

	EnableNotification = conf.NFStatusNotify.Enable

	OverloadProtection = conf.OverloadProtection

	if conf.OverloadControl != nil {
		if overloadLogChangeHandler != nil {
			overloadLogChangeHandler(conf.OverloadControl.EnableDebugLog)
		}
		TrafficRateLimitPerNfInstance = conf.OverloadControl.TrafficRateLimitPerNfInstance
		RetryAfterRangeStart = conf.OverloadControl.RetryAfterRange.Start
		RetryAfterRangeEnd = conf.OverloadControl.RetryAfterRange.End
	}

	if conf.Heartbeat != nil {
		HeartBeatTimerMin = conf.Heartbeat.Min
		HeartBeatTimerOffset = conf.Heartbeat.Offset
	}
}

// ShowInternalConf shows the value of global variable
func (conf *InternalMgmtConf) ShowInternalConf() {
	fmt.Printf("home NRF forward retry time : %d\n", HomeNrfForwardRetryTime)
	fmt.Printf("home NRF forward retry wait time: %d\n", HomeNrfForwardRetryWaitTime)
	fmt.Printf("http server idleTimeout : %v\n", httpserver.GetIdleTimeout())
	fmt.Printf("http server activeTimeout : %v\n", httpserver.GetActiveTimeout())
	fmt.Printf("http header with x-version : %v\n", HTTPWithXVersion)
	fmt.Printf("DbproxyConnectionNum  : %d\n", DbproxyConnectionNum)
	fmt.Printf("DbproxyGrpcCtxTimeout : %d\n", DbproxyGrpcCtxTimeout)
	fmt.Printf("NrfInfoCheckInterval : %d\n", NrfInfoCheckInterval)
	fmt.Printf("EnableNotification : %v\n", EnableNotification)
	fmt.Printf("OverloadControlEnabled = %v\n", OverloadProtection.Enabled)
	fmt.Printf("OverloadControlLevel = %v, OverloadTriggerLatencyThreshold = %v ms, OverloadControlLatencyThreshold = %v ms, OverloadTriggerSampleWindow = %v, OverloadControlSampleWindow = %v, OverloadTriggerTimeSampleWindow = %v, OverloadControlTimeSampleWindow = %v, IdleInterval = %v ms, IdleRecoverRatio = %v, CounterReportInterval = %v ms, OverloadAlarmClearWindow = %v\n", OverloadProtection.OverloadControlLevel, OverloadProtection.OverloadTriggerLatencyThreshold, OverloadProtection.OverloadControlLatencyThreshold, OverloadProtection.OverloadTriggerSampleWindow, OverloadProtection.OverloadControlSampleWindow, OverloadProtection.OverloadTriggerTimeSampleWindow, OverloadProtection.OverloadControlTimeSampleWindow, OverloadProtection.IdleInterval, OverloadProtection.IdleRecoverRatio, OverloadProtection.CounterReportInterval, OverloadProtection.OverloadAlarmClearWindow)
	for _, engine := range OverloadProtection.WorkEngines {
		fmt.Printf("GroupPriority = %v, QueueCapacity = %v, WorkerNumber = %v\n", engine.GroupPriority, engine.QueueCapacity, engine.WorkerNumber)
	}

	fmt.Printf("The priority Policy is %v", conf.PriorityPolicy)
	fmt.Printf("TrafficRateLimitPerNfInstance : %v\n", TrafficRateLimitPerNfInstance)
	fmt.Printf("RetryAfterRange Start : %v\n", RetryAfterRangeStart)
	fmt.Printf("RetryAfterRange End : %v\n", RetryAfterRangeEnd)

	fmt.Printf("The HeartBeatTimerMin is %v", HeartBeatTimerMin)
	fmt.Printf("The HeartBeatTimerOffset is %v", HeartBeatTimerOffset)
}

// RegistOverloadLevelChangeHandler
func RegistOverloadLevelChangeHandler(f OverloadLogChangeHandler) {
	overloadLogChangeHandler = f
}
