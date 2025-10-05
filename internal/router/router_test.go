package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lumiiam/internal/config"
	"lumiiam/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gorm.DB, http.Handler) {
	t.Helper()
	// in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	// migrate
	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Token{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg := &config.Config{
		AppPort:               0,
		AppEnv:                "test",
		PgHost:                "",
		PgPort:                0,
		PgUser:                "",
		PgPassword:            "",
		PgDB:                  "",
		PgSSLMode:             "",
		AccessTokenTTLMinutes: 15,
		RefreshTokenTTLDays:   7,
		PasswordBcryptCost:    4, // faster in tests
	}

	r := New(cfg, db)
	return db, r
}

func TestCreateUser_Success(t *testing.T) {
	_, r := setupTestRouter(t)

	body := map[string]string{
		"email":    "u1@example.com",
		"username": "user1",
		"password": "pass1234",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["email"] != "u1@example.com" || resp["username"] != "user1" {
		t.Fatalf("unexpected resp: %+v", resp)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	_, r := setupTestRouter(t)

	payload := func(email, username string) *httptest.ResponseRecorder {
		body := map[string]string{
			"email":    email,
			"username": username,
			"password": "pass1234",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	w1 := payload("dup@example.com", "dupuser")
	if w1.Code != http.StatusCreated {
		t.Fatalf("first create expected 201, got %d", w1.Code)
	}
	w2 := payload("dup@example.com", "dupuser2")
	if w2.Code == http.StatusCreated {
		t.Fatalf("duplicate email should not be 201")
	}
	w3 := payload("dup2@example.com", "dupuser")
	if w3.Code == http.StatusCreated {
		t.Fatalf("duplicate username should not be 201")
	}
}
