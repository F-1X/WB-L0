package repository

import (
	"context"
	"time"
	"wb/backend/internal/domain/entity"
)

type OrderCache interface {
	GetOrder(key string) (entity.Order, error)
	SetOrder(order entity.Order, duration time.Duration)
	WarmingCache(ctx context.Context, rows []entity.Order)
}
