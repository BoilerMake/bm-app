package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// Current status of app
	Status string

	// A generic place to put unstructured data
	Data interface{}

	// Values to be put back into a form when shown to a user again
	// For example, when they log in with the wrong password we want
	// the email they tried to log in with to persist
	FormRefill interface{}

	// The user's email, blank if user not logged in
	Email           string
	IsAuthenticated bool
}

func NewPage(title string, status string, r *http.Request) (*Page, bool) {
	session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
	if !ok {
		return nil, false
	}

	email, ok := session.Values["EMAIL"].(string)
	if !ok {
		// It's ok to ignore if this errors (for example when a user doesn't have a
		// session) because email will just default to the empty string.
	}

	return &Page{
		Title:           title,
		Status:          status,
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

	// HTML templates to render
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
	r.Get("/about", h.getAbout())
	r.Get("/faq", h.getFaq())

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

	/* EXEC ROUTES */
	r.Group(func(r chi.Router) {
		r.Use(middleware.MustBeExec)

		r.Get("/exec", h.getExec())
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
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake", status, r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "home", p)
	}
}

// getHackers renders the hackers template.
func (h *Handler) getHackers() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Hackers", status, r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "hackers", p)
	}
}

// getAbout renders the about template.
func (h *Handler) getAbout() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - About", status, r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "about", p)
	}
}

// getFaq renders the faq template.
func (h *Handler) getFaq() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - FAQ", status, r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "faq", p)
	}
}

// get404 handles requests that couldn't find a valid route by rendering the
// 404 template.
func (h *Handler) get404() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - 404", status, r)
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
func mustGetEnv(varName string) (value string) {
	value, ok := os.LookupEnv(varName)
	if !ok {
		log.Fatalf("environment variable not set: %v", varName)
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
		}

		return "/404"
	}
}

// onSeasonOnly checks to make sure status is between 2 and 4, meaning
// that the event is either happening right now or is upcoming.  If it is not
// then the route is redirected to a 404.  If the application status is not on
// season then this will return an error.
func onSeasonOnly(status string) error {
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return err
	}

	if statusInt < 2 || statusInt > 4 {
		return fmt.Errorf("application is currently not in season")
	}

	return nil
}
