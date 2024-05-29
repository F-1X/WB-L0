package cache

import (
	"context"
	"wb/backend/internal/domain/entity"
)

// WarmingCache - прогрев кеша
func (c *cache) WarmingCache(ctx context.Context, rows []entity.Order) {
	done := make(chan struct{})
	go func() {
		defer close(done)

		for _, row := range rows {
			select {
			case <-ctx.Done():
				return
			default:
				c.SetOrder(row, 0)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:
		return
	}

}
