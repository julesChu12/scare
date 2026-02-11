package jwt

import "testing"

func TestManagerGenerateAndParse(t *testing.T) {
	manager := NewManager("secret", 1, 2)

	// 测试 B端 Token（包含 identities）
	token, err := manager.GenerateToken(42, "b_end", 7, []string{"admin", "staff"}, "admin")
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
	if len(claims.Identities) != 2 || claims.Identities[0] != "admin" || claims.Identities[1] != "staff" {
		t.Fatalf("expected identities [admin staff], got %v", claims.Identities)
	}
	if claims.Primary != "admin" {
		t.Fatalf("expected primary admin, got %s", claims.Primary)
	}
	if claims.Type != "b_end" {
		t.Fatalf("expected type b_end, got %s", claims.Type)
	}
	if claims.StationID != 7 {
		t.Fatalf("expected station id 7, got %d", claims.StationID)
	}

	// 测试 Roles() 兼容方法
	if len(claims.Roles()) != 2 {
		t.Fatalf("expected Roles() to return identities, got %v", claims.Roles())
	}

	if _, err := manager.ParseToken("invalid.token.value"); err == nil {
		t.Fatalf("expected invalid token error")
	}
}

func TestManagerCEndToken(t *testing.T) {
	manager := NewManager("secret", 1, 2)

	// 测试 C端 Token（可以有身份类型，如 elderly）
	token, err := manager.GenerateToken(100, "c_end", 0, []string{"elderly"}, "elderly")
	if err != nil {
		t.Fatalf("generate c_end token failed: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse c_end token failed: %v", err)
	}
	if claims.UserID != 100 {
		t.Fatalf("expected user id 100, got %d", claims.UserID)
	}
	if len(claims.Identities) != 1 || claims.Identities[0] != "elderly" {
		t.Fatalf("expected identities [elderly], got %v", claims.Identities)
	}
	if claims.Primary != "elderly" {
		t.Fatalf("expected primary elderly, got %s", claims.Primary)
	}
	if claims.Type != "c_end" {
		t.Fatalf("expected type c_end, got %s", claims.Type)
	}
}

func TestManagerRefreshToken(t *testing.T) {
	manager := NewManager("secret", 1, 2)

	// 测试刷新 Token
	token, err := manager.GenerateRefreshToken(42, "b_end", 7, []string{"admin"}, "admin")
	if err != nil {
		t.Fatalf("generate refresh token failed: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse refresh token failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
}
