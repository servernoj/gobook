package hit

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

type Option func(client *Client)

func WithRPS(RPS int) Option {
	return func(client *Client) {
		client.RPS = RPS
	}
}
func WithConcurrency(concurrency int) Option {
	return func(client *Client) {
		client.Concurrency = concurrency
	}
}
func WithClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func Do(ctx context.Context, url string, numberOfRequests int, opts ...Option) (*Stat, error) {
	requestTemplate, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request template: %w", err)
	}
	client := &Client{
		RequestTemplate:  requestTemplate,
		NumberOfRequests: numberOfRequests,
		Concurrency:      runtime.NumCPU(),
		RPS:              0,
	}
	for _, o := range opts {
		o(client)
	}
	return client.Do(ctx), nil
}
