package users

import (
	customErrors "gomoney-mock-epl/errors"
	"math"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type SignUpIntent struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
}

func (i *SignUpIntent) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	i.PasswordHash = string(hash)
	i.Password = ""

	return nil
}

func (i SignUpIntent) Validate() (*customErrors.ValidationError, error) {
	err := v.ValidateStruct(&i,
		v.Field(&i.Email, v.Required.Error("Email is required"), is.Email),
		v.Field(&i.FirstName, v.Required.Error("Administrator first name is required"), v.Length(1, 100)),
		v.Field(&i.LastName, v.Required.Error("Administrator last name is required"), v.Length(1, 100)),
		v.Field(&i.Password, v.Required.Error("Password is required"), v.Length(6, math.MaxUint8)),
	)

	return customErrors.ToValidationError(err,
		"Parts of the data supplied are invalid.",
		"accounts/invalid_org_info")
}
