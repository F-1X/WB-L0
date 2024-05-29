package cache_test

import (
	"sync"
	"testing"
	"time"
	"wb/backend/internal/cache"
	"wb/backend/internal/config"
	"wb/backend/internal/domain/entity"

	"github.com/go-playground/assert"
)

func TestConcurrentSetGetCache(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration:   time.Second * 60,
		IntervalGC:   time.Second * 60,
		MaxItems:     1000,
		MaxItemSize:  1024 * 1024,     // 1MB
		MaxCacheSize: 1 * 1024 * 1024, // 1MB
		MaxKeySize:   500,             // 500 bytes
	}

	cache := cache.New(cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	cache.SetOrder(order, time.Minute)

	var wg sync.WaitGroup
	numRoutines := 10

	readFromCache := func() {
		defer wg.Done()

		for i := 0; i < 100; i++ {
			_, err := cache.GetOrder("order-1")
			if err != nil {
				t.Errorf("Failed to get order from cache: %v", err)
				return
			}
		}
	}

	newOrder := entity.Order{
		OrderUID: "order-2",
	}
	writeToCache := func() {
		defer wg.Done()

		for i := 0; i < 100; i++ {

			cache.SetOrder(newOrder, time.Minute)

		}
	}

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go readFromCache()
		wg.Add(1)
		go writeToCache()
	}

	wg.Wait()
}

func TestConcurrentSetCache(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration:   200,
		IntervalGC:   600,
		MaxItems:     1000,
		MaxItemSize:  1024 * 1024,     // 1MB
		MaxCacheSize: 1 * 1024 * 1024, // 1MB
		MaxKeySize:   500,             // 500 bytes
	}

	cache := cache.New(cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	cache.SetOrder(order, time.Minute)

	var wg sync.WaitGroup
	numRoutines := 10

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 100; i++ {
				newOrder := entity.Order{
					OrderUID: "order-2",
				}
				cache.SetOrder(newOrder, time.Minute)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetCache(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration:   200, // не важно какое значение
		IntervalGC:   600, 
	}

	cache := cache.New(cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	cache.SetOrder(order, time.Second*10)

	var wg sync.WaitGroup
	numRoutines := 10

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 100; i++ {
				result, err := cache.GetOrder(order.OrderUID)
				if err != nil {
					t.Errorf("Failed to get order in cache: %v", err)
					return
				}
				assert.Equal(t, order, result)
			}
		}()
	}

	wg.Wait()
}

// Ставим ГЦ на 2 секунды а время жизни ключа на 1 секунду. Ждем 2 секунды и ГЦ должен удалить ключ.
func TestGCInterval(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: time.Second * 2, // short live
		IntervalGC: time.Second * 1, // 2 second for GC interval
	}

	c := cache.New(cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	c.SetOrder(order, 0)

	// ждем пока ГЦ отработает
	time.Sleep(time.Second * 2)

	result, err := c.GetOrder(order.OrderUID)

	assert.Equal(t, cache.ErrExpired, err)
	assert.Equal(t, entity.Order{}, result)
}

func TestExpirationInSetOrder(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: 1, // по дефолту определим 1 секунду жизни
		IntervalGC: 5, // ставим побольше тут не важно
	}

	c := cache.New(cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	c.SetOrder(order, time.Duration(time.Millisecond))

	// ждем пока истечения срока жизни
	time.Sleep(time.Millisecond)

	result, err := c.GetOrder(order.OrderUID)

	// должно быть истечено временя жизни
	assert.Equal(t, cache.ErrExpired, err)
	assert.Equal(t, entity.Order{}, result)
}
