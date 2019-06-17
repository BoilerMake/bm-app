package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/mailgun/mailgun-go/v3"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

var (
	jwtCookie string // Name for the JWT's cookie.  TODO Better name?
	mgDomain  string
	mgAPIKey  string
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
	mgDomain, ok = os.LookupEnv("MAILGUN_DOMAIN")
	if !ok {
		log.Fatalf("environment variable not set: %v", "MAILGUN_DOMAIN")
	}
	mgAPIKey, ok = os.LookupEnv("MAILGUN_API_KEY")
	if !ok {
		log.Fatalf("environment variable not set: %v", "MAILGUN_API_KEY")
	}

	r.Use(middleware.SetContentTypeJSON) // All responses from here will be JSON
	r.Use(middleware.WithJWT)

	r.Post("/signup", h.postSignup())
	r.Post("/login", h.postLogin())
	r.Post("/forgot", h.postForgotPassword())
	r.Post("/reset-password", h.postResetPassword())

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.getSelf())
	})

	h.Mux = r
	return &h
}

// postSignup tries to sign up a user.
func (h *Handler) postSignup() http.HandlerFunc {
	jwtIssuer, jwtSigningKey := mustGetJWTConfig()

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {
	jwtIssuer, jwtSigningKey := mustGetJWTConfig()

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// postForgotPassword sends the password reset email
func (h *Handler) postForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get info from body
		var emailModel models.EmailModel
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&emailModel)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := h.UserService.GetPasswordReset(emailModel.Email)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// This will be formatted better once the front end is setup for the link
		recipient := emailModel.Email
		subject := "Password Reset"
		body := "Your reset token is: " + token

		// Sending the email (can be made into a method)
		mg := mailgun.NewMailgun(mgDomain, mgAPIKey)
		// Maybe make this env var
		sender := "team@boilermake.org"
		message := mg.NewMessage(sender, subject, body, recipient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		// Send the message	with a 10 second timeout
		resp, id, err := mg.Send(ctx, message)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}
}

// postResetPassword resets the password with a valid token
func (h *Handler) postResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get info from body
		var passwordResetInfo models.PasswordResetPayload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&passwordResetInfo)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = h.UserService.ResetPassword(passwordResetInfo.UserToken, passwordResetInfo.NewPassword)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// getSelf returns the user as JSON.  It will redact some fields (like
// password).
// TODO Right now it sets the password value to blank but keeps the now blank
// field in the JSON response.  Consider even removing that field.
// TODO Same as above, but for passwordConfirm
func (h *Handler) getSelf() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			http.Error(w, "unauthorized", http.StatusUnauthorized)
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

		json.NewEncoder(w).Encode(u)
		return
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
