package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getRanking(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}