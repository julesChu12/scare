package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

// NewPermissionHandler 创建 PermissionHandler
func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// GetPermissionTree 获取权限树
// @Summary      获取权限树
// @Description  获取系统所有权限的树形结构
// @Tags         b_role
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "获取成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/permissions/tree [get]
func (h *PermissionHandler) GetPermissionTree(c *gin.Context) {
	tree, err := h.permissionService.GetPermissionTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取权限树失败",
			"data": gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"tree": tree,
		},
	})
}
