package server

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/handler"
	"github.com/traP-jp/1m25_11/server/internal/repository"
	"github.com/traP-jp/1m25_11/server/pkg/config"
)

type Server struct {
	Handler *handler.Handler
}

func Inject(db *sqlx.DB) *Server {
	repo := repository.New(db)

	cache := &handler.UserCache{}
	botToken := os.Getenv("BOT_TOKEN_KEY")
	if botToken == "" {
		if !config.IsDevelopment() {
			log.Fatal("UserCache: BOT_TOKEN_KEY is required in production")
		}
	} else if err := cache.Refresh(botToken); err != nil {
		if config.IsDevelopment() {
			log.Printf("UserCache: initial refresh failed: %v", err)
		} else {
			log.Fatalf("UserCache: initial refresh failed: %v", err)
		}
	}

	h := handler.New(repo, cache)

	return &Server{
		Handler: h,
	}
}

func (d *Server) SetupRoutes(g *echo.Group) {
	d.Handler.SetupRoutes(g)
}
