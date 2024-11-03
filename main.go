package main

import (
	"github.com/gin-gonic/gin"
	"lumiiam/api"
	"net/http"
)

type PostTokenReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var userDB = map[string]string{
	"admin": "password",
}

func PostToken(c *gin.Context) {
	var user PostTokenReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:   0,
			Msg:    "",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
		return
	}

	// 验证用户名和密码
	if password, exists := userDB[user.Username]; exists {
		if password == user.Password {
			// 生成 token，这里可以替换为实际的 token 生成逻辑
			c.JSON(http.StatusOK, api.HttpResp{
				Code:   200,
				Msg:    "登录成功",
				Data:   map[string]string{"token": DefaultValidToken, "timeout": "3600"}, // 将 token 放入 Data 字段
				Total:  0,
				Errors: nil,
			})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, api.HttpResp{
		Code:   0,
		Msg:    "",
		Data:   nil,
		Total:  0,
		Errors: nil,
	})
}

func GetToken(c *gin.Context) {
	var user PostTokenReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:   0,
			Msg:    "",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
		return
	}

	if password, exists := userDB[user.Username]; exists {
		if password == user.Password {
			c.JSON(http.StatusOK, api.HttpResp{
				Code:   0,
				Msg:    "",
				Data:   nil,
				Total:  0,
				Errors: nil,
			})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, api.HttpResp{
		Code:   0,
		Msg:    "",
		Data:   nil,
		Total:  0,
		Errors: nil,
	})
}

// 假设我们有一个函数来验证 token 的有效性
func isValidToken(token string) bool {
	// 这里实现你的 token 验证逻辑
	// 示例：验证 token 是否与某个预定义值匹配
	return token == DefaultValidToken // 这里用你实际的 token 验证逻辑替换
}

// Token 验证请求体
type ValidateTokenReq struct {
	Token string `json:"token" binding:"required"`
}

func PostValidateToken(c *gin.Context) {
	var req ValidateTokenReq

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:  400,
			Msg:   "请求参数错误",
			Data:  nil,
			Total: 0,
			Errors: []api.ErrorDetail{api.ErrorDetail{
				Msg: err.Error(),
			},
			},
		})
		return
	}

	// 验证 token
	if isValidToken(req.Token) {
		c.JSON(http.StatusOK, api.HttpResp{
			Code:   200,
			Msg:    "Token 有效",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
	} else {
		c.JSON(http.StatusUnauthorized, api.HttpResp{
			Code:   401,
			Msg:    "Token 无效",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
	}
}

const DefaultValidToken = "token-abc"

func main() {
	router := gin.Default()

	router.POST("/api/v1/tokens", PostToken)
	router.GET("/api/v1/tokens", GetToken)
	router.POST("/api/v1/tokens/validate", PostValidateToken)

	router.Run(":8011")
}
