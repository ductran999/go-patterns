package main

import (
	"bufio"
	_ "embed"
	"log"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed mock_data.csv
var embeddedCSVData string

type WorkerPool struct {
	numWorkers int
	jobs       chan string
	results    chan error
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan string, 1000),
		results:    make(chan error, 1000),
	}
}

func (wp *WorkerPool) StreamJobFromEmbedded(data string) error {
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		wp.jobs <- line
	}

	close(wp.jobs)

	return scanner.Err()
}

func (wp *WorkerPool) SpawnWorkers() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for job := range wp.jobs {
				if err := processRecord(id, job); err != nil {
					wp.results <- err
				}
			}
		}(i)
	}

	go func() {
		wp.wg.Wait()
		close(wp.results)
	}()
}

func (wp *WorkerPool) CollectResult() {
	for err := range wp.results {
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func processRecord(workerID int, record string) error {
	slog.Info("Worker read", "worker_id", workerID, "record", record)
	time.Sleep(1 * time.Millisecond)
	return nil
}

func main() {
	s := time.Now()

	numberOfWorker := runtime.NumCPU() * 2
	wp := NewWorkerPool(numberOfWorker)

	go func() {
		if err := wp.StreamJobFromEmbedded(embeddedCSVData); err != nil {
			log.Fatalln("failed to stream embedded data", err)
		}
	}()

	wp.SpawnWorkers()

	wp.CollectResult()

	log.Println("Processed time:", time.Since(s))
}
