package services

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"lumiiam/internal/config"
	"lumiiam/internal/models"

	"gorm.io/gorm"
)

type AuthService struct {
	db   *gorm.DB
	cfg  *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
}

func (s *AuthService) Login(identifier, password string) (*LoginResult, error) {
	var user models.User
	if err := s.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := models.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	access, refresh, err := s.issue_tokens(user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{AccessToken: access, RefreshToken: refresh, UserID: user.ID}, nil
}

func (s *AuthService) issue_tokens(user_id uint) (string, string, error) {
	access_plain, err := models.NewOpaqueToken(32)
	if err != nil { return "", "", err }
	refresh_plain, err := models.NewOpaqueToken(32)
	if err != nil { return "", "", err }

	access_hash := s.hash_token(access_plain)
	refresh_hash := s.hash_token(refresh_plain)

	now := time.Now()
	access_exp := now.Add(time.Duration(s.cfg.AccessTokenTTLMinutes) * time.Minute)
	refresh_exp := now.Add(time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour)

	tx := s.db.Begin()
	if err := tx.Create(&models.Token{UserID: user_id, Kind: models.TokenKindAccess, Hash: access_hash, ExpiresAt: access_exp}).Error; err != nil { tx.Rollback(); return "", "", err }
	if err := tx.Create(&models.Token{UserID: user_id, Kind: models.TokenKindRefresh, Hash: refresh_hash, ExpiresAt: refresh_exp}).Error; err != nil { tx.Rollback(); return "", "", err }
	if err := tx.Commit().Error; err != nil { return "", "", err }
	return access_plain, refresh_plain, nil
}

func (s *AuthService) Refresh(refresh_token string) (*LoginResult, error) {
	h := s.hash_token(refresh_token)
	var t models.Token
	if err := s.db.Where("hash = ? AND kind = ?", h, models.TokenKindRefresh).First(&t).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if t.RevokedAt.Valid || time.Now().After(t.ExpiresAt) {
		return nil, errors.New("refresh token expired or revoked")
	}
	access, refresh, err := s.issue_tokens(t.UserID)
	if err != nil { return nil, err }
	return &LoginResult{AccessToken: access, RefreshToken: refresh, UserID: t.UserID}, nil
}

func (s *AuthService) Logout(access_token string) error {
	h := s.hash_token(access_token)
	return s.db.Model(&models.Token{}).Where("hash = ? AND kind = ?", h, models.TokenKindAccess).Update("revoked_at", sql.NullTime{Time: time.Now(), Valid: true}).Error
}

func (s *AuthService) hash_token(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return fmt_bytes(sum[:])
}

func fmt_bytes(b []byte) string {
	const hextable = "0123456789abcdef"
	res := make([]byte, len(b)*2)
	for i, v := range b {
		res[i*2] = hextable[v>>4]
		res[i*2+1] = hextable[v&0x0f]
	}
	return string(res)
}

func (s *AuthService) ValidateAccess(access_token string) (uint, error) {
	h := s.hash_token(access_token)
	var t models.Token
	if err := s.db.Where("hash = ? AND kind = ?", h, models.TokenKindAccess).First(&t).Error; err != nil {
		return 0, errors.New("invalid access token")
	}
	if t.RevokedAt.Valid || time.Now().After(t.ExpiresAt) {
		return 0, errors.New("access token expired or revoked")
	}
	return t.UserID, nil
}
