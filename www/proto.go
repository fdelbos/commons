package www

import (
	"encoding/json"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type (
	Status string

	Response struct {
		Status  Status      `json:"status"`
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	Entries struct {
		Name string
		Data interface{}
	}

	Body map[string]interface{}
)

const (
	Success Status = "success"
	Fail    Status = "fail"
	Error   Status = "error"
)

func ParseData[T any](r io.Reader) (*T, error) {
	res := struct {
		Data *T `json:"data,omitempty"`
	}{}

	err := json.NewDecoder(r).Decode(&res)
	return res.Data, err
}

func Obj(key string, data interface{}) Entries {
	return Entries{Name: key, Data: data}
}

func respondError(c *fiber.Ctx, statusCode int, message string) error {
	c.Status(statusCode)
	return c.JSON(Response{
		Status:  Error,
		Message: message,
	})
}

func respondData(c *fiber.Ctx, status Status, statusCode int, body interface{}) error {
	c.Status(statusCode)
	return c.JSON(Response{
		Status: status,
		Data:   body,
	})
}

func respondOk(c *fiber.Ctx, statusCode int, body interface{}) error {
	return respondData(c, Success, statusCode, body)
}

func respondKo(c *fiber.Ctx, statusCode int, body interface{}) error {
	return respondData(c, Fail, statusCode, body)
}

// returns a 200 OK response with the given body
func Ok(c *fiber.Ctx, body interface{}) error {
	return respondOk(c, fiber.StatusOK, body)
}

// returns a 201 Created response with the given body
func Created(c *fiber.Ctx, body interface{}) error {
	return respondOk(c, fiber.StatusCreated, body)
}

// returns a 400 BadRequest response with the given body
func BadRequest(c *fiber.Ctx, body interface{}) error {
	return respondKo(c, fiber.StatusBadRequest, body)
}

func Conflict(c *fiber.Ctx, body interface{}) error {
	return respondKo(c, fiber.StatusConflict, body)
}

func ErrNotFound(c *fiber.Ctx) error {
	return respondError(c,
		fiber.StatusNotFound,
		utils.StatusMessage(fiber.StatusNotFound))
}

func ErrInternal(c *fiber.Ctx) error {
	return respondError(c,
		fiber.StatusInternalServerError,
		utils.StatusMessage(fiber.StatusInternalServerError))
}

func ErrToManyRequests(c *fiber.Ctx) error {
	return respondError(c,
		fiber.StatusTooManyRequests,
		utils.StatusMessage(fiber.StatusTooManyRequests))
}

func ErrUnauthorized(c *fiber.Ctx) error {
	return respondError(c,
		fiber.StatusUnauthorized,
		utils.StatusMessage(fiber.StatusUnauthorized))
}

func ErrMethodNotAllowed(c *fiber.Ctx) error {
	return respondError(c,
		fiber.StatusMethodNotAllowed,
		utils.StatusMessage(fiber.StatusMethodNotAllowed))
}

func ErrConflict(c *fiber.Ctx, msg ...string) error {
	errMsg := utils.StatusMessage(fiber.StatusTooManyRequests)
	if len(msg) > 0 {
		errMsg = msg[0]
	}
	return respondError(c, fiber.StatusTooManyRequests, errMsg)
}
