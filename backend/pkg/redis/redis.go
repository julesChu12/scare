package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type Client struct {
	*goredis.Client
}

func InitRedis(cfg Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, errors.New("redis host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 6379
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
