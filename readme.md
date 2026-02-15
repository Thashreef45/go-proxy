# go-proxy

A simple **reverse proxy** written in Go using `httputil.ReverseProxy`.  

It forwards incoming HTTP requests to a backend server and streams responses efficiently.  

## Current Features

- Reverse proxy to a backend defined in `config.json`
- Connection pooling with configurable idle connections
- Error handling for backend failures
- Response streaming with configurable flush interval

## Requirements

- Go 1.20+ (or latest stable)
- A backend server to proxy requests to

## Configuration

The proxy uses a JSON configuration file (`config.json`) for setup. Example:

```json
{
    "listen_port": 8080,
    "backend": "http://localhost:8001",
    "max_idle_conns": 500,
    "max_idle_conns_per_host": 500,
    "flush_interval_ms": 50,
    "idle_conn_timeout_sec": 60
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


