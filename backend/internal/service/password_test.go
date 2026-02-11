package service

import "testing"

func TestPasswordHashVerify(t *testing.T) {
	password := "Test@123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if err := VerifyPassword(hash, password); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if err := VerifyPassword(hash, "wrong"); err == nil {
		t.Fatalf("expected verify to fail")
	}
}
