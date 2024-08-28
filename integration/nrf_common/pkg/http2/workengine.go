package http2

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"math/rand"
	"time"
)

func NewWorkEngine(groupPriority, messagePriorityStart, messagePriorityEnd int, overloadControlLevel int64) *WorkEngine {

	engine := new(WorkEngine)

	engine.GroupPriority = groupPriority
	engine.MessagePriorityStart = messagePriorityStart
	engine.MessagePriorityEnd = messagePriorityEnd
	engine.CurrentOverloadLevel = overloadControlLevel

	engine.RandOverload = make(chan int64, 512)
	go func() {
		randOL := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			if len(engine.RandOverload) < cap(engine.RandOverload) {
				engine.RandOverload <- randOL.Int63n(EngineManager.OverloadControlLevel)
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	log.Warningf("Group Priority %v [%v , %v]", engine.GroupPriority, engine.MessagePriorityStart, engine.MessagePriorityEnd)
	return engine
}

func (engine *WorkEngine) Push(task Task) bool {
	//if EngineManager.Status == Overload {
	//	enterTime := int64(task.EnterTime().Nanosecond())
	//	if enterTime&(EngineManager.OverloadControlLevel-1) >= engine.CurrentOverloadLevel {
	//		return false
	//	}
	//}
	if EngineManager.Status == Overload {
		value := <-engine.RandOverload
		if value >= engine.CurrentOverloadLevel {
			return false
		}
	}
	return true
}
