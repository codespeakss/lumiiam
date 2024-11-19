package repo

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"lumiiam/internal/model"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateItem(item model.User) (*model.User, error) {
	result := r.db.Create(&item)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user item: %w", result.Error)
	}
	return &item, nil
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var user model.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepo) CheckUserPass(name string, pass string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to query user: %w", err)
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (r *UserRepo) buildQueryByWhereConditions(ctx context.Context, modelType interface{}, whereConditions map[string]interface{}) *gorm.DB {
	query := r.db.WithContext(ctx).Model(modelType)
	for key, value := range whereConditions {
		query = query.Where(key, value)
	}
	return query
}

func (r *UserRepo) UpdateItemsByWhereConditions(ctx context.Context, modelType interface{}, whereConditions map[string]interface{}, fieldsToUpdate map[string]interface{}) (int64, error) {
	result := r.buildQueryByWhereConditions(ctx, modelType, whereConditions).Updates(fieldsToUpdate)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to update items by where conditions: %w", result.Error)
	}
	return result.RowsAffected, nil
}

func (r *UserRepo) InitData() error {
	var count int64
	if err := r.db.Model(model.User{}).Count(&count).Error; err != nil {
		fmt.Println("Error checking user table:", err)
		return nil
	}
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		defaultUser := model.User{
			Id:       "u000000",
			Name:     "admin",
			Email:    "x@y.com",
			CreateAt: time.Now().Unix(),
			Password: string(hashedPassword),
		}
		if err := r.db.Create(&defaultUser).Error; err != nil {
			fmt.Println("Error creating default user:", err)
		} else {
			fmt.Println("Default user created successfully.")
		}
	}
	return nil
}
