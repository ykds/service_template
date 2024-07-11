package ratelimiter

import (
	"fmt"
	"math/rand"
	"service_template/pkg/cache"
	"testing"
)

func BenchmarkNewRedisRateLimiter(b *testing.B) {
	rdb, err := cache.NewRedis(cache.Option{
		Host: "192.168.92.153",
		Port: 6379,
	})
	if err != nil {
		panic(err)
	}
	limiter := NewRedisRateLimiter(rdb.(*cache.Redis), 10, 60)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewLocalRateLimiter(b *testing.B) {
	limiter := NewLocalRateLimiter(10, 60)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewLocalRateLimiter2(b *testing.B) {
	limiter1 := NewLocalRateLimiter(10, 60)
	limiter2 := NewLocalRateLimiter(10, 60)
	list := []RateLimiter{limiter1, limiter2}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			limiter := list[rand.Intn(2)]
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewLocalRateLimiter3(b *testing.B) {
	limiter1 := NewLocalRateLimiter(10, 60)
	limiter2 := NewLocalRateLimiter(10, 60)
	limiter3 := NewLocalRateLimiter(10, 60)
	list := []RateLimiter{limiter1, limiter2, limiter3}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			limiter := list[rand.Intn(3)]
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewLocalRateLimiter4(b *testing.B) {
	limiter1 := NewLocalRateLimiter(10, 60)
	limiter2 := NewLocalRateLimiter(10, 60)
	limiter3 := NewLocalRateLimiter(10, 60)
	limiter4 := NewLocalRateLimiter(10, 60)
	list := []RateLimiter{limiter1, limiter2, limiter3, limiter4}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			limiter := list[rand.Intn(4)]
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkNewLocalRateLimiter5(b *testing.B) {
	limiter1 := NewLocalRateLimiter(10, 60)
	limiter2 := NewLocalRateLimiter(10, 60)
	limiter3 := NewLocalRateLimiter(10, 60)
	limiter4 := NewLocalRateLimiter(10, 60)
	limiter5 := NewLocalRateLimiter(10, 60)
	list := []RateLimiter{limiter1, limiter2, limiter3, limiter4, limiter5}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			int63 := rand.Int63()
			limiter := list[rand.Intn(5)]
			_, err := limiter.CanPass(fmt.Sprintf("%d", int63))
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
