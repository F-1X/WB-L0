package server

import (
	"net/http"
	"wb/backend/internal/app"

	"github.com/nats-io/stan.go"
)

type Router struct {
	Mux          http.ServeMux
	orderService app.OrderService
	stan         stan.Conn
}

func NewRouter(orderService app.OrderService, stan stan.Conn, frontendPath string) *Router {
	router := &Router{
		Mux:          *http.NewServeMux(),
		stan:         stan,
		orderService: orderService,
	}

	router.Mux.Handle("GET /", http.FileServer(http.Dir(frontendPath)))
	router.Mux.HandleFunc("POST /order", router.orderHandler)

	return router
}
