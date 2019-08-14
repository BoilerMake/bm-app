package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/BoilerMake/new-backend/internal/http"
	"github.com/BoilerMake/new-backend/internal/mail"
	"github.com/BoilerMake/new-backend/internal/postgres"
	"github.com/BoilerMake/new-backend/internal/s3"
	"github.com/BoilerMake/new-backend/pkg/env"
	"github.com/rollbar/rollbar-go"
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

	// Setup rollbar
	rollbar.SetToken("cb55bdedcb2042f7a6a3687e1c266e08")
	rollbar.SetEnvironment("production")                 // defaults to "development"
	rollbar.SetCodeVersion("v2")                         // optional Git hash/branch/tag (required for GitHub integration)
	rollbar.SetServerHost("web.1")                       // optional override; defaults to hostname
	rollbar.SetServerRoot("github.com/heroku/myproject") // path of project (required for GitHub integration and non-project stacktrace collapsing)

	rollbar.Critical(errors.New("Test error"))
	rollbar.Info("Message body goes here")
	rollbar.Wait()

	us := &postgres.UserService{db}
	as := &postgres.ApplicationService{db}
	mailer := mail.NewMailer()
	S3 := s3.NewS3()
	h := http.NewHandler(us, as, mailer, S3)

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
