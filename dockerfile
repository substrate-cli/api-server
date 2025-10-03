# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-server-cli ./cmd/app

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache netcat-openbsd

# Copy binary and entrypoint
COPY --from=builder /app/api-server-cli .
COPY entrypoint.sh ./

RUN chmod +x entrypoint.sh

# Run entrypoint
CMD ["/app/entrypoint.sh"]
