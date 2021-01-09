package web

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
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

		// ensure only participants can reach the /raffle page
		user, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			if err == sql.ErrNoRows { // user has not submitted an application
				h.Error(w, r, models.ErrRaffleAccessDenied, "/dashboard")
				return
			}
			h.Error(w, r, err, "")
			return
		}
		if user.Decision != 3 {
			h.Error(w, r, models.ErrRaffleAccessDenied, "/dashboard") // user has not been accepted
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

		// ensure only participants can claim a raffle
		user, err := h.ApplicationService.GetByUserID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				h.Error(w, r, models.ErrRaffleAccessDenied, "/dashboard")
				return
			}
			h.Error(w, r, err, "")
			return
		}
		if user.Decision != 3 {
			h.Error(w, r, models.ErrRaffleAccessDenied, "/dashboard")
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

		// check time stamps
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
		ra.FromFormData(r) // err is always nil

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

// addTickets is endpoint for exec adding points to a specific hacker
func (h *Handler) addTickets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)
		// get email address and points to add
		email := r.FormValue("email")
		points := r.FormValue("points")
		pointsInt, err := strconv.Atoi(points)
		if err != nil {
			h.Error(w, r, models.ErrInvalidPointsToAdd, "/exec")
			return
		}

		// get user by email
		user, err := h.UserService.GetByEmail(email)
		if err != nil {
			if err == sql.ErrNoRows {
				h.Error(w, r, models.ErrEmailNotFound, "/exec")
				return
			}
			h.Error(w, r, err, "/exec")
			return
		}

		id := user.ID
		err = h.ApplicationService.AddPointsToUser(id, pointsInt)
		if err != nil {
			h.Error(w, r, err, "/exec")
			return
		}

		// flash success message
		successMessage := fmt.Sprintf("%d points have been successfully added to %s", pointsInt, email)
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: successMessage,
		})
		session.Save(r, w)

		// redirect back to exec page
		http.Redirect(w, r, "/exec", http.StatusSeeOther)
	}
}
