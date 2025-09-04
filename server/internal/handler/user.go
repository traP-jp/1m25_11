package handler

import (
	"errors"
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

// --- スキーマ定義 ---
type (
	GetUsersResponse []GetUserResponse

	GetUserResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateUserResponse struct {
		ID uuid.UUID `json:"id"`
	}
)

// GET /api/v1/users
func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.repo.GetUsers(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "ユーザー一覧の取得中にエラーが発生しました。",
			Code:    "internal_server_error",
		})
	}

	res := make(GetUsersResponse, len(users))
	for i, user := range users {
		res[i] = GetUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return c.JSON(http.StatusOK, res)
}

// POST /api/v1/users
func (h *Handler) CreateUser(c echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Message: "リクエストボディの形式が不正です。",
			Code:    "invalid_request_body",
		})
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Message: "リクエストのバリデーションに失敗しました。",
			Code:    "validation_failed",
			Details: err.Error(),
		})
	}

	userID, err := h.repo.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, &ErrorResponse{
				Message: "指定されたEmailのユーザーは既に存在します。",
				Code:    "email_already_exists",
			})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "ユーザー作成中にエラーが発生しました。",
			Code:    "internal_server_error",
		})
	}

	res := CreateUserResponse{
		ID: userID,
	}

	return c.JSON(http.StatusCreated, res)
}

// GET /api/v1/users/:userID
func (h *Handler) GetUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{
			Message: "ユーザーIDの形式が不正です。",
			Code:    "invalid_user_id",
		})
	}

	user, err := h.repo.GetUser(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{
				Message: "指定されたユーザーが見つかりません。",
				Code:    "user_not_found",
			})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "ユーザー情報の取得中にエラーが発生しました。",
			Code:    "internal_server_error",
		})
	}

	res := GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.JSON(http.StatusOK, res)
}
