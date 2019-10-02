package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/pkg/flash"

	"github.com/gorilla/sessions"
)

var (
	SessionCtxKey = contextKey("Session")
)

// MustBeAuthenticated enforces that a user sending a request is logged in.
// It checks this by seeing if the session has a non empty email. If the
// session does not have an email then that means the session is not valid
// and so the request is redirected to the login page.
func MustBeAuthenticated(h http.Handler) http.Handler {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	store := createCookieStore()

	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionCookieName)

		email, ok := session.Values["EMAIL"].(string)
		if !ok || email == "" {
			session.AddFlash(flash.Flash{
				Type:    flash.Info,
				Message: models.ErrNotLoggedIn.Error(),
			})

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// MustBeExec only allows execs and admin roles to access a route
func MustBeExec(h http.Handler) http.Handler {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	store := createCookieStore()

	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionCookieName)

		role, ok := session.Values["ROLE"].(int)
		if !ok {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		if role != models.RoleExec && role != models.RoleAdmin {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// mustGetEnv looks up and sets an environment variable.  If the environment
// variable is not found, it panics.
func mustGetEnv(varName string) (value string) {
	value, ok := os.LookupEnv(varName)
	if !ok {
		log.Fatalf("environment variable not set: %v", varName)
	}
	return value
}

func createCookieStore() *sessions.CookieStore {
	sessionSecret := mustGetEnv("SESSION_SECRET")
	mode := mustGetEnv("ENV_MODE")

	key := []byte(sessionSecret)
	store := sessions.NewCookieStore(key)

	// Prevents CSRF attacks (on browsers that support SameSite)
	store.Options.SameSite = http.SameSiteStrictMode

	// Prevents XSS attacks (JS isn't allowed to access cookie)
	store.Options.HttpOnly = true

	if mode != "development" {
		// Only transfer cookie over https
		store.Options.Secure = true
	}

	return store
}
