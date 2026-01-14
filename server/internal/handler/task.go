package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/traP-jp/1m25_11/server/internal/repository"
	"log"
	"net/http"
	"os"
)
type StampStatus struct {  
    TotalCount int `json:"totalCount"`  
    Count      int `json:"count"`  
}  

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

	stampTotalCount := make(map[uuid.UUID]int)
	rawCount := make(map[uuid.UUID]int)
	allStamps, err := h.repo.GetStampSummaries(ctx)
	if err != nil {
		log.Printf("Error retrieving all stamps: %v", err)
		return
	}
	log.Print("Retrieved all stamps, starting to fetch stats...")
	for _, stamp := range allStamps {

		var statsData StampStatus

		stampID := stamp.ID
		statsURL := fmt.Sprintf("https://q.trap.jp/api/v3/stamps/%s/stats", stampID)

		statsReq, err := http.NewRequestWithContext(ctx, "GET", statsURL, nil)
		if err != nil {
			log.Printf("Error creating request for stamp stats (%s): %v", stampID, err)
			return
		}
		statsReq.Header.Set("accept", "application/json")
		statsReq.Header.Set("Authorization", "Bearer "+bot_key)

		statsResp, err := client.Do(statsReq)
		if err != nil {
			log.Printf("Error sending request for stamp stats (%s): %v", stampID, err)
			return
		}
		defer statsResp.Body.Close()
		if statsResp.StatusCode != http.StatusOK {
			stampTotalCount[stamp.ID] = 0
			rawCount[stamp.ID] = 0
			log.Printf("Non-200 status code for stamp stats (%s): %d", stampID, statsResp.StatusCode)
			continue
		} else {
			if err := json.NewDecoder(statsResp.Body).Decode(&statsData); err != nil {
				log.Printf("Error decoding stats response for stamp (%s): %v", stampID, err)
				return
			}
			stampTotalCount[stamp.ID] = statsData.TotalCount
			rawCount[stamp.ID] = statsData.Count
		}

	}
	log.Print("Fetched all stamp stats, starting to update database...")
	if err := h.repo.UpdateCount(ctx, stampTotalCount, rawCount); err != nil {
		log.Printf("Error updating total count for stamps: %v", err)
		return
	}
	log.Println("Successfully updated total counts for all stamps")
	return

}
