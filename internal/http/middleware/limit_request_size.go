package middleware

import (
	"net/http"
)

// LimitRequestSize prevents requests from being over 32MiB
func LimitRequestSize(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
