package shortener

import (
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
			fmt.Errorf("request body field 'key' cannot be an empty string"),
			http.StatusBadRequest,
		}
	}
	if len(link.URL) == 0 {
		return &AppError{
			fmt.Errorf("request body field 'url' cannot be an empty string"),
			http.StatusBadRequest,
		}
	}
	u, err := url.Parse(link.URL)
	if err != nil {
		return &AppError{
			fmt.Errorf("unable to parse URL %q: %w", link.URL, err),
			http.StatusBadRequest,
		}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return &AppError{
			fmt.Errorf("invalid URL scheme %q, must be 'http[s]'", u.Scheme),
			http.StatusBadRequest,
		}
	}
	short.Set(&link)
	EncodeAndSend(w, http.StatusOK, &link)
	return nil
}

func handlerResolve(w http.ResponseWriter, r *http.Request) *AppError {
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
}

func (s *Server) Init() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/health", HandlerWrapper(handlerHealth))
	serveMux.Handle("/error", HandlerWrapper(handlerError))
	serveMux.Handle(
		"/short",
		MiddleWareAllowMethod(
			HandlerWrapper(handlerCreate),
			http.MethodPost,
		),
	)
	serveMux.Handle(
		"/r/",
		MiddleWareAllowMethod(
			HandlerWrapper(handlerResolve),
			http.MethodGet,
		),
	)
	s.Handler = serveMux
}
