package http

import (
	"fmt"
	"net/http"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/http/web"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/internal/s3"

	"github.com/go-chi/chi"
)

// A Handler ties together the handlers found in api and web.
type Handler struct {
	*chi.Mux

	WebHandler *web.Handler
}

// NewHandler creates a handler that wraps the subhandlers for the entire app.
func NewHandler(us models.UserService, as models.ApplicationService, mailer mail.Mailer, S3 s3.S3) *Handler {
	h := Handler{
		WebHandler: web.NewHandler(us, as, mailer, S3),
	}
	r := chi.NewRouter()

	// Limit body request size
	r.Use(middleware.LimitRequestSize)
	// Log some stuff for each request
	r.Use(middleware.Logging)
	// Recover from panics
	r.Use(middleware.Recoverer)

	r.Mount("/", h.WebHandler)

	r.Get("/robots.txt", h.getRobotsTxt())

	h.Mux = r
	return &h
}

// getRobotsTxt serves the plain text robots.txt file for web crawlers.
// The file itself is defined as a string in the handler.
func (h *Handler) getRobotsTxt() http.HandlerFunc {
	robotsTxt := `User-agent: *
Disallow:
Crawl-delay: 10`

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, robotsTxt)
	}
}
