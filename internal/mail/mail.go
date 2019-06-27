// Package mail defines an interface through which mail can be sent.
package mail

import (
	"log"
	"os"

	"github.com/BoilerMake/new-backend/internal/mail/mailgun"
	"github.com/BoilerMake/new-backend/internal/mail/mock"
)

// A Mailer defines an interface for sending mail.
type Mailer interface {
	Send(to string, subject string, body string) (err error)
}

// NewMailer creates a new mailer based on the environment mode given to it. In
// development mode it will return a mock mailer and in every other mode
// (really just production) it will return a mailgun mailer.
func NewMailer() (m Mailer) {
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "DB_PASSWORD")
	}

	if mode == "development" {
		m = mock.NewMailer()
	} else {
		sender, ok := os.LookupEnv("MAILGUN_ADDRESS")
		if !ok {
			log.Fatalf("environment variable not set: %v. Did you update your .env file?", "MAILGUN_ADDRESS")
		}
		domain, ok := os.LookupEnv("MAILGUN_DOMAIN")
		if !ok {
			log.Fatalf("environment variable not set: %v. Did you update your .env file?", "MAILGUN_DOMAIN")
		}
		apiKey, ok := os.LookupEnv("MAILGUN_API_KEY")
		if !ok {
			log.Fatalf("environment variable not set: %v. Did you update your .env file?", "MAILGUN_API_KEY")
		}

		m = mailgun.NewMailer(sender, domain, apiKey)
	}

	return m
}
