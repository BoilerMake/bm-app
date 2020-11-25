package web

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
	"net/http"
	"strconv"
	"time"
)

// getRaffle generates the page for users to enter a raffle
func (h *Handler) getRaffle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

		id, ok := session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		p, ok := h.NewPage(w, r, "BoilerMake - Raffle")
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		user, err := h.ApplicationService.GetByUserID(id)
		if err != nil { // TODO: Make check that user is checked-in here
			h.Error(w, r, err, "")
			return
		}
		userPoints := user.Points
		p.Data = map[string]interface{}{
			"TicketsCount": userPoints,
		}

		h.Templates.RenderTemplate(w, "raffle", p)
	}
}

// postRaffle tries to claim a raffle for a user
func (h *Handler) postRaffle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

		id, ok := session.Values["ID"].(int)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		// get code
		var code string
		code = r.FormValue("raffle") // this will be from the hacker page
		if code == "" {
			h.Error(w, r, models.ErrRaffleEmpty, "")
			return
		}

		// get the raffle
		raffle, err := h.RaffleService.GetById(code)
		if err != nil {
			if err == sql.ErrNoRows {
				h.Error(w, r, models.ErrInvalidRaffle, "/raffle")
				return
			} else {
				h.Error(w, r, err, "")
				return
			}
		}

		// check time stamps0
		start, err := strconv.ParseInt(raffle.StartTime, 10, 64)
		if err != nil {
			h.Error(w, r, err, "") // should never reach here
		}
		end, err := strconv.ParseInt(raffle.EndTime, 10, 64)
		if err != nil {
			h.Error(w, r, err, "") // should never reach here
		}
		now := time.Now().UnixNano() / 1000000 // convert current time to milliseconds since epoch
		if now < start || now >= end {
			h.Error(w, r, models.ErrTime, "/raffle")
			return
		}
		// check if user is checked in // maybe move this check to the get

		// claim raffle
		err = h.RaffleService.ClaimRaffle(id, code)
		if err != nil {
			h.Error(w, r, err, "/raffle")
			return
		}

		// add points to user
		points, err := strconv.Atoi(raffle.Points)
		if err != nil {
			h.Error(w, r, err, "") // should never reach here
			return
		}
		err = h.ApplicationService.AddPointsToUser(id, points)
		if err != nil {
			h.Error(w, r, err, "/raffle")
			return
		}

		// redirect to /raffle with success message
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Raffle was successfully claimed!",
		})
		session.Save(r, w)

		http.Redirect(w, r, "/raffle", http.StatusSeeOther)
	}
}

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
