package main

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		cacheKey string
		cacheVal string
		ttl      time.Duration
	}{
		{
			name:     "Valid SET with PX expiry",
			input:    "*5\r\n$3\r\nSET\r\n$4\r\npear\r\n$6\r\norange\r\n$2\r\npx\r\n$3\r\n100\r\n",
			expected: "+OK\r\n",
			cacheKey: "pear",
			cacheVal: "orange",
			ttl:      100 * time.Millisecond,
		},
		{
			name:     "Valid SET without expiry",
			input:    "*3\r\n$3\r\nSET\r\n$4\r\napple\r\n$5\r\nfruit\r\n",
			expected: "+OK\r\n",
			cacheKey: "apple",
			cacheVal: "fruit",
			ttl:      0, // no expiry
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCacheWithTTL()
			reader := bufio.NewReader(strings.NewReader(tt.input))
			parser := NewRESPParser(reader)

			req, err := parser.Parse()
			if err != nil {
				t.Fatalf("Failed to parse request: %v", err)
			}

			resp, err := HandleRequest(cache, req)
			if err != nil {
				t.Fatalf("Failed to handle request: %v", err)
			}

			// Validate the response output
			if resp != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, resp)
			}

			// Validate the cache value before TTL expiry
			val, found := cache.Get(tt.cacheKey)
			if !found || val != tt.cacheVal {
				t.Fatalf("Expected cache value %s, got %v", tt.cacheVal, val)
			}

			// Check cache expiry if TTL is set
			if tt.ttl > 0 {
				time.Sleep(tt.ttl + 10*time.Millisecond) // Wait for TTL to expire
				_, found = cache.Get(tt.cacheKey)
				if found {
					t.Fatalf("Expected cache key %s to be expired", tt.cacheKey)
				}
			}
		})
	}
}
