package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	creatorID, ok := c.Get(userIDContextKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	user, err := h.repo.GetUser(c.Request().Context(), creatorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, user)
}
