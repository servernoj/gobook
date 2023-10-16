package hit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientDo(t *testing.T) {

	// test server
	requestCounter := atomic.Int64{}
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				requestCounter.Add(1)
			},
		),
	)
	defer server.Close()
	// request template
	requestTemplate, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		t.Fatalf("unable to create request template")
	}
	// client instantiation
	client := &Client{
		RequestTemplate:  requestTemplate,
		Concurrency:      2,
		NumberOfRequests: 9,
	}
	// client Do the job
	stat := client.Do(context.Background())
	t.Logf("\n%s\n", stat)
	// assess the results
	if got, wanted := stat.Count, client.NumberOfRequests; got != wanted {
		t.Errorf("numbers of planned/sent requests don't match")
	}
	if got, wanted := requestCounter.Load(), int64(client.NumberOfRequests); got != wanted {
		t.Errorf("numbers of sent/received requests don't match")
	}
	if got, wanted := stat.Errors, 0; got != wanted {
		t.Errorf("no http errors were expected")
	}
}
