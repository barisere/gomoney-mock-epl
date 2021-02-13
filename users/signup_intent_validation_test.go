package users

import (
	"testing"

	"gomoney-mock-epl/errors"

	"github.com/stretchr/testify/assert"
)

func TestRegistrationIntent_Validate(t *testing.T) {
	var org = SignUpIntent{
		Email:     "admin@example.com",
		FirstName: "First name",
		LastName:  "Last name",
		Password:  "very strong password",
	}

	failOnInternalError := func(t *testing.T, internalError error) {
		assert.Nil(t, internalError, "Validation caused an internal error")
	}

	t.Run("Requires administrator's email, first name, last name, and password", func(t *testing.T) {
		intentCopy := SignUpIntent{}
		validationError, internalError := intentCopy.Validate()
		failOnInternalError(t, internalError)
		assert.Equal(t, "accounts/invalid_org_info", validationError.Code)
		assert.Contains(t, validationError.Details, errors.ValidationErrorDetails{
			Field:   "email",
			Message: "Email is required",
		})
		assert.Contains(t, validationError.Details, errors.ValidationErrorDetails{
			Field:   "first_name",
			Message: "Administrator first name is required",
		})
		assert.Contains(t, validationError.Details, errors.ValidationErrorDetails{
			Field:   "last_name",
			Message: "Administrator last name is required",
		})
		assert.Contains(t, validationError.Details, errors.ValidationErrorDetails{
			Field:   "password",
			Message: "Password is required",
		})
	})
}
