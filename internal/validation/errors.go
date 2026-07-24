package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {

		// Return the first validation error
		e := validationErrors[0]

		field := e.Field()

		switch e.Tag() {

		case "required":
			return fmt.Errorf("%s is required", field)

		case "email":
			return fmt.Errorf("%s must be a valid email address", field)

		case "min":
			return fmt.Errorf("%s must be at least %s characters", field, e.Param())

		case "max":
			return fmt.Errorf("%s cannot exceed %s characters", field, e.Param())

		case "gt":
			return fmt.Errorf("%s must be greater than %s", field, e.Param())

		case "len":
			return fmt.Errorf("%s must be exactly %s characters", field, e.Param())

		case "numeric":
			return fmt.Errorf("%s must contain only numbers", field)

		default:
			return fmt.Errorf("%s is invalid", field)
		}
	}

	return err
}
