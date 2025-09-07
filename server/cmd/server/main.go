package main

import (
	"context"
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traP-jp/1m25_11/server/cmd/server/server"
	"github.com/traP-jp/1m25_11/server/pkg/config"
	"github.com/traP-jp/1m25_11/server/pkg/database"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Dynamic CORS (credentials allowed) based on ALLOWED_ORIGINS env
	allowed := config.AllowedOrigins()
	e.Logger.Infof("CORS allowed origins: %v", allowed)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowed,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS, echo.HEAD},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           600,
	}))

	// connect to and migrate database
	db, err := database.Setup(config.MySQL())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	s := server.Inject(db)

	// Setup existing v1 API routes
	v1API := e.Group("/api/v1")
	s.SetupRoutes(v1API)

	//gocron
	ss, er := gocron.NewScheduler()
	if er != nil {
		log.Fatal(er)
	}

	_, err = ss.NewJob(
		gocron.CronJob("0 19 * * *", false),
		gocron.NewTask(s.Handler.CronJobTask, context.Background()),
	)

	if err != nil {
		log.Fatal(err)
	}
	ss.Start()

	e.Logger.Fatal(e.Start(config.AppAddr()))

}
