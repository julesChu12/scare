package storage

import (
	"errors"
	"strings"

	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/storage/local"
	"community-elderly-care-platform/internal/storage/oss"
)

func NewProvider(cfg config.StorageConfig) (Provider, error) {
	driver := strings.ToLower(cfg.Driver)
	switch driver {
	case "local":
		return local.New(cfg.Local.BasePath, cfg.Local.BaseURL)
	case "oss":
		return oss.New(oss.Config{
			Endpoint:        cfg.OSS.Endpoint,
			Bucket:          cfg.OSS.Bucket,
			AccessKeyID:     cfg.OSS.AccessKeyID,
			AccessKeySecret: cfg.OSS.AccessKeySecret,
			BaseURL:         cfg.OSS.BaseURL,
		})
	default:
		return nil, errors.New("unsupported storage driver")
	}
}
