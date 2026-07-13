package gravitino

import "log"

type Client interface {
	WithTarget(baseURL, token string) TargetAPI
}

type clientImpl struct{}

func NewClient() Client {
	return &clientImpl{}
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

func (c *clientImpl) doRequest(path string) {
	log.Println("[HTTP CALL] Executing request to: ", path)
}
