#!/bin/bash
set -e
DB_URL="${DB_URL:-postgresql://postgres.fkavaynkpqgzifbftdvx:hzYWFL%24vSa9P%2Aii@aws-1-eu-west-1.pooler.supabase.com:6543/postgres?sslmode=require}"
echo "Resetting database migrations..."
migrate -path infrastructures/db/migration -database "$DB_URL" down -all
migrate -path infrastructures/db/migration -database "$DB_URL" up
echo "Database migrations reset successfully."
