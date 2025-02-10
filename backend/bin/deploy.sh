#!/usr/bin/env bash
set -euox pipefail

if [ -z "${DIR_SUFFIX:-}" ]; then
    echo "DIR_SUFFIX is not set. Exiting."
    exit 1
fi

if [ -z "${NGINX_CONF:-}" ]; then
    echo "NGINX_CONF is not set. Exiting."
    exit 1
fi

APP_NAME="sampay-api"
APP_DIR="/home/ec2-user/sampay/backend"
WORKER_DIR="$APP_DIR"
BLUE_DIR="/home/ec2-user/sampay/backend-blue"
GREEN_DIR="/home/ec2-user/sampay/backend-green"
NEW_DIR="/home/ec2-user/sampay/backend-$DIR_SUFFIX"

BLUE_PORT=8080
GREEN_PORT=8081

if systemctl is-active --quiet "${APP_NAME}-blue"; then
    ACTIVE_ENV="blue"
    DEPLOY_ENV="green"
    DEPLOY_PORT=$GREEN_PORT
    ACTIVE_PORT=$BLUE_PORT
    APP_DIR="$GREEN_DIR"
else
    ACTIVE_ENV="green"
    DEPLOY_ENV="blue"
    DEPLOY_PORT=$BLUE_PORT
    ACTIVE_PORT=$GREEN_PORT
    APP_DIR="$BLUE_DIR"
fi

PREVIOUS_VERSION_LINK=$(readlink -f "$APP_DIR" || true)
if [ "$PREVIOUS_VERSION_LINK" = "$APP_DIR" ]; then
    PREVIOUS_VERSION_LINK=""
fi

echo "Deploying to $DEPLOY_ENV environment on port $DEPLOY_PORT..."

echo "Preparing database..."
export PACKAGE_ROOT="$NEW_DIR"
if ! ( "$NEW_DIR/build/db-create" && "$NEW_DIR/build/db-migrate" && "$NEW_DIR/build/db-seed" ); then
    echo "Error: Database preparation failed. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

echo "Update symlink to new version..."
if ! (ln -sfn "$NEW_DIR" "$APP_DIR" && ln -sfn "$NEW_DIR" "$WORKER_DIR"); then
    echo "Error: Failed to update symlink. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

function rollback() {
    echo "Rolling back to previous version..."
    if [ -n "$PREVIOUS_VERSION_LINK" ]; then
        ln -sfn "$PREVIOUS_VERSION_LINK" "$APP_DIR"
        ln -sfn "$PREVIOUS_VERSION_LINK" "$WORKER_DIR"
    fi
    sudo systemctl restart "${APP_NAME}-${ACTIVE_ENV}" sampay-worker
    rm -rf "$NEW_DIR"
    exit 1
}

echo "Restarting ${APP_NAME}-${DEPLOY_ENV} and sampay-worker services..."
if ! sudo systemctl restart "${APP_NAME}-${DEPLOY_ENV}" sampay-worker; then
    echo "Error: Failed to restart ${APP_NAME}-${DEPLOY_ENV} and sampay-worker services."
    rollback
fi

echo "Waiting for sampay-worker service to become active..."
if timeout 10s bash -c 'until systemctl is-active --quiet sampay-worker; do sleep 1; done'; then
    echo "sampay-worker service is now active."
else
    echo "Timed out waiting for sampay-worker service to become active."
    rollback
fi

echo "Waiting for ${APP_NAME}-${DEPLOY_ENV} service to become active..."
retry_count=0
max_retries=3
retry_interval=5
while [ $retry_count -lt $max_retries ]; do
    if systemctl is-active --quiet sampay-api; then
        echo "${APP_NAME}-${DEPLOY_ENV} service is active."
        if wget -q --spider "http://localhost:${DEPLOY_PORT}/api/health"; then
            echo "Health check passed successfully."
            break
        else
            echo "Health check failed. Retrying..."
        fi
    else
        echo "${APP_NAME}-${DEPLOY_ENV} service is not yet active. Waiting..."
    fi

    retry_count=$((retry_count + 1))

    if [ $retry_count -lt $max_retries ]; then
        echo "Retrying in $retry_interval seconds... (Attempt $retry_count of $max_retries)"
        sleep $retry_interval
    fi

    if [ $retry_count -eq $max_retries ]; then
        echo "Error: Failed to start ${APP_NAME}-${DEPLOY_ENV} service."
        rollback
    fi
done

echo "Updating Nginx to route traffic to port $DEPLOY_PORT..."
sudo sed -i "s/server 127.0.0.1:$ACTIVE_PORT/server 127.0.0.1:$DEPLOY_PORT/" "$NGINX_CONF"
sudo systemctl reload nginx

echo "Stopping previous service: ${APP_NAME}-${ACTIVE_ENV}..."
sudo systemctl stop "${APP_NAME}-${ACTIVE_ENV}"

echo "Removing old version directory: $PREVIOUS_VERSION_LINK"
rm -rf "$PREVIOUS_VERSION_LINK"

echo "Deployment to $DEPLOY_ENV on port $DEPLOY_PORT completed successfully."
