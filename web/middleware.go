package web

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jwtSigningKey = []byte("R4Hw7tAIUqDVDmOx6Cd64+73PIbHCelQjeAo4eh+PuqKK5G+QhjjKXQAjJoBs8Pu/HTJBpN9OoDhpGIhmpbVIzc1Ygzj+m5Ze+8HfcEEsVq1q9Ec6l+DWWc17Zd730k")

const unauthorizedErrorCode = "auth/unauthorized"
const errorCodeForbidden = "auth/restricted-action"

var jwtMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: jwtSigningKey,
	ErrorHandler: func(e error) error {
		return echo.NewHTTPError(http.StatusUnauthorized, errorDto(unauthorizedErrorCode, e.Error()))
	},
})

var onlyAdmins echo.MiddlewareFunc = func(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		isAdmin := claims["is_admin"].(bool)
		if !isAdmin {
			return echo.NewHTTPError(http.StatusForbidden,
				errorDto(errorCodeForbidden, "You're not allowed to access this resource."))
		}
		return hf(c)
	}
}
