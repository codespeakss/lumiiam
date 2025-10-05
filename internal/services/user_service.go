package services

import (
	"errors"

	"lumiiam/internal/models"
	"gorm.io/gorm"
)

type UserService struct { db *gorm.DB }

func NewUserService(db *gorm.DB) *UserService { return &UserService{db: db} }

func (s *UserService) Create(email, username, password string) (*models.User, error) {
	if err := models.CreateUser(s.db, email, username, password); err != nil { return nil, err }
	var u models.User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil { return nil, err }
	return &u, nil
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := s.db.First(&u, id).Error; err != nil { return nil, err }
	return &u, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := s.db.Where("email = ?", email).First(&u).Error; err != nil { return nil, err }
	return &u, nil
}

func (s *UserService) List(limit, offset int) ([]models.User, int64, error) {
	if limit <= 0 { limit = 20 }
	var users []models.User
	var total int64
	s.db.Model(&models.User{}).Count(&total)
	if err := s.db.Limit(limit).Offset(offset).Order("id asc").Find(&users).Error; err != nil { return nil, 0, err }
	return users, total, nil
}

func (s *UserService) UpdatePassword(id uint, new_password string) error {
	var u models.User
	if err := s.db.First(&u, id).Error; err != nil { return err }
	if new_password == "" { return errors.New("empty password") }
	return models.CreateUserPasswordUpdate(s.db, &u, new_password)
}
