package internalconf

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

func TestParseDiscInternalConf(t *testing.T) {
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
        "connectionNum" : 30,
        "grpcContextTimeout" : 31
    },
    "regionNrf" : {
        "forward":{
            "retryTime": 40,
            "retryWaitTime": 41
              },
        "redirect":{
            "retryTime": 50,
            "retryWaitTime": 51
              }
    },
    "plmnNrf": {
        "forward":{
            "retryTime": 60,
            "retryWaitTime": 61
              },
        "statusCode": [500, 604, 704]
    },
    "discCache": {
        "discCacheEnable": true,
        "discCacheTimeThreshold": 70
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

	internalConf := &InternalDiscConf{}

	err := json.Unmarshal(jsonContent, internalConf)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	internalConf.ParseConf()

	if HomeNrfForwardRetryTime != 10 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if HomeNrfForwardRetryWaitTime != 11 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if httpserver.GetIdleTimeout().Minutes() != 20 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if httpserver.GetActiveTimeout().Minutes() != 21 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if !HTTPWithXVersion {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if DbproxyConnectionNum != 30 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if DbproxyGrpcCtxTimeout != 31 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if RegionNrfForwardRetryTime != 40 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if RegionNrfForwardRetryWaitTime != 41 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if RegionNrfRedirectRetryTime != 50 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if RegionNrfRedirectRetryWaitTime != 51 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if PlmnNrfForwardRetryTime != 60 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if PlmnNrfForwardRetryWaitTime != 61 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}



	if !DiscCacheEnable {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if DiscCacheTimeThreshold != 70 {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}

	if (&OverloadProtection) == nil {
		t.Fatalf("InternalDiscConf.ParseConf parse error")
	}
}
