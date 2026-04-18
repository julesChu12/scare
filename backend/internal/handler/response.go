package handler

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Respond 返回标准 JSON 响应
func Respond(c *gin.Context, status int, msg string, data any) {
	c.PureJSON(status, Response{
		Msg:  msg,
		Data: data,
	})
}

type PageData struct {
	Items    any   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// RespondPage 返回分页 JSON 响应
func RespondPage(c *gin.Context, status int, msg string, items any, page, pageSize int, total int64) {
	c.PureJSON(status, Response{
		Msg: msg,
		Data: PageData{
			Items:    items,
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// RespondError 返回错误 JSON 响应
func RespondError(c *gin.Context, status int, msg string) {
	c.PureJSON(status, Response{
		Msg:  msg,
		Data: nil,
	})
}
