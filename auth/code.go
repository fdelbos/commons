package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	tmplHTML "html/template"
	"strings"
	tmplText "text/template"
	"time"

	"github.com/dchest/uniuri"
)

type (
	Code struct {
		Digest []byte
		Until  time.Time
	}
)

const (
	Digits                    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultDigitsTextTemplate = `Your code is {{.Code}}`
	DefaultDigitsEmailSubject = "Your code"
	SaltLenght                = 16
	DefaultCodeValidity       = 5 * time.Minute
)

var (
	ErrInvalidCode = errors.New("invalid code")
)

type (
	// CodeStore is the interface to store and retrieve codes (ie: the database).
	CodeStore interface {
		NewCode(ctx context.Context, code *Code) error                 // store a new unique code for the email until the given date
		GetCode(ctx context.Context, codeDigest []byte) (*Code, error) // return the code or an error if the digest is not found
		Use(ctx context.Context, codeDigest []byte) error              // mark the code as used, it should never be returned by GetCode again
	}

	// Code is the service to send and validate authenticaition codes by email.
	Codes struct {
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

// NewCodes creates a new code authentication service.
// The code is sent by email.
func NewCodes(mailer Mailer, codeStore CodeStore, opts ...func(*Codes)) (*Codes, error) {
	code := &Codes{
		mailer:   mailer,
		store:    codeStore,
		nbDigits: 8,
		validity: DefaultCodeValidity,
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
func (c *Codes) Send(ctx context.Context, to string) error {
	digits, code := c.NewCode(to)
	data := CodeTemplateData{
		Code:  digits,
		Until: code.Until,
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

	if err := c.store.NewCode(ctx, code); err != nil {
		return err
	}

	err := c.mailer.Send(
		ctx,
		to,
		c.emailSubject,
		textBuff,
		htmlBuff)
	return err
}

// Validate checks if the given code is valid for the given email.
func (c *Codes) Validate(ctx context.Context, digits, email string) error {
	digest := GenDigest(email, digits)
	code, err := c.store.GetCode(ctx, digest)
	if err != nil {
		return err
	}
	if code.Until.Before(time.Now()) {
		return ErrInvalidCode
	}
	if err := c.store.Use(ctx, digest); err != nil {
		return err
	}

	return nil
}

// WithCodeNbDigits sets the number of digits of the code. Default is 8.
func WithCodeNbDigits(nbDigits int) func(*Codes) {
	return func(c *Codes) {
		c.nbDigits = nbDigits
	}
}

// WithCodeValidity sets the validity duration of the code. Default is 5 minutes.
func WithCodeValidity(validity time.Duration) func(*Codes) {
	return func(c *Codes) {
		c.validity = validity
	}
}

// WithCodeTemplates sets the templates used to send the code.
// The text template and subject are mandatory.
func WithCodeTemplates(subject, textTemplate, htmlTemplate string) func(*Codes) {
	return func(c *Codes) {
		c.emailSubject = subject
		c.textTemplate = textTemplate
		c.HTMLTemplate = htmlTemplate
	}
}

func GenDigest(email, digits string) []byte {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	digits = strings.TrimSpace(digits)
	digits = strings.ToUpper(digits)

	sha := sha256.New()
	sha.Write([]byte(email))
	sha.Write([]byte(digits))
	return sha.Sum(nil)
}

func (c *Codes) NewCode(email string) (string, *Code) {
	code := &Code{
		Until: time.Now().Add(c.validity),
	}
	digits := uniuri.NewLenChars(c.nbDigits, []byte(Digits))
	code.Digest = GenDigest(email, digits)

	return digits, code
}
