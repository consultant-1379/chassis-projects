package internalconf

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

func TestParseNotifInternalConf(t *testing.T) {
	jsonContent := []byte(`{
  "httpServer": {
    "idleTimeout": 20,
    "activeTimeout": 21,
    "httpWithXVersion": true
  },
  "dbProxy": {
    "connectionNum": 30,
    "grpcContextTimeout": 31
  },
  "NFStatusNotify": {
    "enable": false,
    "callBackURISupportIPv6": false,
    "notificationAlwaysWithFullNfProfile": true,
    "notificationMaxJobQueue": 2500,
    "notificationMaxJobWorker": 200,
    "notificationMaxNotifyWorker": 50,
    "notificationTimeout": 15,
    "enableDebugLogForOverload": true
  }
}`)

	internalConf := &InternalNotifyConf{}

	err := json.Unmarshal(jsonContent, internalConf)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	internalConf.ParseConf()

	if httpserver.GetIdleTimeout().Minutes() != 20 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if httpserver.GetActiveTimeout().Minutes() != 21 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if !HTTPWithXVersion {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if DbproxyConnectionNum != 30 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if DbproxyGrpcCtxTimeout != 31 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if EnableNotification {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if CallBackURISupportIPv6 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if !NotificationAlwaysWithFullNfProfile {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if NotificationMaxJobQueue != 2500 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if NotificationMaxJobWorker != 200 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if NotificationMaxNotifyWorker != 50 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if NotificationTimeout != 15 {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}

	if !EnableDebugLogForOverload {
		t.Fatalf("InternalNotifyConf.ParseConf parse error")
	}
}
