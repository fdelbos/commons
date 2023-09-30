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

	SessionCookieName = "session_auth"
)

func FilterSession(sessions *auth.Sessions) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionID := ""

		header := c.GetReqHeaders()["Authorization"]
		if header != "" {
			id, found := strings.CutPrefix(header, "Bearer")
			if found {
				sessionID = strings.TrimSpace(id)
			}
		} else {
			sessionID = c.Cookies(SessionCookieName)
		}

		if sessionID == "" {
			return ErrUnauthorized(c)
		}
		session, err := sessions.Get(sessionID)
		if err != nil {
			return ErrUnauthorized(c)
		}

		c.Locals(sessionCtx, &auth.Session{
			UserID: session.UserID,
			Until:  session.Until,
		})
		return c.Next()
	}
}

func GetSession(f *fiber.Ctx) *auth.Session {
	obj := f.Locals(sessionCtx)
	if obj == nil {
		log.Fatal("session not found in fiber context")
	}
	return obj.(*auth.Session)
}
