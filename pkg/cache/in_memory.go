package cache

import (
	"context"
	"errors"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type InMemory struct {
	*gocache.Cache
}

func NewInmemory(opt Option) (Cache, error) {
	return &InMemory{
		Cache: gocache.New(5*time.Minute, 5*time.Minute),
	}, nil
}

func (im *InMemory) Get(_ context.Context, key string) (string, error) {
	value, ok := im.Cache.Get(key)
	if ok {
		return value.(string), nil
	}
	return "", errors.New("not found")
}

func (im *InMemory) Del(_ context.Context, key ...string) (int64, error) {
	for _, k := range key {
		im.Cache.Delete(k)
	}
	return int64(len(key)), nil
}

func (im *InMemory) SetEx(_ context.Context, key string, value interface{}, expiration time.Duration) error {
	im.Cache.Set(key, value, expiration)
	return nil
}

func (im *InMemory) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	value, _ := im.Get(ctx, key)
	im.Cache.Set(key, value, expiration)
	return true, nil
}

func (im *InMemory) Close() error {
	return nil
}

func (im *InMemory) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals := make([]string, 0, len(keys))
	for _, key := range keys {
		v, _ := im.Get(ctx, key)
		vals = append(vals, v)
	}
	return vals, nil
}

func (im *InMemory) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := im.Cache.Get(key)
	return ok, nil
}

func (im *InMemory) IsOk() bool {
	panic("not implemented") // TODO: Implement
}

func (im *InMemory) OccurErr(_ error) {
	panic("not implemented") // TODO: Implement
}

func (im *InMemory) Error() error {
	panic("not implemented") // TODO: Implement
}
