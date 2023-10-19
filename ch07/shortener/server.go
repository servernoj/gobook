package shortener

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/servernoj/gobook/ch07/short"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) *AppError {
	w.Write([]byte("OK"))
	return nil
}

func handlerError(w http.ResponseWriter, r *http.Request) *AppError {
	return &AppError{
		fmt.Errorf("test error"),
		http.StatusInternalServerError,
	}
}

func handlerCreate(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method != "POST" {
		return &AppError{
			fmt.Errorf("method %q not allowed", r.Method),
			http.StatusMethodNotAllowed,
		}
	}
	var link short.Link
	err := Decode(r.Body, &link)
	if err != nil {
		return &AppError{
			fmt.Errorf("unable to decode request body: %w", err),
			http.StatusBadRequest,
		}
	}
	link.Key = strings.TrimSpace(link.Key)
	link.URL = strings.TrimSpace(link.URL)
	if len(link.Key) == 0 {
		return &AppError{
			fmt.Errorf("request body field 'key' can't be an empty string"),
			http.StatusBadRequest,
		}
	}
	if len(link.URL) == 0 {
		return &AppError{
			fmt.Errorf("request body field 'URL' can't be an empty string"),
			http.StatusBadRequest,
		}
	}
	u, err := url.Parse(link.URL)
	if err != nil {
		return &AppError{
			fmt.Errorf("unable to parse URL field %q: %w", link.URL, err),
			http.StatusBadRequest,
		}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return &AppError{
			fmt.Errorf("invalid URL scheme %q, must be 'http[s]'", u.Scheme),
			http.StatusBadRequest,
		}
	}
	// -- store data
	short.Set(&link)
	// -- echo requestv data

	w.Header().Add("content-type", "application/json")
	Encode(w, &link)
	return nil
}

func handlerResolve(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Method != "GET" {
		return &AppError{
			fmt.Errorf("method %q not allowed", r.Method),
			http.StatusMethodNotAllowed,
		}
	}
	key := r.URL.Path[3:]
	link := short.Get(key)
	if link == nil {
		return &AppError{
			fmt.Errorf("key %q not found", key),
			http.StatusNotFound,
		}
	}
	http.Redirect(w, r, link.URL, http.StatusFound)
	return nil
}

type Server struct {
	http.Handler
	ctx context.Context
}

func (s *Server) Init() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/health", HandlerWrapper(handlerHealth))
	serveMux.Handle("/error", HandlerWrapper(handlerError))
	serveMux.Handle("/short", HandlerWrapper(handlerCreate))
	serveMux.Handle("/r/", HandlerWrapper(handlerResolve))
	s.Handler = serveMux
	s.ctx = context.Background()
}
