package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Makefolder/cynero/internal/config"
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/Makefolder/cynero/internal/log"
	"github.com/Makefolder/cynero/internal/routes"
	"github.com/Makefolder/cynero/internal/service"
	"github.com/Makefolder/cynero/pkg/http"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultPath string = "./config/default.yaml"
)

func main() {
	// Config init
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultPath
	}

	cfg, err := config.New(path)
	if err != nil {
		panic(err)
	}

	// Log init
	log, err := log.New(cfg.Server.Environment)
	if err != nil {
		panic(err)
	}

	// DB init
	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}

	// App init
	app := fiber.New(fiber.Config{
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "CyneroServer",
		AppName:       "Cynero",
	})
	api := app.Group("/api").Use(cors.New())

	httpClient := http.New(nil, nil, cfg.HTTP.Timeout)
	s := service.New(log, httpClient, db)
	h := handler.New(s)
	routes.Setup(api, h)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Infoln("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
