package service

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 算法对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword 验证密码与哈希值是否匹配
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
