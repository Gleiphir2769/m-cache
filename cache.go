package m_cache

import (
	"fmt"
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
		tw := timewheel.New(ci, 3600)
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

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if c.items.Len() >= c.policy.Capacity() {
		ek := c.policy.NowEvict()
		if c.tw != nil {
			c.tw.RemoveJob(ek)
		}
		c.items.Remove(ek)
	}
	c.items.Put(k, x)
	c.policy.Promote(k)
	if c.tw != nil {
		c.tw.AddJob(k, d, func() {
			c.items.Remove(k)
			c.policy.Evict(k)
		})
	}
}

func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

func (c *cache) Add(k string, x interface{}, d time.Duration) error {
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if c.items.Len() >= c.policy.Capacity() {
		ek := c.policy.NowEvict()
		if c.tw != nil {
			c.tw.RemoveJob(ek)
		}
		c.items.Remove(ek)
	}
	if c.items.PutIfAbsent(k, x) == 0 {
		return fmt.Errorf("item %s already exists", k)
	}
	c.policy.Promote(k)
	if c.tw != nil {
		c.tw.AddJob(k, d, func() {
			c.items.Remove(k)
			c.policy.Evict(k)
		})
	}
	return nil
}

func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if c.items.Len() >= c.policy.Capacity() {
		ek := c.policy.NowEvict()
		if c.tw != nil {
			c.tw.RemoveJob(ek)
		}
		c.items.Remove(ek)
	}
	if c.items.PutIfExists(k, x) == 0 {
		return fmt.Errorf("item %s doesn't exists", k)
	}
	c.policy.Promote(k)
	if c.tw != nil {
		c.tw.AddJob(k, d, func() {
			c.items.Remove(k)
			c.policy.Evict(k)
		})
	}
	return nil
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.policy.PromoteIfExist(k)
	return c.items.Get(k)
}

// Delete an item from the m-cache. Does nothing if the key is not in the m-cache.
func (c *cache) Delete(k string) {
	ek := c.policy.NowEvict()
	if c.tw != nil {
		c.tw.RemoveJob(ek)
	}
	c.items.Remove(ek)
}
