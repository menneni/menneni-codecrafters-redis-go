package main

import (
	"testing"
	"time"
)

func TestHandleSet(t *testing.T) {
	cache := NewCacheWithTTL()

	// Test basic SET command
	req := &RESPRequest{Args: []string{"SET", "foo", "bar"}}
	resp, err := HandleRequest(cache, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp != "+OK\r\n" {
		t.Fatalf("Expected '+OK', got %v", resp)
	}

	// Check if value is set in the cache
	value, found := cache.Get("foo")
	if !found || value != "bar" {
		t.Fatalf("Expected 'bar' in cache, got %v", value)
	}

	// Test SET command with expiry (px)
	req = &RESPRequest{Args: []string{"SET", "key", "value", "px", "100"}}
	resp, err = HandleRequest(cache, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp != "+OK\r\n" {
		t.Fatalf("Expected '+OK', got %v", resp)
	}

	// Wait for expiry and check if value is removed
	time.Sleep(150 * time.Millisecond)
	_, found = cache.Get("key")
	if found {
		t.Fatalf("Expected value to expire, but found it in cache")
	}
}

func TestHandleGet(t *testing.T) {
	cache := NewCacheWithTTL()
	cache.Set("foo", "bar", 0)

	// Test GET existing key
	req := &RESPRequest{Args: []string{"GET", "foo"}}
	resp, err := HandleRequest(cache, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := "$3\r\nbar\r\n"
	if resp != expected {
		t.Fatalf("Expected '%v', got %v", expected, resp)
	}

	// Test GET non-existing key
	req = &RESPRequest{Args: []string{"GET", "unknown"}}
	resp, err = HandleRequest(cache, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected = "$-1\r\n"
	if resp != expected {
		t.Fatalf("Expected '%v', got %v", expected, resp)
	}
}

func TestHandlePing(t *testing.T) {
	// Test PING without message
	req := &RESPRequest{Args: []string{"PING"}}
	resp, err := HandleRequest(nil, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := "+PONG\r\n"
	if resp != expected {
		t.Fatalf("Expected '%v', got %v", expected, resp)
	}

	// Test PING with message
	req = &RESPRequest{Args: []string{"PING", "Hello"}}
	resp, err = HandleRequest(nil, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected = "+Hello\r\n"
	if resp != expected {
		t.Fatalf("Expected '%v', got %v", expected, resp)
	}
}

func TestHandleEcho(t *testing.T) {
	// Test ECHO with a message
	req := &RESPRequest{Args: []string{"ECHO", "Hello, world!"}}
	resp, err := HandleRequest(nil, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected := "$13\r\nHello, world!\r\n"
	if resp != expected {
		t.Fatalf("Expected '%v', got %v", expected, resp)
	}
}

func TestInvalidCommand(t *testing.T) {
	// Test unknown command
	req := &RESPRequest{Args: []string{"INVALID"}}
	_, err := HandleRequest(nil, req)
	if err == nil {
		t.Fatal("Expected error for unknown command, got none")
	}
	expectedErr := "unknown command: INVALID"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error '%v', got %v", expectedErr, err.Error())
	}
}
