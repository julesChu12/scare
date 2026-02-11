package middleware

import (
	"context"
	"net/http"
	"strings"

	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *jwt.Manager, blacklistService *service.TokenBlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":  "missing authorization",
				"data": nil,
			})
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":  "invalid authorization",
				"data": nil,
			})
			return
		}

		claims, err := jwtManager.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":  "invalid token",
				"data": nil,
			})
			return
		}

		ctx := context.Background()

		// 检查 token 是否在黑名单中（单个 token 黑名单）
		if blacklistService != nil && claims.ID != "" {
			isBlacklisted, err := blacklistService.IsBlacklisted(ctx, claims.ID)
			if err == nil && isBlacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"msg":  "token revoked",
					"data": nil,
				})
				return
			}
		}

		// 检查用户的所有 token 是否被撤销（角色变更时）
		if blacklistService != nil && claims.IssuedAt != nil {
			isRevoked, err := blacklistService.IsUserTokenRevoked(ctx, claims.UserID, claims.Type, claims.IssuedAt.Unix())
			if err == nil && isRevoked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"msg":  "user tokens revoked due to role change",
					"data": nil,
				})
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_identities", claims.Identities) // 用户的所有身份（数组）
		c.Set("user_primary", claims.Primary)        // 主身份
		c.Set("user_type", claims.Type)              // 端类型(c_end/b_end)
		c.Set("station_id", claims.StationID)
		c.Set("token_id", claims.ID)                 // Token ID for logout
		c.Set("token_expires_at", claims.ExpiresAt.Time)

		c.Next()
	}
}
