package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	const (
		maxWorkers   = 4
		maxQueueSize = 20
		port         = ":8081"
	)

	jobQueue := make(chan Job, maxQueueSize)
	dispatcher := NewDispatcher(jobQueue, maxWorkers)
	dispatcher.Run()

	http.HandleFunc("/fibonacci", func(writer http.ResponseWriter, request *http.Request) {
		requestHandler(writer, request, jobQueue)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func requestHandler(writer http.ResponseWriter, request *http.Request, jobQueue chan Job) {
	if request.Method != "POST" {
		writer.Header().Set("Allow", "POST")
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	delay, err := time.ParseDuration(request.FormValue("delay"))
	if err != nil {
		http.Error(writer, "Invalid Delay", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(request.FormValue("value"))
	if err != nil {
		http.Error(writer, "Invalid Value", http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	if name == "" {
		http.Error(writer, "Invalid Name", http.StatusBadRequest)
		return
	}

	job := Job{Name: name, Delay: delay, Number: number}
	jobQueue <- job
	writer.WriteHeader(http.StatusCreated)
}
