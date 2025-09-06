package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type ResponseFromLlm struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	TagNames    []string  `json:"keywords"`
}
type TagInfo struct {
	Name string `json:"name" db:"name"`
}
type Tag struct {
	ID   uuid.UUID `json:"tag_id"`
	Name string    `json:"tag_name"`
}
type StampMetaAddition struct {
	ID          uuid.UUID   `json:"id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	Description string      `json:"description"`
}

func main() {
	bulk_tags_url := "http://localhost:8080/api/v1/bulk/tags"
	bulk_stamps_meta_url := "http://localhost:8080/api/v1/bulk/stamps-meta"
	tags_url := "http://localhost:8080/api/v1/tags"

	responseBites, err := os.ReadFile("response.json")
	if err != nil {
		log.Fatalf("failed to read response.json: %v", err)
	}
	var responses []ResponseFromLlm
	if err := json.Unmarshal(responseBites, &responses); err != nil {
		log.Fatalf("failed to unmarshal response.json: %v", err)
	}

	tagNamesSet := make(map[string]struct{})
	for _, resp := range responses {
		for _, tagName := range resp.TagNames {
			tagNamesSet[tagName] = struct{}{}
		}
	}
	tagInfos := []TagInfo{}
	for tagName := range tagNamesSet {
		tagInfos = append(tagInfos, TagInfo{Name: tagName})
	}
	tagBody, err := json.Marshal(tagInfos)
	if err != nil {
		log.Fatalf("failed to marshal tagInfos: %v", err)
	}
	resp, err := http.Post(bulk_tags_url, "application/json", bytes.NewReader(tagBody))
	if err != nil {
		log.Fatalf("failed to send bulk tags request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("unexpected status code from bulk tags request: %d", resp.StatusCode)
	}

	var tags []Tag
	resp, err = http.Get(tags_url)
	if err != nil {
		log.Fatalf("failed to get tags: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}
	fmt.Print(string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &tags); err != nil {
		log.Fatalf("failed to decode tags: %v", err)
	}

	tagNameToID := make(map[string]uuid.UUID)
	for _, tag := range tags {
		tagNameToID[tag.Name] = tag.ID
	}

	stampMetaAdditions := []StampMetaAddition{}
	for _, resp := range responses {
		tagIDs := []uuid.UUID{}
		for _, tagName := range resp.TagNames {
			if tagID, ok := tagNameToID[tagName]; ok {
				tagIDs = append(tagIDs, tagID)
			} else {
				log.Printf("warning: tag name %s not found in tagNameToID map", tagName)
			}
		}
		stampMetaAdditions = append(stampMetaAdditions, StampMetaAddition{
			ID:          resp.ID,
			TagIDs:      tagIDs,
			Description: resp.Description,
		})
	}

	stampMetaBody, err := json.Marshal(stampMetaAdditions)
	if err != nil {
		log.Fatalf("failed to marshal stampMetaAdditions: %v", err)
	}
	fmt.Print(string(stampMetaBody))
	resp, err = http.Post(bulk_stamps_meta_url, "application/json", bytes.NewReader(stampMetaBody))
	if err != nil {
		log.Fatalf("failed to send bulk stamps meta request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code from bulk stamps meta request: %d", resp.StatusCode)
	}
}
