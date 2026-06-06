package main

import (
	"os"
	"time"

	"github.com/mkfolder/moxie/internal/config"
	"github.com/mkfolder/moxie/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	cfg, err := config.New(path)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Error),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Merchant{}, &models.Order{}); err != nil {
		panic(err)
	}
}
