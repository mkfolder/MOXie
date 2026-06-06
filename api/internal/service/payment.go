package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/helius"
	"github.com/mkfolder/moxie/internal/models"
	"github.com/mkfolder/moxie/pkg/http"
	"github.com/mr-tron/base58/base58"
)

func (s *Service) HandleWebhook(ctx context.Context, transacitons []helius.Transaction) {
	for _, tx := range transacitons {
		go s.processTransaction(&tx)
	}
}

func (s *Service) processTransaction(tx *helius.Transaction) {
	for _, instruction := range tx.Instructions {
		if instruction.ProgramID != memoProgramID {
			continue
		}

		// helius returns any data as base58
		b, err := base58.Decode(instruction.Data)
		if err != nil {
			s.log.Errorf("failed to decode base58 data of memo program: %v", err)
			break
		}

		// order id is encoded into base58 by this API
		orderID, err := base58.Decode(string(b))
		if err != nil {
			s.log.Errorf("failed to decode base58 order id: %v", err)
			break
		}

		uid, err := uuid.Parse(string(orderID))
		if err != nil {
			s.log.Errorf("failed to parse order uuid: %v", err)
			break
		}

		s.log.Debug("tx:\t%s\ndata:\t%s", tx.Signature, string(b))
		go s.processOrder(uid, tx)
		break
	}
}

func (s *Service) processOrder(orderID uuid.UUID, tx *helius.Transaction) {
	var paidAmount uint64

	order, err := s.orders.Find(context.Background(), orderID)
	if err != nil {
		s.log.Errorf("failed to find order with id %s: %v", orderID.String(), err)
		return
	}

	for _, account := range tx.AccountData {
		if account.Account != order.Address {
			continue
		}

		if account.NativeBalanceChange < 0 {
			s.log.Errorf(
				"transaction %s has negative balance change: %d",
				tx.Signature, account.NativeBalanceChange)
			return
		}

		paidAmount = uint64(account.NativeBalanceChange)
		break
	}

	// todo!: save transactions instead of fields in order
	// !      we can sum raw paid amount by querying all "TRANSFER" transactions
	if order.RawPaidAmount == nil {
		order.RawPaidAmount = new(uint64)
		*order.RawPaidAmount = paidAmount
	} else {
		*order.RawPaidAmount += paidAmount
	}

	paidAt := time.Unix(int64(tx.Timestamp), 0).UTC()
	order.PaidAt = &paidAt

	// todo!: save transactions instead of fields in order
	// !      since users may want to split payment
	order.TxHash = &tx.Signature

	if err := s.orders.Update(context.Background(), order); err != nil {
		s.log.Errorf("failed to update order with id %s: %v", orderID.String(), err)
	}

	if *order.RawPaidAmount < order.RawRequestedAmount {
		s.log.Debugf(
			"order payment %s has not been satisfied (expected at least %d, got %d)",
			tx.Signature, order.RawRequestedAmount, *order.RawPaidAmount)
		return
	}

	go s.notify(order)
}

// sends webhook request with order details to our merchants
func (s *Service) notify(order *models.Order) {
	merchant := order.Merchant

	if merchant.WebhookURL == nil {
		s.log.Infof("merchant %s has no webhook url", merchant.ID.String())
		return
	}

	ctx := context.Background()
	s.log.Debug("sending webhook to merchant %s", merchant.ID.String())

	res, err := s.http.Post(ctx, *merchant.WebhookURL, nil, order)
	if err != nil {
		s.log.Errorf("failed to send webhook to merchant %s: %v", merchant.ID.String(), err)
		return
	}

	if !http.IsOK(res) {
		s.log.Errorf("bad response from merchant %s: %s", merchant.ID.String(), res.Status)
		return
	}

	s.log.Infof("successfully notified merchant %s", merchant.ID.String())
}
