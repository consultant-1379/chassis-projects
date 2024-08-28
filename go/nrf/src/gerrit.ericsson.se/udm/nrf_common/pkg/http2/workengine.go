package http2

import (
	"time"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

func NewWorkEngine(groupPriority, messagePriorityStart, messagePriorityEnd, queueCapacity, workerNumber int, overloadControlLevel uint64) *WorkEngine {

	engine := new(WorkEngine)

	engine.GroupPriority = groupPriority
	engine.MessagePriorityStart = messagePriorityStart
	engine.MessagePriorityEnd = messagePriorityEnd
	engine.QueueCapacity = queueCapacity
	engine.WorkerNumber = workerNumber
	engine.CurrentOverloadLevel = overloadControlLevel

	engine.Workers = make([]Worker, workerNumber)
	engine.RequestQueue = make(chan Task, queueCapacity)

	log.Warningf("Group Priority %v [%v , %v] : Queue Capacity = %v, Worker Number = %v", engine.GroupPriority, engine.MessagePriorityStart, engine.MessagePriorityEnd, engine.QueueCapacity, engine.WorkerNumber)
	return engine
}

func (engine *WorkEngine) Start() {

	for i, _ := range engine.Workers {
		engine.Workers[i].Running = true
		go engine.Workers[i].Handle(engine)
	}
}

func (engine *WorkEngine) Stop() {

	for i, _ := range engine.Workers {
		engine.Workers[i].Running = false
	}
}

func (worker *Worker) Handle(engine *WorkEngine) {

	for worker.Running {
		select {
		case task := <-engine.RequestQueue:

			enterTime := task.EnterTime()
			arrivalTime := time.Now()
			waitTime := arrivalTime.Sub(enterTime).Nanoseconds()

			if EngineManager.Status == Overload && float64(waitTime)/1000000.0 >= EngineManager.OverloadControlLatencyThreshold {

				deniedTask := new(HttpTask)
				deniedTask.sc = task.ServerConnection()
				deniedTask.rw = task.ResponseWriter()
				deniedTask.req = task.Request()
				deniedTask.handler = overloadHandler
				deniedTask.enterTime = time.Now()

				if EngineManager.PushDeniedRequest(deniedTask) == false {
					go task.ServerConnection().runHandler(task.ResponseWriter(), task.Request(), overloadHandler)
				}
				continue
			}

			latency := new(TrafficLatency)
			latency.WaitTime = waitTime
			latency.ProcessTime = task.Execute()
			latency.GroupPriority = engine.GroupPriority
			latency.Count = 1
			EngineManager.PushTrafficLatency(latency)
		}
	}
}

func (engine *WorkEngine) Push(task Task) bool {

	if len(engine.RequestQueue) >= engine.QueueCapacity {
		log.Warningf("Work Engine [%v, %v] is full", engine.MessagePriorityStart, engine.MessagePriorityEnd)
		return false
	}

	if EngineManager.Status == Overload {
		enterTime := uint64(task.EnterTime().Nanosecond())
		if enterTime&(EngineManager.OverloadControlLevel-1) >= engine.CurrentOverloadLevel {
			return false
		}
	}

	engine.RequestQueue <- task
	return true
}
