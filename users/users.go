package users

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id" bson:"_id"`
	Email        string `json:"email" bson:"email"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	PasswordHash string `json:"-" bson:"password_hash"`
}

func SignUpUser(ctx context.Context, intent SignUpIntent, db UsersDB) (*User, error) {
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
	user := User{
		Email:        intent.Email,
		FirstName:    intent.FirstName,
		LastName:     intent.LastName,
		PasswordHash: intentCopy.PasswordHash,
	}
	err := db.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func LoginAsUser(ctx context.Context, db UsersDB, dto LoginDto) (*jwt.Token, error) {
	user, err := db.ByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrIncorrectLogin
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))
	if err != nil {
		return nil, ErrIncorrectLogin
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"is_admin": false,
		"id":       user.ID,
	})
	return token, nil
}
