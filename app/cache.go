package main

import (
	"sync"
	"time"
)

type CacheWithTTL struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value  interface{}
	expiry time.Time
	hasTTL bool
}

// NewCacheWithTTL creates a new cache instance
func NewCacheWithTTL() *CacheWithTTL {
	return &CacheWithTTL{
		data: make(map[string]cacheItem),
	}
}

// Set adds a value to the cache with an optional TTL
func (c *CacheWithTTL) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := time.Time{}
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	}

	c.data[key] = cacheItem{
		value:  value,
		expiry: expiry,
		hasTTL: ttl > 0,
	}
}

// Get retrieves a value from the cache
func (c *CacheWithTTL) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[key]
	if !found {
		return nil, false
	}

	// If the item has TTL and is expired, remove it
	if item.hasTTL && time.Now().After(item.expiry) {
		go c.Delete(key) // Non-blocking cleanup
		return nil, false
	}

	return item.value, true
}

// Delete removes an item from the cache
func (c *CacheWithTTL) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
