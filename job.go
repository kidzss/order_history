package main

type Job struct {
	C     *MongoConf
	Query *Result
}

type JobQueue chan Job
