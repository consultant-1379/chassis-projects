package internalconf

import (
	"encoding/json"
	"io/ioutil"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

const (
	defaultHomeNrfForwardRetryTime     = 2   //second
	defaultHomeNrfForwardRetryWaitTime = 200 //Millisecond

	defaultHttpServerIdleTimeout   = 1
	defaultHttpServerActiveTimeout = 2
	defaultHTTPWithXVersion        = false

	defaultDbproxyConnectionNum  = 100
	defaultDbproxyGrpcCtxTimeout = 3
)

var (
	// HomeNrfForwardRetryTime identifies how many times to retry after
	// forwarding to home NRF failed (in second)
	HomeNrfForwardRetryTime = defaultHomeNrfForwardRetryTime

	// HomeNrfForwardRetryWaitTime identifies how long to wait to retry (in Millisecond)
	HomeNrfForwardRetryWaitTime = defaultHomeNrfForwardRetryWaitTime

	// DbproxyConnectionNum idetifies how many connections is setup with dbproxy
	DbproxyConnectionNum = defaultDbproxyConnectionNum

	// DbproxyGrpcCtxTimeout idetifies timeout (in second) for grpc context
	DbproxyGrpcCtxTimeout = defaultDbproxyGrpcCtxTimeout

	// HTTPWithXVersion is to control whether include X-Version in the http header
	HTTPWithXVersion = defaultHTTPWithXVersion

	// OverloadProtection is to control the service quality under high traffic load
	OverloadProtection OverloadProtectionConf
)

// InternalConf is implemented by :
// InternalMgmtConf, InternalDiscConf, InternalProvConf, InternalNotifyConf
type InternalConf interface {
	ParseConf()
	ShowInternalConf()
}

// InternalNrfConf is format of internalConf
type InternalNrfConf struct {
	FileName     string
	ConfInstance InternalConf
}

// HomeNrfConf is format of object "homeNrf" in internalConf.json
type HomeNrfConf struct {
	Forward ForwardConf `json:"forward"`
}

// ForwardConf is format of object "forward" in internalConf.json
type ForwardConf struct {
	RetryTime     int `json:"retryTime"`
	RetryWaitTime int `json:"retryWaitTime"`
}

// HttpServerConf is format of object "httpServer" in internalConf.json
type HttpServerConf struct {
	IdleTimeout      int  `json:"idleTimeout"`
	ActiveTimeout    int  `json:"activeTimeout"`
	HTTPWithXVersion bool `json:"httpWithXVersion"`
}

// DbproxyConf is format of object "dbProxy" in internalConf.json
type DbproxyConf struct {
	ConnectionNum  int `json:"connectionNum"`
	GrpcCtxTimeout int `json:"grpcContextTimeout"`
}

// OverloadProtection is to control the service quality under high traffic load
type OverloadProtectionConf struct {
	Enabled                         bool         `json:"Enabled"`
	OverloadControlLevel            int64        `json:"OverloadControlLevel"`
	OverloadTriggerLatencyThreshold float64      `json:"OverloadTriggerLatencyThreshold"`
	OverloadControlLatencyThreshold float64      `json:"OverloadControlLatencyThreshold"`
	OverloadTriggerSampleWindow     uint64       `json:"OverloadTriggerSampleWindow"`
	OverloadControlSampleWindow     uint64       `json:"OverloadControlSampleWindow"`
	OverloadTriggerTimeSampleWindow uint64       `json:"OverloadTriggerTimeSampleWindow"`
	OverloadControlTimeSampleWindow uint64       `json:"OverloadControlTimeSampleWindow"`
	WorkerNumber                    int          `json:"WorkerNumber"`
	WorkerQueueCapacity             int          `json:"WorkerQueueCapacity"`
	DeniedWorkerNumber              int          `json:"DeniedWorkerNumber"`
	DeniedWorkerQueueCapacity       int          `json:"DeniedWorkerQueueCapacity"`
	IdleInterval                    uint64       `json:"IdleInterval"`
	IdleRecoverRatio                int64        `json:"IdleRecoverRatio"`
	CounterReportInterval           uint64       `json:"CounterReportInterval"`
	OverloadAlarmClearWindow        uint64       `json:"OverloadAlarmClearWindow"`
	WorkEngines                     []WorkEngine `json:"WorkEngine"`
}

// WorkEngine is the message queue + goroutine pool
type WorkEngine struct {
	GroupPriority int `json:"GroupPriority"`
}

// SetFileName set FileName
func (conf *InternalNrfConf) SetFileName(fileName string) {
	conf.FileName = fileName
}

// GetFileName return FileName
func (conf *InternalNrfConf) GetFileName() string {
	return conf.FileName
}

// LoadConf set the value of global variable from configuration file
func (conf *InternalNrfConf) LoadConf() {
	log.Debugf("Start to load config file %s", conf.FileName)
	bytes, err := ioutil.ReadFile(conf.FileName)
	if err != nil {
		log.Warnf("Failed to read config file:%s. %s", conf.FileName, err.Error())
		return
	}

	err = json.Unmarshal(bytes, conf.ConfInstance)
	if err != nil {
		log.Warnf("Failed to unmarshal config file:%s. %s", conf.FileName, err.Error())
		return
	}

	conf.ParseConf()
	log.Debugf("Successfully load config file %s", conf.FileName)

	conf.ShowInternalConf()
}

// ReloadConf reset the value of variable when configuration file changed
func (conf *InternalNrfConf) ReloadConf() {
	log.Debugf("reload config file")
	conf.LoadConf()
}

// ParseConf sets the value of global variable from struct
func (conf *InternalNrfConf) ParseConf() {
	if conf.ConfInstance != nil {
		conf.ConfInstance.ParseConf()
	}
}

// ShowInternalConf shows the value of global variable
func (conf *InternalNrfConf) ShowInternalConf() {
	if conf.ConfInstance != nil {
		conf.ConfInstance.ShowInternalConf()
	}
}
