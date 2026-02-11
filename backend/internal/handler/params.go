package handler

import (
	"errors"
	"strconv"
)

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
