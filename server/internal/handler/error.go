package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	res := ErrorResponse{
		Message: "Internal Server Error",
		Code:    "internal_server_error",
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		switch msg := he.Message.(type) {
		case ErrorResponse:
			res = msg
		case *ErrorResponse:
			res = *msg
		case string:
			res.Message = msg
		default:
			res.Message = "An unexpected error occurred"
		}
	} else {
		c.Logger().Error(err)
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, res)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
