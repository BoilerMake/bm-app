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
			if app.Decision == 0 {
				// Show flash that app has been saved
				session.AddFlash(flash.Flash{
					Type:    flash.Success,
					Message: "Your application has already been submitted! Feel free to update it here until applications close.",
				})
			} else {
				session.AddFlash(flash.Flash{
					Type:    flash.Success,
					Message: "Your application has already been submitted!",
				})
			}
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
		var a models.Application

		session := h.getSession(r)

		id, ok := session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		app, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				app = &models.Application{}
			} else {
				h.Error(w, r, err, "")
				return
			}
		}

		// Don't let users update app after decision has been made
		if app.Decision != 0 {
			session.AddFlash(flash.Flash{
				Type:    flash.Error,
				Message: "You can no longer update your application.",
			})
			session.Save(r, w)

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		err = a.FromFormData(r)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		a.UserID = id

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

		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Your application has been submitted! Check here later to see your decision.",
		})
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
