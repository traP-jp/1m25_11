package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
	"context"
	"encoding/json"
	"os"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type (
	Message struct {
		ID        uuid.UUID    `json:"id"`
		UserID    uuid.UUID    `json:"userId"`
		ChannelID uuid.UUID    `json:"channelId"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Pinned    bool      `json:"pinned"`
		Stamps    []Stamp   `json:"stamps"`
		ThreadID  uuid.UUID    `json:"threadId"`
	}
	Stamp struct {
		UserID    uuid.UUID    `json:"userId"`
		StampID   uuid.UUID   `json:"stampId"`
		Count     int       `json:"count"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

func (h *Handler) MonthlyCount(c echo.Context) error {
	after := "2025-08-05T09:00:00Z"
	before := "2025-09-05T09:00:00Z"
	url := fmt.Sprintf("https://q.trap.jp/api/v3/stamps?after=%s&before=%s&limit=100&sort=createdAt", after, before)
	message, err := h.fetchMessage(c, url)
	if err != nil {
		log.Printf("error fetching messages from API:%v", err)
		return c.String(http.StatusInternalServerError, "error fetching messages from API")
	}
	log.Println("successfully fetched messages from API")

	stampCount := h.reactionStampCount(c, message)
	log.Println("successfully counted reaction stamps")
	err = h.repo.UpdateMonthlyCount(c.Request().Context(), stampCount)
	if err != nil {
		log.Printf("error updating monthly count in database:%v", err)
		return c.String(http.StatusInternalServerError, "error updating monthly count in database")
	}

	return c.String(http.StatusOK, "pong")
}

func (h *Handler) fetchMessage(c echo.Context, url string) ([]Message, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	bot_key, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		log.Println("BOT_TOKEN_KEY not found in environment variables")

		return nil, fmt.Errorf("BOT_TOKEN_KEY not found in environment variables")
	}
	req.Header.Set("Authorization", "Bearer "+bot_key)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request to API:%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching messages: %v", resp.Status)
		return nil, fmt.Errorf("error fetching messages: %v", resp.Status)
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	return messages, nil
}
func (h *Handler) reactionStampCount(c echo.Context, message []Message) map[uuid.UUID]int {
	return nil
}
