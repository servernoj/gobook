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
	// client Do the job
	numberOfRequests := 100
	stat, err := Do(
		context.Background(), server.URL, numberOfRequests,
		WithConcurrency(1),
		WithRPS(10),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s\n", stat)
	// assess the results
	if got, wanted := stat.Count, numberOfRequests; got != wanted {
		t.Errorf("numbers of planned/sent requests don't match")
	}
	if got, wanted := requestCounter.Load(), int64(numberOfRequests); got != wanted {
		t.Errorf("numbers of sent/received requests don't match")
	}
	if got, wanted := stat.Errors, 0; got != wanted {
		t.Errorf("no http errors were expected")
	}
}
