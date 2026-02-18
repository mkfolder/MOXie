package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Find(ctx context.Context, id uuid.UUID) (*T, error) {
	var t T
	err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error
	return &t, err
}

func (r *GormRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var ts []T
	err := r.db.WithContext(ctx).Find(&ts).Error
	return ts, err
}

func (r *GormRepository[T]) Create(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *GormRepository[T]) Update(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(nil, "id = ?", id).Error
}
