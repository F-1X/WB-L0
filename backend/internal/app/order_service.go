package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb/backend/internal/domain/entity"
	"wb/backend/internal/domain/repository"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type OrderService struct {
	repo  repository.OrderDB
	cache repository.OrderCache
	stan  stan.Conn
}

func NewOrderService(repo repository.OrderDB, cache repository.OrderCache, stan stan.Conn) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
		stan:  stan,
	}
}

// HandleHTTPReq - подписка на запросы от HHTP
func (os *OrderService) HandleHTTPReq() {

	_, err := os.stan.NatsConn().Subscribe("request_http.subject", func(msg *nats.Msg) {
		
		id := string(msg.Data)
		log.Println("Received ID:", id)

		data, err := os.cache.GetOrder(id)
		if err == nil {
			orderData, _ := json.Marshal(data)
			if err := msg.Respond(orderData); err != nil {
				log.Println("err", err)
			}
			log.Println("get order from cache:", data.OrderUID)
			return
		} else {
			log.Println("er", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		data, err = os.repo.GetOrder(ctx, id)
		if err == nil {
			orderData, _ := json.Marshal(data)
			if err := msg.Respond(orderData); err != nil {
				log.Println("err", err)
			}
			go os.cache.SetOrder(data, 0)
			return
		}

		errorResponse := []byte(fmt.Sprintf(`{"error":"%s"}`, err))
		if err := msg.Respond(errorResponse); err != nil {
			log.Println("err", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}

// HandleNATSStreaming - подписка на стрим от источника заказов (пополнение в базу, новые заказы)
func (os *OrderService) HandleNATSStreaming() {
	validate := validator.New()

	_, err := os.stan.Subscribe("orders", func(msg *stan.Msg) {
		log.Println(string(msg.Data))
		var order entity.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println(err)
			return
		}

		if err := validate.Struct(order); err != nil {
			log.Println("Invalid order received:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		go os.cache.SetOrder(order, 0)

		if err := os.repo.InsertOrder(ctx, order); err != nil {
			log.Println("failed to insert order, reason:", err)
		}

		log.Println("handled order:", order.OrderUID, msg.Subject)

	}, stan.StartWithLastReceived(), stan.DurableName("durable-subscription"))
	// опции -> начала получения с последнего принятого сообщения и длительная подписка

	if err != nil {
		log.Fatal(err)
	}

}
