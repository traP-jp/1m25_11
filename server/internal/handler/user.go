package handler

import (
	"net/http"


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

	
)



func (h *Handler) GetUser(c echo.Context) error {
	userID:= uuid.Nil

	user, err := h.repo.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	// stampsUserOwned := make([]StampSummary, len(user.StampsUserOwned))
	// for i, s := range user.StampsUserOwned {
	// 	stampsUserOwned[i] = StampSummary{
	// 		Id:     s.ID,
	// 		Name:   s.Name,
	// 		FileId: s.FileID,
	// 	}
	// }

	// tagsUserCreated := make([]TagSummary, len(user.TagsUserCreated))
	// for i, t := range user.TagsUserCreated {
	// 	tagsUserCreated[i] = TagSummary{
	// 		Id:   t.ID,
	// 		Name: t.Name,
	// 	}
	// }

	// stampsUserTagged := make([]struct {
	// 	Stamp StampSummary `json:"stamp"`
	// 	Tag   TagSummary   `json:"tag"`
	// }, len(user.StampsUserTagged))
	// for i, st := range user.StampsUserTagged {
	// 	stampsUserTagged[i] = struct {
	// 		Stamp StampSummary `json:"stamp"`
	// 		Tag   TagSummary   `json:"tag"`
	// 	}{
	// 		Stamp: StampSummary{
	// 			Id:     st.Stamp.ID,
	// 			Name:   st.Stamp.Name,
	// 			FileId: st.Stamp.FileID,
	// 		},
	// 		Tag: TagSummary{
	// 			Id:   st.Tag.ID,
	// 			Name: st.Tag.Name,
	// 		},
	// 	}
	// }

	// descriptionsUserCreated := make([]struct {
	// 	Stamp         StampSummary `json:"stamp"`
	// 	DescriptionID uuid.UUID    `json:"description_id"`
	// }, len(user.DescriptionsUserCreated))
	// for i, d := range user.DescriptionsUserCreated {
	// 	descriptionsUserCreated[i] = struct {
	// 		Stamp         StampSummary `json:"stamp"`
	// 		DescriptionID uuid.UUID    `json:"description_id"`
	// 	}{
	// 		Stamp: StampSummary{
	// 			Id:     d.Stamp.ID,
	// 			Name:   d.Stamp.Name,
	// 			FileId: d.Stamp.FileID,
	// 		},
	// 		DescriptionID: d.DescriptionID,
	// 	}
	// }

	// res := GetUserResponse{
	// 	ID:                      user.ID,
	// 	IsAdmin:                 user.IsAdmin,
	// 	StampsUserOwned:         stampsUserOwned,
	// 	TagsUserCreated:         tagsUserCreated,
	// 	StampsUserTagged:        stampsUserTagged,
	// 	DescriptionsUserCreated: descriptionsUserCreated,
	// }

	return c.JSON(http.StatusOK, user)
}
