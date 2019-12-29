package web

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
)

// getRSVP renders the RSVP template.
func (h *Handler) getRSVP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

		id, ok := session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		// First make sure they've submitted an application
		app, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				session.AddFlash(flash.Flash{
					Type:    flash.Error,
					Message: "Please submit an application first.",
				})
				session.Save(r, w)

				http.Redirect(w, r, "/apply", http.StatusSeeOther)
				return
			} else {
				h.Error(w, r, err, "")
				return
			}
		}

		// Now make sure RSVP has not expired
		if time.Now().Sub(app.AcceptedAt) > models.RSVPExpiryTime {
			session.AddFlash(flash.Flash{
				Type:    flash.Error,
				Message: "Your RSVP has expired.",
			})
			session.Save(r, w)

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		rsvp, err := h.RSVPService.GetByUserID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				rsvp = &models.RSVP{}
			} else {
				h.Error(w, r, err, "")
				return
			}
		} else {
			// Show flash that we have RSVP on file
			session.AddFlash(flash.Flash{
				Type:    flash.Success,
				Message: "Your RSVP has been submitted!",
			})
			session.Save(r, w)
		}

		p, ok := h.NewPage(w, r, "BoilerMake - RSVP")
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		p.FormRefill = rsvp

		h.Templates.RenderTemplate(w, "rsvp", p)
	}
}

// postRSVP renders the RSVP template.
func (h *Handler) postRSVP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		var rsvp models.RSVP

		err := rsvp.FromFormData(r)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		session := h.getSession(r)

		rsvp.UserID, ok = session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		err = h.RSVPService.CreateOrUpdate(&rsvp)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		// Show flash that we got RSVP
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Your RSVP has been submitted!",
		})
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
