# Golang Gin Boilerplate — RealWorld API

[![CI](https://github.com/rifkibyte/golang-gin-boilerplate-realworld/actions/workflows/ci.yml/badge.svg)](https://github.com/rifkibyte/golang-gin-boilerplate-realworld/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/rifkibyte/golang-gin-boilerplate-realworld/blob/main/LICENSE)

> Backend API (Golang/Gin) following the [RealWorld API Spec](https://realworld-docs.netlify.app/docs/specs/backend-specs/endpoints) — CRUD, JWT authentication, pagination, and more.

This codebase is a **boilerplate** for building a REST API with **Golang + Gin + GORM v2 + SQLite**, built on top of [gothinkster/golang-gin-realworld-example-app](https://github.com/gothinkster/golang-gin-realworld-example-app).

## Features

- Authentication & authorization using **JWT** (golang-jwt/jwt/v5)
- **GORM v2** with SQLite (dev/test), ready to switch to MySQL/PostgreSQL
- Input validation with **go-playground/validator/v10**
- Routing, pagination, eager loading (`Preload`) without N+1 queries
- High test coverage + CI pipeline (unit test, lint, gosec, API integration test)
- Configuration via environment variables / `.env` file

## Tech Stack

| Package | Version |
|---------|---------|
| [gin-gonic/gin](https://github.com/gin-gonic/gin) | v1.10.0 |
| [gorm.io/gorm](https://gorm.io/) | v1.25.12 |
| [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | v5.2.1 |
| [go-playground/validator/v10](https://github.com/go-playground/validator) | v10.24.0 |
| [joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 |

## Directory Structure

```
.
├── main.go            # Entry point (routes, migration, server)
├── common/            # Shared utilities
│   ├── database.go    # DB connection (SQLite) & test helpers
│   ├── utils.go       # JWT, validation, helpers
│   └── test_helpers.go
├── users/             # User module (auth, profile, follow)
├── articles/          # Article module (CRUD, comments, favorites, tags)
├── api/               # Postman collection for API integration tests
└── scripts/           # Build/test/CI scripts
```

## Quick Start

### Prerequisites

- Go **1.21+** ([download](https://go.dev/dl/))

### Run the Server

```bash
# 1. Copy the environment configuration
cp .env.example .env

# 2. Install dependencies
go mod download

# 3. Run the server (default port 8080)
go run main.go
```

The server will automatically create a SQLite database file at `./data/gorm.db` and run migrations.

### Configuration

Environment variables are read from `.env` (if present) or from the shell. See [.env.example](.env.example) for the full list.

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `debug` | Gin mode (`debug` / `release`) |
| `DB_PATH` | `./data/gorm.db` | SQLite database path |
| `JWT_SECRET` | *(dev fallback)* | Secret for JWT signing — **required in production!** |

> **Important:** In `release` mode (`GIN_MODE=release`), the server will refuse to start if `JWT_SECRET` is not set. Generate a secret with:
>
> ```bash
> openssl rand -base64 32
> ```

## API Endpoints

Following the [RealWorld API Spec](https://realworld-docs.netlify.app/docs/specs/backend-specs/endpoints):

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Register |
| POST | `/api/users/login` | Login |
| GET / PUT | `/api/user` | Get / update current user |
| GET / POST / DELETE | `/api/profiles/:username/follow` | View / follow / unfollow profile |
| GET / POST | `/api/articles` | List / create articles |
| GET / PUT / DELETE | `/api/articles/:slug` | Get / update / delete article |
| GET | `/api/articles/feed` | Article feed from followed users |
| POST / GET / DELETE | `/api/articles/:slug/comments` | Manage comments |
| POST / DELETE | `/api/articles/:slug/favorite` | Favorite / unfavorite article |
| GET | `/api/tags` | List tags |
| GET | `/api/ping` | Health check |

## Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Scripts

| Script | Purpose |
|--------|---------|
| `scripts/gofmt.sh` | Check code formatting |
| `scripts/coverage.sh` | Generate combined coverage report |
| `scripts/run-api-tests.sh` | API integration tests with Newman + Postman collection |

## License

[MIT](LICENSE). Original codebase from [gothinkster/golang-gin-realworld-example-app](https://github.com/gothinkster/golang-gin-realworld-example-app) (MIT).