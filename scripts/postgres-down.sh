#!/usr/bin/env sh
set -e

if docker ps -a --format '{{.Names}}' | grep -q '^postgres$'; then
  echo "[lumiiam] stopping and removing postgres container..."
  docker rm -f postgres >/dev/null 2>&1 || true
else
  echo "[lumiiam] postgres container not found."
fi

# keep the named volume for persistence; to wipe data uncomment below
# docker volume rm pg_data || true

echo "[lumiiam] done."
