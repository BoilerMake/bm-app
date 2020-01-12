package web

import (
	"errors"
	"net/http"
)

// getExec renders the exec home page.
func (h *Handler) getExec() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Exec")
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		applicationCount := h.ApplicationService.GetApplicationCount()
		userCount := h.UserService.GetUserCount()

		p.Data = map[string]interface{}{
			"ApplicationCount": applicationCount,
			"UserCount":        userCount,
		}

		h.Templates.RenderTemplate(w, "exec", p)

	}

}
