package auth

import (
	"context"
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
		New(ctx context.Context, session Session) error
		Get(ctx context.Context, digest []byte) (*Session, error)
		Close(ctx context.Context, digest []byte) error
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

func (s *Sessions) NewSession(ctx context.Context, userID uuid.UUID, duration time.Duration) (string, error) {
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
		until := time.Now().Add(duration).UTC()
		session.Until = &until
	}

	if err := s.store.New(ctx, *session); err != nil {
		return "", err
	}

	return key, nil
}

func (s *Sessions) Get(ctx context.Context, sessionID string) (*Session, error) {
	digest, err := DigestFromAPIKey(sessionID)
	if err != nil {
		return nil, ErrInvalidSession
	}

	session, err := s.store.Get(ctx, digest)
	if err != nil {
		return nil, ErrInvalidSession
	}

	if session.Until != nil {
		if session.Until.Before(time.Now()) {
			defer s.store.Close(ctx, digest)
			return nil, ErrInvalidSession
		}
	}

	return session, nil
}

func (s *Sessions) Close(ctx context.Context, sessionID string) error {
	digest, err := DigestFromAPIKey(sessionID)
	if err != nil {
		return ErrInvalidSession
	}
	return s.store.Close(ctx, digest)
}
