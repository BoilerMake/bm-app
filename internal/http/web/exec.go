package web

import (
	"net/http"
)

// getExec renders the exec home page.
func (h *Handler) getExec() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)

		p, ok := NewPage(w, r, "BoilerMake - Exec", session)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "exec", p)
	}
}
