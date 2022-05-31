package m_cache

import (
	"fmt"
	"time"
)

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
