package auth_test

import (
	"errors"
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

	duration := 1 * time.Hour
	userID := uuid.New()

	// test new session
	var mockDigest []byte
	store.
		On("New", mock.Anything).
		Return(func(session Session) error {
			assert.Equal(t, userID, session.UserID)
			assert.True(t, session.Until.After(time.Now()))
			assert.True(t, session.Until.Before(time.Now().Add(duration+time.Minute)))
			assert.NotEmpty(t, session.Digest)
			mockDigest = session.Digest
			return nil
		}).
		Once()

	key, err := sessions.NewSession(userID, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	digest, err := DigestFromAPIKey(key)
	assert.NoError(t, err)
	assert.Equal(t, mockDigest, digest)

	// test get session
	until := time.Now().Add(duration)
	store.
		On("Get", mock.Anything).
		Return(func(digest []byte) (*Session, error) {
			assert.Equal(t, mockDigest, digest)
			return &Session{
				Digest: digest,
				UserID: userID,
				Until:  &until,
			}, nil
		}).
		Once()

	session, err := sessions.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, until, *session.Until)
	assert.Equal(t, mockDigest, session.Digest)

	// test close session
	store.
		On("Close", mock.Anything).
		Return(func(digest []byte) error {
			assert.Equal(t, mockDigest, digest)
			return nil
		}).
		Once()

	err = sessions.Close(key)
	assert.NoError(t, err)

	// test invalid session
	store.
		On("Get", mock.Anything).
		Return(func(digest []byte) (*Session, error) {
			assert.Equal(t, mockDigest, digest)
			return nil, errors.New("not found")
		}).
		Once()

	_, err = sessions.Get(key)
	assert.Error(t, err)

	store.AssertExpectations(t)
}
