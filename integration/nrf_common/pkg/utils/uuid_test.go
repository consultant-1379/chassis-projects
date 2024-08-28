package utils

import (
	"testing"
)

func TestUUUID(t *testing.T) {
	if _, err := GetUUIDString(); err != nil {
		t.Fatalf("Can not get uuid %s", err.Error())
	}
}
