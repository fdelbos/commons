package www

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
)

type (
	ServerConfig struct {
		BaseURL   string
		Port      int
		Views     fiber.Views
		Timeout   time.Duration
		BodyLimit int
	}

	Route func(r fiber.Router)

	Option func(*ServerConfig)
)

// Serve starts an instance of the http server.
func Serve(route Route, config ServerConfig) error {
	app := fiber.New(fiber.Config{
		Views:                 config.Views,
		DisableStartupMessage: true,
		JSONDecoder:           sonic.Unmarshal,
		JSONEncoder:           sonic.Marshal,
		CaseSensitive:         true,
		StreamRequestBody:     true,
		ReadTimeout:           config.Timeout,
		WriteTimeout:          config.Timeout,
		BodyLimit:             config.BodyLimit,
	})
	app.Use(logger.New())

	group := app.Group(config.BaseURL)

	route(group)
	addr := fmt.Sprintf(":%d", config.Port)
	log.Info().
		Str("addr", addr+config.BaseURL).
		Msg("starting http server")
	return app.Listen(addr)
}
