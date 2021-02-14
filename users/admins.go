package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Administrator struct {
	ID           string `json:"id" bson:"_id"`
	Email        string `json:"email" bson:"email"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	PasswordHash string `json:"-" bson:"password_hash"`
}

func SignUpAdmin(ctx context.Context, intent SignUpIntent, db AdminsDB) (*Administrator, error) {
	validationErr, internalErr := intent.Validate()
	if validationErr != nil {
		return nil, validationErr
	}
	if internalErr != nil {
		return nil, fmt.Errorf("An unknown error occurred. Please try again. %w", internalErr)
	}
	intentCopy := intent
	if err := (&intentCopy).HashPassword(); err != nil {
		return nil, err
	}
	admin := &Administrator{
		Email:        intent.Email,
		FirstName:    intent.FirstName,
		LastName:     intent.LastName,
		PasswordHash: intentCopy.PasswordHash,
	}
	admin, err := db.Create(ctx, *admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var unimplemented = fmt.Errorf("unimplemented")

var ErrIncorrectLogin = errors.New("incorrect login credentials")

func LoginAsAdmin(ctx context.Context, db AdminsDB, dto LoginDto) (*jwt.Token, error) {
	admin, err := db.ByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrIncorrectLogin
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(dto.Password))
	if err != nil {
		return nil, ErrIncorrectLogin
	}
	return makeJWT(JwtRequest{subject: admin.ID, IsAdmin: true}), nil
}
