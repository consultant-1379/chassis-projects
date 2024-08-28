package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/utils"
)

var (
	TerminateWaitingTime time.Duration = 5
)

func (s *Server) handleSignals() {
	if cm.Opts.NoSigs {
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
				s.processSignalUSR1()
			case syscall.SIGHUP:
				log.Infof("Server signal sighup")
			case syscall.SIGTERM:
				log.Infof("Server signal terminate")
				s.processSignalSIGTERM()
				s.Terminate <- 1
			}
		}
	}()
}

func (s *Server) processSignalUSR1() {
	utils.GenCpuMemPprof(time.Duration(cm.Opts.PprofTime)*time.Second, "")
}

func (s *Server) processSignalSIGTERM() {
	log.Warningf("processSignalSIGTERM... sleep %v seconds", TerminateWaitingTime)
	time.Sleep(TerminateWaitingTime * time.Second)
	log.Warningf("ready for terminate the process")
}
