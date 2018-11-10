package app

import (
	"http"
)

// This is just an example handler
// Probably don't commit it

func (s *server) handleSomething() http.HandlerFunc {
	// Do whatever inits you need to use in the handler here
	return func(w http.ResponseWriter, r *http.Request) {
		...
	}
}
