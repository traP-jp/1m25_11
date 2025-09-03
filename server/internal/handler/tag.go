package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "タグ一覧の取得中にエラーが発生しました。",
			Code:    "internal_server_error",
		})
	}

	return c.JSON(http.StatusOK, tagSummaries)
}

func (h *Handler) createTags(c echo.Context) error {
	var body PostTagsJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "リクエストボディが不正です。", Code: "invalid_request_body"})
	}
	userID := uuid.Nil

	newTag, err := h.repo.CreateTags(c.Request().Context(), repository.CreateTagParams{
		Name:      body.Name,
		CreatorID: userID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, &ErrorResponse{Message: "指定された名前のタグは既に存在します。", Code: "tag_already_exists"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグの作成中にエラーが発生しました。", Code: "internal_server_error"})
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
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	tagDetails, err := h.repo.GetTagDetails(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたタグが見つかりません。", Code: "tag_not_found"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグ詳細の取得中にエラーが発生しました。", Code: "internal_server_error"})
	}

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
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	var body PutTagsTagIdJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "リクエストボディが不正です。", Code: "invalid_request_body"})
	}

	err = h.repo.UpdateTags(c.Request().Context(), tagID, body.Name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたタグが見つかりません。", Code: "tag_not_found"})
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, &ErrorResponse{Message: "指定された名前のタグは既に存在します。", Code: "tag_already_exists"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグの更新中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteTags(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	err = h.repo.DeleteTags(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたタグが見つかりません。", Code: "tag_not_found"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグの削除中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) getStampsByTag(c echo.Context) error {
	tagIDStr := c.Param("tagId")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	stampSummaries, err := h.repo.GetStampsByTagID(c.Request().Context(), tagID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたタグが見つかりません。", Code: "tag_not_found"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "スタンプ一覧の取得中にエラーが発生しました。", Code: "internal_server_error"})
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

