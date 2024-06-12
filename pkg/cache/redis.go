package cache

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "redis: nil")
}

type Option struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

type Redis struct {
	*redis.Client
	ok      atomic.Bool
	err     error
	checkCh chan struct{}
}

func NewRedis(opt Option) (*Redis, error) {
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
