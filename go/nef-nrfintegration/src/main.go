package main

import (
	"gerrit.ericsson.se/nef/nef-golangcommon/pkg/log"
	"time"
)

var config = new(Config)
var logger = log.GetLogger("gerrit.ericsson.se/nef/nef-nrfintegration")

func main() {
	if err := config.LoadConfigMap(); err != nil {
		logger.Critical(err.Error())
	}
	client, err := NewKubeClient()
	if err != nil {
		logger.Critical(err.Error())
	}
	monitorAndReportServices := func() {
		for clusterSvcName, serviceName := range config.serviceNameMap {
			if serviceName =="dummy-service-fornrfagent" {
				if err := sendHeartbeat(serviceName); err != nil {
					logger.Error(err.Error())
				}
			}else {
				go MonitorAndReportSrvWithRetry(config.prefix+"-"+clusterSvcName, serviceName, client)
			}
		}
	}
	go runProbe("80")
	monitorAndReportServices()
	for range time.NewTicker(time.Duration(config.heartbeatInterval) * time.Second).C {
		monitorAndReportServices()
	}
}
