# Go-Store — Cloud-Based File Storage System

A production-grade REST API for uploading, downloading, and sharing files securely in the cloud.

Built with **Go 1.21**, **Gin**, **PostgreSQL**, **Amazon S3 / Google Cloud Storage**, and **JWT authentication**.

---

## Architecture

```
cmd/api/
└── main.go                   Entry point — loads config, starts server

internal/
├── config/config.go          Reads .env + env vars; validates required fields
├── database/database.go      Opens + pings the *sql.DB connection pool
├── auth/
│   ├── jwt.go                GenerateToken / ValidateToken (golang-jwt/jwt v5)
│   ├── anonymous.go          Anonymous session ID generator
│   └── middleware.go         CORS, request logger, JWTAuthMiddleware
├── models/                   Pure data structs — NO database calls
│   ├── user.go
│   ├── file.go
│   └── permission.go
├── repository/               database/sql queries — one struct per table
│   ├── user_repository.go
│   ├── file_repository.go
│   └── permission_repository.go
├── storage/                  Cloud storage abstraction
│   ├── storage.go            Storage interface
│   ├── s3_storage.go         AWS S3 implementation
│   └── gcs_storage.go        Google Cloud Storage implementation
├── services/                 Business logic — orchestrates repo + storage
│   ├── auth_service.go
│   ├── user_service.go
│   └── file_service.go
├── controllers/              Thin HTTP handlers — parse input, call service, write response
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── file_controller.go
├── merrors/merrors.go        Consistent JSON error responses
├── utils/response.go         OK / Created helpers
└── server/
    ├── routes.go             Dependency wiring + route registration
    └── server.go             HTTP server lifecycle, graceful shutdown
```

---

## Prerequisites

- Go 1.21+
- PostgreSQL 14+ (or MySQL 8+)
- AWS account with S3 **or** Google Cloud account with GCS
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (for running migrations)
- [Task](https://taskfile.dev) (optional, replaces Make)

---

## Setup

### 1. Clone

```bash
git clone https://github.com/souvik03-136/Go-Store.git
cd Go-Store
```

### 2. Environment variables

```bash
cp .env.example .env
# Edit .env with your real database credentials, JWT secret, and storage keys.
```

Key variables:

| Variable | Description | Default |
|---|---|---|
| `SERVER_PORT` | HTTP listen port | `8080` |
| `APP_ENV` | `development` or `production` | `development` |
| `DB_DRIVER` | `postgres` or `mysql` | `postgres` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_DATABASE` | Database name | `gostore` |
| `DB_USERNAME` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | *(required)* |
| `JWT_SECRET_KEY` | HMAC signing secret — keep long & random | *(required)* |
| `STORAGE_PROVIDER` | `s3` or `gcs` | `s3` |

### 3. Run migrations

```bash
./scripts/migrate.sh up
# or: task migrate:up
```

### 4. Run the server

```bash
go run ./cmd/api
# or: task run
```

The API is available at `http://localhost:8080`.

### 5. Docker Compose (optional)

```bash
docker compose up --build
# or: task docker:up
```

---

## API Reference

All protected routes require these headers:

```
Authorization: Bearer <token>
X-Salt: <salt>
```

Both `token` and `salt` are returned by the login / register endpoints.

### Auth

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/v1/auth/oauth/register` | No | Register with username + email + password |
| `POST` | `/v1/auth/oauth/login` | No | Login; returns `token` + `salt` |
| `POST` | `/v1/auth/anonymous/register` | No | Create an anonymous session |
| `POST` | `/v1/auth/logout` | No | Client-side logout |
| `GET`  | `/v1/auth/validate?token=&salt=` | No | Validate a token |

**Register example:**
```bash
curl -X POST http://localhost:8080/v1/auth/oauth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"s3cur3pass"}'
```

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST`   | `/v1/users`     | Yes | Create a user |
| `GET`    | `/v1/users/:id` | Yes | Get user by ID |
| `PUT`    | `/v1/users/:id` | Yes | Update user |
| `DELETE` | `/v1/users/:id` | Yes | Delete user |

### Files

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST`   | `/v1/files`     | Yes | Upload a file (`multipart/form-data`, field: `file`) |
| `GET`    | `/v1/files`     | Yes | List all files owned by the authenticated user |
| `GET`    | `/v1/files/:id` | Yes | Get file metadata by ID |
| `PUT`    | `/v1/files/:id` | Yes | Update file metadata (name, content_type) |
| `DELETE` | `/v1/files/:id` | Yes | Delete file (removes cloud object + metadata) |

**Upload example:**
```bash
curl -X POST http://localhost:8080/v1/files \
  -H "Authorization: Bearer <token>" \
  -H "X-Salt: <salt>" \
  -F "file=@/path/to/document.pdf"
```

### Health

```
GET /healthz   →  {"status":"ok"}  or  {"status":"unhealthy","db":"..."}
```

---

## Development

```bash
task test          # run tests with race detector
task lint          # golangci-lint
task fmt           # gofmt
task build         # compile to bin/api
task migrate:up    # apply migrations
task migrate:down  # roll back last migration
```

---

## Security notes

- Passwords are hashed with bcrypt (cost 10).
- Each JWT is signed with `HMAC-SHA256` using a **per-token salt** derived key — a leaked token cannot be used to forge others.
- The `X-Salt` header is required to validate any token. Clients must store it alongside the token.
- File ownership is enforced at the service layer; only the owner can update or delete their files.
- The Docker image runs as a non-root user (`nonroot:nonroot`).

---

## License

MIT — see [LICENSE](./LICENSE).