package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb/backend/internal/app"
	"wb/backend/internal/cache"
	"wb/backend/internal/config"
	"wb/backend/internal/database"
	"wb/backend/internal/server"
	"wb/backend/internal/stanClient"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	ctx1 := context.Background()

	ctx2, cancel := context.WithTimeout(ctx1, time.Duration(time.Second*1))
	defer cancel()

	// REPO
	postgresClient, err := database.NewPostgresClient(ctx2, cfg.Database.URL)
	if err != nil {
		log.Fatal("failed to connect to postgres, err:", err)
	}
	orderRepo, err := database.NewPostgesRepository(postgresClient)
	if err != nil {
		log.Fatal("failed to create postgres repository, err:", err)
	}

	// CACHE
	cacheRepository := cache.New(cfg.CacheConfig)
	ctx3, cancel := context.WithTimeout(ctx1, time.Duration(time.Second*3))
	defer cancel()
	warmData, err := orderRepo.GetOrdersWithLimitByOrder(ctx3, 100, "", "") // 100 записей отсортировных по времени создания (дефолт)
	if err != nil {
		log.Println("failed to warm cache, continue")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
		defer cancel()

		cacheRepository.WarmingCache(ctx, warmData)

		log.Println("cacheRepository successfully warmed")
	}

	// STAN
	natsOpts := []nats.Option{nats.PingInterval(5 * time.Minute), nats.MaxPingsOutstanding(10), nats.MaxReconnects(100)}
	stanOpts := []stan.Option{}

	client := stanClient.New(cfg.OrderService.Addr, natsOpts, cfg.OrderService.Subscriber.ClusterID, cfg.OrderService.Subscriber.ClientID, stanOpts)
	defer client.Close()

	// SERVICE
	orderService := app.NewOrderService(orderRepo, cacheRepository, client)
	go orderService.HandleHTTPReq()
	go orderService.HandleNATSStreaming()

	// HTTP
	router := server.NewRouter(*orderService, client, cfg.FrontendPath)
	server := server.NewServer(&router.Mux, cfg.HTTPServer)
	go server.Run()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("[!] Graceful shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	postgresClient.Close()

	if err := client.Close(); err != nil {
		log.Println("failed to close stan client, err:", err)
	}

	log.Println("Server exiting")
}
