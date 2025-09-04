package repository

import (
	"errors"
	"net/http"
	"os"
)

func (r *Repository) updateUsage() error {
	botKey, ok := os.LookupEnv("BOT_TOKEN_KEY")
	if !ok {
		return errors.New("BOT_TOKEN_KEY not found in environment variables")
	}

	const url = "https://q.trap.jp/api/v3/stamps"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return nil
}
