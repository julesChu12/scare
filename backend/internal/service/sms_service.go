package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"community-elderly-care-platform/pkg/redis"
)

var (
	ErrRateLimitMinute = errors.New("发送过于频繁，请1分钟后再试")
	ErrRateLimitDaily  = errors.New("今日发送次数已达上限（10次）")
	ErrCodeInvalid     = errors.New("验证码错误或已过期")
	ErrSMSNotEnabled   = errors.New("短信服务暂未开通，请联系工作人员")
)

// SMSService 短信验证码服务
type SMSService struct {
	rdb      *redis.Client
	fallback *expiringStore
	env      string // "development" or "production"
}

// NewSMSService 创建短信服务
func NewSMSService(rdb *redis.Client, env string) *SMSService {
	svc := &SMSService{
		rdb: rdb,
		env: env,
	}
	if rdb == nil {
		svc.fallback = newExpiringStore()
	}
	return svc
}

// SetTestCode 设置测试验证码（仅用于单元测试，fallback 模式下直接存储验证码）
func (s *SMSService) SetTestCode(phone, code string) {
	if s.fallback != nil {
		s.fallback.set("sms_code:"+phone, code, 300*time.Second)
	}
}

// SendCode 发送验证码
func (s *SMSService) SendCode(phone string) error {
	// 1. 检查频率限制
	if err := s.checkRateLimit(phone); err != nil {
		return err
	}

	// 2. 生成6位随机验证码
	code, err := s.generateCode()
	if err != nil {
		return err
	}

	// 3. 存储到Redis
	if err := s.storeCode(phone, code); err != nil {
		return err
	}

	// 4. 发送短信（Mock）
	if err := s.sendSMS(phone, code); err != nil {
		return err
	}

	// 5. 更新频率限制计数
	if err := s.updateRateLimit(phone); err != nil {
		return err
	}

	return nil
}

// VerifyCode 验证验证码
func (s *SMSService) VerifyCode(phone, code string) error {
	// 测试环境：万能验证码
	if s.env == "development" || s.env == "debug" || s.env == "test" {
		if code == "000000" {
			return nil // 万能验证码直接通过
		}
	}

	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)

	if s.rdb == nil {
		storedCode, ok := s.fallback.get(key)
		if !ok || storedCode != code {
			return ErrCodeInvalid
		}
		s.fallback.del(key)
		return nil
	}

	// 从Redis获取验证码
	storedCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return ErrCodeInvalid
	}

	// 对比验证码
	if storedCode != code {
		return ErrCodeInvalid
	}

	// 验证成功，删除验证码
	_ = s.rdb.Del(ctx, key)

	return nil
}

// checkRateLimit 检查频率限制
func (s *SMSService) checkRateLimit(phone string) error {
	ctx := context.Background()

	// 检查分钟级限制
	minuteKey := fmt.Sprintf("sms_rate:%s:minute", phone)
	if s.rdb == nil {
		if s.fallback.exists(minuteKey) {
			return ErrRateLimitMinute
		}

		dailyKey := fmt.Sprintf("sms_rate:%s:daily", phone)
		if count, ok := s.fallback.get(dailyKey); ok {
			dailyCount, _ := strconv.Atoi(count)
			if dailyCount >= 10 {
				return ErrRateLimitDaily
			}
		}
		return nil
	}

	exists, err := s.rdb.Exists(ctx, minuteKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return ErrRateLimitMinute
	}

	// 检查每日限制
	dailyKey := fmt.Sprintf("sms_rate:%s:daily", phone)
	count, err := s.rdb.Get(ctx, dailyKey).Result()
	if err == nil {
		// 已有计数
		dailyCount, _ := strconv.Atoi(count)
		if dailyCount >= 10 {
			return ErrRateLimitDaily
		}
	}

	return nil
}

// generateCode 生成6位随机验证码
func (s *SMSService) generateCode() (string, error) {
	// 生成 100000 到 999999 之间的随机数
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	code := n.Int64() + 100000
	return fmt.Sprintf("%06d", code), nil
}

// storeCode 存储验证码到Redis
func (s *SMSService) storeCode(phone, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	if s.rdb == nil {
		s.fallback.set(key, code, 300*time.Second)
		return nil
	}
	return s.rdb.SetEx(ctx, key, code, 300*time.Second).Err() // 5分钟有效期
}

// sendSMS 发送短信（Mock实现）
func (s *SMSService) sendSMS(phone, code string) error {
	// 开发/调试环境：打印到控制台
	if s.env == "development" || s.env == "debug" || s.env == "test" {
		fmt.Printf("\n========== 短信验证码（Mock） ==========\n")
		fmt.Printf("手机号：%s\n", phone)
		fmt.Printf("验证码：%s\n", code)
		fmt.Printf("有效期：5分钟\n")
		fmt.Printf("======================================\n\n")
		return nil
	}

	// 生产环境：返回错误（短信服务未开通）
	return ErrSMSNotEnabled
}

// updateRateLimit 更新频率限制计数
func (s *SMSService) updateRateLimit(phone string) error {
	ctx := context.Background()

	// 设置分钟级限制
	minuteKey := fmt.Sprintf("sms_rate:%s:minute", phone)
	if s.rdb == nil {
		s.fallback.set(minuteKey, "1", 60*time.Second)

		dailyKey := fmt.Sprintf("sms_rate:%s:daily", phone)
		if s.fallback.exists(dailyKey) {
			_, err := s.fallback.incr(dailyKey)
			return err
		}
		s.fallback.set(dailyKey, "1", 24*time.Hour)
		return nil
	}

	if err := s.rdb.SetEx(ctx, minuteKey, "1", 60*time.Second).Err(); err != nil {
		return err
	}

	// 增加每日计数
	dailyKey := fmt.Sprintf("sms_rate:%s:daily", phone)
	exists, err := s.rdb.Exists(ctx, dailyKey).Result()
	if err != nil {
		return err
	}

	if exists > 0 {
		// 已存在，递增
		return s.rdb.Incr(ctx, dailyKey).Err()
	} else {
		// 不存在，创建并设置过期时间
		if err := s.rdb.Set(ctx, dailyKey, "1", 0).Err(); err != nil {
			return err
		}
		return s.rdb.Expire(ctx, dailyKey, 24*time.Hour).Err()
	}
}
