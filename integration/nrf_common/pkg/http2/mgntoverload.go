package http2

import (
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

type MgntService struct {
}

func (service *MgntService) MonitorTrafficLatency(manager *WorkEngineManager) {

	var interval uint64 = 100
	var idleCount uint64 = 0
	var sampleCount uint64 = 0
	var waitTime int64 = 0
	var processTime int64 = 0
	var millisecond uint64 = 1000000
	var exceed uint64 = 0
	var limit uint64 = 1
	count := make(map[int]uint64)
	for _, engine := range manager.WorkEngines {
		count[engine.GroupPriority] = 0
	}

	lastTimeInSecond := manager.getLocalTime()
	minSampleCount := uint64(10) * manager.OverloadTriggerTimeSampleWindow
	for {
		select {
		case latency := <-manager.Statistics:
			idleCount = 0
			sampleCount += latency.Count
			waitTime += latency.WaitTime
			processTime += latency.ProcessTime
			count[latency.GroupPriority] += latency.Count

			if manager.Status == Normal {
				if sampleCount >= minSampleCount && (manager.getLocalTime()-lastTimeInSecond) >= int64(manager.OverloadTriggerTimeSampleWindow) {
					averageWaitTime := float64(waitTime) / float64(sampleCount*millisecond)
					averageResponseTime := float64(waitTime+processTime) / float64(sampleCount*millisecond)
					pendingLength := len(manager.Statistics)
					pendingRequestNumber := manager.GetPendingRequestNumber()
					requestStatistics := manager.GetRequestStatistics(count)
					if averageResponseTime > manager.OverloadTriggerLatencyThreshold {
						if exceed >= limit || averageResponseTime >= manager.OverloadTriggerLatencyThreshold*3/2 {
							exceed = 0
							manager.TriggerOverload(count, 0.325, sampleCount)
						} else {
							exceed++
						}
						log.Warningf("Normal, Overload Trigger Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms , Pending Latency Statistics = %v , Pending Request Number = %v", sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)
					} else if exceed > 0 {
						exceed--
					}

					sampleCount = 0
					waitTime = 0
					processTime = 0
					for _, engine := range manager.WorkEngines {
						count[engine.GroupPriority] = 0
					}
					lastTimeInSecond = manager.getLocalTime()
				}
			} else if manager.Status == Overload {
				if (manager.getLocalTime() - lastTimeInSecond) >= int64(manager.OverloadControlTimeSampleWindow) {
					averageWaitTime := float64(waitTime) / float64(sampleCount*millisecond)
					averageResponseTime := float64(waitTime+processTime) / float64(sampleCount*millisecond)
					pendingLength := len(manager.Statistics)
					pendingRequestNumber := manager.GetPendingRequestNumber()
					requestStatistics := manager.GetRequestStatistics(count)
					log.Warningf("Overload, Overall Pass Pate %v Overload Control Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms, Pending Latency Statistics = %v , Pending Request Number = %v", manager.GetPassRate(), sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)
					if averageResponseTime > manager.OverloadControlLatencyThreshold {
						// Increase Overload
						manager.IncreaseOverloadForMgm(count, averageResponseTime, sampleCount)
					} else {
						// Decrease Overload
						manager.DecreaseOverloadForMgm(count, averageResponseTime, sampleCount)
					}
					sampleCount = 0
					waitTime = 0
					processTime = 0
					for _, engine := range manager.WorkEngines {
						count[engine.GroupPriority] = 0
					}
					lastTimeInSecond = manager.getLocalTime()
				}
			} else {
				idleCount = 0
				sampleCount = 0
				waitTime = 0
				processTime = 0
				for _, engine := range manager.WorkEngines {
					count[engine.GroupPriority] = 0
				}
				log.Warningf("SHOULD NOT HAPPEN, invalid status")
			}
		default:
			idleCount++
			time.Sleep(time.Duration(interval) * time.Millisecond)

			if idleCount == manager.IdleInterval/interval {

				if manager.Status == Normal {
					// Do Nothing
				} else if manager.Status == Overload {
					// Need to recover work engine
					manager.RecoverFromOverload()
				} else {
					log.Warningf("SHOULD NOT HAPPEN, invalid status")
				}

				idleCount = 0
				sampleCount = 0
				waitTime = 0
				processTime = 0
				for _, engine := range manager.WorkEngines {
					count[engine.GroupPriority] = 0
				}
			}

		}
	}
}

func (manager *WorkEngineManager) DecreaseOverloadForMgm(count map[int]uint64, averageResponseTime float64, window uint64) {

	for _, engine := range manager.WorkEngines {
		if engine.CurrentOverloadLevel < manager.OverloadControlLevel {
			passRate := float64(engine.CurrentOverloadLevel) / float64(manager.OverloadControlLevel) * 100.0
			ratio := manager.calculateDecreaceOverloadRatio(averageResponseTime, passRate)
			step := int64(float64(manager.OverloadControlLevel) * ratio)
			if engine.CurrentOverloadLevel+step <= manager.OverloadControlLevel {
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

func (manager *WorkEngineManager) IncreaseOverloadForMgm(count map[int]uint64, averageResponseTime float64, window uint64) {

	for i := len(manager.WorkEngines) - 1; i >= 0; i-- {
		groupPriority := manager.WorkEngines[i].GroupPriority
		if count[groupPriority] != 0 {
			passRate := float64(manager.WorkEngines[i].CurrentOverloadLevel) / float64(manager.OverloadControlLevel) * 100.0
			ratio := manager.calculateIncreaceOverloadRatio(averageResponseTime, passRate)
			step := int64(float64(manager.OverloadControlLevel) * ratio)
			if manager.WorkEngines[i].CurrentOverloadLevel < step {
				step -= manager.WorkEngines[i].CurrentOverloadLevel
				manager.WorkEngines[i].CurrentOverloadLevel = 0
			} else {
				manager.WorkEngines[i].CurrentOverloadLevel -= step
				return
			}
		}
	}
	for _, engine := range manager.WorkEngines {
		groupPriority := engine.GroupPriority
		if count[groupPriority] != 0 {
			engine.CurrentOverloadLevel = manager.OverloadControlLevel / 50
			break
		}
	}
	manager.Status = Overload
}

func (manager *WorkEngineManager) calculateIncreaceOverloadRatio(averageResponseTime float64, passRate float64) float64 {
	ratio := 0.0
	times := averageResponseTime / manager.OverloadControlLatencyThreshold
	if times >= 5.0 {
		if passRate >= 60.0 {
			ratio = 1.0 / 8.0
		} else if passRate >= 20.0 {
			ratio = 1.0 / 8.0
		} else if passRate >= 5.0 {
			ratio = 1.0 / 16.0
		} else {
			ratio = 1.0 / 64.0
		}
	} else if times >= 1.5 {
		if passRate >= 60.0 {
			ratio = 1.0 / 16.0
		} else if passRate >= 20.0 {
			ratio = 1.0 / 32.0
		} else if passRate >= 5.0 {
			ratio = 1.0 / 64.0
		} else {
			ratio = 1.0 / 128.0
		}
	} else {
		if passRate >= 60.0 {
			ratio = 1.0 / 32.0
		} else if passRate >= 20.0 {
			ratio = 1.0 / 64.0
		} else if passRate >= 5.0 {
			ratio = 1.0 / 128.0
		} else {
			ratio = 1.0 / 256.0
		}
	}
	log.Warningf("Overload, times is %v, passRate is %v, ratio is %v", times, passRate, ratio)
	return ratio
}

func (manager *WorkEngineManager) calculateDecreaceOverloadRatio(averageResponseTime float64, passRate float64) float64 {
	ratio := 0.0
	times := averageResponseTime / manager.OverloadControlLatencyThreshold

	if passRate >= 60.0 {
		ratio = 1.0 / 32.0
	} else if passRate >= 20.0 {
		ratio = 1.0 / 32.0
	} else if passRate >= 5.0 {
		ratio = 1.0 / 64.0
	} else {
		ratio = 1.0 / 128.0
	}

	log.Warningf("Overload, times is %v, passRate is %v, ratio is %v", times, passRate, ratio)
	return ratio
}
