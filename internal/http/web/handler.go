package web

import (
	"encoding/json"
	"errors"
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
	"github.com/BoilerMake/new-backend/pkg/flash"
	"github.com/BoilerMake/new-backend/pkg/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/rollbar/rollbar-go"
)

type ErrorReporter func(interfaces ...interface{})

// A Page is all the data needed to render a page.
type Page struct {
	Title string

	// Current status of app
	Status string

	// A generic place to put unstructured data
	Data interface{}

	// An array of messages to show the user
	Flashes []flash.Flash

	// Values to be put back into a form when shown to a user again
	// For example, when they log in with the wrong password we want
	// the email they tried to log in with to persist
	FormRefill interface{}

	// The user's email, blank if user not logged in
	Email           string
	IsAuthenticated bool
}

func NewPage(w http.ResponseWriter, r *http.Request, title string, status string, session *sessions.Session) (*Page, bool) {
	email, ok := session.Values["EMAIL"].(string)
	if !ok {
		// It's ok to ignore if this errors (for example when a user doesn't have a
		// session) because email will just default to the empty string.
	}

	var flashes []flash.Flash
	flashesInterface := session.Flashes()
	session.Save(r, w)

	for _, e := range flashesInterface {
		f, ok := e.(flash.Flash)
		if ok {
			flashes = append(flashes, f)
		}
	}

	return &Page{
		Title:           title,
		Status:          status,
		Email:           email,
		IsAuthenticated: email != "",
		Flashes:         flashes,
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

	// An ErrorReport reports errors to something like rollbar
	ErrReporter ErrorReporter

	// Stores session cookies
	SessionStore *sessions.CookieStore

	// Name cookie that stores sessions
	SessionCookieName string
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

		/* EXEC ROUTES */
		r.Group(func(r chi.Router) {
			r.Use(middleware.MustBeExec)

			r.Get("/exec", h.getExec())
		})
	})

	if mode == "development" {
		// In prod we serve static items through a CDN, in development just serve
		// out of web/static/ at /static/
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("web/static")))
		r.Get("/static/*", fs.ServeHTTP)
	}

	// Only log to rollbar in production
	rollbarEnv := mustGetEnv("ROLLBAR_ENVIRONMENT")

	if rollbarEnv == "production" {
		h.ErrReporter = rollbarReportError
	} else {
		// If we're not in production just print out the errors
		h.ErrReporter = logReportError
	}

	// Set up session store
	sessionSecret := mustGetEnv("SESSION_SECRET")

	sessionKey := []byte(sessionSecret)
	store := sessions.NewCookieStore(sessionKey)

	// Prevents CSRF attacks (on browsers that support SameSite)
	store.Options.SameSite = http.SameSiteStrictMode

	// Prevents XSS attacks (JS isn't allowed to access cookie)
	store.Options.HttpOnly = true

	if mode != "development" {
		// Only transfer cookie over https
		store.Options.Secure = true
	}
	h.SessionStore = store

	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	h.SessionCookieName = sessionCookieName

	r.NotFound(h.get404())

	h.Mux = r
	return &h
}

// getRoot renders the index template.
func (h *Handler) getRoot() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
			return
		}

		h.Templates.RenderTemplate(w, "home", p)
	}
}

// getHackers renders the hackers template.
func (h *Handler) getHackers() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake - Hackers", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
			return
		}

		h.Templates.RenderTemplate(w, "hackers", p)
	}
}

// getAbout renders the about template.
func (h *Handler) getAbout() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake - About", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
			return
		}

		h.Templates.RenderTemplate(w, "about", p)
	}
}

// getFaq renders the faq template.
func (h *Handler) getFaq() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake - FAQ", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
			return
		}

		h.Templates.RenderTemplate(w, "faq", p)
	}
}

// get404 handles requests that couldn't find a valid route by rendering the
// 404 template.
func (h *Handler) get404() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake - 404", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
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

// rollbarReportError reports an error to rollbar and logs it locally.  It
// should only be reporting errors in production.  You should not call this
// function directly, instead call h.Error(...) and let that handle it.
func rollbarReportError(interfaces ...interface{}) {
	rollbar.Error(interfaces...)
	rollbar.Wait()

	// Also log the error locally
	log.Println("ERROR:", interfaces)
}

// logReportError logs an error locally.  In production rollbarReportError
// should be used instead.  You should not call this function directly, instead
// call h.Error(...) and let that handle it.
func logReportError(interfaces ...interface{}) {
	// Also log the error locally
	log.Println("ERROR:", interfaces)
}

// Error checks an error given to it.  If it's a known error that we've made
// we can show it to the user as a flash.  If it's unknown then we should tell
// the user that something went wrong and report the error to rollbar.
func (h *Handler) Error(w http.ResponseWriter, r *http.Request, err error, interfaces ...interface{}) {
	switch err.(type) {
	case *models.ModelError:
		modelError := err.(*models.ModelError)

		// This is an error we know about and should let the user know what happened
		session, _ := h.SessionStore.Get(r, h.SessionCookieName)

		session.AddFlash(flash.Flash{
			Type:    modelError.GetType(),
			Message: modelError.Error(),
		})
		session.Save(r, w)

		// Redirect to previous page
		http.Redirect(w, r, r.URL.RequestURI(), http.StatusSeeOther)
	default:
		// Because we don't know how this error happened, we should report it on rollbar.
		h.ErrReporter(append([]interface{}{err}, interfaces...)...)

		// This error could have come from anywhere, so we should just tell the user
		// something went wrong so that we don't accidently expose something
		// sensitive
		h.Templates.RenderTemplate(w, "500", nil)
	}
}
