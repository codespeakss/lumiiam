package service

import (
	"fmt"
	"gorm.io/gorm"
	"lumiiam/api"
	"lumiiam/config"
	"lumiiam/internal/repo"
	"lumiiam/pkg/cache"
	"lumiiam/pkg/util"
	"time"
)

type TokenService struct {
	repo  *repo.UserRepo
	redis *cache.RedisTokenStore
}

func NewTokenService(db *gorm.DB, redis *cache.RedisTokenStore) *TokenService {
	return &TokenService{
		repo:  repo.NewUserRepo(db),
		redis: redis,
	}
}

func (s *TokenService) CreateToken(item *api.PostTokenReq) (*api.PostTokenResp, error) {
	// check valid
	userInfo, e := s.repo.CheckUserPass(item.Name, item.Password)
	if e != nil {
		return nil, fmt.Errorf("s.repo.CheckUserPass: %w", e)
	}
	if userInfo != nil {
		refreshToken, e := util.GenerateToken()
		if e != nil {
			return nil, fmt.Errorf("")
		}

		e = s.redis.Set(refreshToken, userInfo.Id, config.RefreshTokenTimeoutSecond*time.Second)

		accessToken, e := util.GenerateToken()
		if e != nil {
			return nil, fmt.Errorf("")
		}
		e = s.redis.Set(accessToken, userInfo.Id, config.AccessTokenTimeoutSecond*time.Second)

		return &api.PostTokenResp{
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
			Name:         userInfo.Name,
			ExpiresAt:    time.Now().UnixMilli() + config.RefreshTokenTimeoutSecond*1000,
		}, nil
	}

	return nil, e
}

func (s *TokenService) GetTokenInfo(tokenReq *api.ValidateTokenReq) (*api.ValidateTokenResp, error) {
	userId, ok := s.redis.Get(tokenReq.Token)
	if ok {
		return &api.ValidateTokenResp{
			Id:   userId,
			Name: "",
		}, nil
	}

	return nil, fmt.Errorf("s.redis.Get !ok ")
}
