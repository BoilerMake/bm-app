// Package mock provides a Mailer interface to test sending mail
package mock

import "log"

type Mailer struct{}

func NewMailer() Mailer {
	return Mailer{}
}

func (m Mailer) Send(to string, subject string, body string) error {
	log.Printf("\n******MOCK MAILER******\nTo: %v\nSubject: %v\nMessage: %s", to, subject, body)
	return nil
}
