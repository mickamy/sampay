#!/bin/bash
set -euo pipefail

ENV_FILE="/app/sampay/.env"
ENV_POSTGRES_FILE="/app/sampay/.env.postgres"

get_secret() {
  aws secretsmanager get-secret-value \
    --region "$REGION" \
    --secret-id "$1" \
    --query SecretString \
    --output text
}

if grep -q '^SM_PREFIX=' "$ENV_FILE"; then
  SM_PREFIX=$(grep '^SM_PREFIX=' "$ENV_FILE" | cut -d= -f2-)
elif [ -n "${SM_PREFIX:-}" ]; then
  echo "SM_PREFIX=$SM_PREFIX" >> "$ENV_FILE"
else
  echo "Error: SM_PREFIX is not set in $ENV_FILE or in the environment" >&2
  exit 1
fi

# --- DB secrets -> .env.postgres ---
DB_JSON=$(get_secret "${SM_PREFIX}/db")
{
  echo "POSTGRES_USER=$(echo "$DB_JSON" | jq -r .POSTGRES_USER)"
  echo "POSTGRES_PASSWORD=$(echo "$DB_JSON" | jq -r .POSTGRES_PASSWORD)"
} > "$ENV_POSTGRES_FILE"
chmod 600 "$ENV_POSTGRES_FILE"

# --- KVS password -> .env (append) ---
KVS_JSON=$(get_secret "${SM_PREFIX}/kvs")
KVS_PASSWORD=$(echo "$KVS_JSON" | jq -r .KVS_PASSWORD)
sed -i '/^KVS_PASSWORD=/d' "$ENV_FILE"
echo "KVS_PASSWORD=${KVS_PASSWORD}" >> "$ENV_FILE"

# --- Session secret -> .env (append) ---
APP_JSON=$(get_secret "${SM_PREFIX}/app")
SESSION_SECRET=$(echo "$APP_JSON" | jq -r .SESSION_SECRET)
sed -i '/^SESSION_SECRET=/d' "$ENV_FILE"
echo "SESSION_SECRET=${SESSION_SECRET}" >> "$ENV_FILE"

chmod 600 "$ENV_FILE" "$ENV_POSTGRES_FILE"
