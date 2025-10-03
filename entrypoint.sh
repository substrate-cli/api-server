#!/bin/sh
# ==============================
# entrypoint.sh - Alpine-compatible
# Starts all dependencies before launching API server CLI
# ==============================

# Function to wait for a host:port to be ready
wait_for() {
  host=$1
  port=$2

  echo "Waiting for $host:$port..."
  while ! nc -z $host $port; do
    sleep 1
  done
  echo "$host:$port is up!"
}

# Wait for all dependencies
wait_for rabbitmq 5672
wait_for redis 6379
wait_for llm-node 3000
wait_for consumer-service 8090

# All services ready, start API server CLI
echo "Starting API server CLI..."
./api-server-cli
