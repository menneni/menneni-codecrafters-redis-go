package main

import (
	"sync"
	"time"
)

type CacheWithTtl struct {
	cache map[string]interface{}
	mutex sync.Mutex
}

func NewCacheWithTtl() *CacheWithTtl {
	return &CacheWithTtl{
		cache: make(map[string]interface{}),
	}
}

func (c *CacheWithTtl) Set(key string, val interface{}, ttl time.Duration) {
	c.mutex.Lock()
	c.cache[key] = val
	c.mutex.Unlock()

	// Schedule the key to be deleted after the TTL expires
	if ttl > 0 {
		time.AfterFunc(ttl, func() {
			c.mutex.Lock()
			delete(c.cache, key)
			c.mutex.Unlock()
		})
	}
}

func (c *CacheWithTtl) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	val, ok := c.cache[key]
	return val, ok
}
