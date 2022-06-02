package dict

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkSimple(b *testing.B) {
	d := MakeSimpleDict()
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			if _, exist = d.Get(key); !exist {
				d.Put(key, key)
			}
		}
	})
}

// BenchmarkSimpleMoreRead set read:write as 9:1
func BenchmarkSimpleMoreRead(b *testing.B) {
	d := MakeSimpleDict()
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p == 0 {
				if _, exist = d.Get(key); !exist {
					d.Put(key, key)
				}
			} else {
				d.Get(key)
			}
		}
	})
}

// BenchmarkSimpleMoreRead set read:write as 9:1
func BenchmarkSimpleMoreWrite(b *testing.B) {
	d := MakeSimpleDict()
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p != 0 {
				if _, exist = d.Get(key); !exist {
					d.Put(key, key)
				}
			} else {
				d.Get(key)
			}
		}
	})
}

func BenchmarkShard(b *testing.B) {
	d := MakeShardDict(16)
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			if _, exist = d.Get(key); !exist {
				d.Put(key, key)
			}
		}
	})
}

// BenchmarkShardMoreRead set read:write as 9:1
func BenchmarkShardMoreRead(b *testing.B) {
	d := MakeShardDict(16)
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p == 0 {
				if _, exist = d.Get(key); !exist {
					d.Put(key, key)
				}
			} else {
				d.Get(key)
			}
		}
	})
}

// BenchmarkShardMoreRead set read:write as 9:1
func BenchmarkShardMoreWrite(b *testing.B) {
	d := MakeShardDict(16)
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	var exist bool
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p != 0 {
				if _, exist = d.Get(key); !exist {
					d.Put(key, key)
				}
			} else {
				d.Get(key)
			}
		}
	})
}

func BenchmarkSyncMap(b *testing.B) {
	d := sync.Map{}
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			d.LoadOrStore(key, key)
		}
	})
}

func BenchmarkSyncMapMoreRead(b *testing.B) {
	d := sync.Map{}
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p == 0 {
				d.LoadOrStore(key, key)
			} else {
				d.Load(key)
			}
		}
	})
}

func BenchmarkSyncMapMoreWrite(b *testing.B) {
	d := sync.Map{}
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			key := strconv.Itoa(r.Intn(100000000000))
			p := rand.Intn(10)
			if p != 0 {
				d.LoadOrStore(key, key)
			} else {
				d.Load(key)
			}
		}
	})
}
