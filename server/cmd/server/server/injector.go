package server

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/handler"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type Server struct {
	Handler *handler.Handler
}

func Inject(db *sqlx.DB) *Server {
	repo := repository.New(db)

	cache := &handler.UserCache{}
	botToken := os.Getenv("BOT_TOKEN_KEY")
	if botToken != "" {
		if err := cache.Refresh(botToken); err != nil {
			log.Printf("UserCache: initial refresh failed: %v", err)
		}
	} else {
		log.Println("UserCache: BOT_TOKEN_KEY not set, cache will be empty")
	}

	h := handler.New(repo, cache)

	return &Server{
		Handler: h,
	}
}

func (d *Server) SetupRoutes(g *echo.Group) {
	d.Handler.SetupRoutes(g)
}
