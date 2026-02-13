package crypto

import (
	"encoding/base64"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
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

			decrypted, err := Decrypt(encrypted, key)
			if err != nil {
				t.Fatalf("Decrypt 失败: %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("解密结果不匹配: got %q, want %q", decrypted, tt.plaintext)
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

	// 两个密文都应能正确解密
	dec1, _ := Decrypt(enc1, key)
	dec2, _ := Decrypt(enc2, key)
	if dec1 != plaintext || dec2 != plaintext {
		t.Error("不同密文解密结果不一致")
	}
}

func TestInvalidKey(t *testing.T) {
	shortKey := []byte("too-short")

	_, err := Encrypt("test", shortKey)
	if err != ErrInvalidKey {
		t.Errorf("短密钥应返回 ErrInvalidKey, got: %v", err)
	}

	_, err = Decrypt("dGVzdA==", shortKey)
	if err != ErrInvalidKey {
		t.Errorf("短密钥解密应返回 ErrInvalidKey, got: %v", err)
	}
}

func TestDecryptTampered(t *testing.T) {
	key := []byte("scare-test-id-card-key-32-bytes!")
	encrypted, _ := Encrypt("110101199001011234", key)

	// 篡改密文
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	data[len(data)-1] ^= 0xFF
	tampered := base64.StdEncoding.EncodeToString(data)

	_, err := Decrypt(tampered, key)
	if err != ErrDecryptFailed {
		t.Errorf("篡改密文应返回 ErrDecryptFailed, got: %v", err)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := []byte("scare-test-id-card-key-32-bytes!")
	key2 := []byte("wrong-test-id-card-key-32-bytes!")

	encrypted, _ := Encrypt("110101199001011234", key1)

	_, err := Decrypt(encrypted, key2)
	if err == nil {
		t.Error("错误密钥解密应返回错误")
	}
}
