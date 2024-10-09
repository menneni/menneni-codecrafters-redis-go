package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidCommand = errors.New("invalid RESP command")
)

// RESPRequest structure to hold arguments
type RESPRequest struct {
	Args interface{} // Use slice to hold multiple arguments
}

// RESPParser structure
type RESPParser struct {
	reader *bufio.Reader
}

// NewRESPParser creates a new RESP parser
func NewRESPParser(r io.Reader) *RESPParser {
	return &RESPParser{
		reader: bufio.NewReader(r),
	}
}

// readLine reads a line from the reader and trims any whitespace.
func (p *RESPParser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

// Parse parses the incoming RESP command
func (p *RESPParser) Parse() (*RESPRequest, error) {
	firstByte, err := p.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	// Determine the type of command
	switch firstByte {
	case '*': // Array
		return p.parseArray()
	case '+': // Simple String
		return p.parseSimpleString()
	case '$': // Bulk String
		return p.parseBulkString()
	case ':': // Integer
		return p.parseInteger()
	case '-': // Error
		return p.parseErrorString()
	default:
		return nil, ErrInvalidCommand
	}
}

// parseArray parses a RESP array
func (p *RESPParser) parseArray() (*RESPRequest, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		parsedResult, err := p.Parse()
		if err != nil {
			return nil, err
		}
		result[i] = parsedResult.Args
	}

	return &RESPRequest{Args: result}, nil
}

func (p *RESPParser) parseSimpleString() (*RESPRequest, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	return &RESPRequest{
		Args: line,
	}, nil
}

// parseBulkString parses a bulk string
func (p *RESPParser) parseBulkString() (*RESPRequest, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(line)
	if length == -1 {
		return &RESPRequest{Args: []string{""}}, nil // Null bulk string
	}

	// Read the actual bulk string
	bulkString := make([]byte, length+2) // +2 to account for CRLF
	if _, err := io.ReadFull(p.reader, bulkString); err != nil {
		return nil, err
	}

	return &RESPRequest{Args: string(bulkString[:length])}, nil
}

func (p *RESPParser) parseErrorString() (*RESPRequest, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	return &RESPRequest{
		Args: line,
	}, nil
}

func (p *RESPParser) parseInteger() (*RESPRequest, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	result, err := strconv.ParseInt(line, 10, 64)
	return &RESPRequest{
		Args: result,
	}, err
}
