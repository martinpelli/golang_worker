package main

import (
	"fmt"
	"time"
)

type Job struct {
	Name   string
	Delay  time.Duration
	Number int
}

type Worker struct {
	Id         int
	JobQueue   chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	worker := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: worker,
	}
}

func (dispatcher *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-dispatcher.JobQueue:
			go func() {
				workerJobQueue := <-dispatcher.WorkerPool
				workerJobQueue <- job
			}()
		}
	}
}

func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (worker Worker) Start() {
	go func() {
		for {
			worker.WorkerPool <- worker.JobQueue
			select {
			case job := <-worker.JobQueue:
				fmt.Printf("Worker with id %d started \n", worker.Id)
				fibonacci := Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d finished with result %d\n", worker.Id, fibonacci)
			case <-worker.QuitChan:
				fmt.Printf("Worker with id %d stopped \n", worker.Id)
			}
		}
	}()
}

func (worker Worker) Stop() {
	go func() {
		worker.QuitChan <- true
	}()
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + (n - 2)
}
