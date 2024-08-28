package utils

import (
	"testing"
)

func TestRandomInt(t *testing.T) {
	if RandomInt(5, 20) < 5 || RandomInt(5, 20) > 20 {
		t.Fatal("RandomInt Error")
	}
}
