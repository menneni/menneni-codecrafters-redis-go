package main

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var myMap = make(map[string]interface{})

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
		cmd, result, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Parsed Result: %v\n", result)
		}

		var resultStr string

		switch cmd {
		case Set:
			mySlice, ok := result.([]string)
			if ok {
				myMap[mySlice[0]] = mySlice[1]
				fmt.Printf("Setting %s to %s\n", mySlice[0], mySlice[1])
			} else {
				fmt.Println("result is not a slice of string", result)
			}
			resultStr = sendOk()
		case Get:
			mySlice, ok := result.([]string)
			if ok {
				resultStr = mySlice[0]
			} else {
				fmt.Println("result is not a slice of string", result)
			}
			if val, ok := myMap[resultStr]; ok {
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
			if resultStr, ok = result.(string); !ok {
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
