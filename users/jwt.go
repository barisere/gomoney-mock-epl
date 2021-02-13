package users

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	JWTIssuer = "mock-epl.io"
)

var JWTSigningMethod = jwt.SigningMethodHS256

type JwtRequest struct {
	subject string
	IsAdmin bool `json:"is_admin"`
}

type jwtClaims struct {
	*jwt.StandardClaims
	*JwtRequest
}

func makeJWT(request JwtRequest) *jwt.Token {
	oneHourFromNow := time.Now().Add(time.Hour * 1).Unix()
	claims := jwtClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: oneHourFromNow,
			IssuedAt:  time.Now().Unix(),
			Issuer:    JWTIssuer,
			Subject:   request.subject,
		},
		JwtRequest: &JwtRequest{
			IsAdmin: request.IsAdmin,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
