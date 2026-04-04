.PHONY: run build test dev lint docker-up docker-down migrate-up migrate-down clean help

## run: Start the application
run:
	go run main.go

## build: Build the application binary
build:
	go build -o bin/fiberstarter .

## test: Run all tests
test:
	go test ./... -v

## dev: Run with hot reload (requires Air: go install github.com/cosmtrek/air@latest)
dev:
	air

## lint: Run linter (requires golangci-lint)
lint:
	golangci-lint run

## docker-up: Start all services with Docker Compose
docker-up:
	docker-compose up --build

## docker-down: Stop all Docker Compose services
docker-down:
	docker-compose down

## migrate-up: Run database migrations
migrate-up:
	mysql -u $(DB_USER) -p $(DB_NAME) < migrations/m_00001_Up.sql

## migrate-down: Rollback database migrations
migrate-down:
	mysql -u $(DB_USER) -p $(DB_NAME) < migrations/m_00001_Down.sql

## clean: Remove build artifacts
clean:
	rm -rf bin/ tmp/

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
