package web

import (
	"fmt"
	"net/http"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/go-chi/chi"
)

type Handler struct {
	*chi.Mux

	UserService models.UserService
	//EventService models.EventService
}

func NewHandler(us models.UserService) *Handler {
	h := Handler{UserService: us}
	r := chi.NewRouter()

	r.Get("/", h.getRoot())

	h.Mux = r
	return &h
}

func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ya like web tests?")
	}
}
