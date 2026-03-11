package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	orders, err := s.orders.FindAll(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}

	for idx := range orders {
		orders[idx].Memo = orders[idx].ID.String()
	}

	return orders, nil
}

func (s *Service) CreateOrder(
	ctx context.Context,
	address string,
	rawAmount uint64,
	customData json.RawMessage,
) (*models.Order, error) {
	o := models.Order{
		Address:    address,
		RawAmount:  rawAmount,
		CustomData: customData,
	}

	if err := s.orders.Create(ctx, &o); err != nil {
		return nil, fmt.Errorf("failed to create new order record: %w", err)
	}

	o.Memo = o.ID.String()
	return &o, nil
}

func (s *Service) CountOrdersBefore(ctx context.Context, before time.Time) (int64, error) {
	var c int64
	return c, s.db.WithContext(ctx).Model(&models.Order{}).Count(&c).Where("created_at < ?", before).Error
}

func (s *Service) DeleteOrdersBefore(ctx context.Context, before time.Time) error {
	return s.db.WithContext(ctx).Delete(&models.Order{}, "created_at < ?", before).Error
}
