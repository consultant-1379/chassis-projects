package http2

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
	"sync/atomic"
	"time"
)

type DiscService struct {
	idleCount               uint64
	sampleCount             uint64
	deniedCount             uint64
	previousDeniedCount     uint64
	waitTime                int64
	processTime             int64
	lastAverageResponseTime float64
	exceedCount             float64
	olExceedCount           int64
	millisecond             uint64

	overloadAlarmSentFlag bool

	count   map[int]uint64
	manager *WorkEngineManager
}

func (service *DiscService) init(manager *WorkEngineManager) {
	if manager == nil {
		return
	}

	service.idleCount = 0
	service.sampleCount = 0
	service.deniedCount = 0
	service.previousDeniedCount = 0
	service.waitTime = 0
	service.processTime = 0
	service.lastAverageResponseTime = 0
	service.exceedCount = 0
	service.olExceedCount = 0
	service.millisecond = 1000000

	service.overloadAlarmSentFlag = false

	service.manager = manager
	service.count = make(map[int]uint64)
	for _, engine := range service.manager.WorkEngines {
		service.count[engine.GroupPriority] = 0
	}
}

func (service *DiscService) clearPeriodData() {
	service.sampleCount = 0
	service.previousDeniedCount = service.deniedCount
	service.waitTime = 0
	service.processTime = 0
	for _, engine := range service.manager.WorkEngines {
		service.count[engine.GroupPriority] = 0
	}
}

func (service *DiscService) onLatency(latency *TrafficLatency) {
	if latency == nil {
		return
	}

	service.sampleCount += latency.Count
	service.waitTime += latency.WaitTime
	service.processTime += latency.ProcessTime
	service.count[latency.GroupPriority] += latency.Count
}

func (service *DiscService) calcAverageResponseTime() float64 {
	if service.sampleCount == 0 {
		return 0
	}
	return float64(service.waitTime+service.processTime) / float64(service.sampleCount*service.millisecond)
}

func (service *DiscService) printWorkloadInfo(status string, avRspTime float64, deltaDenied uint64) {
	averageWaitTime := float64(service.waitTime) / float64(service.sampleCount*service.millisecond)
	pendingLength := len(service.manager.Statistics)
	pendingRequestNumber := service.manager.GetPendingRequestNumber()
	requestStatistics := service.manager.GetRequestStatistics(service.count)

	log.Warningf("%v, AveLatency=%5.2f ms, secCount=%4d[%v], deltaDenied:%d, AllRate %v exceed=%3.1f, olexceed=%2d, AveWaitTime=%.2f ms, PendingStat=%v, PendingRequest=%v", status, avRspTime, service.sampleCount, requestStatistics, deltaDenied, service.manager.GetPassRate(), service.exceedCount, service.olExceedCount, averageWaitTime, pendingLength, pendingRequestNumber)
}

func (service *DiscService) onTicker() {
	service.deniedCount = atomic.LoadUint64(&service.manager.DeniedRequestNumber)
	deltaDenied := service.deniedCount - service.previousDeniedCount
	if deltaDenied > 0 {
		pm.Add(float64(deltaDenied), service.manager.TotalCounterName, constvalue.NfDiscovery, "-", "-", "-")
		pm.Add(float64(deltaDenied), service.manager.DeniedCounterName, constvalue.NfDiscovery, "-", "-", "-", "503", "NF_CONGESTION")
	}

	if service.manager.Status == Normal {
		service.olExceedCount = 0
		if service.sampleCount >= service.manager.OverloadTriggerSampleWindow {
			averageResponseTime := service.calcAverageResponseTime()
			service.printWorkloadInfo("Normal", averageResponseTime, deltaDenied)
			if service.sampleCount <= service.manager.OverloadControlSampleWindow {
				if averageResponseTime >= service.manager.OverloadTriggerLatencyThreshold {
					if averageResponseTime >= service.lastAverageResponseTime {
						service.exceedCount += 1.0
					} else {
						service.exceedCount = 0
					}

					if service.exceedCount > 5 {
						//enter overload status
						service.manager.TriggerOverload(service.count, 0.25, service.sampleCount)
						service.exceedCount = 0
					}
				} else {
					service.exceedCount = 0
				}
			} else {
				if averageResponseTime >= service.manager.OverloadTriggerLatencyThreshold {
					service.exceedCount += float64(service.sampleCount) / float64(service.manager.OverloadControlSampleWindow)
				} else {
					service.exceedCount = 0
				}

				if service.exceedCount > 5 {
					//enter overload status
					changeRate := service.getIncreaceRate(service.sampleCount)
					service.manager.TriggerOverload(service.count, 0.20*changeRate, service.sampleCount)
					service.exceedCount = 0
				}
			}
			service.lastAverageResponseTime = averageResponseTime
		}
	} else if service.manager.Status == Overload {
		averageResponseTime := service.calcAverageResponseTime()
		service.printWorkloadInfo("Overload", averageResponseTime, deltaDenied)
		if service.sampleCount < service.manager.OverloadTriggerSampleWindow {
			if averageResponseTime < service.manager.OverloadControlLatencyThreshold {
				service.olExceedCount--
				if service.olExceedCount <= -3 {
					service.decreaceOverloadRatio(service.sampleCount, averageResponseTime)
					service.olExceedCount = 0
				}
			}
		} else {
			if averageResponseTime < service.manager.OverloadControlLatencyThreshold {
				service.olExceedCount--
				if service.olExceedCount <= -3 {
					service.decreaceOverloadRatio(service.sampleCount, averageResponseTime)
					service.olExceedCount = 0
				}
			} else if averageResponseTime >= service.manager.OverloadTriggerLatencyThreshold {
				service.olExceedCount++
				if service.olExceedCount >= 3 {
					service.increaceOverloadRatio(service.sampleCount, averageResponseTime)
					service.olExceedCount = 0
				}
			}
		}
	} else {
		log.Warningf("SHOULD NOT HAPPEN, invalid status")
	}
	service.clearPeriodData()

	if service.overloadAlarmSentFlag == false && service.manager.Status == Overload {
		fm.SendNRFDiscoveryOverloadAlarm(true)
		service.overloadAlarmSentFlag = true
	} else if service.overloadAlarmSentFlag == true && service.manager.Status == Normal {
		fm.SendNRFDiscoveryOverloadAlarm(false)
		service.overloadAlarmSentFlag = false
	}
}

func (service *DiscService) getIncreaceRate(splCount uint64) float64 {
	var ratio float64
	if splCount < service.manager.OverloadTriggerSampleWindow {
		ratio = 0
	} else {
		if splCount <= service.manager.OverloadControlSampleWindow/2 {
			ratio = 0.5
		} else if splCount <= service.manager.OverloadControlSampleWindow*2 {
			ratio = float64(splCount) / float64(service.manager.OverloadControlSampleWindow)
		} else {
			ratio = 2.0
		}
	}
	return ratio
}

func (service *DiscService) getDecreaceRate(splCount uint64) float64 {
	var ratio float64
	if splCount < service.manager.OverloadTriggerSampleWindow {
		if splCount <= service.manager.OverloadTriggerSampleWindow/2 {
			ratio = 8.0
		} else {
			ratio = float64(service.manager.OverloadTriggerSampleWindow) / float64(splCount) * 2.0
		}
	} else {
		if splCount > service.manager.OverloadControlSampleWindow*2 {
			ratio = 0.5
		} else if splCount > service.manager.OverloadControlSampleWindow/2 {
			ratio = float64(service.manager.OverloadControlSampleWindow) / float64(splCount)
		} else {
			ratio = 2.0
		}
	}
	return ratio
}

func (service *DiscService) increaceOverloadRatio(splCount uint64, averRespTime float64) {
	if splCount == 0 {
		log.Error("SHOULD NOT HAPPEN, invalid splCount")
		return
	}

	changeRate := service.getIncreaceRate(splCount)
	diff := averRespTime - service.manager.OverloadTriggerLatencyThreshold
	ratio := 0.0
	if diff >= 40.0 {
		ratio = 1.0 / 8.0
	} else if diff >= 4.0 {
		ratio = 1.0 / 16.0
	} else if diff >= 3.0 {
		ratio = 1.0 / 32.0
	} else if diff >= 2.0 {
		ratio = 1.0 / 64.0
	} else if diff >= 1.0 {
		ratio = 1.0 / 128.0
	} else {
		ratio = 1.0 / 256.0
	}
	service.manager.IncreaseOverload(service.count, ratio*changeRate, splCount)
}

func (service *DiscService) decreaceOverloadRatio(splCount uint64, averRespTime float64) {
	if splCount == 0 {
		service.manager.RecoverFromOverload()
		return
	}
	changeRate := service.getDecreaceRate(splCount)
	diff := service.manager.OverloadControlLatencyThreshold - averRespTime
	ratio := 0.0
	if diff >= 4.0 {
		ratio = 1.0 / 16.0 //32.0
	} else if diff >= 3.0 {
		ratio = 1.0 / 32.0 //64.0
	} else if diff >= 2.0 {
		ratio = 1.0 / 64.0 //128.0
	} else if diff >= 1.0 {
		ratio = 1.0 / 128.0 //256.0
	} else {
		ratio = 1.0 / 256.0 //512.0
	}
	service.manager.DecreaseOverload(service.count, ratio*changeRate, splCount)
}

func (service *DiscService) MonitorTrafficLatency(manager *WorkEngineManager) {
	if manager == nil {
		log.Error("manager is nil. exit overload control monitor")
		return
	}
	service.init(manager)

	olTicker := time.NewTicker(time.Second * 1)
	for {
		select {
		case latency := <-manager.Statistics:
			service.onLatency(latency)
			select {
			case <-olTicker.C:
				service.onTicker()
			default:
				service.idleCount = 0
			}
		case <-olTicker.C:
			service.onTicker()
			select {
			case latency := <-manager.Statistics:
				service.onLatency(latency)
			default:
				service.idleCount = 0
			}
		}
	}
}
