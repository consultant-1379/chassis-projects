package log

import (
	"os"
	"testing"
)

func TestJsonFormatter(t *testing.T) {
	SetServiceID("nrf_mgmt")
	SetNF("nrf")
	SetPodIP("10.10.10.10")
	SetFormatter(&JSONFormatter{})
	SetOutput(os.Stdout)
	SetLevel(DebugLevel)
	Debugf("debug log")
	Infof("info log")
	Warnf("warn log")
	Errorf("error log")
	//Fatal("fatal log")
}
