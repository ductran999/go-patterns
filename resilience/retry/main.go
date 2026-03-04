package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"patterns/resilience/backend/pkg/request"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ductran999/letobserv/pkg/httpclient"
)

func main() {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.Multiplier = 2.0
	b.MaxInterval = 10 * time.Second
	b.MaxElapsedTime = 1 * time.Minute

	c := httpclient.New()

	infer := func() error {
		resp, err := request.DoInference(c)

		// 1. Handle Network Errors (Timeout, DNS, Connection Refused)
		if err != nil {
			slog.Error("Network error occurred", "error", err.Error())
			return err // Triggers Retry
		}
		defer resp.Body.Close()

		// 2. Handle Server Errors (500, 502, 503, 504)
		// These are usually temporary. We SHOULD retry.
		if resp.StatusCode >= 500 {
			errMessage := fmt.Errorf("server side error: %d", resp.StatusCode)
			slog.Warn("Retrying due to server error", "status", resp.StatusCode)
			return errMessage // Triggers Retry
		}

		// 3. Handle Client Errors (400, 401, 403, 404)
		// These mean our request is wrong. Retrying won't help.
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			errMessage := fmt.Errorf("critical client error: %d", resp.StatusCode)
			slog.Error("Permanent error, stopping retry", "status", resp.StatusCode)

			// Use backoff.Permanent to STOP all further retry attempts immediately.
			return backoff.Permanent(errMessage)
		}

		// 4. Success (Status 200-299)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err // Retry if body reading fails
		}

		log.Println("Success:", string(bodyBytes))
		return nil // Success, stop retrying.
	}

	// Start the retry process
	if err := backoff.Retry(infer, b); err != nil {
		slog.Error("Operation failed after all retries", "final_error", err)
	}
}
