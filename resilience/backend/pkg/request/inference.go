package request

import (
	"context"
	"net/http"
	"time"

	"github.com/ductran999/letobserv/pkg/httpclient"
)

func DoInference(c httpclient.Client) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	url := "http://localhost:8080/v1/chat/completions"
	payload := []byte(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "Hello High Performance Go!"}]
	}`)

	return c.Post(ctx, url, payload, httpclient.ReqHeader{
		"Content-Type":  "application/json",
		"Authorization": "Bearer ABC",
	})
}

func DoInferenceError(c httpclient.Client) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	url := "http://localhost:8080/v1/completions"
	payload := []byte(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "Hello High Performance Go!"}]
	}`)

	return c.Post(ctx, url, payload, httpclient.ReqHeader{
		"Content-Type":  "application/json",
		"Authorization": "Bearer ABC",
	})
}
