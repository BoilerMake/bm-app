package web

import (
	"database/sql"
	"net/http"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/gorilla/sessions"
)

// getApply renders the apply template.
func (h *Handler) getApply() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")
	err := onSeasonOnly(status)
	if err != nil {
		return h.get404()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO error handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		id, ok := session.Values["ID"].(int)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "invalid session value", http.StatusInternalServerError)
			return
		}

		app, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			// If the error was that there is no application for this user, just render
			// the blank application form
			if err == sql.ErrNoRows {
				app = &models.Application{}
			} else {
				// TODO error handling
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		p, ok := NewPage("BoilerMake - Apply", status, r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		p.FormRefill = app

		// Otherwise we can show the apply form with the data already filled in
		h.Templates.RenderTemplate(w, "apply", p)
	}
}

// postApply tries to create an application from a post request.
func (h *Handler) postApply() http.HandlerFunc {
	status := mustGetEnv("APP_STATUS")
	err := onSeasonOnly(status)
	if err != nil {
		return h.get404()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var a models.Application
		err := a.FromFormData(r)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO error handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		a.UserID, ok = session.Values["ID"].(int)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "invalid session value", http.StatusInternalServerError)
			return
		}

		err = h.ApplicationService.CreateOrUpdate(&a)
		if err != nil {
			// TODO error handling
			// TODO once session tokens are updated this should show a failure flash
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If creating/saving application worked, upload resume only if a new resume
		// was sent to us.
		if a.ResumeFile != "" {
			err = h.S3.UploadResume(a.UserID, a.Resume)
			if err != nil {
				// TODO error handling
				// TODO once session tokens are updated this should show a failure flash
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Redirect back to application page if successful
		// TODO once session tokens are updated this should show success and give a date for when apps are locked
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}
