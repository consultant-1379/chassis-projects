package utils

import (
	"fmt"
	"testing"
)

func TestSequence(t *testing.T) {
	SetPrefix("host1")
	for i := 1; i < 10; i++ {
		sequenceId := GetSequenceId()
		id := 1000 + i
		compareSequenceId := fmt.Sprintf("host1_%d", id)
		if sequenceId != compareSequenceId {
			t.Fatalf("Get wrong sequenceid %s", sequenceId)
		}
	}

}
