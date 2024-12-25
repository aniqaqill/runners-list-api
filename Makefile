# to enter the container terminal 
clean-port:
	@echo "Cleaning port 3000..."
	@lsof -ti :3000 | xargs kill -9 || echo "Port 3000 is already free."

enter-dock: clean-port
	docker compose run --service-ports web bash

run-dock: clean-port
	docker compose up --build