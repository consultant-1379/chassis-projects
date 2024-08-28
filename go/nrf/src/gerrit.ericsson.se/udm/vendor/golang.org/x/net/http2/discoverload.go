package http2

import (
	"time"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

type DiscService struct {

}

func (service *DiscService)MonitorTrafficLatency(manager *WorkEngineManager)  {
	var interval uint64 = 100
	var idleCount uint64 = 0
	var sampleCount uint64 = 0
	var waitTime int64 = 0
	var processTime int64 = 0
	var millisecond uint64 = 1000000
	var exceed uint64 = 0
	var limit uint64 = 1
	var quickOverloadCount uint64 = 0
	count := make(map[int]uint64)
	for _, engine := range manager.WorkEngines {
		count[engine.GroupPriority] = 0
	}

	for {
		select {
		case latency := <-manager.Statistics:
			idleCount = 0
			sampleCount += latency.Count
			waitTime += latency.WaitTime
			processTime += latency.ProcessTime
			count[latency.GroupPriority] += latency.Count

			if manager.Status == Normal {
				if sampleCount >= manager.OverloadTriggerSampleWindow {
					averageWaitTime := float64(waitTime)/float64(sampleCount*millisecond)
					averageResponseTime := float64(waitTime+processTime)/float64(sampleCount*millisecond)
					pendingLength := len(manager.Statistics)
					pendingRequestNumber := manager.GetPendingRequestNumber()
					requestStatistics := manager.GetRequestStatistics(count)
					log.Warningf("Normal, Overload Trigger Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms , Pending Latency Statistics = %v , Pending Request Number = %v", sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)
					if averageResponseTime > manager.OverloadTriggerLatencyThreshold {
						if exceed >= limit || averageResponseTime >= manager.OverloadTriggerLatencyThreshold*3/2 {
							exceed = 0
							manager.TriggerOverload(count, 0.325, sampleCount)
						} else {
							exceed++
						}
					} else if exceed > 0{
						exceed--
					}

					sampleCount = 0
					waitTime = 0
					processTime = 0
					for _, engine := range manager.WorkEngines {
						count[engine.GroupPriority] = 0
					}
				} else if sampleCount%(manager.OverloadTriggerSampleWindow/10) == 0 {
					averageResponseTime := float64(waitTime+processTime)/float64(sampleCount*millisecond)
					if averageResponseTime >= manager.OverloadTriggerLatencyThreshold*2.0 {
						quickOverloadCount++
						if quickOverloadCount >= 2 {
							averageWaitTime := float64(waitTime)/float64(sampleCount*millisecond)
							pendingLength := len(manager.Statistics)
							pendingRequestNumber := manager.GetPendingRequestNumber()
							requestStatistics := manager.GetRequestStatistics(count)
							log.Warningf("Normal, Overload Trigger Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms , Pending Latency Statistics = %v , Pending Request Number = %v", sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)
							manager.TriggerOverload(count, 0.325, sampleCount)

							sampleCount = 0
							waitTime = 0
							processTime = 0
							for _, engine := range manager.WorkEngines {
								count[engine.GroupPriority] = 0
							}
						}

					} else {
						quickOverloadCount = 0
					}
				}
			} else if manager.Status == Overload {
				if sampleCount >= manager.OverloadControlSampleWindow {
					averageWaitTime := float64(waitTime)/float64(sampleCount*millisecond)
					averageResponseTime := float64(waitTime+processTime)/float64(sampleCount*millisecond)
					pendingLength := len(manager.Statistics)
					pendingRequestNumber := manager.GetPendingRequestNumber()
					requestStatistics := manager.GetRequestStatistics(count)
					log.Warningf("Overload, Overall Pass Pate %v Overload Control Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms, Pending Latency Statistics = %v , Pending Request Number = %v", manager.GetPassRate(), sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)
					if averageResponseTime > manager.OverloadControlLatencyThreshold {
						// Increase Overload
						diff := averageResponseTime - manager.OverloadControlLatencyThreshold
						ratio := 0.0
						if diff >= 20.0 {
							ratio = 1.0/16.0
						} else if diff >= 4.0 {
							ratio = 1.0/32.0
						} else if diff >= 3.0 {
							ratio = 1.0/64.0
						} else if diff >= 2.0 {
							ratio = 1.0/128.0
						} else if diff >= 1.0 {
							ratio = 1.0/256.0
						} else {
							ratio = 1.0/512.0
						}
						manager.IncreaseOverload(count, ratio, sampleCount)
					} else {
						// Decrease Overload
						diff := manager.OverloadControlLatencyThreshold - averageResponseTime
						ratio := 0.0
						if diff >= 4.0 {
							ratio = 1.0/32.0
						} else if diff >= 3.0 {
							ratio = 1.0/256.0
						} else if diff >= 1.0 {
							ratio = 1.0/1024.0
						}
						manager.DecreaseOverload(count, ratio, sampleCount)
					}
					sampleCount = 0
					waitTime = 0
					processTime = 0
					for _, engine := range manager.WorkEngines {
						count[engine.GroupPriority] = 0
					}
				} else if sampleCount%(manager.OverloadControlSampleWindow/10) == 0 {
					averageResponseTime := float64(waitTime+processTime)/float64(sampleCount*millisecond)
					if averageResponseTime >= manager.OverloadControlLatencyThreshold*2.0 {
						averageWaitTime := float64(waitTime)/float64(sampleCount*millisecond)
						pendingLength := len(manager.Statistics)
						pendingRequestNumber := manager.GetPendingRequestNumber()
						requestStatistics := manager.GetRequestStatistics(count)
						log.Warningf("Overload, Overall Pass Pate %v Overload Control Sample Window = %v[%v], Average Wait Time = %.2f ms , Average Response Latency = %.2f ms, Pending Latency Statistics = %v , Pending Request Number = %v", manager.GetPassRate(), sampleCount, requestStatistics, averageWaitTime, averageResponseTime, pendingLength, pendingRequestNumber)

						manager.IncreaseOverload(count, 1.0/8.0, sampleCount)

						sampleCount = 0
						waitTime = 0
						processTime = 0
						for _, engine := range manager.WorkEngines {
							count[engine.GroupPriority] = 0
						}
					}
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