package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	// Testing with a multi-line RESP input string
	input := strings.NewReader("*2\r\n$4\r\nECHO\r\n$5\r\nHello\r\n")

	// Initialize parser with input reader
	parser := NewRESPParser(input)

	// Parse and output the result
	req, err := parser.parse()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Parsed Result: %v\n", req)
	}

	// Test other RESP types and commands
	tests := []string{
		"*1\r\n$4\r\nPING\r\n",
		"*2\r\n$4\r\nPING\r\n$5\r\nHello\r\n",
		"*2\r\n$4\r\nECHO\r\n$5\r\nworld\r\n",
		"*2\r\n$4\r\nECHO\r\n$9\r\npineapple\r\n",
		"+OK\r\n",
		":1000\r\n",
		"$6\r\nfoobar\r\n",
		"-Error message\r\n",
	}

	for _, test := range tests {
		testInput := strings.NewReader(test)
		parser := NewRESPParser(testInput)

		req, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Parsed Result: %v\n", req)
		}
	}
}

func TestParserSetCmd(t *testing.T) {
	// Testing with a multi-line RESP input string
	input := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nFoo\r\n$3\r\nBar\r\n")

	// Initialize parser with input reader
	parser := NewRESPParser(input)

	// Parse and output the result
	req, err := parser.parse()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Parsed Result: %v %v\n", req.Cmd, req.Args)
	}

	// Test other RESP types and commands
	tests := []string{
		"*3\r\n$3\r\nSET\r\n$3\r\nFoo\r\n$3\r\nBar2\r\n",
		"*2\r\n$3\r\nGET\r\n$3\r\nFoo\r\n",
	}

	for _, test := range tests {
		testInput := strings.NewReader(test)
		parser := NewRESPParser(testInput)

		req, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Parsed Result: %v %v\n", req.Cmd, req.Args)
		}
	}
}
