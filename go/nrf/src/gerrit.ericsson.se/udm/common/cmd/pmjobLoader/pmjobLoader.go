package main

import (
	"os"

	"gerrit.ericsson.se/udm/common/app/pmjobLoader"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

func main() {
	initLog()
	log.Debug("pmjobLoader starts...")

	if len(os.Args) <= 1 {
		log.Error("pmjobLoader miss parameter")
		os.Exit(1)
	}

	if os.Args[1] != pmjobLoader.LoadPm && os.Args[1] != pmjobLoader.UnLoadPm {
		log.Error("pmjobLoader miss parameter")
		os.Exit(1)
	}

	pmjobLoader.Init()
	if os.Args[1] == pmjobLoader.LoadPm {
		pmjobLoader.LoadPmJob()
	}

	if os.Args[1] == pmjobLoader.UnLoadPm {
		pmjobLoader.UnLoadPmJob()
	}

	return
}

func initLog() {
	log.SetLevel(log.Level(5))
	log.SetOutput(os.Stdout)
	log.SetServiceID("pmJobLoader")
	log.SetNF(os.Getenv("NF_TYPE"))
	log.SetPodIP(os.Getenv("POD_IP"))
	log.SetFormatter(&log.JSONFormatter{})
}
