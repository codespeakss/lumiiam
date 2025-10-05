package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lumiiam/internal/services"
)

type UserHandler struct { users *services.UserService }

func NewUserHandler(users *services.UserService) *UserHandler { return &UserHandler{users: users} }

func (h *UserHandler) GetMe(c *gin.Context) {
	uid := c.GetUint("user_id")
	u, err := h.users.GetByID(uid)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.users.List(limit, offset)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}
