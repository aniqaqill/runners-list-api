.PHONY: help build run test clean dev-up dev-down dev-clean docker-up docker-down docker-clean generate-mocks unit-test integration-test coverage clean-port

# Default target
.DEFAULT_GOAL := help

# Variables
GOPATH := $(shell go env GOPATH)
MOCKGEN := $(GOPATH)/bin/mockgen
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

## help: Show this help message
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev-up          - Start development environment with Docker"
	@echo "  make dev-down        - Stop development environment"
	@echo "  make dev-clean       - Stop and remove volumes (clean database)"
	@echo "  make dev-logs        - Show development logs"
	@echo ""
	@echo "Docker (Production):"
	@echo "  make docker-up       - Start production environment"
	@echo "  make docker-down     - Stop production environment"
	@echo "  make docker-clean    - Stop and remove volumes"
	@echo ""
	@echo "Building:"
	@echo "  make build           - Build the Go application"
	@echo "  make run             - Run the application locally"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run all tests"
	@echo "  make unit-test       - Run unit tests with coverage"
	@echo "  make test-verbose    - Run tests with verbose output"
	@echo "  make coverage        - Generate and open coverage report"
	@echo "  make generate-mocks  - Generate mock files for testing"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean           - Clean build artifacts and coverage files"
	@echo "  make clean-port      - Kill process on port 8080"
	@echo "  make fmt             - Format Go code"
	@echo "  make lint            - Run linter"
	@echo ""

## Development Commands

## dev-up: Start development environment
dev-up:
	@echo "Starting development environment..."
	docker compose -f docker-compose.dev.yml up --build

## dev-down: Stop development environment
dev-down:
	@echo "Stopping development environment..."
	docker compose -f docker-compose.dev.yml down

## dev-clean: Stop and remove volumes (clean database)
dev-clean:
	@echo "Cleaning development environment..."
	docker compose -f docker-compose.dev.yml down -v

## dev-logs: Show development logs
dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

## Docker Production Commands

## docker-up: Start production environment
docker-up:
	@echo "Starting production environment..."
	docker compose up --build -d

## docker-down: Stop production environment
docker-down:
	@echo "Stopping production environment..."
	docker compose down

## docker-clean: Stop and remove volumes
docker-clean:
	@echo "Cleaning production environment..."
	docker compose down -v

## Build Commands

## build: Build the Go application
build:
	@echo "Building application..."
	go build -o bin/api ./cmd/main.go

## run: Run the application locally
run:
	@echo "Running application..."
	go run ./cmd/main.go

## Testing Commands

## generate-mocks: Generate mock files for testing
generate-mocks:
	@echo "Generating mock files..."
	@if [ ! -f "$(MOCKGEN)" ]; then \
		echo "Installing mockgen..."; \
		go install github.com/golang/mock/mockgen@latest; \
	fi
	@mkdir -p internal/port/mocks
	@echo "  - Generating events mock..."
	@$(MOCKGEN) -source=internal/port/events.go -destination=internal/port/mocks/events_mock.go -package=mocks
	@echo "  - Generating users mock..."
	@$(MOCKGEN) -source=internal/port/users.go -destination=internal/port/mocks/users_mock.go -package=mocks
	@echo "✅ Mocks generated successfully"

## test: Run all tests
test: generate-mocks
	@echo "Running all tests..."
	go test -v ./...

## unit-test: Run unit tests with coverage
unit-test: generate-mocks
	@echo "Running unit tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) -coverpkg=./internal/... ./internal/...
	@echo "Tests completed"
	@echo "Coverage summary:"
	@go tool cover -func=$(COVERAGE_FILE) | grep total

## test-verbose: Run tests with verbose output
test-verbose: generate-mocks
	@echo "Running tests (verbose)..."
	go test -v -race ./...

## coverage: Generate and open coverage report
coverage: unit-test
	@echo "Generating coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Opening coverage report..."
	@xdg-open $(COVERAGE_HTML) 2>/dev/null || open $(COVERAGE_HTML) 2>/dev/null || start $(COVERAGE_HTML) 2>/dev/null || echo "Please open $(COVERAGE_HTML) manually"

## Utility Commands

## clean: Clean build artifacts and coverage files
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/api
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -rf tmp/
	@echo "Cleaned"

## clean-port: Kill process on port 8080
clean-port:
	@echo "Cleaning port 8080..."
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || echo "Port 8080 is already free"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

## lint: Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## mod: Tidy and verify dependencies
mod:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify
	@echo "Dependencies verified"