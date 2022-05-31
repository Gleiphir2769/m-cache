package m_cache

import (
	"m_cache/dict"
	"m_cache/policies"
	"m_cache/timewheel"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	*cache
	// If this is confusing, see the comment at the bottom of New()
}

type cache struct {
	defaultExpiration time.Duration
	items             dict.ConcurrentMap
	onEvicted         func(string, interface{})
	tw                *timewheel.TimeWheel
	policy            policies.EvictionPolicy
}

func New(defaultExpiration, cleanupInterval time.Duration, m dict.ConcurrentMap, p policies.EvictionPolicy) *Cache {
	return newCacheWithTimer(defaultExpiration, cleanupInterval, m, p)
}

func newCacheWithTimer(de time.Duration, ci time.Duration, m dict.ConcurrentMap, p policies.EvictionPolicy) *Cache {
	c := newCache(de, m, p)
	C := &Cache{c}
	if ci > 0 {
		tw := timewheel.New(ci, 10)
		tw.Start()
		c.tw = tw
	}
	return C
}

func newCache(de time.Duration, m dict.ConcurrentMap, p policies.EvictionPolicy) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             m,
		policy:            p,
	}
	return c
}
