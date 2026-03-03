package main

import (
	"log/slog"
	"sync"
)

func connectDB() {
	slog.Info("db connected")
}

func main() {
	var wg sync.WaitGroup

	connectOnce := sync.OnceFunc(connectDB)
	for range 5 {
		wg.Go(func() {
			defer wg.Done()
			connectOnce()
		})
	}

	wg.Wait()
}
