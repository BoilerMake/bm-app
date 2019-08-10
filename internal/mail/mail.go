// Package mail defines an interface through which mail can be sent.
package mail

import (
	"log"
	"os"
	"path"

	"github.com/BoilerMake/new-backend/internal/mail/mailgun"
	"github.com/BoilerMake/new-backend/internal/mail/mock"
	"github.com/BoilerMake/new-backend/pkg/template"
)

// A Mailer defines an interface for sending mail.
type Mailer interface {
	Send(to string, subject string, body string) (err error)
	SendTemplate(to string, subject string, tmpl string, data interface{}) (err error)
}

// NewMailer creates a new mailer based on the environment mode given to it. In
// development mode it will return a mock mailer and in every other mode
// (really just production) it will return a mailgun mailer.
func NewMailer() (m Mailer) {
	tmplPath, ok := os.LookupEnv("TEMPLATES_PATH")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "TEMPLATES_PATH")
	}
	tmplPath = path.Join(tmplPath, "email")
	tmpls, err := template.NewTemplate(tmplPath, nil)
	if err != nil {
		log.Fatalf("failed to load templates: %s", err)
	}

	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "ENV_MODE")
	}

	if mode == "development" {
		m = mock.NewMailer(tmpls)
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

		m = mailgun.NewMailer(sender, domain, apiKey, tmpls)
	}

	return m
}
