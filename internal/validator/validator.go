package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type requestValidator struct {
	v *validator.Validate
}

func NewValidator() *requestValidator {
	return &requestValidator{
		v: validator.New(),
	}
}

func (v *requestValidator) Validate(in any) error {
	if err := v.v.Struct(in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
