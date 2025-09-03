package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	User struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		DisplayName string    `json:"displayName"`
		IconFileID  uuid.UUID `json:"iconFileId"`
		Bot         bool      `json:"bot"`
		State       int       `json:"state"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	ResponseUser struct {
		ID          uuid.UUID `json:"user_id"`
		Name        string    `json:"traq_id"`
		DisplayName string    `json:"user_display_name"`
		IconFileID  uuid.UUID `json:"user_icon_file_id"`
	}
)

func (h *Handler) getUsersList(c echo.Context) error {

	bot_key, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		err := fmt.Errorf("BOT_TOKEN_KEY not found in environment variables")
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "サーバーの設定に問題があります。",
			Code:    "server_configuration_error",
		})
	}

	const url = "https://q.trap.jp/api/v3/users"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "外部APIへのリクエスト作成に失敗しました。",
			Code:    "internal_server_error",
		})
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+bot_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "外部APIへの接続に失敗しました。",
			Code:    "external_api_connection_failed",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("traQ API returned non-200 status code: %d", resp.StatusCode)
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "外部APIから予期せぬ応答がありました。",
			Code:    "external_api_error",
		})
	}

	var users []*User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{
			Message: "外部APIの応答の解析に失敗しました。",
			Code:    "response_decode_error",
		})
	}

	var resUsers []*ResponseUser
	for _, s := range users {
		if !s.Bot {
			resUsers = append(resUsers, &ResponseUser{
				ID:          s.ID,
				Name:        s.Name,
				DisplayName: s.DisplayName,
				IconFileID:  s.IconFileID,
			})
		}
	}

	return c.JSON(http.StatusOK, resUsers)
}
