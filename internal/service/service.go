package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Makefolder/cynero/internal/db"
	"github.com/Makefolder/cynero/internal/models"
	"github.com/Makefolder/cynero/pkg/http"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log    *zap.SugaredLogger
	http   *http.Client
	db     *gorm.DB
	orders db.Repository[models.Order]
}

func New(log *zap.SugaredLogger, client *http.Client, gormDB *gorm.DB) *Service {
	ordersRepository := db.NewGormRepository[models.Order](gormDB)
	return &Service{
		log:    log,
		http:   client,
		db:     gormDB,
		orders: ordersRepository,
	}
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.db.WithContext(ctx).Raw("SELECT 1").Error
}

func (s *Service) FindAll(ctx context.Context) ([]models.Order, error) {
	return s.orders.FindAll(ctx)
}

func (s *Service) CreateOrder(
	ctx context.Context,
	rawAmount uint64,
	customData json.RawMessage,
) (*models.Order, error) {
	o := models.Order{
		Address:    "test.merchant_solana_address",
		Memo:       "test",
		RawAmount:  rawAmount,
		CustomData: customData,
	}

	if err := s.orders.Create(ctx, &o); err != nil {
		return nil, fmt.Errorf("failed to create new order record: %w", err)
	}

	return &o, nil
}
