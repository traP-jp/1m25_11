package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := h.getUserID(c); err != nil {
			return err
		}

		return next(c)
	}
}
