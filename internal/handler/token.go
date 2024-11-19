package handler

import (
	"github.com/gin-gonic/gin"
	"lumiiam/api"
	"net/http"
)

func (h *Handler) PostToken(c *gin.Context) {
	var postTokenReq api.PostTokenReq
	if err := c.ShouldBindJSON(&postTokenReq); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:   0,
			Msg:    "",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
		return
	}

	postTokenResp, e := h.tokenService.CreateToken(&postTokenReq)
	if e != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:   0,
			Msg:    "",
			Total:  0,
			Errors: []api.ErrorDetail{},
		})
		return
	}

	c.JSON(http.StatusOK, api.HttpResp{
		Code:   0,
		Msg:    "",
		Data:   postTokenResp,
		Total:  0,
		Errors: nil,
	})
	return
}

func (h *Handler) GetToken(c *gin.Context) {

}

func (h *Handler) DeleteToken(c *gin.Context) {
	var req api.DeleteTokenReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:  400,
			Msg:   "请求参数错误",
			Data:  nil,
			Total: 0,
			Errors: []api.ErrorDetail{
				api.ErrorDetail{
					Msg: err.Error(),
				},
			},
		})
		return
	}

	err := h.tokenService.DeleteToken(&req)
	if err != nil {
		c.JSON(http.StatusNotImplemented, api.HttpResp{
			Code:  401,
			Msg:   "Token 无效",
			Data:  nil,
			Total: 0,
			Errors: []api.ErrorDetail{
				api.ErrorDetail{
					Msg: err.Error(),
				},
			},
		})
		return
	}

	c.JSON(http.StatusOK, api.HttpResp{
		Code:   200,
		Total:  0,
		Errors: nil,
	})
	return
}

func (h *Handler) PostValidateToken(c *gin.Context) {
	var req api.ValidateTokenReq

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.HttpResp{
			Code:  400,
			Msg:   "请求参数错误",
			Data:  nil,
			Total: 0,
			Errors: []api.ErrorDetail{
				api.ErrorDetail{
					Msg: err.Error(),
				},
			},
		})
		return
	}

	// 验证 token
	tokenResp, e := h.tokenService.GetTokenInfo(&req)
	if e != nil {
		c.JSON(http.StatusUnauthorized, api.HttpResp{
			Code:   401,
			Msg:    "Token 无效",
			Data:   nil,
			Total:  0,
			Errors: nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, api.HttpResp{
			Code:   200,
			Msg:    "Token 有效",
			Data:   tokenResp,
			Total:  0,
			Errors: nil,
		})
		return
	}
}
