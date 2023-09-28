package auth

import (
	"context"
	"io"
)

type (
	Mailer interface {
		Send(ctx context.Context, to, subject string, text, html io.Reader) error
	}
)
