package main

import (
	"fmt"
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
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "DB_HOST")
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "DB_NAME")
	}
	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "DB_USER")
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "DB_PASSWORD")
	}
	dbOptions, _ := os.LookupEnv("DB_OPTIONS")

	// Bring together all our config bits and try to connect
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s %s", dbHost, dbName, dbUser, dbPassword, dbOptions)
	db, err := postgres.Open(connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	us := &postgres.UserService{db}
	h := http.NewHandler(us)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalf("environment variable not set: %v", "PORT")
	}

	addr := ":" + port
	server := http.NewServer(addr, h)

	err = server.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
