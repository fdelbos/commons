package auth

import (
	"bytes"
	"context"
	"errors"
	tmplHTML "html/template"
	"strings"
	tmplText "text/template"
	"time"

	"github.com/dchest/uniuri"
)

const (
	Digits                    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultDigitsTextTemplate = `Your code is {{.Code}}`
	DefaultDigitsEmailSubject = "Your code"
)

var (
	ErrInvalidCode = errors.New("invalid code")
)

type (
	// CodeStore is the interface to store and retrieve codes (ie: the database).
	CodeStore interface {
		NewCode(ctx context.Context, email, code string, until time.Time) error // store a new unique code for the email until the given date
		GetCode(ctx context.Context, code string) (string, time.Time, error)    // return the email, the until date and an error if the id is not found
		Use(ctx context.Context, code string) error                             // mark the code as used, it should never be returned by GetCode again
	}

	// Code is the service to send and validate authenticaition codes by email.
	Code struct {
		mailer       Mailer
		store        CodeStore
		nbDigits     int
		validity     time.Duration
		textTemplate string
		HTMLTemplate string
		emailSubject string

		tmplText *tmplText.Template
		tmplHTML *tmplHTML.Template
	}

	// CodeTemplateData is the data used to render the code templates.
	CodeTemplateData struct {
		Code  string    // the code
		Until time.Time // the date until the code is valid
	}
)

// NewCode creates a new code authentication service.
// The code is sent by email.
func NewCode(mailer Mailer, codeStore CodeStore, opts ...func(*Code)) (*Code, error) {
	code := &Code{
		mailer:   mailer,
		store:    codeStore,
		nbDigits: 8,
		validity: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(code)
	}

	if code.textTemplate == "" {
		code.textTemplate = DefaultDigitsTextTemplate
	}

	var err error
	code.tmplText, err = tmplText.New("code").Parse(code.textTemplate)
	if err != nil {
		return nil, err
	}

	if code.HTMLTemplate != "" {
		code.tmplHTML, err = tmplHTML.New("code").Parse(code.HTMLTemplate)
		if err != nil {
			return nil, err
		}
	}

	if code.emailSubject == "" {
		code.emailSubject = DefaultDigitsEmailSubject
	}

	return code, nil
}

// Send sends a code to the given email.
func (c *Code) Send(ctx context.Context, to string) error {
	data := CodeTemplateData{
		Code:  uniuri.NewLenChars(c.nbDigits, []byte(Digits)),
		Until: time.Now().Add(c.validity),
	}
	textBuff := &bytes.Buffer{}
	if err := c.tmplText.Execute(textBuff, data); err != nil {
		return err
	}

	htmlBuff := &bytes.Buffer{}
	if c.tmplHTML != nil {
		buff := &bytes.Buffer{}
		if err := c.tmplHTML.Execute(buff, data); err != nil {
			return err
		}
	}

	if err := c.store.NewCode(ctx, to, data.Code, data.Until); err != nil {
		return err
	}

	return c.mailer.Send(
		ctx,
		to,
		c.emailSubject,
		textBuff,
		htmlBuff)
}

// Validate checks if the given code is valid for the given email.
func (c *Code) Validate(ctx context.Context, code, email string) error {
	code = strings.ToUpper(code)
	code = strings.TrimSpace(code)
	if code == "" {
		return ErrInvalidCode
	}

	mail, until, err := c.store.GetCode(ctx, code)
	if err != nil {
		return ErrInvalidCode
	}
	if mail != email {
		return ErrInvalidCode
	}
	if until.Before(time.Now()) {
		return ErrInvalidCode
	}
	if err := c.store.Use(ctx, code); err != nil {
		return err
	}

	return nil
}

// WithCodeNbDigits sets the number of digits of the code. Default is 8.
func WithCodeNbDigits(nbDigits int) func(*Code) {
	return func(c *Code) {
		c.nbDigits = nbDigits
	}
}

// WithCodeValidity sets the validity duration of the code. Default is 5 minutes.
func WithCodeValidity(validity time.Duration) func(*Code) {
	return func(c *Code) {
		c.validity = validity
	}
}

// WithCodeTemplates sets the templates used to send the code.
// The text template and subject are mandatory.
func WithCodeTemplates(subject, textTemplate, htmlTemplate string) func(*Code) {
	return func(c *Code) {
		c.emailSubject = subject
		c.textTemplate = textTemplate
		c.HTMLTemplate = htmlTemplate
	}
}
