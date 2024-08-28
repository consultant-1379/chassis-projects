package utils

import (
	"testing"
)

func TestTlsOne(t *testing.T) {
	if _, err := GenTlsConfig(false, "./server.crt", "./server.key", nil); err != nil {
		t.Fatalf("Can not get tlsconfig %s", err.Error())
	}
}
