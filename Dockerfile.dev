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
