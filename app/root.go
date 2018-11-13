package app

import (
	"fmt"
	"net/http"
)

func (s *Server) getRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "what it do")
	}
}
