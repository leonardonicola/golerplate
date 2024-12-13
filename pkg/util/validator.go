package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationResponse represents the error response structure
type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationError formats validation errors into a consistent response
func HandleValidationError(err error) *ValidationResponse {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationError

		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: getValidationErrorMessage(e),
			})
		}

		return &ValidationResponse{
			Errors: errors,
		}
	}

	// If it's not a validation error, return a generic error
	return &ValidationResponse{
		Errors: []ValidationError{
			{
				Field:   "request",
				Message: err.Error(),
			},
		},
	}
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", err.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", err.Param())
	case "cpf":
		return "Invalid CPF format"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	default:
		return err.Error()
	}
}
