package handler

import (
	"fmt"
	"net/http"

	"github.com/traP-jp/1m25_11/server/internal/repository"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	GetUsersResponse []GetUserResponse

	GetUserResponse struct {
		ID               uuid.UUID      `json:"user_id"`
		IsAdmin          bool           `json:"is_admin"`
		StampsUserOwned  []StampSummary `json:"stamps_user_owned"`
		TagsUserCreated  []TagSummary   `json:"tags_user_created"`
		StampsUserTagged []struct {
			Stamp StampSummary `json:"stamp"`
			Tag   TagSummary   `json:"tag"`
		} `json:"stamps_user_tagged"`
		DescriptionsUserCreated []struct {
			Stamp         StampSummary `json:"stamp"`
			DescriptionID uuid.UUID    `json:"description_id"`
		} `json:"descriptions_user_created"`
	}

	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateUserResponse struct {
		ID uuid.UUID `json:"id"`
	}
)

func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.repo.GetUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	res := make(GetUsersResponse, len(users))
	for i, user := range users {
		res[i] = GetUserResponse{
			ID: user.ID,
			// Name:  user.Name,
			// Email: user.Email,
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateUser(c echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err)).SetInternal(err)
	}

	userID, err := h.repo.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	res := CreateUserResponse{
		ID: userID,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userID").SetInternal(err)
	}

	user, err := h.repo.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	stampsUserOwned := make([]StampSummary, len(user.StampsUserOwned))
	for i, s := range user.StampsUserOwned {
		stampsUserOwned[i] = StampSummary{
			Id:     s.ID,
			Name:   s.Name,
			FileId: s.FileID,
		}
	}

	tagsUserCreated := make([]TagSummary, len(user.TagsUserCreated))
	for i, t := range user.TagsUserCreated {
		tagsUserCreated[i] = TagSummary{
			Id:   t.ID,
			Name: t.Name,
		}
	}

	stampsUserTagged := make([]struct {
		Stamp StampSummary `json:"stamp"`
		Tag   TagSummary   `json:"tag"`
	}, len(user.StampsUserTagged))
	for i, st := range user.StampsUserTagged {
		stampsUserTagged[i] = struct {
			Stamp StampSummary `json:"stamp"`
			Tag   TagSummary   `json:"tag"`
		}{
			Stamp: StampSummary{
				Id:     st.Stamp.ID,
				Name:   st.Stamp.Name,
				FileId: st.Stamp.FileID,
			},
			Tag: TagSummary{
				Id:   st.Tag.ID,
				Name: st.Tag.Name,
			},
		}
	}

	descriptionsUserCreated := make([]struct {
		Stamp         StampSummary `json:"stamp"`
		DescriptionID uuid.UUID    `json:"description_id"`
	}, len(user.DescriptionsUserCreated))
	for i, d := range user.DescriptionsUserCreated {
		descriptionsUserCreated[i] = struct {
			Stamp         StampSummary `json:"stamp"`
			DescriptionID uuid.UUID    `json:"description_id"`
		}{
			Stamp: StampSummary{
				Id:     d.Stamp.ID,
				Name:   d.Stamp.Name,
				FileId: d.Stamp.FileID,
			},
			DescriptionID: d.DescriptionID,
		}
	}

	res := GetUserResponse{
		ID:                      user.ID,
		IsAdmin:                 user.IsAdmin,
		StampsUserOwned:         stampsUserOwned,
		TagsUserCreated:         tagsUserCreated,
		StampsUserTagged:        stampsUserTagged,
		DescriptionsUserCreated: descriptionsUserCreated,
	}

	return c.JSON(http.StatusOK, res)
}
