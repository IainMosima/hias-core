# HIAS Core - Go Backend

## Architecture
- DDD + Clean Architecture
- Domain layer: INTERFACES ONLY (entities, repository interfaces, service interfaces, DTOs)
- Infrastructure layer: All implementations (DB repos, cache, queue, workers, scheduler)
- Service layer: Business logic implementations
- API Gateway: REST handlers, routes, middleware + gRPC server

## Tech Stack
- Go 1.24+, Gin, gRPC, SQLC (pgx/v5), golang-migrate, Watermill+SQS, PASETO, Redis, PostgreSQL 16

## Key Commands
- `make run` — Start the server
- `make sqlc` — Regenerate SQLC code (always run after changing query/*.sql)
- `make migrateup` / `make migratedown` — Run migrations
- `make swagger` — Generate Swagger docs
- `make test` — Run tests
- `make seed` — Seed demo data
- `make lint` — Run golangci-lint

## Conventions
- Money: BIGINT (cents), int64 in Go. 8000 KES = 800000
- UUIDs: google/uuid.UUID
- Timestamps: time.Time (TIMESTAMPTZ in DB)
- Status fields: typed string constants in shared/types.go
- All service methods return `*schema.ServiceResponse[T]` (from `domains/identity/schema`)
- Number formats: CLM-YYYY-NNNNNN (claims), POL-YYYY-NNNNNN (policies)

## File Naming
- Entity: `domains/<domain>/entity/<name>.go`
- Repo interface: `domains/<domain>/repository/<name>_repository.go`
- Service interface: `domains/<domain>/service/<name>_service.go`
- Schema (DTOs): `domains/<domain>/schema/<name>_request.go`, `<name>_response.go`
- Repo impl: `infrastructures/repository/<name>_repository.go`
- Service impl: `services/<domain>/<name>_service_impl.go`
- Handler: `services/api-gateway/handlers/<name>_handler.go`
- Routes: `services/api-gateway/routes/routes.go`
- DI wiring: `services/api-gateway/main.go`
- Migrations: `infrastructures/db/migration/NNNNNN_<desc>.{up,down}.sql`
- SQLC queries: `infrastructures/db/query/<entity>.sql`
- pgtype helpers: `infrastructures/repository/pgtype_helpers.go`

## Agent
- Use `/ddd-codegen` for scaffolding new features/domains following established DDD patterns
