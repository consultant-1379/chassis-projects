package jobworker

import (
	"testing"
	"time"
)

type task struct {
	b bool
}

func (t *task) Handler() {
	t.b = true
}

func TestDispatcher(t *testing.T) {
	a := &task{
		b: false,
	}
	j := NewJobWorker(3, 10)
	j.AddJob(a)
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	if a.b == false {
		t.Errorf("Dispatcher can not work")
	}
}
