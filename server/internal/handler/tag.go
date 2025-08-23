package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/api"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

func (h *Handler) getTags(c echo.Context) error {
	tagSummaries, err := h.repo.GetTags(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("failed to get tags: %s", err.Error()),
		})
	}

	response := make([]api.TagSummary, len(tagSummaries))
	for i, tag := range tagSummaries {
		response[i] = api.TagSummary{
			Id:   tag.ID,
			Name: tag.Name,
		}
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) createTags(c echo.Context) error {
	var body api.PostTagsJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid request body.",
		})
	}
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, api.Error{
			Message: "Authentication failed.",
		})
	}

	newTag, err := h.repo.CreateTags(c.Request().Context(), repository.CreateTagParams{
		Name:     body.Name,
		CreatorID: userID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrTagConflict) {
			return echo.NewHTTPError(http.StatusConflict, api.Error{
				Message: "Tag with this name already exists.",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("failed to create tag: %s", err.Error()),
		})
	}
	response := api.Tag{
		Id:        newTag.ID,
		Name:      newTag.Name,
		CreatorId: newTag.CreatorID,
		CreatedAt: newTag.CreatedAt,
		UpdatedAt: newTag.UpdatedAt,
		Count:     0,
		Stamps:    []api.StampSummary{},
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) updateTags(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid tag ID format.",
		})
	}

	var body api.PutTagsTagIdJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid request body.",
		})
	}

	err = h.repo.UpdateTags(c.Request().Context(), tagID, body.Name)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, api.Error{
				Message: "Tag not found.",
			})
		}
		if errors.Is(err, repository.ErrTagConflict) {
			return echo.NewHTTPError(http.StatusConflict, api.Error{
				Message: "Tag with this name already exists.",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("failed to update tag: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteTags(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid tag ID format.",
		})
	}

	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, api.Error{
			Message: "Authentication failed.",
		})
	}

	isAdminRaw, err := h.repo.IsAdmin(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: "failed to check admin status.",
		})
	}
	isAdmin, ok := isAdminRaw.(bool)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: "failed to assert admin status.",
		})
	}
	if !isAdmin {
		return echo.NewHTTPError(http.StatusForbidden, api.Error{
			Message: "You do not have permission to delete this tag.",
		})
	}

	err = h.repo.DeleteTags(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, api.Error{
				Message: "Tag not found.",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("failed to delete tag: %s", err.Error()),
		})
	}
	return c.NoContent(http.StatusNoContent)
}
