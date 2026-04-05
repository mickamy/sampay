#!/bin/bash
set -euo pipefail

ENV_FILE="/app/sampay/.env"
NETWORK="sampay_default"

# Resolve migration image: argument or derive from .env
if [ -n "${1:-}" ]; then
  IMAGE="$1"
else
  ECR_BACKEND_URL=$(grep '^ECR_BACKEND_URL=' "$ENV_FILE" | cut -d= -f2-)
  IMAGE_TAG=$(grep '^IMAGE_TAG=' "$ENV_FILE" | cut -d= -f2-)
  IMAGE="${ECR_BACKEND_URL}:migration-${IMAGE_TAG}"
fi

docker pull "$IMAGE"

REGION=ap-northeast-1

get_secret() {
  aws secretsmanager get-secret-value \
    --region "$REGION" \
    --secret-id "$1" \
    --query SecretString \
    --output text
}

# Read SM_PREFIX from .env.compose
if grep -q '^SM_PREFIX=' "$ENV_FILE"; then
  SM_PREFIX=$(grep '^SM_PREFIX=' "$ENV_FILE" | cut -d= -f2-)
else
  echo "Error: SM_PREFIX is not set in $ENV_FILE" >&2
  exit 1
fi

APP_JSON=$(get_secret "${SM_PREFIX}/app")

run_migration() {
  docker run --rm --network "$NETWORK" \
    -e DOCKER=1 \
    -e MODULE_ROOT=/ \
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
