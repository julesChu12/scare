package handler

import (
	"context"
	"net/http"
	"time"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// LogoutHandler 登出处理器
type LogoutHandler struct {
	blacklistService *service.TokenBlacklistService
}

func NewLogoutHandler(blacklistService *service.TokenBlacklistService) *LogoutHandler {
	return &LogoutHandler{
		blacklistService: blacklistService,
	}
}

// Logout 登出接口（将当前 token 加入黑名单）
// @Summary      用户登出
// @Description  将当前访问令牌加入黑名单，使其失效
// @Tags         b_auth,c_auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object} APIResponse "登出成功"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      500  {object} APIResponse "服务器错误"
// @Router       /b/auth/logout [post]
// @Router       /c/auth/logout [post]
func (h *LogoutHandler) Logout(c *gin.Context) {
	tokenID, exists := c.Get("token_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "missing token id")
		return
	}

	expiresAt, exists := c.Get("token_expires_at")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "missing token expiration")
		return
	}

	// 计算 token 剩余有效期
	expiresTime := expiresAt.(time.Time)
	ttl := time.Until(expiresTime)
	if ttl <= 0 {
		// Token 已过期，无需加入黑名单
		Respond(c, http.StatusOK, "logout successful", nil)
		return
	}

	ctx := context.Background()
	if err := h.blacklistService.AddToBlacklist(ctx, tokenID.(string), ttl); err != nil {
		RespondError(c, http.StatusInternalServerError, "logout failed")
		return
	}

	Respond(c, http.StatusOK, "logout successful", nil)
}
