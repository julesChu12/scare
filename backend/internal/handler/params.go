package handler

import (
	"errors"
	"strconv"
)

// parseInt64Param 将字符串参数解析为 int64
func parseInt64Param(value string) (int64, error) {
	if value == "" {
		return 0, errors.New("empty")
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}
