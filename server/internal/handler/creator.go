package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/pkg/config"
)

var traqHTTPClient = &http.Client{Timeout: 10 * time.Second}

// UserCache は traQ ID → UUID のインメモリキャッシュ
type UserCache struct {
	mu           sync.RWMutex
	traqIDToUUID map[string]uuid.UUID
}

type traqUser struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Bot  bool      `json:"bot"`
}

// Refresh は traQ API からユーザー一覧を取得してキャッシュを更新する
func (uc *UserCache) Refresh(botToken string) error {
	req, err := http.NewRequest("GET", "https://q.trap.jp/api/v3/users", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	q := req.URL.Query()
	q.Add("include-suspended", "true")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+botToken)

	resp, err := traqHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("traQ API returned %d", resp.StatusCode)
	}

	var users []traqUser
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return fmt.Errorf("decode users: %w", err)
	}

	newMap := make(map[string]uuid.UUID, len(users))
	for _, u := range users {
		if !u.Bot {
			newMap[u.Name] = u.ID
		}
	}

	uc.mu.Lock()
	uc.traqIDToUUID = newMap
	uc.mu.Unlock()

	log.Printf("UserCache: refreshed %d users", len(newMap))

	return nil
}

// GetUUID は traQ ID から UUID を返す
func (uc *UserCache) GetUUID(traqID string) (uuid.UUID, bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	id, ok := uc.traqIDToUUID[traqID]

	return id, ok
}

// Size はキャッシュに登録されているユーザー数を返す
func (uc *UserCache) Size() int {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	return len(uc.traqIDToUUID)
}

// RefreshUserCache は UserCache を traQ API から再取得する（cron から呼ばれる）
func (h *Handler) RefreshUserCache() {
	botToken := os.Getenv("BOT_TOKEN_KEY")
	if botToken == "" {
		log.Println("RefreshUserCache: BOT_TOKEN_KEY not set, skipping")

		return
	}
	if err := h.userCache.Refresh(botToken); err != nil {
		log.Printf("RefreshUserCache: failed: %v", err)
	}
}

// getUserID はリクエストから認証済みユーザーの UUID を取得する。
// 本番: X-Forwarded-User ヘッダー（NeoShowcase が付与）を使用。
// 開発: APP_ENV=development のとき DEV_USER 環境変数にフォールバック。
func (h *Handler) getUserID(c echo.Context) (uuid.UUID, error) {
	traqID := c.Request().Header.Get("X-Forwarded-User")
	if traqID == "" && config.IsDevelopment() {
		traqID = os.Getenv("DEV_USER")
	}
	if traqID == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	id, ok := h.userCache.GetUUID(traqID)
	if !ok {
		log.Printf("getUserID: unknown traQ ID %q (cache size=%d)", traqID, h.userCache.Size())

		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found in cache")
	}

	return id, nil
}
