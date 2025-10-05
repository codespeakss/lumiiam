# build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git build-base
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/lumiiam ./cmd/server

# runtime stage
FROM alpine:3.20
WORKDIR /app
ENV app_port=8080 \
    app_env=prod \
    pg_host=postgres \
    pg_port=5432 \
    pg_user=postgres \
    pg_password=postgres \
    pg_db=lumiiam \
    pg_sslmode=disable \
    access_token_ttl_minutes=15 \
    refresh_token_ttl_days=7 \
    password_bcrypt_cost=12
COPY --from=builder /out/lumiiam /app/lumiiam
COPY web /app/web
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/lumiiam"]
