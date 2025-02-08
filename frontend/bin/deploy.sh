#!/usr/bin/env bash
set -euox pipefail

if [ -z "${DIR_SUFFIX:-}" ]; then
    echo "DIR_SUFFIX is not set. Exiting."
    exit 1
fi

APP_NAME="sampay-frontend"
APP_DIR="/home/ec2-user/sampay/frontend"
BLUE_DIR="/home/ec2-user/sampay/frontend-blue"
GREEN_DIR="/home/ec2-user/sampay/frontend-green"

BLUE_PORT=3000
GREEN_PORT=3001

if systemctl is-active --quiet "${APP_NAME}-blue"; then
    ACTIVE_ENV="blue"
    DEPLOY_ENV="green"
    DEPLOY_PORT=$GREEN_PORT
    ACTIVE_PORT=$BLUE_PORT
    NEW_DIR="$GREEN_DIR"
else
    ACTIVE_ENV="green"
    DEPLOY_ENV="blue"
    DEPLOY_PORT=$BLUE_PORT
    ACTIVE_PORT=$GREEN_PORT
    NEW_DIR="$BLUE_DIR"
fi

PREVIOUS_VERSION_LINK=$(readlink -f "$APP_DIR" || echo "")

echo "Deploying to $DEPLOY_ENV environment on port $DEPLOY_PORT..."

mkdir -p "$NEW_DIR"
cd "$NEW_DIR" || exit 1

echo "Building application..."
if ! ( npm ci && npm run build ); then
    echo "Error: Build failed. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

echo "Updating symlink to point to $DEPLOY_ENV..."
if ! ln -sfn "$NEW_DIR" "$APP_DIR"; then
    echo "Error: Failed to update symlink. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

function rollback() {
    echo "Rolling back to previous version..."
    if [ -n "$PREVIOUS_VERSION_LINK" ]; then
        ln -sfn "$PREVIOUS_VERSION_LINK" "$APP_DIR"
        sudo systemctl restart "${APP_NAME}-${ACTIVE_ENV}"
    fi
    rm -rf "$NEW_DIR"
    exit 1
}

echo "Restarting ${APP_NAME}-${DEPLOY_ENV} service..."
if ! sudo systemctl restart "${APP_NAME}-${DEPLOY_ENV}"; then
    echo "Error: Failed to restart ${APP_NAME}-${DEPLOY_ENV} service."
    rollback
fi

echo "Waiting for ${APP_NAME}-${DEPLOY_ENV} service to become active on port $DEPLOY_PORT..."
retry_count=0
max_retries=3
retry_interval=5

while [ $retry_count -lt $max_retries ]; do
    if systemctl is-active --quiet "${APP_NAME}-${DEPLOY_ENV}"; then
        echo "${APP_NAME}-${DEPLOY_ENV} service is active."
        if wget -q --spider "http://localhost:${DEPLOY_PORT}/health"; then
            echo "Health check passed successfully on port $DEPLOY_PORT."
            break
        else
            echo "Health check failed on port $DEPLOY_PORT. Retrying..."
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
        echo "Error: Health check failed after $max_retries attempts on port $DEPLOY_PORT."
        rollback
    fi
done

echo "Updating Nginx to route traffic to port $DEPLOY_PORT..."
sudo sed -i "s/server 127.0.0.1:$ACTIVE_PORT\+/server 127.0.0.1:$DEPLOY_PORT/" "$NGINX_CONF"
sudo systemctl reload nginx

echo "Stopping previous service: ${APP_NAME}-${ACTIVE_ENV}..."
sudo systemctl stop "${APP_NAME}-${ACTIVE_ENV}"

echo "Removing old version directory: $PREVIOUS_VERSION_LINK"
rm -rf "$PREVIOUS_VERSION_LINK"

echo "Deployment to $DEPLOY_ENV on port $DEPLOY_PORT completed successfully."
