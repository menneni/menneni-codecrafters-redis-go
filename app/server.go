package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 6379")

	for {
		// block till we receive incoming connection from client
		c, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	// keep reading data and responding with pong from same connection
	for {
		// read data
		parser := NewRESPParser(c)

		// Parse and output the result
		result, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Parsed Result: %v\n", result)
		}

		if resultStr, ok := result.(string); ok {
			_, err = c.Write([]byte(resultStr))
			if err != nil {
				fmt.Println("Error writing to connection: ", err.Error())
			}
		} else {
			fmt.Println("Error writing to connection: ", resultStr)
		}

		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
	}
}
