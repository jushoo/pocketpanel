package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// New creates a new validator instance that implements fiber.StructValidator interface.
func New() fiber.StructValidator {
	return &fiberValidator{validate: validator.New()}
}

type fiberValidator struct {
	validate *validator.Validate
}

// Validate validates the given struct.
func (v *fiberValidator) Validate(i any) error {
	return v.validate.Struct(i)
}
