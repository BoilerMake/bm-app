package middleware

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func SessionMiddleware(h http.Handler) http.Handler {
	var (
		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		sessionSecret = mustGetEnv("SESSION_SECRET")
		key           = []byte(sessionSecret)
		store         = sessions.NewCookieStore(key)
	)
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")

	fn := func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, sessionCookieName)
		if err != nil {
			// TODO Error Handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), SessionCtxKey, session))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

var (
	SessionCtxKey = "session"
)

var (
	SessionSecret = "SESSION_SECRET"
)

var (
	SessionCookieName = "SESSION_COOKIE_NAME"
)

var (
	JWTCtxKey      = contextKey("JWT")
	JWTErrorCtxKey = contextKey("JWTError")
)

// mustGetEnv looks up and sets an environment variable.  If the environment
// variable is not found, it panics.
func mustGetEnv(var_name string) (value string) {
	value, ok := os.LookupEnv(var_name)
	if !ok {
		log.Fatalf("environment variable not set: %v", var_name)
	}
	return value
}
