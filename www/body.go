package www

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fdelbos/commons/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	Pagination struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}
)

const DefaultMaxOffset = int64(math.MaxInt16)
const DefaultMaxSize = int64(100)

func validationErrors(c *fiber.Ctx, err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	res := map[string]interface{}{}
	errors := err.(validator.ValidationErrors)
	for _, err := range errors {
		res[err.Field()] = err.Translate(validation.Translator())
	}

	return BadRequest(c, Body{
		"validation": res,
	})
}

// Parser is a middleware that parse the body of the request and validates it.
// An invalid validation returns a 400 Bad Request with a JSON body containing the validation errors.
// Here is an example response: {"status":"fail","data":{"validation":{"name":"name is a required field"}}}
func Parser[T any](next func(*fiber.Ctx, *T) error) fiber.Handler {
	return func(c *fiber.Ctx) error {

		body := new(T)

		if err := c.BodyParser(&body); err != nil {
			return respondError(c,
				fiber.StatusBadRequest,
				"invalid encoding")
		}

		if err := validation.Validator().Struct(body); err != nil {
			return validationErrors(c, err)
		}

		return next(c, body)
	}
}

func parseInt64QueryParam(c *fiber.Ctx, key string, min, max, defaultValue int64) (int64, error) {
	value := c.Query(key)

	// If the parameter is not found, the default value is returned.
	if value == "" {
		return defaultValue, nil
	}

	// If the parameter is found but is not an integer, a 400 Bad Request is returned.
	nb, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0,
			BadRequest(c, fmt.Sprintf("invalid query parameter '%s', must be an integer between %d and %d", key, min, max))
	}

	if nb < min || nb > max {
		return 0,
			BadRequest(c, fmt.Sprintf("invalid query parameter '%s', must be an integer between %d and %d", key, min, max))
	}
	return nb, nil
}

// Paginated is a middleware that parses the limit and offset parameters from the query string.
// takes two optional parameters: the maximum limit and the maximum offset.
//
// Only GET requests are allowed to use this middleware
func Paginated(next func(*fiber.Ctx, Pagination) error, params ...int64) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Method() != fiber.MethodGet {
			return BadRequest(c, "invalid method")
		}

		pagination := Pagination{}
		var err error

		maxSize := DefaultMaxSize
		if len(params) > 0 {
			maxSize = params[0]
		}
		pagination.Limit, err = parseInt64QueryParam(c, "limit", 1, maxSize, 10)
		if err != nil {
			return err
		}

		maxOffset := DefaultMaxOffset
		if len(params) > 1 {
			maxOffset = params[1]
		}
		pagination.Offset, err = parseInt64QueryParam(c, "offset", 0, maxOffset, 0)
		if err != nil {
			return err
		}
		return next(c, pagination)
	}
}

func Int64(c *fiber.Ctx, key string) (int64, error) {
	value := c.Params(key)
	if value == "" {
		return 0, BadRequest(c, fmt.Sprintf("invalid parameter '%s'", key))
	}
	nb, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, BadRequest(c, fmt.Sprintf("invalid parameter '%s'", key))
	}
	return nb, nil
}
