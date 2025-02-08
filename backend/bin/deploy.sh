#!/usr/bin/env bash
set -euox pipefail

if [ -z "${DIR_SUFFIX:-}" ]; then
    echo "DIR_SUFFIX is not set. Exiting."
    exit 1
fi

NEW_DIR="/home/ec2-user/sampay/backend-$DIR_SUFFIX"

APP_DIR="/home/ec2-user/sampay/backend"
PREVIOUS_VERSION_LINK=$(readlink -f "$APP_DIR" || echo "")

if [ "$PREVIOUS_VERSION_LINK" = "$APP_DIR" ]; then
    PREVIOUS_VERSION_LINK=""
fi


echo "Preparing database..."
export PACKAGE_ROOT="$NEW_DIR"
if ! ( "$NEW_DIR/build/db-create" && "$NEW_DIR/build/db-migrate" && "$NEW_DIR/build/db-seed" ); then
    echo "Error: Database preparation failed. Exiting."
#    rm -rf "$NEW_DIR"
    exit 1
fi

echo "Update symlink to new version..."
if ! ln -sfn "$NEW_DIR" "$APP_DIR"; then
    echo "Error: Failed to update symlink. Exiting."
#    rm -rf "$NEW_DIR"
    exit 1
fi

function rollback() {
    echo "Rolling back to previous version..."
    if [ -n "$PREVIOUS_VERSION_LINK" ]; then
        ln -sfn "$PREVIOUS_VERSION_LINK" "$APP_DIR"
    fi
    sudo systemctl restart sampay-api sampay-worker
#    rm -rf "$NEW_DIR"
    exit 1
}

echo "Restarting sampay-api and sampay-worker services..."
if ! sudo systemctl restart sampay-api sampay-worker; then
    echo "Error: Failed to restart sampay-api and sampay-worker services."
    rollback
fi

echo "Waiting for sampay-worker service to become active..."
if timeout 10s bash -c 'until systemctl is-active --quiet sampay-worker; do sleep 1; done'; then
    echo "sampay-worker service is now active."
else
    echo "Timed out waiting for sampay-worker service to become active."
    rollback
fi

echo "Waiting for sampay-api service to become active..."
retry_count=0
max_retries=3
retry_interval=5
while [ $retry_count -lt $max_retries ]; do
    if systemctl is-active --quiet sampay-api; then
        echo "sampay-api service is active."
        if wget -q --spider http://localhost:8080/api/health; then
            echo "Health check passed successfully."
            break
        else
            echo "Health check failed. Retrying..."
        fi
    else
        echo "sampay-api service is not yet active. Waiting..."
    fi

    retry_count=$((retry_count + 1))

    if [ $retry_count -lt $max_retries ]; then
        echo "Retrying in $retry_interval seconds... (Attempt $retry_count of $max_retries)"
        sleep $retry_interval
    fi

    if [ $retry_count -eq $max_retries ]; then
        echo "Error: Failed to start sampay-api service."
        rollback
    fi
done

echo "Deployment completed successfully."
echo "Removing old version directory: $PREVIOUS_VERSION_LINK"
rm -rf "$PREVIOUS_VERSION_LINK"
