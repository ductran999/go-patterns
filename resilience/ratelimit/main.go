package main

import (
	"patterns/resilience/backend/pkg/request"
	"sync"

	"github.com/ductran999/letobserv/pkg/httpclient"
)

func main() {
	c := httpclient.New()
	wg := new(sync.WaitGroup)

	for range 20 {
		wg.Go(func() {
			request.DoInference(c)
		})
	}

	wg.Wait()
}
