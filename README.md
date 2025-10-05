# lumiiam

A minimal IAM service for user/login/permission management built with Go 1.20, PostgreSQL, Gorm, and a small Vue 3 test page.

- REST style APIs under `/api/v1`
- URL rules use lowercase and kebab-case (no underscores, no camelCase)
- No JWT. Uses opaque Access Token + Refresh Token stored in DB
- Variable names in payloads avoid camelCase (use snake_case in JSON)

## stack

- Go 1.20
- Gin (HTTP server)
- Gorm + PostgreSQL
- Vue 3 (CDN) test page under `/web`

## structure

```
cmd/server/main.go
internal/config/
internal/db/
internal/models/
internal/services/
internal/handlers/
internal/middleware/
internal/router/
web/
```

## setup

1. Create `.env` from example:

```
cp .env.example .env
```

2. Ensure PostgreSQL is running and database exists:

```
createdb lumiiam
```

3. Fetch deps:

```
go mod tidy
```

4. Run server (local without Docker):

```
go run ./cmd/server
```

## docker compose (recommended)

Build and start PostgreSQL + app:

```
docker compose up -d --build
```

Open: http://localhost:8080/

Stop and remove:

```
docker compose down -v
```

## makefile targets

Common tasks:

```
make tidy         # fetch deps
make build        # build local binary to bin/lumiiam
make run          # run locally without Docker
make docker-build # build Docker image lumiiam:latest
make up           # docker compose up -d --build
make down         # docker compose down -v
make logs         # tail app logs
```

## scripts

Helper scripts under `scripts/`:

```
./scripts/build.sh  # tidy, build binary, build image
./scripts/up.sh     # docker compose up -d --build
./scripts/down.sh   # docker compose down -v
```

If needed, make them executable:

```
chmod +x scripts/*.sh
```

4. Run server:

```
go run ./cmd/server
```

Server starts at `http://localhost:8080`.

- Visit `http://localhost:8080/` to open the test page.
- Default seeded admin user: `admin@example.com` / `admin123`.

## api

Base: `/api/v1`

- POST `/auth/login`
  - body: `{ "identifier": "email-or-username", "password": "..." }`
  - returns: `{ "access_token", "refresh_token", "user_id" }`

- POST `/auth/refresh`
  - body: `{ "refresh_token": "..." }`
  - returns new tokens

- POST `/auth/logout`
  - header: `Authorization: Bearer <access_token>`

- GET `/users/me`
  - header: `Authorization: Bearer <access_token>`

- GET `/users?limit=&offset=`
  - header: `Authorization: Bearer <access_token>`

## notes

- Tokens are opaque random strings. Only their SHA-256 hashes are stored in DB with kind and expiry.
- Access tokens are short-lived, refresh tokens are longer-lived. Both can be revoked.
- JSON fields avoid camelCase and use snake_case via struct tags.
- URLs avoid underscores and camelCase, using kebab-case segments.

## development

- Adjust timeouts and bcrypt cost in `.env`.
- Extend models for roles/permissions. Create role/permission handlers and middleware for RBAC checks.
# lumiiam
