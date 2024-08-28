package log

import (
	"testing"
)

func TestSetLevel(t *testing.T) {
	SetLevel(DebugLevel)
	Debugf("debug level")
	Infof("info level")
	Warnf("warn level")
	Errorf("error level")

	SetLevel(InfoLevel)
	Debugf("debug level")
	Infof("info level")
	Warnf("warn level")
	Errorf("error level")

	SetLevel(WarnLevel)
	Debugf("debug level")
	Infof("info level")
	Warnf("warn level")
	Errorf("error level")

	SetLevel(ErrorLevel)
	Debugf("debug level")
	Infof("info level")
	Warnf("warn level")
	Errorf("error level")
}

func TestGetLevel(t *testing.T) {
	SetLevel(DebugLevel)
	if GetLevel() != DebugLevel {
		t.Fatalf("log level should be DebugLevel, but not!")
	}

	SetLevel(InfoLevel)
	if GetLevel() != InfoLevel {
		t.Fatalf("log level should be InfoLevel, but not!")
	}

	SetLevel(WarnLevel)
	if GetLevel() != WarnLevel {
		t.Fatalf("log level should be WarnLevel, but not!")
	}

	SetLevel(ErrorLevel)
	if GetLevel() != ErrorLevel {
		t.Fatalf("log level should be ErrorLevel, but not!")
	}
}
