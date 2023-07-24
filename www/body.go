package www

import (
	"github.com/fdelbos/commons/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

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

func Parser[T any](next func(*fiber.Ctx, *T) error) func(*fiber.Ctx) error {
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
