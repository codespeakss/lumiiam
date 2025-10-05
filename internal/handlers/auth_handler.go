package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lumiiam/internal/services"
)

type AuthHandler struct { auth *services.AuthService }

func NewAuthHandler(auth *services.AuthService) *AuthHandler { return &AuthHandler{auth: auth} }

func (h *AuthHandler) PostLogin(c *gin.Context) {
	type req struct { Identifier string `json:"identifier"`; Password string `json:"password"` }
	var body req
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	res, err := h.auth.Login(body.Identifier, body.Password)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) PostRefresh(c *gin.Context) {
	type req struct { RefreshToken string `json:"refresh_token"` }
	var body req
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
	res, err := h.auth.Refresh(body.RefreshToken)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) PostLogout(c *gin.Context) {
	tok := c.GetString("access_token")
	if tok == "" { c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"}); return }
	if err := h.auth.Logout(tok); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
