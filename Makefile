# to enter the container terminal 
clean-port:
	@echo "Cleaning port 8080..."
	@lsof -ti :8080 | xargs kill -9 || echo "Port 8080 is already free."

enter-dock: clean-port
	docker compose run --service-ports web bash

run-dev: clean-port
	docker compose up --build