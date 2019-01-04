package http

import (
	"net/http"
	"strings"

	"github.com/BoilerMake/new-backend/internal/http/api"
	"github.com/BoilerMake/new-backend/internal/http/web"
)

// Handler should tie together the handlers found in rest and web.
type Handler struct {
	APIHandler *api.Handler
	WebHandler *web.Handler
}

// ServeHTTP routes a request to the appropriate handler.
// TODO seems like this can be replaced by a chi router
// TODO historically we have used subdomains (like api.boilermake.org)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		h.APIHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/") {
		h.WebHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
