package policies

import (
	"container/list"
	"sync"
)

type LRU struct {
	maxCap       int
	pendingQueue *list.List
	mu           sync.Mutex
}

func NewLRU(maxCap int) *LRU {
	return &LRU{
		maxCap:       maxCap,
		pendingQueue: list.New(),
	}
}

func (L *LRU) SetCapacity(capacity int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	L.maxCap = capacity
}

func (L *LRU) Capacity() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return L.maxCap
}

func (L *LRU) Promote(key string) {
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

func (L *LRU) PromoteIfExist(key string) {
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

func (L *LRU) Evict(key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	var e *list.Element
	for e = L.pendingQueue.Front(); e != nil && e.Value.(string) != key; e = e.Next() {
	}
	if e != nil && e.Value.(string) == key {
		L.pendingQueue.Remove(e)
	}
}

func (L *LRU) Ban(key string) {
	L.Evict(key)
}

func (L *LRU) NowEvict() (key string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	e := L.pendingQueue.Back()
	if e != nil {
		L.pendingQueue.Remove(e)
		return e.Value.(string)
	}
	return ""
}
