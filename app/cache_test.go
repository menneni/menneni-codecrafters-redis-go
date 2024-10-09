package main

import (
	"testing"
	"time"
)

// TestExpiringMapSetAndGet tests setting and getting a value from the expiring map.
func TestExpiringMapSetAndGet(t *testing.T) {
	cache := NewCacheWithTTL()

	// Set a key with a 2-second TTL
	cache.Set("foo", "bar", 2*time.Second)

	// Immediately check that the key is present
	value, found := cache.Get("foo")
	if !found {
		t.Errorf("expected to find key 'foo', but it was not found")
	}
	if value != "bar" {
		t.Errorf("expected value 'bar', but got '%v'", value)
	}

	// Wait for 3 seconds (longer than the TTL)
	time.Sleep(3 * time.Second)

	// After TTL expiration, key should not be found
	_, found = cache.Get("foo")
	if found {
		t.Errorf("expected key 'foo' to be expired, but it was still found")
	}
}
