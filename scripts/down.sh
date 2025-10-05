#!/usr/bin/env sh
set -e

echo "[lumiiam] stopping and removing services..."
docker compose down -v

echo "[lumiiam] done."
