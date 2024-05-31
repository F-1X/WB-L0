package server

import (
	"net/http"
	"wb/backend/internal/app"

	"github.com/nats-io/stan.go"
)

type Handler struct {
	Mux          http.ServeMux
	orderService *app.OrderService
	stan         stan.Conn
}

func NewHandler(orderService *app.OrderService, stan stan.Conn, frontendPath string) *Handler {
	router := &Handler{
		stan:         stan,
		orderService: orderService,
	}

	router.Mux.Handle("GET /", http.FileServer(http.Dir(frontendPath)))
	router.Mux.HandleFunc("GET /order", router.OrderHandler)

	return router
}
