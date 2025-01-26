#!/usr/bin/env bash
set -euo pipefail

cd ~/sampay
git pull

cd frontend
sudo systemctl restart sampay-frontend
