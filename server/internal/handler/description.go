package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) createDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func (h *Handler) getDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func (h *Handler) updateDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func (h *Handler) deleteDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}