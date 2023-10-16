package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type SenderFunc func(r *http.Request) *Result

type Client struct {
	RequestTemplate  *http.Request
	RPS              int
	Concurrency      int
	NumberOfRequests int
	httpClient       *http.Client
}

func (c *Client) getSender() SenderFunc {
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: c.Concurrency,
			},
		}
	}
	return func(r *http.Request) (result *Result) {
		start := time.Now()
		response, err := c.httpClient.Do(r)
		result = &Result{
			Duration: time.Since(start),
		}
		if err != nil {
			result.Err = fmt.Errorf("unable to make http request: %w", err)
			return nil
		}
		if response.StatusCode >= http.StatusBadRequest {
			result.Err = fmt.Errorf("http error (%d)", response.StatusCode)
			return nil
		}
		body := response.Body
		defer body.Close()
		buf, err := io.ReadAll(body)
		if err != nil {
			result.Err = fmt.Errorf("unable to parse response data: %w", err)
			return nil
		}
		result.Bytes = len(buf)
		return
	}
}

func (c *Client) Do(ctx context.Context) *Stat {
	stat := &Stat{}
	sender := c.getSender()
	defer c.httpClient.CloseIdleConnections()
	from := time.Now()
	p := Produce(
		ctx,
		c.NumberOfRequests,
		c.RequestTemplate,
	)
	if c.RPS > 0 {
		p = Throttle(
			p,
			time.Second/time.Duration(c.RPS*c.Concurrency),
		)
	}
	r := SplitAndSend(
		p,
		c.Concurrency,
		sender,
	)
	for result := range r {
		stat.Process(result)
	}

	stat.PostProcess(time.Since(from))
	return stat
}

func FakeSend(r *http.Request) *Result {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	result := &Result{
		Duration: time.Since(start),
		Bytes:    100,
		Err:      nil,
	}
	return result
}

func Produce(ctx context.Context, n int, requestTemplate *http.Request) chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- requestTemplate.Clone(ctx):
			}
		}
	}()

	return out
}

func Throttle(in chan *http.Request, delay time.Duration) chan *http.Request {
	out := make(chan *http.Request)
	ticker := time.NewTicker(delay)

	go func() {
		for r := range in {
			<-ticker.C
			out <- r
		}
		ticker.Stop()
		close(out)
	}()

	return out
}

func SplitAndSend(in chan *http.Request, concurrency int, send SenderFunc) chan *Result {
	out := make(chan *Result)

	go func() {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				for r := range in {
					out <- send(r)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(out)
	}()

	return out
}
