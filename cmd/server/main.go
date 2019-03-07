package main

import (
	"log"
	"os"

	"github.com/BoilerMake/new-backend/internal/http"
	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/internal/postgres"
	"github.com/BoilerMake/new-backend/pkg/env"
)

func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO I *really* don't think this is the ideal way to configure things, plz help ;(
	// Maybe we could have a config package that exports all the configs we need
	// and it has an init that loads env and sets those exported fields?
	// JWT configuring
	signingString, ok := os.LookupEnv("JWT_SIGNING_KEY")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_SIGNING_KEY")
	}
	models.JWTSigningKey = []byte(signingString)

	issuer, ok := os.LookupEnv("JWT_ISSUER")
	if !ok {
		log.Fatalf("environment variable not set: %v", "JWT_ISSUER")
	}
	models.JWTIssuer = issuer

	// Initialize databse
	connStr, ok := os.LookupEnv("DB_CONN")
	if !ok {
		log.Fatalf("environment variable not set: %v", "HOST")
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
