package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/rollbar/rollbar-go"
)

// Recoverer recovers from a panic that a was hit while handling a request.  Ideally,
// it should never be run because ideally nothing will panic.  In the case that does
// happen, this will prevent the entire server from shutting down.  It will also
// log the error (and report it to rollbar), and render the status code 500 page.
func Recoverer(next http.Handler) http.Handler {
	rollbarEnv := mustGetEnv("ROLLBAR_ENVIRONMENT")

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				if rollbarEnv == "production" {
					rollbar.Error("PANIC:", rvr)
					rollbar.Wait()
				}

				log.Printf("PANIC: %+v\n", rvr)
				debug.PrintStack()

				http.Redirect(w, r, "/500", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
