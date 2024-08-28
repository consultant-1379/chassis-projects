package signal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/probe"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	terminate = make(chan int)
)

// HoldUntilSIGTERM returns once process receives SIGTERM
func HoldUntilSIGTERM() {
	select {
	case <-terminate:
		break
	}
}

// HandleSignals catches signals and handles them
func HandleSignals(NoSigs bool, PprofTime int) {
	if NoSigs {
		return
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		for sig := range c {
			log.Debugf("Trapped %q signal", sig)
			switch sig {
			case syscall.SIGINT:
				log.Warningf("Server signal Exiting..")
				os.Exit(0)
			case syscall.SIGUSR1:
				processSignalUSR1(PprofTime)
			case syscall.SIGHUP:
				log.Infof("Server signal sighup")
			case syscall.SIGTERM:
				log.Infof("Server signal terminate")
				processSignalSIGTERM()
				terminate <- 1
			}
		}
	}()

}

func processSignalUSR1(PprofTime int) {
	utils.GenCpuMemPprof(time.Duration(PprofTime)*time.Second, "")
}

func processSignalSIGTERM() {
	log.Infof("processSignalSIGTERM... ")
	probe.SetShutDownFlag(true)
	time.Sleep(constvalue.TerminateWaitingTime * time.Second)
}
