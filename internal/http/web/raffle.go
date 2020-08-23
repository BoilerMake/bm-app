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

//// postRaffle tries to claim a raffle for user
//func (h *Handler) postRaffle() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var code string
//		code = r.FormValue()
//		// get code
//
//		// validate with db by getting the raffle. if it doesn't exist see error
//		// if it does exist, check time stamps
//
//		// try to add to user, check if user checked in then check if already claimed this raffle
//
//		// flash a success or report an error
//
//		// redirect to /raffle -> this should regenerate hte tickets count
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
			w.WriteHeader(http.StatusBadRequest)
			h.Error(w, r, err, "/exec")
			return
		}
		fmt.Println(&ra) // this prints fine
		fmt.Printf("%p\n", &ra) // this prints fine as well

		// send to db
		err = h.RaffleService.Create(&ra) // changing all the fields to strings but dont think won't fix
		if err != nil {
			fmt.Println(err)
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