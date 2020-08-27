package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BoilerMake/bm-app/internal/http"
	"github.com/BoilerMake/bm-app/internal/mail"
	"github.com/BoilerMake/bm-app/internal/postgres"
	"github.com/BoilerMake/bm-app/internal/s3"
	"github.com/BoilerMake/bm-app/pkg/env"
	"github.com/BoilerMake/bm-app/pkg/flash"

	"github.com/rollbar/rollbar-go"
)

func main() {
	err := env.Load(true)
	if err != nil {
		log.Fatalln(err)
	}

	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "ENV_MODE")
	}

	// Set up logging with a new file each time the server starts
	if mode == "production" {
		logPath, ok := os.LookupEnv("LOG_PATH")
		if !ok {
			log.Fatalf("environment variable not set: %v", "LOG_PATH")
		}

		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		logFile, err := os.Create(logPath + "/" + timestamp)
		if err != nil {
			log.Fatalf("failed to create log file: %v", err)
		}

		log.SetOutput(logFile)
	}

	// Register flash struct so it can be serialized later
	gob.Register(flash.Flash{})

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
	as := &postgres.ApplicationService{db}
	rs := &postgres.RSVPService{db}
	anns := &postgres.AnnouncementService{db}
	ras := &postgres.RaffleService{db}
	mailer := mail.NewMailer()
	S3 := s3.NewS3()
	h := http.NewHandler(us, as, rs, anns, ras, mailer, S3)

	rollbarEnv, ok := os.LookupEnv("ROLLBAR_ENVIRONMENT")
	if !ok {
		log.Fatalf("environment variable not set: %v", "ROLLBAR_ENVIRONMENT")
	}

	if rollbarEnv == "production" {
		rollbarToken, ok := os.LookupEnv("ROLLBAR_TOKEN")
		if !ok {
			log.Fatalf("environment variable not set: %v", "ROLLBAR_TOKEN")
		}
		rollbarRoot, ok := os.LookupEnv("ROLLBAR_SERVER_ROOT")
		if !ok {
			log.Fatalf("environment variable not set: %v", "ROLLBAR_SERVER_ROOT")
		}

		rollbar.SetToken(rollbarToken)
		rollbar.SetEnvironment(rollbarEnv)
		rollbar.SetCodeVersion("master")
		rollbar.SetServerHost("web.1")
		rollbar.SetServerRoot(rollbarRoot)
	}

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
