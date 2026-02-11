package middleware

import (
	"net/http"
	"strings"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限检查中间件（仅对B端用户验证权限）
func PermissionMiddleware(permService *service.PermissionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// C端用户跳过权限检查
		userType, _ := c.Get("user_type")
		if userType != "b_end" {
			c.Next()
			return
		}

		// 从上下文获取用户身份列表（由 AuthMiddleware 设置）
		identitiesValue, exists := c.Get("user_identities")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg":  "unauthorized",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 类型断言为 []string
		userIdentities, ok := identitiesValue.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "invalid identities format",
				"data": nil,
			})
			c.Abort()
			return
		}

		// Admin 身份跳过权限检查
		for _, identity := range userIdentities {
			if identity == "admin" {
				c.Next()
				return
			}
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		// 公共 API 跳过检查
		if permService.IsPublicAPI(path, method) {
			c.Next()
			return
		}

		// 检查权限（使用身份作为角色来匹配权限）
		if ok, _ := permService.CheckAPIPermission(userIdentities, path, method); !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "forbidden",
				"data": gin.H{
					"identities": userIdentities,
					"path":       path,
					"method":     method,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MatchPath 路径匹配，支持通配符 *
// 例如: /api/v1/b/tasks/*/claim 匹配 /api/v1/b/tasks/123/claim
func MatchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if part == "*" {
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}
