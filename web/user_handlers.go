package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"gomoney-mock-epl/users"

	"github.com/labstack/echo"
)

func adminSignUpHandler(db users.AdminsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var signupIntent users.SignUpIntent
		if err := c.Bind(&signupIntent); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		admin, err := users.SignUpAdmin(c.Request().Context(), signupIntent, db)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated,
			dataResponse("auth/admin-account", "Admin account created", admin))
	}
}

type loginResponse struct {
	Token string `json:"token"`
}

const adminLoginResponseType = "auth/admin-login"

func adminLoginHandler(db users.AdminsDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginDto users.LoginDto
		if err := c.Bind(&loginDto); err != nil {
			return err
		}
		token, err := users.LoginAsAdmin(c.Request().Context(), db, loginDto)
		if err != nil {
			if errors.Is(err, users.ErrIncorrectLogin) {
				return echo.NewHTTPError(http.StatusUnauthorized, errorDto("auth/unauthorised", err.Error()))
			}
			return err
		}
		tokenString, err := token.SignedString(jwtSigningKey)
		if err != nil {
			return err
		}
		response := dataResponse(adminLoginResponseType, "Logged in successfully", loginResponse{
			Token: tokenString,
		})
		return c.JSON(http.StatusOK, response)
	}
}

func AdminSignupRoute(db users.AdminsDB) RouteProvider {
	return func(e *echo.Echo) {
		e.POST("/signup/admins/", adminSignUpHandler(db), jwtMiddleware)
		e.POST("/login/admins/", adminLoginHandler(db))
	}
}
