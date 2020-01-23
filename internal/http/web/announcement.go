package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
	"github.com/go-chi/chi"
)

type idMessage struct {
	ID int `json:"id"`
}

func (h *Handler) postAnnouncement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		message = r.FormValue("message")
		// Check if message is empty or not
		if message == "" {
			w.WriteHeader(http.StatusBadRequest)
			h.Error(w, r, models.ErrAnnouncementMessageEmpty, "")
			return
		}

		err := h.AnnouncementService.Create(message)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		session := h.getSession(r)
		// Show flash that everything went well
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Announcement posted, you're doing great!",
		})
		session.Save(r, w)

		// Redirect to exec page
		http.Redirect(w, r, "/exec", http.StatusSeeOther)
	}
}

func (h *Handler) getAnnouncement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentAnnouncement, err := h.AnnouncementService.GetCurrent()
		if err == models.ErrNoAnnouncements {
			w.WriteHeader(http.StatusNoContent)
			h.Error(w, r, err, "")
			return
		}

		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentAnnouncement)
	}
}

func (h *Handler) getAnnouncementWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			h.Error(w, r, err, "")
			return
		}

		ann, err := h.AnnouncementService.GetByID(id)
		if err == models.ErrNoAnnouncements {
			w.WriteHeader(http.StatusNoContent)
			h.Error(w, r, err, "")
			return
		}

		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ann)
	}
}

func (h *Handler) deleteAnnouncement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m idMessage
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		// Check if id is empty
		if m.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			h.Error(w, r, models.ErrAnnouncementIDEmpty, "")
			return
		}

		err = h.AnnouncementService.DeleteByID(m.ID)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
