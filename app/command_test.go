package main

import (
	"strings"
	"testing"
	"time"
)

// Mock connection for testing RESP parser
type MockConn struct {
	readBuffer  *strings.Reader
	writeBuffer *strings.Builder
}

// Implementing the Read method for our mock connection
func (mc *MockConn) Read(b []byte) (int, error) {
	return mc.readBuffer.Read(b)
}

// Implementing the Write method for our mock connection
func (mc *MockConn) Write(b []byte) (int, error) {
	return mc.writeBuffer.Write(b)
}

// Implementing the Close method to satisfy the net.Conn interface
func (mc *MockConn) Close() error {
	return nil
}

// Mock connection creator
func NewMockConn(input string) *MockConn {
	return &MockConn{
		readBuffer:  strings.NewReader(input),
		writeBuffer: new(strings.Builder),
	}
}

// Test case for parsing SET command
func TestRESPParser_SetCommand(t *testing.T) {
	conn := NewMockConn("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	parser := NewRESPParser(conn)

	req, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error parsing request: %v", err)
	}

	if req.Cmd != "SET" {
		t.Fatalf("Expected SET command, got %v", req.Cmd)
	}

	args, ok := req.Args.([]string)
	if !ok || len(args) != 2 || args[0] != "key" || args[1] != "value" {
		t.Fatalf("Expected [key value], got %v", req.Args)
	}
}

// Test case for handling SET and GET commands
func TestHandleRequest_SetAndGet(t *testing.T) {
	cache := NewCacheWithTTL()

	// Test SET command
	setReq := &RESPRequest{
		Cmd:  "SET",
		Args: []string{"foo", "bar"},
	}
	_, err := handleRequest(cache, setReq)
	if err != nil {
		t.Fatalf("Error handling SET request: %v", err)
	}

	// Test GET command
	getReq := &RESPRequest{
		Cmd:  "GET",
		Args: []string{"foo"},
	}
	result, err := handleRequest(cache, getReq)
	if err != nil {
		t.Fatalf("Error handling GET request: %v", err)
	}
	expected := "$3\r\nbar\r\n"
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

// Test case for handling missing key in GET command
func TestHandleRequest_GetMissingKey(t *testing.T) {
	cache := NewCacheWithTTL()

	// Test GET command for non-existing key
	getReq := &RESPRequest{
		Cmd:  "GET",
		Args: []string{"missing"},
	}
	result, err := handleRequest(cache, getReq)
	if err != nil {
		t.Fatalf("Error handling GET request: %v", err)
	}
	expected := "$-1\r\n"
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

// Test case for handling SET command with TTL
func TestHandleRequest_SetWithTTL(t *testing.T) {
	cache := NewCacheWithTTL()

	// Test SET command with TTL
	setReq := &RESPRequest{
		Cmd:  "SET",
		Args: []string{"foo", "bar", "px", "100"},
	}
	_, err := handleRequest(cache, setReq)
	if err != nil {
		t.Fatalf("Error handling SET request with TTL: %v", err)
	}

	// Test GET command immediately after setting
	getReq := &RESPRequest{
		Cmd:  "GET",
		Args: []string{"foo"},
	}
	result, err := handleRequest(cache, getReq)
	if err != nil {
		t.Fatalf("Error handling GET request: %v", err)
	}
	expected := "$3\r\nbar\r\n"
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}

	// Sleep for 150ms to wait for TTL expiry
	time.Sleep(150 * time.Millisecond)

	// Test GET command after TTL expires
	result, err = handleRequest(cache, getReq)
	if err != nil {
		t.Fatalf("Error handling GET request after TTL expiry: %v", err)
	}
	expected = "$-1\r\n"
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
