package main

import (
	"log"
	"time"
)

// GenerateNumbers returns a channel that sends numbers
func GenerateNumbers() <-chan int {
	ch := make(chan int)

	// Spawn a go routine to generate numbers and send them to the channel
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		// Close the channel after sending all numbers
		close(ch)
	}()

	return ch
}

func main() {
	for num := range GenerateNumbers() {
		log.Println("received number:", num)
	}
}
