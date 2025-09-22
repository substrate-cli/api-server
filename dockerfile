# api-server-cli/Dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-server-cli .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/api-server-cli .

CMD ["./api-server-cli"]
