package auth

import (
	"errors"
	"time"

	"github.com/dchest/uniuri"
	"github.com/google/uuid"
)

type (
	SessionType string

	SessionsStore interface {
		New(session Session) error
		Get(sessionID string) (*Session, error)
		All(userID uuid.UUID) ([]*Session, error)
		Close(sessionID string) error
	}

	Sessions struct {
		store SessionsStore
	}

	Session struct {
		SessionID string
		UserID    uuid.UUID
		Until     *time.Time
		Type      SessionType
	}
)

const (
	SessionIDLen = 28

	SessionTypeHTML SessionType = "html"
	SessionTypeAPI  SessionType = "api"
)

var ErrInvalidSession = errors.New("invalid session")

func NewSessions(store SessionsStore) *Sessions {
	s := &Sessions{
		store: store,
	}
	return s
}

func (s *Sessions) genSessionID() string {
	return uniuri.NewLen(SessionIDLen)
}

func (s *Sessions) NewHTMLSession(userID uuid.UUID, until *time.Time) (*Session, error) {
	session := &Session{
		SessionID: s.genSessionID(),
		UserID:    userID,
		Until:     until,
		Type:      SessionTypeHTML,
	}
	if err := s.store.New(*session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Sessions) NewAPISession(userID uuid.UUID, until *time.Time) (*Session, error) {
	session := &Session{
		SessionID: s.genSessionID(),
		UserID:    userID,
		Until:     until,
		Type:      SessionTypeAPI,
	}
	if err := s.store.New(*session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Sessions) Get(sessionID string) (*Session, error) {
	session, err := s.store.Get(sessionID)
	if err != nil {
		return nil, ErrInvalidSession
	}
	if session.Until != nil {
		if session.Until.Before(time.Now()) {
			defer s.store.Close(sessionID)
			return nil, ErrInvalidSession
		}
	}

	return session, nil
}

func (s *Sessions) Close(sessionID string) error {
	return s.store.Close(sessionID)
}

func (s *Sessions) CloseAll(userID uuid.UUID, except string) error {
	sessions, err := s.store.All(userID)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if session.SessionID != except {
			if err := s.store.Close(session.SessionID); err != nil {
				return err
			}
		}
	}
	return nil
}
