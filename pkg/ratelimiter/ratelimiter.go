package ratelimiter

import (
	"context"
	"service_template/pkg/cache"
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter interface {
	CanPass(key string) (bool, error)
}

func NewRedisRateLimiter(rdb *cache.Redis, max int64, interval int) RateLimiter {
	return &redisRateLimiter{rdb, max, interval}
}

type redisRateLimiter struct {
	rdb      *cache.Redis
	max      int64
	interval int
}

func (r *redisRateLimiter) CanPass(key string) (bool, error) {
	if !r.rdb.IsOk() {
		return false, r.rdb.Error()
	}
	script := `
local res = redis.call('INCR', KEYS[1])
if res == 1
then 
	redis.call('EXPIRE', KEYS[1], ARGV[1])
end 
return res
`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	count, err := r.rdb.Eval(ctx, script, []string{key}, r.interval).Int64()
	if err != nil {
		r.rdb.OccurErr(err)
		return false, err
	}
	return count <= r.max, nil
}

func NewLocalRateLimiter(max int64, interval int) RateLimiter {
	return &localRateLimiter{
		m:        sync.Mutex{},
		max:      max,
		interval: interval,
		entries:  make(map[string]*Entry, 2048),
	}
}

type Entry struct {
	count     atomic.Int64
	ExpiredAt time.Time
}

type localRateLimiter struct {
	m        sync.Mutex
	max      int64
	interval int
	entries  map[string]*Entry
}

func (l *localRateLimiter) CanPass(key string) (bool, error) {
	l.m.Lock()
	if e, ok := l.entries[key]; ok {
		l.m.Unlock()
		if e.ExpiredAt.Before(time.Now()) {
			e.count.Store(0)
			e.ExpiredAt = time.Now().Add(time.Second * time.Duration(l.interval))
		}
		e.count.Add(1)
		return e.count.Load() <= l.max, nil
	}
	e := &Entry{ExpiredAt: time.Now().Add(time.Second * time.Duration(l.interval))}
	e.count.Store(1)
	l.entries[key] = e
	l.m.Unlock()
	return true, nil
}
