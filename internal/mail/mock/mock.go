// Package mock provides a Mailer interface to test sending mail.
package mock

import (
	"log"

	"github.com/BoilerMake/new-backend/pkg/template"
)

type Mailer struct {
	templates *template.Template
}

func NewMailer(tmpls *template.Template) Mailer {
	return Mailer{
		templates: tmpls,
	}
}

func (m Mailer) Send(to string, subject string, body string) error {
	log.Printf("\n******MOCK MAILER******\nTo: %v\nSubject: %v\nMessage: %s", to, subject, body)
	return nil
}

func (m Mailer) SendTemplate(to string, subject string, tmpl string, data interface{}) error {
	body, err := m.templates.RenderToString(tmpl, data)
	if err != nil {
		return err
	}

	log.Printf("\n******HTML MOCK MAILER******\nTo: %v\nSubject: %v\nMessage: %s", to, subject, body)
	return nil
}
