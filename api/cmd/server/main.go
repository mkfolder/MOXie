package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/mkfolder/moxie/internal/common"
	"github.com/mkfolder/moxie/internal/config"
	"github.com/mkfolder/moxie/internal/handler"
	"github.com/mkfolder/moxie/internal/helius"
	"github.com/mkfolder/moxie/internal/log"
	"github.com/mkfolder/moxie/internal/routes"
	"github.com/mkfolder/moxie/internal/service"
	"github.com/mkfolder/moxie/internal/workers"
	"github.com/mkfolder/moxie/pkg/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Config init
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		panic("HELIUS_API_KEY environment variable not set")
	}

	cfg, err := config.New(path)
	if err != nil {
		panic(err)
	}

	webhookURL := cfg.Server.WebhookURL
	heliusNet := helius.HeliusNetMainnet
	if cfg.Server.Environment == common.EnvironmentDevelopment {
		heliusNet = helius.HeliusNetDevnet
	}

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

	api := app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4321", "http://localhost"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	httpClient := http.New(nil, nil, cfg.HTTP.Timeout)

	params := service.NewServiceParams{
		Log:    log,
		HTTP:   httpClient,
		GormDB: gormDB,
		RPC:    rpc.DevNet_RPC,
		Auth: service.AuthConfig{
			JWTSecret:       cfg.Auth.JWTSecret,
			AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
			RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
		},
	}

	s := service.New(params)
	h := handler.New(s, cfg.Auth)
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
