.PHONY: build run test lint docker-build docker-run openapi tidy

# ── Build ─────────────────────────────────────────────────────────────────────
build:
	CGO_ENABLED=0 go build -o bin/runners-list-api ./cmd

# ── Run locally (requires a .env file) ────────────────────────────────────────
run: build
	./bin/runners-list-api

# ── Tests ─────────────────────────────────────────────────────────────────────
test:
	go test ./... -v -count=1

# ── Lint (requires golangci-lint) ─────────────────────────────────────────────
lint:
	golangci-lint run

# ── Docker ────────────────────────────────────────────────────────────────────
docker-build:
	docker build -t runners-list-api:local .

docker-run:
	docker run --rm --env-file .env -p 8080:8080 runners-list-api:local

# ── OpenAPI ───────────────────────────────────────────────────────────────────
# Regenerate the openapi.yaml spec. Requires swag CLI:
#   go install github.com/swaggo/swag/cmd/swag@latest
openapi:
	swag init -g cmd/main.go -o docs --outputTypes yaml

# ── Tidy ──────────────────────────────────────────────────────────────────────
tidy:
	go mod tidy
