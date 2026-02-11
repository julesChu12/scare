package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// BAuthHandler B端认证处理器
type BAuthHandler struct {
	authService       *service.AuthService
	userRepo          *repository.UserRepository
	userIdentityRepo  *repository.UserIdentityRepository
	permissionService *service.PermissionService
}

func NewBAuthHandler(authService *service.AuthService, userRepo *repository.UserRepository, userIdentityRepo *repository.UserIdentityRepository, permissionService *service.PermissionService) *BAuthHandler {
	return &BAuthHandler{
		authService:       authService,
		userRepo:          userRepo,
		userIdentityRepo:  userIdentityRepo,
		permissionService: permissionService,
	}
}

type bLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type bRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login B端登录接口
// @Summary      B端用户登录
// @Description  B端工作人员使用手机号和密码登录
// @Tags         b_auth
// @Accept       json
// @Produce      json
// @Param        request body bLoginRequest true "登录请求"
// @Success      200  {object} APIResponse "登录成功，返回token和用户信息"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "用户名或密码错误"
// @Failure      403  {object} APIResponse "用户已禁用或无B端角色"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/auth/login [post]
func (h *BAuthHandler) Login(c *gin.Context) {
	var req bLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// B端登录
	tokens, user, err := h.authService.LoginBEnd(req.Phone, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			RespondError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if err == service.ErrUserInactive {
			RespondError(c, http.StatusForbidden, "user inactive")
			return
		}
		if err == service.ErrNoRoleForBEnd {
			RespondError(c, http.StatusForbidden, "user has no B-end identity")
			return
		}
		RespondError(c, http.StatusInternalServerError, "login failed")
		return
	}

	// 获取用户的所有B端身份
	identities, err := h.userIdentityRepo.GetBEndIdentities(user.ID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get identities")
		return
	}
	var userIdentities []string
	for _, identity := range identities {
		userIdentities = append(userIdentities, identity.IdentityType)
	}
	if len(userIdentities) == 0 {
		RespondError(c, http.StatusForbidden, "user has no B-end identity")
		return
	}

	// 获取主身份
	var primary string
	primaryIdentity, err := h.userIdentityRepo.GetPrimaryIdentity(user.ID)
	if err == nil {
		primary = primaryIdentity.IdentityType
	} else {
		primary = userIdentities[0]
	}

	data := gin.H{
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"user_id":       user.ID,
		"identities":    userIdentities,   // B端所有身份
		"primary":       primary,          // 主身份
		"type":          "b_end",          // 端类型
		"station_id":    user.StationID,
		"name":          user.Name,
		"phone":         user.Phone,
		"status":        user.Status,
	}
	Respond(c, http.StatusOK, "ok", data)
}

// Refresh B端刷新Token
// @Summary      刷新访问令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         b_auth
// @Accept       json
// @Produce      json
// @Param        request body bRefreshRequest true "刷新令牌请求"
// @Success      200  {object} APIResponse "刷新成功"
// @Failure      400  {object} APIResponse "请求参数错误"
// @Failure      401  {object} APIResponse "令牌无效或已过期"
// @Router       /b/auth/refresh [post]
func (h *BAuthHandler) Refresh(c *gin.Context) {
	var req bRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	tokens, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "invalid token")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"token":         tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// Me 获取当前登录用户信息（B端）含权限并集
// @Summary      获取当前B端用户信息
// @Description  获取当前登录的B端用户详细信息及权限列表
// @Tags         b_auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "获取成功，返回用户信息和权限"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      403  {object} APIResponse "无B端角色"
// @Failure      404  {object} APIResponse "用户不存在"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/auth/me [get]
func (h *BAuthHandler) Me(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		RespondError(c, http.StatusUnauthorized, "missing user")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	// 获取用户的所有B端身份
	identities, err := h.userIdentityRepo.GetBEndIdentities(user.ID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get identities")
		return
	}
	var userIdentities []string
	for _, identity := range identities {
		userIdentities = append(userIdentities, identity.IdentityType)
	}
	if len(userIdentities) == 0 {
		RespondError(c, http.StatusForbidden, "user has no B-end identity")
		return
	}

	// 获取主身份
	var primary string
	primaryIdentity, err := h.userIdentityRepo.GetPrimaryIdentity(user.ID)
	if err == nil {
		primary = primaryIdentity.IdentityType
	} else {
		primary = userIdentities[0]
	}

	// 计算权限并集（包含 authenticated 角色）
	permissions, err := h.permissionService.GetUserPermissions(userIdentities)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to get permissions")
		return
	}

	Respond(c, http.StatusOK, "ok", gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"phone":      user.Phone,
			"identities": userIdentities, // 所有身份
			"primary":    primary,        // 主身份
			"station_id": user.StationID,
			"status":     user.Status,
		},
		"permissions": permissions, // 权限并集
	})
}
