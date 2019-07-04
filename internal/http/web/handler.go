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
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

// A Handler will route requests to their appropriate HandlerFunc.
type Handler struct {
	// A Mux handles all routing and middleware.
	*chi.Mux

	// A UserService is the interface with the database.
	UserService models.UserService

	// A Mailer is used to send emails
	Mailer mail.Mailer

	templates       *template.Template
	templateBufPool *sync.Pool
}

// NewHandler creates a handler for web requests.
func NewHandler(us models.UserService, mailer mail.Mailer) *Handler {
	h := Handler{
		UserService: us,
		Mailer:      mailer,
	}

	r := chi.NewRouter()

	// Set up templates
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v", "ENV_MODE")
	}

	if mode == "development" {
		// In development mode, reload templates on every request
		r.Use(h.reloadTemplates)
	} else {
		// In prod only load templates once
		err := h.loadTemplates()

		// And fail if they can't be loaded
		if err != nil {
			log.Fatalf("failed to load templates: %s", err)
		}
	}

	// Set up pool of buffers used for rendering templates.
	h.templateBufPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	r.Use(middleware.WithJWT)

	r.Get("/", h.getRoot())

	r.Get("/signup", h.getSignup())
	r.Post("/signup", h.postSignup())

	r.Get("/activate/{code}", h.getActivate())

	r.Post("/forgot", h.postForgotPassword())
	r.Post("/reset-password/{token}", h.postResetPassword())

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
	domain := mustGetEnv("DOMAIN")

	mode := mustGetEnv("ENV_MODE")
	if mode == "development" {
		domain += ":" + mustGetEnv("PORT")
	}

	jwtIssuer := mustGetEnv("JWT_ISSUER")
	jwtSigningKey := []byte(mustGetEnv("JWT_SIGNING_KEY"))
	jwtCookie := mustGetEnv("JWT_COOKIE_NAME")

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO check if login is valid (i.e. account exists), if so log them in
		var u models.User
		u.FromFormData(r)

		id, confirmationCode, err := h.UserService.Signup(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.ID = id

		// Build confirmation email
		to := u.Email
		subject := "Confirm your email"
		link := domain + "/activate/" + confirmationCode
		body := "Please click the following link to confirm your email address: " + link

		err = h.Mailer.Send(to, subject, body)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

// getActivate activates the account that corresponds to the activation code
// if there is such an account.
func (h *Handler) getActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Temporarily ignoring claims returned from getClaimsFromCtx
		_, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		code := chi.URLParam(r, "code")

		u, err := h.UserService.GetByCode(code)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u.IsActive = true
		u.ConfirmationCode = ""
		err = h.UserService.Update(u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// TODO should this redirect to an "activate successful" page?
		// Redirect to homepage if activation was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// postForgotPassword sends the password reset email.
func (h *Handler) postForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		token, err := h.UserService.GetPasswordReset(u.Email)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO This will need to be formatted better once the front end is setup for the link
		to := u.Email
		subject := "Password Reset"
		body := "Your reset token is: " + token

		err = h.Mailer.Send(to, subject, body)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// postResetPassword resets the password with a valid token
func (h *Handler) postResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var passwordResetInfo models.PasswordResetPayload
		passwordResetInfo.UserToken = chi.URLParam(r, "code")
		passwordResetInfo.NewPassword = r.FormValue("new-password")

		err := h.UserService.ResetPassword(passwordResetInfo.UserToken, passwordResetInfo.NewPassword)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	jwtIssuer := mustGetEnv("JWT_ISSUER")
	jwtSigningKey := []byte(mustGetEnv("JWT_SIGNING_KEY"))
	jwtCookie := mustGetEnv("JWT_COOKIE_NAME")

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
			return
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
		claims, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		u, err := h.UserService.GetByEmail(claims["email"].(string))
		if err != nil {
			// TODO error handling
			// This can fail either because the DB is messed up or nothing is found
			// So be sure to deal with that
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
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
