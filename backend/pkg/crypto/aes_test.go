package crypto

import (
	"encoding/base64"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := []byte("scare-test-id-card-key-32-bytes!")

	tests := []struct {
		name      string
		plaintext string
	}{
		{"18位身份证号", "110101199001011234"},
		{"15位身份证号", "110101900101123"},
		{"空字符串", ""},
		{"含字母X", "11010119900101123X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.plaintext, key)
			if err != nil {
				t.Fatalf("Encrypt 失败: %v", err)
			}

			// 验证是合法 base64
			if _, err := base64.StdEncoding.DecodeString(encrypted); err != nil {
				t.Fatalf("加密结果不是合法 base64: %v", err)
			}
		})
	}
}

func TestEncryptNonDeterministic(t *testing.T) {
	key := []byte("scare-test-id-card-key-32-bytes!")
	plaintext := "110101199001011234"

	enc1, _ := Encrypt(plaintext, key)
	enc2, _ := Encrypt(plaintext, key)

	if enc1 == enc2 {
		t.Error("两次加密结果相同，nonce 未随机化")
	}
}

func TestInvalidKey(t *testing.T) {
	shortKey := []byte("too-short")

	_, err := Encrypt("test", shortKey)
	if err != ErrInvalidKey {
		t.Errorf("短密钥应返回 ErrInvalidKey, got: %v", err)
	}
}
