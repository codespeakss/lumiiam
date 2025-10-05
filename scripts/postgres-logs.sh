#!/usr/bin/env sh
set -e

docker logs -f --tail=200 postgres
