package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Makefolder/moxie/internal/config"
	"github.com/Makefolder/moxie/internal/handler"
	"github.com/Makefolder/moxie/internal/helius"
	"github.com/Makefolder/moxie/internal/log"
	"github.com/Makefolder/moxie/internal/routes"
	"github.com/Makefolder/moxie/internal/service"
	"github.com/Makefolder/moxie/internal/workers"
	"github.com/Makefolder/moxie/pkg/http"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	apiKey    string           = "569dbc7b-d716-4219-aa73-7580f65b3011"
	heliusNet helius.HeliusNet = helius.HeliusNetMainnet
)

func main() {
	// Config init
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	cfg, err := config.New(path)
	if err != nil {
		panic(err)
	}

	webhookURL := cfg.Server.WebhookURL

	// Log init
	log, err := log.New(cfg.Server.Environment)
	if err != nil {
		panic(err)
	}

	// DB init
	gormDB, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{
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
		ServerHeader:  "MOXieServer",
		AppName:       "MOXie",
	})
	api := app.Use(cors.New())

	httpClient := http.New(nil, nil, cfg.HTTP.Timeout)
	hc := helius.NewClient(httpClient, webhookURL, apiKey, heliusNet)

	params := service.NewServiceParams{
		Log:    log,
		HTTP:   httpClient,
		HC:     hc,
		GormDB: gormDB,
	}

	s := service.New(params)
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
