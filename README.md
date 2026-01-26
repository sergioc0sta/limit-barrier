# Rate Limiter in Go

A rate limiter in Go that limits requests per second by client IP or by access token, using Redis as the persistence backend.

## Features

- Limits requests per second based on client IP.
- Limits requests per second based on an access token (`API_KEY` header).
- Configurable via `.env` file or environment variables.
- Uses Redis to store counters and block state.
- Supports a persistence strategy (easy to swap Redis for another backend).
- HTTP middleware that can be injected into any web server.
- Returns `429 Too Many Requests` when the limit is exceeded.

## Prerequisites

- Go 1.20+
- Docker and Docker Compose

## How to run

1. Clone the repository:

   ```bash
   git clone https://github.com/sergioc0sta/limit-barrier.git
   cd limit-barrier
   ```

2. Create a `.env` file in the project root and adjust the values:

   ```env
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   REDIS_DB=0
   STORAGE_DRIVER=redis

   IP_MAX_REQ=5
   TOKEN_MAX_REQ=100
   BLOCK_TIME=300
   RATE_LIMIT_DUR=1s
   TOKEN_LIMITS_PATH=configs/token_limits.json
   ```

3. Run everything with Docker Compose (app + Redis):

   ```bash
   docker-compose up --build
   ```

4. The server will be available on port `8080`.

## Tests

Run unit tests:

```bash
go test ./...
```

## Token limits (per token)

You can configure per-token limits via a JSON file (default: `configs/token_limits.json`).
Override the path with `TOKEN_LIMITS_PATH`.
If a request has a token that is not listed in the JSON (or no token at all), the limiter falls back to IP-based limits.

Example `configs/token_limits.json`:

```json
[
  { "token": "TOKEN_BASIC", "limit": 100 },
  { "token": "TOKEN_SILVER", "limit": 200 },
  { "token": "TOKEN_GOLD", "limit": 300 },
  { "token": "TOKEN_PLATINUM", "limit": 400 },
  { "token": "TOKEN_ULTRA", "limit": 500 }
]
```

## How to use

The rate limiter is applied as middleware on protected routes.

### IP-based rate limiting

Just make requests to a protected route; the rate limiter uses the client IP.

Example with `curl`:

```bash
curl http://localhost:8080/ping
```

If the limit is exceeded, the response will be:

```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

With status `429 Too Many Requests`.

### Token-based rate limiting

Send the token in the `API_KEY` header.

Example:

```bash
curl -H "API_KEY: TOKEN_SILVER" http://localhost:8080/ping
```

The rate limiter will use the configured limit for that token (if any), which overrides the IP limit.

## Load testing

Use a simple script to test under load (e.g., `scripts/load-test.sh`).

Basic example:

```bash
for i in {1..20}; do
  curl -s -o /dev/null -w "%{http_code}
" http://localhost:8080/ping &
done
wait
```

## Project structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── configs/
│   └── load_config_redis.go
│   └── token_limits.json
├── internal/
│   ├── limiter/
│   ├── middleware/
│   └── storage/
│       ├── storage.go
│       └── redis/
│           ├── rate_limiter.go
│           └── store.go
├── Dockerfile
├── docker-compose.yml
└── README.md
```
