package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getTags(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *Handler) createTags(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *Handler) updateTags(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *Handler) deleteTags(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}