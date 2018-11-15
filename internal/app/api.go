package app

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO these should be moved into middleware that controls default headers
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode("pong")
	}
}
