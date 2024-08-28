package cm

import (
	"strings"
	"testing"
)

func TestSetDefaultForManagementService(t *testing.T) {

	SetDefaultForManagementService()

	if GetMgmtService().Heartbeat.Default != defaultHeartbeatTimer {
		t.Fatal("Heartbeat.Default should be 120, but not !")
	}

	if GetMgmtService().Heartbeat.GracePeriod != defaultHeartbeatTimerGracePeriod {
		t.Fatal("Heartbeat.GracePeriod should be 5, but not !")
	}

	if GetMgmtService().TrafficRateLimitPerInstance != defaultTrafficRateLimitPerInstance {
		t.Fatal("TrafficRateLimitPerInstance should be 10, but not !")
	}

	if ValidityPeriodOfSubscription != defaultValidityPeriodOfSubscription {
		t.Fatalf("ValidityPeriodOfSubscription should be 604800, but not !")
	}
}

func TestExpectedHeartBeat(t *testing.T) {
	heartbeat := &THeartbeat{Default: 120, GracePeriod: 5}
	heartbeat.DefaultPerNfType = append(heartbeat.DefaultPerNfType, TDefaultPerNfType{NfType: "udm", Value: 10})
	managementProfileIns := &TManagementService{
		SubscriptionExpiredTime: 604800,
		Heartbeat:               heartbeat,
	}

	managementProfileIns.ParseConf()

	if managementProfileIns.Heartbeat.Default != 120 {
		t.Fatal("Heartbeat.Default should be 120, but not !")
	}

	if managementProfileIns.Heartbeat.GracePeriod != 5 {
		t.Fatal("Heartbeat.GracePeriod should be 5, but not !")
	}

	if ValidityPeriodOfSubscription != 604800 {
		t.Fatalf("ValidityPeriodOfSubscription should be 604800, but not !")
	}
}

func TestUnExpectedSubExpiredTime(t *testing.T) {
	managementProfileIns := &TManagementService{
		SubscriptionExpiredTime: -6048,
	}

	managementProfileIns.ParseConf()

	if ValidityPeriodOfSubscription != defaultValidityPeriodOfSubscription {
		t.Fatalf("ValidityPeriodOfSubscription should be 604800, but not !")
	}
}

func TestGetDefaultHeartbeatTimer(t *testing.T) {
	heartbeat := &THeartbeat{Default: 120, GracePeriod: 5}
	heartbeat.DefaultPerNfType = append(heartbeat.DefaultPerNfType, TDefaultPerNfType{NfType: "udm", Value: 10})
	managementProfileIns := &TManagementService{
		SubscriptionExpiredTime: 604800,
		Heartbeat:               heartbeat,
	}

	managementProfileIns.ParseConf()

	if heartbeat.GetDefaultHeartbeatTimer("AUSF") != 120 {
		t.Fatalf("GetDefaultHeartbeatTimer of AUSF should be 120, but not !")
	}

	if heartbeat.GetDefaultHeartbeatTimer("UDM") != 10 {
		t.Fatalf("GetDefaultHeartbeatTimer of UDM should be 10, but not !")
	}
}

func TestTManagementServicetoUpper(t *testing.T) {
	managementService := TManagementService{
		Heartbeat: &THeartbeat{
			Default: 120,
			DefaultPerNfType: []TDefaultPerNfType{
				TDefaultPerNfType{
					NfType: "amf",
					Value:  120,
				},
				TDefaultPerNfType{
					NfType: "ausf",
					Value:  120,
				},
			},
		},
	}

	originalManagementService := managementService

	managementService.toUpper()

	for index := range managementService.Heartbeat.DefaultPerNfType {
		if managementService.Heartbeat.DefaultPerNfType[index].NfType != strings.ToUpper(originalManagementService.Heartbeat.DefaultPerNfType[index].NfType) {
			t.Fatalf("TManagementService.toUpper didn't return right value !")
		}
	}
}
