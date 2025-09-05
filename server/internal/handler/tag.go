package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type Error struct {
	Message string `json:"message"`
}

type TagSummary struct {
	Id   uuid.UUID `json:"tag_id"`
	Name string    `json:"tag_name"`
}

type StampSummary struct {
	Id     uuid.UUID `json:"stamp_id"`
	Name   string    `json:"stamp_name"`
	FileId uuid.UUID `json:"file_id"`
}

type TagDetails struct {
	ID        uuid.UUID `json:"tag_id"`
	Name      string    `json:"tag_name"`
	CreatorID uuid.UUID `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stamps    []struct {
		ID     uuid.UUID `json:"stamp_id"`
		Name   string    `json:"stamp_name"`
		FileID uuid.UUID `json:"file_id"`
	}
}

type Tag struct {
	Id        uuid.UUID      `json:"tag_id"`
	Name      string         `json:"tag_name"`
	CreatorId uuid.UUID      `json:"creator_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Count     int            `json:"count"`
	Stamps    []StampSummary `json:"stamps"`
}

type PostTagsJSONRequestBody struct {
	Name string `json:"name"`
}

type PutTagsTagIdJSONRequestBody struct {
	Name string `json:"name"`
}

func (h *Handler) getTags(c echo.Context) error {
	tagSummaries, err := h.repo.GetTags(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to get tags: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, tagSummaries)
}

func (h *Handler) createTags(c echo.Context) error {
	var body PostTagsJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid request body.",
		})
	}
	creatorID, err := h.getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}

	newTag, err := h.repo.CreateTags(c.Request().Context(), repository.CreateTagParams{
		Name:      body.Name,
		CreatorID: creatorID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrTagConflict) {
			return echo.NewHTTPError(http.StatusConflict, Error{
				Message: "Tag with this name already exists.",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to create tag: %s", err.Error()),
		})
	}

	response := TagSummary{
		Id:   newTag,
		Name: body.Name,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) getTagDetails(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid tag ID format.",
		})
	}

	tagDetailsRaw, err := h.repo.GetTagDetails(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, Error{
				Message: "Tag not found.",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to get tag details: %s", err.Error()),
		})
	}

	if tagDetailsRaw == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: "internal error: required tag details were not found",
		})
	}
	tagDetails := tagDetailsRaw

	stamps := make([]StampSummary, len(tagDetails.Stamps))
	for i, s := range tagDetails.Stamps {
		stamps[i] = StampSummary{
			Id:     s.ID,
			Name:   s.Name,
			FileId: s.FileID,
		}
	}

	response := Tag{
		Id:        tagDetails.ID,
		Name:      tagDetails.Name,
		CreatorId: tagDetails.CreatorID,
		CreatedAt: tagDetails.CreatedAt,
		UpdatedAt: tagDetails.UpdatedAt,
		Count:     len(stamps),
		Stamps:    stamps,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) updateTags(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid tag ID format.",
		})
	}

	var body PutTagsTagIdJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid request body.",
		})
	}

	err = h.repo.UpdateTags(c.Request().Context(), tagID, body.Name)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, Error{
				Message: "Tag not found.",
			})
		}
		if errors.Is(err, repository.ErrTagConflict) {
			return echo.NewHTTPError(http.StatusConflict, Error{
				Message: "Tag with this name already exists.",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to update tag: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteTags(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid tag ID format.",
		})
	}

	// isAdmin := false
	// if !isAdmin {
	// 	return echo.NewHTTPError(http.StatusForbidden, Error{
	// 		Message: "You do not have permission to delete this tag.",
	// 	})
	// }

	err = h.repo.DeleteTags(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, Error{
				Message: "Tag not found.",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to delete tag: %s", err.Error()),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) getStampsByTag(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Error{
			Message: "Invalid tag ID format.",
		})
	}

	stampSummaries, err := h.repo.GetStampsByTagID(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrTagNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, Error{
				Message: "Tag not found.",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed to get stamps by tag: %s", err.Error()),
		})
	}

	response := make([]StampSummary, len(stampSummaries))
	for i, s := range stampSummaries {
		response[i] = StampSummary{
			Id:     s.ID,
			Name:   s.Name,
			FileId: s.FileID,
		}
	}

	return c.JSON(http.StatusOK, response)
}
