package lock

import (
	"context"
	"service_template/pkg/cache"
	"sync"
	"time"
)

type Lock interface {
	TryLock(key string) (bool, *LockEntry, error)
	UnLock(*LockEntry) error
}

type LockEntry struct {
	timer     *time.Timer
	done      bool
	key       string
	expiredAt time.Time
}

func NewLocalLock() Lock {
	return &localLock{
		entries: sync.Map{},
	}
}

type localLock struct {
	entries sync.Map
}

func (l *localLock) TryLock(key string) (bool, *LockEntry, error) {
	_, loaded := l.entries.LoadOrStore(key, struct{}{})
	return !loaded, &LockEntry{key: key}, nil
}

func (l *localLock) UnLock(entry *LockEntry) error {
	if entry == nil {
		return nil
	}
	l.entries.Delete(entry.key)
	return nil
}

func NewRedisLock(rdb *cache.Redis, defaultTtl int) Lock {
	if defaultTtl <= 0 {
		panic("lock ttl must be positive")
	}
	size := uint(1024 * 1024)
	l := &redisLock{
		Redis:     rdb,
		ttl:       defaultTtl,
		m:         sync.Mutex{},
		size:      size,
		renewCh:   make(chan *LockEntry, 4096),
		renewList: make([]*LockEntry, size),
		sig:       make(chan struct{}, 1),
	}
	l.renewLock()
	return l
}

type redisLock struct {
	*cache.Redis
	ttl               int
	m                 sync.Mutex
	renewCh           chan *LockEntry
	renewList         []*LockEntry
	overflow          []*LockEntry
	sig               chan struct{}
	w, r, size, count uint
}

func (r *redisLock) add(entry *LockEntry) {
	r.m.Lock()
	if r.count >= r.size {
		r.overflow = append(r.overflow, r.renewList...)
		r.r = 0
		r.w = 0
		r.count = 0
	}
	r.renewList[r.w] = entry
	r.w++
	if r.w == r.size {
		r.w = 0
	}
	r.count++
	r.m.Unlock()
}

func (r *redisLock) pull() []*LockEntry {
	r.m.Lock()
	n := make([]*LockEntry, len(r.overflow)+int(r.count))
	if len(r.overflow) > 0 {
		copy(n, r.overflow)
		r.overflow = nil
	}
	if r.count > 0 {
		copy(n[len(r.overflow):], r.renewList[r.r:r.w])
		r.r = 0
		r.w = 0
		r.count = 0
	}
	r.m.Unlock()
	return n
}

// 续约
func (r *redisLock) renewLock() {
	go func() {
		for entry := range r.renewCh {
			r.add(entry)
			select {
			case r.sig <- struct{}{}:
			default:
			}
		}
	}()
	go func() {
		for {
			<-r.sig
			toRenew := r.pull()
			for _, e := range toRenew {
				if e.done {
					continue
				}
				e.timer = time.NewTimer(time.Until(e.expiredAt))
				<-e.timer.C
				if e.done {
					continue
				}
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				ok, err := r.Expire(ctx, e.key, time.Duration(r.ttl)*time.Second)
				if err != nil {
					r.OccurErr(err)
				}
				if ok {
					e.expiredAt = time.Now().Add(time.Second * time.Duration(r.ttl/2))
					r.renewCh <- e
				}
				cancel()
				e.timer.Stop()
			}
		}
	}()
}

func (r *redisLock) TryLock(key string) (bool, *LockEntry, error) {
	if !r.IsOk() {
		return false, nil, r.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ok, err := r.SetNX(ctx, key, "", time.Duration(r.ttl)*time.Second).Result()
	if err != nil {
		r.OccurErr(err)
		return false, nil, err
	}
	if !ok {
		return false, nil, nil
	}
	entry := &LockEntry{
		key:       key,
		expiredAt: time.Now().Add(time.Second * time.Duration(r.ttl/2)),
	}
	r.renewCh <- entry
	return true, entry, nil
}

func (r *redisLock) UnLock(entry *LockEntry) error {
	if entry == nil {
		return nil
	}
	entry.done = true
	if entry.timer != nil {
		entry.timer.Reset(0)
	}
	if !r.IsOk() {
		return r.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := r.Del(ctx, entry.key)
	if err != nil {
		r.OccurErr(err)
	}
	return err
}

func NewCombinationLock(rdb *cache.Redis, defaultTtl int) Lock {
	return &combinationLock{
		rdb:       rdb,
		redisLock: NewRedisLock(rdb, defaultTtl),
		localLock: NewLocalLock(),
	}
}

type combinationLock struct {
	rdb       *cache.Redis
	redisLock Lock
	localLock Lock
}

func (c *combinationLock) TryLock(key string) (bool, *LockEntry, error) {
	ok, entry, err := c.redisLock.TryLock(key)
	if err != nil {
		return c.localLock.TryLock(key)
	}
	return ok, entry, nil
}

func (c *combinationLock) UnLock(entry *LockEntry) error {
	if entry == nil {
		return nil
	}
	if c.rdb.IsOk() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		i, err := c.rdb.Del(ctx, entry.key)
		if i != 0 {
			return nil
		}
		if err != nil {
			c.rdb.OccurErr(err)
		}
	}
	_ = c.localLock.UnLock(entry)
	return nil
}
