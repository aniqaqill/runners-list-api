# Development
up-dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up --build

down-dev:
	@echo "Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

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
#export PATH=$PATH:$(go env GOPATH)/bin <- prev not working always hv to initalize but once added to zhrchrc it will work 
#export PATH=$PATH:$(go env GOPATH)/bin >> ~/.zshrc
unit-test: generate-mocks
	@echo "Running tests and generating coverage report..."
	-go test -tags=test -coverprofile=coverage.out -coverpkg=./internal/adapter/middleware,./internal/core/service,./internal/core/domain ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report in the default browser..."
	@open coverage.html || xdg-open coverage.html || start coverage.html || echo "Failed to open coverage.html"

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
	@echo "  generate-mocks - Generate mock files for testing."
	@echo "  unit-test    - Run unit tests and generate coverage report."
	@echo "  clean-port   - Clean unwanted port bindings 8080."
	@echo "  clean-mocks  - Clean up mock files."
	@echo "  help         - Show this help message."