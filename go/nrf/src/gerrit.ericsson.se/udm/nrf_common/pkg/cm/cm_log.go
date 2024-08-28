package cm

import (
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

// ParseConf is to parse log info level and set log level for service or pod
func (conf *TNfServiceLog) ParseConf() {
	podID := PodIP
	usePodLevel := false

	for i := 0; i < len(conf.PodLogs); i++ {
		if podID == conf.PodLogs[i].PodID {
			if conf.PodLogs[i].Severity != "INHERIT" {
				log.SetLevel(log.LevelUint(conf.PodLogs[i].Severity))
				usePodLevel = true
				break
			} else {
				break
			}
		}
	}

	if !usePodLevel {
		log.SetLevel(log.LevelUint(conf.Severity))
	}
}

// Show is to print log level info
func (conf *TNfServiceLog) Show() {
	fmt.Printf("log level : %s\n", log.LevelToString(log.GetLevel()))
}
