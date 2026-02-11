package handler

import (
	"context"
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	permissionService   *service.PermissionService
	userIdentityService *service.UserIdentityService
}

func NewRoleHandler(permissionService *service.PermissionService, userIdentityService *service.UserIdentityService) *RoleHandler {
	return &RoleHandler{
		permissionService:   permissionService,
		userIdentityService: userIdentityService,
	}
}

// GetRolePermissions 获取角色的权限列表
// @Summary      获取角色权限
// @Description  获取指定角色的权限列表
// @Tags         b_role
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        role path string true "角色标识"
// @Success      200  {object} APIResponse "获取成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/roles/{role}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色参数不能为空",
			"data": nil,
		})
		return
	}

	permissions, err := h.permissionService.GetRolePermissions(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取角色权限失败",
			"data": gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"role":        role,
			"permissions": permissions,
		},
	})
}

// UpdateRolePermissions 更新角色权限
// @Summary      更新角色权限
// @Description  更新指定角色的权限列表，受影响用户需重新登录
// @Tags         b_role
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        role path string true "角色标识"
// @Param        request body object{permissions=[]string} true "权限列表"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/roles/{role}/permissions [put]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色参数不能为空",
			"data": nil,
		})
		return
	}

	// 解析请求体
	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": gin.H{"error": err.Error()},
		})
		return
	}

	// 更新权限（含 token 撤销）
	ctx := context.Background()
	affectedUsers, err := h.permissionService.UpdateRolePermissions(ctx, role, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新角色权限失败",
			"data": gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "权限更新成功，相关用户需要重新登录",
		"data": gin.H{
			"role":           role,
			"affected_users": affectedUsers,
			"tokens_revoked": affectedUsers > 0,
		},
	})
}

// UpdateUserIdentities 更新用户身份
// @Summary      更新用户身份
// @Description  授予或撤销用户身份，用户需重新登录
// @Tags         b_role
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "用户ID"
// @Param        request body object{identities=[]string} true "身份列表"
// @Success      200  {object} APIResponse "更新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无权限"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/users/{id}/identities [put]
func (h *RoleHandler) UpdateUserIdentities(c *gin.Context) {
	userID, err := parseInt64Param(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		Identities []string `json:"identities" binding:"required"`
		StationID  int64    `json:"station_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	operatorID, _ := GetUserID(c)

	if err := h.userIdentityService.SyncIdentities(userID, req.Identities, req.StationID, operatorID); err != nil {
		RespondError(c, http.StatusInternalServerError, "update identities failed: "+err.Error())
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"user_id":    userID,
		"identities": req.Identities,
	})
}
