package service

import (
	"fmt"
	"gorm.io/gorm"
	"lumiiam/api"
	"lumiiam/config"
	"lumiiam/internal/repo"
	"lumiiam/pkg/cache"
	"lumiiam/pkg/util"
	"strings"
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

		e = s.redis.Set(refreshToken, userInfo.Id+", "+userInfo.Name, config.RefreshTokenTimeoutSecond*time.Second)

		accessToken, e := util.GenerateToken()
		if e != nil {
			return nil, fmt.Errorf("")
		}
		e = s.redis.Set(accessToken, userInfo.Id+", "+userInfo.Name, config.AccessTokenTimeoutSecond*time.Second)

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
	// 从 Redis 获取 token 对应的值
	userInfo, ok := s.redis.Get(tokenReq.Token)
	if !ok {
		return nil, fmt.Errorf("token not found or expired")
	}

	// 解析 Redis 存储的用户信息，格式为 "userId, userName"
	parts := strings.SplitN(userInfo, ", ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid userInfo format in Redis")
	}

	userId := parts[0]
	userName := parts[1]

	return &api.ValidateTokenResp{
		Id:   userId,
		Name: userName,
	}, nil
}
