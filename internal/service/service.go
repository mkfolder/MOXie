package service

import (
	"context"
	"log"
	"net/url"

	"github.com/Makefolder/cynero/internal/db"
	"github.com/Makefolder/cynero/internal/helius"
	"github.com/Makefolder/cynero/internal/models"
	"github.com/Makefolder/cynero/pkg/http"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const memoProgramID = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"

type Service struct {
	log        *zap.SugaredLogger
	db         *gorm.DB
	hc         *helius.HeliusClient
	http       *http.Client
	webhookURL *url.URL
	orders     db.Repository[models.Order]
	merchants  db.Repository[models.Merchant]
}

type NewServiceParams struct {
	Log        *zap.SugaredLogger
	WebhookURL string
	HC         *helius.HeliusClient
	HTTP       *http.Client
	GormDB     *gorm.DB
}

func New(params NewServiceParams) *Service {
	if params.Log == nil {
		panic("invalid service.New arguments: logger is nil")
	}

	if params.HC == nil {
		panic("invalid service.New arguments: helius client is nil")
	}

	if params.GormDB == nil {
		panic("invalid service.New arguments: gorm db is nil")
	}

	if params.WebhookURL == "" {
		log.Fatal("invalid service.New arguments: webhook url is not set")
	}

	webhook, err := url.Parse(params.WebhookURL)
	if err != nil {
		log.Fatalf("failed to parse webhook url: %v", err)
	}

	ordersRepository := db.NewGormRepository[models.Order](params.GormDB)
	merchantsRepository := db.NewGormRepository[models.Merchant](params.GormDB)
	return &Service{
		log:        params.Log,
		hc:         params.HC,
		http:       params.HTTP,
		db:         params.GormDB,
		webhookURL: webhook,
		orders:     ordersRepository,
		merchants:  merchantsRepository,
	}
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.db.WithContext(ctx).Raw("SELECT 1").Error
}
