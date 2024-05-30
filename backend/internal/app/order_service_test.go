package app_test

// import (
// 	"testing"
// 	"time"
// 	"wb/backend/internal/cache"
// 	"wb/backend/internal/config"
// 	"wb/backend/internal/domain/entity"

// 	"github.com/go-playground/assert"
// )
// func TestHandleHTTPReq(t *testing.T) {
// 	// Создаем тестовыйstan.Conn
// 	stanConn := &stan.Conn{}

// 	// Создаем тестовую реализацию OrderDB
// 	orderDB := &mockOrderDB{}

// 	// Создаем тестовую реализацию OrderCache
// 	orderCache := &mockOrderCache{}

// 	// Создаем экземпляр OrderService
// 	os := app.NewOrderService(orderDB, orderCache, stanConn)

// 	// Подписываемся на запросы от HTTP
// 	os.HandleHTTPReq()

// 	// Проверяем, что запрос обрабатывается корректно
// 	id := "test-id"
// 	msg := &nats.Msg{Data: []byte(id)}
// 	os.stan.NatsConn().Publish("request", msg.Data)
// 	time.Sleep(time.Millisecond) // Ждем обработку запроса
// }
