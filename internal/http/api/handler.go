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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"

	"github.com/mailgun/mailgun-go/v3"
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
	jwtCookie = getEnv("JWT_COOKIE_NAME")
	h := Handler{UserService: us}
	r := chi.NewRouter()

	// TODO See cmd/server/main.go for more about config. This doesn't seem ideal.

	r.Use(middleware.SetContentTypeJSON) // All responses from here will be JSON
	r.Use(middleware.WithJWT)

	r.Post("/signup", h.postSignup())
	r.Post("/login", h.postLogin())
	r.Post("/activate/{code}", h.postActivate())

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.getSelf())
	})

	h.Mux = r
	return &h
}

// postSignup tries to sign up a user.
func (h *Handler) postSignup() http.HandlerFunc {
	// Get necessary environment variables for mailgun
	sender := getEnv("EMAIL_ADDRESS")
	domain := getEnv("DOMAIN_NAME")
	private_key := getEnv("MAILGUN_KEY")

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

		id, confirmationCode, err := h.UserService.Signup(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.ID = id

		// Create an instance of the Mailgun Client
		mg := mailgun.NewMailgun(domain, private_key)

		// Build confirmation email
		subject := "Confirm your email"
		link := "localhost:8080/api/activate/" + string(confirmationCode)
		body := "Please click the following link to confirm your email address: " + link
		recipient := u.Email
		message := mg.NewMessage(sender, subject, body, recipient)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Send the message	with a 10 second timeout
		resp, mid, err := mg.Send(ctx, message)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Printf("ID: %s Resp: %s\n", mid, resp)

		jwt, err := u.GetJWT()
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

		jwt, err := u.GetJWT()
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

func (h *Handler) postActivate() http.HandlerFunc {
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

// getEnv looks up and sets an environment variable
func getEnv(var_name string) (value string) {
	value, ok := os.LookupEnv(var_name)
	if !ok {
		log.Fatalf("environment variable not set: %v", var_name)
	}
	return value
}
