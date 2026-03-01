package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

func produceJobs(data []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, v := range data {
			out <- v
		}
		close(out)
	}()

	return out
}

func worker(id int, jobs <-chan string) {
	for job := range jobs {
		log.Printf("Worker #%d processing: %s\n", id, job)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	tasks := []string{
		"homework",
		"write unit test",
		"write ci/cd pipeline",
		"implement user authentication",
		"refactor legacy database logic",
		"write unit tests for auth module",
		"integrate redis for caching",
		"optimize sql queries for reports",
		"fix memory leak in stream processor",
		"document api using swagger",
		"write github actions workflow",
		"configure jenkins pipeline",
		"build docker images for production",
		"optimize dockerfile layers",
		"setup sonarqube code analysis",
		"configure argocd sync policy",
		"write helm charts for microservices",
		"implement canary release strategy",
		"write terraform for vpc setup",
		"provision rds database on aws",
		"configure s3 bucket permissions",
		"setup kubernetes ingress controller",
		"audit iam roles and policies",
		"implement secrets management using vault",
		"configure prometheus alerts",
		"create grafana dashboard for latency",
		"setup elk stack for log aggregation",
		"configure heartbeat for uptime monitoring",
		"implement distributed tracing with jaeger",
		"run database migration v2.1",
		"backup production postgresql",
		"rotate cloud access keys",
		"cleanup unused docker volumes",
	}

	jobChannel := produceJobs(tasks)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU() * 2
	for i := range numWorkers {
		wg.Go(func() {
			worker(i, jobChannel)
		})
	}

	wg.Wait()
}
