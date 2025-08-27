package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

func (h *Handler) createStampTags(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		if errors.Is(err, repository.ErrStampNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	creatorID := uuid.Nil // 仮でNil UUIDを用いている
	err = h.repo.CreateStampTags(c.Request().Context(), repository.CreateStampTagParams{
		StampID:   stampID,
		TagID:     tagID,
		CreatorID: creatorID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrTagAlreadyAdded) {
			return echo.NewHTTPError(http.StatusConflict).SetInternal(err)
		}
		if errors.Is(err, repository.ErrUnauthorized) {
			return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
		}
		if errors.Is(err, repository.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
func (h *Handler) deleteStampTags(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		if errors.Is(err, repository.ErrStampNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	err = h.repo.DeleteStampTags(c.Request().Context(), stampID, tagID)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotLinked) {
			return echo.NewHTTPError(http.StatusNotFound).SetInternal(err)
		}
		if errors.Is(err, repository.ErrUnauthorized) {
			return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
		}
		if errors.Is(err, repository.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
