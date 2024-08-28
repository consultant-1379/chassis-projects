package http2

import (
	"net/http"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"os"
)

const (
	Normal = iota
	Overload
)

var (
	EnableOverloadDebugLog  = false
	OverloadRetryAfterStart = 3
	OverloadRetryAfterEnd   = 10
	WorkMode                = os.Getenv("WORK_MODE")
)

type Task interface {
	Execute() int64       // The return value is the duration of request handling by nanoseconds
	EnterTime() time.Time // The return value is the timestamp when the requests arrive in the queue of work engine
	ServerConnection() *serverConn
	ResponseWriter() *responseWriter
	Request() *http.Request
	GourpPriority() int
}

type MicroService interface {
	MonitorTrafficLatency(manager *WorkEngineManager)
}

type TrafficLatency struct {
	WaitTime      int64
	ProcessTime   int64
	GroupPriority int
	Count         uint64
}

type WorkEngine struct {
	GroupPriority        int
	MessagePriorityStart int
	MessagePriorityEnd   int
	CurrentOverloadLevel int64
	RandOverload         chan int64
}

type WorkEngineManager struct {
	OverloadControlLevel            int64
	OverloadTriggerLatencyThreshold float64
	OverloadControlLatencyThreshold float64
	OverloadTriggerSampleWindow     uint64
	OverloadControlSampleWindow     uint64
	OverloadTriggerTimeSampleWindow uint64
	OverloadControlTimeSampleWindow uint64
	IdleInterval                    uint64
	IdleRecoverRatio                int64
	DefaultMessagePriority          int
	CounterReportInterval           uint64
	OverloadAlarmClearWindow        uint64
	ProcessRequestWorkerNumber      int
	ProcessRequestQueueCapacity     int
	DeniedRequestWorkerNumber       int
	DeniedRequestQueueCapacity      int
	StatisticsQueueCapacity         int

	Status                 int
	ProcessedRequestNumber uint64
	DeniedRequestNumber    uint64
	LocaltimeInSecond      int64
	WorkEnginePrioritys    []int

	WorkEngines     []*WorkEngine
	Statistics      chan *TrafficLatency
	ProcessRequests chan Task
	DeniedRequests  chan Task

	TotalCounterName  string
	DeniedCounterName string

	ServiceWorkMode MicroService

	OverloadRedirectEnabled bool
	OverloadRedirectAddr    string
}

func logInfo(format string, args ...interface{}) {

	if EnableOverloadDebugLog {
		log.Errorf(format, args...)
	}
}

func SetEnableOverloadDebugLog(enable bool) {
	EnableOverloadDebugLog = enable
}

func SetOverloadRetryAfter(start, end int) {
	OverloadRetryAfterStart = start
	OverloadRetryAfterEnd = end
}
