package web

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
)

// getApply renders the apply template.
func (h *Handler) getApply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

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
		} else {
			// Show flash that app has been saved
			session.AddFlash(flash.Flash{
				Type:    flash.Success,
				Message: "You're application has been submitted! feel free to update it here until applications close.",
			})
			session.Save(r, w)
		}

		p, ok := h.NewPage(w, r, "BoilerMake - Apply")
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
	return func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		var a models.Application

		err := a.FromFormData(r)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		session := h.getSession(r)

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

		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}
