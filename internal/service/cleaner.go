package service

import (
	"context"
	"time"

	"github.com/Makefolder/cynero/internal/models"
)

func (s *Service) DeleteOrdersBefore(ctx context.Context, before time.Time) (int64, error) {
	tx := s.db.WithContext(ctx).
		Where("created_at < ? AND paid_at IS NULL", before).
		Delete(&models.Order{})
	return tx.RowsAffected, tx.Error
}
