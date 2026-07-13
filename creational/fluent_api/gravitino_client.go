package gravitino

import "fmt"

type Client interface {
	WithTarget(baseURL, token string) TargetAPI
}

type clientImpl struct{}

func NewClient() Client {
	return &clientImpl{}
}

func (c *clientImpl) doRequest(path string) {
	fmt.Printf("[HTTP CALL] Executing request to: %s\n", path)
}

func (c *clientImpl) WithTarget(baseURL, token string) TargetAPI {
	return &targetImpl{
		client:  c,
		baseURL: baseURL,
		token:   token,
	}
}

type targetImpl struct {
	client  *clientImpl
	baseURL string
	token   string
}
