package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Makefolder/cynero/internal/db"
	"github.com/Makefolder/cynero/internal/helius"
	"github.com/Makefolder/cynero/internal/models"
	"github.com/mr-tron/base58/base58"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const memoProgramID = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"

type Service struct {
	log    *zap.SugaredLogger
	db     *gorm.DB
	hc     *helius.HeliusClient
	orders db.Repository[models.Order]
}

func New(log *zap.SugaredLogger, hc *helius.HeliusClient, gormDB *gorm.DB) *Service {
	ordersRepository := db.NewGormRepository[models.Order](gormDB)
	return &Service{
		log:    log,
		db:     gormDB,
		orders: ordersRepository,
		hc:     hc,
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

func (s *Service) CreateWebhook(ctx context.Context, webhookURL *url.URL, addresses []string) error {
	return s.hc.CreateWebhook(ctx, webhookURL, addresses)
}

func (s *Service) HandleWebhook(ctx context.Context, body []byte) {
	var transacitons []helius.Transaction
	s.log.Debugf("request body: %s", string(body))

	if err := json.Unmarshal(body, &transacitons); err != nil {
		s.log.Errorf("failed to unmarshal body: %v", err)
	}

	for _, tx := range transacitons {
		for _, instruction := range tx.Instructions {
			if instruction.ProgramID == memoProgramID {
				b, err := base58.Decode(instruction.Data)

				if err != nil {
					s.log.Errorf("failed to decode base58 data of memo program: %v", err)
					break
				}

				s.log.Debug("tx:\t%s\ndata:\t%s", tx.Signature, string(b))
				break
			}
		}
	}
}
