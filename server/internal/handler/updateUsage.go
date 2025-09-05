package handler

import (
	"errors"
	"net/http"
	"os"
	"time"
)

func (h *Handler) getDailyUsage(since time.Time, until time.Time) error {
	botKey, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		return errors.New("BOT_TOKEN_KEY not found in environment variables")
	}

	const url = "https://q.trap.jp/api/v3/messages"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+botKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("API returned status code: " + resp.Status)
	}

	return nil
}
