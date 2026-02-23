# URL Shortener

A REST API for shortening URLs, built with Go as a learning project. Covers layered architecture, PostgreSQL with pgx, graceful shutdown, Docker, and unit testing.

---

## Architecture

```
url-shortener/
├── cmd/api/main.go                   # Entrypoint: wiring, server lifecycle
├── config/config.go                  # Environment variable loading
├── internal/
│   ├── handler/url_handler.go        # HTTP layer: parse requests, write responses
│   ├── service/url_service.go        # Business logic: validation, code generation
│   └── repository/
│       ├── url_repository.go         # Interface + in-memory implementation
│       └── postgres_repository.go    # PostgreSQL implementation
├── pkg/response/response.go          # JSON/error response helpers
├── db/
│   ├── migrations/001_initial_schema.sql
│   └── queries/urls.sql
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

The `URLRepository` interface is the key design decision: `service` depends on the interface, not on a concrete implementation. Swapping storage requires changing one line in `main.go`.

---

## Requirements

- Go 1.23+
- Docker and Docker Compose (for the database)

---

## Setup

**1. Clone and enter the project**

```bash
git clone <repo>
cd url-shortener
```

**2. Configure environment**

```bash
cp .env.example .env
# Edit .env if needed (defaults work with docker-compose)
```

**3. Start the database**

```bash
make db-up
```

This starts a PostgreSQL container on port 5432 and runs the migration from `db/migrations/` automatically.

**4. Download dependencies and run**

```bash
go mod tidy
make run
```

The server starts at `http://localhost:8080`.

**Alternative: run everything with Docker**

```bash
make docker-up
```

---

## Running Tests

Tests are in `internal/service/` and use the in-memory repository, so no database is required.

```bash
make test

# Or directly:
go test -v -race ./internal/...
```

The `-race` flag enables the data race detector. Always use it during development.

---

## API Reference

### Create a short URL

```
POST /shorten
```

Request body:

```json
{
  "url": "https://golang.org/doc/effective_go",
  "custom_code": "effective-go"
}
```

`custom_code` is optional. If omitted, a 6-character alphanumeric code is generated.

Response `201 Created`:

```json
{
  "code": "effective-go",
  "original_url": "https://golang.org/doc/effective_go",
  "short_url": "http://localhost:8080/effective-go"
}
```

---

### Redirect to original URL

```
GET /{code}
```

Returns `301 Moved Permanently` to the original URL and increments the click counter.

```bash
curl -L http://localhost:8080/effective-go
```

---

### Get URL statistics

```
GET /stats/{code}
```

Response `200 OK`:

```json
{
  "code": "effective-go",
  "original_url": "https://golang.org/doc/effective_go",
  "clicks": 14
}
```

---

### List all URLs

```
GET /urls
```

Response `200 OK`:

```json
[
  {
    "Code": "effective-go",
    "OriginalURL": "https://golang.org/doc/effective_go",
    "Clicks": 14
  }
]
```

---

### Delete a short URL

```
DELETE /urls/{code}
```

Response `204 No Content` on success.

---

### Health check

```
GET /health
```

Response `200 OK`:

```json
{"status": "healthy"}
```

Returns `503 Service Unavailable` if the database is unreachable.

---

## Error responses

All errors return JSON:

```json
{"error": "description of the problem"}
```

| Status | Meaning |
|--------|---------|
| 400 | Invalid URL or malformed request body |
| 404 | Code not found |
| 409 | Custom code already in use |
| 500 | Internal server error |

---

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port |
| `BASE_URL` | `http://localhost:8080` | Used to build the `short_url` in responses |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `urlshortener` | Database name |
| `DB_SSLMODE` | `disable` | `disable` for local, `require` for production |

---

## Go concepts covered

- **Interfaces** (`URLRepository`): decouples business logic from storage. Tests use in-memory; production uses PostgreSQL. Swapping requires one line.
- **`sync.RWMutex`**: safe concurrent map access. `RLock` allows multiple concurrent reads; `Lock` gives exclusive write access.
- **`errors.Is`**: idiomatic error comparison. Errors are values in Go, not exceptions.
- **`context.WithTimeout`**: every database call has a deadline. Prevents slow queries from blocking goroutines indefinitely.
- **Graceful shutdown**: waits for in-flight requests to complete before exiting. Triggered by `SIGINT` (Ctrl+C) or `SIGTERM` (Docker stop, Kubernetes).
- **`log/slog`**: structured logging from Go 1.21. Outputs key-value pairs instead of unstructured strings.
- **`pgxpool`**: PostgreSQL connection pool. Reuses connections instead of opening a new one per request.

---

## Next steps

1. Add URL expiration using the `expires_at` column already in the schema
2. Add pagination to `GET /urls`
3. Add rate limiting per IP using `golang.org/x/time/rate`
4. Add JWT authentication to protect write endpoints
5. Add integration tests using `testcontainers-go` to spin up a real PostgreSQL in tests