package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lumiiam/internal/service"
	"lumiiam/pkg/cache"
)

type Handler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewHandler(db *gorm.DB, cache *cache.RedisTokenStore) *Handler {
	userService := service.NewUserService(
		db,
		cache,
	)
	userService.InitServiceData()
	tokenService := service.NewTokenService(
		db,
		cache,
	)

	return &Handler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// v1 user
	userV1 := r.Group("api/v1/users")
	{
		userV1.POST("", h.postUser)
		userV1.GET("/:id", h.getUserByID)
		userV1.PUT("/:id", h.putUserByID)
	}

	tokenV1 := r.Group("api/v1/tokens")
	{
		tokenV1.POST("", h.PostToken)
		tokenV1.GET("", h.GetToken)
		tokenV1.POST("/validate", h.PostValidateToken)
	}
}
