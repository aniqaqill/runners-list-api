# Development
up-dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up --build

down-dev:
	@echo "Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

# Production
build-prod:
	@echo "Building production image..."
	docker build -t runners-list-api-prod --target prod .

up-prod:
	@echo "Starting production environment..."
	docker-compose -f docker-compose.prod.yml up --build

down-prod:
	@echo "Stopping production environment..."
	docker-compose -f docker-compose.prod.yml down

# Testing
up-test:
	@echo "Starting testing environment..."
	docker-compose -f docker-compose.test.yml up --build

generate-mocks:
	@echo "Generating mock files..."
	@if ! command -v mockgen &> /dev/null; then \
		echo "mockgen not found. Installing mockgen..."; \
		go install github.com/golang/mock/mockgen@latest; \
	fi
	@mkdir -p internal/port/mocks
	mockgen -source=internal/port/events.go -destination=internal/port/mocks/events_mock.go -package=mocks
	mockgen -source=internal/port/users.go -destination=internal/port/mocks/users_mock.go -package=mocks


#if mock not availabe ensure it install in you gopath
#go install github.com/golang/mock/mockgen@v1.6.0
#export PATH=$PATH:$(go env GOPATH)/bin
unit-test: generate-mocks
unit-test: generate-mocks
	@echo "Running tests and generating coverage report..."
	-go test -tags=test -coverprofile=coverage.out -coverpkg=./internal/adapter/middleware,./internal/core/service,./internal/core/domain ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report in the default browser..."
	@open coverage.html || xdg-open coverage.html || start coverage.html || echo "Failed to open coverage.html"

down-test:
	@echo "Stopping testing environment..."
	docker-compose -f docker-compose.test.yml down

# Clean up unused Docker resources
clean:
	@echo "Cleaning up unused Docker resources..."
	docker system prune -f

# Clean unwanted port bindings (incase port had cache)
clean-port:
	@echo "Cleaning port 8080..."
	@lsof -ti :8080 | xargs kill -9 || echo "Port 8080 is already free."

# Clean up mock files
clean-mocks:
	@echo "Cleaning up mock files..."
	@rm -rf internal/port/mocks

# Help command
help:
	@echo "Available commands:"
	@echo "  up-dev      - Start the development environment."
	@echo "  down-dev    - Stop the development environment."
	@echo "  build-prod  - Build the production image."
	@echo "  up-prod     - Start the production environment."
	@echo "  down-prod   - Stop the production environment."
	@echo "  up-test     - Start the testing environment."
	@echo "  down-test   - Stop the testing environment."
	@echo "  generate-mocks - Generate mock files for testing."
	@echo "  unit-test    - Run unit tests and generate coverage report."
	@echo "  clean        - Clean up unused Docker resources."
	@echo "  clean-port   - Clean unwanted port bindings 8080."
	@echo "  clean-mocks  - Clean up mock files."
	@echo "  help         - Show this help message."