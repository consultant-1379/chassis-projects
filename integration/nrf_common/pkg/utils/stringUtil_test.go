package utils

import (
	"testing"
)

func TestIsDigit(t *testing.T) {
	s := "123456"
	if !IsDigit(s) {
		t.Fatalf("should be digit, failed")
	}

	s = "123wer"
	if IsDigit(s) {
		t.Fatalf("should not be digit, failed")
	}

}
