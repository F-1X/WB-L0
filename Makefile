build-server:
	@go build -o bin/wb backend/cmd/*

go-test:
	go clean -testcache
	@go test -v ./...
	
run-server: build-server
	@./bin/wb


docker-up:
	sudo docker-compose up

docker-down:
	sudo docker-compose down

install-psql:
	sudo apt-get install -y postgresql-client 

CURDIR=$(shell pwd)

create-db:
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -f "$(CURDIR)/create_database.sql"


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

