# =============================================================================
# Stage 1 — build
# =============================================================================
FROM golang:1.24-alpine AS build

WORKDIR /app

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a fully static binary — no libc dependency, safe for
# distroless/scratch base images.
RUN CGO_ENABLED=0 GOOS=linux go build -o runners-list-api ./cmd

# =============================================================================
# Stage 2 — runtime (distroless)
# =============================================================================
# gcr.io/distroless/static-debian12:nonroot ships:
#   - CA certificates       (needed for TLS to Supabase)
#   - tzdata                (needed for Asia/Singapore timezone)
#   - no shell, no package manager, minimal attack surface
#
# The ":nonroot" tag sets the default USER to uid 65532 (nonroot).
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=build /app/runners-list-api .

# Cloud Run injects $PORT at runtime (default 8080). Our Go binary reads it
# from config.Load() → cfg.Port. Declaring it here is documentation only.
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/runners-list-api"]
