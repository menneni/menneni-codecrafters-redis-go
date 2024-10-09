package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting Redis-like server on 127.0.0.1:6379")
	listener, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Create the cache with TTL
	cache := NewCacheWithTTL()
	startServer(listener, cache)
}

func startServer(listener net.Listener, cache *CacheWithTTL) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, cache)
	}
}

func handleConnection(conn net.Conn, cache *CacheWithTTL) {
	defer conn.Close()

	parser := NewRESPParser(conn)
	for {
		req, err := parser.Parse()
		if err != nil {
			fmt.Println("Error parsing request:", err)
			return
		}

		response, err := HandleRequest(cache, req)
		if err != nil {
			fmt.Println("Error handling request:", err)
			return
		}

		writeResponse(conn, err, response)
	}
}

func writeResponse(conn net.Conn, err error, response string) {
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}
