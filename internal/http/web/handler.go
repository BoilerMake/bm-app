package web

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

var (
	jwtCookie string // Name for the JWT's cookie.  TODO Better name?
)

type Handler struct {
	*chi.Mux

	UserService models.UserService

	templates       *template.Template
	templateBufPool *sync.Pool
}

func NewHandler(us models.UserService) *Handler {
	h := Handler{UserService: us}
	r := chi.NewRouter()

	// TODO See cmd/server/main.go for more about config. This doesn't seem ideal.
	var ok bool
	jwtCookie, ok = os.LookupEnv("JWT_COOKIE_NAME")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_COOKIE_NAME")
	}

	// Set up templates
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v", "ENV_MODE")
	}

	if mode == "development" {
		// In development mode, reload templates on every request
		r.Use(h.reloadTemplates)
	} else {
		// In prod only load them once
		err := h.loadTemplates()

		// And fail if they can't be loaded
		if err != nil {
			log.Fatalf("failed to load templates: %s", err)
		}
	}

	// Set up pool of buffers used for rendering templates
	h.templateBufPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	r.Use(middleware.WithJWT)

	r.Get("/", h.getRoot())

	r.Get("/signup", h.getSignup())
	r.Post("/signup", h.postSignup())
	r.Get("/login", h.getLogin())
	r.Post("/login", h.postLogin())
	r.Get("/account", h.getAccount())

	r.NotFound(h.get404())

	h.Mux = r
	return &h
}

// getRoot renders the index template.
func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.renderTemplate(w, "index", nil)
	}
}

// getSignup renders the signup template.
func (h *Handler) getSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.renderTemplate(w, "signup", nil)
	}
}

// postSignup tries to signup a user from a post request.
func (h *Handler) postSignup() http.HandlerFunc {
	jwtIssuer, jwtSigningKey := mustGetJWTConfig()
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		id, err := h.UserService.Signup(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.ID = id

		jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// TODO right now this is only valid on the domain it's sent from, if we do
		// subdomains (seems likely) then we'll need to change that.
		http.SetCookie(w, &http.Cookie{
			Name:       jwtCookie,
			Value:      jwt,
			Path:       "/",
			RawExpires: "0",
		})

		// Redirect to homepage if signup was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getLogin renders the login template.
func (h *Handler) getLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.renderTemplate(w, "login", nil)
	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {
	jwtIssuer, jwtSigningKey := mustGetJWTConfig()

	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		err := h.UserService.Login(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// TODO Right now this is only valid on the domain it's sent from, if we do
		// subdomains (seems likely) then we'll need to change that.
		http.SetCookie(w, &http.Cookie{
			Name:       jwtCookie,
			Value:      jwt,
			Path:       "/",
			RawExpires: "0",
		})

		// Redirect to homepage if login was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getAccount shows a user their account.
func (h *Handler) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO redirect to login if they're not logged in
		claims, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}

		u, err := h.UserService.GetByEmail(claims["email"].(string))
		if err != nil {
			// TODO error handling
			// This can fail either because the DB is messed up or nothing is found
			// So be sure to deal with that
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}

		data := map[string]interface{}{
			"Email":       u.Email,
			"FirstName":   u.FirstName,
			"LastName":    u.LastName,
			"Phone":       u.Phone,
			"ProjectIdea": u.ProjectIdea,
			"TeamMembers": u.TeamMembers,
		}

		h.renderTemplate(w, "account", data)
	}
}

// get404 handles requests that couldn't find a valid route by rendering the
// 404 template.
func (h *Handler) get404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.renderTemplate(w, "404", nil)
	}
}

// mustGetJWTConfig tries to get JWT configuration variables from the
// environment. It will panic if those variables are not set.
func mustGetJWTConfig() (string, []byte) {
	jwtIssuer, ok := os.LookupEnv("JWT_COOKIE_NAME")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_ISSUER")
	}

	jwtSigningKeyString, ok := os.LookupEnv("JWT_SIGNING_KEY")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_SIGNING_KEY")
	}
	jwtSigningKey := []byte(jwtSigningKeyString)

	return jwtIssuer, jwtSigningKey
}

// getClaimsFromCtx returns the claims of a Context's JWT or an error.
func getClaimsFromCtx(ctx context.Context) (claims jwt.MapClaims, err error) {
	// Always make sure the error field is nil
	err, _ = ctx.Value(middleware.JWTErrorCtxKey).(error)
	if err != nil {
		return nil, err
	}

	// Make sure the token is not nil
	token, ok := ctx.Value(middleware.JWTCtxKey).(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("missing jwt in context")
	}

	claims = token.Claims.(jwt.MapClaims)
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, err
}
