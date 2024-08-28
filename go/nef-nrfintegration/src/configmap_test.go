package main

import (
	"os"
	"testing"
)

func TestRequireNotEmpty(t *testing.T) {
	varName := "NAMESPACE"
	_ = os.Setenv(varName, "default")
	val := requireNotEmpty(varName)
	if val != "default" {
		t.Error("Test: requireNotEmpty failed")
	}
}

func TestRequireNotEmptyPanic(t *testing.T) {
	varName := "NAMESPACE"
	_ = os.Unsetenv(varName)
	defer assertPanic(t)
	_ = requireNotEmpty(varName)
}

func assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Error()
	}
}

func TestConfig_LoadConfigMap(t *testing.T) {
	setAllEnv := func() {
		_ = os.Setenv("NAMESPACE", "default")
		_ = os.Setenv("GLOBAL_PREFIX", "eric")
		_ = os.Setenv("NRF_AGENT_URI", "http://nrf-register-agent")
		_ = os.Setenv("NF_CM_PROFILE_NAME", "ericsson-nef")
		_ = os.Setenv("NF_PROFILE_PATH", "ericsson-nef:nef-function")
		_ = os.Setenv("HEARTBEAT_INTERVAL", "10")
		_ = os.Setenv("RETRY_INTERVAL", "1")
		_ = os.Setenv("RETRY_TIMES", "1")
		_ = os.Setenv("SERVICE_NAME_MAP", `{"nef-svc-sim":"nnef-ee","nef-svc-sim2":"nnef-ee","nef-svc-sim3":"nnef-ee"}`)
	}
	tests := []struct {
		name    string
		test    func()
		wantErr bool
	}{
		{
			"all",
			func() {},
			false,
		},
		{
			"InvalidHeartbeatInterval",
			func() { _ = os.Setenv("HEARTBEAT_INTERVAL", "10.2") },
			true,
		},
		{
			"InvalidRetryInterval",
			func() { _ = os.Setenv("RETRY_INTERVAL", "10.2") },
			true,
		},
		{
			"InvalidRetryInterval2",
			func() { _ = os.Setenv("RETRY_INTERVAL", "-1") },
			true,
		},
		{
			"InvalidRetry_Times",
			func() { _ = os.Setenv("RETRY_TIMES", "10.2") },
			true,
		},
		{
			"EmptyMap",
			func() { _ = os.Setenv("SERVICE_NAME_MAP", "{}") },
			false,
		},
		{
			"InvalidJson",
			func() { _ = os.Setenv("SERVICE_NAME_MAP", "{\"aaa\"}") },
			true,
		},
		{
			"InvalidCombinationOfIntervals",
			func() {
				_ = os.Setenv("HEARTBEAT_INTERVAL", "10")
				_ = os.Setenv("RETRY_INTERVAL", "9")
				_ = os.Setenv("RETRY_TIMES", "2")
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAllEnv()
			tt.test()
			if err := config.LoadConfigMap(); (err != nil) != tt.wantErr {
				t.Errorf("Config.LoadConfigMap() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil {
				t.Log(err)
			}
		})
	}
}
