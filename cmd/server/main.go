package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Makefolder/cynero/internal/config"
	"github.com/Makefolder/cynero/internal/handler"
	"github.com/Makefolder/cynero/internal/helius"
	"github.com/Makefolder/cynero/internal/log"
	"github.com/Makefolder/cynero/internal/routes"
	"github.com/Makefolder/cynero/internal/service"
	"github.com/Makefolder/cynero/internal/workers"
	"github.com/Makefolder/cynero/pkg/http"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultPath string           = "./config/default.yaml"
	apiKey      string           = "569dbc7b-d716-4219-aa73-7580f65b3011"
	heliusNet   helius.HeliusNet = helius.HeliusNetMainnet
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
		Logger:  logger.Default.LogMode(logger.Error),
		NowFunc: func() time.Time { return time.Now().UTC() },
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
	hc := helius.NewClient(httpClient, apiKey, heliusNet)
	s := service.New(log, hc, db)
	h := handler.New(s)
	routes.Setup(api, h)

	go func() {
		if err := s.PingDB(context.Background()); err != nil {
			log.Fatal(err)
		}

		if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	// Workers init
	cw := workers.NewCleanerWorker(log, cfg.Workers.CleanerInterval, cfg.Workers.OrderExpiration, s)

	if err := cw.Start(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		cw.Stop()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Infoln("Gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
