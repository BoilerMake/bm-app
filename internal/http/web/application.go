package web

import (
	"database/sql"
	"net/http"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/rollbar/rollbar-go"

	"github.com/gorilla/sessions"
)

// getApply renders the apply template.
func (h *Handler) getApply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			rollbar.Error("invalid session value")
			rollbar.Wait()
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		id, ok := session.Values["ID"].(int)
		if !ok {
			rollbar.Error("invalid session value")
			rollbar.Wait()
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
				rollbar.Error(err)
				rollbar.Wait()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		p, ok := NewPage("BoilerMake - Apply", r)
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
	return func(w http.ResponseWriter, r *http.Request) {
		var a models.Application
		err := a.FromFormData(r)
		if err != nil {
			rollbar.Error(err, r)
			rollbar.Wait()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			rollbar.Error(err, r)
			rollbar.Wait()
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		a.UserID, ok = session.Values["ID"].(int)
		if !ok {
			rollbar.Error(err, r)
			rollbar.Wait()
			http.Error(w, "invalid session value", http.StatusInternalServerError)
			return
		}

		err = h.ApplicationService.CreateOrUpdate(&a)
		if err != nil {
			rollbar.Error(err)
			rollbar.Wait()
			// TODO once session tokens are updated this should show a failure flash
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If creating/saving application worked, upload resume only if a new resume
		// was sent to us.
		if a.ResumeFile != "" {
			err = h.S3.UploadResume(a.UserID, a.Resume)
			if err != nil {
				rollbar.Error(err)
				rollbar.Wait()
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
