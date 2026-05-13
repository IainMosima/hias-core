DB_URL = $(DB_SOURCE_HIAS)

postgres:
	docker run --name hias-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=supersecret -d postgres:16-alpine

createdb:
	psql postgresql://root:supersecret@localhost:5432/postgres -c "CREATE DATABASE hias_db;"

dropdb:
	psql postgresql://root:supersecret@localhost:5432/postgres -c "DROP DATABASE IF EXISTS hias_db;"

migrateup:
	migrate -path infrastructures/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path infrastructures/db/migration -database "$(DB_URL)" -verbose down

migrateforce:
	migrate -path infrastructures/db/migration -database "$(DB_URL)" force $(VERSION)

sqlc:
	sqlc generate

swagger:
	swag init -g services/api-gateway/main.go -o docs/swagger --parseDependency --parseInternal

proto:
	protoc \
		--proto_path=services/api-gateway/grpc/proto \
		--go_out=services/api-gateway/grpc/pkg \
		--go-grpc_out=services/api-gateway/grpc/pkg \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		services/api-gateway/grpc/proto/*.proto

run:
	go run services/api-gateway/*.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

seed:
	DB_URL="$(DB_URL)" bash scripts/seed-data.sh

migrate-reset:
	bash scripts/migrate-reset.sh

.PHONY: postgres createdb dropdb migrateup migratedown migrateforce sqlc swagger proto run test lint seed migrate-reset
