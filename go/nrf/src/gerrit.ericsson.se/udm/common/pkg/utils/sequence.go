package utils

import (
	"fmt"
	"sync"
)

var seq = New()

func SetPrefix(prefix string) {
	seq.SetPrefix(prefix)
}

func GetSequenceId() string {
	return seq.GetSequenceId()
}

type Sequence struct {
	Id     uint64
	Prefix string
	Mu     sync.Mutex
}

func New() *Sequence {
	return &Sequence{
		Id: 1000,
	}
}

func (seq *Sequence) SetPrefix(prefix string) {
	seq.Mu.Lock()
	defer seq.Mu.Unlock()
	seq.Prefix = prefix
}

func (seq *Sequence) GetSequenceId() string {
	seq.Mu.Lock()
	defer seq.Mu.Unlock()
	seq.Id++
	sequenceId := fmt.Sprintf("%s_%d", seq.Prefix, seq.Id)
	return sequenceId
}
