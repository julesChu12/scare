package service

import (
	"context"
	"fmt"
	"time"

	"community-elderly-care-platform/pkg/redis"
)

// TokenBlacklistService JWT 黑名单服务
type TokenBlacklistService struct {
	redis *redis.Client
}

func NewTokenBlacklistService(redis *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{
		redis: redis,
	}
}

// AddToBlacklist 将 token 加入黑名单
// tokenID: token 的唯一标识（jti）
// expiresIn: token 的剩余有效期
func (s *TokenBlacklistService) AddToBlacklist(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("token:blacklist:%s", tokenID)
	return s.redis.Set(ctx, key, "1", expiresIn).Err()
}

// IsBlacklisted 检查 token 是否在黑名单中
func (s *TokenBlacklistService) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", tokenID)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// RevokeUserTokens 撤销用户的所有 token（角色变更时使用）
// userID: 用户ID
// endType: 端类型（c_end/b_end）
func (s *TokenBlacklistService) RevokeUserTokens(ctx context.Context, userID int64, endType string) error {
	// 使用用户维度的黑名单标记
	// 格式：user:revoked:{userID}:{endType}
	key := fmt.Sprintf("user:revoked:%d:%s", userID, endType)
	// 设置24小时过期（超过 refresh token 的有效期）
	return s.redis.Set(ctx, key, time.Now().Unix(), 24*time.Hour).Err()
}

// IsUserTokenRevoked 检查用户的 token 是否已被撤销
// userID: 用户ID
// endType: 端类型
// issuedAt: token 的签发时间
func (s *TokenBlacklistService) IsUserTokenRevoked(ctx context.Context, userID int64, endType string, issuedAt int64) (bool, error) {
	key := fmt.Sprintf("user:revoked:%d:%s", userID, endType)
	revokedAt, err := s.redis.Get(ctx, key).Int64()
	if err != nil {
		// Key 不存在说明没有被撤销
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	// 如果 token 签发时间早于撤销时间，则认为已被撤销
	return issuedAt < revokedAt, nil
}
