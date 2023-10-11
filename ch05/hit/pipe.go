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
}

func (c *Client) Do() *Stat {
	stat := &Stat{}
	from := time.Now()

	p := Produce(
		c.NumberOfRequests,
		func() *http.Request {
			return c.RequestTemplate.Clone(context.TODO())
		},
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
		Send,
	)
	for result := range r {
		stat.Process(result)
	}

	stat.PostProcess(time.Since(from))
	return stat
}

func Send(r *http.Request) (result *Result) {
	client := http.Client{}
	start := time.Now()
	response, err := client.Do(r)
	result = &Result{
		Duration: time.Since(start),
	}
	if err != nil {
		result.Err = fmt.Errorf("unable to make http request: %w", err)
		return
	}
	if response.StatusCode >= http.StatusBadRequest {
		result.Err = fmt.Errorf("http error (%d)", response.StatusCode)
		return
	}
	body := response.Body
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		result.Err = fmt.Errorf("unable to parse response data: %w", err)
		return
	}
	result.Bytes = len(buf)
	return
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

func Produce(n int, fn func() *http.Request) chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		for i := 0; i < n; i++ {
			r := fn()
			out <- r
		}
		close(out)
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
