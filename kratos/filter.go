package kratos

import (
	"github.com/fdelbos/commons/www"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type (
	kratosCtx string
)

const (
	UserHeader = "X-User-ID"
	userCtx    = kratosCtx("utils/kratos/user")
)

func Filter(c *fiber.Ctx) error {
	if userId := c.Get(UserHeader); userId == "" {
		return www.ErrUnauthorized(c)

	} else if id, err := uuid.Parse(userId); err != nil {
		return www.ErrUnauthorized(c)

	} else {
		c.Locals(userCtx, id)
		return c.Next()
	}
}

func GetUserID(f *fiber.Ctx) uuid.UUID {
	obj := f.Locals(userCtx)
	if obj == nil {
		log.Panic().
			Str("url", f.OriginalURL()).
			Msg("kratos user ID not available in http.Request context")
	}
	return obj.(uuid.UUID)
}
