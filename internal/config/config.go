package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                 int
	AppEnv                  string
	PgHost                  string
	PgPort                  int
	PgUser                  string
	PgPassword              string
	PgDB                    string
	PgSSLMode               string
	AccessTokenTTLMinutes   int
	RefreshTokenTTLDays     int
	PasswordBcryptCost      int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	get_int := func(key string, def int, fallbacks ...string) int {
		v := getenv_multi(def_str(def), key, fallbacks...)
		i, err := strconv.Atoi(v)
		if err != nil { return def }
		return i
	}

	cfg := &Config{
		AppPort:               get_int("app_port", 8080),
		AppEnv:                getenv("app_env", "dev"),
		PgHost:                getenv_multi("127.0.0.1", "pg_host", "POSTGRES_HOST"),
		PgPort:                get_int("pg_port", 5432, "POSTGRES_PORT"),
		PgUser:                getenv_multi("admin", "pg_user", "POSTGRES_USER"),
		PgPassword:            getenv_multi("randompass", "pg_password", "POSTGRES_PASSWORD"),
		PgDB:                  getenv_multi("diam", "pg_db", "POSTGRES_DB"),
		PgSSLMode:             getenv_multi("disable", "pg_sslmode", "POSTGRES_SSLMODE"),
		AccessTokenTTLMinutes: get_int("access_token_ttl_minutes", 15),
		RefreshTokenTTLDays:   get_int("refresh_token_ttl_days", 7),
		PasswordBcryptCost:    get_int("password_bcrypt_cost", 12),
	}
	return cfg, nil
}

func (c *Config) Addr() string { return fmt.Sprintf(":%d", c.AppPort) }

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.PgHost, c.PgUser, c.PgPassword, c.PgDB, c.PgPort, c.PgSSLMode,
	)
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenv_multi(def string, primary string, fallbacks ...string) string {
	if v := os.Getenv(primary); v != "" { return v }
	for _, k := range fallbacks {
		if v := os.Getenv(k); v != "" { return v }
	}
	return def
}

func def_str(i int) string { return strconv.Itoa(i) }
