package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

var (
	jwtCookie string // Name for the JWT's cookie.  TODO Better name?
)

// A Handler will route requests to their appropriate HandlerFunc.
type Handler struct {
	// A Mux handles all routing and middleware.
	*chi.Mux

	// A UserService is the interface with the database.
	UserService models.UserService
}

// NewHandler creates a handler for API requests.
func NewHandler(us models.UserService) *Handler {
	h := Handler{UserService: us}
	r := chi.NewRouter()

	// TODO See cmd/server/main.go for more about config. This doesn't seem ideal.
	var ok bool
	jwtCookie, ok = os.LookupEnv("JWT_COOKIE_NAME")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_COOKIE_NAME")
	}

	r.Use(middleware.SetContentTypeJSON) // All responses from here will be JSON
	r.Use(middleware.WithJWT)

	r.Post("/signup", h.postSignup())
	r.Post("/login", h.postLogin())

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.getSelf())
	})

	h.Mux = r
	return &h
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// postSignup tries to sign up a user.
func (h *Handler) postSignup() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO check if login is valid (i.e. account exists), if so log them in

		//TODO replace with env variable
		session, sess_err := store.Get(r, "session-cookie-name")
		if sess_err != nil {
			// TODO error handling
			http.Error(w, sess_err.Error(), http.StatusInternalServerError)
			return
		}

		var u models.User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&u)

		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := h.UserService.Signup(&u)
		// if err != nil {
		// 	// TODO error handling
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		//
		// 	return
		// }
		u.ID = id
		// No other errors
		session.Values["ID"] = u.ID
		println(session.Values["ID"])
		err = session.Save(r, w) // TODO error checking

	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO replace cookie name with an environment variable
		session, sess_err := store.Get(r, "session-cookie-name")

		if sess_err != nil {
			// TODO error handling
			http.Error(w, sess_err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert JSON user details in request to a user struct
		var u models.User

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&u)

		if err != nil {

			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.UserService.Login(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// No other errors
		session.Values["ID"] = u.ID
		err = session.Save(r, w) // TODO error checking
	}
}

// getSelf returns the user as JSON.  It will redact some fields (like
// password).
// TODO Right now it sets the password value to blank but keeps the now blank
// field in the JSON response.  Consider even removing that field.
// TODO Same as above, but for passwordConfirm
func (h *Handler) getSelf() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO replace cookie name with env name
		session, err := store.Get(r, "session-cookie-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		spew.Dump(session)
		json.NewEncoder(w).Encode(session)
		return
	}
}
