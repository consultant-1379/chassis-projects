package http2

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var EngineManager *WorkEngineManager = nil

type HttpTask struct {
	sc        *serverConn
	rw        *responseWriter
	req       *http.Request
	handler   func(http.ResponseWriter, *http.Request)
	enterTime time.Time
}

func (task *HttpTask) Execute() int64 {
	start := time.Now()
	task.sc.runHandler(task.rw, task.req, task.handler)
	return time.Since(start).Nanoseconds()
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

// Run on its own goroutine.
func overloadHandler(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Content-Type", "application/problem+json")

	retryAfter := fmt.Sprintf("%d", time.Now().Nanosecond()&(OverloadRetryAfterEnd-OverloadRetryAfterStart)+OverloadRetryAfterStart)
	rw.Header().Set("Retry-After", retryAfter)

	rw.WriteHeader(http.StatusServiceUnavailable)
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

	task := new(HttpTask)
	task.sc = sc
	task.rw = rw
	task.req = req
	task.handler = handler
	task.enterTime = time.Now()
	engine := EngineManager.GetWorkEngine(priority)
	if nil == engine || engine.Push(task) == false {
		task.handler = overloadHandler
		if EngineManager.PushDeniedRequest(task) == false {
			go sc.runHandler(rw, req, handler)
		}
	}
}
