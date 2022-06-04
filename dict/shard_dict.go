package dict

import (
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand"
	"os"
	"sync"
	"sync/atomic"
)

type ShardDict struct {
	table    []*Shard
	count    int32
	seed     uint32
	hashAlgo func(seed uint32, k string) uint32
}

type Shard struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	} else {
		return int(n + 1)
	}
}

func MakeShardDict(shardCount int) *ShardDict {
	shardCount = computeCapacity(shardCount)
	table := make([]*Shard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &Shard{
			m: make(map[string]interface{}),
		}
	}
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint32
	if err != nil {
		os.Stderr.Write([]byte("WARNING: m-cache's MakeShardDict failed to read from the system CSPRNG (/dev/urandom or equivalent.) Your system's security may be compromised. Continuing with an insecure seed.\n"))
		seed = insecurerand.Uint32()
	} else {
		seed = uint32(rnd.Uint64())
	}
	d := &ShardDict{
		count:    0,
		table:    table,
		seed:     seed,
		hashAlgo: djb33,
	}
	return d
}

// djb2 with better shuffling. 5x faster than FNV with the hash.Hash overhead.
func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}

func (dict *ShardDict) spread(hashcode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}
	dictSize := uint32(len(dict.table))
	return (dictSize - 1) & hashcode
}

func (dict *ShardDict) getShared(index uint32) *Shard {
	return dict.table[index]
}

func (dict *ShardDict) addCount() {
	atomic.AddInt32(&dict.count, 1)
}

func (dict *ShardDict) decreaseCount() {
	atomic.AddInt32(&dict.count, -1)
}

func (dict *ShardDict) Get(key string) (val interface{}, exists bool) {
	if dict == nil {
		panic("dict is nil")
	}
	hashcode := dict.hashAlgo(dict.seed, key)
	index := dict.spread(hashcode)
	shared := dict.getShared(index)
	shared.mutex.Lock()
	defer shared.mutex.Unlock()

	val, exists = shared.m[key]
	return
}

func (dict *ShardDict) Put(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashcode := dict.hashAlgo(dict.seed, key)
	index := dict.spread(hashcode)
	shared := dict.getShared(index)
	shared.mutex.Lock()
	defer shared.mutex.Unlock()

	if _, ok := shared.m[key]; ok {
		shared.m[key] = val
		return 0
	} else {
		shared.m[key] = val
		dict.addCount()
		return 1
	}
}

func (dict *ShardDict) Len() (length int) {
	for _, s := range dict.table {
		length += len(s.m)
	}
	return length
}

// PutIfAbsent if the key has existed, the value will not be replaced.
func (dict *ShardDict) PutIfAbsent(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashcode := dict.hashAlgo(dict.seed, key)
	index := dict.spread(hashcode)
	shared := dict.getShared(index)
	shared.mutex.Lock()
	defer shared.mutex.Unlock()

	if _, ok := shared.m[key]; ok {
		return 0
	} else {
		shared.m[key] = val
		dict.addCount()
		return 1
	}
}

// PutIfExists the value will only be put when key has existed
func (dict *ShardDict) PutIfExists(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashcode := dict.hashAlgo(dict.seed, key)
	index := dict.spread(hashcode)
	shared := dict.getShared(index)
	shared.mutex.Lock()
	defer shared.mutex.Unlock()

	if _, ok := shared.m[key]; ok {
		shared.m[key] = val
		return 1
	} else {
		return 0
	}
}

func (dict *ShardDict) Remove(key string) (val interface{}, existed bool) {
	if dict == nil {
		panic("dict is nil")
	}
	hashcode := dict.hashAlgo(dict.seed, key)
	index := dict.spread(hashcode)
	shared := dict.getShared(index)
	shared.mutex.Lock()
	defer shared.mutex.Unlock()

	if v, ok := shared.m[key]; ok {
		delete(shared.m, key)
		dict.decreaseCount()
		return v, true
	} else {
		return nil, false
	}
}

func (dict *ShardDict) ForEach(recall RecallFunc) {
	if dict == nil {
		return
	}
	for _, t := range dict.table {
		t.mutex.RLock()
		func() {
			defer t.mutex.RUnlock()
			for k, v := range t.m {
				if recall(k, v) {
					return
				}
			}
		}()
	}
}
