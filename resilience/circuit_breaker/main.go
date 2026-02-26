package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "PythonBackend"
	st.MaxRequests = 1
	st.Interval = 0
	st.Timeout = 30 * time.Second

	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func callPythonBackend() ([]byte, error) {
	body, err := cb.Execute(func() (any, error) {
		resp, err := http.Get("http://localhost:8000/v1/chat/completions")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	})

	if err != nil {
		return nil, err
	}
	return body.([]byte), nil
}

func main() {
	for i := range 20 {
		log.Println(i, "---> call api")
		_, err := callPythonBackend()
		if err != nil {
			slog.Error(err.Error())
		}
		time.Sleep(time.Second * time.Duration(i))
	}
}
