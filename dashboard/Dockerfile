FROM golang:1.24.4 AS builder 

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download && go mod verify 

COPY . .

RUN go build -o dashboard-service ./cmd/

FROM debian:bookworm-slim

WORKDIR /app 

COPY --from=builder /app/dashboard-service .
COPY --from=builder /app/.env .

CMD [ "./dashboard-service" ]