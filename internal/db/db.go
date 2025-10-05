package db

import (
	"log"

	"lumiiam/internal/config"
	"lumiiam/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Token{},
	); err != nil {
		return nil, err
	}
	seed_admin(db)
	return db, nil
}

func seed_admin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}
	log.Println("seeding admin user: admin@example.com / admin123")
	if err := models.CreateUser(db, "admin@example.com", "admin", "admin123"); err != nil {
		log.Printf("seed admin failed: %v", err)
	}
}
