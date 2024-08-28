package utils

import (
	"os"
	"testing"
	"time"
)

const (
	cpufile = "/tmp/cpu.pprof"
	memfile = "/tmp/mem.pprof"
)

func TestGenCpuMemPprof(t *testing.T) {
	_ = os.Remove(cpufile)
	_ = os.Remove(memfile)

	GenCpuMemPprof(5*time.Second, "")
	if GenCpuMemPprof(5*time.Second, "") {
		t.Errorf("Can not generate more pprof")
	}
	time.Sleep(8 * time.Second)
	if _, err := os.Stat(cpufile); err != nil {
		t.Errorf(err.Error())
	}
	if _, err := os.Stat(memfile); err != nil {
		t.Errorf(err.Error())
	}
}
