// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSimple(b *testing.B) {
	d := MakeSimpleDict()

	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for pb.Next() {
			key := strconv.Itoa(r.Intn(1000000))
			d.Put(key, 0)
			d.Get(key)
		}
	})
}

// BenchmarkSimpleMoreRead set read:write as 9:1
func BenchmarkSimpleMoreRead(b *testing.B) {
	d := MakeSimpleDict()
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		lastKey := ""
		for pb.Next() {
			key := strconv.Itoa(r.Intn(1000000))
			p := rand.Intn(10)
			if p == 0 {
				d.Put(key, 0)
				lastKey = key
			} else {
				d.Get(lastKey)
			}
		}
	})
}

func BenchmarkShard(b *testing.B) {
	d := MakeShardDict(16)

	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for pb.Next() {
			key := strconv.Itoa(r.Intn(1000000))
			d.Put(key, 0)
			d.Get(key)
		}
	})
}

// BenchmarkShardMoreRead set read:write as 9:1
func BenchmarkShardMoreRead(b *testing.B) {
	d := MakeShardDict(16)
	rand.Seed(time.Now().Unix())
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		lastKey := ""
		for pb.Next() {
			key := strconv.Itoa(r.Intn(1000000))
			p := rand.Intn(10)
			if p == 0 {
				d.Put(key, 0)
				lastKey = key
			} else {
				d.Get(lastKey)
			}
		}
	})
}
