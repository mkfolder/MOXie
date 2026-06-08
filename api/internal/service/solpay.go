package service

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/pkg/soltxn"
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

	if order.Merchant.Address == nil {
		return nil, errors.New("merchant's address is not available")
	}

	recipient, err := solana.PublicKeyFromBase58(*order.Merchant.Address)
	if err != nil {
		s.log.Errorf("failed to parse recipient address: %w", err)
		return nil, err
	}

	recent, err := s.rpc.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := soltxn.GetNativeTransferBuilder().
		SetSender(sender).
		SetRecipient(recipient).
		SetAmount(order.RawRequestedAmount).
		SetMemo(order.Memo).
		Build(recent.Value.Blockhash)
	if err != nil {
		s.log.Errorf("failed to build transfer txn: %w", err)
		return nil, err
	}

	return tx, nil
}
