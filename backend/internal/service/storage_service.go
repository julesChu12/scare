package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"community-elderly-care-platform/internal/storage"
)

var ErrStorageInvalid = errors.New("invalid storage input")

type StorageService struct {
	provider storage.Provider
}

func NewStorageService(provider storage.Provider) *StorageService {
	return &StorageService{provider: provider}
}

func (s *StorageService) Upload(ctx context.Context, module string, file *multipart.FileHeader) (string, string, error) {
	if s.provider == nil || file == nil {
		return "", "", ErrStorageInvalid
	}
	if module == "" {
		module = "common"
	}
	module = sanitizeModule(module)
	key := buildObjectKey(module, file.Filename)

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	url, err := s.provider.Put(ctx, key, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return "", "", err
	}
	return url, key, nil
}

func buildObjectKey(module, filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}
	date := time.Now().Format("20060102")
	name := fmt.Sprintf("%d", time.Now().UnixNano())
	return fmt.Sprintf("%s/%s/%s%s", module, date, name, ext)
}

func sanitizeModule(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "..", "")
	value = strings.ReplaceAll(value, "\\", "/")
	value = strings.ReplaceAll(value, "//", "/")
	value = strings.Trim(value, "/")
	if value == "" {
		return "common"
	}
	return value
}
