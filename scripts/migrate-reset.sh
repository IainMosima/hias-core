#!/bin/bash
set -e
DB_URL="${DB_URL:-postgresql://root:supersecret@localhost:5432/hias_db?sslmode=disable}"
echo "Resetting database migrations..."
migrate -path infrastructures/db/migration -database "$DB_URL" down -all
migrate -path infrastructures/db/migration -database "$DB_URL" up
echo "Database migrations reset successfully."
