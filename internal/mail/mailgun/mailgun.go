// Package mailgun provides a Mailer interface to send mail
package mailgun

import (
	"context"
	"time"

	"github.com/BoilerMake/new-backend/pkg/template"

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
		return err
	}

	return err
}
