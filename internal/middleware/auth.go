package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"lumiiam/internal/services"
)

type AuthMiddleware struct { auth *services.AuthService }

func NewAuthMiddleware(auth *services.AuthService) *AuthMiddleware { return &AuthMiddleware{auth: auth} }

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"}); return }
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization"}); return }
		uid, err := m.auth.ValidateAccess(parts[1])
		if err != nil { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()}); return }
		c.Set("user_id", uid)
		c.Set("access_token", parts[1])
		c.Next()
	}
}
