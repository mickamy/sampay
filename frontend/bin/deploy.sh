#!/usr/bin/env bash
set -euox pipefail

cd ~/sampay
git pull

cd frontend
sudo systemctl restart sampay-frontend
