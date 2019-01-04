package main

import (
	"log"
	"os"

	"github.com/BoilerMake/new-backend/internal/http"
	"github.com/BoilerMake/new-backend/internal/http/api"
	"github.com/BoilerMake/new-backend/internal/http/web"
	"github.com/BoilerMake/new-backend/internal/postgres"
	"github.com/BoilerMake/new-backend/pkg/env"
)

func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	connStr, ok := os.LookupEnv("DB_CONN")
	if !ok {
		log.Fatalf("environment variable not set: %v", "HOST")
	}

	db, err := postgres.Open(connStr)
	defer db.Close()
	us := &postgres.UserService{db}

	apiHandler := api.NewHandler()
	webHandler := web.NewHandler()

	apiHandler.UserService = us
	webHandler.UserService = us
	// TODO should this just be NewHandler(us)?
	h := &http.Handler{apiHandler, webHandler}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatalf("environment variable not set: %v", "HOST")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalf("environment variable not set: %v", "PORT")
	}

	addr := host + ":" + port
	srv := http.NewServer(addr, h)

	err = srv.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
