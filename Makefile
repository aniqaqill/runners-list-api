up-dev:
	@echo "Starting development environment..."
	docker compose -f docker-compose.dev.yml up --build

down-dev:
	@echo "Stopping development environment..."
	docker compose -f docker-compose.dev.yml down


remove-old-network:
	@echo "Removing old Docker network..."
	@if docker network inspect dev-network >/dev/null 2>&1; then \
		docker network rm dev-network || echo "Network dev-network does not exist."; \
	else \
		echo "Network dev-network does not exist."; \
	fi

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
unit-test: generate-mocks
	@echo "Running tests and generating coverage report..."
	-go test -tags=test -coverprofile=coverage.out -coverpkg=./internal/adapter/middleware,./internal/core/service,./internal/core/domain ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report in the default browser..."
	@open coverage.html || xdg-open coverage.html || start coverage.html || echo "Failed to open coverage.html"

# Clean unwanted port bindings (incase port had cache)
clean-port:
	@echo "Cleaning port 8080..."
	@lsof -ti :8080 | xargs kill -9 || echo "Port 8080 is already free."


# Help command
help:
	@echo "Available commands:"
	@echo "  up-dev      - Start the development environment."
	@echo "  down-dev    - Stop the development environment."
	@echo "  generate-mocks - Generate mock files for testing."
	@echo "  unit-test    - Run unit tests and generate coverage report."
	@echo "  clean-port   - Clean unwanted port bindings 8080."
	@echo "  help         - Show this help message."