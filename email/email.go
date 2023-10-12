package email

import (
	"context"
	"io"
	"log"

	"github.com/wneessen/go-mail"
)

type (
	SMTPEmail struct {
		client *mail.Client

		// SMTP
		smtpUser    string
		smtpPass    string
		smtpFrom    string
		optionalTLS bool
	}

	ConsoleEmail struct{}
)

func NewSMTP(host string, port int, opts ...func(e *SMTPEmail)) (*SMTPEmail, error) {
	email := &SMTPEmail{}
	for _, opt := range opts {
		opt(email)
	}

	var err error
	tls := mail.TLSMandatory
	if email.optionalTLS {
		tls = mail.TLSOpportunistic
	}
	email.client, err = mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPolicy(tls),
		mail.WithUsername(email.smtpUser),
		mail.WithPassword(email.smtpPass),
	)
	if err != nil {
		return nil, err
	}

	return email, nil
}

func (e *SMTPEmail) Send(ctx context.Context, to, subject string, textReader, htmlReader io.Reader) error {
	msg := mail.NewMsg()

	msg.From(e.smtpFrom)
	msg.To(to)
	msg.Subject(subject)

	if raw, err := io.ReadAll(textReader); err == nil {
		msg.SetBodyString(mail.TypeTextPlain, string(raw))
	}

	if raw, err := io.ReadAll(htmlReader); err == nil {
		msg.SetBodyString(mail.TypeTextHTML, string(raw))
	}

	return e.client.DialAndSendWithContext(ctx, msg)
}

func WithPlainAuth(user, pass string) func(e *SMTPEmail) {
	return func(e *SMTPEmail) {
		e.smtpUser = user
		e.smtpPass = pass
	}
}

func WithFrom(from string) func(e *SMTPEmail) {
	return func(e *SMTPEmail) {
		e.smtpFrom = from
	}
}

func WithOptionalTLS(optional bool) func(e *SMTPEmail) {
	return func(e *SMTPEmail) {
		e.optionalTLS = optional
	}
}

func (c ConsoleEmail) Send(ctx context.Context, to, subject string, textReader, htmlReader io.Reader) error {
	if raw, err := io.ReadAll(textReader); err == nil {
		log.Printf(`Sending email to="%s" subject="%s" body="%s"`, to, subject, string(raw))
	} else {
		log.Printf(`Sending email to="%s" subject="%s"`, to, subject)
	}
	return nil
}
