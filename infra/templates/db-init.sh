#!/bin/bash
set -euo pipefail

IMDS_TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60")
REGION=$(curl -s -H "X-aws-ec2-metadata-token: $IMDS_TOKEN" http://169.254.169.254/latest/meta-data/placement/region)
ENV_COMPOSE="/app/.env.compose"
COMPOSE_FILE="/app/compose.prod.yaml"

get_secret() {
  aws secretsmanager get-secret-value \
    --region "$REGION" \
    --secret-id "$1" \
    --query SecretString \
    --output text
}

# Read SM_PREFIX from .env.compose
if grep -q '^SM_PREFIX=' "$ENV_COMPOSE"; then
  SM_PREFIX=$(grep '^SM_PREFIX=' "$ENV_COMPOSE" | cut -d= -f2-)
else
  echo "Error: SM_PREFIX is not set in $ENV_COMPOSE" >&2
  exit 1
fi

APP_JSON=$(get_secret "${SM_PREFIX}/app")

DB_ADMIN_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_ADMIN_PASSWORD)
DB_WRITER_USER=$(echo "$APP_JSON" | jq -r .DB_WRITER_USER)
DB_WRITER_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_WRITER_PASSWORD)
DB_READER_USER=$(echo "$APP_JSON" | jq -r .DB_READER_USER)
DB_READER_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_READER_PASSWORD)
DB_NAME=$(echo "$APP_JSON" | jq -r .DB_NAME)

psql_cmd() {
  docker compose --env-file "$ENV_COMPOSE" -f "$COMPOSE_FILE" exec -T -e PGPASSWORD="$DB_ADMIN_PASSWORD" postgres psql -U postgres "$@"
}

# Wait for postgres to be ready
for i in $(seq 1 30); do
  if docker compose --env-file "$ENV_COMPOSE" -f "$COMPOSE_FILE" exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "Error: postgres not ready after 60s" >&2
    exit 1
  fi
  echo "Waiting for postgres... ($i/30)"
  sleep 2
done

# Create database if not exists
psql_cmd -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
  psql_cmd -c "CREATE DATABASE $DB_NAME"

# Create/update users with passwords from SM
psql_cmd -d "$DB_NAME" -c \
  "DO \$\$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$DB_WRITER_USER') THEN
      CREATE USER $DB_WRITER_USER WITH PASSWORD '$DB_WRITER_PASSWORD';
    ELSE
      ALTER USER $DB_WRITER_USER WITH PASSWORD '$DB_WRITER_PASSWORD';
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$DB_READER_USER') THEN
      CREATE USER $DB_READER_USER WITH PASSWORD '$DB_READER_PASSWORD';
    ELSE
      ALTER USER $DB_READER_USER WITH PASSWORD '$DB_READER_PASSWORD';
    END IF;
  END \$\$;"

# Grant privileges
psql_cmd -d "$DB_NAME" -c \
  "GRANT USAGE ON SCHEMA public TO $DB_WRITER_USER, $DB_READER_USER"

echo "DB initialization complete"
