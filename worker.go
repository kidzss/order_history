package main

import (
//"log"
//"net/http"
//"io/ioutil"
)

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				mongo := &MongoClient{MongoConf: *(job.C)}
				mongo.Execute(job.Query)
			case <-w.Quit:
				return
			}
		}
	}()
}
func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
