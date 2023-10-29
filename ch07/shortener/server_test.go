package shortener

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/servernoj/gobook/ch07/short"
)

type fakeLinkStore struct {
	create   func(context.Context, short.Link) error
	retrieve func(context.Context, string) (*short.Link, error)
}

func (fls *fakeLinkStore) Create(ctx context.Context, link short.Link) error {
	if fls.create == nil {
		return nil
	}
	return fls.create(ctx, link)
}
func (fls *fakeLinkStore) Retrieve(ctx context.Context, key string) (*short.Link, error) {
	if fls.retrieve == nil {
		return nil, nil
	}
	return fls.retrieve(ctx, key)
}

func TestServerHandlers(t *testing.T) {

	type inData struct {
		service     *Service
		handlerName string
		body        string
		path        string
		method      string
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
				handlerName: "HandlerHealth",
				body:        "",
				path:        "/health",
				method:      http.MethodGet,
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
				handlerName: "HandlerError",
				body:        "",
				path:        "/error",
				method:      http.MethodGet,
			},
			out: outData{
				status: http.StatusInternalServerError,
				body:   "internal server error",
				headers: http.Header{
					"content-type": []string{"application/json"},
				},
			},
		},
		"retrieve not found": {
			in: inData{
				service: &Service{
					LinkStore: &fakeLinkStore{
						retrieve: func(context.Context, string) (*short.Link, error) {
							return nil, nil
						},
					},
				},
				handlerName: "HandlerResolve",
				body:        "",
				path:        "/r/boo",
				method:      http.MethodGet,
			},
			out: outData{
				status: http.StatusNotFound,
			},
		},
		"retrieve internal error": {
			in: inData{
				service: &Service{
					LinkStore: &fakeLinkStore{
						retrieve: func(context.Context, string) (*short.Link, error) {
							return nil, errors.New("boo")
						},
					},
				},
				handlerName: "HandlerResolve",
				body:        "",
				path:        "/r/boo",
				method:      http.MethodGet,
			},
			out: outData{
				status: http.StatusInternalServerError,
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
			server := &Server{
				service: tt.in.service,
			}
			handler := HandlerWrapper(
				func(w http.ResponseWriter, r *http.Request) *AppError {
					result := reflect.ValueOf(server).MethodByName(tt.in.handlerName).Call([]reflect.Value{
						reflect.ValueOf(w),
						reflect.ValueOf(r),
					})[0].Interface().(*AppError)
					return result
				},
			)
			handler.ServeHTTP(recorder, request)
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
