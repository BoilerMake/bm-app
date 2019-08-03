package api

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type User struct {
	Username      string
	Authenticated bool
}

// store will hold all session data
var store *sessions.CookieStore

// tpl holds all parsed templates
var tpl *template.Template

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})

	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// A Handler will route requests to their appropriate HandlerFunc.
type Handler struct {
	// A Mux handles all routing and middleware.
	*chi.Mux

	// A UserService is the interface with the database.
	UserService models.UserService
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

// NewHandler creates a handler for API requests.
func NewHandler(us models.UserService) *Handler {
	h := Handler{UserService: us}
	r := chi.NewRouter()

	// Need to configure this better

	// TODO See cmd/server/main.go for more about config. This doesn't seem ideal.
	// var ok bool
	// jwtCookie, ok = os.LookupEnv("JWT_COOKIE_NAME")
	// if !ok {
	// 	log.Fatalf("environment variable not set: %v", "JWT_COOKIE_NAME")
	// }

	r.Use(middleware.SetContentTypeJSON) // All responses from here will be JSON
	// r.Use(middleware.sessionMiddleware)

	r.Post("/signup", h.postSignup())
	r.Post("/login", h.postLogin())

	r.Route("/users", func(r chi.Router) {
		// r.Get("/", h.getSelf())
	})

	h.Mux = r
	return &h
}

// postSignup tries to sign up a user.
func (h *Handler) postSignup() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//TODO replace with env variable
		session, sess_err := store.Get(r, "session-cookie-name")
		if sess_err != nil {
			// TODO error handling
			http.Error(w, sess_err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO check if login is valid (i.e. account exists), if so log them in

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

		// jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		// if err != nil {
		// 	// TODO error handling
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }

		// TODO right now this is only valid on the domain it's sent from, if we do
		// subdomains (seems likely) then we'll need to change that.
		// http.SetCookie(w, &http.Cookie{
		// 	Name:       jwtCookie,
		// 	Value:      jwt,
		// 	Path:       "/",
		// 	RawExpires: "0",
		// })

		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Save(r, w)

	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Convert JSON user details in request to a user struct
		var u models.User

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&u)

		if err != nil {

			// TODO error handling

			session, err := store.Get(r, "cookie-name")
			if err != nil {

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Where authentication could be done
			if r.FormValue("code") != "code" {
				if r.FormValue("code") == "" {
					session.AddFlash("Must enter a code")
				}
				session.AddFlash("The code was incorrect")
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, "/forbidden", http.StatusFound)
				return
			}

			// No other errors
			session.Values["ID"] = u.ID
			err = session.Save(r, w) // TODO error checking

			username := r.FormValue("username")

			user := &User{
				Username:      username,
				Authenticated: true,
			}

			session.Values["user"] = user

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/secret", http.StatusFound)

		}
	}

	// getSelf returns the user as JSON.  It will redact some fields (like
	// password).
	// TODO Right now it sets the password value to blank but keeps the now blank
	// field in the JSON response.  Consider even removing that field.
	// TODO Same as above, but for passwordConfirm

	// func (h *Handler) getSelf() http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		// TODO replace cookie name with env name
	// 		session, err := store.Get(r, "session-cookie-name")
	// 		if err != nil {
	// 			// TODO error handling
	// 			http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 			return
	// 		}
	//
	// 		u, err := h.UserService.GetByEmail(claims["email"].(string))
	// 		if err != nil {
	// 			// TODO error handling
	// 			// This can fail either because the DB is messed up or nothing is found
	// 			// So be sure to deal with that
	// 			http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 			return
	// 		}
	//
	// 		json.NewEncoder(w).Encode(u)
	// 		return
	// 	}
	// }

	// mustGetJWTConfig tries to get JWT configuration variables from the
	// environment. It will panic if those variables are not set.
	// func mustGetJWTConfig() (string, []byte) {
	// 	jwtIssuer, ok := os.LookupEnv("JWT_COOKIE_NAME")
	// 	if !ok {
	// 		log.Fatalf("environment variable not set: %v", "JWT_ISSUER")
	// 	}
	//
	// 	jwtSigningKeyString, ok := os.LookupEnv("JWT_SIGNING_KEY")
	// 	if !ok {
	// 		log.Fatalf("environment variable not set: %v", "JWT_SIGNING_KEY")
	// 	}
	// 	jwtSigningKey := []byte(jwtSigningKeyString)
	//
	// 	return jwtIssuer, jwtSigningKey
	// }

	// getClaimsFromCtx returns the claims of a Context's JWT or an error.
	// func getClaimsFromCtx(ctx context.Context) (claims jwt.MapClaims, err error) {
	// 	// Always make sure the error field is nil
	// 	err, _ = ctx.Value(middleware.JWTErrorCtxKey).(error)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	// Make sure the token is not nil
	// 	token, ok := ctx.Value(middleware.JWTCtxKey).(*jwt.Token)
	// 	if !ok {
	// 		return nil, fmt.Errorf("missing jwt in context")
	// 	}
	//
	// 	claims = token.Claims.(jwt.MapClaims)
	// 	if err = claims.Valid(); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	return claims, err
	// }
}
