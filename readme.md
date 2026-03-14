# go-proxy

A lightweight **reverse proxy** written in Go using `httputil.ReverseProxy`.

It forwards incoming HTTP requests to a backend server while supporting configurable
connection pooling, response streaming, and **config-driven route caching**.

---

# Features

- Reverse proxy using Go's `httputil.ReverseProxy`
- Config-driven backend configuration
- Connection pooling and idle connection tuning
- Response streaming with configurable flush interval
- Graceful backend error handling
- Route-based response caching
- Rate limiting middleware
- Request logging middleware
- CORS middleware
- Modular middleware architecture
- Unit tests for middleware components


---

# Requirements

- Go **1.20+** (or latest stable)
- A backend server to proxy requests to

---

# Configuration

The proxy uses a JSON configuration file (`config.json`) for setup.

Example:

```json
{
  "listen_port": 8080,
  "backend": "http://localhost:8001",
  "max_idle_conns": 500,
  "max_idle_conns_per_host": 500,
  "flush_interval_ms": 50,
  "idle_conn_timeout_sec": 60,
  "cache": {
      "routes": [
          {
              "path": "/api/products",
        "ttl": 60
      }
    ],
    "capacity": 50
  }
}
```

---

# Caching

The proxy supports **route-based caching** configured through `config.json`.

Example configuration:

```json
{
  "cache": {
    "routes": [
      {
        "path": "/api/products",
        "ttl": 60
      }
    ],
    "capacity": 50
  }
}
```

### Cache Behavior

| Request | Behavior |
|--------|----------|
| First request | Cache miss → backend called |
| Next requests | Cache hit → response served from cache |

Caching is typically applied to **GET requests** to reduce backend load and improve response time.

---

# Running the Proxy

Start the proxy with:

```bash
go run main.go
```

The proxy will start listening on the configured port and forward requests to the backend server.

---

# Example

Start a backend server:

```bash
python3 -m http.server 8001
```

Run the proxy:

```bash
go run main.go
```

Now send requests through the proxy:

```bash
curl http://localhost:8080
```

The request will be forwarded to the backend server.

---

## Project Structure

```
proxy-server/
├── config.json           # Application configuration
├── demo.go               # Demo backend server for testing proxy
├── go.mod
├── go.sum
├── LICENSE
├── logs/
│   └── proxy.log         # Proxy logs
├── README.md
└── src/
    ├── cmd/
    │   └── main.go       # Application entry point
    │
    ├── boot/
    │   └── server.go     # Server bootstrap and initialization
    │
    ├── app/
    │   ├── handler/
    │   │   └── proxy-handler.go   # Reverse proxy handler
    │   │
    │   └── middleware/
    │       ├── cors-middleware.go
    │       ├── logger-middleware.go
    │       │
    │       ├── cache/              # Response caching middleware
    │       │   ├── index.go
    │       │   ├── response-writer.go
    │       │   ├── usecase.go
    │       │   ├── z_types.go
    │       │   └── z_test.go
    │       │
    │       └── ratelimiter/        # Rate limiting middleware
    │           ├── index.go
    │           ├── usecase.go
    │           ├── z_type.go
    │           └── z_test.go
    │
    └── internal/
        └── model/
            ├── config-model.go
            └── logs-model.go
```


---

# Testing

Unit tests verify middleware behavior including caching logic.

Run tests using:

```bash
go test ./...
```

Test cases include:

- Cached routes
- Non-cached routes
- Cache hit behavior

---

# License

This project is licensed under the **MIT License**.

See the [LICENSE](LICENSE) file for details.
