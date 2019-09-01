package http

import (
	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/http/web"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
)

// A Handler ties together the handlers found in api and web.
type Handler struct {
	*chi.Mux

	WebHandler *web.Handler
}

// NewHandler creates a handler that wraps the subhandlers for the entire app.
func NewHandler(us models.UserService, mailer mail.Mailer) *Handler {
	h := Handler{
		WebHandler: web.NewHandler(us, mailer),
	}
	r := chi.NewRouter()

	r.Use(middleware.WithSession)

	r.Mount("/", h.WebHandler)

	h.Mux = r
	return &h
}
