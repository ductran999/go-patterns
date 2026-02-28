package main

import (
	"io"
	"log"
	"log/slog"
	"patterns/resilience/backend/pkg/request"
	"sync"

	"github.com/ductran999/letobserv/pkg/httpclient"
)

func main() {
	c := httpclient.New()
	wg := new(sync.WaitGroup)

	for range 20 {
		wg.Go(func() {
			resp, err := request.DoInference(c)
			if err != nil {
				slog.Error("do inference request failed", "error_detail", err.Error())
				return
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("failed to read body", "error", err)
				return
			}

			log.Println(string(bodyBytes))
		})
	}

	wg.Wait()
}
