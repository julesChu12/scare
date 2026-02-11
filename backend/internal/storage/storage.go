package storage

import (
	"context"
	"io"
	"time"
)

type Provider interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error)
	SignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}
