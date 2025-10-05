#!/usr/bin/env sh
set -e

# This script provisions a standalone PostgreSQL container matching the provided spec.
# It writes .pg.env (not .env) to avoid clobbering app env variables.

PG_ENV_FILE=.pg.env

cat > "$PG_ENV_FILE" <<EOF
# PostgreSQL 环境变量
POSTGRES_USER=admin
POSTGRES_PASSWORD=randompass
POSTGRES_DB=diam
EOF

echo "[lumiiam] starting postgres with env-file $PG_ENV_FILE ..."

docker run -d \
  --name postgres \
  --restart=always \
  --env-file "$PG_ENV_FILE" \
  -p 5432:5432 \
  -v pg_data:/var/lib/postgresql/data \
  postgres:16.4

echo "[lumiiam] postgres is starting. View logs with: ./scripts/postgres-logs.sh"
