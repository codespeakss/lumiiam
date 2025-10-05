APP_NAME=lumiiam
APP_BIN=bin/$(APP_NAME)
GOFLAGS=
BASE?=http://localhost:8080
CONCURRENCY?=20
DURATION?=30s
IDENTIFIER?=admin@example.com
PASSWORD?=admin123

.PHONY: all tidy build run clean docker-build up down logs bench pg-up pg-down pg-logs up-external

all: tidy build

 tidy:
	go mod tidy

 build:
	GO111MODULE=on go build $(GOFLAGS) -o $(APP_BIN) ./cmd/server

 run:
	go run ./cmd/server

 clean:
	rm -rf bin

 docker-build:
	docker build -t $(APP_NAME):latest .

 up:
	docker compose up -d --build

 down:
	docker compose down -v

 logs:
	docker compose logs -f --tail=200 app

 bench:
	go run ./cmd/bench \
		--base $(BASE) \
		--concurrency $(CONCURRENCY) \
		--duration $(DURATION) \
		--identifier $(IDENTIFIER) \
		--password $(PASSWORD)

 pg-up:
	./scripts/postgres-up.sh

 pg-down:
	./scripts/postgres-down.sh

 pg-logs:
	./scripts/postgres-logs.sh

 up-external:
	docker compose -f docker-compose.external.yml up -d --build
