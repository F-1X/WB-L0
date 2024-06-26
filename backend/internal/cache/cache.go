package cache

import (
	"context"
	"log"
	"sync"
	"time"
	"wb/backend/internal/config"
	"wb/backend/internal/domain/entity"
	"wb/backend/internal/domain/repository"
)

type cache struct {
	expiration time.Duration // дефолтное значение для хранения ключа
	intervalGC time.Duration // время афк перед обновлением гц
	memory     sync.Map
	cancelGC   context.CancelFunc
}

type Item struct {
	Value      entity.Order
	Expiration int64
	Duration   time.Duration // доп. время хранения ключа. нужно для обновления время жизни
}

func New(ctx context.Context, cfg config.CacheConfig) repository.OrderCache {
	if cfg.Expiration == 0 {
		cfg.Expiration = time.Minute
	}

	cache := cache{
		expiration: cfg.Expiration,
		intervalGC: cfg.IntervalGC,
	}

	if cfg.IntervalGC > 0 {
		go cache.GC()
	}

	go func() {
		<-ctx.Done()
		if cache.cancelGC != nil { // если GC не был запущен
			cache.cancelGC()
		}

		log.Println("[!] cache shutdown")
	}()

	return &cache
}

func (c *cache) GetOrder(key string) (entity.Order, error) {
	value, found := c.memory.Load(key)
	if !found {
		return entity.Order{}, ErrNotFound
	}

	item := value.(Item)
	if time.Now().UnixNano() > item.Expiration {
		c.Delete(key)
		return entity.Order{}, ErrExpired
	}

	c.UpdateOrder(key, Item{
		Value:      item.Value,
		Expiration: time.Now().Add(item.Duration).UnixNano(),
		Duration:   item.Duration,
	})

	return item.Value, nil
}

func (c *cache) SetOrder(order entity.Order, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = c.expiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	item := Item{
		Value:      order,
		Expiration: expiration,
		Duration:   duration,
	}

	c.memory.Store(order.OrderUID, item)

}

func (c *cache) UpdateOrder(orderUID string, item Item) error {
	c.memory.Store(orderUID, item)
	return nil
}

func (c *cache) Delete(key string) error {
	if _, found := c.memory.Load(key); !found {
		return ErrNotFound
	}

	c.memory.Delete(key)
	return nil
}

func (c *cache) GC() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelGC = cancel
	go func() {
		ticker := time.NewTicker(c.intervalGC)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("GC stopped")
				return
			case <-ticker.C:
				log.Println("GC started")
				c.cleanUp()
			}
		}
	}()
}

func (c *cache) cleanUp() {
	c.memory.Range(func(k, v interface{}) bool {
		item := v.(Item)
		if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
			log.Println("GC delete key:", k.(string))
			c.memory.Delete(k)
		}
		return true
	})
}

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
