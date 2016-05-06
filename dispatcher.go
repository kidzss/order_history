package main

import (
	"log"
)

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	Sli_worker []Worker
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: maxWorkers}
}

func (d *Dispatcher) Run(jobQueue JobQueue) {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		d.Sli_worker = append(d.Sli_worker, worker)
		//log.Println("create", i, "workers success")
		worker.Start()
		//log.Println("worker", i, "started")
	}
	log.Printf("create %d worker\n", d.MaxWorkers)
	go d.Dispatch(jobQueue)
}

func (d *Dispatcher) Dispatch(jobQueue JobQueue) {
	for {
		select {
		case job := <-jobQueue:
			//log.Println("get ", job, "Job from JobQueue")
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.Sli_worker {
		worker.Stop()
	}
}
