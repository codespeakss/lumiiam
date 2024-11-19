package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) postUser(c *gin.Context) {
	//var dto model.User
	//if err := c.BindJSON(&dto); err != nil {
	//	logrus.Error("bad request in Post(): ", err)
	//	c.JSON(http.StatusOK, api.HttpRespBody{
	//		Code: api.HTTP_ERR_FORMAT,
	//		Msg:  "ERR_FORMAT: " + err.Error(),
	//	})
	//	return
	//}
	//
	//user, e := h.userService.CreateUser(dto)
	//if e != nil {
	//	c.JSON(http.StatusOK, api.HttpRespBody{
	//		Code: api.HTTP_ERR_OTHER,
	//		Msg:  "ERR_OTHER: " + e.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, api.HttpRespBody{
	//	Code: 0,
	//	Msg:  "ok",
	//	Data: user,
	//})
	//return
}

func (h *Handler) getUserByID(c *gin.Context) {
	//idStr := c.Param("id")
	//
	//user, e := h.userService.GetUserByID(idStr)
	//if e != nil {
	//	c.JSON(http.StatusOK, api.HttpRespBody{
	//		Code: api.HTTP_ERR_OTHER,
	//		Msg:  "ERR_OTHER: " + e.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, api.HttpRespBody{
	//	Code: 0,
	//	Msg:  "ok",
	//	Data: user,
	//})
	//return
}

func (h *Handler) putUserByID(c *gin.Context) {
	//idStr := c.Param("id")
	//
	//user, e := h.userService.GetUserByID(idStr)
	//if e != nil {
	//	c.JSON(http.StatusOK, api.HttpRespBody{
	//		Code: api.HTTP_ERR_OTHER,
	//		Msg:  "ERR_OTHER: " + e.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, api.HttpRespBody{
	//	Code: 0,
	//	Msg:  "ok",
	//	Data: user,
	//})
	//return
}
