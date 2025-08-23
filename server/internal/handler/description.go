package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type descriptionPayload struct {
	Description string `json:"description"`
}

func (h *Handler) createDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	creatorID := uuid.Nil // 仮でNil UUIDを用いている
	payload := new(descriptionPayload)
	if err = c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if payload.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(errors.New("description cannot be empty"))
	}
	err = h.repo.CreateDescriptions(c.Request().Context(), repository.CreateDescriptionParams{
		StampID:     stampID,
		Description: payload.Description,
		CreatorID:   creatorID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) getDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	descriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, descriptions)
}

func (h *Handler) updateDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	creatorID := uuid.Nil // 仮でNil UUIDを用いている
	payload := new(descriptionPayload)
	if err = c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if payload.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(errors.New("description cannot be empty"))
	}
	if err = h.repo.UpdateDescriptions(c.Request().Context(), stampID, creatorID, payload.Description); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	creatorID := uuid.Nil // 仮でNil UUIDを用いている
	if err = h.repo.DeleteDescriptions(c.Request().Context(), stampID, creatorID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
