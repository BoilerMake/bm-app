package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/internal/s3"
	"github.com/BoilerMake/new-backend/pkg/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

// A Page is all the data needed to render a page.
type Page struct {
	Title string

	// A generic place to put unstructured data
	Data interface{}

	// The user's email, blank if user not logged in
	Email           string
	IsAuthenticated bool
}

func NewPage(title string, r *http.Request) (*Page, bool) {
	session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
	if !ok {
		return nil, false
	}

	email, ok := session.Values["EMAIL"].(string)
	if !ok {
		// Without this, email is set TODO
		email = ""
	}

	return &Page{
		Title:           title,
		Email:           email,
		IsAuthenticated: email != "",
	}, true
}

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

	// Cookiestore for session
	CookieStore *sessions.CookieStore
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
	mode := mustGetEnv("ENV_MODE")
	tmplPath := mustGetEnv("TEMPLATES_PATH")
	tmplFuncs := map[string]interface{}{
		"static_path": staticFileReplace(mode),
	}

	tmpls, err := template.NewTemplate(tmplPath, tmplFuncs)
	if err != nil {
		log.Fatalf("failed to load templates: %s", err)
	}
	h.Templates = tmpls

	if mode == "development" {
		// In development mode, reload templates on every request
		r.Use(h.Templates.ReloadTemplates)
	}

	/* WEB ROUTES */
	r.Get("/", h.getRoot())
	r.Get("/hackers", h.getHackers())

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

	r.Get("/logout", h.getLogout())

	// Must have auth
	r.Group(func(r chi.Router) {
		r.Use(middleware.MustBeAuthenticated)

		r.Get("/account", h.getAccount())

		/* APPLICATION ROUTES */
		r.Get("/apply", h.getApply())
		r.Post("/apply", h.postApply())
	})

	if mode == "development" {
		// In prod we serve static items through a CDN, in development just serve
		// out of web/static/ at /static/
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("web/static")))
		r.Get("/static/*", fs.ServeHTTP)
	}

	r.NotFound(h.get404())

	h.Mux = r
	return &h
}

// getRoot renders the index template.
func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", p)

		h.Templates.RenderTemplate(w, "index", p)
	}
}

// getHackers renders the hackers template.
func (h *Handler) getHackers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Hackers", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "hackers", p)
	}
}

// get404 handles requests that couldn't find a valid route by rendering the
// 404 template.
func (h *Handler) get404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - 404", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "404", p)
	}
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

// staticFileReplace is a template helper that reads in a manifest file and
// reroutes requests accordingly.  It's useful when you upload versioned files
// to a CDN for cache invalidation.  The manifest file used is made by gulp
// when running the prod config.
func staticFileReplace(mode string) func(path string) string {
	if mode == "development" {
		// No need to change path in development
		return func(path string) string {
			return "/static/" + path
		}
	}

	// In prod change path to cloudfront
	cloudfrontURL := mustGetEnv("CLOUDFRONT_URL")

	file, err := ioutil.ReadFile("web/static/rev-manifest.json")
	if err != nil {
		log.Fatalf("failed to read static manifest file: %v", err)
	}

	var manifest map[string]interface{}
	err = json.Unmarshal(file, &manifest)
	if err != nil {
		log.Fatalf("failed to parse static manifest file: %v", err)
	}

	return func(path string) string {
		if manifest[path] != nil {
			return cloudfrontURL + "/" + manifest[path].(string)
		} else {
			return "/404"
		}
	}
}
