#!/bin/bash
set -euo pipefail

IMDS_TOKEN=$(curl -s -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60")
REGION=$(curl -s -H "X-aws-ec2-metadata-token: $IMDS_TOKEN" http://169.254.169.254/latest/meta-data/placement/region)
ENV_COMPOSE="/app/.env.compose"
ENV_POSTGRES="/app/.env.postgres"

get_secret() {
  aws secretsmanager get-secret-value \
    --region "$REGION" \
    --secret-id "$1" \
    --query SecretString \
    --output text
}

if grep -q '^SM_PREFIX=' "$ENV_COMPOSE"; then
  SM_PREFIX=$(grep '^SM_PREFIX=' "$ENV_COMPOSE" | cut -d= -f2-)
elif [ -n "${SM_PREFIX:-}" ]; then
  echo "SM_PREFIX=$SM_PREFIX" >> "$ENV_COMPOSE"
else
  echo "Error: SM_PREFIX is not set in $ENV_COMPOSE or in the environment" >&2
  exit 1
fi

# --- DB secrets -> .env.postgres ---
DB_JSON=$(get_secret "${SM_PREFIX}/db")
{
  echo "POSTGRES_USER=$(echo "$DB_JSON" | jq -r .POSTGRES_USER)"
  echo "POSTGRES_PASSWORD=$(echo "$DB_JSON" | jq -r .POSTGRES_PASSWORD)"
} > "$ENV_POSTGRES"
chmod 600 "$ENV_POSTGRES"

# --- KVS password -> .env.compose (append) ---
KVS_JSON=$(get_secret "${SM_PREFIX}/kvs")
KVS_PASSWORD=$(echo "$KVS_JSON" | jq -r .KVS_PASSWORD)
sed -i '/^KVS_PASSWORD=/d' "$ENV_COMPOSE"
echo "KVS_PASSWORD=${KVS_PASSWORD}" >> "$ENV_COMPOSE"

# --- Session secret -> .env.compose (append) ---
APP_JSON=$(get_secret "${SM_PREFIX}/app")
SESSION_SECRET=$(echo "$APP_JSON" | jq -r .SESSION_SECRET)
sed -i '/^SESSION_SECRET=/d' "$ENV_COMPOSE"
echo "SESSION_SECRET=${SESSION_SECRET}" >> "$ENV_COMPOSE"

chmod 600 "$ENV_COMPOSE" "$ENV_POSTGRES"
