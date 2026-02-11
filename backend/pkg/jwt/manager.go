package jwt

import (
	"errors"
	"time"

	"github.com/google/uuid"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID     int64    `json:"uid"`
	Identities []string `json:"identities,omitempty"` // 用户的所有身份类型
	Primary    string   `json:"primary,omitempty"`    // 主身份类型
	Type       string   `json:"type"`                 // 端类型(c_end/b_end)
	StationID  int64    `json:"station_id"`
	jwtlib.RegisteredClaims
}

// Roles 兼容旧代码，返回 Identities
func (c *Claims) Roles() []string {
	return c.Identities
}

type Manager struct {
	secret           []byte
	expiresIn        time.Duration
	refreshExpiresIn time.Duration
}

func NewManager(secret string, expiresInHours, refreshExpiresInHours int) *Manager {
	return &Manager{
		secret:           []byte(secret),
		expiresIn:        time.Duration(expiresInHours) * time.Hour,
		refreshExpiresIn: time.Duration(refreshExpiresInHours) * time.Hour,
	}
}

// GenerateToken 生成访问令牌
// identities: 用户身份类型数组
// primary: 主身份类型
func (m *Manager) GenerateToken(userID int64, endType string, stationID int64, identities []string, primary string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:     userID,
		Identities: identities,
		Primary:    primary,
		Type:       endType,
		StationID:  stationID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.expiresIn)),
			IssuedAt:  jwtlib.NewNumericDate(now),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateRefreshToken 生成刷新令牌
func (m *Manager) GenerateRefreshToken(userID int64, endType string, stationID int64, identities []string, primary string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:     userID,
		Identities: identities,
		Primary:    primary,
		Type:       endType,
		StationID:  stationID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.refreshExpiresIn)),
			IssuedAt:  jwtlib.NewNumericDate(now),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	parsed, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
