# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=showroom-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Main targets
.PHONY: all build clean test coverage deps run dev-setup docker-build docker-run migration seed lint

all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Development setup
dev-setup:
	$(GOMOD) download
	$(GOMOD) tidy
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Run database migrations
migration:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME) --migrate-only

# Seed database with initial data
seed:
	$(GOBUILD) -o seed -v ./scripts/seed.go
	./seed
	rm -f seed

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Docker targets
docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Integration tests
test-integration:
	$(GOTEST) -v -tags=integration ./tests/integration/...

# End-to-end tests  
test-e2e:
	$(GOTEST) -v -tags=e2e ./tests/e2e/...

# Database operations
db-reset:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 5
	$(MAKE) migration
	$(MAKE) seed

# Generate mocks
generate-mocks:
	mockgen -source=internal/repositories/interfaces/user.go -destination=tests/mocks/user_repository.go
	mockgen -source=internal/services/auth.go -destination=tests/mocks/auth_service.go

# Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/golang/mock/mockgen@latest

# Watch for changes and restart server
watch:
	air

# API documentation
docs:
	swag init -g cmd/server/main.go

# Security check
security:
	gosec ./...

# Performance benchmarks
benchmark:
	$(GOTEST) -bench=. -benchmem ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  deps           - Download dependencies"
	@echo "  run            - Build and run the application"
	@echo "  dev-setup      - Setup development environment"
	@echo "  migration      - Run database migrations"
	@echo "  seed           - Seed database with initial data"
	@echo "  lint           - Lint code"
	@echo "  fmt            - Format code"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  docker-logs    - View Docker logs"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e       - Run end-to-end tests"
	@echo "  db-reset       - Reset database"
	@echo "  help           - Show this help message"