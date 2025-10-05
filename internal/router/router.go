package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lumiiam/internal/config"
	"lumiiam/internal/handlers"
	"lumiiam/internal/middleware"
	"lumiiam/internal/services"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// serve static test frontend
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/web/index.html") })

	auth_svc := services.NewAuthService(db, cfg)
	user_svc := services.NewUserService(db)

	auth_mw := middleware.NewAuthMiddleware(auth_svc)
	auth_h := handlers.NewAuthHandler(auth_svc)
	user_h := handlers.NewUserHandler(user_svc)

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		// auth endpoints (kebab-case, no underscores/camelCase)
		api.POST("/auth/login", auth_h.PostLogin)
		api.POST("/auth/refresh", auth_h.PostRefresh)
		api.POST("/auth/logout", auth_mw.RequireAuth(), auth_h.PostLogout)

		// user endpoints
		api.GET("/users/me", auth_mw.RequireAuth(), user_h.GetMe)
		api.GET("/users", auth_mw.RequireAuth(), user_h.GetUsers)
	}

	return r
}
