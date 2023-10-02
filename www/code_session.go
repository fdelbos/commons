package www

import (
	"context"
	"time"

	"github.com/fdelbos/commons/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	EmailToUUID func(ctx context.Context, email string) (uuid.UUID, error)

	CodesService interface {
		NewCode(email string) (string, *auth.Code)
		Send(ctx context.Context, to string) error
		Validate(ctx context.Context, digits string, email string) error
	}

	SessionsService interface {
		Close(ctx context.Context, sessionID string) error
		Get(ctx context.Context, sessionID string) (*auth.Session, error)
		NewSession(ctx context.Context, userID uuid.UUID, duration time.Duration) (string, error)
	}

	CodeSession struct {
		codes       CodesService
		emailToUUID EmailToUUID
		sessions    SessionsService
		duration    time.Duration
	}

	CodeSessionRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	CodeSessionAnswer struct {
		Email string `json:"email" validate:"required,email"`
		Code  string `json:"code" validate:"required"`
	}

	CodeSessionResponse struct {
		SessionID string `json:"session_id"`
	}
)

const (
	DefaultCodeDuration = time.Minute * 5
)

func NewCodeSession(codes CodesService, emailToUUID EmailToUUID, sessions SessionsService, opts ...func(*CodeSession)) *CodeSession {
	cs := &CodeSession{
		codes:       codes,
		emailToUUID: emailToUUID,
		sessions:    sessions,
		duration:    DefaultCodeDuration,
	}

	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

func (css *CodeSession) Routes(r fiber.Router) {
	r.Post("/send", Parser[CodeSessionRequest](css.Send))
	r.Post("/answer", Parser[CodeSessionAnswer](css.Answer))
}

func (css *CodeSession) Send(c *fiber.Ctx, req *CodeSessionRequest) error {
	err := css.codes.Send(c.Context(), req.Email)
	if err != nil {
		return ErrInternal(c)
	}

	return Ok(c, nil)
}

func (css *CodeSession) Answer(c *fiber.Ctx, req *CodeSessionAnswer) error {
	err := css.codes.Validate(c.Context(), req.Code, req.Email)
	if err != nil {
		return ErrUnauthorized(c)
	}

	userID, err := css.emailToUUID(c.Context(), req.Email)
	if err != nil {
		return ErrUnauthorized(c)
	}

	// lets create a new session
	sessionID, err := css.sessions.NewSession(c.Context(), userID, css.duration)
	if err != nil {
		return ErrInternal(c)
	}
	return Created(c, &CodeSessionResponse{
		SessionID: sessionID,
	})
}

func WithDuration(d time.Duration) func(*CodeSession) {
	return func(css *CodeSession) {
		css.duration = d
	}
}
