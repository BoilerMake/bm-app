package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
	"github.com/go-chi/chi"
)

type idMessage struct {
	ID int `json:"id"`
}

type slackMessage struct {
	Text string `json:"text"`
}

func (h *Handler) postAnnouncement() http.HandlerFunc {
	webhookURL := mustGetEnv("SLACK_ANNOUNCEMENTS_WEBHOOK")
	mode := mustGetEnv("ENV_MODE")

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

		if mode != "development" {
			slackMsg, err := json.Marshal(slackMessage{Text: message})
			req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackMsg))
			if err != nil {
				h.Error(w, r, err, "")
				return
			}
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				h.Error(w, r, err, "")
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			if buf.String() != "ok" {
				h.Error(w, r, err, "")
				return
			}
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
