package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/internal/s3"
	"github.com/BoilerMake/new-backend/pkg/template"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

// A Handler will route requests to their appropriate HandlerFunc.
type Handler struct {
	// A Mux handles all routing and middleware.
	*chi.Mux

	// A UserService is the interface with the database.
	UserService models.UserService

	// An ApplicationService is the interface with the databsae
	ApplicationService models.ApplicationService

	// A Mailer is used to send emails
	Mailer mail.Mailer

	// An S3 handles uploading files to S3
	S3 s3.S3

	Templates *template.Template
}

// NewHandler creates a handler for web requests.
func NewHandler(us models.UserService, as models.ApplicationService, mailer mail.Mailer, S3 s3.S3) *Handler {
	h := Handler{
		UserService:        us,
		ApplicationService: as,
		Mailer:             mailer,
		S3:                 S3,
	}

	r := chi.NewRouter()

	// Set up templates
	tmplPath, ok := os.LookupEnv("TEMPLATES_PATH")
	if !ok {
		log.Fatalf("environment variable not set: %v", "TEMPLATES_PATH")
	}

	tmpls, err := template.NewTemplate(tmplPath)
	if err != nil {
		log.Fatalf("failed to load templates: %s", err)
	}
	h.Templates = tmpls

	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v", "ENV_MODE")
	}

	if mode == "development" {
		// In development mode, reload templates on every request
		r.Use(h.Templates.ReloadTemplates)
	}

	// Set up pool of buffers used for rendering templates.
	r.Use(middleware.WithJWT)

	r.Get("/", h.getRoot())

	/* USER ROUTES */
	r.Get("/signup", h.getSignup())
	r.Post("/signup", h.postSignup())

	r.Get("/activate/{code}", h.getActivate())

	r.Get("/forgot", h.getForgotPassword())
	r.Post("/forgot", h.postForgotPassword())
	r.Get("/reset", h.getResetPassword())
	r.Get("/reset/{token}", h.getResetPasswordWithToken())
	r.Post("/reset/{token}", h.postResetPassword())

	r.Get("/login", h.getLogin())
	r.Post("/login", h.postLogin())

	r.Get("/account", h.getAccount())

	/* APPLICATION ROUTES */
	r.Get("/apply", h.getApply())
	r.Post("/apply", h.postApply())

	r.NotFound(h.get404())

	h.Mux = r
	return &h
}

// getRoot renders the index template.
func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Templates.RenderTemplate(w, "index", nil)
	}
}

// get404 handles requests that couldn't find a valid route by rendering the
// 404 template.
func (h *Handler) get404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Templates.RenderTemplate(w, "404", nil)
	}
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

// mustGetEnv looks up and sets an environment variable.  If the environment
// variable is not found, it panics.
func mustGetEnv(var_name string) (value string) {
	value, ok := os.LookupEnv(var_name)
	if !ok {
		log.Fatalf("environment variable not set: %v", var_name)
	}
	return value
}
