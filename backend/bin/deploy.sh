#!/usr/bin/env bash
set -euox pipefail

cd ~/sampay
git pull

cd backend
make build db-prepare
sudo systemctl restart sampay-api
sudo systemctl restart sampay-worker

if timeout 10s bash -c 'until systemctl is-active --quiet sampay-worker; do sleep 1; done'; then
    echo "sampay-worker service is now active."
else
    echo "Timed out waiting for sampay-worker service to become active."
    exit 1
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
            exit 0
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
done

echo "Failed to verify sampay-api service health after $max_retries attempts."
exit 1
