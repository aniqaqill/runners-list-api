# Base stage
FROM golang:1.23.1 AS base
WORKDIR /usr/src/app

# Development stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
RUN go mod download
RUN go mod tidy
CMD ["air", "-c", ".air.toml"]

# Builder stage (for production and testing)
FROM base AS builder
COPY . .
RUN go mod download
RUN go mod tidy
#not working because
#RUN go build -o app ./cmd/main.go
#sebab ada main.go and routes.go
RUN go build -o app ./cmd/... 
# Build all files in the cmd directory

# Production stage
FROM alpine:latest AS prod
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/app .
COPY --from=builder /usr/src/app/.env .
CMD ["./app"]

# Testing stage
FROM builder AS test
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest
RUN go install github.com/golang/mock/mockgen@latest
WORKDIR /app
CMD ["ginkgo", "-r", "--cover", "--coverprofile=coverage.out", "--trace", "--show-node-events"]