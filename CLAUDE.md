# HIAS Core - Go Backend

## Architecture
- DDD + Clean Architecture following `sample/denning-backend/` patterns
- Domain layer: INTERFACES ONLY (entities, repository interfaces, service interfaces, DTOs)
- Infrastructure layer: All implementations (DB repos, cache, queue, workers, scheduler)
- Service layer: Business logic implementations
- API Gateway: REST handlers, routes, middleware + gRPC server

## Tech Stack
- Go 1.24+, Gin, gRPC, SQLC, golang-migrate, Watermill+SQS, PASETO, Redis, Cognito

## Key Commands
- `make run` — Start the server
- `make sqlc` — Regenerate SQLC code
- `make migrateup` / `make migratedown` — Run migrations
- `make swagger` — Generate Swagger docs
- `make test` — Run tests

## Conventions
- Money: BIGINT (cents), int64 in Go. 8000 KES = 800000
- UUIDs: google/uuid.UUID
- Timestamps: time.Time (TIMESTAMPTZ in DB)
- Status fields: typed string constants in shared/types.go
- All service methods return *schema.ServiceResponse[T]
- Number formats: CLM-YYYY-NNNNNN (claims), POL-YYYY-NNNNNN (policies)
