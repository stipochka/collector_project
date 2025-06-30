FROM golang:1.24.4 AS builder 

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download && go mod verify 

COPY . .

RUN go build -o collector-service ./cmd/

FROM debian:bookworm-slim

WORKDIR /app 

COPY --from=builder /app/collector-service .

CMD ["./collector-service", "-path", "/app/config/config.yaml" ]