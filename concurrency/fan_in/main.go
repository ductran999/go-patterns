package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// MergeChanel Fan-in: Merge multiple input channels into one output channel
func MergeChanel(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	merged := make(chan string)

	for _, ch := range channels {
		wg.Go(func() {
			for msg := range ch {
				merged <- msg
			}
		})
	}

	// Wait for all senders to finish, then close the merged channel
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// LogProducer simulate a service that sends logs into its own channel
func LogProducer(serviceName string, out chan<- string) {
	for i := 1; i <= 3; i++ {
		out <- fmt.Sprintf("[%s] log %d", serviceName, i)
		time.Sleep(time.Millisecond * 300)
	}
	close(out)
}

func main() {
	serviceA := make(chan string)
	serviceB := make(chan string)
	serviceC := make(chan string)

	// Start 3 log producers
	go LogProducer("Service A", serviceA)
	go LogProducer("Service B", serviceB)
	go LogProducer("Service C", serviceC)

	// Fan-in their output
	combined := MergeChanel(serviceA, serviceB, serviceC)

	for msg := range combined {
		log.Println("[INFO]", msg)
	}
}
