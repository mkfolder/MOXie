package service

import (
	"context"
	"time"

	"github.com/Makefolder/cynero/internal/models"
)

func (s *Service) CountOrdersBefore(ctx context.Context, before time.Time) (int64, error) {
	var c int64
	return c, s.db.WithContext(ctx).Model(&models.Order{}).Count(&c).Where("created_at < ?", before).Error
}

func (s *Service) DeleteOrdersBefore(ctx context.Context, before time.Time) error {
	return s.db.WithContext(ctx).Delete(&models.Order{}, "created_at < ?", before).Error
}
