package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		log.Println("---> Pool empty! Init new Buffer...")
		return new(bytes.Buffer)
	},
}

func main() {
	var wg sync.WaitGroup

	// Note:
	// Use pool when gc overhead (profiling to detect it)
	// Always reset before put or after get
	for i := range 10_000 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			buf, ok := bufferPool.Get().(*bytes.Buffer)
			if !ok {
				log.Panic("invalid type")
				return
			}
			defer func() {
				buf.Reset()
				bufferPool.Put(buf)
			}()

			fmt.Fprintf(buf, "Worker %d data", workerID)

			log.Printf("Worker %d done: %s\n", workerID, buf.String())
		}(i)
	}

	wg.Wait()
}
