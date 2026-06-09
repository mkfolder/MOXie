package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/models"
)

func (s *Service) FindAll(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]models.Order, int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.Order{}).Where("merchant_id = ?", merchantID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	var orders []models.Order
	err := s.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find all orders: %w", err)
	}

	for idx := range orders {
		address := "not available"
		if orders[idx].Merchant.Address != nil {
			address = *orders[idx].Merchant.Address
		}
		orders[idx].Address = address
	}

	return orders, total, nil
}

func (s *Service) FindOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	order, err := s.orders.Find(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return order, nil
}

func (s *Service) CreateOrder(
	ctx context.Context,
	merchantID uuid.UUID,
	rawAmount uint64,
	customData json.RawMessage,
) (*models.Order, error) {
	merchant, err := s.merchants.Find(ctx, merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	if merchant.Address == nil {
		return nil, fmt.Errorf("merchant address is not specified")
	}

	address := *merchant.Address
	o := models.Order{
		MerchantID:         merchantID,
		Address:            address,
		RawRequestedAmount: rawAmount,
		CustomData:         customData,
	}

	if err := s.orders.Create(ctx, &o); err != nil {
		return nil, fmt.Errorf("failed to create new order record: %w", err)
	}

	return &o, nil
}
