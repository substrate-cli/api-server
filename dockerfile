# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod/go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build CLI binary
RUN go build -o api-server-cli ./cmd/app

# Run stage
FROM alpine:latest
WORKDIR /app

# Install nc (netcat) for port checking
RUN apk add --no-cache netcat-openbsd

# Copy binary and entrypoint
COPY --from=builder /app/api-server-cli .
COPY entrypoint.sh ./

# Make sure entrypoint is executable
RUN chmod +x entrypoint.sh

# Run entrypoint
CMD ["/app/entrypoint.sh"]
