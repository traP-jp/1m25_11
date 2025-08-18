package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getStamps(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func (h *Handler) getCertainStamps(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}