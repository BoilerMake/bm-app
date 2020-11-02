package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BoilerMake/bm-app/internal/status"
)

// OnSeasonOnly makes sure the APP_STATUS environment variable is between 2 and 4, if
// it's outside of that then the handler returns 404 instead
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

// AppsOpenOnly makes sure the APP_STATUS environment variable is either 2 or 4, if
// it's outside of that then
func AppsOpenOnly(h http.Handler) http.Handler {
	statusString := mustGetEnv("APP_STATUS")
	statusInt, err := strconv.Atoi(statusString)
	if err != nil {
		log.Fatalf("Failed to convert status to int: %v", err)
	}
	st := status.Status(statusInt)

	fn := func(w http.ResponseWriter, r *http.Request) {
		if st == status.ApplicationsOpen || st == status.Live {
			h.ServeHTTP(w, r)
			return
		} else {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}
	}

	return http.HandlerFunc(fn)
}

// LiveOnly makes sure that the APP_STATUS environment variable is only 4, if
// it's outside of that then return 404
func LiveOnly(h http.Handler) http.Handler {
	statusString := mustGetEnv("APP_STATUS")
	statusInt, err := strconv.Atoi(statusString)
	if err != nil {
		log.Fatalf("Faied to convert status to int: %v", err)
	}
	st := status.Status(statusInt)

	fn := func(w http.ResponseWriter, r*http.Request) {
		if st == status.Live {
			h.ServeHTTP(w, r)
			return
		} else {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}
	}

	return http.HandlerFunc(fn)
}