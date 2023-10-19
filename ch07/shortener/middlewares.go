package shortener

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Printf("%s %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		},
	)
}

type HandlerToError func(w http.ResponseWriter, r *http.Request) error

func HandlerWrapper(h HandlerToError) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h(w, r)
			switch {
			case err == nil:
			case errors.Is(err, ErrBadRequest):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, ErrNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, ErrAlreadyExists):
				http.Error(w, err.Error(), http.StatusConflict)
			case errors.Is(err, ErrInvalidMethod):
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("%s %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr, err.Error())
			}
		},
	)
}
