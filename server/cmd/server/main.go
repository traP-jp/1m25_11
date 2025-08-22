package main

import (
	"github.com/traP-jp/1m25_11/server/cmd/server/server"
	"github.com/traP-jp/1m25_11/server/pkg/config"
	"github.com/traP-jp/1m25_11/server/pkg/database"
	"github.com/go-co-op/gocron/v2" 
 	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://1m25-11.trap.show", "http://localhost"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
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

	// Setup OpenAPI routes (at /api root)
	s.SetupAPIRoutes(e) 

	//gocron
	ss, er := gocron.NewScheduler()
	if er != nil{
		log.Fatal(er)
	}

	
	_, err := ss.NewJob(
		gocron.CronJob("0 0 * * *", false),
		gocron.NewTask(h.CronJobTask),
	)

	if err != nil{
		log.Fatal(err)
	}
	ss.Start()

	e.Logger.Fatal(e.Start(config.AppAddr()))

}
