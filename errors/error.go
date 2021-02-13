package errors

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	"strings"
)

type ApplicationError struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

func (a ApplicationError) Error() string {
	return fmt.Sprintf("%s: %s", a.Code, a.Message)
}

type ValidationErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Details []ValidationErrorDetails `json:"details;omitempty"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Code, v.Message)
}

func toErrorDetails(errors v.Errors) []ValidationErrorDetails {
	details := make([]ValidationErrorDetails, 0, len(errors))
	for field, err := range errors {
		messageParts := strings.Split(err.Error(), ";")
		for _, message := range messageParts {
			message = strings.TrimSpace(message)
			if message != "" {
				details = append(details, ValidationErrorDetails{
					Field:   field,
					Message: message,
				})
			}
		}
	}
	return details
}

func ToValidationError(errors error, message, code string) (*ValidationError, error) {
	if e, ok := errors.(v.Errors); ok {
		return &ValidationError{
				Code:    code,
				Message: message,
			Details: toErrorDetails(e),
		}, nil
	}
	return nil, errors
}
