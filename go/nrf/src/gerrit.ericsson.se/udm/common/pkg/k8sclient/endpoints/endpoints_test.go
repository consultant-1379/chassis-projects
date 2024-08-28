package endpoints

import (
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

type testEndpointsUpdateHandler struct {
	ipaddresses     []string
	monitorInterval int
}

func (t *testEndpointsUpdateHandler) HandleIPAddressUpdate(event Event, ipaddresses []string) {
	log.Debugf("endpoints update comes, event: %s, newIPAddresses: %v", GetEventName(event), ipaddresses)
	t.ipaddresses = ipaddresses
}

func (t *testEndpointsUpdateHandler) GetIPAddresses() []string {
	return t.ipaddresses
}

func (t *testEndpointsUpdateHandler) GetInterval() int {
	return t.monitorInterval
}

func init() {
	initLog()
	time.Sleep(time.Duration(3) * time.Second)
}

func initLog() {
	log.SetLevel(log.FatalLevel)
}

func TestMonitorEndpoints(t *testing.T) {

	Convey("TestApplyMethod", t, func() {
		Convey("endpoints with one pod", func() {
			NewWatcher().StartWatcher()

			_ = ApplyFunc(getEndpoints, func(string) []string {
				return []string{"192.168.240.30"}
			})

			handler := &testEndpointsUpdateHandler{
				monitorInterval: 3,
			}

			NewWatcher().AddEndpointsToWatcher("default", "eric-nrf-notification", handler)

			time.Sleep(time.Duration(4) * time.Second)
			currentIPAddresses := handler.GetIPAddresses()

			So(len(currentIPAddresses), ShouldEqual, 1)
			So(currentIPAddresses[0], ShouldEqual, "192.168.240.30")

			NewWatcher().quitOne <- true

			time.Sleep(time.Duration(1) * time.Second)

		})

		Convey("endpoints with two pod", func() {

			patch := ApplyFunc(getEndpoints, func(string) []string {
				return []string{"192.168.240.30", "192.168.240.31"}
			})

			handler := &testEndpointsUpdateHandler{
				monitorInterval: 3,
			}

			NewWatcher().AddEndpointsToWatcher("default", "eric-nrf-notification", handler)

			time.Sleep(time.Duration(4) * time.Second)
			currentIPAddresses := handler.GetIPAddresses()

			So(len(currentIPAddresses), ShouldEqual, 2)

			currentIPAddressesMap := make(map[string]bool)
			for _, ip := range currentIPAddresses {
				currentIPAddressesMap[ip] = true
			}

			So(currentIPAddressesMap["192.168.240.30"], ShouldEqual, true)
			So(currentIPAddressesMap["192.168.240.31"], ShouldEqual, true)

			NewWatcher().StopWatcher()

			time.Sleep(time.Duration(1) * time.Second)

			patch.Reset()
		})
	})
}

func TestGetEventName(t *testing.T) {
	if GetEventName(NOCHANGE) != "NOCHANGE" {
		t.Fatalf("GetEventName didn't return right value")
	}

	if GetEventName(ADD) != "ADD" {
		t.Fatalf("GetEventName didn't return right value")
	}

	if GetEventName(UPDATE) != "UPDATE" {
		t.Fatalf("GetEventName didn't return right value")
	}

	if GetEventName(REMOVE) != "REMOVE" {
		t.Fatalf("GetEventName didn't return right value")
	}
}

func TestCompare(t *testing.T) {
	watcher := &Watcher{}
	// case 1 : ip addresses change from [] to []
	if NOCHANGE != watcher.compare(nil, nil) {
		t.Fatalf("Watcher.compare didn't return right value")
	}

	// case 2 : ip addresses change from [] to [192.168.240.30]
	if ADD != watcher.compare(nil, []string{"192.168.240.30"}) {
		t.Fatalf("Watcher.compare didn't return right value")
	}

	// case 3 : ip addresses change from [192.168.240.30] to [192.168.240.30, 192.168.240.31]
	if UPDATE != watcher.compare([]string{"192.168.240.30"}, []string{"192.168.240.30", "192.168.240.31"}) {
		t.Fatalf("Watcher.compare didn't return right value")
	}

	// case 4 : ip addresses change from [192.168.240.30] to [192.168.240.31]
	if UPDATE != watcher.compare([]string{"192.168.240.30"}, []string{"192.168.240.31"}) {
		t.Fatalf("Watcher.compare didn't return right value")
	}

	// case 5 : ip addresses change from [192.168.240.30] to [192.168.240.30]
	if NOCHANGE != watcher.compare([]string{"192.168.240.30"}, []string{"192.168.240.30"}) {
		t.Fatalf("Watcher.compare didn't return right value")
	}

	// case 6 : ip addresses change from [192.168.240.30] to []
	if REMOVE != watcher.compare([]string{"192.168.240.30"}, nil) {
		t.Fatalf("Watcher.compare didn't return right value")
	}
}
