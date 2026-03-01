package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

func callBackend() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/auth", nil)
	if err != nil {
		return nil, err
	}

	// Add x-api-key header
	req.Header.Set("X-Api-Key", "ak-1234")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func main() {
	wg := sync.WaitGroup{}
	for i := range 20 {
		wg.Go(
			func() {
				resp, err := callBackend()
				if err != nil {
					slog.Error(err.Error())
				} else {
					slog.Info(string(resp), "goroutine", i)
				}
			},
		)

		// Burst 5 request to test singleflight request
		if i > 5 {
			time.Sleep(time.Second * 1)
		}
	}

	wg.Wait()
}
