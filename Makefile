.PHONY: help build run test lint link-check swagger frontend-install frontend-test frontend-build docker-up docker-down clean migrate-user migrate-order run-user run-order run-audit-log test-user test-order test-audit-log

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build all services
	@echo "Building user-service..."
	@cd services/user-service && go build -o ../../bin/user-service ./cmd/main.go
	@echo "Building order-service..."
	@cd services/order-service && go build -o ../../bin/order-service ./cmd/main.go
	@echo "Building audit-log-service..."
	@cd services/audit-log-service && go build -o ../../bin/audit-log-service ./cmd/main.go
	@echo "Build complete!"

run-user: ## Run user service
	@cd services/user-service && go run ./cmd/main.go

run-order: ## Run order service
	@cd services/order-service && go run ./cmd/main.go

run-audit-log: ## Run audit log service
	@cd services/audit-log-service && go run ./cmd/main.go

run: ## Run all services concurrently
	@echo "Starting all services..."
	@make -j3 run-user run-order run-audit-log

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

test-audit-log: ## Run audit log service tests
	@echo "Testing audit-log-service..."
	@cd services/audit-log-service && go test -v -race ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"

link-check: ## Check docs links (requires lychee)
	@echo "Running link check..."
	@command -v lychee >/dev/null 2>&1 && lychee --config .lychee.toml README.md EXAMPLES.md || echo "lychee not installed. See https://github.com/lycheeverse/lychee"

swagger: ## Generate Swagger docs (requires swag)
	@echo "Generating Swagger docs..."
	@SWAG=$$(command -v swag || command -v $$(go env GOPATH)/bin/swag); \
	if [ -z "$$SWAG" ]; then echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@v1.16.4"; exit 1; fi; \
	(cd services/user-service && $$SWAG init -g cmd/main.go -o docs --parseDependency --parseInternal); \
	(cd services/order-service && $$SWAG init -g cmd/main.go -o docs --parseDependency --parseInternal); \
	(cd services/audit-log-service && $$SWAG init -g cmd/main.go -o docs --parseDependency --parseInternal)

frontend-install: ## Install frontend dependencies
	@cd frontend && npm install

frontend-test: ## Run frontend tests
	@cd frontend && npm test

frontend-build: ## Build frontend assets
	@cd frontend && npm run build

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

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete!"
