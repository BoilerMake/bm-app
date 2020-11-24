package web

import (
	"fmt"
	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
	"net/http"
)

// createRaffle is endpoint for creating raffle from exec
func (h *Handler) createRaffle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create a raffle struct
		var ra models.Raffle

		session := h.getSession(r)

		// Populate raffle
		_ = ra.FromFormData(r) // err is always nil

		// Validate raffle
		err := ra.Validate()
		if err != nil {
			h.Error(w, r, err, "/exec")
			return
		}

		// send to db
		err = h.RaffleService.Create(&ra)
		if err != nil {
			h.Error(w, r, err, "/exec")
			return
		}

		// flash success messages

		successMessage := fmt.Sprintf("Raffle: %s was successfully created", ra.Code)
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: successMessage,
		})
		session.Save(r, w)

		// redirect back to exec page
		http.Redirect(w, r, "/exec", http.StatusSeeOther)
	}
}
