package service

import (
	"gorm.io/gorm"
	"lumiiam/internal/model"
	"lumiiam/internal/repo"
	"lumiiam/pkg/cache"
)

type UserService struct {
	repo  *repo.UserRepo
	redis *cache.RedisTokenStore
}

func NewUserService(db *gorm.DB, redis *cache.RedisTokenStore) *UserService {
	return &UserService{
		repo:  repo.NewUserRepo(db),
		redis: redis,
	}
}

func (s *UserService) CreateUser(item model.User) (*model.User, error) {
	//item.Id = idgen.GetIdWithPref("user")
	//item.CreateAt = time.Now().UnixMilli()
	//log.Println("item: ", item)

	return s.repo.CreateItem(item)
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) InitServiceData() error {
	s.repo.InitData()
	return nil
}
