.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

docker-build: ## Build Docker image
	docker build -t community-issue-mapper:latest ..

docker-up: ## Start Docker containers
	docker-compose -f ../docker-compose.yml up -d

docker-down: ## Stop Docker containers
	docker-compose -f ../docker-compose.yml down

docker-logs: ## View Docker logs
	docker-compose -f ../docker-compose.yml logs -f

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgres://issueapp:issueapp_password@localhost:5432/community_issues?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgres://issueapp:issueapp_password@localhost:5432/community_issues?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

dev: ## Run in development mode with auto-reload (requires air)
	air
