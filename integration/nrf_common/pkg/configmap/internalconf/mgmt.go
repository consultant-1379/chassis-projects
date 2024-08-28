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

	// PlmnNrfStatusCheckInterval is the interval to check whether plmn nrf is available, it should be smaller than expireTime of alarm nrfMngtNrfConnectionFailure
	PlmnNrfStatusCheckInterval = 300

	// IsSwitchPlmnForOverloadResp indicates if nrf as consumer try another plmn nrf when it receives a respond with statusCode 503/429
	IsSwitchPlmnForOverloadResp = true

	// ThrottlePercentageForOverloadResp indicates if it will trigger throttling about sending request to plmn nrf for nrfprofile when statusCode(503/429) of response exceed this percent 20%
	ThrottlePercentageForOverloadResp = 20

	// RetryAfterRangeStart is min of RetryAfter Range
	// In case of service overload, the HTTP header field "Retry-After" will be added in the response
	// to indicate how many seconds the NF instance has to wait before making a new request.
	RetryAfterRangeStart = 1

	// RetryAfterRangeEnd is max of RetryAfter Range
	RetryAfterRangeEnd = 120

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
	Heartbeat          *HeartbeatConf         `json:"heartbeat"`
}

// MgmtRegionNrfConf is format of obejct "regionNrf" in internalConf.json
type MgmtRegionNrfConf struct {
	NrfInfoCheckInterval       int `json:"nrfInfoCheckInterval, omitempty"`
	PlmnNrfStatusCheckInterval int `json:"plmnNrfStatusCheckInterval, omitempty"`
}

// OverloadControlConf is for some configuration of overload control
type OverloadControlConf struct {
	EnableDebugLog                    bool        `json:"enableDebugLog"`
	IsSwitchPlmnForOverloadResp       bool        `json:"isSwitchPlmnForOverloadResp"`
	ThrottlePercentageForOverloadResp int         `json:"throttlePercentageForOverloadResp"`
	RetryAfterRange                   NumberRange `json:"retryAfterRange"`
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

	if conf.RegionNrf.PlmnNrfStatusCheckInterval > 0 {
		PlmnNrfStatusCheckInterval = conf.RegionNrf.PlmnNrfStatusCheckInterval
	}

	EnableNotification = conf.NFStatusNotify.Enable

	OverloadProtection = conf.OverloadProtection

	if conf.OverloadControl != nil {
		if overloadLogChangeHandler != nil {
			overloadLogChangeHandler(conf.OverloadControl.EnableDebugLog)
		}
		IsSwitchPlmnForOverloadResp = conf.OverloadControl.IsSwitchPlmnForOverloadResp
		if conf.OverloadControl.ThrottlePercentageForOverloadResp > 0 {
			ThrottlePercentageForOverloadResp = conf.OverloadControl.ThrottlePercentageForOverloadResp
		}

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
	fmt.Printf("OverloadControlEnabled = %v, ", OverloadProtection.Enabled)
	fmt.Printf("OverloadControlLevel = %v, OverloadTriggerLatencyThreshold = %v ms, OverloadControlLatencyThreshold = %v ms, OverloadTriggerSampleWindow = %v, OverloadControlSampleWindow = %v, OverloadTriggerTimeSampleWindow = %v, OverloadControlTimeSampleWindow = %v, WorkerNumber = %v, WorkerQueueCapacity = %v, DeniedWorkerNumber = %v, DeniedWorkerQueueCapacity = %v, IdleInterval = %v ms, IdleRecoverRatio = %v, CounterReportInterval = %v ms, OverloadAlarmClearWindow = %v\n", OverloadProtection.OverloadControlLevel, OverloadProtection.OverloadTriggerLatencyThreshold, OverloadProtection.OverloadControlLatencyThreshold, OverloadProtection.OverloadTriggerSampleWindow, OverloadProtection.OverloadControlSampleWindow, OverloadProtection.OverloadTriggerTimeSampleWindow, OverloadProtection.OverloadControlTimeSampleWindow, OverloadProtection.WorkerNumber, OverloadProtection.WorkerQueueCapacity, OverloadProtection.DeniedWorkerNumber, OverloadProtection.DeniedWorkerQueueCapacity, OverloadProtection.IdleInterval, OverloadProtection.IdleRecoverRatio, OverloadProtection.CounterReportInterval, OverloadProtection.OverloadAlarmClearWindow)
	for _, engine := range OverloadProtection.WorkEngines {
		fmt.Printf("GroupPriority = %v, ", engine.GroupPriority)
	}

	fmt.Printf("RetryAfterRange Start : %v\n", RetryAfterRangeStart)
	fmt.Printf("RetryAfterRange End : %v\n", RetryAfterRangeEnd)

	fmt.Printf("The HeartBeatTimerMin is %v", HeartBeatTimerMin)
	fmt.Printf("The HeartBeatTimerOffset is %v", HeartBeatTimerOffset)
}

// RegistOverloadLevelChangeHandler
func RegistOverloadLevelChangeHandler(f OverloadLogChangeHandler) {
	overloadLogChangeHandler = f
}
