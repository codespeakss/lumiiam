#!/usr/bin/env sh
set -e

echo "[lumiiam] go mod tidy..."
go mod tidy

echo "[lumiiam] building local binary to bin/lumiiam..."
mkdir -p bin
GO111MODULE=on go build -o bin/lumiiam ./cmd/server

echo "[lumiiam] building docker image lumiiam:latest..."
docker build -t lumiiam:latest .

echo "[lumiiam] done."
