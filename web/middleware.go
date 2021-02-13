package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var jwtSigningKey = []byte("R4Hw7tAIUqDVDmOx6Cd64+73PIbHCelQjeAo4eh+PuqKK5G+QhjjKXQAjJoBs8Pu/HTJBpN9OoDhpGIhmpbVIzc1Ygzj+m5Ze+8HfcEEsVq1q9Ec6l+DWWc17Zd730k")

var unauthorizedErrorCode = "auth/unauthorized"

var jwtMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: jwtSigningKey,
	ErrorHandler: func(e error) error {
		return echo.NewHTTPError(http.StatusUnauthorized, errorDto(unauthorizedErrorCode, e.Error()))
	},
})
