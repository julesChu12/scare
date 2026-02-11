package oss

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Provider struct {
	bucket  *aliyunoss.Bucket
	baseURL string
}

type Config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	AccessKeySecret string
	BaseURL         string
}

func New(cfg Config) (*Provider, error) {
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, errors.New("oss config is incomplete")
	}
	client, err := aliyunoss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, err
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return &Provider{bucket: bucket, baseURL: baseURL}, nil
}

func (p *Provider) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error) {
	if key == "" {
		return "", errors.New("object key required")
	}
	options := []aliyunoss.Option{}
	if contentType != "" {
		options = append(options, aliyunoss.ContentType(contentType))
	}
	if err := p.bucket.PutObject(key, r, options...); err != nil {
		return "", err
	}
	if p.baseURL != "" {
		return p.baseURL + "/" + key, nil
	}
	return key, nil
}

func (p *Provider) SignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if key == "" {
		return "", errors.New("object key required")
	}
	seconds := int64(expire.Seconds())
	if seconds <= 0 {
		seconds = 3600
	}
	return p.bucket.SignURL(key, aliyunoss.HTTPGet, seconds)
}

func (p *Provider) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("object key required")
	}
	return p.bucket.DeleteObject(key)
}
