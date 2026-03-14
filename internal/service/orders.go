package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Makefolder/cynero/internal/models"
	"github.com/google/uuid"
)

func (s *Service) FindAll(ctx context.Context) ([]models.Order, error) {
	orders, err := s.orders.FindAll(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}

	for idx := range orders {
		orders[idx].Address = orders[idx].Merchant.Address
	}

	return orders, nil
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

	o := models.Order{
		MerchantID:         merchantID,
		Address:            merchant.Address,
		RawRequestedAmount: rawAmount,
		CustomData:         customData,
	}

	if err := s.orders.Create(ctx, &o); err != nil {
		return nil, fmt.Errorf("failed to create new order record: %w", err)
	}

	return &o, nil
}
