#!/usr/bin/env sh
set -e

# Start PostgreSQL in the background using the official entrypoint
docker-entrypoint.sh postgres &

# Wait for PostgreSQL to be ready
# Use 127.0.0.1 explicitly to avoid IPv6 resolution issues in Alpine
until pg_isready -h 127.0.0.1 -p 5432 -q; do
  sleep 0.5
done

# Start sql-tap proxy in the foreground
exec sql-tapd \
  --driver=postgres \
  --listen=:5433 \
  --upstream=127.0.0.1:5432 \
  --grpc=:9091 \
  --http=:8081
