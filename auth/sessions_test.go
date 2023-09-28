package auth_test

import (
	"testing"
	"time"

	. "github.com/fdelbos/commons/auth"
	"github.com/fdelbos/commons/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSessions(t *testing.T) {
	store := mocks.NewAuthSessionsStore(t)
	sessions := NewSessions(store)

	until := time.Now().Add(time.Minute)
	userID := uuid.New()

	// test new session
	mockSessionID := ""
	store.
		On("New", mock.Anything).
		Return(func(session Session) error {
			assert.Equal(t, userID, session.UserID)
			assert.True(t, session.Until.After(time.Now()))
			assert.Equal(t, SessionTypeAPI, session.Type)
			assert.Equal(t, len(session.SessionID), SessionIDLen)
			mockSessionID = session.SessionID
			return nil
		}).
		Once()

	session, err := sessions.NewAPISession(userID, until)
	assert.NoError(t, err)
	assert.Equal(t, mockSessionID, session.SessionID)

	store.AssertExpectations(t)
}
