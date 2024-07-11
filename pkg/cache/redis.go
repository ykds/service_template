package cache

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cache = (*Redis)(nil)

type Redis struct {
	*redis.Client
	ok      atomic.Bool
	err     error
	checkCh chan struct{}
}

func NewRedis(opt Option) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Password: opt.Password,
		DB:       opt.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	r := &Redis{Client: rdb, checkCh: make(chan struct{}, 1)}
	r.ok.Store(true)
	return r, nil
}

func (r *Redis) IsOk() bool {
	return r.ok.Load()
}

func (r *Redis) OccurErr(err error) {
	select {
	case r.checkCh <- struct{}{}:
		r.ok.Store(false)
		r.err = err
		go func() {
			defer func() {
				<-r.checkCh
			}()
			ticket := time.NewTicker(500 * time.Millisecond)
			for range ticket.C {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				err := r.Client.Ping(ctx).Err()
				if err != nil {
					cancel()
					continue
				}
				cancel()
				ticket.Stop()
				r.ok.Store(true)
				r.err = nil
				return
			}
		}()
	default:
	}
}

func (r *Redis) Error() error {
	return r.err
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, key ...string) (int64, error) {
	return r.Client.Del(ctx, key...).Result()
}

func (r *Redis) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.SetEx(ctx, key, value, expiration).Err()
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(ctx, key, expiration).Result()
}

func (r *Redis) MGet(ctx context.Context, keys ...string) ([]string, error) {
	ret, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	vals := make([]string, 0, len(ret))
	for _, item := range ret {
		v, ok := item.(string)
		if ok {
			vals = append(vals, v)
		}
	}
	return vals, nil
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	i, err := r.Client.Exists(ctx, []string{key}...).Result()
	return i > 0, err
}
