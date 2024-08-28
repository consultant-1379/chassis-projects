package fm

import (
	"errors"
	"os"
	"strings"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/fmproxy"
)

func SendAlarmStub(e error) {
	sendAlarm = func(alarmPara *fmproxy.AlarmInfo, isRaise bool) error {
		return e
	}
}

func TestConnectionStatus(t *testing.T) {
	os.Setenv("POD_NAME", "127.0.0.1")

	t.Run("TestConnectionStatus01", func(t *testing.T) {
		SendAlarmStub(nil)
		ConnectionStatus("nrf-mgmt", false)
		if !alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus01 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-mgmt", true)
		if alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus01 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-disc", false)
		if !alarmRaised["nrf-disc"] {
			t.Errorf("TestConnectionStatus01 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-disc", true)
		if alarmRaised["nrf-disc"] {
			t.Errorf("TestConnectionStatus01 failed, %v", alarmRaised)
		}
	})
	t.Run("TestConnectionStatus03", func(t *testing.T) {
		var err = errors.New("test errors")
		SendAlarmStub(err)
		ConnectionStatus("nrf-mgmt", false)
		if alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus03 failed, %v", alarmRaised)
		}
	})
}

func TestConnectionStatus02(t *testing.T) {
	os.Setenv("POD_NAME", "127.0.0.1")

	t.Run("TestConnectionStatus02", func(t *testing.T) {
		SendAlarmStub(nil)
		ConnectionStatus("nrf-mgmt", false)
		if !alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-mgmt", false)
		if !alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-disc", false)
		if !alarmRaised["nrf-disc"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-disc", false)
		if !alarmRaised["nrf-disc"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-mgmt", true)
		if alarmRaised["nrf-mgmt"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
		ConnectionStatus("nrf-disc", true)
		if alarmRaised["nrf-disc"] {
			t.Errorf("TestConnectionStatus02 failed, %v", alarmRaised)
		}
	})
}

func TestDestinationStatus(t *testing.T) {
	os.Setenv("POD_NAME", "127.0.0.1")

	t.Run("TestDestinationStatus01", func(t *testing.T) {
		SendAlarmStub(nil)
		DestinationStatus("noAvailableDestination", false, "UDR", "AUSF")
		if !strings.Contains(unavailableNfTypeList["UDR"], "AUSF") {
			t.Errorf("TestDestinationStatus01 failed, %v", unavailableNfTypeList)
		}
		DestinationStatus("noAvailableDestination", true, "UDR", "AUSF")
		if strings.Contains(unavailableNfTypeList["UDR"], "AUSF") {
			t.Errorf("TestDestinationStatus01 failed, %v", unavailableNfTypeList)
		}
	})

	t.Run("TestDestinationStatus03", func(t *testing.T) {
		var err = errors.New("test errors")
		SendAlarmStub(err)
		DestinationStatus("noAvailableDestination", false, "UDR", "AUSF")
		if !strings.Contains(unavailableNfTypeList["UDR"], "AUSF") {
			t.Errorf("TestDestinationStatus03 failed, %v", unavailableNfTypeList)
		}
	})
}

func TestDestinationStatus02(t *testing.T) {
	os.Setenv("POD_NAME", "127.0.0.1")
	t.Run("TestDestinationStatus02", func(t *testing.T) {
		SendAlarmStub(nil)
		DestinationStatus("noAvailableDestination", false, "UDR", "AUSF")
		if !strings.Contains(unavailableNfTypeList["UDR"], "AUSF") {
			t.Errorf("TestDestinationStatus02 failed, %v", unavailableNfTypeList)
		}
		DestinationStatus("noAvailableDestination", false, "UDR", "UDM")
		if !strings.Contains(unavailableNfTypeList["UDR"], "AUSF") ||
			!strings.Contains(unavailableNfTypeList["UDR"], "UDM") {
			t.Errorf("TestDestinationStatus02 failed, %v", unavailableNfTypeList)
		}
		DestinationStatus("noAvailableDestination", true, "UDR", "AUSF")
		if strings.Contains(unavailableNfTypeList["UDR"], "AUSF") ||
			!strings.Contains(unavailableNfTypeList["UDR"], "UDM") {
			t.Errorf("TestDestinationStatus02 failed, %v", unavailableNfTypeList)
		}
		DestinationStatus("noAvailableDestination", false, "UDR", "UDM")
		if strings.Contains(unavailableNfTypeList["UDR"], "AUSF") ||
			!strings.Contains(unavailableNfTypeList["UDR"], "UDM") {
			t.Errorf("TestDestinationStatus02 failed, %v", unavailableNfTypeList)
		}
		DestinationStatus("noAvailableDestination", true, "UDR", "UDM")
		if strings.Contains(unavailableNfTypeList["UDR"], "AUSF") ||
			strings.Contains(unavailableNfTypeList["UDR"], "UDM") {
			t.Errorf("TestDestinationStatus02 failed, %v", unavailableNfTypeList)
		}
	})
}
