#!/usr/bin/env bash
set -euox pipefail

if [ -z "${DIR_SUFFIX:-}" ]; then
    echo "DIR_SUFFIX is not set. Exiting."
    exit 1
fi

NEW_DIR="/home/ec2-user/sampay/frontend-$DIR_SUFFIX"

APP_DIR="/home/ec2-user/sampay/frontend"
PREVIOUS_VERSION_LINK=$(readlink -f "$APP_DIR" || echo "")

if [ "$PREVIOUS_VERSION_LINK" = "$APP_DIR" ]; then
    PREVIOUS_VERSION_LINK=""
fi

echo "Building application..."
cd "$NEW_DIR" || exit 1
if ! ( npm ci && npm run build ); then
    echo "Error: Build failed. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

echo "Update symlink to new version..."
if ! ln -sfn "$NEW_DIR" "$APP_DIR"; then
    echo "Error: Failed to update symlink. Exiting."
    rm -rf "$NEW_DIR"
    exit 1
fi

function rollback() {
    echo "Rolling back to previous version..."
    if [ -n "$PREVIOUS_VERSION_LINK" ]; then
        ln -sfn "$PREVIOUS_VERSION_LINK" "$APP_DIR"
    fi
    sudo systemctl restart sampay-api sampay-worker
    rm -rf "$NEW_DIR"
    exit 1
}

echo "Restarting sampay-frontend service..."
if ! sudo systemctl restart sampay-frontend; then
    echo "Error: Failed to restart sampay-frontend service."
    rollback
fi

echo "Waiting for sampay-frontend service to become active..."
if timeout 10s bash -c 'until systemctl is-active --quiet sampay-frontend; do sleep 1; done'; then
    echo "sampay-frontend service is now active."
else
    echo "Timed out waiting for sampay-frontend service to become active."
    rollback
fi

echo "Deployment completed successfully."
echo "Removing old version directory: $PREVIOUS_VERSION_LINK"
rm -rf "$PREVIOUS_VERSION_LINK"
