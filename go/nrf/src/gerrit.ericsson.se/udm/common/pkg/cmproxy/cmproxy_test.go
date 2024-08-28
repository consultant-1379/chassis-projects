package cmproxy

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
)

var (
	postStatus    = 0
	putSubStatus  = 0
	getConfStatus = 0
)

func cmMediatorPostHandler(rw http.ResponseWriter, req *http.Request) {
	switch postStatus {
	default:
		rw.WriteHeader(http.StatusCreated)
	case 409:
		rw.WriteHeader(http.StatusConflict)
	}
}
func cmMediatorPutSHandler(rw http.ResponseWriter, req *http.Request) {
	switch putSubStatus {
	case 404:
		rw.WriteHeader(http.StatusNotFound)
	default:
		rw.WriteHeader(http.StatusOK)
	}
}
func cmMediatorGetCHandler(rw http.ResponseWriter, req *http.Request) {
	switch getConfStatus {
	case 404:
		rw.WriteHeader(http.StatusNotFound)
	default:
		rw.WriteHeader(http.StatusOK)
	}
}

func cmNfProfileHandler(Event, ConfigurationName, format string, RawData []byte) {
	fmt.Println("cmNfProfileHandler")
}

func cmNrfServiceProfilesHandler(Event, ConfigurationName, format string, RawData []byte) {
	fmt.Println("cmNrfServiceProfilesHandler")
}

func TestRegisterAndDeRegisterConf(t *testing.T) {
	msgbus.PendingInitializeFailed = false

	configName1 := "nfProfile"
	configName2 := "nrfServiceProfile"
	jsonPath := "test"
	topic := "test"

	h := httpserver.InitHTTPServer(
		httpserver.Trace(true),
		httpserver.HostPort("", "5003"),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.SetRoute(),
		httpserver.PathFunc("/cm/api/v1.1/subscriptions", "POST", cmMediatorPostHandler),
		httpserver.PathFunc("/cm/api/v1.1/subscriptions"+configName1+"_sub", "PUT", cmMediatorPutSHandler),
		httpserver.PathFunc("/cm/api/v1.1/subscriptions"+configName2+"_sub", "PUT", cmMediatorPutSHandler),
		httpserver.PathFunc("/cm/api/v1.1/configurations", "POST", cmMediatorPostHandler),
		httpserver.PathFunc("/cm/api/v1.1/configurations"+configName1, "GET", cmMediatorGetCHandler),
		httpserver.PathFunc("/cm/api/v1.1/configurations"+configName2, "GET", cmMediatorGetCHandler),
	)
	h.Run()
	defer func() {
		h.Stop()
	}()

	Init("http://127.0.0.1:5003/cm/api/v1.1/")
	Run()
	defer func() {
		Stop()
	}()

	t.Run("TestRegisterConf1", func(t *testing.T) {
		RegisterConf(configName1, jsonPath, topic, cmNfProfileHandler, NtfFormatFull)
		_, existed := cmConfigList[configName1]
		if !existed {
			t.Fatalf("cmConfig %s is not exist, RegisterConf configName1 Failed", configName1)
			delete(cmConfigList, configName1)
		}
	})
	t.Run("TestRegisterConf2", func(t *testing.T) {
		RegisterConf(configName2, jsonPath, topic, cmNfProfileHandler, NtfFormatPatch)
		_, existed := cmConfigList[configName2]
		if !existed {
			t.Fatalf("cmConfig %s is not exist, RegisterConf configName1 Failed", configName2)
			delete(cmConfigList, configName2)
		}
	})
	t.Run("TestRegisterConf1Twice", func(t *testing.T) {
		RegisterConf(configName1, jsonPath, topic, cmNrfServiceProfilesHandler, NtfFormatFull)
	})

	t.Run("TestGetConfigurations", func(t *testing.T) {
		getConfStatus = 0
		getConfigurations()
	})
	t.Run("TestPutSubscriptions", func(t *testing.T) {
		putSubStatus = 0
		putSubscriptions()
	})
	t.Run("TestGetConfigurations", func(t *testing.T) {
		getConfStatus = 404
		defer func() {
			getConfStatus = 0
		}()
		getConfigurations()
	})
	t.Run("TestPutSubscriptions", func(t *testing.T) {
		putSubStatus = 404
		defer func() {
			putSubStatus = 0
		}()
		putSubscriptions()
	})

	t.Run("TestDeRegisterConf1", func(t *testing.T) {
		DeRegisterConf(configName1, topic)
		_, existed := cmConfigList[configName1]
		if existed {
			t.Fatalf("DeRegisterConf configName1 Failed")
		}
	})
	t.Run("TestDeRegisterConf2", func(t *testing.T) {
		DeRegisterConf(configName2, topic)
		_, existed := cmConfigList[configName1]
		if existed {
			t.Fatalf("DeRegisterConf configName2 Failed")
		}
	})
	t.Run("TestDeRegisterConf1Twice", func(t *testing.T) {
		DeRegisterConf(configName1, topic)
	})

	t.Run("TestGetConfigurations", func(t *testing.T) {
		getConfigurations()
	})
	t.Run("TestPutSubscriptions", func(t *testing.T) {
		putSubscriptions()
	})

}
