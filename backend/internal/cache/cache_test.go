package cache_test

import (
	"context"
	"sync"
	"testing"
	"time"
	"wb/backend/internal/cache"
	"wb/backend/internal/config"
	"wb/backend/internal/domain/entity"

	"github.com/go-playground/assert"
)

// проверка время жизни ключа заданному по SetOrder
func TestGetOrder(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: time.Second * 1,
		IntervalGC: time.Second * 1,
	}

	c := cache.New(context.Background(), cfg)
	order := entity.Order{OrderUID: "order-1"}
	c.SetOrder(order, time.Duration(time.Millisecond))

	time.Sleep(time.Millisecond * 1)

	result, err := c.GetOrder(order.OrderUID)

	assert.Equal(t, cache.ErrExpired, err)
	assert.Equal(t, entity.Order{}, result)
}

// конкурентный доступ к кешу (order1 только чтение, order2 только запись, order3 чтение и запись)
func TestConcurrentSetGetCache(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: time.Second * 60,
		IntervalGC: time.Second * 60,
	}

	cache := cache.New(context.Background(), cfg)
	order1 := entity.Order{
		OrderUID: "order-1",
	}

	cache.SetOrder(order1, time.Minute)

	var wg sync.WaitGroup

	readFromCache := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			result, err := cache.GetOrder("order-1")
			assert.Equal(t, nil, err)
			assert.Equal(t, order1, result)
		}
	}

	order2 := entity.Order{
		OrderUID: "order-2",
	}
	writeToCache := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			cache.SetOrder(order2, time.Minute)
		}
	}

	order3 := entity.Order{
		OrderUID: "order-3",
	}

	readOrder3 := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			result, err := cache.GetOrder("order-3")
			assert.Equal(t, nil, err)
			assert.Equal(t, order3, result)
		}

	}

	writeOrder3 := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			cache.SetOrder(order3, time.Minute)
		}

	}

	numRoutines := 10
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go readFromCache()
		wg.Add(1)
		go writeToCache()
		wg.Add(1)
		go readOrder3()
		wg.Add(1)
		go writeOrder3()
	}

	wg.Wait()
}

func TestConcurrentSetCache(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: 200,
		IntervalGC: 600,
	}

	cache := cache.New(context.Background(), cfg)
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
		Expiration: 200, // не важно какое значение
		IntervalGC: 600,
	}

	cache := cache.New(context.Background(), cfg)
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

				assert.Equal(t, nil, err)
				assert.Equal(t, order, result)
			}
		}()
	}

	wg.Wait()
}

// Время ГЦ должно быть больше время жизни ключа. Ждем пока ГЦ отработает и проверяем что ключ не найден
func TestGCIntervalKeyNotFound(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: time.Millisecond * 40,
		IntervalGC: time.Millisecond * 150,
	}

	c := cache.New(context.Background(), cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	c.SetOrder(order, 0)

	// ждем пока ГЦ отработает и истечет время жизни ключа
	time.Sleep(time.Millisecond * 150)

	result, err := c.GetOrder(order.OrderUID)

	assert.Equal(t, cache.ErrNotFound, err)
	assert.Equal(t, entity.Order{}, result)
}

// Проверяем что ключ просрочен (истекло время жизни) а потом не найден.
func TestGCNotFoundKey(t *testing.T) {
	cfg := config.CacheConfig{
		Expiration: time.Second,
		IntervalGC: time.Second,
	}

	c := cache.New(context.Background(), cfg)
	order := entity.Order{
		OrderUID: "order-1",
	}

	c.SetOrder(order, time.Duration(time.Millisecond*100))

	time.Sleep(time.Millisecond * 100)

	// сперва выведется что ключ просрочен
	result, err := c.GetOrder(order.OrderUID)

	assert.Equal(t, cache.ErrExpired, err)
	assert.Equal(t, entity.Order{}, result)

	// затем он удалится, и не будет найден.
	result, err = c.GetOrder(order.OrderUID)

	assert.Equal(t, cache.ErrNotFound, err)
	assert.Equal(t, entity.Order{}, result)
}
