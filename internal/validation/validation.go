package validation

import (
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRegister(register models.Register) error {
	if err := validate.Struct(register); err != nil {
		return FormatValidationError(err)
	}
	return nil
}


func ValidateLogin(student models.Login) error {

	if err := validate.Struct(student); err != nil {
		return FormatValidationError(err)
	}
	return nil
}