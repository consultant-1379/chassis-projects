package main

import (
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/options"
	"gerrit.ericsson.se/udm/nrf_discovery/app"
)

func main() {

	opts := options.Instance()

	server := nrf.NewServer(opts)
	defer server.Stop()

	log.Warningf("Version: %s", options.Version)

	go func() {
		cm.Setup()
	}()

	time.Sleep(5 * time.Second)

	log.Debugf("start the server")

	server.Run()

	select {
	case <-server.Terminate:
		break
	}
}
