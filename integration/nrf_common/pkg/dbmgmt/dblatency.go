package dbmgmt

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
)

//Latency is struct of time
type Latency struct {
	startGrpcTime    int64
	enterDBProxyTime uint64
	leaveDBProxyTime uint64
	endGrpcTime      int64
	FilterStartTime  int64
	FilterEndTime    int64
	ReqStartTime     int64
	ReqEndTime       int64
}

//DbLatency is struct of all kinds of request time statistics
type DbLatency struct {
	groupIDRequestChannel   chan Latency
	instIDRequestChannel    chan Latency
	nfProfileFilterChannel  chan Latency
	InnerFilterChannel      chan Latency
	HandlerChannel          chan Latency
	groupIDTotalCostTime    int64
	groupIDTransmitReqTime  int64
	groupIDTransmitRespTime int64
	groupIDDbProcessTime    int64
	instIDTotalCostTime     int64
	instIDTransmitReqTime   int64
	instIDTransmitRespTime  int64
	instIDDbProcessTime     int64
	filterTotalCostTime     int64
	filterTransmitReqTime   int64
	filterTransmitRespTime  int64
	filterDbProcessTime     int64
	innerFilterTime         int64
	handlerTime             int64
}

// DBLatency is an entrance of all counter
var DBLatency *DbLatency

//StartCalculate is to init channel and start to calculate avg time
func StartCalculate() {
	DBLatency = &DbLatency{}
	DBLatency.groupIDRequestChannel = make(chan Latency, 100)
	DBLatency.instIDRequestChannel = make(chan Latency, 100)
	DBLatency.nfProfileFilterChannel = make(chan Latency, 100)
	DBLatency.InnerFilterChannel = make(chan Latency, 100)
	DBLatency.HandlerChannel = make(chan Latency, 100)
	go DBLatency.start()
}

//start is to start a goroutine to process the event
func (c *DbLatency) start() {
	groupIDCount := 0
	instIDCount := 0
	filterCount := 0
	innerFilerCount := 0
	handlerCount := 0
	for {
		select {
		case latency := <-c.groupIDRequestChannel:
			groupIDCount++
			c.groupIDTransmitReqTime += int64(latency.enterDBProxyTime) - latency.startGrpcTime
			c.groupIDDbProcessTime += int64(latency.leaveDBProxyTime) - int64(latency.enterDBProxyTime)
			c.groupIDTransmitRespTime += latency.endGrpcTime - int64(latency.leaveDBProxyTime)
			c.groupIDTotalCostTime += latency.endGrpcTime - latency.startGrpcTime
			if groupIDCount >= internalconf.StatisticsNum {
				log.Warnf("GROUP ID total cost time=%vms, send request to grpc cost time=%vms, dbproxy process cost time=%vms, send response to discovery cost time=%vms, ", float64(c.groupIDTotalCostTime) / float64(groupIDCount), float64(c.groupIDTransmitReqTime) / float64(groupIDCount), float64(c.groupIDDbProcessTime) / float64(groupIDCount), float64(c.groupIDTransmitRespTime) / float64(groupIDCount))
				c.groupIDTransmitReqTime = 0
				c.groupIDDbProcessTime = 0
				c.groupIDTransmitRespTime = 0
				c.groupIDTotalCostTime = 0
				groupIDCount = 0
			}
		case latency := <-c.instIDRequestChannel:
			instIDCount++
			c.instIDTransmitReqTime += int64(latency.enterDBProxyTime) - latency.startGrpcTime
			c.instIDDbProcessTime += int64(latency.leaveDBProxyTime) - int64(latency.enterDBProxyTime)
			c.instIDTransmitRespTime += latency.endGrpcTime - int64(latency.leaveDBProxyTime)
			c.instIDTotalCostTime += latency.endGrpcTime - latency.startGrpcTime
			if instIDCount >= internalconf.StatisticsNum {
				log.Warnf("INSTANCE ID total cost time=%vms, send request to grpc cost time=%vms, dbproxy process cost time=%vms, send response to discovery cost time=%vms, ", float64(c.instIDTotalCostTime) / float64(instIDCount), float64(c.instIDTransmitReqTime) / float64(instIDCount), float64(c.instIDDbProcessTime) / float64(instIDCount), float64(c.instIDTransmitRespTime) / float64(instIDCount))
				c.instIDTransmitReqTime = 0
				c.instIDTransmitRespTime = 0
				c.instIDDbProcessTime = 0
				c.instIDTotalCostTime = 0
				instIDCount = 0
			}

		case latency := <-c.nfProfileFilterChannel:
			filterCount++
			c.filterTransmitReqTime += int64(latency.enterDBProxyTime) - latency.startGrpcTime
			c.filterDbProcessTime += int64(latency.leaveDBProxyTime) - int64(latency.enterDBProxyTime)
			c.filterTransmitRespTime += latency.endGrpcTime - int64(latency.leaveDBProxyTime)
			c.filterTotalCostTime += latency.endGrpcTime - latency.startGrpcTime
			if filterCount >= internalconf.StatisticsNum {
				log.Warnf("NFPROFILE total cost time=%vms, send request to grpc cost time=%vms, dbproxy process cost time=%vms, send response to discovery cost time=%vms, ", float64(c.filterTotalCostTime) / float64(filterCount), float64(c.filterTransmitReqTime) / float64(filterCount), float64(c.filterDbProcessTime) / float64(filterCount), float64(c.filterTransmitRespTime) / float64(filterCount))
				c.filterTransmitReqTime = 0
				c.filterTransmitRespTime = 0
				c.filterDbProcessTime = 0
				c.filterTotalCostTime = 0
				filterCount = 0
			}
		case latency := <-c.InnerFilterChannel:
			innerFilerCount++
			c.innerFilterTime += latency.FilterEndTime - latency.FilterStartTime
			if innerFilerCount >= internalconf.StatisticsNum {
				log.Warnf("INNER FILTER COST TIME=%vms, ", float64(c.innerFilterTime) / float64(innerFilerCount))
				c.innerFilterTime = 0
				innerFilerCount = 0
			}
		case latency := <-c.HandlerChannel:
			handlerCount++
			c.handlerTime += latency.ReqEndTime - latency.ReqStartTime
			if handlerCount >= internalconf.StatisticsNum {
				log.Warnf("HANDLER TOTAL COST TIME=%vms, ", float64(c.handlerTime) / float64(handlerCount))
				c.handlerTime = 0
				handlerCount = 0
			}

		}
	}
}
