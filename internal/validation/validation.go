package validation

import (
	"errors"
	"strings"

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

func ValidateTasks(tasks models.Task) error {
	if err := validate.Struct(tasks); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

func ValidatePost(post models.SocialPost) error {
	if err := validate.Struct(post); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

func ValidateComment(comment models.Comment) error {
	if comment.Post_id <= 0 {
		return errors.New("invalid post ID")
	}

	if comment.UserID <= 0 {
		return errors.New("invalid user ID")
	}

	content := strings.TrimSpace(comment.Comment)

	if content == "" {
		return errors.New("comment cannot be empty")
	}

	if len(content) > 2000 {
		return errors.New("comment cannot exceed 2000 characters")
	}

	return nil
}

func ValidateNegotiation(negotiate models.SendMessage) error {
	if negotiate.TaskId <= 0 {
		return errors.New("invalid post ID")
	}

	if negotiate.ApplicantID <= 0 {
		return errors.New("invalid user ID")
	}

	content := strings.TrimSpace(negotiate.Message)

	if content == "" {
		return errors.New("message cannot be empty")
	}

	return nil
}
