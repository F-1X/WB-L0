build-server:
	@go build -o bin/wb backend/cmd/order_service/main.go

run-server: build-server
	@./bin/wb

go-test:
	go clean -testcache
	@go test -v ./...

docker-up:
	sudo docker-compose up

docker-down:
	sudo docker-compose down

install-psql:
	sudo apt-get install -y postgresql-client 

CURDIR=$(shell pwd)

include .env
export

mig-create:
	migrate create -ext sql -dir backend/internal/database/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))
%:
	@:
	
mig-up:
	migrate -database "${DATABASE_URL}?sslmode=disable" -path backend/internal/database/migrate/migrations up

mig-down:
	migrate -database "${DATABASE_URL}?sslmode=disable" -path backend/internal/database/migrate/migrations down

mig-drop:
	migrate -database "${DATABASE_URL}?sslmode=disable" -path backend/internal/database/migrate/migrations drop

mig-version:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -path $(MIGRATIONS_PATH) version

clean-db:
	PGPASSWORD=$(DB_PASSWORD) psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"


DB_USER=test
DB_PASSWORD=test
export DB_NAME
CONTAINER_NAME1=test
CONTAINER_NAME2=test2

test.integration:
	docker run --rm -d -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) --name $(CONTAINER_NAME1) postgres:alpine 
	@sleep 3 
	docker exec -i $(CONTAINER_NAME1) psql -U $(DB_USER) -d $(DB_NAME) < ./test-migration.sql
	docker run -p 4223:4223 -p 8223:8223 --name $(CONTAINER_NAME2) nats-streaming:latest -p 4223 -m 8223

kill-containers:
	@docker kill $(CONTAINER_NAME1) $(CONTAINER_NAME2) || true
	@docker rm -f $(CONTAINER_NAME1) $(CONTAINER_NAME2) || true


run-test:
	go test -v ./backend/internal/tests/

run-test-load:
	go test -v ./backend/internal/tests/load_test.go

test.coverage:
	go test --short -coverprofile=cover.out -v ./...
	go tool cover -func=cover.out

lint:
	golangci-lint run