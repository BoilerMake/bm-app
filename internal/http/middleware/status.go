package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BoilerMake/new-backend/internal/status"
)

// OnSeasonOnly make sure the APP_STATUS environment variable is between 2 and 4, if
// it's outside of that then
func OnSeasonOnly(h http.Handler) http.Handler {
	statusString := mustGetEnv("APP_STATUS")
	statusInt, err := strconv.Atoi(statusString)
	if err != nil {
		log.Fatalf("Failed to convert status to int: %v", err)
	}
	st := status.Status(statusInt)

	fn := func(w http.ResponseWriter, r *http.Request) {
		if st < status.ApplicationsOpen || st > status.Live {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
