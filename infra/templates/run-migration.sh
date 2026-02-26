#!/bin/bash
set -euo pipefail

ENV_COMPOSE="/app/.env.compose"
NETWORK="app_default"

# Resolve migration image: argument or derive from .env.compose
if [ -n "${1:-}" ]; then
  IMAGE="$1"
else
  ECR_BACKEND_URL=$(grep '^ECR_BACKEND_URL=' "$ENV_COMPOSE" | cut -d= -f2-)
  IMAGE_TAG=$(grep '^IMAGE_TAG=' "$ENV_COMPOSE" | cut -d= -f2-)
  IMAGE="${ECR_BACKEND_URL}:migration-${IMAGE_TAG}"
fi

docker pull "$IMAGE"

IMDS_TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60")
REGION=$(curl -s -H "X-aws-ec2-metadata-token: $IMDS_TOKEN" http://169.254.169.254/latest/meta-data/placement/region)

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

run_migration() {
  docker run --rm --network "$NETWORK" \
    -e DOCKER=1 \
    -e "DB_HOST=$(echo "$APP_JSON" | jq -r .DB_HOST)" \
    -e "DB_PORT=$(echo "$APP_JSON" | jq -r .DB_PORT)" \
    -e "DB_NAME=$(echo "$APP_JSON" | jq -r .DB_NAME)" \
    -e "DB_TIMEZONE=$(echo "$APP_JSON" | jq -r .DB_TIMEZONE)" \
    -e "DB_ADMIN_USER=$(echo "$APP_JSON" | jq -r .DB_ADMIN_USER)" \
    -e "DB_ADMIN_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_ADMIN_PASSWORD)" \
    -e "DB_WRITER_USER=$(echo "$APP_JSON" | jq -r .DB_WRITER_USER)" \
    -e "DB_WRITER_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_WRITER_PASSWORD)" \
    -e "DB_READER_USER=$(echo "$APP_JSON" | jq -r .DB_READER_USER)" \
    -e "DB_READER_PASSWORD=$(echo "$APP_JSON" | jq -r .DB_READER_PASSWORD)" \
    "$IMAGE" "$1"
}

echo "Running db-create..."
run_migration /bin/db-create

echo "Running db-migrate..."
run_migration /bin/db-migrate

echo "Running db-seed..."
run_migration /bin/db-seed

echo "Migration complete"
