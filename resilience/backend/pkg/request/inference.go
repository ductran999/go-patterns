package request

import (
	"context"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/ductran999/letobserv/pkg/httpclient"
)

func DoInference(c httpclient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	url := "http://localhost:8080/v1/chat/completions"
	payload := []byte(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "Hello High Performance Go!"}]
	}`)

	resp, err := c.Post(ctx, url, payload, httpclient.ReqHeader{
		"Content-Type":  "application/json",
		"Authorization": "Bearer ABC",
	})
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
}
