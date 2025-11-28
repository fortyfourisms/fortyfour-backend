.PHONY: help run test test-coverage migrate-up migrate-down build clean

help:
	@echo "Available commands:"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build artifacts"

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up:
	mysql -u root -p myapp_db < migrations/001_create_users_table.sql
	mysql -u root -p myapp_db < migrations/002_create_posts_table.sql

migrate-down:
	mysql -u root -p myapp_db -e "DROP TABLE IF EXISTS posts; DROP TABLE IF EXISTS users;"

build:
	go build -o bin/server cmd/api/main.go

clean:
	rm -rf bin/
	rm -f coverage.out