package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/traP-jp/1m25_11/server/internal/repository"
)

func (h *Handler) Test(ctx context.Context) {
	log.Println("Starting test")
}

func (h *Handler) CronJobTask(ctx context.Context) {

	bot_key, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		log.Println("BOT_TOKEN_KEY not found in environment variables")

		return
	}
	const url = "https://q.trap.jp/api/v3/stamps"
	log.Println("Starting scheduled job to fetch users...")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+bot_key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request to API:%v", err)

		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-200 status code: %d", resp.StatusCode)

		return
	}

	var apiResp []*repository.ResponseStamp
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("Error decoding response: %v", err)

		return
	}

	if err := h.repo.SaveStamp(ctx, apiResp); err != nil {
		log.Printf("Error saving stamps: %v", err)

		return
	}

	log.Println("successfully cronJobTask")

}
