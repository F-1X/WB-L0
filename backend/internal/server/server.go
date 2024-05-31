package server

import (
	"context"
	"log"
	"net/http"
	"wb/backend/internal/config"
)

type Server struct {
	srv    http.Server
	router *http.ServeMux
}

func NewServer(router *http.ServeMux, cfg config.HTTPServerConfig) *Server {
	return &Server{
		srv: http.Server{
			Addr:        cfg.Addr,
			Handler:     router,
			ReadTimeout: cfg.Timeout,
			IdleTimeout: cfg.Idle_timeout,
		},
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) {

	log.Println("[+] HTTP server started on addr: ", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.srv.Addr, err)
	}

	<-ctx.Done()
	log.Println("[!] Server exiting")
	s.srv.Shutdown(ctx)

}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
