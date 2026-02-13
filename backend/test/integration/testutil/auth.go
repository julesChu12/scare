//go:build integration

package testutil

import (
	"community-elderly-care-platform/pkg/jwt"
)

const jwtSecret = "test-jwt-secret-key-for-integration-tests"

func newJWTManager() *jwt.Manager {
	return jwt.NewManager(jwtSecret, 24, 168)
}

// AdminToken 生成 Admin 角色的 B 端 JWT
func AdminToken() string {
	token, _ := newJWTManager().GenerateToken(1, "b_end", 0, []string{"admin"}, "admin")
	return token
}

// StationManagerToken 生成站点管理员 B 端 JWT
func StationManagerToken() string {
	token, _ := newJWTManager().GenerateToken(2, "b_end", 1, []string{"station_manager"}, "station_manager")
	return token
}

// StaffToken 生成工作人员 B 端 JWT
func StaffToken() string {
	token, _ := newJWTManager().GenerateToken(4, "b_end", 1, []string{"staff"}, "staff")
	return token
}

// CEndElderlyToken 生成 C 端老年人 JWT
func CEndElderlyToken() string {
	token, _ := newJWTManager().GenerateToken(10, "c_end", 0, []string{"elderly"}, "elderly")
	return token
}

// CEndFamilyToken 生成 C 端家属 JWT
func CEndFamilyToken() string {
	token, _ := newJWTManager().GenerateToken(11, "c_end", 0, []string{"family"}, "family")
	return token
}
