.PHONY: help build run test lint docker-up docker-down clean migrate-user migrate-order

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	@echo "Building user-service..."
	@cd services/user-service && go build -o ../../bin/user-service ./cmd/main.go
	@echo "Building order-service..."
	@cd services/order-service && go build -o ../../bin/order-service ./cmd/main.go
	@echo "Building repository-service..."
	@cd services/repository-service && go build -o ../../bin/repository-service ./cmd/main.go
	@echo "Build complete!"

run-user: ## Run user service
	@cd services/user-service && go run ./cmd/main.go

run-order: ## Run order service
	@cd services/order-service && go run ./cmd/main.go

run-repository: ## Run repository service
	@cd services/repository-service && go run ./cmd/main.go

run: ## Run all services concurrently
	@echo "Starting all services..."
	@make -j3 run-user run-order run-repository

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "Test coverage:"
	@go tool cover -func=coverage.txt

test-user: ## Run user service tests
	@echo "Testing user-service..."
	@cd services/user-service && go test -v -race ./...

test-order: ## Run order service tests
	@echo "Testing order-service..."
	@cd services/order-service && go test -v -race ./...

test-repository: ## Run repository service tests
	@echo "Testing repository-service..."
	@cd services/repository-service && go test -v -race ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"

docker-up: ## Start all services with docker-compose
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Waiting for databases to be ready..."
	@sleep 5
	@echo "Services started successfully!"

docker-down: ## Stop all docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped!"

docker-clean: ## Remove all containers, volumes, and images
	@echo "Cleaning up Docker resources..."
	@docker-compose down -v --rmi all
	@echo "Cleanup complete!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.txt
	@echo "Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

migrate-user: ## Run user database migrations
	@echo "Running user service migrations..."
	@cd services/user-service && go run ./cmd/migrate/main.go

migrate-order: ## Run order database migrations
	@echo "Running order service migrations..."
	@cd services/order-service && go run ./cmd/migrate/main.go

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@echo "Tools installed! Make sure $(go env GOPATH)/bin is in your PATH"

dev-user: ## Run user service with hot reload (requires air)
	@cd services/user-service && air

dev-order: ## Run order service with hot reload (requires air)
	@cd services/order-service && air

dev-repository: ## Run repository service with hot reload (requires air)
	@cd services/repository-service && air

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete!"