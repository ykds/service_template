package ratelimiter

import (
	"accelerator_api/pkg/cache"
	"context"
	"sync"
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
	return count < r.max, nil
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
	count     int64
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
		if e.ExpiredAt.Before(time.Now()) {
			e.count = 0
			e.ExpiredAt = time.Now().Add(time.Second * time.Duration(l.interval))
		}
		e.count += 1
		l.m.Unlock()
		return e.count < l.max, nil
	}
	l.entries[key] = &Entry{ExpiredAt: time.Now().Add(time.Second * time.Duration(l.interval))}
	l.m.Unlock()
	return true, nil
}
