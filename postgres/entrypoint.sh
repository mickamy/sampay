#!/usr/bin/env sh
set -e

# Start PostgreSQL in the background using the official entrypoint
docker-entrypoint.sh postgres &
PG_PID=$!

# Wait for PostgreSQL to be ready
# Use 127.0.0.1 explicitly to avoid IPv6 resolution issues in Alpine
TIMEOUT=60
START=$(date +%s)
while :; do
  if ! kill -0 "$PG_PID" 2>/dev/null; then
    echo "PostgreSQL process exited before becoming ready" >&2
    exit 1
  fi
  if pg_isready -h 127.0.0.1 -p 5432 -q; then
    break
  fi
  ELAPSED=$(($(date +%s) - START))
  if [ "$ELAPSED" -ge "$TIMEOUT" ]; then
    echo "Timed out (${TIMEOUT}s) waiting for PostgreSQL" >&2
    exit 1
  fi
  sleep 0.5
done

# Start sql-tap proxy in the foreground
exec sql-tapd \
  --driver=postgres \
  --listen=:5433 \
  --upstream=127.0.0.1:5432 \
  --grpc=:9091 \
  --http=:8081
