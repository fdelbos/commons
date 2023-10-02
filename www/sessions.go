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
	SessionParamName  = "session"
)

// FilterSession is a middleware that checks for a session in the request.
// The session can be provided in the following ways:
// - Authorization header: Bearer <session>
// - Cookie: session_auth=<session>
// - Query param: session=<session>
func FilterSession(sessions *auth.Sessions) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionID := ""

		// check for bearer token
		header := c.GetReqHeaders()["Authorization"]
		if header != "" {
			id, found := strings.CutPrefix(header, "Bearer")
			if found {
				sessionID = strings.TrimSpace(id)
			}
		}

		// check for cookie
		if sessionID == "" {
			sessionID = c.Cookies(SessionCookieName)
		}

		// check for query param
		if sessionID == "" {
			sessionID = c.Query(SessionParamName)
		}

		if sessionID == "" {

			return ErrUnauthorized(c)
		}

		session, err := sessions.Get(c.Context(), sessionID)
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

// GetSession returns the session from the fiber context.
func GetSession(f *fiber.Ctx) *auth.Session {
	obj := f.Locals(sessionCtx)
	if obj == nil {
		log.Fatal("session not found in fiber context")
	}
	return obj.(*auth.Session)
}
