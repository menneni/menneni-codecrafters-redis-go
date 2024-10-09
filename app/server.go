package main

import (
	"flag"
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

	config := InitConfig()
	// Create the cache with TTL
	cache := NewCacheWithTTL()
	startServer(listener, cache, config)
}

type Config struct {
	Dir        string
	DBFilename string
}

// InitConfig initializes the configuration by defining flags and parsing them.
func InitConfig() *Config {
	config := &Config{
		Dir:        "/tmp/redis-files", // Default value
		DBFilename: "dump.rdb",         // Default value
	}

	// Define flags and bind them to the config struct
	flag.StringVar(&config.Dir, "dir", config.Dir, "directory for Redis files")
	flag.StringVar(&config.DBFilename, "dbfilename", config.DBFilename, "filename for the database dump")

	// Parse the command line flags
	flag.Parse()

	return config
}

func startServer(listener net.Listener, cache *CacheWithTTL, config *Config) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, cache, config)
	}
}

func handleConnection(conn net.Conn, cache *CacheWithTTL, config *Config) {
	defer conn.Close()

	parser := NewRESPParser(conn)
	for {
		req, err := parser.Parse()
		if err != nil {
			fmt.Println("Error parsing request:", err)
			return
		}

		response, err := HandleRequest(cache, req, config)
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
