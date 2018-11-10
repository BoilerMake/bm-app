package app

import (
	"database/sql"

	"github.com/go-chi/chi"
)

// server keeps our routes and dependency injections. 
// Keeping dependency injection here allows for easier testing using mocks
type server struct {
	router *chi.Mux
	db *sql.DB
}

func Start(c *config.Config) {
	srv := &{
		router: chi.NewRouter()
	}

	// Register routes in app/routes.go
	srv.routes()
}

func Stop() {

}
