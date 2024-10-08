package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen

func main() {
	fmt.Println("Logs from your program will appear here!")
	listener, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Failed to close listener", err)
		}
	}(listener)
	fmt.Println("Server is listening on port 6379")
	// Create an expiring map
	cache := NewCacheWithTtl()

	for {
		// block till we receive incoming connection from client
		c, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(cache, c)
	}
}

func handleConnection(cache *CacheWithTtl, c net.Conn) {
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
		}
	}(c)
	// keep reading data and responding with pong from same connection
	for {
		// read data
		parser := NewRESPParser(c)

		// Parse and output the result
		req, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
			return
		} else {
			fmt.Printf("Parsed Result: %v\n", req)
		}

		var resultStr string

		switch req.Cmd {
		case Set:
			mySlice, ok := req.Args.([]string)
			if ok && len(mySlice) == 2 {
				cache.Set(mySlice[0], mySlice[1], 0)
				fmt.Printf("Setting %s to %s\n", mySlice[0], mySlice[1])
			} else if ok && len(mySlice) == 4 {
				ttl, _ := strconv.Atoi(mySlice[3])
				if strings.EqualFold(mySlice[2], "px") {
					cache.Set(mySlice[0], mySlice[1], time.Duration(ttl)*time.Millisecond)
				}
				fmt.Printf("Setting %s to %s\n", mySlice[0], mySlice[1])
			} else {
				fmt.Println("result is not a slice of string", req.Args)
			}
			resultStr = sendOk()
		case Get:
			mySlice, ok := req.Args.([]string)
			if ok {
				resultStr = mySlice[0]
			} else {
				fmt.Println("result is not a slice of string", req.Args)
			}
			if val, ok := cache.Get(resultStr); ok {
				resultStr, err = sendBulkStringResp(val)
				fmt.Println("found val for given key", resultStr)
				if err != nil {
					fmt.Println("Error writing to connection: ", err.Error())
				}
			} else {
				resultStr = sendNullBulkString()
			}
		default:
			var ok bool
			if resultStr, ok = req.Args.(string); !ok {
				fmt.Println("Error writing to connection: ", err.Error())
			}
		}

		// write to connection
		_, err = c.Write([]byte(resultStr))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
	}
}

func sendOk() string {
	return "+OK\r\n"
}

func sendNullBulkString() string {
	return "$-1\r\n"
}

func sendSimpleStringResp(val interface{}) (string, error) {
	if valStr, ok := val.(string); ok {
		return fmt.Sprintf("+%s\r\n", valStr), nil
	}
	return "", errors.New("invalid response from server")
}

func sendBulkStringResp(val interface{}) (string, error) {
	if valStr, ok := val.(string); ok {
		return fmt.Sprintf("$%d\r\n%s\r\n", len(valStr), valStr), nil
	}
	return "", errors.New("invalid response from server")
}
