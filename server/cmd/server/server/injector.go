package server

import (
	"github.com/traP-jp/1m25_11/server/internal/handler"
	"github.com/traP-jp/1m25_11/server/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

)

type Server struct {
	Handler *handler.Handler
}

func Inject(db *sqlx.DB) *Server {
	repo := repository.New(db)
	h := handler.New(repo)
	
	return &Server{
		Handler: h,
	}
}

func (d *Server) SetupRoutes(g *echo.Group) {
	// TODO: handler.SetupRoutesを呼び出す or 直接書く？
	d.Handler.SetupRoutes(g)
}


