package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
)

type Handler struct {
	*chi.Mux

	UserService models.UserService
	templates   *template.Template
}

func NewHandler(us models.UserService) *Handler {
	h := Handler{UserService: us}
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
		// In prod only load them once
		err := h.loadTemplates()

		// And fail if they can't be loaded
		if err != nil {
			log.Fatalf("failed to load templates: %s", err)
		}
	}

	// All responses from here will have the Content-Type header be text/html
	r.Use(middleware.SetContentTypeHTML)

	r.Get("/", h.getRoot())

	h.Mux = r
	return &h
}

// Loads templates (files ending in .tmpl) from web/templates/ and all its
// subdirectories.
func (h *Handler) loadTemplates() (err error) {
	h.templates = template.New("")

	err = filepath.Walk("web/templates/", func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl") {
			template.Must(h.templates.ParseFiles(path))
		}

		return nil
	})

	return err
}

// reloadTemplates is a middleware that calls loadTemplates on every request.
// It's useful to live reload templates in development but shouldn't be used in
// production.
func (h *Handler) reloadTemplates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.loadTemplates()
		if err != nil {
			log.Fatalf("failed to reload templates: %s", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.templates.ExecuteTemplate(w, "index", nil)
	}
}
