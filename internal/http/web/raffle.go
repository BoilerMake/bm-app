package web

import (
	"errors"
	"fmt"
	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"
	"net/http"
)

// getRaffle renders the raffle page for hackers.
func (h *Handler) getRaffle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Raffle")
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		ticketsCount := 5 // make a call to db to retrieve points

		p.Data = map[string]interface{}{
			"TicketsCount" : ticketsCount,
		}

		h.Templates.RenderTemplate(w, "raffle", p)

	}

}

//// postClaimRaffle tries to claim a raffle for user
//func (h *Handler) postClaimRaffle() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// get code
//		var code string
//		code = r.FormValue("raffle")
//		if code == "" {
//			h.Error(w,r, models.ErrRaffleEmpty, "") // post call on /raffle url
//			return
//		}
//
//		// validate with db by getting the raffle. if it doesn't exist see error
//		ra, err := h.RaffleService.GetByCode(code)
//		if err != nil {
//			h.Error(w, r, err,"")
//			return
//		}
//		// if it does exist, check time stamps
//		start := ra.StartTime
//		end := ra.EndTime
//		// get current time and check if it's within timestamp
//		// try to add to user, check if user checked in then check if already claimed this raffle
//		// one service call. The service call can make other service calls
//		// Look at Update/Set/Where command in CreateOrUpdate for application service
//		// where call may need to check for appreance in hstore first
//		// pull map from hstore and check. add map back in and update map and points
//
//		// flash a success or report an error
//
//		// redirect to /raffle -> this should regenerate the tickets count
//	}
//}

// createRaffle attempts to create a new raffle
func (h *Handler) createRaffle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create a raffle struct
		var ra models.Raffle

		// validate
		err := ra.FromFormData(r) // maybe change it so the db stores them as strings
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			h.Error(w, r, err, "/exec")
			return
		}

		// send to db
		err = h.RaffleService.Create(&ra)
		if err != nil {
			h.Error(w, r, err, "/exec")
		}
		// flash success messages
		session := h.getSession(r)
		successMessage := fmt.Sprintf("Raffle: %s was successfully created", ra.Code)
		session.AddFlash(flash.Flash{
			Type:		flash.Success,
			Message:	successMessage,
		})
		session.Save(r, w)

		// redirect back to exec page
		http.Redirect(w, r, "/exec", http.StatusSeeOther)
	}
}

//// addPointsToHacker attempts to add raffle points to a specific hacker
//func (h * Handler) addPointsToHacker() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// parse data to get hacker email and points to add
//		// call GetByEmail in application_service.go (needs to be added): returns models.Application
//		// add points
//		// run CreateOrUpdate
//
//	}
//}