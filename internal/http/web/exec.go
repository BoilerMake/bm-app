package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type WebData struct {
	count string
}

// getExec renders the exec home page.
func (h *Handler) getExec() http.HandlerFunc {
	applicationCount := h.ApplicationService.GetApplicationCount()
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Exec")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		applicationCountString := strconv.Itoa(applicationCount)
		println(applicationCountString)
		wd := WebData{
			count: applicationCountString,
		}
		fmt.Print(wd)
		h.Templates.RenderTemplate(w, "exec", p)

	}

}
