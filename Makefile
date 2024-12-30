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
	@echo "  clean       - Clean up unused Docker resources."
	@echo "  clean-port  - Clean unwanted port bindings 8080."
	@echo "  help        - Show this help message."