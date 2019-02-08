package api

import (
	"encoding/json"
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
	r.Get("/users/{id}", h.getUser())

	h.Mux = r
	return &h
}

// getRoot is an example handler endpoint. It server get requests at hostname.com/api
func (h *Handler) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ya like API tests?")
	}
}

// getUser returns a user given an id.  It should only allow a user to get
// themselves
// TODO auth
func (h *Handler) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			// TODO error handling
			http.Error(w, "no id given", http.StatusInternalServerError)
			return
		}

		u, err := h.UserService.GetById(id)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}
