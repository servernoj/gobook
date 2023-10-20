package shortener

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

func Log(ctx context.Context, fmt string, args ...any) {
	if s, ok := ctx.Value(http.ServerContextKey).(*http.Server); ok && s != nil && s.ErrorLog != nil {
		s.ErrorLog.Printf(fmt, args...)
	}
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			Log(r.Context(), "%s %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
		},
	)
}

type HandlerToError func(w http.ResponseWriter, r *http.Request) *AppError

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

func MiddleWareAllowMethod(h http.Handler, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(allowedMethod, r.Method) {
			Log(r.Context(), "%s %s %s method not allowed", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	}
}
