package storage

import (
	"errors"
	"net/url"
	"strings"

	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/internal/storage/local"
	"community-elderly-care-platform/internal/storage/oss"
)

func NewProvider(cfg config.StorageConfig) (Provider, error) {
	driver := strings.ToLower(cfg.Driver)
	switch driver {
	case "local":
		urlPrefix := extractPathPrefix(cfg.Local.BaseURL)
		return local.New(cfg.Local.BasePath, urlPrefix)
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

// extractPathPrefix 从完整 URL 中提取路径部分
// 例如: http://localhost:8080/static -> /static
func extractPathPrefix(fullURL string) string {
	if fullURL == "" {
		return "/static"
	}
	u, err := url.Parse(fullURL)
	if err != nil {
		return "/static"
	}
	return u.Path
}
