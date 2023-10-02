package www_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/fdelbos/commons/internal/mocks"
	. "github.com/fdelbos/commons/www"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCodeSession(t *testing.T) {
	app := fiber.New()

	codesService := mocks.NewWWWCodesService(t)

	expectedEmail := "expected@test.com"
	expectedUUID := uuid.New()
	expectedCode := "123456"

	emailToUUID := func(ctx context.Context, email string) (uuid.UUID, error) {
		assert.Equal(t, expectedEmail, email)
		return expectedUUID, nil
	}

	sessionsService := mocks.NewWWWSessionsService(t)

	NewCodeSession(codesService, emailToUUID, sessionsService).
		Routes(app)

	t.Run("new codes", func(t *testing.T) {
		codesService.
			On("Send", mock.Anything, expectedEmail).
			Return(nil).
			Once()

		body := bytes.NewBufferString(`{"email":"` + expectedEmail + `"}`)
		req := httptest.NewRequest("POST", "/send", body)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("validate codes", func(t *testing.T) {
		codesService.
			On("Validate", mock.Anything, expectedCode, expectedEmail).
			Return(nil).
			Once()

		sessionsService.
			On("NewSession", mock.Anything, expectedUUID, DefaultCodeDuration).
			Return("session_id", nil).
			Once()

		body := bytes.NewBufferString(`{"email":"` + expectedEmail + `", "code":"` + expectedCode + `"}`)
		req := httptest.NewRequest("POST", "/answer", body)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		data, err := ParseData[CodeSessionResponse](resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "session_id", data.SessionID)
	})
}
