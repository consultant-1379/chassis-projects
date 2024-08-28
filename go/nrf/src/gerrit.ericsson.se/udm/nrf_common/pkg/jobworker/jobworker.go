package jobworker

import (
	"fmt"
)

type JobWorker struct {
	maxWorkers int
	workerPool chan chan Job
	jobQueue   chan Job
}

type Job interface {
	Handler()
}

type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func NewJobWorker(maxWorkers, maxQueue int) *JobWorker {
	j := &JobWorker{
		workerPool: make(chan chan Job, maxWorkers),
		maxWorkers: maxWorkers,
		jobQueue:   make(chan Job, maxQueue),
	}
	j.run()

	return j
}

func (j *JobWorker) AddJob(job Job) {
	j.jobQueue <- job
}

func (j *JobWorker) GetJobQueueLength() int {
	return len(j.jobQueue)
}

//add job into queue for disc
func (j *JobWorker) AddJobForDisc(job Job) bool {
	select {
	case j.jobQueue <- job:
		return true
	default:
		return false
	}
}

func (d *JobWorker) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *JobWorker) dispatch() {
	for {
		select {
		case t := <-d.jobQueue:
			jobChannel := <-d.workerPool
			jobChannel <- t
		}
	}
}

func newWorker(workerPool chan chan Job) worker {
	return worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w worker) stop() {
	go func() {
		w.quit <- true
	}()
}

func (w worker) start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel

			select {
			case j := <-w.jobChannel:
				j.Handler()

			case <-w.quit:
				fmt.Println("Worker is exiting")
				return
			}
		}
	}()
}
