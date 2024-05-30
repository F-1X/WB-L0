package server

import (
	"context"
	"log"
	"net/http"
	"wb/backend/internal/config"
)

type server struct {
	srv http.Server
}

func NewServer(router *http.ServeMux, cfg config.HTTPServerConfig) *server {
	return &server{
		srv: http.Server{
			Addr:        cfg.Addr,
			Handler:     router,
			ReadTimeout: cfg.Timeout,
			IdleTimeout: cfg.Idle_timeout,
			// TLSConfig:         &tls.Config{},
			// ReadHeaderTimeout: 5 * time.Second,
			// WriteTimeout:      10 * time.Second,
			// MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
			// TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			// ConnState:         func(conn net.Conn, state http.ConnState) {},
			// ErrorLog:          log.Default(),
			// BaseContext:       func(listener net.Listener) context.Context { return context.Background() },
			// ConnContext:       func(ctx context.Context, c net.Conn) context.Context { return ctx },
		},
	}
}

func (s *server) Run() {

	log.Println("[+] HTTP server started on addr: ", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.srv.Addr, err)
	}

}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
