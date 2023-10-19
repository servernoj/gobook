package shortener

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/servernoj/gobook/ch07/short"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("OK"))
	return nil
}

func handlerError(w http.ResponseWriter, r *http.Request) error {
	return ErrTest
}

func handlerCreate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return ErrInvalidMethod
	}
	var link short.Link
	err := Decode(r.Body, &link)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBadRequest, err)
	}
	link.Key = strings.TrimSpace(link.Key)
	link.URL = strings.TrimSpace(link.URL)
	if len(link.Key) == 0 {
		return fmt.Errorf("%w: Key can't be an empty string", ErrBadRequest)
	}
	if len(link.URL) == 0 {
		return fmt.Errorf("%w: URL can't be an empty string", ErrBadRequest)
	}
	u, err := url.Parse(link.URL)
	if err != nil {
		return fmt.Errorf("%w: unable to parse URL %q", ErrBadRequest, link.URL)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: invalid URL scheme %q", ErrBadRequest, u.Scheme)
	}
	// -- store data
	short.Set(&link)
	// -- echo requestv data

	w.Header().Add("content-type", "application/json")
	Encode(w, &link)
	return nil
}

func handlerResolve(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return ErrInvalidMethod
	}
	key := r.URL.Path[3:]
	link := short.Get(key)
	if link == nil {
		return ErrNotFound
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
