package web

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/pkg/flash"
)

// getApply renders the apply template.
func (h *Handler) getApply() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")
	err := onSeasonOnly(status)
	if err != nil {
		return h.get404()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, sessionCookieName)

		id, ok := session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		app, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			// If the error was that there is no application for this user, just render
			// the blank application form
			if err == sql.ErrNoRows {
				app = &models.Application{}
			} else {
				h.Error(w, r, err, "")
				return
			}
		}

		p, ok := NewPage(w, r, "BoilerMake - Apply", status, session)
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		p.FormRefill = app

		// Otherwise we can show the apply form with the data already filled in
		h.Templates.RenderTemplate(w, "apply", p)
	}
}

// postApply tries to create an application from a post request.
func (h *Handler) postApply() http.HandlerFunc {
	sessionCookieName := mustGetEnv("SESSION_COOKIE_NAME")
	status := mustGetEnv("APP_STATUS")
	err := onSeasonOnly(status)
	if err != nil {
		return h.get404()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		var a models.Application

		err := a.FromFormData(r)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		session, _ := h.SessionStore.Get(r, sessionCookieName)

		a.UserID, ok = session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		err = h.ApplicationService.CreateOrUpdate(&a)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		// If creating/saving application worked, upload resume only if a new resume
		// was sent to us.
		if a.ResumeFile != "" {
			err = h.S3.UploadResume(a.UserID, a.Resume)
			if err != nil {
				h.Error(w, r, err, "")
				return
			}
		}

		// Redirect back to application page if successful
		// Also show a success flash
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Your application has been saved!  Feel free to change it until applications close",
		})
		session.Save(r, w)

		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}
