package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/models"
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
