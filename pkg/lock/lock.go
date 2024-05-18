package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Lock interface {
	TryLock(key string) (bool, error)
	UnLock(key string) error
}

func NewLocalLock() Lock {
	return &localLock{
		entries: sync.Map{},
	}
}

type localLock struct {
	entries sync.Map
}

func (l *localLock) TryLock(key string) (bool, error) {
	_, loaded := l.entries.LoadOrStore(key, struct{}{})
	return !loaded, nil
}

func (l *localLock) UnLock(key string) error {
	l.entries.Delete(key)
	return nil
}

func NewRedisLock(rdb *redis.Client, defaultTtl int) Lock {
	if defaultTtl <= 0 {
		panic("lock ttl must be positive")
	}
	size := uint(1024 * 1024)
	l := &redisLock{
		Client:    rdb,
		ttl:       defaultTtl,
		m:         sync.Mutex{},
		size:      size,
		renewCh:   make(chan *renewEntry, 4096),
		renewList: make([]*renewEntry, size),
		sig:       make(chan *renewEntry, 1),
	}
	l.renewLock()
	return l
}

type renewEntry struct {
	key       string
	expiredAt time.Time
}

type redisLock struct {
	*redis.Client
	ttl               int
	m                 sync.Mutex
	renewCh           chan *renewEntry
	renewList         []*renewEntry // 循环队列, 重复利用内存, 避免不断申请和释放内存
	overflow          []*renewEntry // 溢出队列, 当循环队列满了, 把所有数据迁移到这里, 然后循环队列重置. 长度无限制
	sig               chan *renewEntry
	w, r, size, count uint
}

func (r *redisLock) add(entry *renewEntry) {
	r.m.Lock()
	// 循环队列满了, 数据全转移到移除队列, 重置循环队列
	if r.count >= r.size {
		r.overflow = append(r.overflow, r.renewList...)
		r.r = 0
		r.w = 0
		r.count = 0
	}
	// 添加新元素
	r.renewList[r.w] = entry
	r.w++
	if r.w == r.size {
		r.w = 0
	}
	r.count++
	r.m.Unlock()
}

// 拉取所有待续约 key
func (r *redisLock) pull() []*renewEntry {
	r.m.Lock()
	// 因为溢出队列是在循环队列满了后, 才有数据, 也就是此时的溢出队列数据是比循环队列的数据更早
	// 所以如果溢出队列有数据, 把溢出队列放在前部, 再放循环队列
	n := make([]*renewEntry, len(r.overflow)+int(r.count))
	if len(r.overflow) > 0 {
		copy(n, r.overflow)
		r.overflow = nil
	}
	if len(r.renewList) > 0 {
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
	// renewCh 接收到数据时, 先通过 sig chan 推送一条过去 hang 住等待续约
	// 其余的放到循环队列中等待.
	go func() {
		for entry := range r.renewCh {
			select {
			case r.sig <- entry:
			default:
				r.add(entry)
			}
		}
	}()
	go func() {
		// 从 sig 中取到最新的续约项, 等待到期续约.
		// 完成后再把目前列表所有的待续约项拿出来循环等待续约, 减少对队列的锁竞争.
		for {
			e := <-r.sig
			t := time.NewTimer(time.Until(e.expiredAt))
			<-t.C
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			ok, err := r.Expire(ctx, e.key, time.Duration(r.ttl)*time.Second).Result()
			if err == nil && ok {
				e.expiredAt = time.Now().Add(time.Second * time.Duration(r.ttl/2))
				r.renewCh <- e
			}
			cancel()
			t.Stop()

			toRenew := r.pull()
			for _, e := range toRenew {
				t := time.NewTimer(time.Until(e.expiredAt))
				<-t.C
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				ok, err := r.Expire(ctx, e.key, time.Duration(r.ttl)*time.Second).Result()
				if err == nil && ok {
					e.expiredAt = time.Now().Add(time.Second * time.Duration(r.ttl/2))
					r.renewCh <- e
				}
				cancel()
				t.Stop()
			}
		}
	}()
}

func (r *redisLock) TryLock(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ok, err := r.SetNX(ctx, key, "", time.Duration(r.ttl)*time.Second).Result()
	if err != nil || !ok {
		return false, err
	}
	if ok {
		r.renewCh <- &renewEntry{
			key:       key,
			expiredAt: time.Now().Add(time.Second * time.Duration(r.ttl/2)),
		}
	}
	return true, nil
}

func (r *redisLock) UnLock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return r.Del(ctx, key).Err()
}

func NewCombinationLock(rdb *redis.Client, defaultTtl int) Lock {
	return &combinationLock{
		rdb:       rdb,
		redisLock: NewRedisLock(rdb, defaultTtl),
		localLock: NewLocalLock(),
	}
}

type combinationLock struct {
	rdb       *redis.Client
	redisLock Lock
	localLock Lock
}

func (c *combinationLock) TryLock(key string) (bool, error) {
	ok, err := c.redisLock.TryLock(key)
	if err != nil {
		return c.localLock.TryLock(key)
	}
	return ok, nil
}

func (c *combinationLock) UnLock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	i, err := c.rdb.Del(ctx, key).Uint64()
	if err != nil || i == 0 {
		_ = c.localLock.UnLock(key)
	}
	return nil
}
