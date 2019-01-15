package api

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

	r.Get("/", h.getRoot)

	h.Mux = r
	return &h
}

// getRoot is an example handler endpoint. It server get requests at hostname.com/api
func (h *Handler) getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ya like API tests?")
}
