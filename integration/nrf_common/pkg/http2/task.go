package http2

import (
	"fmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"net/http"
	"strconv"
	"time"
)

var EngineManager *WorkEngineManager = nil

type HttpTask struct {
	gourpPriority int
	sc            *serverConn
	rw            *responseWriter
	req           *http.Request
	handler       func(http.ResponseWriter, *http.Request)
	enterTime     time.Time
}

func (task *HttpTask) getExtraLatency(req *http.Request) int64 {
	var latency int64 = 0
	if vlist, ok := req.Header[constvalue.HeaderExtraLatency]; ok {
		for _, value := range vlist {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				continue
			}
			latency += intValue
		}
	}
	return latency
}

func (task *HttpTask) Execute() int64 {
	start := time.Now()
	task.sc.runHandler(task.rw, task.req, task.handler)
	return time.Since(start).Nanoseconds() - task.getExtraLatency(task.req)
}

func (task *HttpTask) EnterTime() time.Time {
	return task.enterTime
}

func (task *HttpTask) ServerConnection() *serverConn {
	return task.sc
}

func (task *HttpTask) ResponseWriter() *responseWriter {
	return task.rw
}

func (task *HttpTask) Request() *http.Request {
	return task.req
}

func (task *HttpTask) GourpPriority() int {
	return task.gourpPriority
}

// Run on its own goroutine.
func overloadHandler(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "application/problem+json")

	retryAfter := fmt.Sprintf("%d", time.Now().Nanosecond()&(OverloadRetryAfterEnd-OverloadRetryAfterStart)+OverloadRetryAfterStart)
	rw.Header().Set("Retry-After", retryAfter)

	if WorkMode == constvalue.APP_WORKMODE_NRF_DISC && EngineManager.OverloadRedirectEnabled {
		rw.Header().Set("Location", EngineManager.OverloadRedirectAddr+req.RequestURI)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}
	rw.Write([]byte(fmt.Sprintf(`{"title": "%s"}`, "Service Overload")))
}

func Process(sc *serverConn, rw *responseWriter, req *http.Request, handler func(http.ResponseWriter, *http.Request)) {

	if EngineManager == nil {
		go sc.runHandler(rw, req, handler)
		return
	}

	priority := EngineManager.DefaultMessagePriority
	priorityStr := req.Header.Get("3gpp-Sbi-Message-Priority")
	if priorityStr != "" {
		value, err := strconv.Atoi(priorityStr)
		if err == nil {
			priority = value
		}
	}

	// set priority 0 (0 has the highest priority) if the method is "DELETE"
	if req.Method == "DELETE" {
		priority = 0
	}

	task := new(HttpTask)
	task.sc = sc
	task.rw = rw
	task.req = req
	task.handler = handler
	task.enterTime = time.Now()
	engine := EngineManager.GetWorkEngine(priority)
	if nil == engine || engine.Push(task) == false {
		task.gourpPriority = 0
		task.handler = overloadHandler
		if EngineManager.PushDeniedRequest(task) == false {
			//drop this request if deny queue full
			//go sc.runHandler(rw, req, handler)
		}
	} else {
		task.gourpPriority = engine.GroupPriority
		if EngineManager.PushProcessRequest(task) == false {
			task.handler = overloadHandler
			if EngineManager.PushDeniedRequest(task) == false {
				//drop this request if deny queue full
				//go sc.runHandler(rw, req, handler)
			}
		}
	}
}
