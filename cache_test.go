package m_cache

import (
	"container/list"
	"m_cache/dict"
	"m_cache/policies"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	tc := New(DefaultExpiration, 0, dict.MakeSimpleDict(), policies.NewNon())

	a, found := tc.Get("a")
	if found || a != nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	b, found := tc.Get("b")
	if found || b != nil {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	c, found := tc.Get("c")
	if found || c != nil {
		t.Error("Getting C found value that shouldn't exist:", c)
	}

	tc.Set("a", 1, DefaultExpiration)
	tc.Set("b", "b", DefaultExpiration)
	tc.Set("c", 3.5, DefaultExpiration)

	x, found := tc.Get("a")
	if !found {
		t.Error("a was not found while getting a2")
	}
	if x == nil {
		t.Error("x for a is nil")
	} else if a2 := x.(int); a2+2 != 3 {
		t.Error("a2 (which should be 1) plus 2 does not equal 3; value:", a2)
	}

	x, found = tc.Get("b")
	if !found {
		t.Error("b was not found while getting b2")
	}
	if x == nil {
		t.Error("x for b is nil")
	} else if b2 := x.(string); b2+"B" != "bB" {
		t.Error("b2 (which should be b) plus B does not equal bB; value:", b2)
	}

	x, found = tc.Get("c")
	if !found {
		t.Error("c was not found while getting c2")
	}
	if x == nil {
		t.Error("x for c is nil")
	} else if c2 := x.(float64); c2+1.2 != 4.7 {
		t.Error("c2 (which should be 3.5) plus 1.2 does not equal 4.7; value:", c2)
	}
}

func TestCacheTimes(t *testing.T) {
	var found bool

	tc := New(500*time.Millisecond, 100*time.Millisecond, dict.MakeShardDict(16), policies.NewNon())
	tc.Set("a", 1, DefaultExpiration)
	tc.Set("b", 2, NoExpiration)
	tc.Set("c", 3, 200*time.Millisecond)
	tc.Set("d", 4, 700*time.Millisecond)

	<-time.After(250 * time.Millisecond)
	_, found = tc.Get("c")
	if found {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(300 * time.Millisecond)
	_, found = tc.Get("a")
	if found {
		t.Error("Found a when it should have been automatically deleted")
	}

	_, found = tc.Get("b")
	if !found {
		t.Error("Did not find b even though it was set to never expire")
	}

	_, found = tc.Get("d")
	if !found {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(200 * time.Millisecond)
	_, found = tc.Get("d")
	if found {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}

func TestCacheEviction(t *testing.T) {
	var found bool
	tc := New(500*time.Millisecond, 100*time.Millisecond, dict.MakeShardDict(16), policies.NewLRU(3))
	tc.Set("a", 1, NoExpiration)
	tc.Set("b", 2, NoExpiration)
	tc.Set("c", 3, NoExpiration)
	tc.Set("d", 4, NoExpiration)
	_, found = tc.Get("a")
	if found {
		t.Error("Found a when it should have been automatically deleted (LRU)")
	}
	tc.Get("b")
	tc.Set("e", 5, NoExpiration)
	_, found = tc.Get("b")
	if found {
		t.Error("Found b when it should have been automatically deleted (LRU)")
	}
}

func TestRacing(t *testing.T) {
	mockLru := NewMockLRU(100)
	tc := New(500*time.Millisecond, 100*time.Millisecond, dict.MakeShardDict(16), mockLru)
	wg := sync.WaitGroup{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			var lastKey string
			for j := 0; j < 10000; j++ {
				action := rand.Intn(3)
				k := strconv.Itoa(rand.Intn(10000000))
				switch action {
				case 0:
					tc.Set(k, 0, NoExpiration)
					lastKey = k
				case 1:
					tc.Get(lastKey)
				case 2:
					tc.Delete(lastKey)
					lastKey = ""
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	items := tc.items
	pendings := mockLru.pendingQueue

	if items.Len() != pendings.Len() {
		t.Error("Items doesn't equal to pendings")
	}

	for e := pendings.Front(); e != nil; e = e.Next() {
		k := e.Value.(string)
		if _, ok := items.Get(k); !ok {
			t.Error("Element e in pending queue doesn't exist in items")
		}
	}
}

type MockLRU struct {
	maxCap       int64
	pendingQueue *list.List
	mu           sync.Mutex
}

func NewMockLRU(maxCap int64) *MockLRU {
	return &MockLRU{
		maxCap:       maxCap,
		pendingQueue: list.New(),
	}
}

func (L *MockLRU) SetCapacity(capacity int64) {
	L.mu.Lock()
	defer L.mu.Unlock()
	L.maxCap = capacity
}

func (L *MockLRU) Capacity() int64 {
	L.mu.Lock()
	defer L.mu.Unlock()
	return L.maxCap
}

func (L *MockLRU) Promote(key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	var e *list.Element
	for e = L.pendingQueue.Front(); e != nil && e.Value.(string) != key; e = e.Next() {
	}
	if e != nil && e.Value.(string) == key {
		L.pendingQueue.Remove(e)
	}
	L.pendingQueue.PushFront(key)
}

func (L *MockLRU) PromoteIfExist(key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	var e *list.Element
	for e = L.pendingQueue.Front(); e != nil && e.Value.(string) != key; e = e.Next() {
	}
	if e != nil && e.Value.(string) == key {
		L.pendingQueue.Remove(e)
		L.pendingQueue.PushFront(key)
	}
}

func (L *MockLRU) Evict(key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	var e *list.Element
	for e = L.pendingQueue.Front(); e != nil && e.Value.(string) != key; e = e.Next() {
	}
	if e != nil && e.Value.(string) == key {
		L.pendingQueue.Remove(e)
	}
}

func (L *MockLRU) Ban(key string) {
	L.Evict(key)
}

func (L *MockLRU) NowEvict() (key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	e := L.pendingQueue.Back()
	if e != nil {
		L.pendingQueue.Remove(e)
		return e.Value.(string)
	}
	return ""
}
