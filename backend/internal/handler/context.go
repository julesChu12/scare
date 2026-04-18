package handler

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"community-elderly-care-platform/internal/dao/model"

	"github.com/gin-gonic/gin"
)

// GetUserID 获取当前请求上下文中的用户ID
func GetUserID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

// GetUserType 获取当前请求上下文中的用户类型
func GetUserType(c *gin.Context) (string, bool) {
	value, ok := c.Get("user_type")
	if !ok {
		return "", false
	}
	userType, ok := value.(string)
	return userType, ok
}

// GetStationID 获取当前请求上下文中的站点ID
func GetStationID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("station_id")
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

// GetUserIdentities 获取用户身份列表
func GetUserIdentities(c *gin.Context) []string {
	value, ok := c.Get("user_identities")
	if !ok {
		return nil
	}
	identities, ok := value.([]string)
	if !ok {
		return nil
	}
	return identities
}

// GetUserRoles 获取用户角色列表（兼容旧代码，返回身份列表）
// Deprecated: use GetUserIdentities instead
func GetUserRoles(c *gin.Context) []string {
	return GetUserIdentities(c)
}

// imageHost 根据请求的 Host 拼接图像完整 URL
// 数据库中 images 字段存储的是相对路径（/static/...）或旧版完整 URL（http://localhost:8080/...）
// 此函数统一转换为基于当前请求 Host 的完整 URL
func imageHost(c *gin.Context, storedImages string) string {
	if storedImages == "" || storedImages == "null" {
		return storedImages
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	base := scheme + "://" + host

	var paths []string
	if err := json.Unmarshal([]byte(storedImages), &paths); err != nil {
		return storedImages
	}

	for i, p := range paths {
		// 已经是完整 URL（旧数据），提取路径部分后重拼
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			if u, err := url.Parse(p); err == nil {
				paths[i] = base + u.Path
			}
		} else if !strings.HasPrefix(p, "//") {
			// 相对路径，拼上 base
			paths[i] = base + p
		}
	}

	out, err := json.Marshal(paths)
	if err != nil {
		log.Printf("imageHost: marshal images failed: %v", err)
		return storedImages
	}
	return string(out)
}

// enrichImagesURL 将请求体中的 images 字段转换为基于当前 Host 的完整 URL
// 适用于 model.ServiceRequest 类型的请求对象
func enrichImagesURL(c *gin.Context, req *model.ServiceRequest) {
	if req == nil {
		return
	}
	req.Images = imageHost(c, req.Images)
}

