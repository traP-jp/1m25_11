package handler

import (
	"encoding/json"
	"fmt"
	"log"
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
		State       int      `json:"state"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	ResponseUser struct {
		ID          uuid.UUID `json:"user_id"`
		Name        string    `json:"traq_id"`
		DisplayName string    `json:"user_display_name"`
		IconFileID  uuid.UUID `json:"user_icon_file_id"`
		State       int      `json:"user_state"`
	}
)

func (h *Handler) getUsersList(c echo.Context) error {

	bot_key, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		log.Println("BOT_TOKEN_KEY not found in environment variables")

		return fmt.Errorf("BOT_TOKEN_KEY not found in environment variables")
	}
	const url = "https://q.trap.jp/api/v3/users"
	log.Println("Starting get user list")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)

		return fmt.Errorf("error creating request: %w", err)
	}
	q := req.URL.Query()
	q.Add("include-suspended", "true")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+bot_key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request to API:%v", err)

		return fmt.Errorf("error sending request to API:%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-200 status code: %d", resp.StatusCode)

		return fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}
	var users []*User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		log.Printf("error decoding response body: %v", err)

		return fmt.Errorf("error decoding response body: %w", err)
	}
	var resUsers []*ResponseUser
	for _, s := range users {
		if !s.Bot {
			resUsers = append(resUsers, &ResponseUser{
				ID:          s.ID,
				Name:        s.Name,
				DisplayName: s.DisplayName,
				IconFileID:  s.IconFileID,
				State:       s.State,
			})
		}
	}

	return c.JSON(http.StatusOK, resUsers)
}
