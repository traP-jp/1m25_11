package handler

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)




func (h *Handler) CronJobTask(c echo.Context) {
	const token = "";
	const url = "https://q.trap.jp/api/v3/stamps"
	log.Println("Starting scheduled job to fetch users...")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer"+ token)
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

	var apiResp []*repository.StampResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil{
		log.Printf("Error decoding response: %v", err)

		return
	}

	if err := h.repo.UpdateStamp(c.Request().Context(), apiResp); err != nil {
		log.Printf("Error saving stamps: %v", err)

		return
	}


	log.Println("successfully fetched and update stamps")

}
