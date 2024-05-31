package tests

import (
	"context"
	"os"
	"testing"
	"wb/backend/internal/app"
	"wb/backend/internal/cache"
	"wb/backend/internal/config"
	"wb/backend/internal/database"
	"wb/backend/internal/domain/repository"
	"wb/backend/internal/server"
	"wb/backend/internal/stanClient"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/suite"
)

var dbURL string

func init() {
	dbURL = os.Getenv("TEST_DB_URI")
}

type APITestSuite struct {
	suite.Suite

	Server *server.Server
	db     database.PostgresDB

	OrderService *app.OrderService
	repoDB       repository.OrderDB
	repoCache    repository.OrderCache
	stanClient   *stanClient.StanClient
	handler      *server.Handler

}


func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	client, err := database.NewPostgresClient(context.Background(), dbURL)
	if err != nil {
		s.FailNow("Failed to connect to postgres", err)
	}

	s.db = client

	repoDB, err := database.NewPostgesRepository(context.Background(), client)
	if err != nil {
		s.FailNow("Failed to connect init order repo", err)
	}

	memCache := cache.New(context.Background(), config.CacheConfig{})

	s.repoCache = memCache
	s.repoDB = repoDB
	if err := s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}
	stanClient := stanClient.New(context.Background(), "127.0.0.1:4223", []nats.Option{}, "test-cluster", "test-client", []stan.Option{})
	s.stanClient = stanClient

	s.initOrderService()

	s.handler = server.NewHandler(s.OrderService, stanClient, "./frontend/static")
	s.Server = server.NewServer(&s.handler.Mux, config.HTTPServerConfig{})
}

func (s *APITestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}

}

func (s *APITestSuite) initOrderService() {

	OrderService := app.NewOrderService(
		s.repoDB,
		s.repoCache,
		s.stanClient,
	)

	s.OrderService = OrderService
	go s.OrderService.HandleHTTPReq()
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) populateDB() error {

	s.repoDB.InsertOrder(context.Background(), order1)

	s.repoDB.InsertOrder(context.Background(), order2)

	return nil
}
