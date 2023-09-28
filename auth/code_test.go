package auth_test

import (
	"context"
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

	code, err := NewCode(mailer, codeStore)
	assert.NoError(t, err)

	ctx := context.Background()

	// test send code
	mailer.
		On("Send", ctx, "test@example.com", DefaultDigitsEmailSubject, mock.Anything, mock.Anything).
		Return(nil).
		Once()

	codeStore.
		On("NewCode", ctx, "test@example.com", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	err = code.Send(ctx, "test@example.com")
	assert.NoError(t, err)

	// test validate code
	codeStore.
		On("GetCode", ctx, "ABCDEFGH").
		Return("test@example.com", time.Now().Add(time.Minute), nil).
		Once()

	codeStore.
		On("Use", ctx, "ABCDEFGH").
		Return(nil).
		Once()

	err = code.Validate(ctx, "aBcDEFgh", "test@example.com")
	assert.NoError(t, err)

	mailer.AssertExpectations(t)
	codeStore.AssertExpectations(t)
}
