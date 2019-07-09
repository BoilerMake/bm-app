package middleware

import (
	"net/http"
)

// SetContentTypeJSON sets the Content-Type for a response to application/json
func SetContentTypeJSON(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "application/json")
	}

	return http.HandlerFunc(fn)
}
