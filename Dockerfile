FROM golang:1.23.1

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN go mod download
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]