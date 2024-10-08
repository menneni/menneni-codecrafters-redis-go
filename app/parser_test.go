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
	result, err := parser.parse()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Parsed Result: %v\n", result)
	}

	// Test other RESP types and commands
	tests := []string{
		"+OK\r\n",
		":1000\r\n",
		"$6\r\nfoobar\r\n",
		"-Error message\r\n",
		"*2\r\nPING\r\n",
		"*2\r\n$4\r\nPING\r\n$5\r\nHello\r\n",
		"*2\r\n$4\r\nECHO\r\n$5\r\nworld\r\n",
	}

	for _, test := range tests {
		testInput := strings.NewReader(test)
		parser := NewRESPParser(testInput)

		result, err := parser.parse()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Parsed Result: %v\n", result)
		}
	}
}
