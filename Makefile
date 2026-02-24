.PHONY: help run test test-coverage migrate-up migrate-down build clean swagger swagger-ikas swagger-all

help:
	@echo "Available commands:"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make swagger        - Generate swagger docs for main app"
	@echo "  make swagger-ikas   - Generate swagger docs for IKAS service"
	@echo "  make swagger-all    - Generate swagger docs for all services"

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html

test-coverage-minimal:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

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

swagger:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude ikas

swagger-ikas:
	cd ikas && swag init -g cmd/main.go -o docs --parseDependency --parseInternal

swagger-all: swagger swagger-ikas
	@echo "Swagger docs generated for all services"



.PHONY: help build up down restart logs clean ps test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker images
	docker-compose build --no-cache

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs
	docker-compose logs -f

logs-app: ## Show app logs only
	docker-compose logs -f app

logs-mysql: ## Show MySQL logs
	docker-compose logs -f mysql

logs-redis: ## Show Redis logs
	docker-compose logs -f redis

ps: ## Show running containers
	docker-compose ps

clean: ## Remove all containers, volumes, and images
	docker-compose down -v --rmi all

clean-volumes: ## Remove only volumes
	docker-compose down -v

shell-app: ## Open shell in app container
	docker-compose exec app sh

shell-mysql: ## Open MySQL shell
	docker-compose exec mysql mysql -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME}

shell-redis: ## Open Redis CLI
	docker-compose exec redis redis-cli -a ${REDIS_PASSWORD}

test: ## Run tests in container
	docker-compose exec app go test -v ./...

migrate: ## Run database migrations
	docker-compose exec mysql mysql -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME} < migrations/001_create_users_table.sql
	docker-compose exec mysql mysql -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME} < migrations/002_create_posts_table.sql

# Development
dev-up: ## Start services for development
	docker-compose -f docker-compose.dev.yml up -d

dev-down: ## Stop development services
	docker-compose -f docker-compose.dev.yml down

# Production
prod-build: ## Build for production
	docker-compose -f docker-compose.prod.yml build

prod-up: ## Start production services
	docker-compose -f docker-compose.prod.yml up -d

prod-down: ## Stop production services
	docker-compose -f docker-compose.prod.yml down

# Monitoring
stats: ## Show container stats
	docker stats

health: ## Check health of all services
	@echo "Checking service health..."
	@docker-compose ps | grep -E "(Up|healthy)"