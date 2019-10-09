package web

import (
	"errors"
	"net/http"
)

// getExec renders the exec home page.
func (h *Handler) getExec() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)
		p, ok := NewPage(w, r, "BoilerMake - Exec", status, session)

		if !ok {
			h.Error(w, r, errors.New("creating page failed"))
			return
		}

		h.Templates.RenderTemplate(w, "exec", p)
	}
}
