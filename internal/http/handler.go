package http

import (
	"fmt"
	"net/http"

	"github.com/BoilerMake/bm-app/internal/http/middleware"
	"github.com/BoilerMake/bm-app/internal/http/web"
	"github.com/BoilerMake/bm-app/internal/mail"
	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/internal/s3"

	"github.com/go-chi/chi"
)

// A Handler ties together the handlers found in api and web.
type Handler struct {
	*chi.Mux

	WebHandler *web.Handler
}

// NewHandler creates a handler that wraps the subhandlers for the entire app.
func NewHandler(us models.UserService, as models.ApplicationService, rs models.RSVPService, anns models.AnnouncementService, ras models.RaffleService, mailer mail.Mailer, S3 s3.S3) *Handler {
	h := Handler{
		WebHandler: web.NewHandler(us, as, rs, anns, ras, mailer, S3),
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
	r.Get("/service-worker.js", h.getInvalidateServiceWorker())

	h.Mux = r
	return &h
}

// getRobotsTxt serves the plain text robots.txt file for web crawlers.
// The file itself is defined as a string in the handler.
func (h *Handler) getRobotsTxt() http.HandlerFunc {
	robotsTxt := "User-agent: *\n" +
		"Disallow:\n" +
		"Crawl-delay: 10"

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, robotsTxt)
	}
}

// getInvalidateServiceWorker a really bad service worker.  But hey, at least
// it makes things less broken.
func (h *Handler) getInvalidateServiceWorker() http.HandlerFunc {
	badServiceWorker := []byte(`self.addEventListener('install', () => {
	self.skipWaiting();
});`)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(badServiceWorker)
	}
}
