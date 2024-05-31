package repository

import (
	"context"
	"wb/backend/internal/domain/entity"
)

type OrderDB interface {
	GetOrder(ctx context.Context, id string) (entity.Order, error)
	InsertOrder(ctx context.Context, order entity.Order) error
	GetOrdersWithLimitByOrder(ctx context.Context, limit int, order string, direction string) ([]entity.Order, error)
}
