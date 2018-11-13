package main

import (
	"log"

	"github.com/cjenright/new-backend/internal/app"
	"github.com/cjenright/new-backend/pkg/env"
)

// Serve creates and starts a new server.
func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	srv := app.NewServer()
	err = srv.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
