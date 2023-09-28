package www_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/fdelbos/commons/www"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type (
	TestBody1 struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age"`
	}
)

func TestBody(t *testing.T) {
	app := fiber.New()
	app.Post("/", www.Parser[TestBody1](func(c *fiber.Ctx, body *TestBody1) error {
		return www.Created(c, nil)
	}))

	tc := []struct {
		name string
		body string
		code int
	}{
		{"empty", "", fiber.StatusBadRequest},
		{"optional", `{"name":"test"}`, fiber.StatusCreated},
		{"complete", `{"name":"test", "age": 42}`, fiber.StatusCreated},
		{"empty", `{"name":"", "age": 42}`, fiber.StatusBadRequest},
		{"invalid", `{"name": 42}`, fiber.StatusBadRequest},
		{"required", `{"age": 42}`, fiber.StatusBadRequest},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body := bytes.NewBufferString(c.body)
			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, c.code, resp.StatusCode)
		})
	}

}
