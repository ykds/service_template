package lock

import (
	"fmt"
	"math/rand"
	"service_template/pkg/cache"
	"sync"
	"testing"
	"time"
)

func BenchmarkRedisLock(b *testing.B) {
	rdb, err := cache.NewRedis(cache.Option{
		Host: "localhost",
		Port: 6379,
	})
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	wg := sync.WaitGroup{}
	rl := NewRedisLock(rdb, 60)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Intn(50))
		ok, err := rl.TryLock(key)
		if err != nil {
			panic(err)
		}
		if ok {
			wg.Add(1)
			time.AfterFunc(time.Millisecond*time.Duration(rand.Intn(100)+1), func() {
				defer wg.Done()
				_ = rl.UnLock(key)

			})
		}
	}
	wg.Wait()
}

func BenchmarkLocalLock(b *testing.B) {
	lock := NewLocalLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := fmt.Sprintf("%d", rand.Intn(20))
			ok, err := lock.TryLock(key)
			if err != nil {
				panic(err)
			}
			if ok {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)+1))
				err := lock.UnLock(key)
				if err != nil {
					panic(err)
				}
			}
		}
	})
}

func TestLocalLock(t *testing.T) {
	lock := NewLocalLock()
	ok, _ := lock.TryLock("test")
	if ok {
		okk, _ := lock.TryLock("test1")
		if !okk {
			t.Fatal("can't lock with two diff key")
		}
		_ = lock.UnLock("test1")

		ok2, _ := lock.TryLock("test")
		if ok2 {
			t.Fatal("get the same lock")
		}
		_ = lock.UnLock("test")

		ok3, _ := lock.TryLock("test")
		if !ok3 {
			t.Fatal("can't lock after unlock")
		}
		return
	}
	t.Fatal("lock failed")
}

func TestRedisLock(t *testing.T) {
	rdb, err := cache.NewRedis(cache.Option{
		Host: "localhost",
		Port: 6379,
	})
	if err != nil {
		panic(err)
	}
	defer rdb.Close()
	lock := NewRedisLock(rdb, 10)
	ok, _ := lock.TryLock("test")
	if !ok {
		t.Fatal("申请锁失败")
	}

	okk, _ := lock.TryLock("test1")
	if !okk {
		t.Fatal("不同的key获取锁失败")
	}
	_ = lock.UnLock("test1")

	ok2, _ := lock.TryLock("test")
	if ok2 {
		t.Fatal("同key申请锁成功")
	}
	_ = lock.UnLock("test")

	ok3, _ := lock.TryLock("test")
	if !ok3 {
		t.Fatal("解锁后无法申请锁")
	}
	time.Sleep(time.Second * 15)
	ok4, _ := lock.TryLock("test")
	if ok4 {
		t.Fatal("锁没有成功续约")
	}
	_ = lock.UnLock("test")
	ok5, _ := lock.TryLock("test")
	if !ok5 {
		t.Fatal("解锁失败")
	}
}
