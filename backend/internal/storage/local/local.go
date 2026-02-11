package local

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Provider struct {
	basePath string
	baseURL  string
}

func New(basePath, baseURL string) (*Provider, error) {
	if basePath == "" || baseURL == "" {
		return nil, errors.New("base path and base url are required")
	}
	return &Provider{basePath: basePath, baseURL: strings.TrimRight(baseURL, "/")}, nil
}

func (p *Provider) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error) {
	if key == "" {
		return "", errors.New("object key required")
	}
	clean := filepath.Clean(key)
	path := filepath.Join(p.basePath, clean)
	if !strings.HasPrefix(path, filepath.Clean(p.basePath)+string(os.PathSeparator)) {
		return "", errors.New("invalid object key")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, r); err != nil {
		return "", err
	}
	return p.baseURL + "/" + clean, nil
}

func (p *Provider) SignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if key == "" {
		return "", errors.New("object key required")
	}
	clean := filepath.Clean(key)
	return p.baseURL + "/" + clean, nil
}

func (p *Provider) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("object key required")
	}
	clean := filepath.Clean(key)
	path := filepath.Join(p.basePath, clean)
	return os.Remove(path)
}
