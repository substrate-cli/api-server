# api-server
entry point to substrate cli

# environment vars ----

# ***********

```bash
SAFE_ORIGINS="http://localhost:8090, http://localhost:8080, http://localhost:3000"
NODE="api-server"
PORT=8080
MODE="cli"
DEFAULT_USER="987"
SUPPORTED_MODELS="anthropic,openai,gemini"
REDIS_ADDR="localhost:6379"
AMQP_URL="amqp://guest:guest@localhost:5672/"
BUNDLE="server"
```

# ***********

# run api-server

go run ./cmd/app

this is an entry point to substrate-cli, for more informations, follow instrcutions on https://trysubstrate.com/notes

consumer-service - https://github.com/substrate-cli/consumer-service-cli
llm-node - https://github.com/substrate-cli/llm-node-cli


