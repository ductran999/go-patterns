package main

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"patterns/resilience/backend/pkg/request"
	"strconv"
	"sync"
	"time"

	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "PythonBackend"
	st.MaxRequests = 1

	st.Interval = 60 * time.Second // reset counter after 60s
	st.Timeout = 5 * time.Second

	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func callPythonBackend(c httpclient.Client) ([]byte, error) {
	body, err := cb.Execute(func() (any, error) {
		resp, err := request.DoInferenceError(c)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("backend error: status " + strconv.Itoa(resp.StatusCode))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	})

	if err != nil {
		// Check ErrOpenState -> Circuit is closed
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, errors.New("circuit open: backend python is down")
		}
		return nil, err
	}

	return body.([]byte), nil
}

func main() {
	c := httpclient.New()
	var wg sync.WaitGroup

	for range 20 {
		wg.Go(func() {
			resp, err := callPythonBackend(c)
			if err != nil {
				slog.Error(err.Error())
			} else {
				log.Println("resp:", string(resp))
			}
		})

		time.Sleep(time.Second)
	}

	wg.Wait()
}
