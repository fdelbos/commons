package www

import (
	"log"
	"strings"

	"github.com/fdelbos/commons/auth"
	"github.com/gofiber/fiber/v2"
)

type (
	ctx string
)

const (
	sessionCtx = ctx("commons/www/session")

	SessionCookie = "auth"
)

func FilterSession(sessions *auth.Sessions) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(SessionCookie)
		if sessionID == "" {
			header := c.GetReqHeaders()["Authorization"]
			id, found := strings.CutPrefix(header, "Bearer")
			if found {
				sessionID = strings.TrimSpace(id)
			}
		}
		if sessionID == "" {
			return ErrUnauthorized(c)
		}

		session, err := sessions.Get(sessionID)
		if err != nil {
			return ErrUnauthorized(c)
		}

		c.Locals(sessionCtx, session)
		return c.Next()
	}
}

func SetSessionCookie(f *fiber.Ctx, sessionID string) {
	f.Cookie(&fiber.Cookie{
		Name:     SessionCookie,
		Value:    sessionID,
		HTTPOnly: true,
	})
}

func GetSession(f *fiber.Ctx) *auth.Session {
	obj := f.Locals(sessionCtx)
	if obj == nil {
		log.Fatal("session not found in fiber context")
	}
	return obj.(*auth.Session)
}
