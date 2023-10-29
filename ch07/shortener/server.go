package shortener

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/servernoj/gobook/ch07/short"
)

type Server struct {
	http.Handler
	service *Service
}

func (s *Server) handlerHealth(w http.ResponseWriter, r *http.Request) *AppError {
	w.Write([]byte("OK"))
	return nil
}

func (s *Server) handlerError(w http.ResponseWriter, r *http.Request) *AppError {
	return &AppError{
		fmt.Errorf("test error"),
		http.StatusInternalServerError,
	}
}

func (s *Server) handlerCreate(w http.ResponseWriter, r *http.Request) *AppError {
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
	if err := s.service.LinkStore.Create(r.Context(), link); err != nil {
		return &AppError{
			err,
			http.StatusInternalServerError,
		}
	}
	EncodeAndSend(w, http.StatusCreated, link)
	return nil
}

func (s *Server) handlerResolve(w http.ResponseWriter, r *http.Request) *AppError {
	key := r.URL.Path[3:]
	link, err := s.service.LinkStore.Retrieve(r.Context(), key)
	if err != nil {
		return &AppError{
			fmt.Errorf("unable to retrieve link record with key %q", key),
			http.StatusInternalServerError,
		}
	}
	if link == nil {
		return &AppError{
			fmt.Errorf("key %q not found", key),
			http.StatusNotFound,
		}
	}
	http.Redirect(w, r, link.URL, http.StatusFound)
	return nil
}

func (s *Server) Init(service *Service) {
	s.service = service
	serveMux := http.NewServeMux()
	serveMux.Handle("/health", HandlerWrapper(s.handlerHealth))
	serveMux.Handle("/error", HandlerWrapper(s.handlerError))
	serveMux.Handle(
		"/short",
		MiddleWareAllowMethod(
			HandlerWrapper(s.handlerCreate),
			http.MethodPost,
		),
	)
	serveMux.Handle(
		"/r/",
		MiddleWareAllowMethod(
			HandlerWrapper(s.handlerResolve),
			http.MethodGet,
		),
	)
	s.Handler = serveMux
}

type HandlerToError func(w http.ResponseWriter, r *http.Request) *AppError

type contextKey string

const ServiceKey = contextKey("Service")

func HandlerWrapper(h HandlerToError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			Log(r.Context(), "%s %s %s %s (%d)\n", r.Method, r.URL.Path, r.RemoteAddr, err.Error, err.Code)
			if err.Code == http.StatusInternalServerError {
				err.Error = errors.New("internal server error")
			}
			EncodeAndSend(w, err.Code, map[string]string{
				"message": err.Error.Error(),
			})
		}
	}
}
