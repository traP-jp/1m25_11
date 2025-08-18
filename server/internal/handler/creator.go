package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getCreatorDetails(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}