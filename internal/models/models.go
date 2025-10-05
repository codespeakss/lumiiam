package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:128" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index"`
	RoleID uint `gorm:"index"`
}

type RolePermission struct {
	ID           uint `gorm:"primaryKey"`
	RoleID       uint `gorm:"index"`
	PermissionID uint `gorm:"index"`
}

type TokenKind string

const (
	TokenKindAccess  TokenKind = "access"
	TokenKindRefresh TokenKind = "refresh"
)

type Token struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	UserID    uint         `gorm:"index" json:"user_id"`
	Kind      TokenKind    `gorm:"size:16;index" json:"kind"`
	Hash      string       `gorm:"size:128;index" json:"-"`
	ExpiresAt time.Time    `gorm:"index" json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
	CreatedAt time.Time    `json:"created_at"`
}

func CreateUser(db *gorm.DB, email, username, password string) error {
	if email == "" || password == "" || username == "" {
		return errors.New("missing fields")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := &User{Email: email, Username: username, Password: string(hash)}
	return db.Create(u).Error
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewOpaqueToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
