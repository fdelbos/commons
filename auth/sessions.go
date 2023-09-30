package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type (
	Session struct {
		Digest []byte
		UserID uuid.UUID
		Until  *time.Time
	}
	SessionsStore interface {
		New(session Session) error
		Get(digest []byte) (*Session, error)
		Close(digest []byte) error
	}

	Sessions struct {
		store SessionsStore
	}
)

const (
	Forever = 0
)

var ErrInvalidSession = errors.New("invalid session")

func NewSessions(store SessionsStore) *Sessions {
	s := &Sessions{
		store: store,
	}
	return s
}

func (s *Sessions) NewSession(userID uuid.UUID, duration time.Duration) (string, error) {
	key, err := NewApiKey()
	if err != nil {
		return "", err
	}

	digest, err := DigestFromAPIKey(key)
	if err != nil {
		return "", err
	}

	session := &Session{
		Digest: digest,
		UserID: userID,
	}
	if duration != Forever {
		until := time.Now().Add(duration)
		session.Until = &until
	}
	if err := s.store.New(*session); err != nil {
		return "", err
	}

	return key, nil
}

func (s *Sessions) Get(sessionID string) (*Session, error) {
	digest, err := DigestFromAPIKey(sessionID)
	if err != nil {
		return nil, ErrInvalidSession
	}

	session, err := s.store.Get(digest)
	if err != nil {
		return nil, ErrInvalidSession
	}

	if session.Until != nil {
		if session.Until.Before(time.Now()) {
			defer s.store.Close(digest)
			return nil, ErrInvalidSession
		}
	}

	return session, nil
}

func (s *Sessions) Close(sessionID string) error {
	digest, err := DigestFromAPIKey(sessionID)
	if err != nil {
		return ErrInvalidSession
	}
	return s.store.Close(digest)
}
