package email

import (
	"context"
	"io"

	"github.com/wneessen/go-mail"
)

type (
	Email struct {
		client *mail.Client

		// SMTP
		smtpUser    string
		smtpPass    string
		smtpFrom    string
		optionalTLS bool
	}
)

func NewSMTP(host string, port int, opts ...func(e *Email)) (*Email, error) {
	email := &Email{}
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

func (e *Email) Send(ctx context.Context, to, subject string, textReader, htmlReader io.Reader) error {
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

func WithPlainAuth(user, pass string) func(e *Email) {
	return func(e *Email) {
		e.smtpUser = user
		e.smtpPass = pass
	}
}

func WithFrom(from string) func(e *Email) {
	return func(e *Email) {
		e.smtpFrom = from
	}
}

func WithOptionalTLS(optional bool) func(e *Email) {
	return func(e *Email) {
		e.optionalTLS = optional
	}
}
