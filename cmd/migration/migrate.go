package main

import (
	"os"

	"github.com/Makefolder/cynero/internal/config"
	"github.com/Makefolder/cynero/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPath string = "./config/default.yaml"
)

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = defaultPath
	}

	cfg, err := config.New(path)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Order{}); err != nil {
		panic(err)
	}
}
