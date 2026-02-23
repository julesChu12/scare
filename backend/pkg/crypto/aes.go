package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKey        = errors.New("crypto: 密钥长度必须为 32 字节（AES-256）")
	ErrCiphertextTooShort = errors.New("crypto: 密文长度不足")
	ErrDecryptFailed     = errors.New("crypto: 解密失败")
)

// Encrypt 使用 AES-256-GCM 加密明文
// 返回 base64(nonce + ciphertext + tag)
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// nonce + ciphertext(含 tag)
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

