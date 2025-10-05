#!/usr/bin/env sh
set -e

echo "[lumiiam] building and starting services via docker compose..."
docker compose up -d --build

echo "[lumiiam] services are starting. View logs with: ./scripts/logs.sh"
