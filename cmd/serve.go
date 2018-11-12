package cmd

import (
	"github.com/cjenright/new-backend/app"
)

// Serve creates and starts a new server.
func Serve() {
	srv := app.NewServer()
	srv.Start()
}
