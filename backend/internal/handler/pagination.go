package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// GetPagination 解析并返回分页参数
func GetPagination(c *gin.Context) (int, int) {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), defaultPageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
