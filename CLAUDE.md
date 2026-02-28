# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Motico API** is a multi-tenant REST API for inventory management built in Go using a hexagonal architecture pattern.

- **Language**: Go 1.25.5+
- **Framework**: Chi v5 (HTTP router)
- **Database**: PostgreSQL (via pgx driver) hosted on Supabase
- **Authentication**: JWT with Bearer tokens
- **Documentation**: Swagger (Swaggo)

## Quick Start Commands

### Running the Application
```bash
go run cmd/api/main.go
```
Server starts at `http://0.0.0.0:8080` with API routes under `/api/v1`

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/domain/product

# Run single test
go test ./internal/domain/product -run TestValidate
```

### Code Quality
```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Full linter check (configured in .golangci.yml)
golangci-lint run ./...

# Install pre-commit hooks
pre-commit install
```

The pre-commit hooks automatically run `go fmt`, `go vet`, `go test` (short mode), and `golangci-lint` on every commit.

## Architecture

### Hexagonal Pattern Structure

```
internal/
├── domain/          # Business logic (entities, use cases, interfaces)
│   ├── auth/        # Authentication service
│   ├── category/    # Category management
│   ├── product/     # Product management
│   ├── stock/       # Stock management
│   ├── store/       # Store management
│   ├── transfer/    # Stock transfer logic
│   └── tenant/      # Multi-tenancy
├── repository/      # Data access layer (adapters)
│   ├── connection.go    # Database connection pool
│   ├── category.go
│   ├── product.go
│   ├── stock.go
│   ├── store.go
│   ├── transfer.go
│   └── tenant.go
└── rest/           # HTTP handlers (adapters)
    ├── router.go       # Route definitions
    ├── middleware.go   # Auth, tenant, recovery
    └── {domain}/       # Domain-specific handlers
```

### Key Architectural Points

1. **Domain Modules**: Each domain (product, stock, etc.) contains:
   - `entities/` - Domain models
   - `repository.go` - Repository interface
   - `service.go` - Business logic (uses repository interface)

2. **Dependency Injection**: Dependencies are injected at startup in `cmd/api/main.go`:
   - Services depend on repository interfaces
   - Handlers depend on services
   - Router receives all handlers via `RouterDependencies`

3. **Multi-tenancy**: Enforced via:
   - `TenantMiddleware` extracts tenant from request context
   - All repository queries filter by tenant ID
   - Tenant ID required in request headers/context

4. **Middleware Stack** (in order):
   - RequestID, RealIP (Chi defaults)
   - Logger middleware
   - Recovery middleware
   - Content-Type validation
   - Tenant extraction (for protected routes)
   - JWT authentication (for protected routes)

## Configuration

### Environment Variables (.env)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL connection
- `JWT_SECRET_KEY` - Min 32 characters for signing tokens

### Application Config (config/config.json)
- Server host, port, timeouts
- Logging level and format
- Database pool settings

## Common Development Tasks

### Adding a New Domain Entity

1. Create domain module: `internal/domain/{entity}/`
2. Define entity in `entities/{entity}.go`
3. Create repository interface in `repository.go`
4. Implement repository in `internal/repository/{entity}.go` (uses pgx)
5. Create service in `service.go` with business logic
6. Create handler in `internal/rest/{entity}/handler.go`
7. Register routes in `internal/rest/router.go`
8. Add Swagger documentation comments on handlers

### Running Migrations

Migrations are in `migrations/` directory. Execute manually in Supabase SQL editor or via `psql`:
```bash
psql -U $DB_USER -h $DB_HOST -d $DB_NAME < migrations/001_initial_schema.sql
```

### Debugging Database Issues

Use the diagnostic scripts:
```bash
bash scripts/test-connection      # Test DB connectivity
bash scripts/diagnose-connection  # Detailed diagnostics
```

Connection pool is created in `internal/repository/connection.go` and managed in main.go with graceful shutdown.

## Testing Patterns

- Tests use short flag: `go test ./... -short` (for pre-commit)
- Integration tests connect to actual database (not short)
- Repository tests in `internal/repository/*_test.go`
- Mock dependencies where needed (avoid testing database calls)

## Swagger Documentation

- Swagger spec auto-generated from code comments
- Regenerate with: `swag init -g cmd/api/main.go`
- Access at: `http://localhost:8080/swagger/index.html`
- Bearer token auth defined in swagger comments

## Important Files

- `cmd/api/main.go` - Application bootstrap and dependency injection
- `internal/rest/router.go` - Route definitions and middleware order
- `internal/rest/middleware.go` - Auth/tenant/recovery middleware
- `.golangci.yml` - Linter configuration (includes security, style, performance checks)
- `.pre-commit-config.yaml` - Git hooks configuration

## Code Standards

- **Import Organization**: Local imports use `motico-api` prefix (configured in `.golangci.yml`)
- **Error Handling**: Repository errors defined in `internal/repository/errors.go`
- **Logging**: Use `pkg/logger` with structured fields (Zap logger)
- **Graceful Shutdown**: Database pool closes before server shutdown (important for data consistency)
