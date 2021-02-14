package web

import (
	"net/http"

	"gomoney-mock-epl/errors"

	"github.com/labstack/echo"
)

// DefaultErrorHandler formats every error that bubbles up through
// the Echo server, and sends them out using an appropriate status code.
func DefaultErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	switch err.(type) {
	case *errors.ValidationError:
	case errors.ValidationError:
		_ = c.JSON(http.StatusUnprocessableEntity, err)
	case *echo.HTTPError:
		_err := err.(*echo.HTTPError)
		_ = c.JSON(_err.Code, _err.Message)
	default:
		_ = c.JSON(http.StatusInternalServerError, errorDto("internal_error", err.Error()))
	}
}
