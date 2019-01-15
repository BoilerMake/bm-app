package http

import (
	"github.com/BoilerMake/new-backend/internal/http/api"
	"github.com/BoilerMake/new-backend/internal/http/web"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
)

// Handler should tie together the handlers found in rest and web.
type Handler struct {
	*chi.Mux

	APIHandler *api.Handler
	WebHandler *web.Handler
}

func NewHandler(us models.UserService) *Handler {
	h := Handler{
		APIHandler: api.NewHandler(us),
		WebHandler: web.NewHandler(us),
	}
	r := chi.NewRouter()

	// TODO historically we have used subdomains (like api.boilermake.org)
	r.Mount("/api", h.APIHandler)
	r.Mount("/", h.WebHandler)

	h.Mux = r
	return &h
}
