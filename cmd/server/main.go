package main

import (
	"log"
	"os"

	"github.com/BoilerMake/new-backend/internal/http"
	"github.com/BoilerMake/new-backend/internal/postgres"
	"github.com/BoilerMake/new-backend/pkg/env"
)

func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize databse
	connStr, ok := os.LookupEnv("DB_CONN")
	if !ok {
		log.Fatalf("environment variable not set: %v", "DB_CONN")
	}
	db, err := postgres.Open(connStr)
	defer db.Close()

	us := &postgres.UserService{db}
	h := http.NewHandler(us)

	// Server configuring
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatalf("environment variable not set: %v", "HOST")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalf("environment variable not set: %v", "PORT")
	}

	addr := host + ":" + port
	server := http.NewServer(addr, h)

	err = server.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
