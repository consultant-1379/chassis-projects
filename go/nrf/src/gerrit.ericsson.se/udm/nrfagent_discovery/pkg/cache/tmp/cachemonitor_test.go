package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	fmt.Println("test Init function")
}

func TestSuperviseDefaultTTL(t *testing.T) {
	nfInstanceID := "12345678-abcd-ef12-1000-000000000003"
	ttlMonitor := new(ttlMonitor)
	ttlMonitor.init("AUSF", "UDM", homeCache)
	ttlMonitor.superviseDefaultTTL(nfInstanceID)
	time.Sleep(1 * time.Second)

	live, ok := ttlMonitor.leftLive(nfInstanceID)
	if !ok {
		t.Errorf("TestLeftLive: cachemonitor LeftLive return failure.")
	}
	if live != ttlMonitor.defaultTtl-1 {
		t.Errorf("TestLeftLive: cachemonitor LeftLive live time value is more than left time.")
	}
}

func TestSupervise(t *testing.T) {
	fmt.Println("test Supervise function")
}

func TestStop(t *testing.T) {
	fmt.Println("test Supervise function")
}

func TestStopAll(t *testing.T) {
	fmt.Println("test StopAll function")
	nfInstanceID1 := "12345678-abcd-ef12-1000-000000000003"
	nfInstanceID2 := "12345678-abcd-ef12-2000-000000000006"
	var ttlTime uint = 30

	ttlMonitor := new(ttlMonitor)
	ttlMonitor.init("AUSF", "UDM", homeCache)
	ttlMonitor.supervise(nfInstanceID1, ttlTime)
	ttlMonitor.supervise(nfInstanceID2, ttlTime)
	ttlMonitor.stopAll()
	ttlMonitor.deleteAll()

	_, err := ttlMonitor.getTimePoint(nfInstanceID1)
	if err == nil {
		t.Fatalf("After delete all minitor, expect no timePoint for [%s], but not", nfInstanceID1)
	}

	_, err = ttlMonitor.getTimePoint(nfInstanceID2)
	if err == nil {
		t.Fatalf("After delete all minitor, expect no timePoint for [%s], but not", nfInstanceID2)
	}
}

func TestLeftLive(t *testing.T) {
	fmt.Println("test Reset function")
}

func TestDelete(t *testing.T) {
	nfInstanceID := "12345678-abcd-ef12-1000-000000000003"
	var ttlTime uint = 30

	ttlMonitor := new(ttlMonitor)
	ttlMonitor.init("AUSF", "UDM", homeCache)
	ttlMonitor.supervise(nfInstanceID, ttlTime)
	ttlMonitor.delete(nfInstanceID)

	_, err := ttlMonitor.getTimePoint(nfInstanceID)
	if err == nil {
		t.Fatalf("After delete all minitor, expect no timePoint for [%s], but not", nfInstanceID)
	}
}

func TestReset(t *testing.T) {
	nfInstanceID := "12345678-abcd-ef12-1000-000000000003"
	var ttlTime1 uint = 60
	ttlMonitor := new(ttlMonitor)
	ttlMonitor.init("AUSF", "UDM", homeCache)
	ttlMonitor.supervise(nfInstanceID, ttlTime1)
	time.Sleep(1 * time.Second)

	live1, ok1 := ttlMonitor.leftLive(nfInstanceID)
	if !ok1 {
		t.Errorf("TestReset: cachemonitor LeftLive return failure.")
	}
	if live1 != ttlTime1-1 {
		t.Errorf("TestReset: cachemonitor LeftLive live time value is more than left time.")
	}

	var ttlTime2 uint = 30
	ttlMonitor.reset(nfInstanceID, ttlTime2)
	live2, ok2 := ttlMonitor.leftLive(nfInstanceID)
	if !ok2 {
		t.Errorf("TestReset: cachemonitor LeftLive return failure.")
	}
	if live2 != ttlTime2 {
		t.Errorf("TestReset: cachemonitor LeftLive live time value is more than left time.")
	}
}

func TestStart(t *testing.T) {
	fmt.Println("test Start function")
}
