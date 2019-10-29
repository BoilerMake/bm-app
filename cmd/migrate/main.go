package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BoilerMake/new-backend/internal/postgres"
	"github.com/BoilerMake/new-backend/pkg/env"

	"github.com/pressly/goose"
)

func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize database
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

	err = db.Ping()
	for err != nil {
		fmt.Println("Ping to database failed, retrying")
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
