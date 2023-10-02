package auth_test

import (
	"context"
	"io"
	"testing"
	"time"

	. "github.com/fdelbos/commons/auth"
	"github.com/fdelbos/commons/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCode(t *testing.T) {
	mailer := mocks.NewAuthMailer(t)
	codeStore := mocks.NewAuthCodeStore(t)

	codes, err := NewCodes(mailer, codeStore, WithCodeTemplates(DefaultDigitsEmailSubject, "{{.Code}}", ""))
	assert.NoError(t, err)

	ctx := context.Background()

	digitsFromMock := ""
	// test send code

	codeStore.
		On("NewCode", ctx, mock.Anything).
		Return(func(ctx context.Context, code *Code) error {
			assert.True(t, code.Until.After(time.Now()))
			return nil
		}).
		Once()

	mailer.
		On("Send", ctx, "test@example.com", DefaultDigitsEmailSubject, mock.Anything, mock.Anything).
		Return(func(ctx context.Context, to, subject string, text, html io.Reader) error {
			assert.Equal(t, "test@example.com", to)
			assert.Equal(t, DefaultDigitsEmailSubject, subject)
			raw, err := io.ReadAll(text)
			assert.NoError(t, err)
			digitsFromMock = string(raw)
			return nil
		}).
		Once()

	err = codes.Send(ctx, "test@example.com")
	assert.NoError(t, err)

	// test validate code
	digest := GenDigest("test@example.com", digitsFromMock)
	until := time.Now().Add(time.Hour)
	codeStore.
		On("GetCode", ctx, digest).
		Return(&Code{Digest: digest, Until: until}, nil).
		Once()

	codeStore.
		On("Use", ctx, digest).
		Return(nil).
		Once()

	err = codes.Validate(ctx, digitsFromMock, "test@example.com")
	assert.NoError(t, err)

	mailer.AssertExpectations(t)
	codeStore.AssertExpectations(t)
}
