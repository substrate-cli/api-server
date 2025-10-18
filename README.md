# api-server
entry point to substrate cli

# environment variables

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

# run api-server

```bash
go run ./cmd/app
```

# APIs 

```bash 
curl -X POST http://localhost:8080/api/spin-request \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "build me an elegant portfolio website with dark theme for a full stack developer",
    "model": "anthropic",
    "clustername": "folio-demo"
  }'
```

this is an entry point to substrate-cli, for more informations, follow instructions on https://trysubstrate.com/notes.   
consumer-service - https://github.com/substrate-cli/consumer-service-cli.  
llm-node - https://github.com/substrate-cli/llm-node-cli


