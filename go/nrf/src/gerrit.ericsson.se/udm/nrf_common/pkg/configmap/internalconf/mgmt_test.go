package internalconf

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

func TestParseMgmtInternalConf(t *testing.T) {
	jsonContent := []byte(`{
  "homeNrf": {
    "forward": {
      "retryTime": 10,
      "retryWaitTime": 11
    }
  },
  "httpServer": {
    "idleTimeout": 20,
    "activeTimeout": 21,
    "httpWithXVersion": true
  },
  "dbProxy": {
    "connectionNum": 30,
    "grpcContextTimeout": 31
  },
  "regionNrf": {
    "nrfInfoCheckInterval":40
  },
  "NFStatusNotify": {
    "enable": false,
    "callBackURISupportIPv6": true,
    "notificationAlwaysWithFullNfProfile": false,
    "notificationMaxJobQueue": 2500,
    "notificationMaxJobWorker": 20,
    "notificationMaxNotifyWorker": 5,
    "notificationTimeout": 15,
    "enableDebugLogForOverload": false
  },
  "overloadControl": {
    "enableDebugLog": false,
    "trafficRateLimitPerNfInstance": 50,
    "retryAfterRange": {
      "start": 51,
      "end": 52
    }
  },
    "OverloadProtection": {
        "OverloadControlLevel": 32768,
        "OverloadTriggerLatencyThreshold": 10.0,
        "OverloadControlLatencyThreshold": 10.0,
        "OverloadTriggerSampleWindow": 20000,
        "OverloadControlSampleWindow": 20000,
        "IdleInterval": 1000,
        "IdleRecoverRatio": 32,
        "CounterReportInterval": 500,
        "OverloadAlarmClearWindow": 20,
        "WorkEngine": [
                {
                        "GroupPriority": 1,
                        "QueueCapacity": 4096,
                        "WorkerNumber": 256
                },
                {
                        "GroupPriority": 2,
                        "QueueCapacity": 4096,
                        "WorkerNumber": 256
                },
                {
                        "GroupPriority": 3,
                        "QueueCapacity": 4096,
                        "WorkerNumber": 256
                }
        ]
    }
}`)

	internalConf := &InternalMgmtConf{}

	err := json.Unmarshal(jsonContent, internalConf)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	internalConf.ParseConf()

	if HomeNrfForwardRetryTime != 10 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if HomeNrfForwardRetryWaitTime != 11 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if httpserver.GetIdleTimeout().Minutes() != 20 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if httpserver.GetActiveTimeout().Minutes() != 21 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if !HTTPWithXVersion {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if DbproxyConnectionNum != 30 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if DbproxyGrpcCtxTimeout != 31 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if NrfInfoCheckInterval != 40 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if EnableNotification {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if TrafficRateLimitPerNfInstance != 50 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if RetryAfterRangeStart != 51 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if RetryAfterRangeEnd != 52 {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}

	if (&OverloadProtection) == nil {
		t.Fatalf("InternalMgmtConf.ParseConf parse error")
	}
}
