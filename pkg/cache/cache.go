package cache

import (
	"context"
	"strings"
	"time"
)

const (
	CacheTypeRedis    = "redis"
	CacheTypeInmemory = "inmemory"
)

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "redis: nil") || strings.Contains(err.Error(), "not found")
}

type Option struct {
	Type     string `json:"type" yaml:"type"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

type Cache interface {
	IsOk() bool
	OccurErr(error)
	Error() error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, key ...string) ([]string, error)
	Del(ctx context.Context, key ...string) (int64, error)
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Close() error
	Exists(ctx context.Context, key string) (bool, error)
}

func NewCache(opt Option) (Cache, error) {
	switch opt.Type {
	case CacheTypeRedis:
		return NewRedis(opt)
	case CacheTypeInmemory:
		return NewInmemory(opt)
	default:
		panic("not support cache")
	}
}
