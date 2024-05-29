package repository

import (
	"context"
	"wb/backend/internal/domain/entity"
)

//go:generate go run https://github.com/vektra/mockery/v2@v2.43.1
type OrderDB interface {
	GetOrder(ctx context.Context, id string) (entity.Order, error)
	// GetOrders(ctx context.Context) ([]entity.Order, error)
	InsertOrder(ctx context.Context, order entity.Order) error
	GetOrdersWithLimitByOrder(ctx context.Context, limit int, order string, direction string) ([]entity.Order, error)
}
