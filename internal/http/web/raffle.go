package web

import (
	"errors"
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

		ticketsCount := 5

		p.Data = map[string]interface{}{
			"TicketsCount" : ticketsCount,
		}

		h.Templates.RenderTemplate(w, "raffle", p)

	}

}
