package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb/backend/internal/cache"
	"wb/backend/internal/config"
	"wb/backend/internal/database"
	"wb/backend/internal/server"
	"wb/backend/internal/services"
	"wb/backend/internal/stanClient"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	defer log.Println("[!] Application shutdown")

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	ctx1, cancel := context.WithCancel(context.Background())
	go grace(cancel)

	// REPO
	postgresClient, err := database.NewPostgresClient(ctx1, cfg.Database.URL)
	if err != nil {
		log.Fatal("[-] failed to connect to postgres, err:", err)
	}
	
	orderRepo, err := database.NewPostgesRepository(ctx1, postgresClient)
	if err != nil {
		log.Fatal("[-] failed to create postgres repository, err:", err)
	}

	// CACHE
	cacheRepository := cache.New(ctx1, cfg.CacheConfig)
	ctx2, cancel2 := context.WithTimeout(ctx1, time.Duration(time.Second*5))
	defer cancel2()

	warmData, err := orderRepo.GetOrdersWithLimitByOrder(ctx2, 100, "", "") // 100 записей отсортировных по времени создания (дефолт)
	if err != nil {
		log.Println("[-] failed to warm cache, reason:", err)
	} else {
		ctx, cancel := context.WithTimeout(ctx1, time.Duration(time.Second*1))
		defer cancel()

		cacheRepository.WarmingCache(ctx, warmData)

		log.Println("[+] Cache successfully warmed")
	}

	// STAN
	natsOpts := []nats.Option{nats.PingInterval(5 * time.Minute), nats.MaxPingsOutstanding(10), nats.MaxReconnects(100)}
	stanOpts := []stan.Option{}

	client := stanClient.New(ctx1, cfg.OrderService.Addr, natsOpts, cfg.OrderService.Subscriber.ClusterID, cfg.OrderService.Subscriber.ClientID, stanOpts)
	defer client.Close()

	// SERVICE
	orderService := services.NewOrderService(orderRepo, cacheRepository, client)
	go orderService.HandleHTTPReq()
	go orderService.HandleNATSStreaming()

	// HTTP
	router := server.NewHandler(orderService, client, cfg.FrontendPath)
	server := server.NewServer(&router.Mux, cfg.HTTPServer)

	go server.Run(ctx1)

	<-ctx1.Done()
}

func grace(c context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("[*] QUIT signal received")
	c()
}
