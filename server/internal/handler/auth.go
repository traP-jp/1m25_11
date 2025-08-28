package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) login(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func (h *Handler) callback(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}