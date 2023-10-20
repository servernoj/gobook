package shortener

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerHandlers(t *testing.T) {

	type inData struct {
		handler http.Handler
		body    string
		path    string
		method  string
	}

	type outData struct {
		status  int
		body    string
		headers http.Header
	}

	tests := map[string]struct {
		in  inData
		out outData
	}{
		"health endpoint": {
			in: inData{
				handler: HandlerWrapper(handlerHealth),
				body:    "",
				path:    "/health",
				method:  http.MethodGet,
			},
			out: outData{
				status: http.StatusOK,
				body:   "OK",
				headers: http.Header{
					"content-type": []string{"text/plain"},
				},
			},
		},
		"error endpoint": {
			in: inData{
				handler: HandlerWrapper(handlerError),
				body:    "",
				path:    "/error",
				method:  http.MethodGet,
			},
			out: outData{
				status: http.StatusInternalServerError,
				body:   "internal server error",
				headers: http.Header{
					"content-type": []string{"application/json"},
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(
				tt.in.method,
				tt.in.path,
				strings.NewReader(tt.in.body),
			)
			recorder := httptest.NewRecorder()
			tt.in.handler.ServeHTTP(recorder, request)
			if got, want := recorder.Code, tt.out.status; got != want {
				t.Errorf("response status mismatch, got: %d, want: %d", got, want)
			}
			if got, want := recorder.Body.String(), tt.out.body; len(want) > 0 && !strings.Contains(got, want) {
				t.Errorf("response body mismatch, got: %q, want: %q", got, want)
			}
			if tt.out.headers != nil {
				for key, value := range tt.out.headers {
					recorderHeaderValue := recorder.Header().Get(key)
					found := false
					for _, valueItem := range value {
						if strings.Contains(recorderHeaderValue, valueItem) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("header (%q,%q) wanted but was not found in %q", key, value, recorderHeaderValue)
					}
				}
			}
		})

	}

}
