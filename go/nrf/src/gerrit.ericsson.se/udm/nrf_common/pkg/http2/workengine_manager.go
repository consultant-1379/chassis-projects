package http2

import (
	"fmt"
	"sync/atomic"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
)

func NewWorkEngineManager() *WorkEngineManager {
	manager := new(WorkEngineManager)

	manager.Status = Normal
	manager.ProcessedRequestNumber = 0
	manager.DeniedRequestNumber = 0
	manager.DeniedRequestWorkerNumber = 16
	manager.DeniedRequestQueueCapacity = 40960
	manager.StatisticsQueueCapacity = 4096


	if WorkMode == constvalue.APP_WORKMODE_NRF_DISC {
		manager.ServiceWorkMode = &DiscService{}
	} else if WorkMode == constvalue.APP_WORKMODE_NRF_MGMT {
		manager.ServiceWorkMode = &MgntService{}
	} else {
		log.Warningf("unknown work mode:%v", WorkMode)
	}
	return manager
}

func (manager *WorkEngineManager) SetOverloadControlLevel(level uint64) {
	manager.OverloadControlLevel = level
}

func (manager *WorkEngineManager) SetOverloadTriggerLatencyThreshold(value float64) {
	manager.OverloadTriggerLatencyThreshold = value
}

func (manager *WorkEngineManager) SetOverloadControlLatencyThreshold(value float64) {
	manager.OverloadControlLatencyThreshold = value
}

func (manager *WorkEngineManager) SetOverloadTriggerSampleWindow(window uint64) {
	manager.OverloadTriggerSampleWindow = window
}

func (manager *WorkEngineManager) SetOverloadControlSampleWindow(window uint64) {
	manager.OverloadControlSampleWindow = window
}

func (manager *WorkEngineManager) SetOverloadTriggerTimeSampleWindow(window uint64) {
	manager.OverloadTriggerTimeSampleWindow = window
}

func (manager *WorkEngineManager) SetOverloadControlTimeSampleWindow(window uint64) {
	manager.OverloadControlTimeSampleWindow = window
}

func (manager *WorkEngineManager) SetIdleInterval(interval uint64) {
	manager.IdleInterval = interval
}

func (manager *WorkEngineManager) SetIdleRecoverRatio(ratio uint64) {
	manager.IdleRecoverRatio = ratio
}

func (manager *WorkEngineManager) SetDefaultMessagePriority(priority int) {
	manager.DefaultMessagePriority = priority
}

func (manager *WorkEngineManager) SetCounterReportInterval(interval uint64) {
	manager.CounterReportInterval = interval
}

func (manager *WorkEngineManager) SetOverloadAlarmClearWindow(window uint64) {
	manager.OverloadAlarmClearWindow = window
}

func (manager *WorkEngineManager) SetDeniedRequestWorkerNumber(number int) {
	manager.DeniedRequestWorkerNumber = number
}

func (manager *WorkEngineManager) SetDeniedRequestQueueCapacity(capacity int) {
	manager.DeniedRequestQueueCapacity = capacity
}

func (manager *WorkEngineManager) SetStatisticsQueueCapacity(capacity int) {
	manager.StatisticsQueueCapacity = capacity
}

func (manager *WorkEngineManager) RegisterWorkEngine(engine *WorkEngine) bool {
	if engine == nil {
		return false
	}

	manager.WorkEngines = append(manager.WorkEngines, engine)
	return true
}

func (manager *WorkEngineManager) Start() {
	for _, engine := range manager.WorkEngines {
		engine.Start()
	}

	log.Warningf("Overload Control Level = %v , Overload Trigger Latency Threshold = %.2f ms, Overload Control Latency Threshold = %.2f ms, Overload Trigger Sample Window = %v , Overload Control Sample Window = %v , Idle Interval = %v ms, Idle Recover Ratio = %v , Denied Request Worker Number = %v , Denied Request Queue Capacity = %v , Statistics Queue Capacity = %v", manager.OverloadControlLevel, manager.OverloadTriggerLatencyThreshold, manager.OverloadControlLatencyThreshold, manager.OverloadTriggerSampleWindow, manager.OverloadControlSampleWindow, manager.IdleInterval, manager.IdleRecoverRatio, manager.DeniedRequestWorkerNumber, manager.DeniedRequestQueueCapacity, manager.StatisticsQueueCapacity)

	manager.DeniedRequests = make(chan Task, manager.DeniedRequestQueueCapacity)
	manager.Statistics = make(chan *TrafficLatency, manager.StatisticsQueueCapacity)

	go manager.ServiceWorkMode.MonitorTrafficLatency(manager)
	if WorkMode == constvalue.APP_WORKMODE_NRF_MGMT {
		manager.GenerateLocalTime()
	}
	for i := 0; i < manager.DeniedRequestWorkerNumber; i++ {
		go manager.HandleDeniedRequest()
	}
}

func (manager *WorkEngineManager) GetWorkEngine(priority int) *WorkEngine {

	if len(manager.WorkEngines) == 0 {
		return nil
	}

	for _, engine := range manager.WorkEngines {
		if priority >= engine.MessagePriorityStart && priority <= engine.MessagePriorityEnd {
			return engine
		}
	}

	lowest := 0
	for i, _ := range manager.WorkEngines {
		if manager.WorkEngines[i].MessagePriorityEnd > manager.WorkEngines[lowest].MessagePriorityEnd {
			lowest = i
		}
	}
	return manager.WorkEngines[lowest]
}

func (manager *WorkEngineManager) PushTrafficLatency(latency *TrafficLatency) {
	manager.Statistics <- latency
}

func (manager *WorkEngineManager) GetRequestStatistics(count map[int]uint64) string {
	statistics := ""
	for key, value := range count {
		if statistics != "" {
			statistics += ","
		}
		statistics += fmt.Sprintf("GroupPriority %v:%v", key, value)
	}
	return statistics
}

func (manager *WorkEngineManager) GetPassRate() string {
	passRate := ""
	for _, engine := range manager.WorkEngines {
		rate := float64(engine.CurrentOverloadLevel) / float64(manager.OverloadControlLevel) * 100.0
		passRate += fmt.Sprintf("[%v,%v]: %.2f , ", engine.MessagePriorityStart, engine.MessagePriorityEnd, rate)
	}
	return passRate
}

func (manager *WorkEngineManager) GetPendingRequestNumber() int {
	total := 0
	for _, engine := range manager.WorkEngines {
		total += len(engine.RequestQueue)
	}
	return total
}

func (manager *WorkEngineManager) GetTotalConfigLevel() uint64 {
	var total uint64 = 0
	for i := 0; i < len(manager.WorkEngines); i++ {
		total += manager.OverloadControlLevel
	}
	return total
}

func (manager *WorkEngineManager) TriggerOverload(count map[int]uint64, ratio float64, window uint64) {
	step := float64(manager.OverloadControlLevel) * ratio

	for i := len(manager.WorkEngines) - 1; i >= 0; i-- {
		groupPriority := manager.WorkEngines[i].GroupPriority
		if count[groupPriority] != 0 {
			base := (float64(count[groupPriority]) / float64(window)) * float64(manager.OverloadControlLevel)
			capacity := (float64(count[groupPriority]) / float64(window)) * float64(manager.WorkEngines[i].CurrentOverloadLevel)
			if step <= capacity {
				manager.WorkEngines[i].CurrentOverloadLevel -= uint64(step / base * float64(manager.OverloadControlLevel))
				break
			} else {
				manager.WorkEngines[i].CurrentOverloadLevel = 0
				step -= capacity
			}
		}
	}
	manager.Status = Overload
}

func (manager *WorkEngineManager) IncreaseOverload(count map[int]uint64, ratio float64, window uint64) {
	step := uint64(float64(manager.OverloadControlLevel) * ratio)

	for i := len(manager.WorkEngines) - 1; i >= 0; i-- {
		groupPriority := manager.WorkEngines[i].GroupPriority
		if count[groupPriority] != 0 {
			if manager.WorkEngines[i].CurrentOverloadLevel < step {
				step -= manager.WorkEngines[i].CurrentOverloadLevel
				manager.WorkEngines[i].CurrentOverloadLevel = 0
			} else {
				manager.WorkEngines[i].CurrentOverloadLevel -= step
				return
			}
		}
	}
	for i := range manager.WorkEngines {
		groupPriority := manager.WorkEngines[i].GroupPriority
		if count[groupPriority] != 0 {
			manager.WorkEngines[i].CurrentOverloadLevel = manager.OverloadControlLevel / 50
			break
		}
	}
	manager.Status = Overload
}

func (manager *WorkEngineManager) DecreaseOverload(count map[int]uint64, ratio float64, window uint64) {
	step := uint64(float64(manager.OverloadControlLevel) * ratio)

	for _, engine := range manager.WorkEngines {
		if engine.CurrentOverloadLevel < manager.OverloadControlLevel {
			if engine.CurrentOverloadLevel + step <= manager.OverloadControlLevel {
				engine.CurrentOverloadLevel += step
				return
			} else {
				step -= manager.OverloadControlLevel - engine.CurrentOverloadLevel
				engine.CurrentOverloadLevel = manager.OverloadControlLevel
			}
		}
	}
	manager.Status = Normal
	log.Warningf("Decrease Overload to Normal")
}

func (manager *WorkEngineManager) RecoverFromOverload() {
	var step uint64 = manager.GetTotalConfigLevel() / manager.IdleRecoverRatio

	for _, engine := range manager.WorkEngines {
		if engine.CurrentOverloadLevel < manager.OverloadControlLevel {
			if engine.CurrentOverloadLevel + step <= manager.OverloadControlLevel {
				engine.CurrentOverloadLevel += step
				log.Warningf("Recovering, Overall Pass Rate %v", manager.GetPassRate())
				return
			} else {
				step = manager.OverloadControlLevel - engine.CurrentOverloadLevel
				engine.CurrentOverloadLevel = manager.OverloadControlLevel
			}
		}
	}
	manager.Status = Normal
	log.Warningf("Recover from Overload to Normal")
}

func (manager *WorkEngineManager) PushDeniedRequest(task Task) bool {
	if len(manager.DeniedRequests) >= cap(manager.DeniedRequests) {
		log.Warningf("Fail to push denied request because queue is full")
		return false
	}
	manager.DeniedRequests <- task
	return true
}

func (manager *WorkEngineManager) HandleDeniedRequest() {
	for {
		select {
		case task := <-manager.DeniedRequests:
			atomic.AddUint64(&manager.DeniedRequestNumber, 1)
			task.Execute()
		}
	}
}

func (manager *WorkEngineManager) ReportCounterForNFDisc(totalCounterName, unprocessedCounterName string) {
	var previousUnprocessed uint64 = 0
	var currentUnprocessed uint64 = 0
	for {
		currentUnprocessed = manager.DeniedRequestNumber

		delta1 := float64(currentUnprocessed - previousUnprocessed)
		if delta1 > 0 {
			pm.Add(delta1, totalCounterName, constvalue.NfDiscovery, "unknown", "unknown")
			pm.Add(delta1, unprocessedCounterName, constvalue.NfDiscovery, "unknown", "unknown", "503", "NF_CONGESTION")
		}

		logInfo("Report Counter : %v, %v += %v", totalCounterName, unprocessedCounterName, delta1)
		time.Sleep(time.Duration(manager.CounterReportInterval) * time.Millisecond)

		previousUnprocessed = currentUnprocessed
		currentUnprocessed = 0
	}
}

func (manager *WorkEngineManager) ReportCounterForNFMgm(totalCounterName, unprocessedCounterName string) {
	var previousUnprocessed uint64 = 0
	var currentUnprocessed uint64 = 0

	// Alarm related
	var overloadAlarmSentFlag bool = false
	var normalStatusCount uint64 = 0

	for {
		currentUnprocessed = manager.DeniedRequestNumber

		delta1 := float64(currentUnprocessed - previousUnprocessed)
		if delta1 > 0 {
			pm.Add(delta1, totalCounterName, constvalue.NfManagement, "unknown", "unknown")
			pm.Add(delta1, unprocessedCounterName, constvalue.NfManagement, "unknown", "unknown", "503", "Overload")

			normalStatusCount = 0
			//Trigger overload alarm
			if !overloadAlarmSentFlag {
				fm.SendNRFManagementOverloadAlarm(true)
				overloadAlarmSentFlag = true
			}
		} else {
			if overloadAlarmSentFlag {
				//clear overload alarm if condition is met
				normalStatusCount++
				if normalStatusCount >= manager.OverloadAlarmClearWindow {
					normalStatusCount = 0
					fm.SendNRFManagementOverloadAlarm(false)
					overloadAlarmSentFlag = false
				}
			}
		}

		logInfo("Report Counter : %v, %v += %v", totalCounterName, unprocessedCounterName, delta1)

		previousUnprocessed = currentUnprocessed
		currentUnprocessed = 0

		time.Sleep(time.Duration(manager.CounterReportInterval) * time.Millisecond)
	}
}


func (manager *WorkEngineManager) getLocalTime() int64 {
	return manager.LocaltimeInSecond
}

// GenerateLocalTime is to generate the local time every second
func (manager *WorkEngineManager) GenerateLocalTime() {
	ticker := time.NewTicker(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case t := <-ticker.C:
				manager.LocaltimeInSecond = t.Unix()
			}
		}
	}()
}
