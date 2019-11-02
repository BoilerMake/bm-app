// Package mailgun provides a Mailer interface to send mail
package mailgun

import (
	"context"
	"strings"
	"time"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/template"

	"github.com/mailgun/mailgun-go/v3"
)

type Mailer struct {
	mg        *mailgun.MailgunImpl
	sender    string
	templates *template.Template
}

func NewMailer(sender string, domain string, apiKey string, tmpls *template.Template) Mailer {
	return Mailer{
		mg:        mailgun.NewMailgun(domain, apiKey),
		sender:    sender,
		templates: tmpls,
	}
}

func (m Mailer) Send(to string, subject string, body string) error {
	message := m.mg.NewMessage(m.sender, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err := m.mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return err
}

func (m Mailer) SendTemplate(to string, subject string, tmpl string, data interface{}) error {
	body, err := m.templates.RenderToString(tmpl, data)
	if err != nil {
		return err
	}

	message := m.mg.NewMessage(m.sender, subject, body, to)

	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err = m.mg.Send(ctx, message)
	if err != nil {
		// Mailgun's kinda annoying how it handles errors, so we'll just convert it to a string and
		// see if it's the kind we're worried about.  We do a little bit of client side email
		// validation (just what the HTML element with type="email" does) but it seems like
		// mailgun is more strict.
		if strings.Contains(err.Error(), "'to' parameter is not a valid address. please check documentation") {
			return models.ErrInvalidEmail
		}

		return err
	}

	return err
}
