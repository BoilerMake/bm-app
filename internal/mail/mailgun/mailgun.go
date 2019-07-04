// Package mailgun provides a Mailer interface to send mail
package mailgun

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

type Mailer struct {
	mg     *mailgun.MailgunImpl
	sender string
}

func NewMailer(sender string, domain string, apiKey string) Mailer {
	return Mailer{
		mg:     mailgun.NewMailgun(domain, apiKey),
		sender: sender,
	}
}

func (m Mailer) Send(to string, subject string, body string) (err error) {
	message := m.mg.NewMessage(m.sender, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err = m.mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return err
}
