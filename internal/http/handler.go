package http

import (
	"net/http"

	"github.com/BoilerMake/new-backend/internal/http/api"
	"github.com/BoilerMake/new-backend/internal/http/web"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
)

// A Handler ties together the handlers found in api and web.
type Handler struct {
	*chi.Mux

	APIHandler *api.Handler
	WebHandler *web.Handler
}

// NewHandler creates a handler that wraps the subhandlers for the entire app.
func NewHandler(us models.UserService, mailer mail.Mailer) *Handler {
	h := Handler{
		APIHandler: api.NewHandler(us, mailer),
		WebHandler: web.NewHandler(us),
	}
	r := chi.NewRouter()

	// Static file serving out of web/static/
	fs := http.StripPrefix("/static", http.FileServer(http.Dir("web/static")))
	r.Get("/static/*", fs.ServeHTTP)

	// TODO historically we have used subdomains (like api.boilermake.org)
	r.Mount("/api", h.APIHandler)
	r.Mount("/", h.WebHandler)

	h.Mux = r
	return &h
}
