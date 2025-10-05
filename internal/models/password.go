package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserPasswordUpdate(db *gorm.DB, u *User, new_password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil { return err }
	return db.Model(u).Update("password", string(hash)).Error
}
