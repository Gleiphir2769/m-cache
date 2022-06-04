package m_cache

import (
	"container/list"
	"m_cache/dict"
	"m_cache/policies"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkCacheWithLRU(b *testing.B) {
	c := New(time.Second, time.Second, dict.MakeShardDict(16), policies.NewLRU(10000000))
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000))
			if _, exist = c.Get(key); !exist {
				c.Set(key, key, c.defaultExpiration)
			}
		}
	})
}

func BenchmarkLRU(b *testing.B) {
	p := list.New()
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(1000000))
			p.PushBack(key)
		}
	})
}
