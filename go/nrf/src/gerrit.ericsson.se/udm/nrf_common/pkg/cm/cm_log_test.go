package cm

import (
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

func init() {
	log.SetLevel(log.FatalLevel)
}

func TestPodLevel(t *testing.T) {
	logCong := &TNfServiceLog{
		LogID:    "nrf_mgmt",
		Severity: "WARNING",
		PodLogs: []PodLog{
			PodLog{
				PodID:    "192.168.240.22",
				Severity: "ERROR",
			},
			PodLog{
				PodID:    "192.168.240.25",
				Severity: "DEBUG",
			},
			PodLog{
				PodID:    "192.168.240.26",
				Severity: "INFO",
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	PodIP = "192.168.240.22"

	logCong.ParseConf()

	if log.GetLevel() != log.ErrorLevel {
		t.Fatal("failed to parse logconf")
	}

	log.SetLevel(log.FatalLevel)
}

func TestServiceLevel(t *testing.T) {
	logCong := &TNfServiceLog{
		LogID:    "nrf_mgmt",
		Severity: "WARNING",
		PodLogs: []PodLog{
			PodLog{
				PodID:    "192.168.240.22",
				Severity: "ERROR",
			},
			PodLog{
				PodID:    "192.168.240.25",
				Severity: "DEBUG",
			},
			PodLog{
				PodID:    "192.168.240.26",
				Severity: "INFO",
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	PodIP = "192.168.240.01"

	logCong.ParseConf()

	if log.GetLevel() != log.WarnLevel {
		t.Fatal("failed to parse logconf")
	}

	log.SetLevel(log.FatalLevel)
}

func TestLevelIngerit(t *testing.T) {
	logCong := &TNfServiceLog{
		LogID:    "nrf_mgmt",
		Severity: "WARNING",
		PodLogs: []PodLog{
			PodLog{
				PodID:    "192.168.240.22",
				Severity: "ERROR",
			},
			PodLog{
				PodID:    "192.168.240.25",
				Severity: "INHERIT",
			},
			PodLog{
				PodID:    "192.168.240.26",
				Severity: "INFO",
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	PodIP = "192.168.240.25"

	logCong.ParseConf()

	if log.GetLevel() != log.WarnLevel {
		t.Fatal("failed to parse logconf")
	}

	log.SetLevel(log.FatalLevel)
}

func TestUnexpectedLevel(t *testing.T) {
	logCong := &TNfServiceLog{
		LogID:    "nrf_mgmt",
		Severity: "WARNING",
		PodLogs: []PodLog{
			PodLog{
				PodID:    "192.168.240.22",
				Severity: "ERROR",
			},
			PodLog{
				PodID:    "192.168.240.25",
				Severity: "INHERIT",
			},
			PodLog{
				PodID:    "192.168.240.26",
				Severity: "UNKNOWN",
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	PodIP = "192.168.240.26"

	logCong.ParseConf()

	if log.GetLevel() != log.WarnLevel {
		t.Fatal("failed to parse logconf")
	}

	log.SetLevel(log.FatalLevel)
}
