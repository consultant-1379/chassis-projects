package main

import (
	"gerrit.ericsson.se/nef/nef-golangcommon/pkg/log"
	"testing"
)

func TestNewLogger(t *testing.T) {
	if logger.GetLevel() != log.InfoLevel {
		t.Error(logger.GetLevel())
	}
}
