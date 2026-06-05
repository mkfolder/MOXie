package service

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
)

func (s *Service) BuildOrderTransaction(
	ctx context.Context,
	sender solana.PublicKey,
	orderID uuid.UUID,
) (*solana.Transaction, error) {
	order, err := s.orders.Find(ctx, orderID)
	if err != nil {
		return nil, err
	}

	recipient, err := solana.PublicKeyFromBase58(order.Merchant.Address)
	if err != nil {
		s.log.Errorf("failed to parse recipient address: %w", err)
		return nil, err
	}

	tx, err := s.buildTransferTxn(ctx, sender, recipient, order.RawRequestedAmount, &order.Memo)
	if err != nil {
		s.log.Errorf("failed to build transfer txn: %w", err)
		return nil, err
	}

	return tx, nil
}
